# AI-Powered Autonomous Workflows

Use GitHub Actions + Claude API to build autonomous development workflows with human checkpoints.

## Quick Start

### 1. Get API Key

```bash
# Get API key from https://console.anthropic.com/
# Add to GitHub repository secrets as ANTHROPIC_API_KEY
```

### 2. Create First Workflow

```bash
mkdir -p .github/workflows
mkdir -p .github/scripts

# Copy example workflows from ARCHITECTURE.md
```

### 3. Test Dry Run

```yaml
env:
  DRY_RUN: true  # Claude suggests but doesn't commit
```

## Use Cases

### 🤖 Autonomous Tasks
- Daily documentation reviews
- Weekly code quality checks
- Dependency updates
- Test coverage reports
- Performance benchmarking

### 👨‍💻 Multi-Phase Development
- Phase 1: Claude implements feature
- **Human checkpoint:** Review implementation
- Phase 2: Claude writes tests
- **Human checkpoint:** Review tests
- Phase 3: Claude updates docs
- **Human checkpoint:** Final review & merge

### 🔄 Scheduled Maintenance
- Weekly sprint summaries
- Monthly security audits
- Quarterly refactoring reviews
- Automated changelog generation

## Safety Features

✅ Dry-run mode for testing
✅ Human approval checkpoints
✅ Cost monitoring
✅ Rollback mechanisms
✅ Rate limiting

## Example Workflows

### Daily Doc Review
```yaml
on:
  schedule:
    - cron: '0 9 * * *'  # 9 AM daily
```

### Phase-Based Implementation
```yaml
on:
  workflow_dispatch:
    inputs:
      phase:
        type: choice
        options:
          - phase-1
          - phase-2
```

### Code Review Bot
```yaml
on:
  pull_request:
    types: [opened, synchronize]
```

## Architecture

```
GitHub Actions (Scheduler)
    ↓
Claude API (Implementation)
    ↓
Create PR
    ↓
Second Claude (Review)
    ↓
Human Checkpoint (Approve/Reject)
    ↓
Merge & Update Project Board
```

## Files

- `ARCHITECTURE.md` - Complete architecture guide
- `README.md` - This file

## Cost Estimation

Claude 3.5 Sonnet pricing:
- **Input:** $3 per million tokens
- **Output:** $15 per million tokens

Example costs:
- Doc review: ~$0.05 per run
- Phase implementation: ~$0.50 per phase
- Code review: ~$0.10 per PR

## Getting Started

1. Read `ARCHITECTURE.md` for full details
2. Set up API keys in GitHub Secrets
3. Create first workflow (start with doc review)
4. Test in dry-run mode
5. Add human approval checkpoints
6. Monitor costs and usage
7. Expand to more workflows

## See Also

- `/github-projects` skill for project integration
- Claude API docs: https://docs.anthropic.com/
- GitHub Actions docs: https://docs.github.com/actions
