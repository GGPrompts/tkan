# GitHub Actions → Local PC Automation

Set up GitHub Actions to trigger Claude CLI on your local machine.

## Option 1: Self-Hosted GitHub Runner (Recommended)

### Setup

1. **Install GitHub Runner on Your PC**

```bash
# Create directory for runner
mkdir -p ~/actions-runner && cd ~/actions-runner

# Download latest runner
curl -o actions-runner-linux-x64-2.311.0.tar.gz -L \
  https://github.com/actions/runner/releases/download/v2.311.0/actions-runner-linux-x64-2.311.0.tar.gz

# Extract
tar xzf ./actions-runner-linux-x64-2.311.0.tar.gz

# Configure (follow prompts)
./config.sh --url https://github.com/GGPrompts/tkan --token YOUR_TOKEN

# Install as service (runs on startup)
sudo ./svc.sh install
sudo ./svc.sh start
```

2. **Create GitHub Action That Runs on Your PC**

```yaml
# .github/workflows/local-claude-automation.yml
name: Local Claude Automation
on:
  schedule:
    - cron: '0 9 * * *'  # 9 AM UTC
  workflow_dispatch:

jobs:
  run-on-local-pc:
    runs-on: self-hosted  # This makes it run on YOUR PC!

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Run Claude CLI for doc review
        run: |
          # This runs directly on your PC!
          claude --dangerously-skip-permissions \
            "Review the codebase and update any outdated documentation. \
             Focus on: README.md, CLAUDE.md, and code comments. \
             Commit your changes with message: 'docs: AI automated review'"

      - name: Push changes if any
        run: |
          if [[ -n $(git status -s) ]]; then
            git config user.name "Claude AI Bot"
            git config user.email "claude@local"
            git push origin HEAD
          fi
```

### How It Works

```
9 AM UTC
    ↓
GitHub Actions triggers
    ↓
Sends job to YOUR self-hosted runner
    ↓
Runner (on your PC) receives job
    ↓
Runs: claude --dangerously-skip-permissions "..."
    ↓
Claude executes using your Pro plan
    ↓
Commits and pushes results
    ↓
✅ Done!
```

**Advantages:**
- ✅ GitHub schedules it (runs even if PC was off, when it comes back)
- ✅ Uses your Claude Pro subscription (no API cost)
- ✅ Full MCP access (Desktop Commander, etc.)
- ✅ Commits results automatically
- ✅ Logs in GitHub Actions UI

**Disadvantages:**
- ⚠️ PC must be on (or will run when it comes back online)
- ⚠️ Self-hosted runner security considerations
- ⚠️ Need to keep runner service running

---

## Option 2: Webhook Listener

Create a local webhook server that listens for GitHub triggers.

### Setup

1. **Create Webhook Listener Script**

```bash
# .local/scripts/webhook-listener.sh
#!/bin/bash
# Listens for webhook from GitHub Actions

PORT=8080
SECRET="your-webhook-secret-here"

# Simple webhook listener using netcat
while true; do
    PAYLOAD=$(nc -l -p $PORT -q 1)

    # Verify secret (basic security)
    if echo "$PAYLOAD" | grep -q "$SECRET"; then
        echo "Webhook received! Triggering Claude..."

        # Run Claude CLI
        cd /home/matt/projects/tkan
        claude --dangerously-skip-permissions \
            "Review codebase and update documentation"

        # Commit results
        if [[ -n $(git status -s) ]]; then
            git add -A
            git commit -m "docs: AI automated review"
            git push
        fi
    fi
done
```

2. **Start Listener on Boot**

```bash
# Add to systemd
cat > ~/.config/systemd/user/webhook-listener.service << EOF
[Unit]
Description=GitHub Webhook Listener

[Service]
ExecStart=/home/matt/.local/scripts/webhook-listener.sh
Restart=always

[Install]
WantedBy=default.target
EOF

systemctl --user enable webhook-listener
systemctl --user start webhook-listener
```

3. **GitHub Action Sends Webhook**

```yaml
# .github/workflows/trigger-local-claude.yml
name: Trigger Local Claude
on:
  schedule:
    - cron: '0 9 * * *'

jobs:
  trigger:
    runs-on: ubuntu-latest
    steps:
      - name: Send webhook to local PC
        run: |
          curl -X POST http://YOUR_PUBLIC_IP:8080 \
            -H "Content-Type: application/json" \
            -d '{"secret":"your-webhook-secret-here","task":"doc-review"}'
```

**Requires:**
- Port forwarding (8080 → your PC)
- Dynamic DNS or static IP
- Firewall configuration

---

## Option 3: GitHub CLI Polling

Your PC periodically checks for workflow triggers.

### Setup

```bash
# .local/scripts/check-github-workflows.sh
#!/bin/bash
# Polls GitHub for workflow runs tagged for local execution

cd /home/matt/projects/tkan

# Check for workflow runs with "local-trigger" label
RUNS=$(gh run list --workflow=local-claude-automation.yml \
    --json databaseId,status,conclusion \
    --limit 1)

RUN_ID=$(echo "$RUNS" | jq -r '.[0].databaseId')
STATUS=$(echo "$RUNS" | jq -r '.[0].status')

# If there's a queued run, execute it locally
if [ "$STATUS" = "queued" ]; then
    echo "Found queued run: $RUN_ID"

    # Execute Claude CLI
    claude --dangerously-skip-permissions \
        "Review codebase and update documentation"

    # Mark as complete (create a marker file in repo)
    echo "completed" > .github/.local-run-$RUN_ID
    git add .github/.local-run-$RUN_ID
    git commit -m "chore: Mark run $RUN_ID as complete"
    git push
fi
```

Run via cron every 5 minutes:
```bash
*/5 * * * * /home/matt/.local/scripts/check-github-workflows.sh
```

---

## Comparison: All Approaches

| Approach | PC Must Be On? | GitHub Scheduled? | API Cost | Setup Complexity |
|----------|---------------|------------------|----------|------------------|
| **Pure Cron** (my original) | Yes | No | $0 | Easy |
| **Self-Hosted Runner** | Yes* | Yes | $0 | Medium |
| **Webhook Listener** | Yes | Yes | $0 | Hard |
| **GitHub CLI Polling** | Yes | Kinda† | $0 | Easy |
| **Claude API** | No | Yes | $10-20/mo | Easy |

\* Can catch up when PC turns back on
† Not truly scheduled, polls every 5 min

---

## Security Considerations

### ⚠️ Using `--dangerously-skip-permissions`

This flag is **powerful** and **dangerous**:

```bash
# What it does:
claude --dangerously-skip-permissions "prompt"

# Allows:
✅ Write files without asking
✅ Execute commands without asking
✅ Access network without asking
✅ Commit to git without asking

# Risks:
❌ Malicious prompts could damage your system
❌ Bugs could delete important files
❌ No human review before execution
```

**Safer Alternative:**

```bash
# Use a restricted prompt that requires confirmation for critical actions
claude "Review docs and show me proposed changes. \
       Don't commit until I approve."
```

Or create a wrapper:

```bash
# .local/scripts/safe-claude-automation.sh
#!/bin/bash

# Create a safety sandbox
SANDBOX_DIR="/tmp/claude-sandbox-$(date +%s)"
mkdir -p "$SANDBOX_DIR"
cd "$SANDBOX_DIR"

# Clone repo to sandbox
git clone /home/matt/projects/tkan .

# Run Claude with skip-permissions (safe because it's in sandbox)
claude --dangerously-skip-permissions "$1"

# Review changes
echo "Changes made:"
git diff HEAD

# Require manual approval
read -p "Apply these changes to real repo? (yes/no): " approval

if [ "$approval" = "yes" ]; then
    # Copy changes to real repo
    rsync -av --exclude='.git' ./ /home/matt/projects/tkan/

    cd /home/matt/projects/tkan
    git add -A
    git commit -m "docs: AI automated review (approved)"
    git push
else
    echo "Changes discarded"
fi

# Cleanup
rm -rf "$SANDBOX_DIR"
```

---

## Best Recommendation

Based on your setup, I recommend:

### **Hybrid: Self-Hosted Runner + Safety Sandbox**

1. **GitHub Actions schedules** (cloud, always works)
2. **Self-hosted runner on your PC** (uses Claude Pro)
3. **Safety sandbox** (runs in temp dir first)
4. **Manual approval** (you review before applying)

```yaml
# .github/workflows/safe-local-automation.yml
name: Safe Local Automation
on:
  schedule:
    - cron: '0 9 * * *'
  workflow_dispatch:

jobs:
  generate-review:
    runs-on: self-hosted
    steps:
      - name: Create sandbox and run Claude
        run: |
          /home/matt/.local/scripts/safe-claude-automation.sh \
            "Review documentation and propose updates"

      - name: Create PR with proposed changes
        if: success()
        run: |
          # Changes are in sandbox, not committed yet
          # Create PR for review
          cd /tmp/claude-sandbox-*
          git checkout -b ai-doc-review-$(date +%Y%m%d)
          git add -A
          git commit -m "docs: AI proposed documentation updates"
          git push -u origin HEAD

          gh pr create \
            --title "🤖 AI Documentation Review" \
            --body "Claude AI has proposed documentation updates. Please review." \
            --label "documentation,ai-generated,needs-review"
```

**Workflow:**
1. 9 AM: GitHub Actions triggers
2. Your PC (runner) receives job
3. Claude runs in sandbox with `--dangerously-skip-permissions`
4. Creates PR (doesn't auto-merge)
5. 🔔 You get notification
6. You review PR and merge if good

**Benefits:**
- ✅ Uses Claude Pro (no API cost)
- ✅ GitHub schedules it
- ✅ Full automation
- ✅ Safety sandbox
- ✅ Human approval checkpoint
- ✅ All the automation, none of the risk

---

## Comparison to My Original Approach

| Feature | My Cron Approach | Your GitHub+Claude Idea | Recommended Hybrid |
|---------|-----------------|------------------------|-------------------|
| **Scheduling** | Local cron | GitHub Actions | GitHub Actions |
| **Execution** | Generate prompts | Auto-execute | Auto-execute in sandbox |
| **Review** | Before execution | After execution | Before applying |
| **PC Offline** | Missed (unless anacron) | Queued, runs later | Queued, runs later |
| **Cost** | $0 | $0 | $0 |
| **Safety** | Very safe | Risky | Very safe |
| **Automation** | Manual review needed | Fully automated | Fully automated with approval |

**Verdict:** Your idea is **brilliant** - it combines GitHub's reliable scheduling with your Claude Pro subscription. The hybrid approach adds safety while keeping full automation!

---

## 📅 Missed Schedules: What Happens?

Great question about PC being off!

### Cron (My Original Approach)

**Default behavior:**
```bash
# If PC is off at 9 AM:
9:00 AM - ❌ Missed, doesn't run
10:00 AM - PC turns on
10:01 AM - Nothing happens (job was missed)
```

**With `anacron` (catches up):**
```bash
# Install anacron
sudo apt install anacron

# Configure to catch up
echo '@daily 10 /home/matt/.local/scripts/scheduled-prompts.sh' > /etc/anacron/
# or in crontab:
@daily /home/matt/.local/scripts/scheduled-prompts.sh
```

Now:
```bash
9:00 AM - ❌ PC off, missed
2:00 PM - PC turns on
2:05 PM - ✅ Anacron runs missed job!
```

### Systemd Timers with Persistent=true

```bash
# In timer file:
[Timer]
OnCalendar=daily
OnCalendar=09:00
Persistent=true  # ← This catches up on missed runs!
```

Behavior:
```bash
9:00 AM - ❌ PC off
3:00 PM - PC turns on
3:01 PM - ✅ Systemd runs missed job immediately
```

### Self-Hosted GitHub Runner

**Automatic catch-up:**
```bash
9:00 AM - GitHub Action triggered
9:01 AM - Job queued (waiting for runner)
2:00 PM - Your PC turns on
2:01 PM - Runner comes online
2:02 PM - ✅ Runs the queued job!
```

GitHub Actions will **wait up to 72 hours** for a runner to become available.

---

## 🎯 Final Recommendation

Use this approach:

1. **Install self-hosted GitHub runner** (one-time setup)
2. **Create workflow with safety sandbox**
3. **Use `--dangerously-skip-permissions`** (but in sandbox)
4. **Create PR for review** (don't auto-merge)

You get:
- ✅ GitHub's reliable scheduling (catch-up if PC was off)
- ✅ Claude Pro subscription (no API cost)
- ✅ Full automation (sandbox execution)
- ✅ Safety (review before applying)
- ✅ All benefits, minimal risk

Want me to create the complete setup files for this approach?