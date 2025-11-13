#!/bin/bash
# bulk-operations.sh - Perform bulk operations on project items
# Usage: ./bulk-operations.sh <command> [args...]

set -euo pipefail

COMMAND=${1:-}

usage() {
    cat << EOF
Usage: $0 <command> [args...]

Commands:
  move-all <project-id> <from-status> <to-status> [field-cache]
    Move all items from one status to another

  create-batch <owner> <project-number> <file>
    Create multiple items from a file (one title per line)

  archive-done <owner> <project-number> [field-cache]
    Archive all items in "Done" status

  list-by-status <owner> <project-number> [status]
    List all items, optionally filtered by status

  export-csv <owner> <project-number> <output-file>
    Export all items to CSV format

Examples:
  $0 move-all PVT_xxx "Todo" "In Progress"
  $0 create-batch GGPrompts 7 tasks.txt
  $0 list-by-status GGPrompts 7 "In Progress"
  $0 export-csv GGPrompts 7 project.csv

EOF
}

if [[ -z "$COMMAND" ]]; then
    usage
    exit 1
fi

# Get the directory containing this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

case "$COMMAND" in
    move-all)
        PROJECT_ID=${2:-}
        FROM_STATUS=${3:-}
        TO_STATUS=${4:-}
        FIELD_CACHE=${5:-fields-cache.json}

        if [[ -z "$PROJECT_ID" ]] || [[ -z "$FROM_STATUS" ]] || [[ -z "$TO_STATUS" ]]; then
            echo "Error: Missing arguments" >&2
            echo "Usage: $0 move-all <project-id> <from-status> <to-status> [field-cache]" >&2
            exit 1
        fi

        # Get project owner and number from cache
        OWNER=$(jq -r '.owner' "$FIELD_CACHE")
        PROJECT_NUMBER=$(jq -r '.project_number' "$FIELD_CACHE")

        echo "Fetching items in '$FROM_STATUS' status..." >&2
        ITEMS=$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | \
            jq -r --arg status "$FROM_STATUS" \
            '.items[] | select(.status == $status) | .id')

        ITEM_COUNT=$(echo "$ITEMS" | wc -l)
        echo "Found $ITEM_COUNT items to move" >&2

        MOVED=0
        while IFS= read -r ITEM_ID; do
            [[ -z "$ITEM_ID" ]] && continue
            echo "[$((++MOVED))/$ITEM_COUNT] Moving $ITEM_ID..." >&2
            "$SCRIPT_DIR/move-card.sh" "$PROJECT_ID" "$ITEM_ID" "$TO_STATUS" "$FIELD_CACHE" >/dev/null
            sleep 0.5  # Rate limiting courtesy
        done <<< "$ITEMS"

        echo "✓ Moved $MOVED items from '$FROM_STATUS' to '$TO_STATUS'" >&2
        ;;

    create-batch)
        OWNER=${2:-}
        PROJECT_NUMBER=${3:-}
        FILE=${4:-}

        if [[ -z "$OWNER" ]] || [[ -z "$PROJECT_NUMBER" ]] || [[ -z "$FILE" ]]; then
            echo "Error: Missing arguments" >&2
            echo "Usage: $0 create-batch <owner> <project-number> <file>" >&2
            exit 1
        fi

        if [[ ! -f "$FILE" ]]; then
            echo "Error: File not found: $FILE" >&2
            exit 1
        fi

        CREATED=0
        while IFS= read -r TITLE; do
            [[ -z "$TITLE" ]] && continue
            echo "Creating: $TITLE" >&2
            gh project item-create "$PROJECT_NUMBER" --owner "$OWNER" \
                --title "$TITLE" >/dev/null
            ((CREATED++))
            sleep 0.5  # Rate limiting courtesy
        done < "$FILE"

        echo "✓ Created $CREATED items" >&2
        ;;

    archive-done)
        OWNER=${2:-}
        PROJECT_NUMBER=${3:-}
        FIELD_CACHE=${4:-fields-cache.json}

        if [[ -z "$OWNER" ]] || [[ -z "$PROJECT_NUMBER" ]]; then
            echo "Error: Missing arguments" >&2
            echo "Usage: $0 archive-done <owner> <project-number> [field-cache]" >&2
            exit 1
        fi

        PROJECT_ID=$(jq -r '.project_id' "$FIELD_CACHE")

        echo "Fetching items in 'Done' status..." >&2
        ITEMS=$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | \
            jq -r '.items[] | select(.status == "Done") | .id')

        ITEM_COUNT=$(echo "$ITEMS" | wc -l)
        echo "Found $ITEM_COUNT items to archive" >&2

        # Check if Archive status exists
        ARCHIVE_EXISTS=$(jq -r '.lookups.status_field.options[] | select(.name == "Archive") | .id' \
            "$FIELD_CACHE")

        if [[ -z "$ARCHIVE_EXISTS" ]]; then
            echo "Warning: No 'Archive' status found. Available statuses:" >&2
            jq -r '.lookups.status_field.options[] | "  - \(.name)"' "$FIELD_CACHE" >&2
            exit 1
        fi

        ARCHIVED=0
        while IFS= read -r ITEM_ID; do
            [[ -z "$ITEM_ID" ]] && continue
            echo "[$((++ARCHIVED))/$ITEM_COUNT] Archiving $ITEM_ID..." >&2
            "$SCRIPT_DIR/move-card.sh" "$PROJECT_ID" "$ITEM_ID" "Archive" "$FIELD_CACHE" >/dev/null
            sleep 0.5  # Rate limiting courtesy
        done <<< "$ITEMS"

        echo "✓ Archived $ARCHIVED items" >&2
        ;;

    list-by-status)
        OWNER=${2:-}
        PROJECT_NUMBER=${3:-}
        STATUS=${4:-}

        if [[ -z "$OWNER" ]] || [[ -z "$PROJECT_NUMBER" ]]; then
            echo "Error: Missing arguments" >&2
            echo "Usage: $0 list-by-status <owner> <project-number> [status]" >&2
            exit 1
        fi

        if [[ -z "$STATUS" ]]; then
            # List all items grouped by status
            gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | \
                jq -r '
                    .items
                    | group_by(.status)
                    | .[]
                    | "\n\(.[0].status):",
                      (.[] | "  [\(.id | split("_")[1])] \(.content.title)")
                '
        else
            # List items for specific status
            gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json | \
                jq -r --arg status "$STATUS" \
                    '.items[]
                    | select(.status == $status)
                    | "[\(.id | split("_")[1])] \(.content.title)"'
        fi
        ;;

    export-csv)
        OWNER=${2:-}
        PROJECT_NUMBER=${3:-}
        OUTPUT_FILE=${4:-}

        if [[ -z "$OWNER" ]] || [[ -z "$PROJECT_NUMBER" ]] || [[ -z "$OUTPUT_FILE" ]]; then
            echo "Error: Missing arguments" >&2
            echo "Usage: $0 export-csv <owner> <project-number> <output-file>" >&2
            exit 1
        fi

        echo "Exporting project to CSV..." >&2

        # Fetch all items
        ITEMS=$(gh project item-list "$PROJECT_NUMBER" --owner "$OWNER" --format json)

        # Convert to CSV
        echo "$ITEMS" | jq -r '
            ["ID", "Type", "Title", "Status", "Labels", "Assignees", "Created", "Updated"],
            (.items[] | [
                .id,
                .type,
                .content.title,
                .status,
                (.labels | map(.name) | join(";")),
                (.assignees | map(.login) | join(";")),
                .createdAt,
                .updatedAt
            ])
            | @csv
        ' > "$OUTPUT_FILE"

        ITEM_COUNT=$(echo "$ITEMS" | jq '.items | length')
        echo "✓ Exported $ITEM_COUNT items to: $OUTPUT_FILE" >&2
        ;;

    *)
        echo "Error: Unknown command: $COMMAND" >&2
        usage
        exit 1
        ;;
esac
