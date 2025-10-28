# tkan - Terminal Kanban Board

**A dual-view task management TUI with visual board and sortable table views**

Created: 2025-10-28

---

## Overview

**tkan** (Terminal Kanban) is a terminal-based Kanban board application that combines the visual appeal of draggable cards with the power of sortable table views. Built with Go, Bubbletea, and Lipgloss, it provides a modern, keyboard-driven task management experience.

### Core Features

- **Dual View System**: Toggle between visual board and data table
- **Drag & Drop**: Solitaire-style card dragging between columns
- **Sortable Headers**: Click table headers to sort by any field
- **Detail Panel**: Always-visible card details with full metadata
- **Project-Specific**: Each project has its own `.tkan.yaml` board
- **Keyboard-First**: Full keyboard navigation with mouse support
- **Filtering**: Search by text, filter by tags, assignees, or columns

### Design Philosophy

1. **Best of Both Worlds**: Visual board for daily work, table view for planning
2. **Proven Patterns**: Reuse battle-tested components from TUITemplate, TFE, and Solitaire
3. **Local-First**: No servers, no cloud - just YAML files in your project
4. **Fast**: Optimized rendering, instant switching between views
5. **Beautiful**: Carefully designed UI with consistent styling

---

## Visual Design

### View 1: Board View (Primary)

```
┌──────────────────────────────────────────────────────────────────┬─────────────────────┐
│  📋 tkan - ~/projects/MyProject                    Board View    │ ▶ CARD DETAILS      │
├────────────────────────────────────────────────────┬──────────────┤                     │
│  TODO   │ PROGRESS │ REVIEW  │   DONE    │ ARCHIVE │              │ Fix login flow      │
│   (3)   │   (2)    │  (1)    │   (5)     │  (12)   │              │ ━━━━━━━━━━━━━━━━━━━ │
│ ┌──────┐ ┌──────┐  ┌──────┐  ┌──────┐   ┌──────┐  │              │                     │
│ │Fix   │ │Add   │  │Review│  │Setup │   │Old   │  │              │ Description:        │
│ │login │ │auth  │  │PR#42 │  │DB    │   │stuff │  │              │ Users can't auth    │
│ │#p1   │ │#feat │  │#code │  │#done │   │#done │  │              │ via OAuth. Error    │
│ └──────┘ └──────┘  └──────┘  └──────┘   └──────┘  │              │ 401 on refresh...   │
│ ┌──────┐ ┌──────┐            ┌──────┐              │              │                     │
│ │Add   │ │Docs  │            │Tests │              │              │ Tags: #bug #p1      │
│ │tests │ │      │            │      │              │              │ Assigned: @alice    │
│ │#test │ │      │            │      │              │              │ Due: Jan 15 (5d)    │
│ └──────┘ └──────┘            └──────┘              │              │ Created: Jan 1      │
│                                                     │              │ Modified: Jan 10    │
├─────────────────────────────────────────────────────┴──────────────┤                     │
│ Drag cards | Click: Select | V: Table | Tab: Detail | n: New      │ [E]dit [M]ove      │
└──────────────────────────────────────────────────────────────────────┴─────────────────────┘
```

**Features:**
- Compact 10×4 char cards (title + tags only)
- 5 configurable columns (default: TODO, PROGRESS, REVIEW, DONE, ARCHIVE)
- Drag cards between columns with mouse
- Keyboard navigation: ←/→ for columns, ↑/↓ for cards
- Detail panel shows full card info (33% width, toggleable with Tab)

### View 2: Table View (V to toggle)

```
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│  📋 tkan - ~/projects/MyProject                          Table View [V to switch]       │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│  Title ↓              Column      Tags        Assignee   Due Date   Created    Modified │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│  Fix login flow       TODO        #bug #p1    @alice     Jan 15     Jan 1      Jan 10   │
│  Add OAuth support    PROGRESS    #feature    @bob       Jan 20     Jan 5      Jan 11   │
│  Review PR #42        REVIEW      #code       @charlie   Jan 18     Jan 8      Jan 12   │
│  Setup database       DONE        #infra      @dave      -          Dec 28     Jan 2    │
│  Write tests          TODO        #test       @alice     Jan 22     Jan 3      Jan 9    │
│  Add documentation    PROGRESS    #docs       @eve       Jan 25     Jan 6      Jan 11   │
│  Deploy to prod       TODO        #ops        @bob       Jan 30     Jan 7      Jan 8    │
│  Fix CSS bug          DONE        #bug        @alice     -          Jan 1      Jan 5    │
│  Refactor auth        ARCHIVE     #tech-debt  @charlie   -          Dec 15     Dec 20   │
│  Update README        DONE        #docs       @eve       -          Jan 2      Jan 4    │
│                                                                                           │
│ [Showing 10 of 23 cards]                                                                 │
├─────────────────────────────────────────────────────────────────────────────────────────┤
│ Click headers to sort | Enter: Details | /: Filter | n: New | V: Board view            │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

**Features:**
- Sortable columns (click header to sort, click again to reverse)
- All card metadata visible at once
- Better for long lists (archive with 100+ cards)
- Search/filter across all fields
- Bulk operations (future: select multiple cards)

---

## Architecture

### Technology Stack

**Core:**
- **Language**: Go 1.21+
- **TUI Framework**: [Bubbletea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lipgloss](https://github.com/charmbracelet/lipgloss)
- **Configuration**: YAML (gopkg.in/yaml.v3)

**Optional:**
- **Forms**: [Huh](https://github.com/charmbracelet/huh) for card editing
- **Glamour**: Markdown rendering in card descriptions

### Project Structure

```
tkan/
├── main.go              # Entry point (minimal, ~20 lines)
├── types.go             # Data structures (Card, Column, Model)
├── model.go             # Model initialization, layout calculations
├── view.go              # Main view dispatcher
├── view_board.go        # Board view rendering
├── view_table.go        # Table view rendering (TFE pattern)
├── view_detail.go       # Detail panel rendering
├── update.go            # Update dispatcher
├── update_keyboard.go   # Keyboard handling
├── update_mouse.go      # Mouse handling (drag, header clicks)
├── styles.go            # Lipgloss style definitions
├── config.go            # YAML config management
├── persistence.go       # Load/save .tkan.yaml files
├── helpers.go           # Utilities (truncate, format, etc.)
├── go.mod
├── go.sum
├── README.md
└── PLAN.md              # This file
```

### Component Integration

| Component | Source Project | Purpose |
|-----------|----------------|---------|
| Weight-based layout | TUITemplate | Dynamic column widths, detail panel sizing |
| Card drag system | Solitaire | Click detection, drag distance calculation |
| Sortable headers | TFE | Column header rendering, click → sort |
| Detail panel | TUITemplate | Dual-pane layout with accordion mode |
| Persistence | New | YAML storage per project directory |

---

## Data Model

### Core Types (types.go)

```go
package main

import (
    "time"
    "github.com/charmbracelet/bubbles/spinner"
    "github.com/charmbracelet/lipgloss"
)

// Card represents a single task/card
type Card struct {
    ID          string     // UUID
    Title       string     // Short title (required)
    Description string     // Long description (markdown)
    Column      string     // Current column name
    Tags        []string   // #bug, #feature, #p1, etc.
    Assignees   []string   // @alice, @bob, etc.
    DueDate     *time.Time // Optional due date
    Created     time.Time  // Creation timestamp
    Modified    time.Time  // Last modified timestamp
}

// Column represents a board column
type Column struct {
    Name  string          // Column display name
    Cards []Card          // Cards in this column
    Color lipgloss.Color  // Column color theme
}

// Board represents the entire Kanban board
type Board struct {
    Name        string   // Board name (derived from project)
    ProjectPath string   // Absolute path to project directory
    Columns     []Column // Ordered list of columns
}

// ViewMode represents the current view
type ViewMode int

const (
    ViewBoard ViewMode = iota // Visual Kanban board
    ViewTable                  // Sortable data table
)

// Model represents the application state
type Model struct {
    // Terminal dimensions
    width  int
    height int

    // Board data
    board Board

    // View state
    viewMode      ViewMode
    showDetails   bool // Toggle detail panel
    detailWidth   int  // Detail panel width (dynamic)

    // Board view state
    focusedColumn int // Which column has focus
    focusedCard   int // Which card in column has focus

    // Table view state (TFE pattern)
    cursor     int    // Selected row in table
    sortBy     string // "title", "column", "assignee", "due", "created", "modified"
    sortAsc    bool   // Sort direction
    tableScroll int   // Vertical scroll offset

    // Mouse drag state (Solitaire pattern)
    draggingCard   *Card
    draggingCards  []Card // For future: drag multiple cards
    dragFromColumn int
    dragFromIndex  int
    mousePressX    int
    mousePressY    int

    // Double-click detection
    lastClickTime time.Time
    lastClickX    int
    lastClickY    int

    // Selection state
    selectedCard   *Card
    selectedColumn int
    selectedIndex  int

    // Search/Filter
    searchMode    bool
    searchQuery   string
    filterTags    []string   // Active tag filters
    filterAssignees []string // Active assignee filters

    // Dialog state
    showDialog    bool
    dialogType    string // "edit", "new", "confirm", "message"
    dialogInput   string

    // Status message
    statusMsg   string
    statusError bool
    statusTime  time.Time

    // Loading spinner
    spinner spinner.Model
    loading bool
}
```

### File Format (.tkan.yaml)

```yaml
name: MyProject Kanban
project_path: /home/user/projects/myproject
columns:
  - name: TODO
    color: "#FFA500"  # Orange
    cards:
      - id: card-001
        title: Fix login flow
        description: |
          Users unable to authenticate via OAuth.
          Getting 401 error on token refresh endpoint.

          Need to investigate token expiration handling.
        column: TODO
        tags:
          - bug
          - p1
          - auth
        assignees:
          - alice
        due_date: 2024-01-15T00:00:00Z
        created: 2024-01-01T10:00:00Z
        modified: 2024-01-10T15:30:00Z

  - name: PROGRESS
    color: "#4169E1"  # Royal Blue
    cards: []

  - name: REVIEW
    color: "#9370DB"  # Medium Purple
    cards: []

  - name: DONE
    color: "#32CD32"  # Lime Green
    cards: []

  - name: ARCHIVE
    color: "#808080"  # Gray
    cards: []

# Optional: board settings
settings:
  card_width: 10   # Characters
  card_height: 4   # Lines
  max_columns: 6   # Max visible columns before scrolling
```

---

## Implementation Phases

### Phase 1: Foundation (Week 1) ✅ **COMPLETED - 2025-10-28**

**Goal**: Set up project structure, data model, basic board view

**What We Built**:
- ✅ Full Kanban board view with BACKLOG, TODO, PROGRESS, REVIEW, DONE columns
- ✅ 12×5 char cards with wrapped titles (Solitaire-style stacking)
- ✅ Toggleable detail panel (Tab key, 33% width)
- ✅ Toggleable ARCHIVE column ('a' key, hidden by default)
- ✅ Multi-project support with project selector
- ✅ Keyboard navigation (←→↑↓, g/G, p for projects)
- ✅ YAML persistence (.tkan.yaml files)
- ✅ Project scanning (finds all .tkan.yaml in subdirectories)
- ✅ Card stacking (borrowed from ~/projects/TUIClassics/solitaire)

**Files Created** (1,338 lines total):
- main.go, types.go, model.go, view.go, styles.go
- update.go, update_keyboard.go, update_mouse.go
- persistence.go, projects.go

**Tasks:**
1. ✅ Create project with TUITemplate structure
2. ✅ Define types (Card, Column, Board, Model)
3. ✅ Implement YAML persistence (load/save)
4. ✅ Create basic board view with static columns
5. ✅ Render cards as simple bordered boxes

**Deliverable**: ✅ Can view a board with hardcoded cards (and much more!)

### Phase 2: Board View - Drag System (Week 2) 🚧 **NEXT**

**Goal**: Implement Solitaire-style card dragging

**Pattern from Solitaire** (update_mouse.go):
```go
// On press: track position + prepare drag
m.mousePressX, m.mousePressY = msg.X, msg.Y
m.draggingCard = card
m.dragFromColumn = columnIdx
m.dragFromIndex = cardIdx

// On release: distance check
dx := msg.X - m.mousePressX
dy := msg.Y - m.mousePressY
if dx*dx + dy*dy < 4 {
    // Click (select card)
} else {
    // Drag (move card to new column)
    toCol := m.getColumnAtPosition(msg.X, msg.Y)
    m.moveCard(fromCol, toCol, card)
}
```

**Tasks:**
1. Implement `getColumnAtPosition(x, y)` coordinate mapping
2. Add mouse press handler (track press position)
3. Add mouse release handler (distance calculation)
4. Implement `moveCard()` logic
5. Add visual feedback during drag (highlight drop zone)

**Deliverable**: Can drag cards between columns with mouse

### Phase 3: Table View (Week 3) 📅 **PLANNED**

**Goal**: Implement TFE-style table view with sortable headers

**Pattern from TFE** (render_file_list.go):
```go
// Calculate dynamic column widths
titleWidth := availableWidth * 25 / 100
columnWidth := 10
tagsWidth := availableWidth * 15 / 100
// ...

// Render header with sort indicator
sortIndicator := ""
if m.sortAsc {
    sortIndicator = " ↑"
} else {
    sortIndicator = " ↓"
}

titleHeader := "Title"
if m.sortBy == "title" {
    titleHeader += sortIndicator
}

header := fmt.Sprintf("%-*s  %-*s  ...", titleWidth, titleHeader, ...)
```

**Tasks:**
1. Create `renderTableView()` function
2. Calculate dynamic column widths
3. Render header with sort indicators
4. Render card rows with truncation
5. Implement `sortCards()` function

**Deliverable**: Can view cards in sortable table format

### Phase 4: Sortable Headers - Click Detection (Week 4)

**Goal**: Click table headers to sort columns

**Pattern from TFE** (update_mouse.go:528):
```go
if m.viewMode == ViewTable && msg.Y == headerY {
    adjustedX := msg.X - 2  // Account for border

    // Calculate column ranges
    titleEnd := titleWidth
    columnEnd := titleEnd + columnWidth + 2
    tagsEnd := columnEnd + tagsWidth + 2
    // ...

    var newSortBy string
    if adjustedX >= 2 && adjustedX <= titleEnd {
        newSortBy = "title"
    } else if adjustedX > titleEnd && adjustedX <= columnEnd {
        newSortBy = "column"
    }
    // ... (check all columns)

    if newSortBy != "" {
        if m.sortBy == newSortBy {
            m.sortAsc = !m.sortAsc  // Toggle direction
        } else {
            m.sortBy = newSortBy    // New column
            m.sortAsc = true
        }
        m.sortCards()
    }
}
```

**Tasks:**
1. Add header click detection in `handleMouseEvent()`
2. Calculate column X ranges (must match renderTableView)
3. Implement sort toggle logic
4. Re-render table after sort

**Deliverable**: Can click headers to sort in table view

### Phase 5: Detail Panel (Week 5)

**Goal**: Add detail panel showing full card information

**Pattern from TUITemplate** (dual-pane layout):
```go
func (m Model) calculateDetailLayout() (boardWidth, detailWidth int) {
    contentWidth, _ := m.calculateLayout()

    leftWeight, rightWeight := 2, 1  // 66% board, 33% details

    if m.accordionMode && m.focusedPane == "details" {
        leftWeight, rightWeight = 1, 1  // 50/50 when focused
    }

    totalWeight := leftWeight + rightWeight
    boardWidth = (contentWidth * leftWeight) / totalWeight
    detailWidth = contentWidth - boardWidth - 1  // -1 for divider

    return boardWidth, detailWidth
}
```

**Tasks:**
1. Create `renderDetailPanel()` function
2. Implement weight-based layout calculation
3. Show card title, description (word-wrapped), metadata
4. Add action buttons ([E]dit, [M]ove, [D]elete, [C]opy)
5. Support Tab key to toggle panel visibility

**Deliverable**: Detail panel shows full card info, toggleable with Tab

### Phase 6: Keyboard Navigation (Week 6)

**Goal**: Full keyboard control for both views

**Board View:**
- `←/→` or `h/l`: Move focus between columns
- `↑/↓` or `k/j`: Move focus between cards in column
- `Enter`: Select card (show in detail panel)
- `m`: Move card to column (show column picker)
- `e`: Edit card (open dialog)
- `d`: Delete card (confirm)
- `n`: New card in current column

**Table View:**
- `↑/↓` or `k/j`: Navigate rows
- `Enter`: Select card (show in detail panel)
- Same actions as board view (e, d, n, m)

**Both Views:**
- `v`: Toggle between board and table view
- `Tab`: Toggle detail panel visibility
- `/`: Search/filter mode
- `#`: Filter by tag (show tag picker)
- `@`: Filter by assignee (show assignee picker)
- `Esc`: Clear filters / close dialogs
- `q`: Quit

**Tasks:**
1. Implement keyboard handler for board view
2. Implement keyboard handler for table view
3. Add view toggle (v key)
4. Add search mode (/ key)
5. Add filter modes (#, @ keys)

**Deliverable**: Full keyboard navigation in both views

### Phase 7: Card Editing & Dialogs (Week 7)

**Goal**: Add/edit/delete cards with dialogs

**Using Charm's Huh library:**
```go
import "github.com/charmbracelet/huh"

func (m Model) openEditDialog() (tea.Model, tea.Cmd) {
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Title").
                Value(&m.selectedCard.Title),
            huh.NewText().
                Title("Description").
                CharLimit(500).
                Value(&m.selectedCard.Description),
            huh.NewInput().
                Title("Tags (space-separated)").
                Value(&tagsString),
            huh.NewInput().
                Title("Assignees (space-separated)").
                Value(&assigneesString),
        ),
    )

    // Show modal form...
}
```

**Tasks:**
1. Add Huh dependency
2. Create edit dialog (title, description, tags, assignees, due date)
3. Create new card dialog (same fields + column picker)
4. Create delete confirmation dialog
5. Create move dialog (column picker)
6. Save changes to YAML on edit

**Deliverable**: Can create, edit, and delete cards

---

## Future Enhancements

### Phase 8+: Advanced Features

1. **Undo/Redo Stack**
   - Track all moves/edits
   - Ctrl+Z to undo, Ctrl+Y to redo

2. **Multi-Select Cards**
   - Space to select/deselect
   - Move multiple cards at once
   - Bulk tag/assignee updates

3. **Card History/Activity Log**
   - Track all changes to a card
   - Show in detail panel as timeline

4. **Custom Columns**
   - Add/remove/rename columns
   - Reorder columns with drag
   - Set column colors

5. **Filtering Enhancements**
   - Regex search
   - Date range filters (due this week, overdue, etc.)
   - Save filter presets

6. **Reporting/Analytics**
   - Cycle time per card
   - Lead time metrics
   - Burndown charts (ASCII art)
   - Export to CSV/JSON

7. **Multi-Project Management**
   - `tkan --list` to see all projects
   - Quick switch between projects
   - Dashboard view showing all boards

8. **Collaboration Features**
   - Git-based sync (commit .tkan.yaml)
   - Conflict resolution on pull
   - Optional: lightweight HTTP server for real-time sync

9. **Integrations**
   - Import from GitHub Issues
   - Import from Jira
   - Export to markdown checklist
   - Webhook on card state changes

10. **Theming**
    - Customizable color schemes
    - Config file for styles
    - Dark/light mode presets

---

## Key Design Decisions

### Why Dual Views?

**Board View** is great for:
- Daily work (moving cards through workflow)
- Visual overview of project state
- Quick status checks during standups

**Table View** is great for:
- Planning sessions (sorting by priority, due date)
- Finding specific cards (search, filter, sort)
- Reviewing large backlogs (archive with 100+ cards)
- Bulk operations (select multiple cards)

Having both gives users the right tool for each task.

### Why Side Detail Panel?

Alternatives considered:
- **Modal dialog**: Blocks view, requires dismissing
- **Bottom panel**: Wastes horizontal space on wide terminals
- **No panel (inline in table)**: Limited space for description

**Side panel wins because:**
- Always visible (no switching)
- Works in both views
- Doesn't interrupt workflow
- Can show full markdown description
- Easy to toggle (Tab key)

### Why .tkan.yaml Files?

Alternatives considered:
- **SQLite database**: Overkill for simple task lists
- **JSON**: Less human-readable/editable
- **Markdown files**: Hard to preserve metadata, ordering

**YAML wins because:**
- Human-readable and editable
- Git-friendly (readable diffs)
- Easy to parse and generate
- Standard for config files
- Can version control with git

### Why 10×4 Cards in Board View?

- **10 chars wide**: Fits 5 columns + detail panel on 80-col terminal
- **4 lines tall**: Shows title (1-2 lines) + tags (1 line)
- **Compact**: See many cards at once without scrolling
- **Detail panel compensates**: Full info always visible on right

Larger cards (15×5 or 20×7) tested poorly:
- Only 2-3 columns fit on screen
- More scrolling required
- Less "board overview" feel

### Why Solitaire's Drag Pattern?

The Solitaire drag system (update_mouse.go) is perfect because:
- **Distance-based detection**: `dx² + dy² < 4` differentiates click vs drag
- **Proven**: Already works great in Solitaire
- **Simple**: Just 3 state variables (press position, dragging flag)
- **No ghost cards needed**: Terminal updates fast enough

Other drag patterns (continuous tracking, ghost cards) were more complex for no benefit.

### Why TFE's Table Pattern?

TFE's detail view (render_file_list.go) is ideal because:
- **Dynamic column widths**: Adapts to terminal size
- **Sortable headers**: Click to sort already implemented
- **Clean**: Simple ASCII table, no fancy borders
- **Fast**: Renders 100+ rows smoothly

Alternative (Bubble Table library) considered but:
- Added dependency
- Less flexible (harder to customize)
- TFE pattern proven to work great

---

## Technical Notes

### Layout Calculations (from TUITemplate/CLAUDE.md)

**Critical: Always account for borders in height calculations**

```go
func (m model) calculateLayout() (int, int) {
    contentWidth := m.width
    contentHeight := m.height

    if m.config.UI.ShowTitle {
        contentHeight -= 3 // title bar (3 lines)
    }
    if m.config.UI.ShowStatus {
        contentHeight -= 1 // status bar
    }

    // CRITICAL: Account for panel borders
    contentHeight -= 2 // top + bottom borders

    return contentWidth, contentHeight
}
```

**Visual layout:**
```
┌─────────────────────────────────┐  ← Title Bar (3 lines)
│  App Title                      │
│  Subtitle/Info                  │
├─────────────────────────────────┤  ─┐
│ ┌─────────────┬───────────────┐ │   │
│ │             │               │ │   │
│ │   Board     │    Detail     │ │   │ Content Height
│ │             │               │ │   │ (minus borders)
│ │             │               │ │   │
│ └─────────────┴───────────────┘ │   │
├─────────────────────────────────┤  ─┘
│ Status Bar: Help text here      │  ← Status Bar (1 line)
└─────────────────────────────────┘
```

### Mouse Coordinate Mapping

**Key pattern from Solitaire** (update_mouse.go:250):
```go
func (m *Model) getColumnAtPosition(x, y int) (column, cardIndex int) {
    // Calculate padding to match rendering
    leftPadding := (m.width - totalContentWidth) / 2
    topPadding := (m.height - totalContentHeight) / 2

    // Adjust for padding
    relX := x - leftPadding
    relY := y - topPadding

    // Map to column
    columnWidth := 12  // 10 card + 2 spacing
    column = relX / columnWidth

    // Map to card within column
    cardHeight := 5  // 4 card + 1 spacing
    cardIndex = relY / cardHeight

    return column, cardIndex
}
```

**Must match rendering logic exactly** or clicks will be offset!

### Text Truncation

Always truncate text explicitly - never rely on Lipgloss auto-wrapping inside bordered panels:

```go
func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen-1] + "…"
}

// In rendering:
title := truncate(card.Title, cardWidth-4)  // -2 borders, -2 padding
```

### Sorting Logic

```go
func (m *Model) sortCards() {
    cards := m.getAllCards()  // Flatten all columns

    sort.Slice(cards, func(i, j int) bool {
        switch m.sortBy {
        case "title":
            cmp := strings.Compare(cards[i].Title, cards[j].Title)
            if m.sortAsc {
                return cmp < 0
            }
            return cmp > 0

        case "due":
            // Handle nil due dates
            if cards[i].DueDate == nil {
                return !m.sortAsc  // Nils sort last
            }
            if cards[j].DueDate == nil {
                return m.sortAsc
            }
            if m.sortAsc {
                return cards[i].DueDate.Before(*cards[j].DueDate)
            }
            return cards[j].DueDate.Before(*cards[i].DueDate)

        // ... other fields
        }
    })

    // Update cursor to maintain selection
    // (see TFE update_mouse.go:629)
}
```

---

## Project Goals

### Success Criteria

**Must Have (v1.0):**
- ✅ Board view with draggable cards
- ✅ Table view with sortable headers
- ✅ Detail panel with full card info
- ✅ Keyboard navigation (all actions)
- ✅ Create/edit/delete cards
- ✅ YAML persistence per project
- ✅ Search and filter

**Nice to Have (v1.1+):**
- Undo/redo stack
- Custom columns
- Multi-select cards
- Card history/activity log
- Export to CSV/JSON

**Future (v2.0+):**
- Multi-project dashboard
- Analytics/reporting
- GitHub/Jira integration
- Real-time collaboration

### Non-Goals

- Web interface (terminal only)
- Mobile app (Termux support maybe later)
- Cloud sync (Git is enough)
- Team permissions/access control
- Time tracking (use external tools)

---

## Development Workflow

### Getting Started

```bash
# Create project from TUITemplate
cd ~/projects/TUITemplate
./scripts/new_project.sh

# Name: tkan
# Layout: dual_pane
# Components: panel

# Build and run
cd ~/projects/tkan
go mod tidy
go run .
```

### Testing Strategy

1. **Manual testing** during development (TUI hard to unit test)
2. **Test with various terminal sizes**:
   - Narrow: 80×24 (minimum)
   - Medium: 120×40 (laptop)
   - Wide: 200×60 (desktop)
3. **Test scenarios**:
   - Empty board (no cards)
   - Full board (100+ cards)
   - Long card titles/descriptions
   - Many tags/assignees
   - Drag edge cases (drag to same column, invalid drops)
4. **Performance testing**:
   - 1000+ cards in archive
   - Rapid view switching
   - Sorting large tables

### Release Process

1. Tag version: `git tag -a v1.0.0 -m "Initial release"`
2. Build for multiple platforms:
   ```bash
   GOOS=linux GOARCH=amd64 go build -o tkan-linux-amd64
   GOOS=darwin GOARCH=amd64 go build -o tkan-darwin-amd64
   GOOS=windows GOARCH=amd64 go build -o tkan-windows-amd64.exe
   ```
3. Create GitHub release with binaries
4. Update README with download links

---

## Related Projects

**Inspiration:**
- [lazygit](https://github.com/jesseduffield/lazygit) - Accordion layout, weight-based panels
- [TFE](https://github.com/GGPrompts/tfe) - Sortable headers, detail view
- [TUITemplate](https://github.com/GGPrompts/TUITemplate) - Dual-pane layout, architecture
- [TUIClassics/Solitaire](https://github.com/GGPrompts/TUIClassics) - Card drag system

**Similar Tools:**
- [taskwarrior-tui](https://github.com/kdheepak/taskwarrior-tui) - TUI for TaskWarrior
- [mdt](https://github.com/basilvetas/mdt) - Markdown-based tasks
- [dstask](https://github.com/naggie/dstask) - Git-based task tracking

**Differentiators:**
- ✨ Dual-view system (board + table)
- ✨ Drag & drop cards like Solitaire
- ✨ Sortable table headers like TFE
- ✨ Project-specific boards (not global)
- ✨ Beautiful UI with Lipgloss
- ✨ No external dependencies (no TaskWarrior, no database)

---

## License

MIT License - Use freely for any purpose.

---

**Let's build this! 🚀**

Start with Phase 1 and iterate from there. Each phase builds on the previous one, with working software at every step.
