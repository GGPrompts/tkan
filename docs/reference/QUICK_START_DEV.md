# tkan Development Quick Start

## Project At A Glance

- **Type**: Terminal Kanban Board (TUI)
- **Language**: Go 1.21+
- **Framework**: Bubbletea + Lipgloss
- **Status**: Phase 2 complete (drag & drop working), Phase 3 planned
- **Lines of Code**: 2,501 across 13 Go files
- **Architecture**: Model-View-Update (Elm pattern) + Backend abstraction

## Getting Started

### Build & Run

```bash
cd /home/matt/projects/tkan
go build
./tkan                    # Local YAML mode
./tkan --github owner/7   # GitHub Projects mode
```

### Project Structure

```
tkan/
├── main.go              # Entry point (140 lines)
├── model.go             # State management (463 lines) ⭐ Read this first
├── view.go              # Rendering (396 lines)
├── update.go            # Event dispatcher (49 lines)
├── update_keyboard.go   # Keyboard handling (151 lines)
├── update_mouse.go      # Drag & drop (146 lines) ⭐ Well-implemented
├── styles.go            # Lipgloss styling (268 lines)
├── types.go             # Data structures (107 lines)
├── backend.go           # Backend interface (111 lines)
├── backend_github.go    # GitHub Projects (414 lines)
├── persistence.go       # YAML load/save (152 lines)
├── projects.go          # Project scanner (85 lines)
├── update_timer.go      # Drag delay (19 lines)
├── .tkan.yaml           # Example board
└── PLAN.md              # Implementation roadmap (33 KB) ⭐ Essential reading
```

## Key Files to Understand

### 1. **model.go** (Start here!)
- `NewModel()` - Initialize application state
- `calculateLayout()` - Board vs detail panel sizing
- `moveCard()` - Move card between columns (includes save)
- `getDropPosition()` - Calculate where card lands during drag
- `getInsertIndexInColumn()` - Smart insertion (top/bottom half)

### 2. **update_mouse.go** (Drag implementation)
- `handleMousePress()` - Start potential drag
- `handleMouseMotion()` - Update drop target visually
- `handleMouseRelease()` - Drop card or cancel
- **Pattern**: Distance-based detection (`dx² + dy² < 4`)

### 3. **view.go** (Rendering)
- `renderBoardView()` - Main board layout
- `renderColumn()` - Solitaire-style stacking
- `renderDetailPanel()` - Right side panel
- **Key insight**: Stacked cards show only 2 lines except last card (5 lines)

### 4. **styles.go** (Colors & styling)
- Card size: 12×5 characters
- Selected: Magenta border
- Ghost: Subdued gray (during drag)
- Drop indicator: Green line

### 5. **backend.go & backend_github.go** (Persistence)
- Interface-based design
- LocalBackend - YAML files
- GitHubBackend - GitHub Projects API
- All backends implement: LoadBoard, SaveBoard, MoveCard, etc.

## What's Complete ✅

### Phase 1: Foundation
- Board view with columns
- Detail panel (Tab to toggle)
- Keyboard navigation (arrows, vim keys)
- YAML persistence
- Multi-project support

### Phase 2: Drag & Drop
- Click & drag cards
- Visual feedback (ghost + drop indicator)
- Works across columns
- Auto-save on drop
- 150ms drag delay to prevent accidental drags

### Phase 2.5: GitHub Backend
- `--github owner/project-number` flag
- Loads from GitHub Projects
- Updates status when cards move
- Full integration in main.go

## What's NOT Complete ❌

### Phase 3: Table View (Critical)
- `v` key switches to table view
- Shows placeholder message only
- No table rendering code
- **TODO**: Implement using TFE pattern from PLAN.md

### Phase 4: Card Editing (Critical)
- Keyboard shortcuts exist but empty:
  - `n` - New card
  - `e` - Edit card
  - `d` - Delete card
  - `m` - Move card
- **TODO**: Implement dialog forms (use Huh library)

## Quick Wins (1-2 hours each)

1. **Card Creation** - Implement `n` key
   - Use Huh library for form
   - Generate new ID
   - Add to board.Cards
   - Save via backend

2. **Card Deletion** - Implement `d` key
   - Move to ARCHIVE column
   - Save changes

3. **Card Move** - Implement `m` key
   - Show column picker
   - Call moveCard()

4. **Help Screen** - Implement `?` key
   - Show HOTKEYS.md content

## Medium Tasks (4-8 hours)

1. **Table View** (Phase 3)
   - Copy TFE pattern from PLAN.md
   - Calculate dynamic column widths
   - Render sortable headers
   - Implement sort logic

2. **GitHub Field Discovery**
   - Dynamic fetch field IDs instead of hardcoded
   - Cache results

3. **Search & Filter**
   - Implement `/` key
   - Filter by tags, assignees, dates

## Testing Strategy

```bash
# Test various terminal sizes
export COLUMNS=80 LINES=24    # Min size
export COLUMNS=120 LINES=40   # Medium
export COLUMNS=200 LINES=60   # Wide

# Test with many cards
# Edit .tkan.yaml to add 100+ cards

# Test drag edge cases
# Drag to same position (should be no-op)
# Drag to empty column
# Drag to ARCHIVE
```

## Known Issues

1. **TODO in code** (update_keyboard.go:59)
   - "TODO: Show error message"
   - Easy fix if you need better error handling

2. **Hardcoded GitHub field IDs**
   - Should be discovered dynamically
   - Field IDs in backend_github.go

3. **No input validation**
   - Card data should be validated before save

4. **No tests**
   - Would help prevent regressions

## Learning Resources

### In This Project
- **PLAN.md** - Complete implementation guide with code patterns
- **HOTKEYS.md** - All keyboard shortcuts documented
- **CLAUDE.md** - GitHub integration guide for AI
- **README.md** - Feature overview
- **COMPLETE_SETUP_GUIDE.md** - Detailed setup instructions

### External
- **TUITemplate** - Dual-pane layout patterns
- **Solitaire** - Card drag mechanics
- **TFE** - Table view patterns (sortable headers)
- **Bubbletea docs** - Framework reference

## Development Workflow

### Adding a Feature

1. **Update types.go** if new data needed
2. **Update model.go** if state management needed
3. **Update view.go** if UI changes
4. **Update update_keyboard.go** or **update_mouse.go** for input
5. **Test** with `go run .`

### Testing Workflow

1. Run `./tkan` (local mode)
2. Test with `./tkan --github GGPrompts/7` (GitHub mode)
3. Check YAML file after changes
4. Verify in project list view if multiple projects

## Useful Shortcuts

### Vim Navigation (Already Implemented)
- `h/l` - Left/right (← →)
- `j/k` - Down/up (↑ ↓)
- `g/G` - First/last column
- `0/$/^` - Could be future features

### Current Features
- `Tab` - Toggle detail panel
- `a` - Toggle archive column
- `p` - Project picker
- `v` - Switch view (works, shows placeholder)
- `q` - Quit

### Skeleton Code (Waiting to be Implemented)
- `n` - New card
- `e` - Edit card
- `d` - Delete card
- `m` - Move card
- `/` - Search
- `?` - Help

## Code Style

The project follows:
- Go conventions
- One function per concern
- Clear variable names
- Comments for complex logic
- Consistent error handling

## Performance Notes

- Renders efficiently even with many cards
- Stacking algorithm is O(n)
- Drag detection is fast (distance check)
- YAML parsing is straightforward
- GitHub API calls should be cached (not currently done)

## Common Tasks

### Add New Column Type
1. Update types.go (Column struct)
2. Update persistence.go (CreateDefaultBoard)
3. Update backend_github.go (if GitHub support needed)

### Change Card Size
1. Update styles.go (cardWidth, cardHeight)
2. Update view.go (renderCard functions)
3. Update model.go (insertion calculations)

### Add New Color
1. Update styles.go (colorVariable)
2. Update relevant style definitions
3. Test rendering

## When Stuck

1. **Check PLAN.md** - Has implementation patterns for all phases
2. **Read comments** - Key functions are documented
3. **Compare with TFE** - For table view patterns
4. **Compare with Solitaire** - For drag mechanics
5. **Check existing patterns** - In model.go or view.go

## Next Steps

1. Read PLAN.md (understand the full roadmap)
2. Read model.go (understand state management)
3. Pick ONE quick win from above
4. Implement it
5. Test it
6. Commit it
7. Repeat for next feature

Estimated time to v1.0: **16-23 hours** (2-3 development sprints)

Good luck! This is a well-structured codebase with clear patterns. The main work ahead is feature implementation, not architectural changes.

