# GitHub Project Columns REST API - Comprehensive Summary

## Overview

The GitHub Project Columns REST API is part of the **Projects (Classic)** API and is **deprecated** as of May 23, 2024 with removal scheduled for April 1, 2025. GitHub is transitioning users to the new Projects experience (ProjectsV2), which uses **Fields** instead of **Columns**.

### Deprecation Notice
- **Deprecation Date**: May 23, 2024
- **Removal Date**: April 1, 2025
- **Replacement**: Use ProjectsV2 API with Fields instead
- **Details**: See [GitHub Blog](https://github.blog/changelog/2024-05-23-sunset-notice-projects-classic/)

---

## Column Management Operations

### 1. List Project Columns
**Endpoint**: `GET /projects/{project_id}/columns`

Retrieves all columns in a classic project.

**Parameters**:
- `project_id` (required): The project ID
- `per_page` (optional): Results per page (1-100, default: 30)
- `page` (optional): Page number for pagination

**Response**: Array of Project Column objects
- Status: 200 OK
- Includes Link header for pagination

**Possible Status Codes**:
- 200: Success
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden
- 404: Not Found

---

### 2. Create a Project Column
**Endpoint**: `POST /projects/{project_id}/columns`

Creates a new column in a project.

**Parameters**:
- `project_id` (required): The project ID

**Request Body** (required):
```json
{
  "name": "Remaining tasks"
}
```

**Response**: Project Column object
- Status: 201 Created
- Returns the newly created column with ID and metadata

**Possible Status Codes**:
- 201: Created
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden
- 422: Validation Failed

---

### 3. Get a Project Column
**Endpoint**: `GET /projects/columns/{column_id}`

Retrieves details about a specific column.

**Parameters**:
- `column_id` (required): The column ID

**Response**: Single Project Column object

**Possible Status Codes**:
- 200: Success
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden
- 404: Not Found

---

### 4. Update a Project Column
**Endpoint**: `PATCH /projects/columns/{column_id}`

Renames an existing column.

**Parameters**:
- `column_id` (required): The column ID

**Request Body** (required):
```json
{
  "name": "To Do"
}
```

**Response**: Updated Project Column object

**Possible Status Codes**:
- 200: Success
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden

---

### 5. Delete a Project Column
**Endpoint**: `DELETE /projects/columns/{column_id}`

Removes a column from a project.

**Parameters**:
- `column_id` (required): The column ID

**Response**: No content returned

**Possible Status Codes**:
- 204: No Content (success)
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden

---

### 6. Move a Project Column
**Endpoint**: `POST /projects/columns/{column_id}/moves`

Changes the position of a column within a project board.

**Parameters**:
- `column_id` (required): The column ID

**Request Body** (required):
```json
{
  "position": "last"
}
```

**Position Values**:
- `"first"` - Move to the beginning
- `"last"` - Move to the end
- `"after:<column_id>"` - Place after a specific column (e.g., `"after:367"`)

**Response**: Empty object `{}`
- Status: 201 Created

**Possible Status Codes**:
- 201: Created
- 304: Not Modified
- 401: Requires Authentication
- 403: Forbidden
- 422: Validation Failed

---

## Project Column Data Structure

### Project Column Schema

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

**Field Descriptions**:
- `id`: Unique integer identifier for the column
- `node_id`: GraphQL node identifier (Base64 encoded)
- `name`: Display name of the column (string)
- `url`: REST API URL for this specific column
- `project_url`: REST API URL for the parent project
- `cards_url`: REST API URL to manage cards in this column
- `created_at`: ISO 8601 timestamp of creation
- `updated_at`: ISO 8601 timestamp of last modification

---

## Classic Projects vs ProjectsV2: Key Differences

### Classic Projects (DEPRECATED)

**Structure**:
- Uses **Columns** as the primary organizational unit
- Each column represents a workflow stage (e.g., "To Do", "In Progress", "Done")
- Simple, linear structure
- Cards (issues/PRs) are added directly to columns

**Characteristics**:
- One-dimensional workflow (columns only)
- Limited customization
- Simpler API
- Fixed structure
- End-of-life: April 1, 2025

**API Endpoints**:
- `/projects/{project_id}/columns` - List/Create columns
- `/projects/columns/{column_id}` - Get/Update/Delete single column
- `/projects/columns/{column_id}/moves` - Reorder columns
- `/projects/columns/{column_id}/cards` - Manage cards in column

---

### ProjectsV2 (CURRENT STANDARD)

**Structure**:
- Uses **Fields** instead of columns (more flexible than classic columns)
- Fields can be of multiple types: status, single select, text, number, date, etc.
- Multi-dimensional organization (can have multiple fields for different purposes)
- Items (cards) have properties for each field

**Characteristics**:
- Flexible, customizable fields
- Multiple field types available
- Supports custom field values
- More powerful filtering and sorting
- Supported for new projects
- Active development and improvements

**Field Types Available**:
- `assignees` - Team member assignments
- `labels` - GitHub labels
- `milestone` - Release milestones
- `repository` - Source repository
- `title` - Item title
- `text` - Free-form text
- `single_select` - Dropdown options (like status)
- `number` - Numeric values
- `date` - Date fields
- `iteration` - Sprint/iteration tracking
- `issue_type` - Type classification
- `parent_issue` - Hierarchy support
- `sub_issues_progress` - Sub-task tracking
- `reviewers` - Pull request reviewers
- `linked_pull_requests` - Related PRs

**API Endpoints**:
- `/users/{username}/projectsV2/{project_number}` - Project operations
- `/orgs/{org}/projectsV2/{project_number}` - Org project operations
- `/users/{username}/projectsV2/{project_number}/fields` - List/Create/Update fields
- `/users/{username}/projectsV2/{project_number}/items` - List/Create/Update items

---

## How Columns Relate to Status Fields in ProjectsV2

### Migration Mapping

**Classic Projects Columns** map to **ProjectsV2 Status Fields**:

| Classic (Columns) | ProjectsV2 (Fields) | Notes |
|---|---|---|
| Column name (e.g., "To Do") | Single Select Option | Each column becomes a status option |
| Column position | Field value option order | Column order is preserved as option order |
| Card in column | Item with status value | Card location becomes field value |

### Example Mapping

**Classic Project**:
```
Column 1: "Backlog"
Column 2: "In Progress"  
Column 3: "Review"
Column 4: "Done"
```

**Becomes in ProjectsV2**:
```
Field: "Status" (single_select type)
Options:
  - "Backlog"
  - "In Progress"
  - "Review"
  - "Done"
```

### Key Differences in Implementation

| Aspect | Classic Columns | ProjectsV2 Fields |
|---|---|---|
| **Concept** | Physical columns hold cards | Fields store values on items |
| **API Pattern** | Move cards between columns | Update item's field value |
| **Query Complexity** | Simple (which column?) | Rich (any field, multiple values) |
| **Customization** | Limited to column names | Full control per field type |
| **Automation** | Webhook-based | GitHub Actions + API |
| **Filtering** | Limited | Advanced with multiple fields |

---

## Example Workflows

### Classic Columns Workflow

```bash
# List columns in a project
curl -H "Authorization: token TOKEN" \
  https://api.github.com/projects/120/columns

# Create a new column
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Review"}' \
  https://api.github.com/projects/120/columns

# Move a column to the end
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"position":"last"}' \
  https://api.github.com/projects/columns/367/moves

# Delete a column
curl -X DELETE -H "Authorization: token TOKEN" \
  https://api.github.com/projects/columns/367
```

### ProjectsV2 Equivalent

```bash
# List fields in a ProjectV2
curl -H "Authorization: token TOKEN" \
  https://api.github.com/users/USERNAME/projectsV2/PROJ_NUM/fields

# Create a status field
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name":"Status",
    "data_type":"single_select",
    "options":[
      {"name":"Backlog"},
      {"name":"In Progress"},
      {"name":"Done"}
    ]
  }' \
  https://api.github.com/users/USERNAME/projectsV2/PROJ_NUM/fields

# Update an item's status field
curl -X PATCH -H "Authorization: token TOKEN" \
  https://api.github.com/users/USERNAME/projectsV2/PROJ_NUM/items/ITEM_ID
```

---

## API Compatibility Matrix

| Feature | Classic Projects | ProjectsV2 | GitHub Apps | Notes |
|---|---|---|---|---|
| List Columns/Fields | ✓ Columns | ✓ Fields | ✓ | Both support GitHub Apps |
| Create Column/Field | ✓ Columns | ✓ Fields | ✓ | Both support GitHub Apps |
| Update Column/Field | ✓ Rename only | ✓ Full edit | ✓ | V2 is more flexible |
| Delete Column/Field | ✓ | ✓ | ✓ | Both support deletion |
| Reorder Columns | ✓ Move endpoint | ✗ Not applicable | ✓ | V2 uses field value order |
| Custom Field Types | ✗ | ✓ 14+ types | ✓ | Major V2 advantage |
| Webhook Support | ✓ | ✓ | ✓ | Different event types |

---

## Important Notes for Migration

1. **Timeline**: Classic Projects are sunsetting April 1, 2025. Start migration planning now.

2. **Backwards Compatibility**: Old Classic Project endpoints will stop working after April 1, 2025.

3. **Data Preservation**: GitHub may provide migration tools, but it's unclear if automatic migration will occur.

4. **Authentication**: Both APIs support:
   - OAuth tokens
   - Personal access tokens
   - GitHub App authentication

5. **Rate Limits**: Standard GitHub API rate limits apply (60 requests/hour for unauthenticated, 5000/hour for authenticated)

6. **Error Handling**: Both APIs follow standard GitHub error responses with appropriate HTTP status codes

---

## Summary of Capabilities

### Classic Columns API Capabilities
- Create, read, update, delete columns
- Reorder columns on a board
- Basic workflow organization
- Simple API for basic use cases
- **Status: DEPRECATED - use ProjectsV2 instead**

### ProjectsV2 Fields API Capabilities
- 14+ field types for rich data modeling
- Create complex workflows with multiple dimensions
- Support for hierarchy (parent/sub issues)
- Iteration planning with sprint support
- Advanced filtering and sorting
- More powerful automation possibilities
- **Status: CURRENT STANDARD - actively maintained**

---

## Recommendations for tkan Development

Given the deprecation timeline:

1. **For GitHub Backend**: The CLAUDE.md mentions using ProjectsV2 API, which is correct.

2. **For Status Management**: 
   - Classic: Move cards between columns
   - ProjectsV2: Update item's status field value

3. **For Column/Field Updates**:
   - Classic: Use `/projects/columns/{id}` PATCH endpoint
   - ProjectsV2: Use `/projects/{number}/fields/{id}` endpoints

4. **Recommended Path**:
   - Continue with ProjectsV2 implementation
   - Use Fields instead of Columns concept
   - Map board columns to status field options
