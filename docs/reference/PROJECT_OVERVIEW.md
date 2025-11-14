# tkan (Terminal Kanban Board) - Comprehensive Project Overview

## Executive Summary

**tkan** is a feature-rich terminal-based Kanban board application built with Go, Bubbletea, and Lipgloss. It provides a dual-view task management experience with visual drag-and-drop cards and planned sortable table views. The project is well-architected, partially complete, and ready for continued development.

**Project Status**: Phase 2 Complete, Phase 3 Planned
**Code Size**: 2,501 lines across 13 Go files
**Latest Update**: November 13, 2025 (GitHub Projects integration + hotkeys)
**License**: MIT

---

## Project Structure & Architecture

### Core Files Overview

| File | Lines | Purpose | Status |
|------|-------|---------|--------|
| `model.go` | 463 | Core state management, layout calculations | ✅ Complete |
| `backend_github.go` | 414 | GitHub Projects API integration | ✅ Complete |
| `view.go` | 396 | Main rendering logic for all views | ✅ Partial |
| `styles.go` | 268 | Lipgloss styling definitions | ✅ Complete |
| `persistence.go` | 152 | YAML load/save operations | ✅ Complete |
| `update_keyboard.go` | 151 | Keyboard event handling | ✅ Partial |
| `update_mouse.go` | 146 | Mouse/drag event handling | ✅ Complete |
| `main.go` | 140 | Entry point, CLI flags, backend selection | ✅ Complete |
| `backend.go` | 111 | Backend interface definition | ✅ Complete |
| `types.go` | 107 | Core data structures (Card, Column, Board, Model) | ✅ Complete |
| `projects.go` | 85 | Project discovery and scanning | ✅ Complete |
| `update.go` | 49 | Update dispatcher | ✅ Complete |
| `update_timer.go` | 19 | Drag delay timer implementation | ✅ Complete |

### Architecture Patterns

The project follows proven patterns from other projects in the ecosystem:

1. **TUITemplate Pattern** - Dual-pane layout with weight-based sizing
   - Uses `calculateLayout()` to compute board vs detail panel widths
   - Dynamic 67%/33% split when details shown, 100% when hidden

2. **Solitaire Pattern** - Card drag mechanics
   - Distance-based drag detection: `dx² + dy² < 4` pixels distinguishes click vs drag
   - Ghost card rendering during drag
   - Drop indicator showing insertion position
   - Drag delay timer (150ms) for responsive feel

3. **Backend Abstraction**
   - Interface-based design allows swapping implementations
   - LocalBackend for `.tkan.yaml` files
   - GitHubBackend for GitHub Projects API

---

## Feature Completeness Analysis

### Phase 1: Foundation ✅ COMPLETED

**Board View Features:**
- ✅ Visual Kanban board with 5 columns (BACKLOG, TODO, PROGRESS, REVIEW, DONE)
- ✅ Solitaire-style stacked cards (12×5 characters)
- ✅ Card titles with word wrapping (no truncation)
- ✅ Toggleable detail panel (Tab key, 33% width, shows full metadata)
- ✅ Toggleable ARCHIVE column ('a' key)
- ✅ Column headers with card counts
- ✅ YAML persistence (`.tkan.yaml` files)
- ✅ Project scanning (auto-discover `.tkan.yaml` in subdirectories)
- ✅ Multi-project support (project list view with navigation)

**Keyboard Navigation:**
- ✅ ←/→ or h/l - Navigate between columns
- ✅ ↑/↓ or k/j - Navigate cards within column
- ✅ g/G - Jump to first/last column
- ✅ Tab - Toggle detail panel visibility
- ✅ a/A - Toggle archive column visibility
- ✅ p/P - Return to project list (multi-project)
- ✅ v/V - Switch to table view (placeholder)
- ✅ q/Ctrl+C - Quit application

**Detail Panel:**
- ✅ Shows selected card's full title
- ✅ Shows description (word-wrapped)
- ✅ Shows tags as colored badges
- ✅ Shows assignee
- ✅ Shows due date
- ✅ Shows creation and modification timestamps

### Phase 2: Drag & Drop ✅ COMPLETED

**Implemented:**
- ✅ Click & drag to move cards between columns
- ✅ Reorder cards within same column
- ✅ Visual drop indicator (green horizontal line)
- ✅ Ghost card effect (dragged card appears faded at source)
- ✅ Smart insertion (hover top/bottom half of card to insert before/after)
- ✅ Empty column support
- ✅ Auto-save on card move
- ✅ Drag delay (150ms) to prevent accidental drags
- ✅ Works across all columns including ARCHIVE

**Technical Implementation:**
- `handleMousePress()` - Starts potential drag, shows immediate selection
- `handleMouseMotion()` - Updates drop target position during drag
- `handleMouseRelease()` - Drops card or cancels drag
- `moveCard()` - Updates data model and saves changes
- `getDropPosition()` - Calculates column and insertion index
- `getInsertIndexInColumn()` - Precise position based on card hover

### Phase 2.5: GitHub Projects Backend ✅ COMPLETED

**Implemented:**
- ✅ GitHub Projects API integration
- ✅ `--github owner/project-number` flag
- ✅ Field mapping (Status ↔ Columns)
- ✅ Load board from GitHub Projects
- ✅ Move card updates GitHub Projects status
- ✅ Backend abstraction (Interface pattern)
- ✅ CLAUDE.md documentation for AI integration
- ✅ Error handling and authentication checks

**Usage:**
```bash
./tkan                    # Local YAML mode
./tkan --github owner/7   # GitHub Project mode
```

### Phase 3: Table View 📅 PLANNED (Not Started)

**Not Implemented:**
- ❌ Table view rendering (TFE-style)
- ❌ Sortable column headers
- ❌ Click-to-sort functionality
- ❌ Table scrolling/pagination
- ❌ Search and filtering in table view

**Partially Visible in Hotkeys but Not Functional:**
- View mode toggle works (v key switches mode)
- Table view placeholder message shown
- `handleTableKeyMsg()` exists but empty

### Phase 4: Card Editing 📅 PLANNED (Partially Started)

**Keyboard Shortcuts Defined But Not Implemented:**
- `n` - New card (exists, does nothing)
- `e` - Edit card (exists, does nothing)
- `d` - Delete card (exists, does nothing)
- `m` - Move card (exists, does nothing)
- `/` - Search/filter (exists, does nothing)
- `?` - Help (exists, does nothing)

**Code Present but Unimplemented:**
```go
case "n":
    // New card
    return m, nil  // TODO: Implement

case "e":
    // Edit card
    return m, nil  // TODO: Implement
```

**Noted in HOTKEYS.md but Nonfunctional:**
- Batch operations (Space to select multiple cards)
- GitHub sync commands (g s, g p, g f)

---

## Data Model

### Core Types (types.go)

```go
type Card struct {
    ID          string    // UUID for card
    Title       string    // Card title (required)
    Description string    // Long description
    Tags        []string  // #bug, #feature, etc.
    Assignee    string    // @username
    DueDate     string    // Due date (YYYY-MM-DD)
    CreatedAt   time.Time // Creation timestamp
    ModifiedAt  time.Time // Last modification
    Column      string    // Which column (TODO, PROGRESS, etc.)
}

type Column struct {
    Name  string   // Column display name
    Cards []*Card  // Cards in this column
}

type Board struct {
    Name        string    // Board/project name
    Description string    // Optional description
    Columns     []Column  // All columns
    Cards       []*Card   // All cards (flat list)
    CreatedAt   time.Time
    ModifiedAt  time.Time
}

type Model struct {
    // Board data
    board          *Board
    backend        Backend  // Persistence implementation
    
    // UI state
    viewMode       ViewMode // Board/Table/ProjectList
    selectedColumn int      // Current column index
    selectedCard   int      // Current card in column
    showDetails    bool     // Detail panel visible
    showArchive    bool     // Archive column visible
    
    // Mouse drag state
    draggingCard   *Card    // Card being dragged
    dragFromColumn int      // Source column
    dragFromIndex  int      // Source position
    dropTargetColumn int    // Where it would drop
    dropTargetIndex  int    // Position in target column
    mousePressX, Y  int     // Press position for drag detection
    
    // Layout
    width, height int       // Terminal dimensions
    boardWidth    int       // Calculated width (67% or 100%)
    detailWidth   int       // Calculated width (33% or 0%)
}
```

### File Format (`.tkan.yaml`)

```yaml
name: My Project
description: A sample Kanban board
columns:
  - name: BACKLOG
  - name: TODO
  - name: PROGRESS
  - name: REVIEW
  - name: DONE
  - name: ARCHIVE
cards:
  - id: "1"
    title: Fix login flow
    description: Users can't authenticate via OAuth...
    tags: [bug, p1]
    assignee: '@alice'
    due_date: "2025-01-15"
    created_at: 2025-10-18T00:00:00Z
    modified_at: 2025-10-28T22:13:06Z
    column: TODO
```

---

## UI/UX Implementation

### Board View Rendering

**Layout (from view.go):**
```
┌─────────────────────────────────────────────────┬──────────────────┐
│ 📋 tkan - My Project                   Board    │ CARD DETAILS     │
├──────────┬──────────┬──────────┬──────┬─────────┤                  │
│BACKLOG(1)│TODO(2)   │PROGRESS(5)│REV(1)│DONE(1)  │ Fix login flow   │
│┌────────┐│┌────────┐│┌────────┐│┌────┐│┌──────┐ │ ━━━━━━━━━━━━━━━━ │
││New      ││Fix     ││Card 1    ││Rev ││Set     │ │                  │
││feature  ││Write   ││Card 2    ││PR   ││DB      │ │ Description:     │
│└────────┘│└────────┘│┌────────┐│└────┘│└──────┐ │ Users can't auth │
│          │          ││Card 3   │      │Set     │ │ via OAuth. Error │
│          │          ││         │      │        │ │ 401 on refresh   │
│          │          │└────────┘│      │└──────┘ │                  │
│          │          │          │      │         │ Tags: #bug #p1   │
│          │          │          │      │         │ @alice           │
├──────────┴──────────┴──────────┴──────┴─────────┤ Jan 15 (5 days)  │
│ ← →: Columns | ↑ ↓: Cards | Tab: Details | q   │                  │
└──────────────────────────────────────────────────┴──────────────────┘
```

**Stacking Effect:**
- Each stacked card (except last) shows only 2 lines of border
- Last card in column shows full 5 lines
- Allows seeing more cards in compact space
- Detail panel compensates with full information on right

### Rendering Functions (styles.go)

- `renderCard()` - Full 5-line card with border
- `renderCardTopLines()` - Top 2 lines for stacking
- `renderCardGhost()` - Faded version during drag
- `renderCardTopLinesGhost()` - Faded top lines

### Styling (Lipgloss)

- **Card**: 12 chars wide × 5 lines (bordered, rounded corners)
- **Selected**: Bright magenta border and text
- **Ghost**: Subdued gray (dragging)
- **Drop indicator**: Green horizontal line
- **Detail panel**: 33% of width, bordered
- **Tags**: Cyan text on dark background with padding
- **Colors**: Mostly professional grays with accent colors

---

## Backend Systems

### Local Backend (YAML Files)

**Features:**
- Loads/saves `.tkan.yaml` files in project directory
- Auto-discovers projects in subdirectories (3 levels max)
- Simple, human-readable format
- No server required
- Version control friendly

**Methods:**
- `LoadBoard()` - Read YAML file
- `SaveBoard()` - Write YAML file
- `MoveCard()` - Update card's column field
- `UpdateCard()` - Replace card in board
- `CreateCard()` - Add new card to board
- `DeleteCard()` - Move to ARCHIVE column

### GitHub Projects Backend

**Features:**
- Fetches project from GitHub via GraphQL + REST APIs
- Maps GitHub Project columns to tkan columns
- Supports field mapping (Status ↔ Column names)
- Live updates through GitHub CLI (`gh` command)
- No local storage needed

**Implementation:**
- Uses `gh project item-list` and REST APIs
- Maps standard column names (Todo, In Progress, Done, etc.)
- Stores GitHub item IDs for sync
- Error handling with authentication checks

**Field Mapping:**
```
GitHub Status  →  tkan Column
"Todo"         →  "TODO"
"In Progress"  →  "PROGRESS"
"In Review"    →  "REVIEW"
"Done"         →  "DONE"
```

---

## Known Issues & Gaps

### Critical Gaps

1. **Table View Not Implemented**
   - View mode toggle works but shows placeholder
   - No rendering of table layout
   - No sortable headers
   - No table-specific keyboard controls
   - TFE pattern documented in PLAN.md but not coded

2. **Card Editing Not Implemented**
   - Keyboard shortcuts exist but do nothing (n, e, d, m)
   - No dialog/form system
   - No input validation
   - No card creation flow
   - No delete confirmation
   - Huh library not imported (form framework)

3. **Search & Filter**
   - `/` key defined but does nothing
   - No search mode state
   - No filter system
   - Filter state exists in Model but unused

4. **GitHub Backend Incomplete**
   - One TODO comment: "Show error message" (update_keyboard.go:59)
   - Field IDs hardcoded (should be dynamic discovery)
   - No rate limit handling
   - No caching of project metadata

### Documentation Gaps

1. **Setup Guide** - COMPLETE_SETUP_GUIDE.md exists (extensive)
2. **Hotkeys** - HOTKEYS.md exists but many commands not implemented
3. **API Reference** - GITHUB_COLUMNS_API_REFERENCE.md provides details
4. **Development Plan** - PLAN.md has implementation details for all phases
5. **Architecture** - CLAUDE.md explains integration for AI

### Code Issues

1. **No Unit Tests** - No `*_test.go` files found
2. **Limited Error Handling** - Some errors silently ignored
3. **Hard-coded Field IDs** - GitHub field IDs not dynamically discovered
4. **Partial Keyboard Bindings** - Many shortcuts defined but not implemented
5. **No Input Validation** - No checks for card data

---

## Recent Changes & Commit History

### Latest Commits (last 10)

1. **9ed7c72** (Nov 13) - Merge branch 'master' (GitHub integration)
2. **1f93563** (Oct 28) - Add comprehensive setup guide and GitHub Projects API reference
3. **7d65036** (Oct 28) - Add comprehensive HOTKEYS.md reference guide
4. **886f602** (Oct 28) - **GitHub Projects backend and drag delay improvements** ✅ Major feature
5. **9df1ff1** (Oct 28) - **Complete Phase 2 - Drag & Drop with Visual Feedback** ✅ Major feature
6. **d7e068a** (Oct 28) - Implement Phase 1 - Foundation complete ✅ Major feature
7. **4afcaff** (Oct 28) - Docs: Add session kickoff prompt for implementation
8. **a7766cf** (Oct 28) - Docs: Add AI card creation system design
9. **a1e12c5** (Oct 28) - Docs: Add project README
10. **7f7dc55** (Oct 28) - Docs: Add comprehensive project plan

### Version Timeline

- **v0.1.0** (Oct 28) - Initial release with Phase 1 foundation
- **v0.2.0** (Oct 28) - Multi-project support and BACKLOG column
- **v0.3.0** (Oct 28) - Card display improvements (Solitaire-style stacking)
- **v0.4.0** (Oct 28) - Phase 2 completed (drag & drop, visual feedback)
- **v0.5.0** (Oct 28) - GitHub Projects backend integration

---

## Development Readiness

### What's Ready to Build Next

**Phase 3: Table View** (Well-documented, clear pattern)
- Pattern from TFE project already documented in PLAN.md
- Layout calculations defined
- Sort logic pseudocode provided
- Easy win for next developer

**Quick Wins** (1-2 hour tasks)
- Implement card creation (n key) with simple dialog
- Implement card deletion (d key) to archive
- Implement move card (m key) to column picker
- Add help screen (? key)

**Medium Tasks** (4-8 hours)
- Table view rendering (copy TFE pattern)
- Sortable headers (click detection + sort)
- Search/filter implementation
- Test with large datasets (100+ cards)

### Testing Strategy Needed

- Manual testing on various terminal sizes (80×24, 120×40, 200×60)
- Drag & drop edge cases (drag to same position, drag across columns)
- Large dataset performance (1000+ cards)
- GitHub API rate limiting handling
- YAML parsing with edge cases (special characters, etc.)

### Performance Considerations

- Current codebase renders efficiently
- No obvious bottlenecks
- Stacking algorithm scales well
- YAML loading is straightforward
- GitHub API calls should be cached

---

## Integration Points

### GitHub Integration Ready

- **CLAUDE.md** - Complete integration guide for Claude AI
- **GitHub CLI** - Can manage tasks via `gh` commands
- **Field Discovery** - Should be implemented (currently hardcoded)
- **Sync Commands** - Documented but not implemented

### TUI Skills Available

- **bubbletea/** - Complete Bubble Tea framework documentation
- **github-projects/** - GitHub Projects API reference
- **Layout patterns** - Weight-based, dual-pane, accordion examples
- **Component examples** - Keyboard/mouse handling patterns

---

## Code Quality Assessment

### Strengths

1. ✅ Clean architecture (Backend interface abstraction)
2. ✅ Well-organized files (One concern per file)
3. ✅ Consistent naming (model, view, update patterns)
4. ✅ Good documentation (PLAN.md, HOTKEYS.md, CLAUDE.md)
5. ✅ Proven patterns (TUITemplate, Solitaire, TFE)
6. ✅ Drag implementation is solid (150ms delay works well)
7. ✅ Layout calculations correct (matches rendering)

### Weaknesses

1. ❌ No tests (no test files found)
2. ❌ TODO comment left in code (update_keyboard.go:59)
3. ❌ Many unimplemented shortcuts (skeleton code)
4. ❌ Limited error messages
5. ❌ No input validation
6. ❌ Hard-coded field IDs (GitHub)
7. ❌ Table view placeholder but no implementation

### Code Statistics

- **Total Lines**: 2,501 across 13 Go files
- **Largest File**: model.go (463 lines)
- **Smallest File**: update_timer.go (19 lines)
- **Average File**: 192 lines
- **Files Fully Complete**: ~8 (70%)
- **Files Partially Complete**: ~4 (30%)
- **Files Not Started**: ~1 (10%)

---

## Next Steps Recommendations

### Immediate (Next Developer)

1. **Implement Table View** (Phase 3)
   - Use TFE pattern from PLAN.md
   - 4-6 hours estimated
   - Clear documentation exists

2. **Implement Card Editing** (Phase 4)
   - Create dialog system (using Huh library)
   - Implement n/e/d/m keyboard shortcuts
   - 6-8 hours estimated

3. **Add Tests**
   - Unit tests for card movements
   - Integration tests for YAML save/load
   - 2-3 hours estimated

### Short Term (This Sprint)

1. **Search & Filter**
   - Implement `/` search mode
   - Filter by tags and assignees
   - 3-4 hours estimated

2. **Field Discovery for GitHub**
   - Dynamically fetch field IDs from GitHub
   - Remove hardcoded values
   - 2 hours estimated

3. **Performance Testing**
   - Test with 1000+ cards
   - Optimize rendering if needed
   - 2-3 hours estimated

### Medium Term (Next Sprints)

1. **Undo/Redo Stack**
   - Track moves and edits
   - Ctrl+Z/Ctrl+Y support
   - 4-6 hours estimated

2. **Custom Columns**
   - Add/remove/rename columns
   - Persist to YAML
   - 4-5 hours estimated

3. **Multi-Select & Bulk Operations**
   - Space to select multiple cards
   - Bulk move/tag updates
   - 6-8 hours estimated

4. **Enhanced Filtering**
   - Regex search
   - Date range filters
   - Save filter presets
   - 6-8 hours estimated

---

## Key Files for Understanding

**For Architecture Understanding:**
- `model.go` - Core state management and layout
- `backend.go` + `backend_github.go` - Backend abstraction
- `types.go` - Data model definition

**For UI Implementation:**
- `view.go` - Rendering system
- `styles.go` - Color and styling
- `update_keyboard.go` + `update_mouse.go` - Event handling

**For Documentation:**
- `PLAN.md` - Complete implementation roadmap with patterns
- `CLAUDE.md` - AI integration guide
- `HOTKEYS.md` - All keyboard shortcuts

---

## Conclusion

tkan is a well-structured, partially complete TUI application with solid fundamentals. The first two phases (foundation and drag & drop) are fully implemented with good code quality. The project is ready for continued development with clear documentation for the next phases.

The codebase demonstrates good architectural practices, proper separation of concerns, and reuse of proven patterns from related projects. The main gaps are feature implementation (table view, card editing, search) rather than architectural issues.

Next developer should focus on implementing Phase 3 (Table View) and Phase 4 (Card Editing) using the provided patterns in PLAN.md. The project would benefit from adding test coverage and implementing the TODO/FIXME items before moving to production.

**Estimated Time to v1.0**: 2-3 sprints with 1-2 developers working on feature implementation.

