# Classic Columns API vs ProjectsV2 API - Complete Comparison

## Quick Reference Table

| Aspect | Classic Columns (Deprecated) | ProjectsV2 (Current) |
|--------|-----|-----|
| **Status** | Deprecated (removal: Apr 2025) | Active & maintained |
| **Primary Unit** | Column (represents a stage) | Field (stores item properties) |
| **URL Pattern** | `/projects/{id}/columns` | `/users/{user}/projectsV2/{num}` |
| **Organization** | 1D (columns only) | Multi-dimensional (multiple fields) |
| **Field Types** | None (just names) | 14+ types (single_select, text, date, etc.) |
| **Item Location** | Physical position (which column) | Property value (field value on item) |
| **Reordering** | Move endpoint | Field option order |
| **Customization** | Minimal | Extensive |
| **Automation** | Webhooks | GitHub Actions + API |
| **API Stability** | Frozen (will be removed) | Actively developed |

---

## Endpoint Comparison

### 1. Listing Columns/Fields

**Classic Columns**:
```bash
GET /projects/{project_id}/columns
```
Returns array of column objects

**ProjectsV2**:
```bash
GET /users/{username}/projectsV2/{project_number}/fields
GET /orgs/{org}/projectsV2/{project_number}/fields
```
Returns array of field objects

**Key Difference**: ProjectsV2 separates user projects from organization projects

---

### 2. Creating Columns/Fields

**Classic Columns**:
```bash
POST /projects/{project_id}/columns
Content-Type: application/json

{
  "name": "Todo"
}
```

Response:
```json
{
  "id": 367,
  "name": "Todo",
  "url": "https://api.github.com/projects/columns/367",
  "project_url": "https://api.github.com/projects/120",
  "cards_url": "https://api.github.com/projects/columns/367/cards",
  "created_at": "2016-09-05T14:18:44Z",
  "updated_at": "2016-09-05T14:22:28Z"
}
```

**ProjectsV2**:
```bash
POST /users/{username}/projectsV2/{project_number}/fields
Content-Type: application/json

{
  "name": "Status",
  "data_type": "single_select"
}
```

Response:
```json
{
  "id": 123456,
  "name": "Status",
  "data_type": "single_select",
  "created_at": "2022-04-28T12:00:00Z",
  "updated_at": "2022-04-28T12:00:00Z",
  "project_url": "https://api.github.com/users/username/projectsV2/7"
}
```

**Key Differences**:
- ProjectsV2 requires `data_type` specification
- ProjectsV2 fields are more metadata-focused
- ProjectsV2 returns different ID structure

---

### 3. Updating Columns/Fields

**Classic Columns**:
```bash
PATCH /projects/columns/{column_id}
Content-Type: application/json

{
  "name": "In Progress"
}
```
Only supports renaming

**ProjectsV2**:
```bash
PATCH /users/{username}/projectsV2/{project_number}/fields/{field_id}
Content-Type: application/json

{
  "name": "Status",
  "data_type": "single_select"
}
```
Supports full field editing

**Key Differences**:
- ProjectsV2 can update more properties
- Classic is read-only except for name

---

### 4. Deleting Columns/Fields

**Classic Columns**:
```bash
DELETE /projects/columns/{column_id}
```
Returns 204 No Content

**ProjectsV2**:
```bash
DELETE /users/{username}/projectsV2/{project_number}/fields/{field_id}
```
Returns 204 No Content

**Key Difference**: Similar behavior, different URL structure

---

### 5. Reordering

**Classic Columns** (Special feature):
```bash
POST /projects/columns/{column_id}/moves
Content-Type: application/json

{
  "position": "last"
}
```

Options:
- `"first"` - Move to start
- `"last"` - Move to end
- `"after:367"` - After specific column

Returns: 201 Created (empty response)

**ProjectsV2**:
No dedicated reorder endpoint. Instead:
- Update the order of `single_select` options when modifying field
- Items reference field values, not positions

---

## Data Structure Comparison

### Classic Column Object

```json
{
  "url": "https://api.github.com/projects/columns/367",
  "project_url": "https://api.github.com/projects/120",
  "cards_url": "https://api.github.com/projects/columns/367/cards",
  "id": 367,
  "node_id": "MDEzOlByb2plY3RDb2x1bW4zNjc=",
  "name": "To Do",
  "created_at": "2016-09-05T14:18:44Z",
  "updated_at": "2016-09-05T14:22:28Z"
}
```

**Properties**:
- `id`: Integer ID
- `name`: Column name only
- `cards_url`: Direct link to cards endpoint
- Includes REST URLs for navigation
- Timestamp tracking

---

### ProjectsV2 Field Object

```json
{
  "id": 123456789,
  "node_id": "PVT_kwDOANJ1Fs4AZPrL",
  "project_url": "https://api.github.com/users/username/projectsV2/7",
  "name": "Status",
  "data_type": "single_select",
  "options": [
    {
      "id": "f75ad846",
      "name": "Todo"
    },
    {
      "id": "47fc9ee4",
      "name": "In Progress"
    },
    {
      "id": "98236657",
      "name": "Done"
    }
  ],
  "created_at": "2022-04-28T12:00:00Z",
  "updated_at": "2022-04-28T12:00:00Z"
}
```

**Properties**:
- `id`: Large integer ID
- `name`: Field name
- `data_type`: Type specification (crucial difference)
- `options`: Available values (for single_select)
- No cards_url (items don't "belong to" a field)
- Timestamp tracking

---

## Operational Model Differences

### Classic Columns Model (Container-based)

```
Project 120
├── Column 367: "Todo"
│   └── Card 1: Issue #1
│   └── Card 2: Issue #2
├── Column 368: "In Progress"
│   └── Card 3: Issue #3
└── Column 369: "Done"
    └── Card 4: Issue #4

Action: "Move Card 3 from Column 368 to Column 369"
```

**Operations**:
- Add card to column
- Move card between columns
- Remove card from column

---

### ProjectsV2 Model (Property-based)

```
Project 7 (User owned)
├── Field: "Status" (type: single_select)
│   ├── Option: "Todo" (id: f75ad846)
│   ├── Option: "In Progress" (id: 47fc9ee4)
│   └── Option: "Done" (id: 98236657)
├── Field: "Priority" (type: single_select)
│   ├── Option: "Low"
│   ├── Option: "Medium"
│   └── Option: "High"
└── Items
    ├── Item 1 (Issue #1)
    │   ├── Status: "Todo"
    │   └── Priority: "High"
    ├── Item 2 (Issue #2)
    │   ├── Status: "In Progress"
    │   └── Priority: "Medium"
```

**Operations**:
- Set item's Status field value
- Set item's Priority field value
- Query items by field values
- Create/update/delete fields
- Modify field options

---

## Workflow Example Comparison

### Scenario: Move a task to "In Progress"

**Classic Columns Approach**:
1. Get current column of card
2. Call: `POST /projects/columns/{new_column_id}/cards`
3. Card physically moves to new column

```bash
# 1. Find the "In Progress" column
curl https://api.github.com/projects/120/columns

# 2. Move card 456 to that column
curl -X POST https://api.github.com/projects/columns/368/cards \
  -d '{"content_id": 456, "content_type": "Issue"}'
```

**ProjectsV2 Approach**:
1. Get the Status field ID
2. Update item's Status field value

```bash
# 1. Get Status field
curl https://api.github.com/users/myuser/projectsV2/7/fields

# 2. Update item's Status field
curl -X PATCH https://api.github.com/users/myuser/projectsV2/7/items/item_id \
  -d '{
    "field_values": {
      "field_id": 123456789,
      "value": "47fc9ee4"
    }
  }'
```

---

## Feature Parity Matrix

| Feature | Classic | ProjectsV2 | Notes |
|---------|---------|-----------|-------|
| **Basic CRUD** | ✓ | ✓ | Both support create/read/update/delete |
| **Reordering** | ✓ (move) | ✗ | Only Classic has dedicated move endpoint |
| **Single Status** | ✓ | ✓ | Both support status-like fields |
| **Multiple Fields** | ✗ | ✓ | ProjectsV2 exclusive |
| **Field Types** | ✗ (names only) | ✓ (14+ types) | Major V2 advantage |
| **Custom Options** | ✗ | ✓ | Add field options dynamically |
| **Filtering** | Limited | Advanced | V2 supports complex queries |
| **Automation** | Webhooks | Webhooks + Actions | V2 is more automatable |
| **Hierarchy** | ✗ | ✓ | V2 supports parent/child issues |
| **Team Collaboration** | Basic | Advanced | V2 has richer perms |

---

## Migration Path for tkan

### Step 1: Identify What tkan Uses
If using Classic (unlikely given CLAUDE.md), map:
- Columns → Status Field
- Column names → Status options
- Card location → Item field value

### Step 2: Update Backend Implementation
Replace:
```go
// OLD
func GetColumns(projectID int) // GET /projects/{id}/columns

// NEW
func GetFields(projectNumber int) // GET /users/{owner}/projectsV2/{num}/fields
```

### Step 3: Update Data Model
Replace:
```go
// OLD
type Column struct {
    ID    int
    Name  string
    Cards []Card
}

// NEW
type Field struct {
    ID       int
    Name     string
    DataType string
    Options  []Option
}

type Option struct {
    ID   string
    Name string
}
```

### Step 4: Update Item Movement Logic
Replace:
```go
// OLD: Move card between columns
MoveCard(cardID, toColumnID)

// NEW: Update item's field value
UpdateItemField(itemID, fieldID, valueID)
```

---

## Error Handling Comparison

### Classic Columns Errors

```
401: Requires Authentication
403: Forbidden (no access to project)
404: Column not found
422: Validation Failed (invalid position)
```

### ProjectsV2 Errors

```
401: Requires Authentication
403: Forbidden (no access to project)
404: Field/Item not found
422: Validation Failed (invalid data_type)
```

**Similarity**: Both follow standard GitHub API error patterns

---

## Authentication & Permissions

Both APIs support:
- OAuth tokens with project scope
- Personal access tokens
- GitHub App authentication

**Scopes**:
- `read:project` - Read projects
- `write:project` - Create and manage

---

## Deprecation Timeline

### Classic Projects (Columns)
- **Deprecated**: May 23, 2024
- **Removal Date**: April 1, 2025
- **Time Remaining**: ~6 months from Oct 2024
- **Action**: Migrate now

### ProjectsV2
- **Status**: Current standard
- **Future**: Actively maintained
- **Investment**: Safe long-term

---

## Summary for Implementation

1. **Choose ProjectsV2** - It's the future of GitHub Projects
2. **Think in Fields** - Not columns or physical locations
3. **Use Status Field** - For kanban-like workflows
4. **Leverage Multiple Fields** - For richer data models
5. **Plan for Options** - Manage selectable values per field

tkan's current approach using ProjectsV2 in CLAUDE.md is the correct choice for longevity and feature richness.
