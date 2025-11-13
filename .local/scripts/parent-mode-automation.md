# Parent Mode: Multi-Context Development Automation

For developers who code in 5-minute bursts between diaper changes and snack time.

## 🎯 The Problem

**Your current workflow:**
```
8:00 PM - Kids asleep! Time to code!
8:05 PM - Open 5 terminals
8:10 PM - Create 5 git worktrees
8:15 PM - Copy same prompt to each terminal
8:20 PM - Start watching them work
8:25 PM - "DAAAADDY!" 😭
8:30 PM - Back to coding... wait, which terminal was doing what?
8:35 PM - Try to remember context
8:40 PM - "DAAAADDY I NEED WATER!" 💧
8:45 PM - Give up, close 3 terminals
9:00 PM - Collapse in exhaustion
```

**Time spent coding:** 20 minutes
**Time spent context switching:** 40 minutes

---

## ✨ The Solution: Parent Mode

Let your phone do the parallel work while you handle the kids!

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ 8:00 PM - You trigger "parent-mode" command                 │
│ Takes 30 seconds of your time                               │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│ Your Phone (Termux) Creates 5 Parallel Contexts             │
│ • Feature A: calendar-view branch                           │
│ • Feature B: table-view branch                              │
│ • Feature C: export-feature branch                          │
│ • Refactor: cleanup-backend branch                          │
│ • Docs: update-readme branch                                │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│ 8:00 PM - 8:40 PM (While You're With Kids)                  │
│ Claude works on all 5 branches in parallel                  │
│ • Implements features                                       │
│ • Writes tests                                              │
│ • Updates docs                                              │
│ • Creates commits                                           │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│ 8:40 PM - Kids Finally Asleep                               │
│ 🔔 "Parent Mode Complete: 5 features ready for review"      │
│ • Feature A: ✅ Implemented + tested                        │
│ • Feature B: ✅ Implemented + tested                        │
│ • Feature C: ⚠️  Partial (needs clarification)              │
│ • Refactor: ✅ Complete                                     │
│ • Docs: ✅ Updated                                          │
└──────────────────┬──────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────┐
│ 8:45 PM - You Review & Merge                                │
│ • Quick review of each branch (5 min each)                  │
│ • Merge the good ones                                       │
│ • Add notes to the partial one                              │
│ Total review time: 25 minutes                               │
└─────────────────────────────────────────────────────────────┘
```

**Result:**
- **Code time:** 25 minutes (focused review)
- **AI time:** 40 minutes (parallel implementation)
- **Context switching:** 0 minutes (phone handled it)
- **Productivity:** 5 features instead of 1!

---

## 🚀 Implementation

### 1. Parent Mode Command (Termux)

```bash
#!/data/data/com.termux/files/usr/bin/bash
# ~/.termux/automation/parent-mode.sh
# Run multiple features in parallel while you handle kids

set -euo pipefail

PROJECT_DIR="$HOME/projects/tkan"
WORK_DIR="$HOME/tkan-parent-mode"
LOG_DIR="$HOME/.termux/logs/parent-mode"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

mkdir -p "$LOG_DIR"

log() {
    echo "[$(date +'%H:%M:%S')] $*" | tee -a "$LOG_DIR/parent-mode-$TIMESTAMP.log"
}

notify() {
    termux-notification \
        --title "Parent Mode: $1" \
        --content "$2" \
        --button1 "View" \
        --button1-action "termux-open $LOG_DIR/parent-mode-$TIMESTAMP.log"
}

log "🏁 Parent Mode Starting..."
notify "Started" "Running $TASK_COUNT parallel tasks"

# Read task list from config
TASKS_FILE="${1:-$HOME/.termux/parent-mode-tasks.json}"

if [ ! -f "$TASKS_FILE" ]; then
    log "❌ No tasks file found. Creating template..."
    cat > "$TASKS_FILE" << 'EOF'
{
  "tasks": [
    {
      "name": "calendar-view",
      "branch": "feature/calendar-view",
      "prompt": "Implement calendar view for tkan. See .claude/skills/tkan-calendar-view/IMPLEMENTATION.md for details. Create calendar.go, view_calendar.go, and update types.go. Write tests.",
      "priority": "high"
    },
    {
      "name": "table-view",
      "branch": "feature/table-view",
      "prompt": "Add table view mode to tkan. Create a sortable table showing all cards with columns: Title, Status, Assignee, Due Date, Tags. Use Bubble Tea table component.",
      "priority": "medium"
    },
    {
      "name": "export-csv",
      "branch": "feature/export-csv",
      "prompt": "Add CSV export functionality. Command: tkan export --format csv --output report.csv. Include all card metadata.",
      "priority": "low"
    },
    {
      "name": "refactor-backend",
      "branch": "refactor/backend-cleanup",
      "prompt": "Refactor backend_github.go to implement dynamic field discovery instead of hardcoded IDs. See .claude/skills/github-projects/references/field-discovery.md",
      "priority": "high"
    },
    {
      "name": "update-docs",
      "branch": "docs/update-readme",
      "prompt": "Update README.md with new features: calendar view, table view, export. Add screenshots and usage examples.",
      "priority": "medium"
    }
  ]
}
EOF
    log "📝 Template created at $TASKS_FILE"
    log "📝 Edit it with your tasks, then run again"
    exit 0
fi

# Parse tasks
TASK_COUNT=$(jq '.tasks | length' "$TASKS_FILE")
log "📋 Found $TASK_COUNT tasks"

# Create work directory with worktrees
mkdir -p "$WORK_DIR"
cd "$PROJECT_DIR"

# Clean up old worktrees from previous runs
git worktree prune

log "🔨 Setting up git worktrees..."

# Process each task in parallel
PIDS=()
for i in $(seq 0 $((TASK_COUNT - 1))); do
    TASK_NAME=$(jq -r ".tasks[$i].name" "$TASKS_FILE")
    BRANCH=$(jq -r ".tasks[$i].branch" "$TASKS_FILE")
    PROMPT=$(jq -r ".tasks[$i].prompt" "$TASKS_FILE")
    PRIORITY=$(jq -r ".tasks[$i].priority" "$TASKS_FILE")

    WORKTREE_DIR="$WORK_DIR/$TASK_NAME"

    log "📦 Setting up: $TASK_NAME ($PRIORITY priority)"

    # Create worktree
    if [ -d "$WORKTREE_DIR" ]; then
        rm -rf "$WORKTREE_DIR"
    fi

    git worktree add "$WORKTREE_DIR" -b "$BRANCH" 2>/dev/null || \
        git worktree add "$WORKTREE_DIR" "$BRANCH"

    # Run Claude in this worktree (in background)
    (
        cd "$WORKTREE_DIR"

        log "🤖 Starting Claude on $TASK_NAME..."

        # Run Claude with the task prompt
        claude --dangerously-skip-permissions "$PROMPT" \
            > "$LOG_DIR/$TASK_NAME-$TIMESTAMP.log" 2>&1

        EXIT_CODE=$?

        # Commit if changes were made
        if [[ -n $(git status -s) ]]; then
            git add -A
            git commit -m "feat($TASK_NAME): AI implementation

$PROMPT

🤖 Generated by Parent Mode automation
$(date)"
            log "✅ $TASK_NAME: Completed with changes"

            # Push to remote
            git push -u origin "$BRANCH" 2>&1 | tee -a "$LOG_DIR/$TASK_NAME-$TIMESTAMP.log"

        elif [ $EXIT_CODE -eq 0 ]; then
            log "✅ $TASK_NAME: Completed (no changes needed)"
        else
            log "⚠️  $TASK_NAME: Completed with errors (exit code: $EXIT_CODE)"
        fi
    ) &

    PIDS+=($!)
done

log "⏳ Waiting for all tasks to complete..."
log "💡 Go spend time with your kids! I'll notify when done."

# Initial notification
notify "In Progress" "Working on $TASK_COUNT tasks. Go be with the kids! 👶👶"

# Wait for all background jobs
for pid in "${PIDS[@]}"; do
    wait $pid
done

log "🎉 All tasks complete!"

# Generate summary
SUMMARY="Parent Mode Summary\n\n"
SUCCESS_COUNT=0
PARTIAL_COUNT=0
FAILED_COUNT=0

for i in $(seq 0 $((TASK_COUNT - 1))); do
    TASK_NAME=$(jq -r ".tasks[$i].name" "$TASKS_FILE")
    WORKTREE_DIR="$WORK_DIR/$TASK_NAME"

    cd "$WORKTREE_DIR"

    if [[ -n $(git log origin/master..HEAD --oneline 2>/dev/null) ]]; then
        SUMMARY+="✅ $TASK_NAME: Ready for review\n"
        ((SUCCESS_COUNT++))
    elif [[ -n $(git status -s) ]]; then
        SUMMARY+="⚠️  $TASK_NAME: Has uncommitted changes\n"
        ((PARTIAL_COUNT++))
    else
        SUMMARY+="❌ $TASK_NAME: No changes made\n"
        ((FAILED_COUNT++))
    fi
done

SUMMARY+="\nResults: $SUCCESS_COUNT ready, $PARTIAL_COUNT partial, $FAILED_COUNT failed"

log "$SUMMARY"

# Final notification with action buttons
termux-notification \
    --id parent-mode \
    --title "🎉 Parent Mode Complete!" \
    --content "$SUCCESS_COUNT tasks ready for review" \
    --priority high \
    --button1 "Review" \
    --button1-action "termux-open $LOG_DIR/parent-mode-$TIMESTAMP.log" \
    --button2 "GitHub" \
    --button2-action "termux-open-url https://github.com/GGPrompts/tkan/branches"

log "📊 Full log: $LOG_DIR/parent-mode-$TIMESTAMP.log"
log "🌲 Worktrees in: $WORK_DIR"
log "📝 Review branches on GitHub or locally"
```

### 2. Quick Trigger Command

```bash
# ~/.termux/shortcuts/parent-mode
#!/data/data/com.termux/files/usr/bin/bash
# Termux shortcut - appears in widget

# Run parent mode with default tasks
~/.termux/automation/parent-mode.sh

# Or with custom task file
# ~/.termux/automation/parent-mode.sh ~/my-urgent-tasks.json
```

### 3. PC Quick Review Script

```bash
# .local/scripts/review-parent-mode.sh
#!/bin/bash
# Quick review script for parent mode branches

set -euo pipefail

REPO_DIR="/home/matt/projects/tkan"
cd "$REPO_DIR"

echo "🔍 Parent Mode Review Dashboard"
echo "================================"
echo ""

# Fetch latest
git fetch --all --prune

# List all feature branches
BRANCHES=$(git branch -r | grep -E "feature/|refactor/|docs/" | sed 's/origin\///')

if [ -z "$BRANCHES" ]; then
    echo "No parent mode branches found"
    exit 0
fi

echo "📋 Available branches from parent mode:"
echo ""

# Show each branch with status
while IFS= read -r branch; do
    branch=$(echo "$branch" | xargs)  # trim whitespace

    # Get commit count ahead of master
    AHEAD=$(git rev-list --count master..origin/$branch 2>/dev/null || echo "0")

    # Get last commit message
    LAST_COMMIT=$(git log origin/$branch --oneline -1 2>/dev/null | cut -d' ' -f2- || echo "No commits")

    # Get file changes
    FILES_CHANGED=$(git diff master..origin/$branch --name-only 2>/dev/null | wc -l)

    echo "Branch: $branch"
    echo "  Commits: $AHEAD ahead of master"
    echo "  Files: $FILES_CHANGED changed"
    echo "  Last: $LAST_COMMIT"
    echo ""
done <<< "$BRANCHES"

echo "================================"
echo ""
echo "Options:"
echo "  1. Review branch"
echo "  2. Create PRs for all"
echo "  3. Merge all (dangerous!)"
echo "  4. Exit"
echo ""
read -p "Choose (1-4): " choice

case $choice in
    1)
        echo ""
        read -p "Branch name: " branch
        git checkout "$branch"
        git log master..HEAD --oneline
        echo ""
        git diff master..HEAD --stat
        echo ""
        read -p "Create PR? (y/n): " pr
        if [ "$pr" = "y" ]; then
            gh pr create --fill
        fi
        ;;

    2)
        echo "Creating PRs for all branches..."
        while IFS= read -r branch; do
            branch=$(echo "$branch" | xargs)
            echo "Creating PR for: $branch"
            git checkout "$branch"
            gh pr create --fill --label "ai-generated,parent-mode" || true
        done <<< "$BRANCHES"
        echo "✅ PRs created!"
        ;;

    3)
        echo "⚠️  This will merge all branches to master!"
        read -p "Are you SURE? (type 'yes'): " confirm
        if [ "$confirm" = "yes" ]; then
            git checkout master
            while IFS= read -r branch; do
                branch=$(echo "$branch" | xargs)
                echo "Merging: $branch"
                git merge --no-ff "origin/$branch" -m "Merge parent-mode: $branch"
            done <<< "$BRANCHES"
            git push
            echo "✅ All merged!"
        fi
        ;;

    4)
        echo "Bye!"
        ;;
esac
```

---

## 📋 Usage Workflow

### Morning (5 minutes)

```bash
# On phone (Termux) or via SSH from PC:

# 1. Edit task list
nano ~/.termux/parent-mode-tasks.json

# Add your 5 features for the day:
{
  "tasks": [
    {
      "name": "calendar-view",
      "branch": "feature/calendar-view",
      "prompt": "Implement calendar view...",
      "priority": "high"
    },
    // ... 4 more tasks
  ]
}

# 2. Trigger parent mode
~/.termux/automation/parent-mode.sh

# 3. Put phone back in pocket
# Go handle the kids! 👶👶
```

### Evening (25 minutes)

```bash
# When kids are asleep and you get notification:
# "🎉 Parent Mode Complete! 5 tasks ready"

# On PC:
cd ~/projects/tkan
./.local/scripts/review-parent-mode.sh

# Review each branch (5 min each):
# 1. Read what Claude did
# 2. Test it quickly
# 3. Create PR or merge
# 4. Move to next

# Total time: 25 minutes of focused review
# vs 2 hours of context switching!
```

---

## 🎯 Real-World Example

### Your Current Friday Night

```
8:00 PM - Kids asleep, open laptop
8:05 PM - Terminal 1: calendar view worktree
8:10 PM - Copy prompt about calendar view
8:15 PM - Terminal 2: table view worktree
8:17 PM - "DADDY I HAD A BAD DREAM!" 😭
8:25 PM - Back... wait, what was I doing?
8:30 PM - Terminal 3... or was it 4?
8:35 PM - Kid needs water 💧
8:40 PM - Give up on terminals 3-5
8:45 PM - Watch terminal 1 and 2
9:00 PM - One feature done, context switch to terminal 2
9:05 PM - Wife needs help with something
9:15 PM - Forget where I was
9:30 PM - Give up, go to bed
```

**Results:**
- 1 feature implemented
- 1 feature half-done
- 3 features abandoned
- Mental exhaustion 😵

### With Parent Mode

```
8:00 PM - Kids asleep, grab phone
8:02 PM - Edit task list (5 features)
8:05 PM - Run: parent-mode.sh
8:06 PM - Put phone down
8:07 PM - Go watch TV with wife 📺
         [Phone works on 5 features in parallel]
9:00 PM - 🔔 "Parent Mode Complete!"
9:05 PM - Open laptop
9:10 PM - Review branch 1: ✅ Looks good, merge
9:15 PM - Review branch 2: ✅ Needs tiny fix, done
9:20 PM - Review branch 3: ✅ Perfect, merge
9:25 PM - Review branch 4: ⚠️  Partial, add notes
9:30 PM - Review branch 5: ✅ Great, merge
9:35 PM - All done! Close laptop
```

**Results:**
- 4 features merged
- 1 feature has clear notes for next time
- Watched TV with wife
- Not mentally exhausted 😊

---

## 🔋 Battery Optimization

Since you're running 5 parallel Claude instances:

```bash
# Add to parent-mode.sh

# Check battery before starting
BATTERY=$(termux-battery-status | jq -r '.percentage')
CHARGING=$(termux-battery-status | jq -r '.status')

if [ "$BATTERY" -lt 50 ] && [ "$CHARGING" != "CHARGING" ]; then
    log "⚠️  Battery low ($BATTERY%). Please plug in phone."
    notify "Battery Low" "Plug in phone before starting parent mode"
    exit 1
fi

# Acquire wakelock
termux-wake-lock

# At end of script
termux-wake-unlock
```

---

## 💡 Pro Tips

### 1. Weekend Deep Work Mode

```json
// weekend-tasks.json
{
  "tasks": [
    {
      "name": "major-refactor-1",
      "branch": "refactor/backend-phase-1",
      "prompt": "Phase 1 of backend refactor...",
      "priority": "high"
    },
    // ... more ambitious tasks
  ]
}

// Run Friday night, review Sunday morning
```

### 2. Bug Fix Mode

```json
// bugs.json
{
  "tasks": [
    {
      "name": "fix-issue-42",
      "branch": "fix/issue-42",
      "prompt": "Fix: Cards not moving in drag mode. See issue #42.",
      "priority": "high"
    },
    {
      "name": "fix-issue-43",
      "branch": "fix/issue-43",
      "prompt": "Fix: Calendar view crashes on empty month.",
      "priority": "high"
    }
  ]
}
```

### 3. Documentation Sprint

```json
// docs.json
{
  "tasks": [
    {
      "name": "readme",
      "branch": "docs/update-readme",
      "prompt": "Update README with all new features",
      "priority": "high"
    },
    {
      "name": "api-docs",
      "branch": "docs/api-docs",
      "prompt": "Generate API documentation for all public functions",
      "priority": "medium"
    },
    {
      "name": "tutorials",
      "branch": "docs/tutorials",
      "prompt": "Write getting started tutorial",
      "priority": "low"
    }
  ]
}
```

---

## 🚨 Emergency: "Kid Woke Up Mid-Review"

```bash
# Quick save script
# .local/scripts/save-and-pause.sh

#!/bin/bash
# Save current state and pause

CURRENT_BRANCH=$(git branch --show-current)

# Stash any uncommitted work
git add -A
git stash save "WIP: Kid emergency at $(date)"

# Switch back to master
git checkout master

# Leave note for later
echo "Paused on $CURRENT_BRANCH at $(date)" >> ~/.parent-mode-pause.log

echo "✅ Saved! Go handle the kid. Resume later with:"
echo "   git checkout $CURRENT_BRANCH"
echo "   git stash pop"
```

Add to your shell:
```bash
alias emergency='~/.local/scripts/save-and-pause.sh'
```

Now when kid wakes up:
```bash
$ emergency
✅ Saved! Go handle the kid.
```

---

## 📊 Time Savings Analysis

### Without Parent Mode (Your Current Workflow)

| Activity | Time | Productivity |
|----------|------|--------------|
| Setup 5 terminals | 15 min | 0% |
| Context switching | 20 min | 0% |
| Actual coding | 25 min | 100% |
| Kid interruptions | 30 min | 0% |
| **Total** | **90 min** | **28%** |
| **Features completed** | **1-2** | |

### With Parent Mode

| Activity | Time | Productivity |
|----------|------|--------------|
| Setup task list | 5 min | 100% |
| AI working (in background) | 40 min | 0% (you're with kids!) |
| Review & merge | 25 min | 100% |
| **Total** | **30 min of YOUR time** | **100%** |
| **Features completed** | **4-5** | |

**Savings:**
- **60 minutes saved** per session
- **3-4x more features** completed
- **0 context switching stress**
- **Time with kids instead of frustrated coding**

---

## 🎁 Bonus: Termux Widget

Put parent mode on your home screen:

```bash
# 1. Install Termux:Widget from F-Droid

# 2. Create shortcut
mkdir -p ~/.shortcuts
cat > ~/.shortcuts/ParentMode << 'EOF'
#!/data/data/com.termux/files/usr/bin/bash
~/.termux/automation/parent-mode.sh
EOF

chmod +x ~/.shortcuts/ParentMode

# 3. Add Termux widget to home screen
# 4. Tap "ParentMode" to trigger!
```

Now you can literally tap your phone once and go back to the kids! 👶👶

---

## 🎉 Summary

**Before Parent Mode:**
- 5 terminals open
- 5 git worktrees manually created
- Copy/paste prompts 5 times
- Watch and context-switch
- Get interrupted constantly
- 90 minutes for 1-2 features
- Mental exhaustion

**After Parent Mode:**
- 1 phone tap
- Phone creates 5 worktrees automatically
- All 5 features work in parallel
- You go be with kids
- 🔔 Notification when done
- 25 minutes to review 4-5 features
- Happy wife, happy life 😊

**Your $200/month Claude Pro plan just became your AI development team that works while you're doing bedtime stories!** 🚀👶👶

---

**Files Created:**
- `.termux/automation/parent-mode.sh` - Main parent mode script
- `.local/scripts/review-parent-mode.sh` - Quick review dashboard
- Template: `~/.termux/parent-mode-tasks.json`

Want me to create any of these actual files for you?
