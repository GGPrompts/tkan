# AI-Powered Autonomous Workflows with Claude

## Overview

Build GitHub Actions workflows that trigger Claude to perform autonomous development tasks on a schedule, with human approval checkpoints between phases.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                        GitHub Actions (Scheduler)                    │
│  Cron: '0 9 * * *' - Every day at 9 AM                              │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Phase 1: Documentation Review                     │
│  • Claude scans codebase for outdated docs                          │
│  • Updates README, API docs, comments                               │
│  • Creates PR with changes                                          │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Human Checkpoint: Review & Approve PR                   │
│  Notification sent via:                                              │
│  • GitHub notification                                               │
│  • Email                                                             │
│  • Slack/Discord webhook                                             │
│  • GitHub Project card with "Needs Review" status                   │
└──────────────────────────┬──────────────────────────────────────────┘
                           │ (After human approval)
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Phase 2: Code Implementation                      │
│  • Claude implements next project phase                             │
│  • Runs tests, ensures they pass                                    │
│  • Creates PR for implementation                                    │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Phase 3: Automated Code Review                          │
│  • Second Claude instance reviews first Claude's work               │
│  • Checks for security, performance, best practices                 │
│  • Comments on PR with findings                                     │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Human Checkpoint: Final Review & Merge                  │
│  • Review Claude's implementation + code review comments            │
│  • Approve or request changes                                       │
│  • Merge when ready                                                 │
└──────────────────────────┬──────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Phase 4: Update Project Board                     │
│  • Move completed tasks to "Done"                                   │
│  • Create new tasks for next phase                                  │
│  • Update sprint progress                                           │
└─────────────────────────────────────────────────────────────────────┘
```

## Implementation Options

### Option 1: Claude API (Official, Best)

Use Anthropic's official API to trigger Claude programmatically.

**Pros:**
- Official, supported method
- Full control over prompts
- Can use different models (Claude 3.5 Sonnet, Opus, etc.)
- Streaming responses

**Cons:**
- Requires API key (paid)
- Need to handle code execution yourself

### Option 2: Aider (AI Pair Programmer)

Use [Aider](https://github.com/paul-gauthier/aider) - An AI coding assistant that works in your terminal.

**Pros:**
- Designed for autonomous coding
- Handles git operations automatically
- Works with multiple files
- Free with your own API key

**Cons:**
- Still requires Claude API key
- Another tool to maintain

### Option 3: GitHub Copilot Workspace (Future)

GitHub is building autonomous coding agents, but not available yet for Actions.

## Recommended Approach: Claude API + GitHub Actions

### Prerequisites

1. **Anthropic API Key**
   ```bash
   # Get API key from https://console.anthropic.com/
   # Add to GitHub Secrets as ANTHROPIC_API_KEY
   ```

2. **GitHub Personal Access Token**
   ```bash
   # For creating PRs and issues
   # Add to GitHub Secrets as GH_TOKEN
   ```

3. **Project Configuration**
   ```bash
   # Add project details to secrets
   # PROJECT_OWNER, PROJECT_NUMBER, etc.
   ```

## Example Workflows

### Workflow 1: Daily Documentation Review

```yaml
# .github/workflows/ai-doc-review.yml
name: AI Documentation Review
on:
  schedule:
    - cron: '0 9 * * *'  # 9 AM UTC daily
  workflow_dispatch:  # Manual trigger

jobs:
  review-docs:
    runs-on: ubuntu-latest
    permissions:
      contents: write
      pull-requests: write
      issues: write

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Python
        uses: actions/setup-python@v4
        with:
          python-version: '3.11'

      - name: Install dependencies
        run: |
          pip install anthropic PyGithub

      - name: Run AI Documentation Review
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          REPO_NAME: ${{ github.repository }}
        run: |
          python .github/scripts/ai-doc-review.py

      - name: Create PR if changes found
        if: success()
        run: |
          if [[ -n $(git status -s) ]]; then
            git config user.name "Claude AI Bot"
            git config user.email "ai-bot@example.com"
            git checkout -b ai-doc-review-$(date +%Y%m%d)
            git add -A
            git commit -m "docs: AI-automated documentation review

            🤖 This PR was created by Claude AI during scheduled documentation review.

            Changes:
            - Updated outdated code examples
            - Fixed broken links
            - Improved clarity in README

            Please review before merging."

            git push origin HEAD

            gh pr create \
              --title "🤖 AI Documentation Review - $(date +%Y-%m-%d)" \
              --body "Automated documentation updates from Claude AI. Please review." \
              --label "documentation,ai-generated"
          fi
```

### Workflow 2: Multi-Phase Project Implementation

```yaml
# .github/workflows/ai-project-phases.yml
name: AI Project Implementation
on:
  workflow_dispatch:
    inputs:
      phase:
        description: 'Project phase to implement'
        required: true
        type: choice
        options:
          - phase-1-foundation
          - phase-2-features
          - phase-3-polish

jobs:
  implement-phase:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Get phase details from project
        id: phase-info
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Get phase tasks from GitHub Project
          PHASE_TASKS=$(gh project item-list 7 \
            --owner GGPrompts \
            --format json | \
            jq -r '.items[] | select(.labels[] | contains("${{ inputs.phase }}"))')

          echo "tasks=$PHASE_TASKS" >> $GITHUB_OUTPUT

      - name: Execute phase with Claude
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          PHASE_TASKS: ${{ steps.phase-info.outputs.tasks }}
        run: |
          python .github/scripts/ai-phase-executor.py \
            --phase "${{ inputs.phase }}" \
            --tasks "$PHASE_TASKS"

      - name: Run tests
        run: |
          # Run your test suite
          go test ./...

      - name: Create implementation PR
        if: success()
        run: |
          git config user.name "Claude AI Bot"
          git config user.email "ai-bot@example.com"

          BRANCH="ai-${{ inputs.phase }}-$(date +%Y%m%d)"
          git checkout -b $BRANCH
          git add -A
          git commit -m "feat: Implement ${{ inputs.phase }}

          🤖 This PR was created by Claude AI for phase implementation.

          Phase: ${{ inputs.phase }}

          Changes include:
          - Feature implementations from project tasks
          - Corresponding tests
          - Documentation updates

          ⚠️ HUMAN REVIEW REQUIRED before merging."

          git push origin HEAD

          gh pr create \
            --title "🤖 Phase Implementation: ${{ inputs.phase }}" \
            --body "$(cat .github/templates/phase-pr-template.md)" \
            --label "ai-generated,needs-review" \
            --assignee "${{ github.actor }}"

      - name: Request code review from second Claude
        if: success()
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Get PR number
          PR_NUMBER=$(gh pr view --json number -q '.number')

          # Run AI code review
          python .github/scripts/ai-code-review.py \
            --pr $PR_NUMBER

      - name: Update project board
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          # Move phase tasks to "In Review"
          python .github/scripts/update-project-status.py \
            --phase "${{ inputs.phase }}" \
            --status "In Review"

      - name: Send notification
        env:
          SLACK_WEBHOOK: ${{ secrets.SLACK_WEBHOOK }}
        run: |
          curl -X POST $SLACK_WEBHOOK \
            -H 'Content-Type: application/json' \
            -d '{
              "text": "🤖 Claude has completed ${{ inputs.phase }}!",
              "blocks": [{
                "type": "section",
                "text": {
                  "type": "mrkdwn",
                  "text": "*Phase Implementation Complete*\n\nPhase: ${{ inputs.phase }}\nPR: <link>\n\n⚠️ Human review required before proceeding."
                }
              }]
            }'
```

### Workflow 3: Scheduled Code Review & Documentation Sync

```yaml
# .github/workflows/ai-weekly-maintenance.yml
name: AI Weekly Maintenance
on:
  schedule:
    - cron: '0 0 * * 0'  # Sunday midnight

jobs:
  weekly-review:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for analysis

      - name: Analyze code changes this week
        id: changes
        run: |
          CHANGED_FILES=$(git diff --name-only HEAD~7..HEAD)
          echo "files=$CHANGED_FILES" >> $GITHUB_OUTPUT

      - name: AI-powered documentation sync
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          CHANGED_FILES: ${{ steps.changes.outputs.files }}
        run: |
          python .github/scripts/ai-doc-sync.py \
            --changed-files "$CHANGED_FILES"

      - name: Generate weekly report
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          python .github/scripts/ai-weekly-report.py \
            --output weekly-report.md

      - name: Create issue with report
        run: |
          gh issue create \
            --title "📊 AI Weekly Report - $(date +%Y-%m-%d)" \
            --body-file weekly-report.md \
            --label "report,ai-generated"
```

## Python Scripts for Claude Integration

### Script 1: AI Documentation Review (`ai-doc-review.py`)

```python
#!/usr/bin/env python3
import os
import anthropic
from pathlib import Path

def review_documentation():
    client = anthropic.Anthropic(
        api_key=os.environ.get("ANTHROPIC_API_KEY")
    )

    # Read current codebase
    codebase_files = []
    for ext in ['*.go', '*.md', '*.yaml']:
        codebase_files.extend(Path('.').rglob(ext))

    # Build context
    context = ""
    for file in codebase_files[:20]:  # Limit to avoid token limits
        try:
            content = file.read_text()
            context += f"\n\n## {file}\n```\n{content}\n```\n"
        except:
            continue

    # Prompt Claude
    message = client.messages.create(
        model="claude-3-5-sonnet-20241022",
        max_tokens=4096,
        messages=[{
            "role": "user",
            "content": f"""You are a technical documentation expert. Review this codebase and update any outdated documentation.

Current codebase:
{context}

Tasks:
1. Find discrepancies between code and documentation
2. Update README.md with current functionality
3. Fix broken examples
4. Update API documentation
5. Add missing docstrings to functions

IMPORTANT:
- Only output the COMPLETE updated file contents
- Use proper markdown formatting
- Be concise but accurate
- Format as: FILE_PATH followed by content

Output format:
FILE: path/to/file.md
CONTENT:
<complete file content here>
---
"""
        }]
    )

    # Parse Claude's response
    response_text = message.content[0].text

    # Extract file updates
    updates = parse_file_updates(response_text)

    # Apply updates
    for file_path, content in updates.items():
        Path(file_path).write_text(content)
        print(f"Updated: {file_path}")

def parse_file_updates(text):
    """Parse Claude's response into file updates"""
    updates = {}
    current_file = None
    current_content = []

    for line in text.split('\n'):
        if line.startswith('FILE:'):
            if current_file:
                updates[current_file] = '\n'.join(current_content)
            current_file = line.replace('FILE:', '').strip()
            current_content = []
        elif line.startswith('CONTENT:'):
            continue
        elif line == '---':
            if current_file:
                updates[current_file] = '\n'.join(current_content)
            current_file = None
            current_content = []
        elif current_file:
            current_content.append(line)

    return updates

if __name__ == "__main__":
    review_documentation()
```

### Script 2: AI Phase Executor (`ai-phase-executor.py`)

```python
#!/usr/bin/env python3
import os
import sys
import json
import anthropic
import argparse
from pathlib import Path

def execute_phase(phase_name, tasks):
    client = anthropic.Anthropic(
        api_key=os.environ.get("ANTHROPIC_API_KEY")
    )

    # Read current codebase
    codebase = read_codebase()

    # Build implementation prompt
    prompt = f"""You are an expert Go developer working on the tkan Terminal Kanban Board project.

PHASE: {phase_name}

TASKS TO IMPLEMENT:
{json.dumps(tasks, indent=2)}

CURRENT CODEBASE:
{codebase}

Your job:
1. Implement the features described in the tasks
2. Follow the existing code style and architecture
3. Add appropriate tests for new functionality
4. Update documentation where needed
5. Ensure all tests pass

IMPORTANT CONSTRAINTS:
- Only modify/create files necessary for this phase
- Follow Go best practices
- Use Bubbletea patterns for TUI code
- Add comments explaining complex logic
- Keep functions under 50 lines
- Write table-driven tests

Output format:
For each file you create or modify, use this format:

FILE: path/to/file.go
ACTION: create|modify
CONTENT:
<complete file content>
---

Begin implementation:
"""

    message = client.messages.create(
        model="claude-3-5-sonnet-20241022",
        max_tokens=8192,
        messages=[{"role": "user", "content": prompt}]
    )

    # Parse and apply changes
    response_text = message.content[0].text
    apply_code_changes(response_text)

def read_codebase():
    """Read relevant codebase files"""
    files = list(Path('.').glob('*.go'))

    context = ""
    for file in files:
        if file.name.endswith('_test.go'):
            continue
        content = file.read_text()
        context += f"\n\n## {file}\n```go\n{content}\n```\n"

    return context

def apply_code_changes(response):
    """Parse Claude's response and apply file changes"""
    # Similar to parse_file_updates but with ACTION field
    current_file = None
    current_action = None
    current_content = []

    for line in response.split('\n'):
        if line.startswith('FILE:'):
            if current_file and current_content:
                apply_change(current_file, current_action, '\n'.join(current_content))
            current_file = line.replace('FILE:', '').strip()
            current_content = []
        elif line.startswith('ACTION:'):
            current_action = line.replace('ACTION:', '').strip()
        elif line.startswith('CONTENT:'):
            continue
        elif line == '---':
            if current_file and current_content:
                apply_change(current_file, current_action, '\n'.join(current_content))
            current_file = None
            current_action = None
            current_content = []
        elif current_file:
            current_content.append(line)

def apply_change(file_path, action, content):
    """Apply a single file change"""
    path = Path(file_path)

    if action == 'create':
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content)
        print(f"Created: {file_path}")
    elif action == 'modify':
        path.write_text(content)
        print(f"Modified: {file_path}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--phase', required=True)
    parser.add_argument('--tasks', required=True)
    args = parser.parse_args()

    tasks = json.loads(args.tasks)
    execute_phase(args.phase, tasks)
```

### Script 3: AI Code Review (`ai-code-review.py`)

```python
#!/usr/bin/env python3
import os
import sys
import anthropic
import subprocess
import argparse

def review_pr(pr_number):
    client = anthropic.Anthropic(
        api_key=os.environ.get("ANTHROPIC_API_KEY")
    )

    # Get PR diff
    diff = subprocess.check_output(
        f"gh pr diff {pr_number}",
        shell=True,
        text=True
    )

    # Get PR description
    pr_info = subprocess.check_output(
        f"gh pr view {pr_number} --json title,body",
        shell=True,
        text=True
    )

    # Prompt for code review
    prompt = f"""You are a senior software engineer performing a code review.

PR Information:
{pr_info}

Code Changes:
```diff
{diff}
```

Perform a thorough code review focusing on:
1. Security vulnerabilities
2. Performance issues
3. Code quality and maintainability
4. Test coverage
5. Documentation completeness
6. Adherence to Go best practices
7. Potential bugs or edge cases

Provide specific, actionable feedback.

Output format:
## Summary
[Overall assessment]

## Critical Issues
- [Issue 1]
- [Issue 2]

## Suggestions
- [Suggestion 1]
- [Suggestion 2]

## Positive Observations
- [Good thing 1]
- [Good thing 2]

## Verdict
[APPROVE | REQUEST_CHANGES | COMMENT]
"""

    message = client.messages.create(
        model="claude-3-5-sonnet-20241022",
        max_tokens=4096,
        messages=[{"role": "user", "content": prompt}]
    )

    review_comment = message.content[0].text

    # Post review as comment
    subprocess.run(
        f"gh pr comment {pr_number} --body '{review_comment}'",
        shell=True,
        check=True
    )

    print(f"✅ Code review posted to PR #{pr_number}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument('--pr', type=int, required=True)
    args = parser.parse_args()

    review_pr(args.pr)
```

## Human Approval Checkpoints

### Option 1: GitHub PR Review

```yaml
# Workflow waits for PR approval before continuing
- name: Wait for approval
  uses: trstringer/manual-approval@v1
  with:
    secret: ${{ github.TOKEN }}
    approvers: ${{ github.actor }}
    minimum-approvals: 1
    issue-title: "Approve AI changes for ${{ inputs.phase }}"
```

### Option 2: GitHub Issues for Approval

```yaml
- name: Create approval issue
  run: |
    gh issue create \
      --title "🤖 Approve Phase ${{ inputs.phase }}" \
      --body "Claude has completed work. Review PR #$PR and comment 'APPROVED' to continue." \
      --label "needs-approval"

- name: Wait for approval
  run: |
    while true; do
      COMMENTS=$(gh issue view $ISSUE_NUM --json comments -q '.comments[].body')
      if echo "$COMMENTS" | grep -q "APPROVED"; then
        break
      fi
      sleep 60
    done
```

### Option 3: Slack/Discord Approval

```yaml
- name: Request approval via Slack
  run: |
    curl -X POST $SLACK_WEBHOOK -d '{
      "text": "Claude needs approval to proceed",
      "attachments": [{
        "text": "Phase complete. Reply APPROVE to continue",
        "callback_id": "ai_approval",
        "actions": [{
          "name": "approve",
          "text": "Approve",
          "type": "button",
          "value": "approved"
        }]
      }]
    }'
```

## Safety Features

### 1. Dry Run Mode

```yaml
env:
  DRY_RUN: true  # Claude suggests changes but doesn't apply them
```

### 2. Rollback Mechanism

```yaml
- name: Create rollback point
  run: |
    git tag "ai-rollback-$(date +%Y%m%d-%H%M%S)"
    git push --tags
```

### 3. Rate Limiting

```python
import time
from functools import wraps

def rate_limit(max_calls=10, period=3600):
    """Limit Claude API calls"""
    calls = []

    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            now = time.time()
            calls[:] = [c for c in calls if c > now - period]

            if len(calls) >= max_calls:
                wait = period - (now - calls[0])
                raise Exception(f"Rate limit exceeded. Wait {wait}s")

            calls.append(now)
            return func(*args, **kwargs)
        return wrapper
    return decorator
```

### 4. Cost Monitoring

```yaml
- name: Estimate API cost
  run: |
    # Claude 3.5 Sonnet: $3/million input, $15/million output tokens
    ESTIMATED_TOKENS=100000
    ESTIMATED_COST=$(echo "scale=2; $ESTIMATED_TOKENS * 15 / 1000000" | bc)
    echo "Estimated cost: \$$ESTIMATED_COST"

    if (( $(echo "$ESTIMATED_COST > 5" | bc -l) )); then
      echo "Cost too high! Aborting."
      exit 1
    fi
```

## Example: Complete End-to-End Workflow

### Scenario: Weekly Documentation Update

1. **Friday 5 PM:** GitHub Action triggers
2. **Claude scans codebase** for changes since last week
3. **Updates documentation** to match current code
4. **Creates PR** with doc updates
5. **Second Claude reviews** the documentation changes
6. **Sends Slack notification:** "📝 Weekly docs update ready for review"
7. **You review over weekend** and approve
8. **Monday 9 AM:** Auto-merge if approved
9. **Updates project board:** Move "Update docs" card to Done
10. **Creates new card:** "Review next week's docs" scheduled for next Friday

## Getting Started Checklist

- [ ] Get Anthropic API key
- [ ] Add API key to GitHub Secrets
- [ ] Create first simple workflow (doc review)
- [ ] Test with dry-run mode
- [ ] Set up approval mechanism
- [ ] Add cost monitoring
- [ ] Create rollback procedures
- [ ] Set up notifications
- [ ] Document the process for team

## Next Steps

Would you like me to create:
1. **Complete working examples** of these workflows for tkan?
2. **A GitHub Action template** repository you can copy?
3. **Monitoring dashboard** to track Claude's work?
4. **Integration with tkan project board** for automatic task management?

This architecture gives you AI-powered autonomous development with human oversight at critical points!
