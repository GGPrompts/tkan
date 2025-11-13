#!/bin/bash
# review-scheduled-prompts.sh - Interactive prompt review tool
# Usage: Run manually to process scheduled prompts

set -euo pipefail

PROMPTS_DIR="/home/matt/projects/tkan/.claude/scheduled-prompts"

if [ ! -d "$PROMPTS_DIR" ]; then
    echo "No scheduled prompts found"
    exit 0
fi

# Count pending prompts
PENDING=$(find "$PROMPTS_DIR" -name "*pending.md" -type f | wc -l)

if [ "$PENDING" -eq 0 ]; then
    echo "✅ No pending prompts to review"
    exit 0
fi

echo "📋 Found $PENDING scheduled prompt(s)"
echo ""

# List prompts with priority
echo "Available prompts:"
find "$PROMPTS_DIR" -name "*pending.md" -type f | while read -r file; do
    PRIORITY=$(grep "^priority:" "$file" | cut -d' ' -f2)
    TITLE=$(basename "$file" .md | cut -d'-' -f2-)

    case $PRIORITY in
        high)   ICON="🔴" ;;
        normal) ICON="🟡" ;;
        low)    ICON="🟢" ;;
        *)      ICON="⚪" ;;
    esac

    echo "  $ICON $(basename "$file")"
done

echo ""
echo "Options:"
echo "  1. Open all prompts in Claude Code"
echo "  2. List prompts"
echo "  3. Mark all as complete"
echo "  4. Cancel"
echo ""
read -p "Choose action (1-4): " choice

case $choice in
    1)
        # Open in Claude Code (or default editor)
        for file in "$PROMPTS_DIR"/*pending.md; do
            echo "Opening: $(basename "$file")"
            cat "$file"
            echo ""
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo ""
        done

        echo ""
        read -p "Mark all as complete? (y/n): " complete
        if [ "$complete" = "y" ]; then
            find "$PROMPTS_DIR" -name "*pending.md" -type f | while read -r file; do
                mv "$file" "${file/pending/completed}"
            done
            echo "✅ All prompts marked as completed"
        fi
        ;;

    2)
        # List with content preview
        find "$PROMPTS_DIR" -name "*pending.md" -type f | while read -r file; do
            echo "═══════════════════════════════════════"
            echo "File: $(basename "$file")"
            echo "═══════════════════════════════════════"
            cat "$file"
            echo ""
        done
        ;;

    3)
        # Mark all complete
        find "$PROMPTS_DIR" -name "*pending.md" -type f | while read -r file; do
            mv "$file" "${file/pending/completed}"
        done
        echo "✅ All prompts marked as completed"
        ;;

    4)
        echo "Cancelled"
        exit 0
        ;;
esac
