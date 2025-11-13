# GraphQL Query Reference

Complete reference of GraphQL queries and mutations for GitHub Projects v2.

## Table of Contents
- [Project Queries](#project-queries)
- [Item Queries](#item-queries)
- [Field Queries](#field-queries)
- [Mutations](#mutations)
- [Pagination](#pagination)
- [Error Handling](#error-handling)

## Project Queries

### Get Project Details (User)

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      id
      number
      title
      shortDescription
      url
      readme
      public
      closed
      createdAt
      updatedAt
    }
  }
}
```

### Get Project Details (Organization)

```graphql
query {
  organization(login: "myorg") {
    projectV2(number: 7) {
      id
      title
      shortDescription
    }
  }
}
```

### List All User Projects

```graphql
query {
  user(login: "GGPrompts") {
    projectsV2(first: 20) {
      nodes {
        id
        number
        title
        url
      }
      pageInfo {
        hasNextPage
        endCursor
      }
    }
  }
}
```

### List All Organization Projects

```graphql
query {
  organization(login: "myorg") {
    projectsV2(first: 20) {
      nodes {
        id
        number
        title
        url
      }
    }
  }
}
```

## Item Queries

### Get All Items (Basic)

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          id
          type
          createdAt
          updatedAt
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }
}
```

### Get Items with Content

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          id
          type
          content {
            ... on DraftIssue {
              id
              title
              body
              createdAt
              updatedAt
            }
            ... on Issue {
              id
              number
              title
              body
              state
              url
              repository {
                name
                owner {
                  login
                }
              }
              labels(first: 10) {
                nodes {
                  name
                  color
                }
              }
              assignees(first: 10) {
                nodes {
                  login
                  avatarUrl
                }
              }
            }
            ... on PullRequest {
              id
              number
              title
              body
              state
              url
              repository {
                name
              }
              isDraft
              merged
            }
          }
        }
      }
    }
  }
}
```

### Get Items with Field Values

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          id
          type
          content {
            ... on DraftIssue {
              title
            }
            ... on Issue {
              title
              number
            }
          }
          fieldValues(first: 20) {
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
              ... on ProjectV2ItemFieldDateValue {
                date
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
              ... on ProjectV2ItemFieldIterationValue {
                title
                startDate
                duration
                field {
                  ... on ProjectV2IterationField {
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
}
```

### Get Single Item by ID

```graphql
query {
  node(id: "PVTI_lADOBd...") {
    ... on ProjectV2Item {
      id
      type
      content {
        ... on DraftIssue {
          title
          body
        }
      }
      fieldValues(first: 20) {
        nodes {
          ... on ProjectV2ItemFieldSingleSelectValue {
            name
          }
        }
      }
    }
  }
}
```

### Get Item by Field Value (Status)

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          id
          content {
            ... on DraftIssue {
              title
            }
          }
          fieldValueByName(name: "Status") {
            ... on ProjectV2ItemFieldSingleSelectValue {
              name
            }
          }
        }
      }
    }
  }
}
```

## Field Queries

### Get All Fields (Complete)

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
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
              completedIterations {
                id
                title
              }
            }
          }
        }
      }
    }
  }
}
```

### Get Specific Field by Name

```graphql
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      field(name: "Status") {
        ... on ProjectV2SingleSelectField {
          id
          name
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
```

## Mutations

### Create Draft Issue

```graphql
mutation {
  addProjectV2DraftIssue(input: {
    projectId: "PVT_kwDOBd..."
    title: "New task"
    body: "Task description"
  }) {
    projectItem {
      id
      content {
        ... on DraftIssue {
          title
          body
        }
      }
    }
  }
}
```

### Add Existing Issue to Project

```graphql
mutation {
  addProjectV2ItemById(input: {
    projectId: "PVT_kwDOBd..."
    contentId: "I_kwDOBd..."  # Issue node ID
  }) {
    item {
      id
    }
  }
}
```

### Update Item Field Value (Single Select)

```graphql
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
    value: {
      singleSelectOptionId: "f75ad846"
    }
  }) {
    projectV2Item {
      id
      fieldValueByName(name: "Status") {
        ... on ProjectV2ItemFieldSingleSelectValue {
          name
        }
      }
    }
  }
}
```

### Update Item Field Value (Text)

```graphql
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
    value: {
      text: "Updated text value"
    }
  }) {
    projectV2Item {
      id
    }
  }
}
```

### Update Item Field Value (Number)

```graphql
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
    value: {
      number: 5.0
    }
  }) {
    projectV2Item {
      id
    }
  }
}
```

### Update Item Field Value (Date)

```graphql
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
    value: {
      date: "2024-12-31"
    }
  }) {
    projectV2Item {
      id
    }
  }
}
```

### Update Item Field Value (Iteration)

```graphql
mutation {
  updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
    value: {
      iterationId: "PVTI_lADOBd..."
    }
  }) {
    projectV2Item {
      id
    }
  }
}
```

### Clear Field Value

```graphql
mutation {
  clearProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
    fieldId: "PVTF_lADOBd..."
  }) {
    projectV2Item {
      id
    }
  }
}
```

### Update Draft Issue

```graphql
mutation {
  updateProjectV2DraftIssue(input: {
    draftIssueId: "DI_kwDOBd..."
    title: "Updated title"
    body: "Updated description"
  }) {
    draftIssue {
      id
      title
      body
      updatedAt
    }
  }
}
```

### Delete Item from Project

```graphql
mutation {
  deleteProjectV2Item(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
  }) {
    deletedItemId
  }
}
```

### Archive Item

```graphql
mutation {
  archiveProjectV2Item(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd..."
  }) {
    item {
      id
      isArchived
    }
  }
}
```

### Batch Mutations (Multiple Operations)

```graphql
mutation {
  item1: updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd1..."
    fieldId: "PVTF_lADOBd..."
    value: {singleSelectOptionId: "47fc9ee4"}
  }) {
    projectV2Item { id }
  }

  item2: updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd2..."
    fieldId: "PVTF_lADOBd..."
    value: {singleSelectOptionId: "47fc9ee4"}
  }) {
    projectV2Item { id }
  }

  item3: updateProjectV2ItemFieldValue(input: {
    projectId: "PVT_kwDOBd..."
    itemId: "PVTI_lADOBd3..."
    fieldId: "PVTF_lADOBd..."
    value: {singleSelectOptionId: "47fc9ee4"}
  }) {
    projectV2Item { id }
  }
}
```

## Pagination

### Forward Pagination

```graphql
query($cursor: String) {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 50, after: $cursor) {
        nodes {
          id
        }
        pageInfo {
          hasNextPage
          endCursor
        }
      }
    }
  }
}
```

### Complete Pagination Example (Bash)

```bash
#!/bin/bash

OWNER="GGPrompts"
PROJECT_NUMBER=7
CURSOR=""
HAS_NEXT=true

while [ "$HAS_NEXT" = "true" ]; do
    RESPONSE=$(gh api graphql -f query='
    query($cursor: String) {
      user(login: "'$OWNER'") {
        projectV2(number: '$PROJECT_NUMBER') {
          items(first: 100, after: $cursor) {
            nodes {
              id
              content {
                ... on DraftIssue { title }
              }
            }
            pageInfo {
              hasNextPage
              endCursor
            }
          }
        }
      }
    }
    ' -f cursor="$CURSOR")

    # Extract items
    echo "$RESPONSE" | jq -r '.data.user.projectV2.items.nodes[] | .id'

    # Update pagination state
    HAS_NEXT=$(echo "$RESPONSE" | jq -r '.data.user.projectV2.items.pageInfo.hasNextPage')
    CURSOR=$(echo "$RESPONSE" | jq -r '.data.user.projectV2.items.pageInfo.endCursor')
done
```

## Error Handling

### Check for Errors in Response

```bash
RESPONSE=$(gh api graphql -f query='...')

if echo "$RESPONSE" | jq -e '.errors' >/dev/null 2>&1; then
    echo "GraphQL Errors:"
    echo "$RESPONSE" | jq -r '.errors[] | "  - [\(.type)] \(.message)"'
    exit 1
fi
```

### Common Error Types

#### Authentication Error
```json
{
  "errors": [{
    "type": "FORBIDDEN",
    "message": "Resource not accessible by integration"
  }]
}
```

**Solution:** Add project scope: `gh auth refresh -h github.com -s project`

#### Not Found Error
```json
{
  "errors": [{
    "type": "NOT_FOUND",
    "message": "Could not resolve to a ProjectV2 with the number 7."
  }]
}
```

**Solution:** Verify project number and owner

#### Invalid Field Error
```json
{
  "errors": [{
    "type": "UNPROCESSABLE",
    "message": "Field does not exist on this project"
  }]
}
```

**Solution:** Re-run field discovery to get current field IDs

#### Rate Limit Error
```json
{
  "errors": [{
    "type": "RATE_LIMITED",
    "message": "API rate limit exceeded"
  }]
}
```

**Solution:** Implement exponential backoff, reduce API calls

## Using Variables

### Query with Variables

```bash
gh api graphql \
  -f query='
    query($owner: String!, $number: Int!) {
      user(login: $owner) {
        projectV2(number: $number) {
          id
          title
        }
      }
    }
  ' \
  -f owner="GGPrompts" \
  -F number=7
```

### Mutation with Variables

```bash
gh api graphql \
  -f query='
    mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
      updateProjectV2ItemFieldValue(input: {
        projectId: $projectId
        itemId: $itemId
        fieldId: $fieldId
        value: {singleSelectOptionId: $optionId}
      }) {
        projectV2Item { id }
      }
    }
  ' \
  -f projectId="PVT_kwDOBd..." \
  -f itemId="PVTI_lADOBd..." \
  -f fieldId="PVTF_lADOBd..." \
  -f optionId="47fc9ee4"
```

## Introspection Queries

### Get Project Schema

```graphql
{
  __type(name: "ProjectV2") {
    name
    fields {
      name
      type {
        name
        kind
      }
    }
  }
}
```

### Get Field Types

```graphql
{
  __type(name: "ProjectV2FieldType") {
    name
    enumValues {
      name
      description
    }
  }
}
```

## Performance Tips

### Request Only What You Need
```graphql
# BAD: Requesting unnecessary data
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          # Fetching all fields even if not needed
          id
          type
          content { ... }
          fieldValues(first: 20) { ... }
        }
      }
    }
  }
}

# GOOD: Minimal query
query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          id
          fieldValueByName(name: "Status") {
            ... on ProjectV2ItemFieldSingleSelectValue {
              name
            }
          }
        }
      }
    }
  }
}
```

### Use Fragments for Reusability

```graphql
fragment ItemFields on ProjectV2Item {
  id
  type
  content {
    ... on DraftIssue {
      title
      body
    }
  }
  fieldValueByName(name: "Status") {
    ... on ProjectV2ItemFieldSingleSelectValue {
      name
    }
  }
}

query {
  user(login: "GGPrompts") {
    projectV2(number: 7) {
      items(first: 100) {
        nodes {
          ...ItemFields
        }
      }
    }
  }
}
```

## See Also

- `../SKILL.md` - Main skill documentation
- `field-discovery.md` - Field discovery implementation guide
- `error-codes.md` - HTTP and GraphQL error reference
