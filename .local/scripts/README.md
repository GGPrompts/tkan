# Local Automation for tkan (No API Required!)

This automation system uses **local scheduling + Desktop Commander MCP** instead of the Claude API. Perfect for leveraging your **Claude Pro subscription** without additional API costs!

## 📋 Overview

### What This Does

Automatically generates **prompts** for you to review in Claude Code, based on scheduled checks:

```
┌─────────────────────────────────────┐
│  Cron/Systemd Timer (Your Computer) │
│  Runs daily at 9 AM                 │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Check Project Health               │
│  • Code changes without doc updates │
│  • Stale tasks (>2 weeks)           │
│  • Missing test coverage            │
│  • Weekly summary (Sundays)         │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Generate Prompt Files              │
│  Saved to .claude/scheduled-prompts │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  Desktop Notification + GitHub Card │
│  "3 scheduled prompts ready"        │
└──────────────┬──────────────────────┘
               │
               ▼ (When you're ready)
┌─────────────────────────────────────┐
│  You Review Prompts in Claude Code  │
│  Run them using your Pro plan       │
│  Complete tasks interactively       │
└─────────────────────────────────────┘
```

### Key Benefits

✅ **No API costs** - Uses your existing Claude Pro subscription
✅ **Runs locally** - Your machine, your control
✅ **Desktop Commander MCP** - Leverage tools you already have
✅ **Interactive** - You review and approve everything
✅ **Scheduled** - Automated checks while you sleep
✅ **Flexible** - Easy to modify schedules

## 🚀 Quick Start

### 1. Install (One Time)

```bash
cd /home/matt/projects/tkan
./.local/scripts/setup-automation.sh
```

Choose option 1 (Cron) for simple setup.

### 2. Test It

```bash
# Run the checks manually
./.local/scripts/scheduled-prompts.sh

# Review the generated prompts
./.local/scripts/review-scheduled-prompts.sh
```

### 3. Wait for Automation

The script now runs automatically at 9 AM daily. When prompts are ready:
- 🔔 Desktop notification appears
- 📋 Card added to GitHub Project #7
- 📄 Prompt files in `.claude/scheduled-prompts/`

### 4. Review & Execute

```bash
# See what's ready
./.local/scripts/review-scheduled-prompts.sh

# Or manually open prompts
cd .claude/scheduled-prompts
cat 2024-10-29-doc-sync.md

# Then execute in Claude Code using your Pro plan!
```

## 📁 Files Created

```
.local/scripts/
├── scheduled-prompts.sh        # Main automation script
├── review-scheduled-prompts.sh # Interactive prompt reviewer
├── setup-automation.sh         # One-time setup tool
└── README.md                   # This file

.claude/scheduled-prompts/      # Generated prompts
├── 2024-10-29-doc-sync.md
├── 2024-10-29-stale-items-review.md
└── 2024-11-03-weekly-summary.md  # Sundays only
```

## 🔍 What Gets Checked

### Daily (9 AM)

#### 1. Documentation Sync Check
**Trigger:** Code changes without doc updates in last 7 days

**Generated Prompt:**
```markdown
Review the codebase changes from the last week and update documentation:

1. Review changed .go files
2. Check if README.md needs updates
3. Update outdated examples
4. Fix broken links
```

#### 2. Stale Items Review
**Trigger:** Tasks "In Progress" for >2 weeks

**Generated Prompt:**
```markdown
Review and update stale project items:

The following tasks have been "In Progress" for over 2 weeks:
- Add GitHub sync
- Implement table view

For each task:
1. Assess current status
2. Move to Done/Todo/break into smaller tasks
```

#### 3. Test Coverage Check
**Trigger:** .go files without corresponding _test.go

**Generated Prompt:**
```markdown
The following files are missing test coverage:
- calendar.go
- backend_github.go

Create table-driven tests for the most critical file.
```

### Weekly (Sundays)

#### 4. Weekly Summary
**Trigger:** Every Sunday

**Generated Prompt:**
```markdown
Generate a weekly summary for tkan project:

1. Review git log from last week
2. Check GitHub Project #7 progress
3. Create summary report:
   - Accomplishments
   - Challenges/blockers
   - Next week's priorities

4. Save as: docs/weekly-summaries/2024-10-29.md
5. Create GitHub Issue for tracking
```

## ⚙️ Customization

### Change Schedule

Edit crontab:
```bash
crontab -e

# Change from 9 AM to 6 PM
0 18 * * * /home/matt/.local/scripts/scheduled-prompts.sh
```

Or with systemd:
```bash
systemctl --user edit tkan-scheduled-prompts.timer

# Change OnCalendar to your preferred time
OnCalendar=18:00
```

### Add Custom Checks

Edit `scheduled-prompts.sh`:

```bash
# Add new function
check_security_issues() {
    # Your custom check here
    create_prompt "security-review" "$(cat <<PROMPT
Check for security issues:
1. Review dependency versions
2. Check for exposed secrets
3. Validate input handling
PROMPT
)" "high"
}

# Call it in main
check_security_issues
```

### Adjust Thresholds

```bash
# In scheduled-prompts.sh

# Change "stale" from 14 days to 7 days
select(.createdAt < (now - 604800 | strftime("%Y-%m-%dT%H:%M:%SZ")))

# Change "recent code changes" from 7 days to 3 days
find . -name "*.go" -mtime -3
```

## 🔧 Advanced Usage

### Option 1: Pure Automation (No Review)

For tasks that don't need Claude, modify scripts to execute directly:

```bash
# In scheduled-prompts.sh, instead of create_prompt():

auto_fix_docs() {
    cd "$TKAN_DIR"

    # Automatically run prettier on markdown
    npx prettier --write "*.md"

    # Auto-commit if changes
    if [[ -n $(git status -s) ]]; then
        git add "*.md"
        git commit -m "docs: Auto-format markdown files"
    fi
}
```

### Option 2: Hybrid with Desktop Commander

Use Desktop Commander MCP tools for autonomous tasks:

```bash
# Example: Auto-update dependencies
auto_update_deps() {
    cd "$TKAN_DIR"

    # Use Desktop Commander to analyze dependencies
    # (This would be integrated with Claude Code MCP)

    # For now, simple version:
    go get -u ./...
    go mod tidy

    # Create PR
    git checkout -b auto-deps-$(date +%Y%m%d)
    git add go.mod go.sum
    git commit -m "deps: Auto-update dependencies"
    gh pr create --title "🤖 Auto dependency update" \
        --body "Automated dependency updates" \
        --label "dependencies"
}
```

### Option 3: Integration with Calendar View

Combine with the calendar view:

```bash
# Check tasks due this week
check_upcoming_deadlines() {
    # Get items with due dates this week
    UPCOMING=$(gh project item-list 7 --owner GGPrompts --format json | \
        jq -r --arg week_end "$(date -d '+7 days' +%Y-%m-%d)" \
        '.items[] | select(.dueDate <= $week_end) | .content.title')

    if [ -n "$UPCOMING" ]; then
        create_prompt "deadline-review" "$(cat <<PROMPT
Tasks due this week:
$UPCOMING

Review each task:
1. Is it on track?
2. Any blockers?
3. Need to adjust timeline?
PROMPT
)" "high"
    fi
}
```

## 📊 Monitoring

### View Logs

```bash
# Cron logs
tail -f ~/.local/logs/tkan-automation.log

# Systemd logs
journalctl --user -u tkan-scheduled-prompts.service -f
```

### Check Schedule

```bash
# Cron
crontab -l

# Systemd
systemctl --user list-timers
```

### Test Manually

```bash
# Run checks now
./.local/scripts/scheduled-prompts.sh

# See what would be created (dry run)
PROMPTS_DIR="/tmp/test-prompts" ./.local/scripts/scheduled-prompts.sh
```

## 🎯 Comparison: Local vs API

| Feature | Local Automation | Claude API |
|---------|-----------------|------------|
| **Cost** | Free (uses Pro plan) | ~$10-20/month |
| **Execution** | Interactive review | Fully autonomous |
| **Control** | Full control | GitHub Actions limits |
| **Speed** | When you're ready | Immediate (scheduled) |
| **Flexibility** | Very flexible | Limited to API |
| **Best For** | Prompt generation | Full automation |

## 🔐 Security

### Advantages
✅ Runs on your machine (no code sent to external APIs)
✅ You review everything before execution
✅ Full git history of changes
✅ No API keys to manage

### Best Practices
- Review generated prompts before running
- Keep scripts executable-only by you: `chmod 700`
- Don't commit sensitive prompts to git
- Use `.claude/scheduled-prompts/` in `.gitignore`

## 🎨 Example Workflow

### Monday Morning Routine

**9:00 AM** - Automated checks run
```
[Cron] Running tkan scheduled checks...
[✓] Documentation sync check
[✓] Stale items review
[✓] Test coverage check
[!] 3 prompts created
```

**9:05 AM** - You see notification
```
🔔 "tkan Scheduled Prompts"
   3 task(s) ready for Claude Code review
```

**9:30 AM** - You review prompts
```bash
$ .local/scripts/review-scheduled-prompts.sh

📋 Found 3 scheduled prompt(s)

Available prompts:
  🔴 2024-10-28-doc-sync.md
  🟡 2024-10-28-stale-items-review.md
  🟢 2024-10-28-add-tests.md
```

**9:35 AM** - Execute in Claude Code
```
# Open Claude Code
# Paste prompt: "Review code changes and update docs..."
# Claude suggests changes
# You review and approve
# Done! ✅
```

**9:45 AM** - Mark complete
```bash
# Prompts automatically marked as completed after review
# Or manually:
$ .local/scripts/review-scheduled-prompts.sh
> Option 3: Mark all as complete
✅ All prompts marked as completed
```

## 💡 Pro Tips

### 1. Use with GitHub Projects
The scripts automatically create cards for review tasks:
```bash
gh project item-create 7 --owner GGPrompts \
    --title "🤖 Review scheduled prompts (3 tasks)"
```

### 2. Prioritize Prompts
Prompts have priority levels (high/normal/low). Review high priority first:
```bash
# High priority = 🔴
# Normal = 🟡
# Low = 🟢
```

### 3. Batch Review
Save multiple prompts and review once a week:
```bash
# Change daily to weekly
0 9 * * 1  # Monday 9 AM only
```

### 4. Combine with Calendar View
Once you implement the calendar view, prompts can reference tasks by date:
```bash
create_prompt "calendar-check" "$(cat <<PROMPT
Review calendar view for overloaded days:
1. Open calendar view (press 'c')
2. Navigate to next 2 weeks
3. Identify days with >5 tasks
4. Reschedule or break down tasks
PROMPT
)"
```

## 🚀 Next Steps

1. **Test the setup**
   ```bash
   ./.local/scripts/setup-automation.sh
   ```

2. **Run first check**
   ```bash
   ./.local/scripts/scheduled-prompts.sh
   ```

3. **Review results**
   ```bash
   ./.local/scripts/review-scheduled-prompts.sh
   ```

4. **Customize schedules** to your workflow

5. **Add custom checks** for your specific needs

## 🤝 Need Help?

- Check logs: `~/.local/logs/tkan-automation.log`
- Test manually: `./.local/scripts/scheduled-prompts.sh`
- Disable: `crontab -e` and comment out lines
- Reset: Delete `.claude/scheduled-prompts/` directory

---

**Created:** 2024-10-28
**Requires:** Linux, cron or systemd, gh CLI, jq
**No API key needed!** ✅
