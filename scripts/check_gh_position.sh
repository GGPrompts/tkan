#!/bin/bash
# Check if GitHub Projects GraphQL API exposes position/order fields
# This uses the GraphQL API directly to see raw data

echo "=== Checking GitHub Projects Position Fields via GraphQL ==="
echo ""
echo "Usage: ./scripts/check_gh_position.sh OWNER PROJECT_NUMBER"
echo ""

if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Example: ./scripts/check_gh_position.sh GGPrompts 7"
    echo ""
    echo "Please provide owner and project number as arguments"
    exit 1
fi

OWNER=$1
PROJECT_NUM=$2

# First, get the project ID
echo "Step 1: Getting project ID..."
PROJECT_ID=$(gh api graphql -f query="
{
  user(login: \"$OWNER\") {
    projectV2(number: $PROJECT_NUM) {
      id
    }
  }
}" 2>/dev/null | jq -r '.data.user.projectV2.id' || \
gh api graphql -f query="
{
  organization(login: \"$OWNER\") {
    projectV2(number: $PROJECT_NUM) {
      id
    }
  }
}" 2>/dev/null | jq -r '.data.organization.projectV2.id')

if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "null" ]; then
    echo "Error: Could not find project ID"
    exit 1
fi

echo "Project ID: $PROJECT_ID"
echo ""

# Now query for items with all available fields
echo "Step 2: Querying items for position/order information..."
gh api graphql -f query="
{
  node(id: \"$PROJECT_ID\") {
    ... on ProjectV2 {
      items(first: 3) {
        nodes {
          id
          content {
            ... on Issue {
              title
              number
            }
            ... on PullRequest {
              title
              number
            }
            ... on DraftIssue {
              title
            }
          }
          fieldValues(first: 10) {
            nodes {
              ... on ProjectV2ItemFieldTextValue {
                text
                field {
                  ... on ProjectV2Field {
                    name
                  }
                }
              }
              ... on ProjectV2ItemFieldNumberValue {
                number
                field {
                  ... on ProjectV2Field {
                    name
                  }
                }
              }
              ... on ProjectV2ItemFieldSingleSelectValue {
                name
                field {
                  ... on ProjectV2SingleSelectField {
                    name
                  }
                }
              }
            }
          }
        }
      }
    }
  }
}" | jq .

echo ""
echo "=== Analysis ==="
echo "Look for fields like:"
echo "  - position"
echo "  - order"
echo "  - rank"
echo "  - priority (custom field)"
echo "  - Any number field that could represent order"
