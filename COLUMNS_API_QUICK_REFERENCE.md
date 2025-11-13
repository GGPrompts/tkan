# GitHub Project Columns API - Quick Reference

## Critical Note: DEPRECATED (Removal: April 1, 2025)
Use **ProjectsV2 API** instead for new projects.

---

## All Column Endpoints

| Method | Endpoint | Purpose | Status |
|--------|----------|---------|--------|
| GET | `/projects/{project_id}/columns` | List all columns | 200 OK |
| POST | `/projects/{project_id}/columns` | Create column | 201 Created |
| GET | `/projects/columns/{column_id}` | Get column details | 200 OK |
| PATCH | `/projects/columns/{column_id}` | Rename column | 200 OK |
| DELETE | `/projects/columns/{column_id}` | Delete column | 204 No Content |
| POST | `/projects/columns/{column_id}/moves` | Reorder column | 201 Created |

---

## Quick Code Examples

### List Columns
```bash
curl -H "Authorization: token TOKEN" \
  https://api.github.com/projects/120/columns
```

### Create Column
```bash
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Stage"}' \
  https://api.github.com/projects/120/columns
```

### Rename Column
```bash
curl -X PATCH -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Review"}' \
  https://api.github.com/projects/columns/367
```

### Delete Column
```bash
curl -X DELETE -H "Authorization: token TOKEN" \
  https://api.github.com/projects/columns/367
```

### Move Column
```bash
# To end
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"position":"last"}' \
  https://api.github.com/projects/columns/367/moves

# After another column
curl -X POST -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"position":"after:365"}' \
  https://api.github.com/projects/columns/367/moves
```

---

## Column Object Structure

```json
{
  "id": 367,
  "node_id": "MDEzOlByb2plY3RDb2x1bW4zNjc=",
  "url": "https://api.github.com/projects/columns/367",
  "project_url": "https://api.github.com/projects/120",
  "cards_url": "https://api.github.com/projects/columns/367/cards",
  "name": "To Do",
  "created_at": "2016-09-05T14:18:44Z",
  "updated_at": "2016-09-05T14:22:28Z"
}
```

---

## HTTP Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| 200 | Success | Operation completed |
| 201 | Created | Resource created successfully |
| 204 | No Content | Deleted successfully |
| 304 | Not Modified | No changes needed (caching) |
| 401 | Unauthorized | Need valid token |
| 403 | Forbidden | No permission to project |
| 404 | Not Found | Column/project doesn't exist |
| 422 | Validation Failed | Bad input (wrong position format, etc.) |

---

## Position Values for Moving

```
"first"          → Move to beginning
"last"           → Move to end
"after:COLUMN_ID" → Place after specific column (e.g., "after:365")
```

---

## Alternative: Use ProjectsV2 Instead

If building new features, use ProjectsV2 which offers:
- Multiple field types (not just columns)
- Better automation support
- Longer-term support
- More powerful APIs

### ProjectsV2 Equivalent

```bash
# Get status field for a ProjectV2
curl -H "Authorization: token TOKEN" \
  https://api.github.com/users/USERNAME/projectsV2/7/fields

# Update item's status field
curl -X PATCH -H "Authorization: token TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "field_values": {
      "field_id": 123456789,
      "value": "IN_PROGRESS"
    }
  }' \
  https://api.github.com/users/USERNAME/projectsV2/7/items/ITEM_ID
```

---

## Common Issues

| Problem | Solution |
|---------|----------|
| 401 Unauthorized | Check token validity: `gh auth status` |
| 403 Forbidden | Verify you have project write permissions |
| 404 Not Found | Verify correct project_id and column_id |
| 422 Invalid position | Use "first", "last", or "after:COLUMN_ID" format |
| Column name already exists | Different names needed in same project |

---

## Authentication

All requests require authentication:

```bash
# Using token
-H "Authorization: token YOUR_TOKEN"

# Using gh CLI (easier)
gh api /projects/{id}/columns
```

---

## Pagination

List endpoint supports pagination:
```bash
curl "https://api.github.com/projects/120/columns?page=2&per_page=50" \
  -H "Authorization: token TOKEN"
```

---

## Rate Limiting

- Unauthenticated: 60 requests/hour
- Authenticated: 5,000 requests/hour

Check headers:
```bash
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 4999
X-RateLimit-Reset: 1234567890
```

---

## Related APIs

**Cards in Columns**:
- `GET /projects/columns/{column_id}/cards`
- `POST /projects/columns/{column_id}/cards`
- `PATCH /projects/columns/cards/{card_id}`
- `DELETE /projects/columns/cards/{card_id}`

**Note**: Cards endpoints also deprecated with Columns

---

## Timeline

- Deprecated: May 23, 2024
- Removal Date: April 1, 2025
- Status: Frozen, no new features

**Action**: Migrate to ProjectsV2 before April 2025
