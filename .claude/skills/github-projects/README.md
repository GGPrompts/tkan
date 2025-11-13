# GitHub Projects REST API Integration Skill

Complete skill for working with GitHub Projects (v2) REST API in kanban board integrations like tkan.

## Quick Start

### 1. Authenticate
```bash
gh auth refresh -h github.com -s project
```

### 2. Discover Field IDs
```bash
./scripts/discover-fields.sh GGPrompts 7
# Creates fields-cache.json with project metadata
```

### 3. Move a Card
```bash
PROJECT_ID=$(jq -r '.project_id' fields-cache.json)
./scripts/move-card.sh "$PROJECT_ID" "PVTI_..." "In Progress"
```

## Skill Structure

```
github-projects/
├── SKILL.md                    # Main skill documentation
├── README.md                   # This file
├── scripts/
│   ├── discover-fields.sh      # Fetch and cache project field metadata
│   ├── move-card.sh            # Move item between status columns
│   └── bulk-operations.sh      # Batch operations (move-all, create-batch, etc.)
└── references/
    ├── field-discovery.md      # Dynamic field discovery implementation
    ├── graphql-queries.md      # Complete GraphQL query reference
    └── error-codes.md          # HTTP and GraphQL error handling
```

## Core Operations

### View Projects
```bash
gh project list --owner @me
gh project list --owner GGPrompts
```

### List Items
```bash
gh project item-list 7 --owner GGPrompts --format json
```

### Create Item
```bash
gh project item-create 7 --owner GGPrompts \
  --title "New task" \
  --body "Description"
```

### Update Status (requires GraphQL)
```bash
# See scripts/move-card.sh for complete implementation
gh api graphql -f query='
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_..."
    itemId: "PVTI_..."
    fieldId: "PVTF_..."
    value: {singleSelectOptionId: "option-id"}
  }) {
    projectV2Item { id }
  }
}
'
```

## Scripts Reference

### discover-fields.sh
**Purpose:** Fetch project metadata and field IDs

**Usage:**
```bash
./scripts/discover-fields.sh <owner> <project-number> [output-file]
```

**Examples:**
```bash
# Discover fields for tkan project
./scripts/discover-fields.sh GGPrompts 7

# Save to custom file
./scripts/discover-fields.sh GGPrompts 7 myproject.json
```

**Output:** Creates `fields-cache.json` with:
- Project ID
- All field IDs and types
- Status field options (for kanban columns)
- Helpful lookup tables

### move-card.sh
**Purpose:** Move a project item to different status

**Usage:**
```bash
./scripts/move-card.sh <project-id> <item-id> <status-name> [field-cache]
```

**Examples:**
```bash
# Move to In Progress
./scripts/move-card.sh PVT_xxx PVTI_yyy "In Progress"

# Using custom cache file
./scripts/move-card.sh PVT_xxx PVTI_yyy "Done" myproject.json
```

**Prerequisites:** Must run `discover-fields.sh` first

### bulk-operations.sh
**Purpose:** Perform batch operations on multiple items

**Commands:**
```bash
# Move all items from one status to another
./scripts/bulk-operations.sh move-all PVT_xxx "Todo" "In Progress"

# Create multiple items from file
./scripts/bulk-operations.sh create-batch GGPrompts 7 tasks.txt

# Archive all done items
./scripts/bulk-operations.sh archive-done GGPrompts 7

# List items by status
./scripts/bulk-operations.sh list-by-status GGPrompts 7 "In Progress"

# Export to CSV
./scripts/bulk-operations.sh export-csv GGPrompts 7 project.csv
```

## Reference Documentation

### field-discovery.md
Complete implementation guide for dynamic field discovery:
- Fetching field metadata
- Cache management
- Go code examples
- Testing patterns
- Migration from hardcoded IDs

**Use when:** Implementing field discovery in tkan or similar projects

### graphql-queries.md
Complete GraphQL query and mutation reference:
- Project queries
- Item queries with field values
- All mutation types
- Pagination examples
- Error handling
- Performance tips

**Use when:** Writing custom GraphQL operations

### error-codes.md
HTTP status codes and GraphQL error types:
- Common error scenarios
- Debugging steps
- Error handling patterns (Bash & Go)
- Rate limiting guidance

**Use when:** Debugging API issues or implementing error handling

## tkan Integration

### Current Status (backend_github.go)
The tkan project has:
- ✅ GraphQL-based item listing
- ✅ Status field updates (MoveCard)
- ⚠️ Hardcoded field IDs (lines 402-415)
- ❌ No dynamic field discovery

### Recommended Improvements

1. **Implement Dynamic Field Discovery** (High Priority)
   - Replace hardcoded field IDs with runtime discovery
   - See: `references/field-discovery.md`

2. **Add Field Caching** (High Priority)
   - Cache field metadata to reduce API calls
   - Implement TTL-based cache invalidation

3. **Support Additional Field Types** (Medium Priority)
   - Currently only handles Status (single-select)
   - Add: Assignees, Labels, Target Date, etc.

4. **Improve Error Handling** (Medium Priority)
   - Better rate limit handling
   - User-friendly error messages
   - Graceful degradation

### Quick Integration Example

```go
// In backend_github.go

// Add field discovery
func (g *GitHubBackend) LoadBoard() (*Board, error) {
    // Fetch field metadata first (with caching)
    if err := g.ensureFieldMetadata(); err != nil {
        return nil, err
    }

    // Continue with existing item loading...
}

// Replace placeholder methods
func (g *GitHubBackend) getStatusFieldID() (string, error) {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return "", err
    }
    return metadata.StatusField.ID, nil
}
```

See `references/field-discovery.md` for complete implementation.

## Common Workflows

### Setup New Project Integration

1. **Authenticate:**
   ```bash
   gh auth refresh -h github.com -s project
   ```

2. **Find your project number:**
   ```bash
   gh project list --owner @me
   # Or check URL: github.com/users/YOU/projects/NUMBER
   ```

3. **Discover fields:**
   ```bash
   ./scripts/discover-fields.sh YOUR_OWNER PROJECT_NUMBER
   ```

4. **Test access:**
   ```bash
   gh project item-list PROJECT_NUMBER --owner YOUR_OWNER
   ```

### Move Multiple Cards

```bash
# Get project ID
PROJECT_ID=$(jq -r '.project_id' fields-cache.json)

# Get all "Todo" items
ITEMS=$(gh project item-list 7 --owner GGPrompts --format json | \
  jq -r '.items[] | select(.status == "Todo") | .id')

# Move each to "In Progress"
for ITEM_ID in $ITEMS; do
  ./scripts/move-card.sh "$PROJECT_ID" "$ITEM_ID" "In Progress"
  sleep 0.5  # Rate limiting courtesy
done
```

### Create Cards from List

```bash
# Create tasks.txt with one task per line
cat > tasks.txt << EOF
Implement user authentication
Add error logging
Write unit tests
Update documentation
EOF

# Bulk create
./scripts/bulk-operations.sh create-batch GGPrompts 7 tasks.txt
```

### Export for Reporting

```bash
# Export to CSV
./scripts/bulk-operations.sh export-csv GGPrompts 7 report.csv

# Open in spreadsheet
libreoffice report.csv  # or `open report.csv` on macOS
```

## Troubleshooting

### "Token has not been granted required scopes"
```bash
gh auth refresh -h github.com -s project
gh auth status  # Verify "project" is in scopes
```

### "Project not found"
```bash
# Verify project number
gh project list --owner GGPrompts

# Check owner name (case-sensitive)
# Try organization if user fails
gh project list --owner myorg
```

### "Field does not exist on this project"
```bash
# Field IDs changed, re-discover
./scripts/discover-fields.sh GGPrompts 7

# Check new field IDs
jq '.fields[] | {name, id}' fields-cache.json
```

### Rate limit exceeded
```bash
# Check remaining calls
gh api rate_limit | jq '.resources.graphql'

# Wait for reset or implement caching
```

## API Limits

- **REST API:** 5,000 requests/hour
- **GraphQL API:** 5,000 points/hour
- **Check limit:** `gh api rate_limit`

**Best Practices:**
- Cache field metadata (TTL: 5 minutes)
- Use batch mutations for multiple updates
- Implement exponential backoff on errors

## Resources

### Documentation
- [Projects REST API](https://docs.github.com/en/rest/projects/projects)
- [Project Items API](https://docs.github.com/en/rest/projects/items)
- [Project Fields API](https://docs.github.com/en/rest/projects/fields)
- [GraphQL ProjectV2](https://docs.github.com/en/graphql/reference/objects#projectv2)

### Related Skills
- `/tui-add-keybinding` - Add keyboard shortcuts to tkan
- `/tui-dynamic-panel` - Enhance tkan UI layouts
- `/bubbletea` - Comprehensive TUI development

### GitHub
- [tkan Repository](https://github.com/GGPrompts/tkan)
- [GitHub CLI Manual](https://cli.github.com/manual/)

## Support

For issues with:
- **This skill:** Open issue in tkan repository
- **GitHub API:** Check [GitHub Status](https://www.githubstatus.com/)
- **GitHub CLI:** Run `gh help` or visit [cli.github.com](https://cli.github.com/)

## License

MIT License - See project root for details

---

**Last Updated:** 2024-10-28
**Skill Version:** 1.0.0
**Tested With:** GitHub CLI 2.57.0, GitHub API 2022-11-28
