#!/bin/bash
# scheduled-prompts.sh - Create scheduled prompts for Claude Code review
# Usage: Run via cron to prepare daily/weekly prompts

set -euo pipefail

TKAN_DIR="/home/matt/projects/tkan"
PROMPTS_DIR="$TKAN_DIR/.claude/scheduled-prompts"
TODAY=$(date +%Y-%m-%d)

mkdir -p "$PROMPTS_DIR"

# Function to create a prompt file
create_prompt() {
    local name=$1
    local content=$2
    local priority=${3:-normal}

    cat > "$PROMPTS_DIR/${TODAY}-${name}.md" << EOF
---
created: $(date -Iseconds)
priority: $priority
status: pending
---

# Scheduled Task: $name

$content

---
**Auto-generated:** $(date)
**Review this prompt and run it when ready**
EOF

    echo "Created prompt: $name"
}

# Check for code changes without doc updates
check_documentation_sync() {
    cd "$TKAN_DIR"

    # Files changed in last 7 days
    CHANGED_GO=$(find . -name "*.go" -mtime -7 -type f | wc -l)
    CHANGED_MD=$(find . -name "*.md" -mtime -7 -type f | wc -l)

    if [ "$CHANGED_GO" -gt 0 ] && [ "$CHANGED_MD" -eq 0 ]; then
        create_prompt "doc-sync" "$(cat <<'PROMPT'
Review the codebase changes from the last week and update documentation:

1. Run: git log --since="1 week ago" --name-only --pretty=format: | sort -u
2. Review changed .go files
3. Check if README.md, CLAUDE.md, or API docs need updates
4. Update any outdated examples or documentation
5. Create a commit with doc updates

Focus on:
- Changed function signatures
- New features that need documenting
- Outdated code examples
- Broken links
PROMPT
)" "high"
    fi
}

# Analyze project health
analyze_project_health() {
    cd "$TKAN_DIR"

    # Get project stats
    TODO_COUNT=$(gh project item-list 7 --owner GGPrompts --format json | \
        jq '[.items[] | select(.status == "Todo")] | length')
    IN_PROGRESS=$(gh project item-list 7 --owner GGPrompts --format json | \
        jq '[.items[] | select(.status == "In Progress")] | length')

    # Check for stale items (created >14 days ago, still in progress)
    STALE_ITEMS=$(gh project item-list 7 --owner GGPrompts --format json | \
        jq -r '.items[] |
            select(.status == "In Progress") |
            select(.createdAt < (now - 1209600 | strftime("%Y-%m-%dT%H:%M:%SZ"))) |
            .content.title')

    if [ -n "$STALE_ITEMS" ]; then
        create_prompt "stale-items-review" "$(cat <<PROMPT
Review and update stale project items:

The following tasks have been "In Progress" for over 2 weeks:
$STALE_ITEMS

For each task:
1. Assess current status
2. Either:
   - Move to "Done" if completed
   - Move back to "Todo" if blocked
   - Break into smaller tasks if too large
   - Add notes about blockers

Then update the project board accordingly.
PROMPT
)" "normal"
    fi
}

# Check for missing tests
check_test_coverage() {
    cd "$TKAN_DIR"

    # Find .go files without corresponding _test.go
    MISSING_TESTS=""
    for gofile in *.go; do
        # Skip test files themselves
        if [[ "$gofile" == *"_test.go" ]]; then
            continue
        fi

        # Skip main.go
        if [[ "$gofile" == "main.go" ]]; then
            continue
        fi

        testfile="${gofile%.go}_test.go"
        if [[ ! -f "$testfile" ]]; then
            MISSING_TESTS="$MISSING_TESTS\n- $gofile"
        fi
    done

    if [ -n "$MISSING_TESTS" ]; then
        create_prompt "add-tests" "$(cat <<PROMPT
The following files are missing test coverage:
$MISSING_TESTS

For each file:
1. Review the functions that need testing
2. Create table-driven tests
3. Aim for >80% coverage
4. Test edge cases and error handling

Create a new test file for the most critical file first.
PROMPT
)" "low"
    fi
}

# Generate weekly summary (only on Sundays)
generate_weekly_summary() {
    if [ "$(date +%u)" -eq 7 ]; then  # Sunday
        cd "$TKAN_DIR"

        # Get activity from last week
        COMMITS_LAST_WEEK=$(git log --since="1 week ago" --oneline | wc -l)

        create_prompt "weekly-summary" "$(cat <<PROMPT
Generate a weekly summary for tkan project:

1. Review git log from last week:
   - Number of commits: $COMMITS_LAST_WEEK
   - Key changes and features added

2. Check GitHub Project #7:
   - Tasks completed this week
   - Tasks still in progress
   - New tasks added

3. Create a summary report covering:
   - Accomplishments
   - Challenges/blockers
   - Next week's priorities
   - Any technical debt identified

4. Save the report as: docs/weekly-summaries/$(date +%Y-%m-%d).md

5. Create a GitHub Issue with the summary for tracking
PROMPT
)" "normal"
    fi
}

# Run all checks
echo "Running scheduled checks for tkan..."
check_documentation_sync
analyze_project_health
check_test_coverage
generate_weekly_summary

# Count pending prompts
PENDING=$(ls -1 "$PROMPTS_DIR"/*-pending.md 2>/dev/null | wc -l)

if [ "$PENDING" -gt 0 ]; then
    echo ""
    echo "📋 $PENDING scheduled prompt(s) ready for review"
    echo "Location: $PROMPTS_DIR"

    # Send desktop notification
    notify-send "tkan Scheduled Prompts" \
        "$PENDING task(s) ready for Claude Code review" \
        -i dialog-information

    # Add to GitHub Project
    gh project item-create 7 --owner GGPrompts \
        --title "🤖 Review scheduled prompts ($PENDING tasks)" \
        --body "Check .claude/scheduled-prompts/ for automated analysis results"
fi
