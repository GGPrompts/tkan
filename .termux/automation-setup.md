# Termux Automation Server Setup

Your phone is the perfect always-on automation server for Claude-powered tasks!

## Why Termux is Perfect

| Feature | Termux on Phone | PC | Cloud API |
|---------|----------------|-----|-----------|
| **Always On** | ✅ (charging overnight) | ❌ (turns off) | ✅ |
| **Claude Pro** | ✅ (uses your plan) | ✅ (uses your plan) | ❌ (separate API) |
| **Cost** | $0 | $0 | $10-20/mo |
| **Cron Jobs** | ✅ | ✅ | N/A |
| **Git Access** | ✅ | ✅ | ✅ |
| **MCP Tools** | Partial | Full | N/A |
| **Background** | ✅ (with Termux:Boot) | ✅ | ✅ |
| **Internet** | ✅ (cellular/wifi) | Depends | ✅ |

**Winner:** Termux! 🏆

---

## 🚀 Quick Setup

### 1. Install Required Packages

```bash
# Update packages
pkg upgrade

# Install essentials
pkg install git gh jq cronie termux-services

# Enable cron service
sv-enable crond

# Install Claude Code (if not already)
# Follow: https://github.com/anthropics/claude-code
```

### 2. Clone Your Repository

```bash
# Create workspace
mkdir -p ~/projects
cd ~/projects

# Clone tkan
git clone https://github.com/GGPrompts/tkan.git
cd tkan

# Configure git
git config user.name "Claude AI Bot"
git config user.email "claude@termux"

# Set up GitHub CLI
gh auth login
gh auth refresh -h github.com -s project
```

### 3. Create Automation Scripts

```bash
# Create scripts directory
mkdir -p ~/.termux/automation

# Create main automation script
cat > ~/.termux/automation/scheduled-tasks.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash
# Termux automation runner

set -euo pipefail

TKAN_DIR="$HOME/projects/tkan"
LOG_DIR="$HOME/.termux/logs"
DATE=$(date +%Y-%m-%d)

mkdir -p "$LOG_DIR"

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_DIR/automation.log"
}

log "Starting scheduled tasks..."

# Update repository
cd "$TKAN_DIR"
git pull origin master

# Run documentation review
doc_review() {
    log "Running documentation review..."

    # Check for code changes without doc updates
    CHANGED_GO=$(git log --since="7 days ago" --name-only --pretty=format: | \
                 grep "\.go$" | sort -u | wc -l)
    CHANGED_MD=$(git log --since="7 days ago" --name-only --pretty=format: | \
                 grep "\.md$" | sort -u | wc -l)

    if [ "$CHANGED_GO" -gt 0 ] && [ "$CHANGED_MD" -eq 0 ]; then
        log "Code changed but docs didn't - running AI review"

        # Run Claude CLI with auto-skip
        claude --dangerously-skip-permissions \
            "Review the codebase changes from the last week. \
             Update README.md, CLAUDE.md, or API documentation as needed. \
             Focus on: \
             1. New features that need documenting \
             2. Changed function signatures \
             3. Outdated examples \
             4. Broken links \
             Commit your changes with message: 'docs: AI automated weekly review'" \
            >> "$LOG_DIR/doc-review-$DATE.log" 2>&1

        # Push if changes were made
        if [[ -n $(git status -s) ]]; then
            git push origin master
            log "Documentation updated and pushed"

            # Send notification
            termux-notification \
                --title "tkan Automation" \
                --content "Documentation updated successfully" \
                --button1 "View" \
                --button1-action "termux-open-url https://github.com/GGPrompts/tkan"
        else
            log "No documentation changes needed"
        fi
    else
        log "Documentation appears up to date"
    fi
}

# Run project health check
health_check() {
    log "Running project health check..."

    # Get project stats
    cd "$TKAN_DIR"

    TODO_COUNT=$(gh project item-list 7 --owner GGPrompts --format json | \
                 jq '[.items[] | select(.status == "Todo")] | length' || echo "0")
    IN_PROGRESS=$(gh project item-list 7 --owner GGPrompts --format json | \
                  jq '[.items[] | select(.status == "In Progress")] | length' || echo "0")

    log "Project stats: $TODO_COUNT todo, $IN_PROGRESS in progress"

    # Check for stale items (>14 days in progress)
    STALE=$(gh project item-list 7 --owner GGPrompts --format json | \
            jq -r '.items[] |
                select(.status == "In Progress") |
                select(.createdAt < (now - 1209600 | strftime("%Y-%m-%dT%H:%M:%SZ"))) |
                .content.title' || echo "")

    if [ -n "$STALE" ]; then
        log "Found stale items: $STALE"

        # Create review task
        gh project item-create 7 --owner GGPrompts \
            --title "🤖 Review stale items (Termux automation)" \
            --body "The following items have been in progress >14 days:
$STALE

Review and either:
- Move to Done if completed
- Move to Todo if blocked
- Break into smaller tasks
- Add notes about blockers"

        termux-notification \
            --title "tkan Health Check" \
            --content "Found stale items needing review"
    fi
}

# Run weekly summary (Sundays only)
weekly_summary() {
    if [ "$(date +%u)" -eq 7 ]; then
        log "Running weekly summary..."

        cd "$TKAN_DIR"

        # Count commits
        COMMITS=$(git log --since="1 week ago" --oneline | wc -l)

        # Run Claude to generate summary
        claude --dangerously-skip-permissions \
            "Generate a weekly summary for the tkan project:

            1. Review git log from last week (found $COMMITS commits)
            2. Check GitHub Project #7 progress
            3. Create a summary in docs/weekly-summaries/$(date +%Y-%m-%d).md covering:
               - Key accomplishments
               - Challenges or blockers
               - Next week's priorities
               - Technical debt identified
            4. Commit with message: 'docs: Weekly summary for $(date +%Y-%m-%d)'" \
            >> "$LOG_DIR/weekly-summary-$DATE.log" 2>&1

        if [[ -n $(git status -s) ]]; then
            git push origin master
            log "Weekly summary created and pushed"

            termux-notification \
                --title "tkan Weekly Summary" \
                --content "Summary for $(date +%Y-%m-%d) created" \
                --priority high
        fi
    fi
}

# Run all checks
doc_review
health_check
weekly_summary

log "Scheduled tasks completed"

# Cleanup old logs (keep last 30 days)
find "$LOG_DIR" -name "*.log" -mtime +30 -delete
EOF

chmod +x ~/.termux/automation/scheduled-tasks.sh
```

### 4. Schedule with Cron

```bash
# Edit crontab
crontab -e

# Add this line (runs daily at 9 AM)
0 9 * * * /data/data/com.termux/files/home/.termux/automation/scheduled-tasks.sh

# Or for testing, every hour:
# 0 * * * * /data/data/com.termux/files/home/.termux/automation/scheduled-tasks.sh
```

### 5. Enable Termux:Boot (Auto-start on Phone Boot)

```bash
# Install Termux:Boot from F-Droid or Play Store

# Create boot script
mkdir -p ~/.termux/boot
cat > ~/.termux/boot/start-automation.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash
# Auto-start services on phone boot

# Start cron
sv-enable crond

# Optional: Run immediate health check
# ~/.termux/automation/scheduled-tasks.sh &
EOF

chmod +x ~/.termux/boot/start-automation.sh
```

---

## 🎯 Usage Patterns

### Pattern 1: Scheduled Automation (Recommended)

Your phone runs tasks automatically while charging overnight:

```
11:00 PM - You plug in phone to charge
            ↓
 9:00 AM - Cron triggers automation
            ↓
           Claude reviews codebase
            ↓
           Updates documentation
            ↓
           Commits and pushes
            ↓
           🔔 Notification: "Docs updated!"
            ↓
 9:30 AM - You wake up, review notification
            ↓
           Check GitHub for changes
            ↓
           ✅ Done!
```

### Pattern 2: On-Demand Execution

Run manually when you want:

```bash
# SSH into your phone from PC
ssh u0_a123@192.168.1.100 -p 8022

# Or use Termux directly
cd ~/projects/tkan
~/.termux/automation/scheduled-tasks.sh
```

### Pattern 3: GitHub Webhook Trigger

Phone listens for GitHub webhooks:

```bash
# Install webhook listener
pkg install webhook

# Create webhook config
cat > ~/.termux/automation/webhook.json << 'EOF'
[
  {
    "id": "tkan-automation",
    "execute-command": "/data/data/com.termux/files/home/.termux/automation/scheduled-tasks.sh",
    "command-working-directory": "/data/data/com.termux/files/home/projects/tkan",
    "pass-arguments-to-command": [],
    "trigger-rule": {
      "match": {
        "type": "value",
        "value": "tkan-automation-trigger",
        "parameter": {
          "source": "payload",
          "name": "action"
        }
      }
    }
  }
]
EOF

# Start webhook listener (in background)
nohup webhook -hooks ~/.termux/automation/webhook.json -port 9000 &

# GitHub Action sends webhook to your phone:
# curl -X POST http://your-phone-ip:9000/hooks/tkan-automation \
#   -d '{"action":"tkan-automation-trigger"}'
```

---

## 🔋 Battery & Performance Optimization

### Keep Phone Charged

```bash
# Check battery status
termux-battery-status

# Automation best practices:
# - Run while charging
# - Schedule for overnight (9 AM = after overnight charge)
# - Disable power-intensive operations
```

### Prevent Termux from Being Killed

```bash
# Acquire wakelock (keeps Termux alive)
termux-wake-lock

# Release when done
termux-wake-unlock

# Add to automation script:
cat >> ~/.termux/automation/scheduled-tasks.sh << 'EOF'
# Acquire wakelock at start
termux-wake-lock

# Release at end (add to end of script)
termux-wake-unlock
EOF
```

### Battery-Friendly Schedule

```bash
# Instead of hourly, run once daily
# Or only when charging:

crontab -e
# Add condition:
0 9 * * * [ $(termux-battery-status | jq -r '.status') = "CHARGING" ] && ~/.termux/automation/scheduled-tasks.sh
```

---

## 📱 Termux-Specific Features

### Send Notifications

```bash
# Success notification
termux-notification \
    --title "tkan Automation" \
    --content "Documentation updated successfully" \
    --button1 "View Diff" \
    --button1-action "termux-open-url https://github.com/GGPrompts/tkan/commits"

# Error notification
termux-notification \
    --title "tkan Automation Failed" \
    --content "Check logs for details" \
    --priority high
```

### Access from PC

```bash
# On Termux (phone):
# Install SSH
pkg install openssh

# Start SSH server
sshd

# Get IP address
ifconfig | grep inet

# On PC:
ssh u0_a123@192.168.1.100 -p 8022

# Now you can:
# - View logs
# - Trigger automation manually
# - Check status
```

### Share Files with PC

```bash
# On Termux:
pkg install rsync

# Sync logs to PC
rsync -av ~/.termux/logs/ user@pc-ip:~/tkan-logs/

# Or use git to share
cd ~/projects/tkan
git add docs/
git commit -m "docs: Update from Termux automation"
git push
```

---

## 🛡️ Security Considerations

### Termux Limitations (Good for Safety!)

Claude on Termux **cannot** access:
- ❌ Chrome DevTools (you mentioned)
- ❌ Desktop applications
- ❌ System-level operations

This is actually **safer** for automation because:
- ✅ Limited damage potential
- ✅ Can't mess with your desktop
- ✅ Isolated to Termux environment

### Recommended Safety Measures

```bash
# 1. Use a separate git worktree
cd ~/projects/tkan
git worktree add ~/tkan-automation master

# Claude works in automation worktree
cd ~/tkan-automation

# 2. Create PR instead of direct push
claude --dangerously-skip-permissions "..."
git checkout -b ai-review-$(date +%Y%m%d)
git push origin HEAD
gh pr create --title "AI Review" --body "From Termux automation"

# 3. Dry-run mode
DRY_RUN=true ~/.termux/automation/scheduled-tasks.sh
```

---

## 🔍 Monitoring & Debugging

### View Logs

```bash
# Real-time log
tail -f ~/.termux/logs/automation.log

# Today's doc review
cat ~/.termux/logs/doc-review-$(date +%Y-%m-%d).log

# Recent activity
ls -lt ~/.termux/logs/ | head -10
```

### Check Cron Status

```bash
# List cron jobs
crontab -l

# Check if cron is running
sv status crond

# View cron logs
cat /data/data/com.termux/files/usr/var/log/cron.log
```

### Test Automation

```bash
# Run manually
~/.termux/automation/scheduled-tasks.sh

# Run with verbose output
bash -x ~/.termux/automation/scheduled-tasks.sh

# Test specific function
cd ~/projects/tkan
source ~/.termux/automation/scheduled-tasks.sh
doc_review  # Run just this function
```

---

## 🎨 Advanced Patterns

### Multi-Project Management

```bash
# Manage multiple projects from phone
cat > ~/.termux/automation/multi-project.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash

PROJECTS=(
    "tkan:GGPrompts:7"
    "another-project:YourOrg:12"
)

for project in "${PROJECTS[@]}"; do
    IFS=':' read -r name owner number <<< "$project"

    cd ~/projects/$name
    git pull

    # Run automation for this project
    claude --dangerously-skip-permissions \
        "Review $name documentation and update as needed"

    [[ -n $(git status -s) ]] && git push
done
EOF
```

### Conditional Automation

```bash
# Only run if connected to WiFi (save cellular data)
NETWORK_TYPE=$(termux-wifi-connectioninfo | jq -r '.supplicant_state')

if [ "$NETWORK_TYPE" = "COMPLETED" ]; then
    # On WiFi, run full automation
    ~/.termux/automation/scheduled-tasks.sh
else
    # On cellular, skip or run lightweight tasks only
    echo "Skipping - not on WiFi"
fi
```

### Integration with Desktop Commander

```bash
# Use MCP tools available in Termux
claude --dangerously-skip-permissions \
    "Use Desktop Commander to:
     1. List recent files in ~/projects/tkan
     2. Search for TODO comments
     3. Update project documentation

     Note: Some tools like Chrome DevTools won't work on Termux,
     but file operations, search, and git work fine!"
```

---

## 🆚 Comparison: All Automation Options

| Approach | Always On | Claude Pro | Setup | Best For |
|----------|----------|------------|-------|----------|
| **Termux (Phone)** | ✅ 24/7 | ✅ | Easy | **Best overall** |
| **PC Cron** | ⚠️ When on | ✅ | Easy | PC-only tasks |
| **Self-Hosted Runner** | ⚠️ When on | ✅ | Medium | GitHub integration |
| **Cloud API** | ✅ 24/7 | ❌ | Easy | Production critical |

**Winner: Termux!** 🏆

Advantages:
- Always on (phone in pocket)
- Uses Claude Pro (no API cost)
- Easy setup
- Background execution
- Push notifications
- Accessible from anywhere (SSH)

---

## 📊 Real-World Example Schedule

Here's how I'd run automation on Termux:

```bash
# Crontab
crontab -e

# Daily at 9 AM (after overnight charge)
0 9 * * * ~/.termux/automation/scheduled-tasks.sh

# Weekly summary (Sunday 10 AM)
0 10 * * 0 ~/.termux/automation/weekly-summary.sh

# Quick health check (every 6 hours, only if charging)
0 */6 * * * [ $(termux-battery-status | jq -r '.status') = "CHARGING" ] && ~/.termux/automation/health-check.sh

# Backup logs (monthly)
0 0 1 * * ~/.termux/automation/backup-logs.sh
```

**Battery impact:** Minimal (runs while charging, quick execution)

---

## 🚀 Quick Start Checklist

- [ ] Install required packages: `pkg install git gh jq cronie termux-services`
- [ ] Enable cron: `sv-enable crond`
- [ ] Clone tkan repo: `git clone ...`
- [ ] Create automation scripts in `~/.termux/automation/`
- [ ] Set up crontab: `crontab -e`
- [ ] Install Termux:Boot (optional)
- [ ] Test: `~/.termux/automation/scheduled-tasks.sh`
- [ ] Check notifications work: `termux-notification ...`
- [ ] Monitor logs: `tail -f ~/.termux/logs/automation.log`

---

## 💡 Pro Tips

1. **Phone as CI/CD server** - Your phone is a legitimate automation server!
2. **Always-on advantage** - Unlike PC, phone stays on 24/7
3. **Use notifications** - Get alerts when automation completes
4. **SSH access** - Manage from PC when needed
5. **Battery-friendly** - Schedule during charging hours
6. **Git worktrees** - Separate automation from main work
7. **Multiple projects** - Manage all repos from phone
8. **Backup to cloud** - Push logs to GitHub/Dropbox

---

## 🎉 Summary

Your Termux setup is **perfect** for automation because:

✅ Always on (phone is always with you)
✅ Uses Claude Pro (no API costs)
✅ Full git access (push/pull)
✅ Background execution (cron)
✅ Push notifications (termux-notification)
✅ Remote access (SSH)
✅ Battery efficient (runs while charging)
✅ Isolated environment (safer than PC)

**This is genuinely better than most cloud CI/CD setups!**

Your $200/month Claude Pro subscription just became a **24/7 automation server** for $0 additional cost! 🚀
