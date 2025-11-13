# Error Codes Reference

Complete reference for HTTP status codes and GraphQL error types in GitHub Projects API.

## HTTP Status Codes

### 2xx Success

#### 200 OK
```
Request succeeded and returned data
```

**Common with:**
- GET requests to REST endpoints
- Successful GraphQL queries

**Response:**
```json
{
  "data": { ... },
  "status": 200
}
```

#### 201 Created
```
Resource successfully created
```

**Common with:**
- POST requests creating new items
- Adding issues to projects

**Example:**
```bash
gh project item-create 7 --owner GGPrompts --title "New task"
# Returns 201 with created item data
```

#### 204 No Content
```
Operation successful, no content returned
```

**Common with:**
- DELETE operations
- Some update operations

### 3xx Redirection

#### 304 Not Modified
```
Cached content is still valid
```

**Common with:**
- Requests with `If-None-Match` or `If-Modified-Since` headers
- API responses when content hasn't changed

**Action:** Use cached data

### 4xx Client Errors

#### 401 Unauthorized
```
Authentication required or invalid credentials
```

**Causes:**
- Missing authentication token
- Invalid or expired token
- Token doesn't have required scopes

**Solutions:**
```bash
# Check auth status
gh auth status

# Re-authenticate
gh auth login

# Refresh with project scope
gh auth refresh -h github.com -s project
```

**Example Error:**
```json
{
  "message": "Requires authentication",
  "documentation_url": "https://docs.github.com/rest/..."
}
```

#### 403 Forbidden
```
Authenticated but not authorized to perform action
```

**Causes:**
- Token lacks required scopes
- Insufficient permissions on project
- Rate limit exceeded (see X-RateLimit headers)
- Resource access denied

**Solutions:**
```bash
# Add project scope
gh auth refresh -h github.com -s project

# Check rate limit
gh api rate_limit

# For organization projects, verify team permissions
```

**Example Error:**
```json
{
  "message": "Resource not accessible by integration",
  "documentation_url": "https://docs.github.com/graphql/..."
}
```

**Rate Limit Headers:**
```
X-RateLimit-Limit: 5000
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1640995200
```

#### 404 Not Found
```
Resource doesn't exist or you don't have access
```

**Causes:**
- Invalid project number
- Wrong owner name
- Item was deleted
- Private project you can't access

**Solutions:**
```bash
# Verify project exists
gh project list --owner GGPrompts

# Check project number in URL
# https://github.com/users/GGPrompts/projects/7
#                                          ^^^ This is the number

# For organizations, ensure you have access
gh project list --owner myorg
```

**Example Error:**
```json
{
  "errors": [{
    "type": "NOT_FOUND",
    "message": "Could not resolve to a ProjectV2 with the number 7."
  }]
}
```

#### 422 Unprocessable Entity
```
Request syntax is valid but semantically incorrect
```

**Causes:**
- Invalid field IDs
- Invalid option IDs
- Type mismatch (e.g., text for number field)
- Missing required fields
- Validation constraints violated

**Solutions:**
```bash
# Re-fetch field metadata
./discover-fields.sh GGPrompts 7

# Verify field IDs are current
jq '.fields[] | {name, id}' fields-cache.json

# Check option IDs for single-select fields
jq '.lookups.status_field.options' fields-cache.json
```

**Example Errors:**
```json
{
  "errors": [{
    "type": "UNPROCESSABLE",
    "message": "Field does not exist on this project"
  }]
}
```

```json
{
  "errors": [{
    "message": "Variable $optionId of type String! was provided invalid value",
    "locations": [{"line": 1, "column": 8}]
  }]
}
```

### 5xx Server Errors

#### 500 Internal Server Error
```
GitHub server encountered an error
```

**Action:** Retry with exponential backoff

#### 502 Bad Gateway
```
GitHub is temporarily unavailable
```

**Action:** Retry after a delay

#### 503 Service Unavailable
```
GitHub is undergoing maintenance or overloaded
```

**Action:** Check [GitHub Status](https://www.githubstatus.com/)

## GraphQL Error Types

### FORBIDDEN

**Cause:** Insufficient permissions

**Example:**
```json
{
  "errors": [{
    "type": "FORBIDDEN",
    "path": ["user", "projectV2"],
    "message": "Resource not accessible by integration"
  }]
}
```

**Solution:**
```bash
gh auth refresh -h github.com -s project
```

### NOT_FOUND

**Cause:** Resource doesn't exist

**Example:**
```json
{
  "errors": [{
    "type": "NOT_FOUND",
    "path": ["user", "projectV2"],
    "message": "Could not resolve to a ProjectV2 with the number 99."
  }]
}
```

**Solution:** Verify project number and owner

### UNPROCESSABLE

**Cause:** Invalid input data

**Example:**
```json
{
  "errors": [{
    "type": "UNPROCESSABLE",
    "path": ["updateProjectV2ItemFieldValue"],
    "message": "Field does not exist on this project",
    "extensions": {
      "fieldId": "invalid-field-id"
    }
  }]
}
```

**Solution:** Re-fetch field metadata

### RATE_LIMITED

**Cause:** Too many requests

**Example:**
```json
{
  "errors": [{
    "type": "RATE_LIMITED",
    "message": "API rate limit exceeded for user ID 12345."
  }]
}
```

**Check Limit:**
```bash
gh api rate_limit
```

**Response:**
```json
{
  "resources": {
    "graphql": {
      "limit": 5000,
      "remaining": 0,
      "reset": 1640995200
    }
  }
}
```

**Solution:** Wait until reset time or implement caching

### VALIDATION_FAILED

**Cause:** Input validation failed

**Example:**
```json
{
  "errors": [{
    "type": "VALIDATION_FAILED",
    "message": "Validation Failed",
    "errors": [
      {
        "field": "title",
        "code": "missing",
        "message": "title is required"
      }
    ]
  }]
}
```

**Solution:** Fix input according to validation message

### INSUFFICIENT_SCOPES

**Cause:** Token missing required OAuth scopes

**Example:**
```json
{
  "message": "This endpoint requires you to be authenticated with the 'project' scope",
  "documentation_url": "https://docs.github.com/graphql"
}
```

**Solution:**
```bash
gh auth refresh -h github.com -s project
```

## Common Error Scenarios

### Scenario 1: Token Has No Project Scope

**Error:**
```
Error: HTTP 401: Requires authentication (https://docs.github.com/rest)
```

or

```json
{
  "message": "Token has not been granted the required scopes"
}
```

**Solution:**
```bash
# Check current scopes
gh auth status

# Add project scope
gh auth refresh -h github.com -s project

# Verify
gh auth status
# Should show "Token scopes: ..., project, ..."
```

### Scenario 2: Hardcoded Field IDs Are Stale

**Error:**
```json
{
  "errors": [{
    "type": "UNPROCESSABLE",
    "message": "Field does not exist on this project"
  }]
}
```

**Why:** Project was recreated or fields were modified

**Solution:**
```bash
# Re-discover fields
./scripts/discover-fields.sh GGPrompts 7

# Update your code/config with new field IDs
jq '.lookups.status_field.id' fields-cache.json
```

### Scenario 3: Invalid Status Option

**Error:**
```json
{
  "errors": [{
    "type": "UNPROCESSABLE",
    "message": "Option does not exist for this field"
  }]
}
```

**Why:** Status option was renamed or deleted

**Solution:**
```bash
# Check available options
jq '.lookups.status_field.options' fields-cache.json

# Output:
# [
#   {"name": "Todo", "id": "f75ad846"},
#   {"name": "In Progress", "id": "47fc9ee4"},
#   {"name": "Done", "id": "98236657"}
# ]

# Use exact option name from output
```

### Scenario 4: Rate Limit Exceeded

**Error:**
```json
{
  "errors": [{
    "type": "RATE_LIMITED",
    "message": "API rate limit exceeded"
  }]
}
```

**Check Status:**
```bash
gh api rate_limit | jq '.resources.graphql'
```

**Output:**
```json
{
  "limit": 5000,
  "remaining": 0,
  "reset": 1640995200,
  "used": 5000
}
```

**Solution:**
```bash
# Calculate wait time
RESET_TIME=$(gh api rate_limit | jq -r '.resources.graphql.reset')
WAIT_SECONDS=$((RESET_TIME - $(date +%s)))
echo "Wait $WAIT_SECONDS seconds until rate limit resets"

# Or implement exponential backoff
sleep 60  # Wait 1 minute and retry
```

**Prevention:**
- Cache field metadata (refresh every 5 minutes instead of every call)
- Use batch mutations for multiple updates
- Implement request throttling

### Scenario 5: Project Not Found

**Error:**
```json
{
  "errors": [{
    "type": "NOT_FOUND",
    "message": "Could not resolve to a ProjectV2 with the number 7."
  }]
}
```

**Debugging Steps:**

1. **Verify project exists:**
```bash
gh project list --owner GGPrompts
```

2. **Check project number:**
```bash
# Project URL: https://github.com/users/GGPrompts/projects/7
#                                                    ^^^^^ ^^^
#                                                    owner  number
```

3. **Try organization context:**
```bash
# If user context fails, try organization
gh project list --owner myorg
```

4. **Verify access:**
```bash
gh project view 7 --owner GGPrompts
```

### Scenario 6: Interface Conversion Error (tkan-specific)

**Error in tkan:**
```
panic: interface conversion: interface {} is nil, not map[string]interface {}
```

**Cause:** Expected data structure doesn't match actual response

**Debugging:**
```bash
# Inspect raw API response
gh project item-list 7 --owner GGPrompts --format json > response.json

# Check structure
jq '.items[0] | keys' response.json

# Check field values structure
jq '.items[0].fieldValues' response.json
```

**Fix:** Update parsing code to handle nil values or different structure

## Error Handling in Code

### Bash Script Pattern

```bash
#!/bin/bash

execute_graphql() {
    local query=$1
    local max_retries=3
    local retry=0

    while [ $retry -lt $max_retries ]; do
        RESPONSE=$(gh api graphql -f query="$query" 2>&1)
        EXIT_CODE=$?

        # Check for errors
        if [ $EXIT_CODE -ne 0 ]; then
            echo "Error: API call failed" >&2
            ((retry++))
            sleep $((2 ** retry))  # Exponential backoff
            continue
        fi

        # Check for GraphQL errors
        if echo "$RESPONSE" | jq -e '.errors' >/dev/null 2>&1; then
            ERROR_TYPE=$(echo "$RESPONSE" | jq -r '.errors[0].type // "UNKNOWN"')

            case $ERROR_TYPE in
                RATE_LIMITED)
                    echo "Rate limited, waiting..." >&2
                    sleep 60
                    ((retry++))
                    continue
                    ;;
                NOT_FOUND|FORBIDDEN|INSUFFICIENT_SCOPES)
                    echo "Fatal error: $ERROR_TYPE" >&2
                    echo "$RESPONSE" | jq -r '.errors[] | .message' >&2
                    return 1
                    ;;
                UNPROCESSABLE)
                    echo "Validation error:" >&2
                    echo "$RESPONSE" | jq -r '.errors[] | .message' >&2
                    return 1
                    ;;
                *)
                    echo "Unknown error: $ERROR_TYPE" >&2
                    ((retry++))
                    sleep $((2 ** retry))
                    continue
                    ;;
            esac
        fi

        # Success
        echo "$RESPONSE"
        return 0
    done

    echo "Max retries exceeded" >&2
    return 1
}
```

### Go Error Handling Pattern

```go
func (g *GitHubBackend) executeGraphQLWithRetry(query string) (interface{}, error) {
    maxRetries := 3
    baseDelay := time.Second

    for attempt := 0; attempt < maxRetries; attempt++ {
        result, err := g.executeGraphQL(query, nil)

        if err == nil {
            return result, nil
        }

        // Check error type
        if graphQLErr, ok := err.(*GraphQLError); ok {
            switch graphQLErr.Type {
            case "RATE_LIMITED":
                // Wait and retry
                delay := baseDelay * time.Duration(1<<uint(attempt))
                log.Printf("Rate limited, waiting %v", delay)
                time.Sleep(delay)
                continue

            case "NOT_FOUND", "FORBIDDEN", "INSUFFICIENT_SCOPES":
                // Fatal errors, don't retry
                return nil, fmt.Errorf("fatal error: %w", err)

            case "UNPROCESSABLE":
                // Validation error, don't retry
                return nil, fmt.Errorf("validation failed: %w", err)

            default:
                // Unknown error, retry
                log.Printf("Unknown error (attempt %d/%d): %v", attempt+1, maxRetries, err)
                time.Sleep(baseDelay * time.Duration(1<<uint(attempt)))
                continue
            }
        }

        // Non-GraphQL error, retry
        time.Sleep(baseDelay * time.Duration(1<<uint(attempt)))
    }

    return nil, fmt.Errorf("max retries exceeded")
}
```

## Monitoring and Debugging

### Log All API Calls

```bash
# Enable GitHub CLI debugging
export GH_DEBUG=api

gh project item-list 7 --owner GGPrompts

# Shows full request/response
```

### Check Rate Limit Before Operations

```bash
#!/bin/bash

check_rate_limit() {
    REMAINING=$(gh api rate_limit | jq -r '.resources.graphql.remaining')
    LIMIT=$(gh api rate_limit | jq -r '.resources.graphql.limit')

    echo "Rate limit: $REMAINING / $LIMIT remaining"

    if [ "$REMAINING" -lt 100 ]; then
        echo "Warning: Low rate limit remaining" >&2
        return 1
    fi

    return 0
}

# Before expensive operation
if ! check_rate_limit; then
    echo "Aborting due to low rate limit"
    exit 1
fi
```

## See Also

- `../SKILL.md` - Main skill documentation
- `graphql-queries.md` - Query examples
- `field-discovery.md` - Field discovery guide
