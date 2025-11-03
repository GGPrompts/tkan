# tkan Hotkeys & Commands

Quick reference for Terminal Kanban Board keyboard shortcuts and common workflows.

## 🎯 Navigation

### Column Navigation
```
← / h - Move left
→ / l - Move right
↑ / k - Move up
↓ / j - Move down
```

### Quick Jumps
```
1 - Jump to BACKLOG column
2 - Jump to TODO column
3 - Jump to PROGRESS column
4 - Jump to REVIEW column
5 - Jump to DONE column
```

### Panel Navigation
```
Tab - Toggle between board and detail panel
```

## 🃏 Card Operations

### Create & Edit
```
n - New card (in current column)
e - Edit selected card
Enter - Edit selected card (alternative)
```

### Move Cards
```
m - Move card to different column (opens picker)
Space - Mark/unmark card (for batch operations)
```

### Delete & Archive
```
d - Delete selected card
a - Toggle archive column visibility
A - Archive selected card
```

### Card Details
```
i - Show card info in detail panel
```

## 🖱️ Mouse Operations

### Drag & Drop
```
Click and drag - Move card within column or across columns
Drop - Release to place card at green indicator line
```

### Click Actions
```
Single click - Select card
Double click - Edit card
Right click - Context menu (coming soon)
```

## 🎨 View Modes

### Display Options
```
v - Toggle table view (coming in Phase 3)
b - Return to board view
```

### Detail Panel
```
Tab - Toggle detail panel
Ctrl+D - Toggle detail panel (alternative)
```

## 📁 Project Management

### Multi-Project
```
p - Project picker (when multiple .tkan.yaml files found)
Ctrl+P - Force project picker
```

### File Operations
```
s - Save current state (auto-saves on changes)
r - Reload from disk
```

## 🔍 Search & Filter (Coming Soon)

```
/ - Search cards
f - Filter by tag
@ - Filter by assignee
# - Filter by priority
```

## ⌨️ Common Workflows

### Quick Add Workflow
```bash
n                    # Create new card
# Type title, press Enter
# Type description, press Ctrl+S
# Card appears in current column
```

### Move Card Workflow
```bash
j/k                  # Navigate to card
m                    # Open move menu
→ or l               # Select target column
Enter                # Confirm move
```

### Drag & Drop Workflow
```bash
# Click on card
# Hold and drag to new position
# Green line shows drop position
# Release to drop
```

### Archive Old Cards
```bash
5                    # Jump to DONE column
j/k                  # Navigate to completed card
A                    # Archive card
a                    # Toggle archive view to verify
```

### Review Mode Workflow
```bash
4                    # Jump to REVIEW column
↓                    # Select card to review
i                    # Show details in panel
Tab                  # Focus detail panel
# Read full description
Tab                  # Back to board
5 or m               # Move to DONE if approved
```

## 🚀 Power User Tips

### Keyboard-Only Navigation
```bash
# Never touch mouse:
hjkl                 # Vim-style movement
12345                # Jump to columns
nm                   # Create and move
Tab                  # Toggle panels
```

### Batch Operations
```bash
Space                # Mark first card
↓ Space              # Mark second card
↓ Space              # Mark third card
# (Batch move/delete coming soon)
```

## 🎯 GitHub Projects Integration

### Sync Commands (Coming Soon)
```bash
g s                  # Sync with GitHub Projects
g p                  # Push changes to GitHub
g f                  # Fetch updates from GitHub
```

### Card Metadata
```yaml
# .tkan.yaml format for GitHub integration:
cards:
  - id: 1
    title: "Fix login bug"
    description: "Users can't authenticate"
    column: "TODO"
    tags: ["bug", "p1"]
    assignee: "@alice"
    gh_issue: 42        # Links to GitHub issue #42
    gh_project: true    # Synced with GH Projects
```

## ⚙️ Configuration

### YAML Structure
```yaml
# example.tkan.yaml
project: "My Project"
columns:
  - BACKLOG
  - TODO
  - PROGRESS
  - REVIEW
  - DONE
cards:
  - id: 1
    title: "Example task"
    column: "TODO"
    tags: ["feature"]
```

## 🛠️ Troubleshooting

### Reset View
```
Ctrl+R - Refresh/reload board
```

### Exit
```
q - Quit tkan
Ctrl+C - Force quit
Esc - Cancel current operation
```

## 📝 Notes

- Cards are saved to `.tkan.yaml` in the project root
- Drag & drop works with mouse or trackpad
- Keyboard navigation uses Vim-style hjkl keys
- Archive column is hidden by default (press 'a' to show)
- Detail panel shows full card metadata
- Ghost cards show where dragged cards came from

---

**Version**: tkan v1.0
**Last Updated**: 2024-11-02
