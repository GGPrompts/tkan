# Execution Pipeline Architecture

**Purpose:** Maximize Claude Max subscription usage by automating planned work execution while you're busy with kids.

**Problem:** Using only 30% of weekly credits in 2 days, wasting 70% of $200/month subscription.

**Solution:** Plan work manually (30% credits), automate execution (70% credits), review results when you have time.

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│ Phase 1: Planning (Monday-Tuesday) - You + Claude              │
│ • Brainstorm features                                           │
│ • Create detailed implementation plans                          │
│ • Break into phases with checkpoints                            │
│ • Generate task queue (JSON)                                    │
│ • Uses: 30% of weekly credits                                   │
└──────────────────┬──────────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│ Phase 2: Execution (Rest of Week) - Termux Automation          │
│ • Reads task queue                                              │
│ • Executes phases sequentially                                  │
│ • Self-review checkpoints between phases                        │
│ • Updates tkan board in real-time                               │
│ • Creates PRs for review                                        │
│ • Uses: 60-70% of weekly credits (otherwise wasted)             │
└──────────────────┬──────────────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────────────┐
│ Phase 3: Review (Next Monday) - You                             │
│ • Review PRs from automation                                    │
│ • Merge good work                                               │
│ • Fix issues from failed checkpoints                            │
│ • Plan next batch                                               │
│ • Time: 25-60 minutes focused review                            │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📋 Task Queue Format

### Queue File Structure

```
~/.termux/task-queue/
├── pending/
│   ├── export-feature.json
│   ├── calendar-view.json
│   └── table-view.json
├── in-progress/
│   └── current-task.json
├── completed/
│   └── archive/
└── failed/
    └── needs-review/
```

### Task Definition Format

```json
{
  "id": "export-feature-2025-01",
  "title": "Add CSV/JSON Export to tkan",
  "github_project": 7,
  "github_issue": 42,
  "created": "2025-01-29T10:00:00Z",
  "priority": "high",
  "estimated_credits": "medium",

  "phases": [
    {
      "phase": 1,
      "name": "Foundation",
      "goal": "Create export interface and stubs",
      "estimated_duration": "30min",
      "tasks": [
        {
          "task_id": "export-1-interface",
          "branch": "feature/export-phase1-foundation",
          "worktree_dir": "~/tkan-work/export-phase1",
          "base_branch": "master",

          "implementation": {
            "prompt": "Create export.go with Exporter interface:\n\ntype Exporter interface {\n    Export(cards []Card, w io.Writer) error\n    Format() string\n}\n\nImplement two stubs:\n- CSVExporter struct\n- JSONExporter struct\n\nFollow existing tkan patterns. Use same error handling style as backend_github.go.",
            "context_files": [
              "types.go",
              "backend_github.go"
            ],
            "max_duration": "20min"
          },

          "validation": {
            "prompt": "Review export.go implementation:\n\n1. Does interface match specification?\n2. Are stubs correctly structured?\n3. Follows tkan code style?\n4. Proper error handling?\n5. Any missing edge cases?\n\nProvide: PASS/FAIL + specific issues if any.",
            "pass_criteria": "All 5 checks pass",
            "max_duration": "10min"
          }
        }
      ],

      "checkpoint": {
        "type": "code_review",
        "reviewer": "claude_agent_2",
        "prompt": "Comprehensive Phase 1 review:\n\n1. Review all code from Phase 1\n2. Check: interface design, implementation quality, tests\n3. Run: go build, go test\n4. Decision: PASS or FAIL\n\nIf FAIL: List blocking issues\nIf PASS: Approve progression to Phase 2",
        "on_pass": {
          "action": "create_pr",
          "pr_title": "feat: Add export interface (Phase 1)",
          "pr_body": "Phase 1 complete: Export interface foundation\n\nReviewed by: Automated Claude agent\nStatus: Ready for human review",
          "next_phase": 2
        },
        "on_fail": {
          "action": "create_issue",
          "issue_title": "Automation Failed: Export Phase 1",
          "stop_execution": true,
          "notify": true
        }
      }
    },

    {
      "phase": 2,
      "name": "Implementation",
      "goal": "Fully implement CSV and JSON exporters",
      "depends_on": "phase_1_approved",
      "tasks": [
        {
          "task_id": "export-2-csv",
          "branch": "feature/export-phase2-csv",
          "implementation": {
            "prompt": "Implement CSVExporter.Export() method:\n\n- Header row: Title, Status, Assignee, Labels, Created, Updated\n- Proper CSV escaping (quotes, commas, newlines)\n- Handle empty fields gracefully\n- Write comprehensive tests in export_csv_test.go\n\nTest with real tkan data from project #7.",
            "max_duration": "45min"
          },
          "validation": {
            "prompt": "Validate CSV export:\n\n1. Run all tests: go test ./...\n2. Manual test: export project #7 to CSV\n3. Open CSV in spreadsheet, verify formatting\n4. Check edge cases: special chars, empty fields\n\nPASS/FAIL + issues",
            "max_duration": "15min"
          }
        },
        {
          "task_id": "export-2-json",
          "branch": "feature/export-phase2-json",
          "parallel_with": "export-2-csv",
          "implementation": {
            "prompt": "Implement JSONExporter.Export() method:\n\n- Pretty-printed JSON with 2-space indent\n- Include all card fields\n- Proper JSON escaping\n- Write tests in export_json_test.go\n\nTest with real tkan data.",
            "max_duration": "45min"
          }
        }
      ],

      "checkpoint": {
        "type": "integration_test",
        "prompt": "Phase 2 integration review:\n\n1. Both exporters implemented?\n2. All tests pass?\n3. Export project #7 to both formats\n4. Verify output quality\n5. Check error handling\n\nPASS/FAIL decision",
        "on_pass": {
          "action": "merge_branches_and_create_pr",
          "merge_into": "feature/export-complete",
          "next_phase": 3
        }
      }
    },

    {
      "phase": 3,
      "name": "Polish",
      "goal": "Error handling, docs, CLI integration",
      "tasks": [
        {
          "task_id": "export-3-cli",
          "branch": "feature/export-phase3-polish",
          "implementation": {
            "prompt": "Add CLI commands:\n\n- tkan export --format csv --output cards.csv\n- tkan export --format json --output cards.json\n- Default output: stdout\n\nUpdate main.go with new commands.\nAdd help text and examples.\nUpdate README.md with export documentation.",
            "max_duration": "30min"
          }
        }
      ],

      "checkpoint": {
        "type": "final_review",
        "prompt": "Final review before PR:\n\n1. All features working?\n2. Tests comprehensive?\n3. Documentation complete?\n4. CLI intuitive?\n5. Ready for production?\n\nCreate final PR if PASS.",
        "on_pass": {
          "action": "create_final_pr",
          "pr_title": "feat: Add CSV and JSON export functionality",
          "pr_body": "Complete export feature implementation.\n\nPhases completed:\n- Phase 1: Interface foundation ✅\n- Phase 2: CSV/JSON implementation ✅  \n- Phase 3: CLI integration & docs ✅\n\nAutomated implementation + review.\nReady for human approval and merge.",
          "labels": ["feature", "automated", "ready-for-review"]
        }
      }
    }
  ],

  "notifications": {
    "on_phase_complete": true,
    "on_checkpoint_fail": true,
    "on_final_complete": true,
    "method": "termux-notification"
  },

  "tkan_integration": {
    "update_status": true,
    "create_subtasks": true,
    "track_progress": true
  }
}
```

---

## 🤖 Execution Engine

### Core Script: `~/.termux/automation/execution-engine.sh`

```bash
#!/data/data/com.termux/files/usr/bin/bash
# Execution pipeline for automated task implementation

set -euo pipefail

QUEUE_DIR="$HOME/.termux/task-queue"
WORK_DIR="$HOME/tkan-work"
LOG_DIR="$HOME/.termux/logs/execution"
PROJECT_DIR="$HOME/projects/tkan"

mkdir -p "$LOG_DIR" "$WORK_DIR"

log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_DIR/engine.log"
}

notify() {
    termux-notification --title "Execution Pipeline: $1" --content "$2"
}

# Main execution loop
process_queue() {
    # Find next pending task
    TASK_FILE=$(ls -1 "$QUEUE_DIR/pending/"*.json 2>/dev/null | head -1)

    if [ -z "$TASK_FILE" ]; then
        log "No pending tasks in queue"
        return 0
    fi

    TASK_ID=$(jq -r '.id' "$TASK_FILE")
    log "Processing task: $TASK_ID"

    # Move to in-progress
    mv "$TASK_FILE" "$QUEUE_DIR/in-progress/"
    TASK_FILE="$QUEUE_DIR/in-progress/$(basename $TASK_FILE)"

    # Execute each phase
    PHASE_COUNT=$(jq '.phases | length' "$TASK_FILE")

    for phase_num in $(seq 1 $PHASE_COUNT); do
        execute_phase "$TASK_FILE" "$phase_num"

        # Check if checkpoint passed
        if ! check_checkpoint "$TASK_FILE" "$phase_num"; then
            log "Phase $phase_num checkpoint failed, stopping execution"
            mv "$TASK_FILE" "$QUEUE_DIR/failed/"
            notify "Pipeline Failed" "Task $TASK_ID failed at Phase $phase_num"
            return 1
        fi

        notify "Phase Complete" "Task $TASK_ID Phase $phase_num ✅"
    done

    # All phases complete
    mv "$TASK_FILE" "$QUEUE_DIR/completed/"
    notify "Task Complete" "All phases of $TASK_ID completed successfully!"
    log "Task $TASK_ID completed successfully"
}

execute_phase() {
    local task_file=$1
    local phase_num=$2

    log "Executing Phase $phase_num"

    # Get phase details
    local phase_name=$(jq -r ".phases[$((phase_num-1))].name" "$task_file")
    local task_count=$(jq ".phases[$((phase_num-1))].tasks | length" "$task_file")

    log "Phase: $phase_name, Tasks: $task_count"

    # Execute each task in the phase
    for task_idx in $(seq 0 $((task_count-1))); do
        execute_task "$task_file" "$phase_num" "$task_idx"
    done
}

execute_task() {
    local task_file=$1
    local phase_num=$2
    local task_idx=$3

    # Extract task details
    local task_id=$(jq -r ".phases[$((phase_num-1))].tasks[$task_idx].task_id" "$task_file")
    local branch=$(jq -r ".phases[$((phase_num-1))].tasks[$task_idx].branch" "$task_file")
    local worktree_dir=$(jq -r ".phases[$((phase_num-1))].tasks[$task_idx].worktree_dir" "$task_file")
    local impl_prompt=$(jq -r ".phases[$((phase_num-1))].tasks[$task_idx].implementation.prompt" "$task_file")
    local val_prompt=$(jq -r ".phases[$((phase_num-1))].tasks[$task_idx].validation.prompt" "$task_file")

    log "Task: $task_id on branch $branch"

    # Expand tilde in worktree_dir
    worktree_dir="${worktree_dir/#\~/$HOME}"

    # Create git worktree
    cd "$PROJECT_DIR"
    if [ -d "$worktree_dir" ]; then
        rm -rf "$worktree_dir"
    fi
    git worktree add "$worktree_dir" -b "$branch" 2>/dev/null || \
        git worktree add "$worktree_dir" "$branch"

    # Implementation phase
    log "Running implementation..."
    cd "$worktree_dir"

    claude --dangerously-skip-permissions "$impl_prompt" \
        > "$LOG_DIR/$task_id-impl.log" 2>&1

    # Commit if changes
    if [[ -n $(git status -s) ]]; then
        git add -A
        git commit -m "feat($task_id): Automated implementation

$impl_prompt

🤖 Automated by execution pipeline
Phase: $phase_num
$(date)"
    fi

    # Validation phase
    log "Running validation..."

    claude --dangerously-skip-permissions "$val_prompt" \
        > "$LOG_DIR/$task_id-validation.log" 2>&1

    # Check validation result
    if grep -qi "PASS" "$LOG_DIR/$task_id-validation.log"; then
        log "Task $task_id: PASSED validation"
        git push -u origin "$branch"
    else
        log "Task $task_id: FAILED validation"
        return 1
    fi

    # Update tkan board
    update_tkan_status "$task_id" "completed"
}

check_checkpoint() {
    local task_file=$1
    local phase_num=$2

    log "Running Phase $phase_num checkpoint..."

    local checkpoint_prompt=$(jq -r ".phases[$((phase_num-1))].checkpoint.prompt" "$task_file")
    local checkpoint_type=$(jq -r ".phases[$((phase_num-1))].checkpoint.type" "$task_file")

    # Run checkpoint review
    cd "$WORK_DIR"

    claude --dangerously-skip-permissions "$checkpoint_prompt" \
        > "$LOG_DIR/checkpoint-phase$phase_num.log" 2>&1

    # Check result
    if grep -qi "PASS" "$LOG_DIR/checkpoint-phase$phase_num.log"; then
        log "Checkpoint Phase $phase_num: PASSED"

        # Execute on_pass action
        local action=$(jq -r ".phases[$((phase_num-1))].checkpoint.on_pass.action" "$task_file")
        execute_checkpoint_action "$task_file" "$phase_num" "pass"

        return 0
    else
        log "Checkpoint Phase $phase_num: FAILED"

        # Execute on_fail action
        execute_checkpoint_action "$task_file" "$phase_num" "fail"

        return 1
    fi
}

execute_checkpoint_action() {
    local task_file=$1
    local phase_num=$2
    local result=$3  # pass or fail

    local action=$(jq -r ".phases[$((phase_num-1))].checkpoint.on_$result.action" "$task_file")

    case $action in
        create_pr)
            local pr_title=$(jq -r ".phases[$((phase_num-1))].checkpoint.on_$result.pr_title" "$task_file")
            local pr_body=$(jq -r ".phases[$((phase_num-1))].checkpoint.on_$result.pr_body" "$task_file")

            cd "$PROJECT_DIR"
            gh pr create --title "$pr_title" --body "$pr_body" --label "automated"
            ;;

        create_issue)
            local issue_title=$(jq -r ".phases[$((phase_num-1))].checkpoint.on_$result.issue_title" "$task_file")

            cd "$PROJECT_DIR"
            gh issue create --title "$issue_title" --body "Automated execution failed. See logs in $LOG_DIR"
            ;;
    esac
}

update_tkan_status() {
    local task_id=$1
    local status=$2

    # Update GitHub Project via gh CLI
    # This updates the tkan board in real-time

    cd "$PROJECT_DIR"

    # Find item in project by title matching task_id
    # Update status field
    # (Implementation depends on GitHub Projects API)

    log "Updated tkan status for $task_id: $status"
}

# Main loop
main() {
    log "Execution engine started"

    # Check battery
    if command -v termux-battery-status &>/dev/null; then
        BATTERY=$(termux-battery-status | jq -r '.percentage')
        CHARGING=$(termux-battery-status | jq -r '.status')

        if [ "$BATTERY" -lt 30 ] && [ "$CHARGING" != "CHARGING" ]; then
            log "Battery too low ($BATTERY%), skipping execution"
            exit 0
        fi

        termux-wake-lock
    fi

    # Process queue
    while process_queue; do
        log "Checking for more tasks..."
        sleep 5
    done

    if command -v termux-wake-unlock &>/dev/null; then
        termux-wake-unlock
    fi

    log "Execution engine finished"
}

main "$@"
```

### Cron Schedule

```bash
# Run execution engine every 6 hours (while charging)
0 */6 * * * [ "$(termux-battery-status | jq -r '.status')" = "CHARGING" ] && \
  ~/.termux/automation/execution-engine.sh
```

---

## 🎯 tkan Integration

### Real-Time Status Updates

As automation executes, update tkan board:

```go
// In tkan: Add automation status field
type Card struct {
    // ... existing fields
    AutomationStatus string    `json:"automation_status,omitempty"`
    AutomationPhase  int        `json:"automation_phase,omitempty"`
    LastAutomated    time.Time  `json:"last_automated,omitempty"`
}

// Status values:
// - "" (empty) = manual task
// - "queued" = in task queue
// - "phase_1" = Phase 1 executing
// - "review" = At checkpoint
// - "complete" = Automation finished
// - "failed" = Automation failed
```

### Visual Indicators in TUI

```go
// In tkan view, show automation status with icons
func renderCard(card Card) string {
    icon := ""
    switch card.AutomationStatus {
    case "queued":
        icon = "⏳"
    case "phase_1", "phase_2", "phase_3":
        icon = "🤖"
    case "review":
        icon = "🔍"
    case "complete":
        icon = "✅"
    case "failed":
        icon = "❌"
    }

    return fmt.Sprintf("%s %s", icon, card.Title)
}
```

### Automation Tab

Add new view mode to tkan:

```
┌─ tkan: Automation Status ─────────────────────────┐
│                                                    │
│ Queued (2)                                         │
│ ⏳ Add export functionality                       │
│ ⏳ Calendar view implementation                   │
│                                                    │
│ In Progress (1)                                    │
│ 🤖 Table view - Phase 2/3                         │
│    Started: 2 hours ago                           │
│    Current: Implementing table rendering          │
│                                                    │
│ Review Needed (1)                                  │
│ 🔍 Search feature - Checkpoint failed             │
│    Issue: Tests failing, needs manual fix         │
│                                                    │
│ Completed Today (3)                                │
│ ✅ Documentation update                            │
│ ✅ Bug fix #42                                     │
│ ✅ Error handling improvements                    │
│                                                    │
│ Credit Usage:                                      │
│ This week: ████████░░ 82%                         │
│ Automation: ██████░░░░ 52%                        │
│                                                    │
└────────────────────────────────────────────────────┘

Press 'a' for automation view, 'b' for board view
```

---

## 🔄 Multi-Agent Review System

### Agent Roles

```
Agent 1 (Implementer):
├─ Reads task prompt
├─ Implements feature
├─ Writes tests
└─ Commits code

Agent 2 (Reviewer):
├─ Reviews Agent 1's code
├─ Runs tests
├─ Checks style/patterns
└─ PASS/FAIL decision

Agent 3 (Integrator):
├─ Merges phase branches
├─ Runs integration tests
├─ Checks for conflicts
└─ Creates final PR
```

### Implementation

```bash
# In execution-engine.sh

execute_task_with_review() {
    local task_file=$1
    local phase_num=$2
    local task_idx=$3

    # Agent 1: Implementation
    log "Agent 1 (Implementer): Starting..."

    claude --dangerously-skip-permissions "$impl_prompt" \
        > "$LOG_DIR/agent1-impl.log" 2>&1

    git add -A
    git commit -m "feat: Agent 1 implementation"

    # Agent 2: Review
    log "Agent 2 (Reviewer): Starting review..."

    REVIEW_PROMPT="You are Agent 2, code reviewer.

Review the changes Agent 1 made:
$(git diff HEAD~1)

Check:
1. Code quality and style
2. Test coverage
3. Error handling
4. Follows tkan patterns
5. Any bugs or issues

Respond: PASS or FAIL + detailed feedback"

    claude --dangerously-skip-permissions "$REVIEW_PROMPT" \
        > "$LOG_DIR/agent2-review.log" 2>&1

    # Check review result
    if grep -qi "PASS" "$LOG_DIR/agent2-review.log"; then
        log "Agent 2: APPROVED implementation"
        git push origin "$branch"
    else
        log "Agent 2: REJECTED implementation"

        # Agent 1 retry with Agent 2 feedback
        RETRY_PROMPT="Agent 2 found issues:

$(cat $LOG_DIR/agent2-review.log)

Please fix these issues."

        claude --dangerously-skip-permissions "$RETRY_PROMPT" \
            > "$LOG_DIR/agent1-retry.log" 2>&1

        git add -A
        git commit -m "fix: Address Agent 2 feedback"
        git push origin "$branch"
    fi
}
```

---

## 📊 Credit Usage Optimization

### Tracking

```json
// ~/.termux/credit-tracking.json
{
  "week_start": "2025-01-27",
  "total_credits": 100,
  "usage": {
    "planning": {
      "credits_used": 30,
      "sessions": [
        {
          "date": "2025-01-27T10:00:00Z",
          "task": "Export feature planning",
          "credits": 15
        },
        {
          "date": "2025-01-28T14:00:00Z",
          "task": "Calendar view planning",
          "credits": 15
        }
      ]
    },
    "automation": {
      "credits_used": 52,
      "tasks": [
        {
          "date": "2025-01-28T20:00:00Z",
          "task": "Export Phase 1",
          "credits": 10
        },
        {
          "date": "2025-01-29T02:00:00Z",
          "task": "Export Phase 2",
          "credits": 25
        },
        {
          "date": "2025-01-30T08:00:00Z",
          "task": "Export Phase 3",
          "credits": 12
        },
        {
          "date": "2025-01-30T14:00:00Z",
          "task": "Code reviews",
          "credits": 5
        }
      ]
    },
    "manual": {
      "credits_used": 8,
      "adhoc_tasks": true
    }
  },
  "utilization": "90%",
  "wasted": "10%",
  "effective_cost_per_month": "$222"
}
```

### Optimization Strategies

1. **Task Sizing**
   - Small tasks (5-10 credits): Quick wins, frequent completion
   - Medium tasks (15-25 credits): Main features
   - Large tasks (30+ credits): Complex features, multiple phases

2. **Checkpoint Efficiency**
   - Quick validation: 2-5 credits
   - Code review: 5-10 credits
   - Integration testing: 10-15 credits

3. **Peak Usage Times**
   - Nights: 8pm - 6am (you're asleep, phone works)
   - Weekends: Saturday/Sunday (you're with kids)
   - Weekdays: During work hours (you're busy)

4. **Credit Budget Per Feature**
   ```
   Small feature (total: 30 credits):
   ├─ Planning: 10 credits
   ├─ Automation: 15 credits
   └─ Review/fixes: 5 credits

   Medium feature (total: 60 credits):
   ├─ Planning: 20 credits
   ├─ Automation: 30 credits
   └─ Review/fixes: 10 credits

   Large feature (total: 100+ credits):
   ├─ Planning: 30 credits
   ├─ Automation: 60 credits
   └─ Review/fixes: 20 credits
   ```

---

## 🚀 Getting Started

### 1. Create Planning Workflow

```bash
# Create alias for planning
alias tkan-plan='claude "Help me plan this feature for automated execution:

Break into 3-5 phases with checkpoints.
Each task: 20-40 min implementation.
Include validation prompts.
Generate JSON task queue format.
Optimize for automation execution.

Feature to plan:"'

# Usage:
tkan-plan "Add CSV/JSON export to tkan"

# Saves output to task queue
```

### 2. Set Up Termux

```bash
# Install packages
pkg install git gh jq cronie termux-services termux-api

# Enable cron
sv-enable crond

# Create directories
mkdir -p ~/.termux/task-queue/{pending,in-progress,completed,failed}
mkdir -p ~/.termux/automation
mkdir -p ~/.termux/logs/execution

# Copy execution-engine.sh
# (Create script from architecture above)

# Add to crontab
crontab -e
# Add: 0 */6 * * * [ "$(termux-battery-status | jq -r '.status')" = "CHARGING" ] && ~/.termux/automation/execution-engine.sh
```

### 3. Test the Pipeline

```bash
# Create a simple test task
cat > ~/.termux/task-queue/pending/test-task.json << 'EOF'
{
  "id": "test-001",
  "title": "Test Pipeline",
  "phases": [
    {
      "phase": 1,
      "tasks": [{
        "task_id": "test-1",
        "branch": "test/pipeline",
        "worktree_dir": "~/tkan-work/test",
        "implementation": {
          "prompt": "Create a file called TEST.md with 'Pipeline works!'"
        },
        "validation": {
          "prompt": "Check if TEST.md exists and contains 'Pipeline works!'. Respond PASS or FAIL."
        }
      }],
      "checkpoint": {
        "prompt": "Review TEST.md. Respond PASS if it exists.",
        "on_pass": {"action": "create_pr"}
      }
    }
  ]
}
EOF

# Run manually
~/.termux/automation/execution-engine.sh

# Check results
ls ~/.termux/task-queue/completed/
```

---

## 💰 Expected ROI

### Before Pipeline
```
Weekly credits: 100
Used: 30 (planning + manual coding)
Wasted: 70
Monthly cost: $200
Effective cost: $666/month
Features/week: 1-2
```

### After Pipeline
```
Weekly credits: 100
Planning: 30 (manual with Claude)
Automation: 60 (while you're busy)
Wasted: 10
Monthly cost: $200
Effective cost: $222/month
Features/week: 5-8

Savings: $444/month effective cost
Time saved: 10-15 hours/week of active coding
Family time: Priceless
```

---

## 📝 Next Steps

1. **Document your planning workflow**
   - How you currently brainstorm features
   - What makes a good phase breakdown
   - Common patterns you use

2. **Create first task queue manually**
   - Pick a small feature
   - Write JSON by hand
   - Test execution

3. **Iterate on format**
   - What works well?
   - What's too verbose?
   - What's missing?

4. **Build tooling gradually**
   - Start with simple bash scripts
   - Add tkan integration when needed
   - Build dashboard when you have time

5. **Optimize for your workflow**
   - Parent of toddlers = interruption-resistant
   - Bubbletea developer = visual feedback important
   - Credit maximization = longer automated tasks

---

## 🎯 Philosophy

**This isn't about replacing you with AI.**

This is about:
- Using credits you're already paying for
- Automating grunt work while you're busy
- Spending time with kids instead of boilerplate
- Reviewing and guiding rather than typing
- Building your library of TUI apps faster

**You're the architect, Claude is the construction crew.**

---

**Status:** Architecture documented, ready for implementation when you have time.

**Integration:** Designed specifically for tkan + Termux + Claude Max + toddler parent lifestyle.
