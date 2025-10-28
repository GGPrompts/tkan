# tkan - Terminal Kanban Board

**A beautiful dual-view task management TUI with visual Kanban board and sortable table views**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Bubbletea](https://img.shields.io/badge/Bubbletea-TUI-ff69b4)](https://github.com/charmbracelet/bubbletea)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## ✨ Features

**Currently Implemented:**
- ✅ **📋 Visual Kanban Board**: BACKLOG, TODO, PROGRESS, REVIEW, DONE columns
- ✅ **🃏 Solitaire-Style Cards**: Stacked cards with wrapped titles (12×5 chars)
- ✅ **🖱️ Drag & Drop**: Mouse-based card dragging with live visual feedback
- ✅ **🎯 Card Reordering**: Drag cards anywhere - between cards or across columns
- ✅ **💚 Drop Indicator**: Green line shows exactly where cards will land
- ✅ **👻 Ghost Cards**: Dragged cards appear faded at source position
- ✅ **📁 Multi-Project Support**: Automatic discovery of `.tkan.yaml` files
- ✅ **🎨 Detail Panel**: Toggleable 33% width panel with full card info
- ✅ **📦 Archive Column**: Toggleable archive (press 'a')
- ✅ **⌨️ Keyboard Navigation**: Full keyboard control (←→↑↓, vim keys)
- ✅ **💾 YAML Persistence**: Plain text `.tkan.yaml` files
- ✅ **🎯 Project Selector**: Choose from multiple projects with ↑/↓

**Coming Soon:**
- 📅 **Card Editing**: Create/edit/delete cards (Phase 2, partial)
- 📅 **Table View**: Sortable data table view (Phase 3)
- 📅 **Search/Filter**: Find cards by text, tags, assignee (Phase 3)

---

## 📸 Screenshots

### Board View (Current)
```
┌────────────────────────────────────────────────┬─────────────────────┐
│ BACKLOG │ TODO │ PROGRESS │ REVIEW │ DONE     │ ▶ CARD DETAILS      │
│   (1)   │ (2)  │   (5)    │  (1)   │  (1)     │                     │
│ ┌──────  ┌────  ┌──────────  ┌────  ┌────     │ Fix login flow      │
│ │New      │Fix    │Card 1      │Rev   │Set      │ ━━━━━━━━━━━━━━━━━━━ │
│ ┌──────  ┌────  ┌──────────  ┌────             │                     │
│ │feature  │Write  │Card 2      │PR             │ Description:        │
│           ┌────  ┌──────────┐                   │ Users can't auth    │
│           │Add    │Card 3    │                  │ via OAuth. Error    │
│                   │          │                  │ 401 on refresh...   │
│                   └──────────┘                  │                     │
│                   ┌──────────┐                  │ Tags: #bug #p1      │
│                   │Card 4    │                  │ Assigned: @alice    │
│                   │          │                  │ Due: Jan 15         │
│                   └──────────┘                  │ Created: Oct 18     │
│                   ┌──────────┐                  │ Modified: Oct 27    │
│                   │Last card │                  │                     │
│                   │visible   │                  │ [E]dit [M]ove      │
│                   └──────────┘                  │                     │
└────────────────────────────────────────────────┴─────────────────────┘
```

**Features shown:**
- Solitaire-style stacking (partial cards + full last card)
- Wrapped titles (no truncation!)
- Detail panel with full metadata
- Cards stack vertically to show more in less space

### Table View (Coming in Phase 3)

---

## 🚀 Quick Start

### Installation

**From source:**
```bash
git clone https://github.com/GGPrompts/tkan.git
cd tkan
go build
sudo mv tkan /usr/local/bin/
```

**Or download binary from [Releases](https://github.com/GGPrompts/tkan/releases)** (coming soon)

### Usage

```bash
# Start tkan (scans current directory for .tkan.yaml files)
tkan

# If no board found, creates .tkan.yaml with sample cards
# If multiple projects found, shows project selector
```

### First Time Setup

1. Navigate to your project directory: `cd ~/projects/myapp`
2. Run: `tkan`
3. tkan will create `.tkan.yaml` in the project root
4. Start adding cards with `n` key
5. Drag cards between columns or press `v` for table view

---

## ⌨️ Keyboard Shortcuts

### Project List
- `↑/↓` or `k/j` - Navigate projects
- `Enter` - Open selected project
- `q` - Quit

### Board View
- `←/→` or `h/l` - Navigate columns
- `↑/↓` or `k/j` - Navigate cards
- `g` - Jump to first column
- `G` - Jump to last column
- `Tab` - Toggle detail panel
- `a` - Toggle archive column visibility
- `p` - Return to project list (if multiple projects)
- `v` - Table view (coming in Phase 3)
- `q` - Quit

### Mouse Controls
- **Click & Drag** - Move cards between columns or reorder within column
- **Visual Feedback** - Green line shows drop position, ghost card at source
- **Precise Positioning** - Hover over top/bottom half of cards to insert before/after

### Coming Soon (Phase 2, partial)
- `n` - New card
- `e` - Edit card
- `d` - Delete card
- `m` - Move card to column (keyboard alternative)
- `/` - Search/filter

---

## 📝 Board Configuration

tkan stores boards in `.tkan.yaml` files in your project directory:

```yaml
name: MyProject Kanban
columns:
  - name: TODO
    color: "#FFA500"
    cards:
      - id: card-001
        title: Fix login flow
        description: Users unable to authenticate
        tags: [bug, p1]
        assignees: [alice]
        due_date: 2024-01-15T00:00:00Z
        created: 2024-01-01T10:00:00Z
        modified: 2024-01-10T15:30:00Z
```

You can edit this file directly or use tkan's UI.

---

## 🏗️ Architecture

Built with proven patterns from:
- **[TUITemplate](https://github.com/GGPrompts/TUITemplate)** - Dual-pane layout system
- **[TFE](https://github.com/GGPrompts/tfe)** - Sortable table headers
- **[Solitaire](https://github.com/GGPrompts/TUIClassics)** - Card drag mechanics

**Technology Stack:**
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Huh](https://github.com/charmbracelet/huh) - Forms (card editing)
- YAML - Data persistence

---

## 🗺️ Roadmap

**Phase 1 - Foundation** ✅ COMPLETED (2025-10-28)
- [x] Project structure and data model
- [x] Board view with stacked cards
- [x] Detail panel (toggleable)
- [x] Keyboard navigation
- [x] Multi-project support
- [x] YAML persistence
- [x] Archive toggle

**Phase 2 - Interactions** ✅ COMPLETED (2025-10-28)
- [x] Drag & drop cards (Solitaire-style)
- [x] Card reordering within columns
- [x] Visual drop indicator (green line)
- [x] Ghost card effect during drag
- [x] Move cards across all columns (including DONE)
- [ ] Card creation (pending)
- [ ] Card editing (pending)
- [ ] Card deletion (pending)
- [ ] Move card keyboard shortcut (pending)

**Phase 3 - Table View** 📅 PLANNED
- [ ] Table view with sortable headers
- [ ] Click to sort
- [ ] Search and filtering
- [ ] Column customization

**v1.1 - Enhanced Features**
- [ ] Undo/redo
- [ ] Custom columns
- [ ] Multi-select cards
- [ ] Card history
- [ ] Export to CSV/JSON

**v2.0 - Advanced Features**
- [ ] Multi-project dashboard
- [ ] Analytics/reporting
- [ ] GitHub Issues integration
- [ ] Jira import

See [PLAN.md](PLAN.md) for detailed implementation plan.

---

## 🤝 Contributing

Contributions welcome! Areas of focus:
- Card editing improvements
- Additional integrations (GitHub, Jira)
- Performance optimizations
- Documentation improvements

---

## 📄 License

MIT License - Use freely for any purpose.

---

## 🙏 Acknowledgments

Built with ❤️ using [Charm](https://charm.sh/) libraries.

Inspired by:
- [lazygit](https://github.com/jesseduffield/lazygit) - Accordion layout
- [taskwarrior-tui](https://github.com/kdheepak/taskwarrior-tui) - Task management
- Physical Kanban boards - Original inspiration

---

**Start managing tasks beautifully in your terminal!** 🚀
