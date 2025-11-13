---
name: github-projects
description: Complete guide for working with GitHub Projects (v2) REST API for kanban board operations. Covers authentication, item management, field operations, status updates, and practical patterns for project automation.
license: MIT
---

# GitHub Projects REST API Integration

Production-ready skill for managing GitHub Projects (v2) using the REST API, optimized for kanban board operations like the tkan project.

## When to Use This Skill

Use this skill when you need to:
- Create, read, update, or delete project items
- Move cards between kanban columns (status changes)
- Discover and work with custom project fields
- Automate project management workflows
- Build integrations that sync with GitHub Projects
- Troubleshoot GitHub API authentication or field mapping issues

## Prerequisites

### Required Tools
- **GitHub CLI (`gh`)** - Installed and authenticated
- **`jq`** - For JSON processing (optional but recommended)

### Authentication & Setup

#### 1. Check Authentication Status
```bash
gh auth status
```

#### 2. Add Project Scope (if missing)
```bash
gh auth refresh -h github.com -s project
```

#### 3. Verify Access
```bash
# List your projects
gh project list --owner @me

# List organization projects
gh project list --owner <org-name>
```

## Core Concepts

### ProjectsV2 Architecture

GitHub Projects v2 uses a property-based model:
- **Projects** contain **Items**
- **Items** have **Fields** with **Values**
- **Status Field** is a special single-select field that maps to kanban columns

### Key IDs You'll Need

Every project operation requires specific IDs:
1. **Project Number** - Visible in URL (e.g., `7` in `github.com/users/GGPrompts/projects/7`)
2. **Project ID** - GraphQL node ID (starts with `PVT_`)
3. **Field ID** - Unique identifier for each field (e.g., `233495315` for Status)
4. **Option ID** - Unique identifier for single-select options (e.g., `f75ad846` for "Todo")
5. **Item ID** - Unique identifier for each project item

## Authentication Methods

### Method 1: GitHub CLI (Recommended)

**Best for:** Quick operations, scripting, development

```bash
# List items
gh project item-list 7 --owner GGPrompts --format json

# Create item
gh project item-create 7 --owner GGPrompts \
  --title "New task" \
  --body "Task description"
```

**Pros:** Automatic authentication, easy syntax, handles both REST and GraphQL
**Cons:** Requires gh CLI installed

### Method 2: Direct REST API

**Best for:** Production integrations, web services

```bash
# Using curl with personal access token
curl -H "Authorization: Bearer ghp_xxxx" \
     -H "Accept: application/vnd.github+json" \
     -H "X-GitHub-Api-Version: 2022-11-28" \
     https://api.github.com/users/GGPrompts/projectsV2/7/items
```

**Pros:** No dependencies, works anywhere
**Cons:** Manual token management

### Method 3: GraphQL API

**Best for:** Complex queries, batch operations, field updates

```bash
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      id
      title
      items(first: 100) {
        nodes {
          id
          fieldValues(first: 20) {
            nodes {
              ... on ProjectV2ItemFieldSingleSelectValue {
                field {
                  ... on ProjectV2SingleSelectField {
                    name
                  }
                }
                name
              }
            }
          }
        }
      }
    }
  }
}
'
```

**Pros:** Powerful, efficient, required for field updates
**Cons:** More complex syntax

## Decision Tree: Which API to Use?

```
Need to...
├─ List items/projects? → gh CLI
├─ Create draft issue? → gh CLI
├─ Get field metadata? → REST API or GraphQL
├─ Update item fields? → GraphQL (required)
└─ Batch operations? → GraphQL
```

## Core Operations

### 1. List All Projects

#### For User
```bash
gh project list --owner @me
gh project list --owner GGPrompts
```

#### For Organization
```bash
gh project list --owner <org-name>
```

### 2. Get Project Details

```bash
# Using gh CLI
gh project view 7 --owner GGPrompts --format json

# Using GraphQL
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      id
      title
      shortDescription
      url
    }
  }
}
'
```

### 3. List All Items in Project

```bash
# Simple list (up to 100 items)
gh project item-list 7 --owner GGPrompts --format json

# With limit
gh project item-list 7 --owner GGPrompts --limit 50 --format json

# Filter by status (using jq)
gh project item-list 7 --owner GGPrompts --format json | \
  jq '.items[] | select(.status == "Todo")'
```

### 4. Create New Item (Draft Issue)

```bash
# Basic creation
gh project item-create 7 --owner GGPrompts \
  --title "Implement feature X" \
  --body "Description of feature X"

# Capture the created item ID
ITEM_ID=$(gh project item-create 7 --owner GGPrompts \
  --title "New task" \
  --body "Details" \
  --format json | jq -r '.id')
```

### 5. Update Item Status (Move Card Between Columns)

This is the most common operation for kanban boards. **Requires GraphQL.**

#### Step-by-step Process

**Step 1:** Get Project ID (only needed once, can cache)
```bash
PROJECT_ID=$(gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      id
    }
  }
}
' | jq -r '.data.user.projectV2.id')
```

**Step 2:** Move card to new status
```bash
# Variables from tkan project
PROJECT_ID="PVT_..."           # From step 1
ITEM_ID="PVTI_..."            # Item to move
FIELD_ID="233495315"          # Status field ID
OPTION_ID="47fc9ee4"          # "In Progress" option ID

gh api graphql -f query='
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
    }
  }
}
'
```

### 6. Update Draft Issue Details

```bash
gh api graphql -f query='
mutation {
  updateProjectV2DraftIssue(input: {
    draftIssueId: "'$ITEM_ID'"
    title: "Updated title"
    body: "Updated description"
  }) {
    draftIssue {
      id
      title
    }
  }
}
'
```

### 7. Delete/Archive Item

```bash
# Delete item from project
gh api graphql -f query='
mutation {
  deleteProjectV2Item(input: {
    projectId: "'$PROJECT_ID'"
    itemId: "'$ITEM_ID'"
  }) {
    deletedItemId
  }
}
'
```

## Field Management

### Understanding Field Types

| Field Type | Description | Value Type | Example Use |
|------------|-------------|------------|-------------|
| `single_select` | One choice from options | Option ID | Status, Priority, Size |
| `text` | Free-form text | String | Notes, Description |
| `number` | Numeric value | Float/Integer | Story Points, Hours |
| `date` | Date value | ISO 8601 string | Due Date, Target |
| `iteration` | Sprint/cycle | Iteration ID | Sprint 1, Q1 2024 |

### Get All Fields for a Project

```bash
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      fields(first: 20) {
        nodes {
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
            }
          }
        }
      }
    }
  }
}
'
```

### Extract Status Field Mapping

```bash
# Get Status field with all options
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      fields(first: 20) {
        nodes {
          ... on ProjectV2SingleSelectField {
            id
            name
            options {
              id
              name
            }
          }
        }
      }
    }
  }
}
' | jq '.data.user.projectV2.fields.nodes[] | select(.name == "Status")'
```

### Cache Field IDs (Best Practice)

Instead of hardcoding field IDs, fetch and cache them:

```bash
# Discovery script (run once, save to config)
cat > discover-fields.sh << 'EOF'
#!/bin/bash
OWNER=$1
PROJECT_NUMBER=$2

gh api graphql -f query='
query {
  user(login: "'$OWNER'") {
    projectV2(number: '$PROJECT_NUMBER') {
      id
      fields(first: 20) {
        nodes {
          ... on ProjectV2SingleSelectField {
            id
            name
            options {
              id
              name
            }
          }
        }
      }
    }
  }
}
' | jq -r '
.data.user.projectV2 | {
  project_id: .id,
  status_field: (.fields.nodes[] | select(.name == "Status") | {
    id,
    options: [.options[] | {name, id}]
  })
}
'
EOF

chmod +x discover-fields.sh
./discover-fields.sh GGPrompts 7
```

## Practical Patterns for Kanban Operations

### Pattern 1: Move Card Between Columns

```bash
#!/bin/bash
# move-card.sh - Move a card to a new status

PROJECT_ID=$1
ITEM_ID=$2
STATUS_NAME=$3  # "Todo", "In Progress", "Done"

# Field IDs for tkan project #7
FIELD_ID="233495315"

# Status option IDs
declare -A STATUS_OPTIONS=(
    ["Todo"]="f75ad846"
    ["In Progress"]="47fc9ee4"
    ["Done"]="98236657"
)

OPTION_ID="${STATUS_OPTIONS[$STATUS_NAME]}"

gh api graphql -f query='
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "'$PROJECT_ID'"
    itemId: "'$ITEM_ID'"
    fieldId: "'$FIELD_ID'"
    value: {singleSelectOptionId: "'$OPTION_ID'"}
  }) {
    projectV2Item { id }
  }
}
'
```

### Pattern 2: Bulk Move Cards

```bash
#!/bin/bash
# bulk-move.sh - Move multiple cards to same status

PROJECT_ID=$1
STATUS=$2
shift 2
ITEM_IDS=("$@")

for ITEM_ID in "${ITEM_IDS[@]}"; do
    echo "Moving $ITEM_ID to $STATUS..."
    ./move-card.sh "$PROJECT_ID" "$ITEM_ID" "$STATUS"
    sleep 0.5  # Rate limiting courtesy
done
```

### Pattern 3: Create Card with Initial Status

```bash
#!/bin/bash
# create-and-move.sh - Create card and set initial status

OWNER=$1
PROJECT_NUMBER=$2
TITLE=$3
STATUS=$4

# Create draft issue
ITEM_ID=$(gh project item-create $PROJECT_NUMBER --owner $OWNER \
  --title "$TITLE" \
  --format json | jq -r '.id')

echo "Created item: $ITEM_ID"

# Get project ID
PROJECT_ID=$(gh api graphql -f query='
query {
  user(login: "'$OWNER'") {
    projectV2(number: '$PROJECT_NUMBER') { id }
  }
}
' | jq -r '.data.user.projectV2.id')

# Move to desired status
./move-card.sh "$PROJECT_ID" "$ITEM_ID" "$STATUS"
```

### Pattern 4: Filter and Display by Status

```bash
#!/bin/bash
# show-column.sh - Display all cards in a specific column

OWNER=$1
PROJECT_NUMBER=$2
STATUS=$3  # "Todo", "In Progress", "Done"

gh project item-list $PROJECT_NUMBER --owner $OWNER --format json | \
jq -r --arg status "$STATUS" '
.items[]
| select(.status == $status)
| "[\(.id)] \(.content.title)"
'
```

## Error Handling

### Common Errors and Solutions

#### Error: "Token has not been granted required scopes"
```bash
# Solution: Refresh authentication with project scope
gh auth refresh -h github.com -s project
```

#### Error: "Project not found"
```bash
# Check owner and number are correct
gh project list --owner GGPrompts

# Verify you have access
gh project view 7 --owner GGPrompts
```

#### Error: "Interface conversion" or "field not found"
```bash
# Solution: Verify field structure
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      fields(first: 20) {
        nodes {
          __typename
          ... on ProjectV2Field {
            id
            name
          }
        }
      }
    }
  }
}
'
```

#### Error: Rate Limiting (403)
```bash
# Check your rate limit status
gh api rate_limit

# Solution: Implement caching, reduce API calls
# GitHub API limit: 5000 requests/hour
```

#### Error: Invalid field ID or option ID
```bash
# Solution: Re-fetch field IDs (they may have changed)
./discover-fields.sh GGPrompts 7
```

## tkan Project Integration

### Current Implementation (backend_github.go)

The tkan project uses GraphQL for most operations:

**File:** `backend_github.go` (lines 1-415)

**Key Functions:**
- `LoadBoard()` - Fetches project items (lines 88-272)
- `MoveCard()` - Updates item status (lines 327-400)
- `getProjectID()` - Placeholder for project ID caching (lines 402-405)
- `getStatusFieldID()` - Placeholder for field ID caching (lines 407-410)
- `getStatusOptionID()` - Placeholder for option ID caching (lines 412-415)

### Known Field IDs for tkan Project #7

```go
// Hardcoded in current implementation (should be dynamic)
const (
    StatusFieldID    = "233495315"
    TodoOptionID     = "f75ad846"
    InProgressID     = "47fc9ee4"
    DoneOptionID     = "98236657"
)
```

### Recommended Improvements

1. **Dynamic Field Discovery** (Priority: High)
   ```go
   // Implement actual field discovery
   func (g *GitHubBackend) getStatusFieldID() string {
       // Fetch and cache from API
       // See references/field-discovery.md
   }
   ```

2. **Project ID Caching** (Priority: High)
   ```go
   // Cache project ID to reduce API calls
   func (g *GitHubBackend) getProjectID() string {
       if g.cachedProjectID != "" {
           return g.cachedProjectID
       }
       // Fetch via GraphQL
   }
   ```

3. **Support Multiple Field Types** (Priority: Medium)
   - Currently only handles Status (single-select)
   - Add support for: Assignees, Labels, Target Date, etc.

4. **Error Recovery** (Priority: Medium)
   - Better handling of API rate limits
   - Graceful degradation when fields missing
   - User-friendly error messages

## Testing Your Integration

### 1. Verify Authentication
```bash
gh auth status
# Should show "✓ Logged in to github.com"
# Should list 'project' in scopes
```

### 2. Test Project Access
```bash
gh project view 7 --owner GGPrompts --format json
# Should return project details without error
```

### 3. Test Item Listing
```bash
gh project item-list 7 --owner GGPrompts --format json
# Should return array of items
```

### 4. Test Field Discovery
```bash
./scripts/discover-fields.sh GGPrompts 7
# Should return field IDs and options
```

### 5. Test Card Movement
```bash
# Create test item
ITEM_ID=$(gh project item-create 7 --owner GGPrompts \
  --title "Test card" \
  --format json | jq -r '.id')

# Move to In Progress
./scripts/move-card.sh <project-id> $ITEM_ID "In Progress"

# Verify (check in web UI or via API)
```

## Performance Best Practices

### 1. Cache Everything Possible
- Project ID (changes rarely)
- Field IDs (changes rarely)
- Option IDs (changes rarely)
- Item list (invalidate on updates)

### 2. Minimize API Calls
```bash
# BAD: Multiple calls for same data
for item in $ITEMS; do
    gh project item-list 7 --owner GGPrompts  # ❌ Repeated
done

# GOOD: Fetch once, process locally
ITEMS=$(gh project item-list 7 --owner GGPrompts --format json)
echo "$ITEMS" | jq '.items[] | ...'  # ✓ Cached
```

### 3. Batch Operations When Possible
GraphQL supports multiple mutations in one request:
```graphql
mutation {
  move1: updateProjectV2ItemFieldValue(input: {...}) { ... }
  move2: updateProjectV2ItemFieldValue(input: {...}) { ... }
  move3: updateProjectV2ItemFieldValue(input: {...}) { ... }
}
```

### 4. Use Pagination for Large Projects
```bash
# For projects with >100 items, use pagination
gh api graphql -f query='
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100, after: "'$CURSOR'") {
        nodes { id }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }
}
'
```

## API Rate Limits

- **REST API**: 5,000 requests/hour (authenticated)
- **GraphQL API**: 5,000 points/hour (queries ~1 point, mutations ~1-10 points)
- **Check your limit**: `gh api rate_limit`

**Strategies:**
- Cache aggressively (project metadata rarely changes)
- Use GraphQL for batch operations (fewer requests)
- Implement exponential backoff on 429 errors
- Consider implementing webhooks for real-time updates

## Reference Files

- `references/field-discovery.md` - Complete field discovery implementation
- `references/graphql-queries.md` - All GraphQL query examples
- `references/error-codes.md` - HTTP status code reference
- `scripts/discover-fields.sh` - Field metadata discovery script
- `scripts/move-card.sh` - Card movement script
- `scripts/bulk-operations.sh` - Batch operation examples

## Additional Resources

### Official Documentation
- [Projects REST API](https://docs.github.com/en/rest/projects/projects?apiVersion=2022-11-28)
- [Project Items API](https://docs.github.com/en/rest/projects/items?apiVersion=2022-11-28)
- [Project Fields API](https://docs.github.com/en/rest/projects/fields?apiVersion=2022-11-28)
- [GraphQL ProjectV2](https://docs.github.com/en/graphql/reference/objects#projectv2)
- [GitHub CLI Manual](https://cli.github.com/manual/)

### Related Skills
- `/tui-add-keybinding` - Add keyboard shortcuts to tkan
- `/tui-dynamic-panel` - Enhance tkan UI layouts
- `/bubbletea` - Comprehensive TUI development

## Quick Reference Card

```
┌─────────────────────────────────────────────────────┐
│ GitHub Projects API Quick Reference                 │
├─────────────────────────────────────────────────────┤
│ List items:                                         │
│   gh project item-list 7 --owner GGPrompts         │
│                                                     │
│ Create item:                                        │
│   gh project item-create 7 --owner GGPrompts \     │
│     --title "Task" --body "Description"            │
│                                                     │
│ Get fields:                                         │
│   gh api /users/GGPrompts/projectsV2/7/fields      │
│                                                     │
│ Move card (GraphQL required):                       │
│   updateProjectV2ItemFieldValue(                   │
│     projectId, itemId, fieldId,                    │
│     value: {singleSelectOptionId}                  │
│   )                                                 │
│                                                     │
│ Status IDs for tkan #7:                            │
│   Field: 233495315                                 │
│   Todo: f75ad846                                   │
│   In Progress: 47fc9ee4                            │
│   Done: 98236657                                   │
└─────────────────────────────────────────────────────┘
```

## Troubleshooting Checklist

- [ ] `gh auth status` shows authenticated
- [ ] Project scope is enabled (`gh auth refresh -h github.com -s project`)
- [ ] Owner and project number are correct
- [ ] Field IDs are current (re-run field discovery)
- [ ] Using GraphQL for field updates (REST API can't update fields)
- [ ] Rate limit not exceeded (`gh api rate_limit`)
- [ ] Item IDs are valid (items not deleted)
- [ ] Proper API version header for REST calls (`X-GitHub-Api-Version: 2022-11-28`)

---

**Last Updated:** 2024-10-28
**Skill Version:** 1.0.0
**Tested With:** GitHub CLI 2.57.0, GitHub API 2022-11-28
