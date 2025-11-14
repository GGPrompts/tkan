#!/bin/bash
# Script to check what fields GitHub Projects returns
# This helps us understand if position/order information is available

echo "=== Checking GitHub Projects Fields ==="
echo ""
echo "Usage: ./scripts/check_gh_fields.sh OWNER PROJECT_NUMBER"
echo ""

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Example: ./scripts/check_gh_fields.sh GGPrompts 7"
    echo ""
    echo "Please provide owner and project number as arguments"
    exit 1
fi

OWNER=$1
PROJECT_NUM=$2

echo "Fetching first item from project $PROJECT_NUM (owner: $OWNER)..."
echo ""

# Get one item and show all available fields
gh project item-list "$PROJECT_NUM" --owner "$OWNER" --format json --limit 1 | jq '
    .items[0] |
    {
        "Top-level keys": keys,
        "Content keys": (.content | keys // []),
        "FieldValues keys": (.fieldValues | keys // []),
        "Full item": .
    }
'

echo ""
echo "=== Looking for position/order fields ==="
gh project item-list "$PROJECT_NUM" --owner "$OWNER" --format json --limit 5 | jq '.items[] | {id, title: .content.title, status: .fieldValues.Status}'
