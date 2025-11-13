# Field Discovery Implementation Guide

Complete implementation guide for dynamic field discovery in GitHub Projects integrations.

## Problem Statement

Hardcoding field IDs makes your integration brittle:
- Field IDs change if project is recreated
- Different projects have different field IDs
- Custom fields can't be supported
- Option IDs for status values are project-specific

## Solution: Dynamic Field Discovery

Fetch field metadata at runtime and cache it for performance.

## Implementation Pattern

### 1. Fetch Field Metadata

```go
type FieldMetadata struct {
    ProjectID    string
    Fields       map[string]Field  // field name -> Field
    StatusField  *SingleSelectField
    CachedAt     time.Time
}

type Field struct {
    ID       string
    Name     string
    Type     string
    Options  []FieldOption  // for single_select fields
}

type FieldOption struct {
    ID    string
    Name  string
    Color string
}

func (g *GitHubBackend) fetchFieldMetadata() (*FieldMetadata, error) {
    query := `
    query($owner: String!, $number: Int!) {
      user(login: $owner) {
        projectV2(number: $number) {
          id
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
                }
              }
            }
          }
        }
      }
    }
    `

    // Execute GraphQL query
    result, err := g.executeGraphQL(query, map[string]interface{}{
        "owner":  g.owner,
        "number": g.projectNumber,
    })
    if err != nil {
        return nil, err
    }

    // Parse response
    metadata := &FieldMetadata{
        Fields:   make(map[string]Field),
        CachedAt: time.Now(),
    }

    project := result.Data.User.ProjectV2
    metadata.ProjectID = project.ID

    for _, fieldNode := range project.Fields.Nodes {
        field := Field{
            ID:   fieldNode.ID,
            Name: fieldNode.Name,
            Type: fieldNode.DataType,
        }

        // Handle single-select fields
        if fieldNode.Typename == "ProjectV2SingleSelectField" {
            for _, opt := range fieldNode.Options {
                field.Options = append(field.Options, FieldOption{
                    ID:    opt.ID,
                    Name:  opt.Name,
                    Color: opt.Color,
                })
            }

            // Save Status field reference
            if field.Name == "Status" {
                metadata.StatusField = &field
            }
        }

        metadata.Fields[field.Name] = field
    }

    return metadata, nil
}
```

### 2. Cache Management

```go
type GitHubBackend struct {
    // ... existing fields
    fieldMetadata *FieldMetadata
    metadataLock  sync.RWMutex
    cacheTTL      time.Duration
}

func (g *GitHubBackend) getFieldMetadata() (*FieldMetadata, error) {
    g.metadataLock.RLock()

    // Check if cache is valid
    if g.fieldMetadata != nil {
        age := time.Since(g.fieldMetadata.CachedAt)
        if age < g.cacheTTL {
            defer g.metadataLock.RUnlock()
            return g.fieldMetadata, nil
        }
    }
    g.metadataLock.RUnlock()

    // Cache expired or doesn't exist, fetch new data
    g.metadataLock.Lock()
    defer g.metadataLock.Unlock()

    // Double-check after acquiring write lock
    if g.fieldMetadata != nil && time.Since(g.fieldMetadata.CachedAt) < g.cacheTTL {
        return g.fieldMetadata, nil
    }

    // Fetch fresh metadata
    metadata, err := g.fetchFieldMetadata()
    if err != nil {
        return nil, err
    }

    g.fieldMetadata = metadata
    return metadata, nil
}

func (g *GitHubBackend) invalidateFieldCache() {
    g.metadataLock.Lock()
    defer g.metadataLock.Unlock()
    g.fieldMetadata = nil
}
```

### 3. Replace Hardcoded Methods

Update the placeholder methods in `backend_github.go`:

```go
// Before (lines 402-415):
func (g *GitHubBackend) getProjectID() string {
    // This would need to be fetched and cached
    return ""
}

func (g *GitHubBackend) getStatusFieldID() string {
    // This would need to be fetched from project fields
    return ""
}

func (g *GitHubBackend) getStatusOptionID(status string) string {
    // This would need to be fetched from field options
    return ""
}

// After:
func (g *GitHubBackend) getProjectID() (string, error) {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return "", err
    }
    return metadata.ProjectID, nil
}

func (g *GitHubBackend) getStatusFieldID() (string, error) {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return "", err
    }
    if metadata.StatusField == nil {
        return "", fmt.Errorf("status field not found")
    }
    return metadata.StatusField.ID, nil
}

func (g *GitHubBackend) getStatusOptionID(status string) (string, error) {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return "", err
    }

    if metadata.StatusField == nil {
        return "", fmt.Errorf("status field not found")
    }

    for _, opt := range metadata.StatusField.Options {
        if opt.Name == status {
            return opt.ID, nil
        }
    }

    return "", fmt.Errorf("status option '%s' not found", status)
}
```

### 4. Update MoveCard to Use Dynamic Fields

```go
func (g *GitHubBackend) MoveCard(cardID string, fromColumn string, toColumn string) error {
    // Get dynamic field IDs
    projectID, err := g.getProjectID()
    if err != nil {
        return fmt.Errorf("failed to get project ID: %w", err)
    }

    fieldID, err := g.getStatusFieldID()
    if err != nil {
        return fmt.Errorf("failed to get status field ID: %w", err)
    }

    optionID, err := g.getStatusOptionID(toColumn)
    if err != nil {
        return fmt.Errorf("failed to get status option ID for '%s': %w", toColumn, err)
    }

    // Execute mutation (same as before, but with dynamic IDs)
    mutation := fmt.Sprintf(`
        mutation {
            updateProjectV2ItemFieldValue(input: {
                projectId: "%s"
                itemId: "%s"
                fieldId: "%s"
                value: {
                    singleSelectOptionId: "%s"
                }
            }) {
                projectV2Item {
                    id
                }
            }
        }
    `, projectID, cardID, fieldID, optionID)

    _, err = g.executeGraphQL(mutation, nil)
    return err
}
```

## File-Based Caching (Alternative)

For CLI applications that restart frequently, consider file-based caching:

```go
func (g *GitHubBackend) loadFieldCacheFromFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return err
    }

    var metadata FieldMetadata
    if err := json.Unmarshal(data, &metadata); err != nil {
        return err
    }

    // Check if cache is stale
    age := time.Since(metadata.CachedAt)
    if age > g.cacheTTL {
        return fmt.Errorf("cache is stale")
    }

    g.fieldMetadata = &metadata
    return nil
}

func (g *GitHubBackend) saveFieldCacheToFile(path string) error {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return err
    }

    data, err := json.MarshalIndent(metadata, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}
```

## Status Mapping

Map between tkan's column names and GitHub Projects status names:

```go
var statusMapping = map[string]string{
    "BACKLOG":     "Todo",
    "TODO":        "Todo",
    "IN_PROGRESS": "In Progress",
    "DONE":        "Done",
}

func (g *GitHubBackend) mapColumnToStatus(column string) string {
    if status, ok := statusMapping[column]; ok {
        return status
    }
    return column  // Use as-is if no mapping
}

func (g *GitHubBackend) mapStatusToColumn(status string) string {
    for column, mappedStatus := range statusMapping {
        if mappedStatus == status {
            return column
        }
    }
    return "BACKLOG"  // Default
}
```

## Error Handling

```go
func (g *GitHubBackend) getStatusOptionID(status string) (string, error) {
    metadata, err := g.getFieldMetadata()
    if err != nil {
        return "", err
    }

    if metadata.StatusField == nil {
        // Fallback: try to fetch fields again
        g.invalidateFieldCache()
        metadata, err = g.getFieldMetadata()
        if err != nil {
            return "", fmt.Errorf("status field not found and refetch failed: %w", err)
        }
    }

    for _, opt := range metadata.StatusField.Options {
        if opt.Name == status {
            return opt.ID, nil
        }
    }

    // Provide helpful error message
    available := make([]string, len(metadata.StatusField.Options))
    for i, opt := range metadata.StatusField.Options {
        available[i] = opt.Name
    }

    return "", fmt.Errorf(
        "status option '%s' not found. Available: %s",
        status,
        strings.Join(available, ", "),
    )
}
```

## Testing

```go
func TestFieldDiscovery(t *testing.T) {
    backend := NewGitHubBackend("GGPrompts", "7")

    // Test field metadata fetching
    metadata, err := backend.getFieldMetadata()
    require.NoError(t, err)
    require.NotNil(t, metadata)
    require.NotEmpty(t, metadata.ProjectID)

    // Test status field exists
    require.NotNil(t, metadata.StatusField)
    require.Equal(t, "Status", metadata.StatusField.Name)
    require.NotEmpty(t, metadata.StatusField.Options)

    // Test status option lookup
    optionID, err := backend.getStatusOptionID("In Progress")
    require.NoError(t, err)
    require.NotEmpty(t, optionID)

    // Test invalid option
    _, err = backend.getStatusOptionID("Invalid Status")
    require.Error(t, err)
}

func TestFieldCaching(t *testing.T) {
    backend := NewGitHubBackend("GGPrompts", "7")
    backend.cacheTTL = 1 * time.Second

    // First fetch
    meta1, err := backend.getFieldMetadata()
    require.NoError(t, err)
    time1 := meta1.CachedAt

    // Second fetch (should use cache)
    meta2, err := backend.getFieldMetadata()
    require.NoError(t, err)
    require.Equal(t, time1, meta2.CachedAt)

    // Wait for cache to expire
    time.Sleep(2 * time.Second)

    // Third fetch (should refetch)
    meta3, err := backend.getFieldMetadata()
    require.NoError(t, err)
    require.NotEqual(t, time1, meta3.CachedAt)
}
```

## Performance Considerations

### Cache TTL Guidelines
- **Development**: 30 seconds (fast iteration)
- **Production CLI**: 5 minutes (balance freshness and performance)
- **Long-running services**: 1 hour (project fields rarely change)

### API Call Reduction
With dynamic field discovery:
- **Before**: 1 API call per card move (mutation only)
- **After (first move)**: 2 API calls (field discovery + mutation)
- **After (cached)**: 1 API call (mutation only)

Net result: Almost no overhead after initial discovery.

## Migration Path

1. **Add field discovery** without removing hardcoded values
2. **Log comparison** between hardcoded and discovered values
3. **Verify correctness** in development
4. **Switch to discovered values** as primary
5. **Remove hardcoded values** after validation period

```go
func (g *GitHubBackend) getStatusFieldID() (string, error) {
    const HARDCODED_ID = "233495315"  // Temporary fallback

    metadata, err := g.getFieldMetadata()
    if err != nil {
        log.Printf("Warning: Field discovery failed, using hardcoded ID: %v", err)
        return HARDCODED_ID, nil
    }

    if metadata.StatusField == nil {
        log.Printf("Warning: Status field not found, using hardcoded ID")
        return HARDCODED_ID, nil
    }

    discoveredID := metadata.StatusField.ID

    // Validation during migration
    if discoveredID != HARDCODED_ID {
        log.Printf("WARNING: Discovered ID (%s) differs from hardcoded (%s)",
            discoveredID, HARDCODED_ID)
    }

    return discoveredID, nil
}
```

## See Also

- `../SKILL.md` - Main skill documentation
- `graphql-queries.md` - GraphQL query examples
- `../scripts/discover-fields.sh` - Command-line field discovery tool
