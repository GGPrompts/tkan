#!/bin/bash
# discover-fields.sh - Discover and cache project field metadata
# Usage: ./discover-fields.sh <owner> <project-number> [output-file]

set -euo pipefail

OWNER=${1:-}
PROJECT_NUMBER=${2:-}
OUTPUT_FILE=${3:-fields-cache.json}

if [[ -z "$OWNER" ]] || [[ -z "$PROJECT_NUMBER" ]]; then
    echo "Usage: $0 <owner> <project-number> [output-file]"
    echo ""
    echo "Examples:"
    echo "  $0 GGPrompts 7"
    echo "  $0 myorg 123 myproject-fields.json"
    exit 1
fi

echo "Discovering fields for $OWNER/Project#$PROJECT_NUMBER..." >&2

# Try user context first, fall back to organization
QUERY_USER='
query {
  user(login: "'$OWNER'") {
    projectV2(number: '$PROJECT_NUMBER') {
      id
      title
      fields(first: 50) {
        nodes {
          __typename
          ... on ProjectV2Field {
            id
            name
            dataType
          }
          ... on ProjectV2SingleSelectField {
            id
            name
            dataType
            options {
              id
              name
              color
              description
            }
          }
          ... on ProjectV2IterationField {
            id
            name
            dataType
            configuration {
              iterations {
                id
                title
                startDate
                duration
              }
            }
          }
        }
      }
    }
  }
}
'

QUERY_ORG='
query {
  organization(login: "'$OWNER'") {
    projectV2(number: '$PROJECT_NUMBER') {
      id
      title
      fields(first: 50) {
        nodes {
          __typename
          ... on ProjectV2Field {
            id
            name
            dataType
          }
          ... on ProjectV2SingleSelectField {
            id
            name
            dataType
            options {
              id
              name
              color
              description
            }
          }
          ... on ProjectV2IterationField {
            id
            name
            dataType
            configuration {
              iterations {
                id
                title
                startDate
                duration
              }
            }
          }
        }
      }
    }
  }
}
'

# Try user first
RESULT=$(gh api graphql -f query="$QUERY_USER" 2>/dev/null) || {
    echo "User context failed, trying organization..." >&2
    RESULT=$(gh api graphql -f query="$QUERY_ORG" 2>/dev/null) || {
        echo "Error: Could not fetch project fields" >&2
        echo "Verify that:" >&2
        echo "  1. You have access to the project" >&2
        echo "  2. The project number is correct" >&2
        echo "  3. You have the 'project' scope: gh auth refresh -h github.com -s project" >&2
        exit 1
    }
}

# Extract project data
PROJECT_DATA=$(echo "$RESULT" | jq -r '
  .data.user.projectV2 // .data.organization.projectV2
')

if [[ "$PROJECT_DATA" == "null" ]]; then
    echo "Error: Project not found" >&2
    exit 1
fi

# Process and format the data
FORMATTED=$(echo "$PROJECT_DATA" | jq '{
  project_id: .id,
  project_title: .title,
  owner: "'$OWNER'",
  project_number: '$PROJECT_NUMBER',
  discovered_at: now | strftime("%Y-%m-%dT%H:%M:%SZ"),
  fields: [
    .fields.nodes[] | {
      id,
      name,
      type: .dataType,
      typename: .__typename,
      options: (
        if .options then
          [.options[] | {
            id,
            name,
            color,
            description
          }]
        else
          null
        end
      ),
      iterations: (
        if .configuration.iterations then
          [.configuration.iterations[] | {
            id,
            title,
            startDate,
            duration
          }]
        else
          null
        end
      )
    }
  ],
  # Helpful lookups
  lookups: {
    status_field: (
      .fields.nodes[]
      | select(.name == "Status")
      | {
          id,
          options: [.options[] | {name, id}]
        }
    ),
    field_by_name: (
      .fields.nodes
      | map({(.name): {id, type: .dataType}})
      | add
    )
  }
}')

# Save to file
echo "$FORMATTED" > "$OUTPUT_FILE"

echo "✓ Fields discovered and saved to: $OUTPUT_FILE" >&2
echo "" >&2

# Display summary
echo "Project: $(echo "$FORMATTED" | jq -r '.project_title')" >&2
echo "Project ID: $(echo "$FORMATTED" | jq -r '.project_id')" >&2
echo "Fields found: $(echo "$FORMATTED" | jq '.fields | length')" >&2
echo "" >&2

# Show Status field if present
if echo "$FORMATTED" | jq -e '.lookups.status_field' >/dev/null 2>&1; then
    echo "Status Field:" >&2
    echo "  ID: $(echo "$FORMATTED" | jq -r '.lookups.status_field.id')" >&2
    echo "  Options:" >&2
    echo "$FORMATTED" | jq -r '.lookups.status_field.options[] | "    - \(.name): \(.id)"' >&2
fi

# Output the JSON to stdout for piping
cat "$OUTPUT_FILE"
