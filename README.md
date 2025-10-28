# tkan - Terminal Kanban Board

**A beautiful dual-view task management TUI with visual Kanban board and sortable table views**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Bubbletea](https://img.shields.io/badge/Bubbletea-TUI-ff69b4)](https://github.com/charmbracelet/bubbletea)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## ✨ Features

- **🎯 Dual View System**: Toggle between visual Kanban board and sortable data table
- **🎮 Drag & Drop**: Solitaire-style card dragging between columns
- **📊 Sortable Headers**: Click table headers to sort by any field
- **📋 Detail Panel**: Always-visible card details with full metadata
- **🎨 Beautiful UI**: Carefully crafted with Lipgloss styling
- **⌨️ Keyboard-First**: Complete keyboard navigation with mouse support
- **📁 Project-Specific**: Each project has its own `.tkan.yaml` board
- **🔍 Filtering**: Search by text, filter by tags, assignees, or columns

---

## 📸 Screenshots

### Board View
```
┌──────────────────────────────────────────────┬─────────────────────┐
│  TODO   │ PROGRESS │ REVIEW  │   DONE        │ ▶ CARD DETAILS      │
│   (3)   │   (2)    │  (1)    │   (5)         │                     │
│ ┌──────┐ ┌──────┐  ┌──────┐  ┌──────┐       │ Fix login flow      │
│ │Fix   │ │Add   │  │Review│  │Setup │       │ ━━━━━━━━━━━━━━━━━━━ │
│ │login │ │auth  │  │PR#42 │  │DB    │       │                     │
│ │#p1   │ │#feat │  │#code │  │#done │       │ Description:        │
│ └──────┘ └──────┘  └──────┘  └──────┘       │ Users can't auth... │
└──────────────────────────────────────────────┴─────────────────────┘
```

### Table View
```
┌─────────────────────────────────────────────────────────────────────┐
│  Title ↓              Column      Tags        Assignee   Due Date   │
├─────────────────────────────────────────────────────────────────────┤
│  Fix login flow       TODO        #bug #p1    @alice     Jan 15     │
│  Add OAuth support    PROGRESS    #feature    @bob       Jan 20     │
│  Review PR #42        REVIEW      #code       @charlie   Jan 18     │
│  Setup database       DONE        #infra      @dave      -          │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Installation

**From source:**
```bash
git clone https://github.com/yourname/tkan.git
cd tkan
go build
sudo mv tkan /usr/local/bin/
```

**Or download binary from [Releases](https://github.com/yourname/tkan/releases)**

### Usage

```bash
# Start tkan in current directory
tkan

# Use specific project
tkan ~/projects/myapp

# List all projects with boards
tkan --list

# Create new board interactively
tkan init
```

### First Time Setup

1. Navigate to your project directory: `cd ~/projects/myapp`
2. Run: `tkan`
3. tkan will create `.tkan.yaml` in the project root
4. Start adding cards with `n` key
5. Drag cards between columns or press `v` for table view

---

## ⌨️ Keyboard Shortcuts

### General
- `v` - Toggle between board and table view
- `Tab` - Toggle detail panel
- `/` - Search/filter mode
- `#` - Filter by tag
- `@` - Filter by assignee
- `q` - Quit
- `?` - Show help

### Board View
- `←/→` or `h/l` - Navigate columns
- `↑/↓` or `k/j` - Navigate cards
- `Enter` - Select card
- `n` - New card
- `e` - Edit card
- `d` - Delete card
- `m` - Move card to column

### Table View
- `↑/↓` or `k/j` - Navigate rows
- `Enter` - Select card
- Click headers to sort

### Mouse
- **Click** - Select card/row
- **Drag** (Board view) - Move card between columns
- **Click header** (Table view) - Sort by column
- **Right-click** - Context menu (future)

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

**v1.0 - Core Features** (In Progress)
- [x] Project structure and data model
- [ ] Board view with draggable cards
- [ ] Table view with sortable headers
- [ ] Detail panel
- [ ] Keyboard navigation
- [ ] Card editing
- [ ] Search and filtering

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
