# tkan Changelog

## v0.5.1 - 2025-11-13

### GitHub Projects Multi-Project Support

**New Features:**
- ✅ **List GitHub Projects** - `--github-owner` flag to list all projects from an owner
- ✅ **Switch Between Projects** - Use 'p' key to switch between GitHub projects
- ✅ **Dynamic Backend Switching** - Automatically switches backend when changing projects
- ✅ **@me Support** - Use `@me` to list your own GitHub projects

**Usage:**
```bash
# List all your GitHub projects (switchable with 'p')
./tkan --github-owner @me

# List projects from a specific owner
./tkan --github-owner GGPrompts

# Use a specific project (single project mode)
./tkan --github GGPrompts/7
```

**Technical Details:**
- Added `ListGitHubProjects()` function to scan GitHub projects
- Updated `loadSelectedProject()` to handle GitHub project paths
- GitHub projects stored in format: `github:owner/project-number`
- Backend dynamically created when switching projects

**Files Modified:**
- backend_github.go: Added `ListGitHubProjects()` and `GitHubProjectInfo` type
- main.go: Added `--github-owner` flag and multi-project support
- model.go: Updated `loadSelectedProject()` to switch backends

---

## v0.5.0 - 2025-11-13

### Card Management & UI Polish

**New Features:**
- ✅ **Help Screen** - Press `?` to view comprehensive keyboard shortcuts and controls
- ✅ **Card Creation** - Press `n` to create new cards with title and description
- ✅ **Card Editing** - Press `e` to edit existing cards (title and description)
- ✅ **Card Deletion** - Press `d` to delete selected cards
- ✅ **Form Modal** - Beautiful centered modal for card creation/editing
  - Tab/↑/↓ to navigate between fields
  - Enter to move to next field (saves on last field)
  - Ctrl+S or Ctrl+Enter to save
  - Esc to cancel

**UI Fixes:**
- ✅ Fixed column header alignment (removed extra padding)
- ✅ Fixed duplicate divider before detail panel
- ✅ Cards now center-aligned within columns to match headers

**Technical Details:**
- Added `github.com/charmbracelet/bubbles/textinput` dependency
- New form modes: `FormCreateCard`, `FormEditCard`
- Form state management in Model
- Unique card ID generation with nanosecond timestamps
- Modal overlay rendering with centered placement

**Keyboard Shortcuts Added:**
- `?` - Toggle help screen
- `n` - Create new card
- `e` - Edit selected card
- `d` - Delete selected card
- In forms: Tab/Enter/Ctrl+S to navigate and save, Esc to cancel

---

## v0.4.0 - 2025-10-29

### GitHub Projects Backend Integration

**New Features:**
- ✅ GitHub Projects (ProjectsV2) as alternative backend
- ✅ `--github` flag to use GitHub Projects instead of local YAML
- ✅ Real-time sync with GitHub web interface
- ✅ Direct API integration via `gh` CLI

**Usage:**
```bash
# Use local YAML files (default)
./tkan

# Use GitHub Project
./tkan --github GGPrompts/7
```

**Technical Details:**
- Backend abstraction interface
- `GitHubBackend` implementation
- Field ID mapping for status updates
- See `CLAUDE.md` for GitHub integration guide

---

## v0.3.0 - 2025-10-28

### Card Display Improvements (Solitaire-style)

**Better Card Rendering:**
- ✅ Card size increased to 12×5 (from 10×4) for better readability
- ✅ Titles now **wrap** instead of truncating with "..."
- ✅ **Removed tags** from card face (tags shown in detail panel only)
- ✅ Solitaire-style stacking: Show only top 2 lines of each card except the last
- ✅ Last card in each column shows full 5 lines
- ✅ Creates cascading visual effect like ~/projects/TUIClassics solitaire
- ✅ "(X more above)" indicator when cards overflow

**Card Format:**
```
┌──────────┐
│Title     │  ← Full title, wrapped across
│wrapped   │     multiple lines (up to 3)
│here      │
└──────────┘
```

**Stacking Effect:**
```
┌──────────  ← Card 1 (top 2 lines only)
│Card 1
┌──────────  ← Card 2 (top 2 lines only)
│Card 2
┌──────────┐ ← Card 3 (full card, last in stack)
│Card 3    │
│          │
└──────────┘
```

**Why This Works:**
- Detail panel shows full card info (description, tags, assignee, etc.)
- Board view focuses on title only - cleaner, easier to scan
- Stacking lets you see more cards in less vertical space
- Wrapped titles eliminate "Implement auth..." truncation

**Technical Details:**
- Borrowed stacking pattern from TUIClassics Solitaire
- `renderCard()` - renders full 5-line card
- `renderCardTopLines()` - renders only top 2 lines for stacking
- Word-wrap algorithm with max 3 lines, ellipsis on overflow

---

## v0.2.0 - 2025-10-28

### New Features

**Column Layout Changes:**
- ✅ Added BACKLOG column on the far left (for ideas/future work)
- ✅ Archive column now toggleable with 'a' key (hidden by default)
- ✅ Default columns: BACKLOG | TODO | PROGRESS | REVIEW | DONE
- ✅ Press 'a' to show/hide ARCHIVE column

**Multi-Project Support:**
- ✅ Project scanning: Automatically finds all `.tkan.yaml` files in current directory and subdirectories
- ✅ Project list view: If multiple projects found, starts with selection screen
- ✅ Navigate projects: ↑/↓ to select, Enter to open
- ✅ Switch projects: Press 'p' to return to project list from board view
- ✅ Smart display: Shows only board if single project found

**Keyboard Shortcuts:**

*Project List View:*
- `↑/↓` or `k/j` - Navigate projects
- `Enter` - Open selected project
- `q` - Quit

*Board View:*
- `←/→` or `h/l` - Navigate columns
- `↑/↓` or `k/j` - Navigate cards
- `Tab` - Toggle detail panel
- `a` - Toggle archive column visibility
- `p` - Return to project list (if multiple projects)
- `v` - Switch to table view (not yet implemented)
- `q` - Quit

### Technical Details

**New Files:**
- `projects.go` - Project scanning and discovery

**Updated Files:**
- `types.go` - Added ViewProjectList mode, Project type, archive toggle
- `model.go` - Added project management methods
- `view.go` - Added project list rendering
- `update_keyboard.go` - Added project list navigation
- `persistence.go` - Updated default board with BACKLOG
- `main.go` - Added project scanning on startup

**Code Stats:**
- 10 Go files
- ~1,100 lines of code

---

## v0.1.0 - 2025-10-28

### Initial Release

**Core Features:**
- ✅ Kanban board view with 5 columns
- ✅ 10×4 char cards with title and tags
- ✅ Toggleable detail panel (33% width)
- ✅ Keyboard navigation (arrows)
- ✅ YAML persistence (.tkan.yaml)
- ✅ Sample board with demo cards

**Architecture:**
- Built with Go, Bubbletea, Lipgloss
- Weight-based layout system (from TUITemplate)
- Card drag patterns (from Solitaire - Phase 2)
- Sortable headers (from TFE - Phase 2)
