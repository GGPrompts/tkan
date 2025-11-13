# tkan Calendar View Skill

Add a calendar view to tkan for visualizing tasks by due date.

## Quick Start

### 1. Add the View Mode

```go
// types.go
const (
    ViewProjectList ViewMode = iota
    ViewBoard
    ViewTable
    ViewCalendar  // Add this
)
```

### 2. Create Calendar Files

```bash
# Create new files
touch calendar.go view_calendar.go calendar_test.go

# Copy implementations from IMPLEMENTATION.md
```

### 3. Test It

```bash
# Build
go build

# Run
./tkan

# Press 'c' to switch to calendar view
```

## Features

✅ Monthly calendar grid with tasks
✅ Navigate months with arrow keys
✅ Jump to today with 't'
✅ View day details
✅ Integration with GitHub Projects date fields

## Files

- `IMPLEMENTATION.md` - Complete implementation guide
- `README.md` - This file

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `c` | Switch to calendar view |
| `b` | Back to board view |
| `←` `→` | Previous/Next month |
| `↑` `↓` | Navigate days |
| `t` | Jump to today |
| `Enter` | View day details |

## See Also

- Main tkan documentation
- `/bubbletea` skill for TUI patterns
- `/tui-add-keybinding` for adding shortcuts
