# Claude Integration Guide for tkan

## Overview
tkan is a Terminal Kanban Board that uses GitHub Projects as its backend, enabling seamless collaboration between human developers and AI assistants like Claude.

## Architecture Decision: GitHub Projects Backend

### Why GitHub Projects?
We chose GitHub Projects over local YAML storage for several key reasons:

1. **Claude Integration**: Claude can directly manage tasks using `gh` CLI commands, eliminating complex workarounds
2. **Multi-Platform Access**: Edit tasks from terminal, web browser, or mobile devices
3. **Real-Time Sync**: Changes from any client are immediately reflected everywhere
4. **Built-in Collaboration**: Teams can work together without file conflicts
5. **No File Watching Needed**: Direct API access instead of file system monitoring

### Trade-offs Considered
- **Local YAML**: Fast, offline-capable, but difficult for Claude to modify
- **GitHub Projects**: Requires internet, slightly slower, but enables true collaboration
- **Hybrid Approach**: Possible future enhancement for offline work with sync

## How Claude Can Help

### View Project Tasks
```bash
# List all items in the project
gh project item-list 7 --owner GGPrompts --format json

# Get project details
gh api /users/GGPrompts/projectsV2/7
```

### Create New Tasks
```bash
# Create a draft issue
gh project item-create 7 --owner GGPrompts \
  --title "Task title" \
  --body "Task description"
```

### Update Task Status
```bash
# Using GitHub REST API
gh api --method PATCH /users/GGPrompts/projectsV2/7/items/{item_id} \
  --field field_id=233495315 \
  --field value=f75ad846  # Todo status
```

### Field IDs for tkan Development Project (#7)
- **Status Field**: `233495315`
  - Todo: `f75ad846`
  - In Progress: `47fc9ee4`
  - Done: `98236657`

## GitHub API Documentation

### Essential GitHub Docs Pages

#### Projects REST API (Recommended)
- [Projects Overview](https://docs.github.com/en/rest/projects/projects?apiVersion=2022-11-28)
- [Project Items](https://docs.github.com/en/rest/projects/items?apiVersion=2022-11-28)
- [Project Fields](https://docs.github.com/en/rest/projects/fields?apiVersion=2022-11-28)
- [Project Columns](https://docs.github.com/en/rest/projects/columns?apiVersion=2022-11-28)

#### GraphQL API (Alternative)
- [ProjectV2 GraphQL](https://docs.github.com/en/graphql/reference/objects#projectv2)
- [Managing Projects with GraphQL](https://docs.github.com/en/issues/planning-and-tracking-with-projects/automating-your-project/using-the-api-to-manage-projects)

#### GitHub CLI
- [gh project commands](https://cli.github.com/manual/gh_project)
- [gh api command](https://cli.github.com/manual/gh_api)

## Implementation Files

### Core Backend Files
- `backend.go` - Backend interface definition
- `backend_github.go` - GitHub Projects implementation
- `main.go` - CLI flags and backend selection

### Key Functions
- `NewGitHubBackend()` - Creates GitHub backend instance
- `LoadBoard()` - Fetches project and converts to tkan format
- `MoveCard()` - Updates item status in GitHub
- `SaveBoard()` - No-op for GitHub (changes are immediate)

## Usage Examples

### For Users
```bash
# Use local YAML files (default)
./tkan

# Use GitHub Project
./tkan --github GGPrompts/7

# Get help
./tkan --help
```

### For Claude
```bash
# See what needs to be done
gh project item-list 7 --owner GGPrompts --format json | jq '.items[] | select(.content.status == "Todo") | .content.title'

# Move a task to In Progress
gh api --method PATCH /users/GGPrompts/projectsV2/7/items/{item_id} \
  --field field_id=233495315 \
  --field value=47fc9ee4

# Create a bug report
gh project item-create 7 --owner GGPrompts \
  --title "Fix: GitHub backend field mapping" \
  --body "The Status field ID needs to be dynamically fetched"
```

## Future Enhancements

1. **Dynamic Field Discovery**: Fetch field IDs at runtime instead of hardcoding
2. **Status Caching**: Cache project metadata to reduce API calls
3. **Bidirectional Sync**: Support offline mode with periodic sync
4. **Custom Fields**: Support additional GitHub Project fields
5. **Auto-refresh**: Poll for changes from other clients
6. **Error Recovery**: Better handling of API rate limits and failures

## Testing the Integration

1. Ensure GitHub CLI is authenticated:
   ```bash
   gh auth status
   gh auth refresh -h github.com -s project  # Add project scope if needed
   ```

2. Test the connection:
   ```bash
   gh project list --owner @me
   ```

3. Run tkan with GitHub backend:
   ```bash
   ./tkan --github GGPrompts/7
   ```

## Troubleshooting

### Common Issues
- **"Token has not been granted required scopes"**: Run `gh auth refresh -h github.com -s project`
- **"Project not found"**: Verify project number and owner are correct
- **"Interface conversion error"**: Check that project has standard Status field
- **Rate limiting**: GitHub API has rate limits; implement caching to reduce calls

### Debug Commands
```bash
# Check authentication
gh auth status

# List your projects
gh project list --owner @me

# View project fields
gh api /users/{owner}/projectsV2/{number}/fields

# View raw item data
gh api /users/{owner}/projectsV2/{number}/items
```

## TUI Development Skills

This project includes specialized Claude skills for TUI development in `.claude/skills/`:

### Available TUI Skills

1. **tui-add-keybinding.md** - Add keyboard shortcuts with proper handlers
   - Updates key handler, help text, and UI
   - Example: `/tui-add-keybinding "t" "toggle-view" "Switch between board and table views"`

2. **tui-add-tab.md** - Add new tabs to the UI
   - Creates tab components with navigation
   - Useful for adding new views or sections

3. **tui-dynamic-panel.md** - Create resizable panels
   - Implements weight-based panel layouts
   - Perfect for detail views and sidebars

4. **tui-fix-responsive.md** - Fix responsive layout issues
   - Handles terminal resizing properly
   - Ensures UI adapts to different sizes

5. **tui-new-layout.md** - Create new layout patterns
   - Build accordion, dual-pane, or custom layouts
   - Based on TUITemplate patterns

6. **bubbletea/** - Comprehensive Bubble Tea framework skill
   - Complete framework documentation
   - Component patterns and examples
   - Scripts for common tasks

### Using TUI Skills

These skills understand the Bubble Tea architecture used in tkan:
- Model-View-Update pattern
- Lipgloss styling
- Component composition
- Event handling

When adding new UI features, reference the appropriate skill:
```
"Use the tui-add-keybinding skill to add 't' for toggling table view"
"Apply the tui-dynamic-panel skill to make the detail panel resizable"
```

### TUI Architecture Files

Key files for TUI modifications:
- `model.go` - Core model and state management
- `view.go` - Rendering and UI components
- `update.go` - Main update handler
- `update_keyboard.go` - Keyboard event handling
- `update_mouse.go` - Mouse/drag interactions
- `styles.go` - Lipgloss styling definitions

## Contributing

When adding new features that Claude should know about:
1. Update this file with new commands/endpoints
2. Include field IDs for any new projects
3. Add examples of how Claude can use the feature
4. Link to relevant GitHub documentation
5. Use TUI skills for UI modifications
6. Update skills if new patterns emerge

---

*Last Updated: 2024-10-28*
*Project: tkan - Terminal Kanban Board*
*GitHub Project: #7 (tkan Development)*
*TUI Framework: Bubble Tea + Lipgloss*