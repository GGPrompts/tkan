# GitHub Project Columns REST API - Documentation Index

## Overview
Complete exploration of the GitHub Project Columns REST API, including comprehensive analysis of Classic Projects vs ProjectsV2 APIs.

**Important**: The Columns API is deprecated (May 23, 2024) and will be removed April 1, 2025. Use ProjectsV2 API instead.

---

## Documentation Files

### 1. GITHUB_COLUMNS_API_REFERENCE.md
**Purpose**: Comprehensive technical reference  
**Audience**: Developers needing full details  
**Contents**:
- Column management operations (6 endpoints)
- Data structure schemas
- Classic vs ProjectsV2 architecture
- Migration guidance
- Example workflows
- API compatibility matrix
- Troubleshooting

**When to use**: When you need deep technical knowledge or troubleshooting

---

### 2. COLUMNS_VS_PROJECTSV2.md
**Purpose**: Detailed comparison between APIs  
**Audience**: Architects and senior developers  
**Contents**:
- Quick reference table
- Endpoint-by-endpoint comparison
- Data structure differences
- Operational model differences
- Feature parity matrix
- Migration path for tkan
- Error handling comparison
- Deprecation timeline

**When to use**: When deciding which API to use or planning migrations

---

### 3. COLUMNS_API_QUICK_REFERENCE.md
**Purpose**: Fast lookup guide  
**Audience**: Developers actively coding  
**Contents**:
- All endpoints in table format
- curl code examples
- Column object structure
- HTTP status codes
- Position values for moving
- Common issues and solutions
- Authentication details
- Rate limiting info

**When to use**: For quick lookups while coding or debugging

---

## Quick Navigation

### By Task

**I want to list columns**
- Quick: COLUMNS_API_QUICK_REFERENCE.md → "List Columns" section
- Details: GITHUB_COLUMNS_API_REFERENCE.md → "List Project Columns" section

**I want to create a column**
- Quick: COLUMNS_API_QUICK_REFERENCE.md → "Create Column" section
- Details: GITHUB_COLUMNS_API_REFERENCE.md → "Create a Project Column" section

**I want to rename a column**
- Quick: COLUMNS_API_QUICK_REFERENCE.md → "Rename Column" section
- Details: GITHUB_COLUMNS_API_REFERENCE.md → "Update a Project Column" section

**I want to reorder columns**
- Quick: COLUMNS_API_QUICK_REFERENCE.md → "Move Column" section
- Details: GITHUB_COLUMNS_API_REFERENCE.md → "Move a Project Column" section

**I want to understand Columns vs ProjectsV2**
- COLUMNS_VS_PROJECTSV2.md → Start with "Quick Reference Table"

**I'm migrating from Columns to ProjectsV2**
- COLUMNS_VS_PROJECTSV2.md → "Migration Path for tkan" section
- GITHUB_COLUMNS_API_REFERENCE.md → "Recommendations for tkan Development" section

**I need code examples**
- COLUMNS_API_QUICK_REFERENCE.md → "Quick Code Examples" section
- GITHUB_COLUMNS_API_REFERENCE.md → "Example Workflows" section

**I'm getting an error**
- COLUMNS_API_QUICK_REFERENCE.md → "Common Issues" section
- GITHUB_COLUMNS_API_REFERENCE.md → "Troubleshooting" section

---

## Key Facts at a Glance

### Columns API (DEPRECATED)

**6 Endpoints**:
1. `GET /projects/{project_id}/columns` - List
2. `POST /projects/{project_id}/columns` - Create
3. `GET /projects/columns/{column_id}` - Get
4. `PATCH /projects/columns/{column_id}` - Update (rename)
5. `DELETE /projects/columns/{column_id}` - Delete
6. `POST /projects/columns/{column_id}/moves` - Reorder

**Status Codes**:
- 200: Success
- 201: Created
- 204: Deleted
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 422: Validation Failed

**Column Structure**:
```json
{
  "id": 367,
  "name": "To Do",
  "url": "https://api.github.com/projects/columns/367",
  "project_url": "https://api.github.com/projects/120",
  "cards_url": "https://api.github.com/projects/columns/367/cards",
  "node_id": "MDEzOlByb2plY3RDb2x1bW4zNjc=",
  "created_at": "2016-09-05T14:18:44Z",
  "updated_at": "2016-09-05T14:22:28Z"
}
```

---

## Comparison Summary

| Feature | Columns (Classic) | ProjectsV2 |
|---------|-------------------|-----------|
| Status | Deprecated (removal Apr 2025) | Current standard |
| Structure | 1D (columns only) | Multi-dimensional (fields) |
| Field Types | None (names only) | 14+ types |
| Item Movement | Move between columns | Update field value |
| Customization | Minimal | Extensive |
| Automation | Webhooks | Webhooks + Actions |
| Timeline | Frozen | Actively developed |

---

## For tkan Project

tkan is **correctly using ProjectsV2** per CLAUDE.md.

**Key mapping**:
- Classic Column → ProjectsV2 Status Field
- Column name → Status option value
- Card in column → Item with status value

**Recommended approach**:
- Continue with ProjectsV2 implementation
- Use Status field for kanban board
- Leverage other field types for additional data

---

## Important Dates

- **May 23, 2024**: Columns API deprecated
- **April 1, 2025**: Columns API removed
- **Today (Oct 28, 2024)**: ~6 months until removal

---

## External Resources

- GitHub REST API Docs: https://docs.github.com/en/rest
- Projects REST API: https://docs.github.com/en/rest/projects
- ProjectsV2 GraphQL: https://docs.github.com/en/graphql/reference/objects#projectv2
- Deprecation Notice: https://github.blog/changelog/2024-05-23-sunset-notice-projects-classic/

---

## File Statistics

| Document | Size | Sections | Focus |
|----------|------|----------|-------|
| GITHUB_COLUMNS_API_REFERENCE.md | 12 KB | 10+ | Comprehensive reference |
| COLUMNS_VS_PROJECTSV2.md | 11 KB | 15+ | Comparison & migration |
| COLUMNS_API_QUICK_REFERENCE.md | 4.8 KB | 10+ | Quick lookup |
| API_DOCUMENTATION_INDEX.md | This file | Navigation | Index & guide |

**Total**: ~28 KB of documentation

---

## How to Use These Documents

1. **New to the topic?**  
   Start with: COLUMNS_VS_PROJECTSV2.md → "Quick Reference Table"

2. **Need quick lookup?**  
   Use: COLUMNS_API_QUICK_REFERENCE.md

3. **Deep dive required?**  
   Read: GITHUB_COLUMNS_API_REFERENCE.md

4. **Planning implementation?**  
   Focus on: COLUMNS_VS_PROJECTSV2.md → "Migration Path for tkan"

5. **Troubleshooting?**  
   Check: COLUMNS_API_QUICK_REFERENCE.md → "Common Issues"

---

## Recommendations

For the tkan project:
- No changes needed (already using ProjectsV2)
- These docs are for reference and education
- File for future architecture discussions
- Share with team members learning about GitHub APIs

For any new GitHub integration:
- Always choose ProjectsV2 over Columns
- Plan for April 2025 Columns removal
- Use these docs for API design patterns

---

Created: October 28, 2024  
Based on GitHub OpenAPI Spec v2022-11-28  
Status: Columns API Deprecated, ProjectsV2 Current

