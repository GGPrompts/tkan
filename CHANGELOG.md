# tkan Changelog

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
