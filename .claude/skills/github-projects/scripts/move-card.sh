#!/bin/bash
# move-card.sh - Move a project item to a different status
# Usage: ./move-card.sh <project-id> <item-id> <status-name> [field-cache]

set -euo pipefail

PROJECT_ID=${1:-}
ITEM_ID=${2:-}
STATUS_NAME=${3:-}
FIELD_CACHE=${4:-fields-cache.json}

if [[ -z "$PROJECT_ID" ]] || [[ -z "$ITEM_ID" ]] || [[ -z "$STATUS_NAME" ]]; then
    echo "Usage: $0 <project-id> <item-id> <status-name> [field-cache]"
    echo ""
    echo "Arguments:"
    echo "  project-id   - GraphQL project ID (starts with PVT_)"
    echo "  item-id      - GraphQL item ID (starts with PVTI_)"
    echo "  status-name  - Target status name (e.g., 'Todo', 'In Progress', 'Done')"
    echo "  field-cache  - Optional: path to field cache JSON (default: fields-cache.json)"
    echo ""
    echo "Examples:"
    echo "  $0 PVT_xxx PVTI_yyy 'In Progress'"
    echo "  $0 PVT_xxx PVTI_yyy 'Done' myproject-fields.json"
    echo ""
    echo "Run discover-fields.sh first to create the field cache."
    exit 1
fi

# Check if field cache exists
if [[ ! -f "$FIELD_CACHE" ]]; then
    echo "Error: Field cache not found: $FIELD_CACHE" >&2
    echo "Run discover-fields.sh first to create the cache." >&2
    exit 1
fi

# Get status field ID and option ID from cache
FIELD_ID=$(jq -r '.lookups.status_field.id // empty' "$FIELD_CACHE")
OPTION_ID=$(jq -r --arg status "$STATUS_NAME" \
    '.lookups.status_field.options[] | select(.name == $status) | .id // empty' \
    "$FIELD_CACHE")

if [[ -z "$FIELD_ID" ]]; then
    echo "Error: Status field not found in cache" >&2
    echo "Available fields:" >&2
    jq -r '.fields[] | "  - \(.name) (ID: \(.id))"' "$FIELD_CACHE" >&2
    exit 1
fi

if [[ -z "$OPTION_ID" ]]; then
    echo "Error: Status option '$STATUS_NAME' not found" >&2
    echo "Available options:" >&2
    jq -r '.lookups.status_field.options[] | "  - \(.name)"' "$FIELD_CACHE" >&2
    exit 1
fi

echo "Moving item to '$STATUS_NAME'..." >&2
echo "  Project ID: $PROJECT_ID" >&2
echo "  Item ID: $ITEM_ID" >&2
echo "  Field ID: $FIELD_ID" >&2
echo "  Option ID: $OPTION_ID" >&2

# Execute the mutation
MUTATION='
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "'$PROJECT_ID'"
    itemId: "'$ITEM_ID'"
    fieldId: "'$FIELD_ID'"
    value: {
      singleSelectOptionId: "'$OPTION_ID'"
    }
  }) {
    projectV2Item {
      id
      fieldValueByName(name: "Status") {
        ... on ProjectV2ItemFieldSingleSelectValue {
          name
        }
      }
    }
  }
}
'

RESULT=$(gh api graphql -f query="$MUTATION" 2>&1) || {
    echo "Error: Failed to update item" >&2
    echo "$RESULT" >&2
    exit 1
}

# Check for errors
if echo "$RESULT" | jq -e '.errors' >/dev/null 2>&1; then
    echo "GraphQL Errors:" >&2
    echo "$RESULT" | jq -r '.errors[] | "  - \(.message)"' >&2
    exit 1
fi

# Extract new status
NEW_STATUS=$(echo "$RESULT" | jq -r \
    '.data.updateProjectV2ItemFieldValue.projectV2Item.fieldValueByName.name // "Unknown"')

echo "✓ Item moved successfully to: $NEW_STATUS" >&2

# Output result for scripting
echo "$RESULT"
