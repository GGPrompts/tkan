# tkan Calendar View Implementation Guide

## Overview

Add a calendar view to tkan that displays tasks by their `DueDate` field in a monthly grid layout.

## Architecture Changes

### 1. Update ViewMode (types.go)

```go
// ViewMode represents the current view
type ViewMode int

const (
    ViewProjectList ViewMode = iota
    ViewBoard
    ViewTable
    ViewCalendar  // NEW
)
```

### 2. Add Calendar State (types.go)

```go
// Model additions
type Model struct {
    // ... existing fields ...

    // Calendar view state
    calendarMonth time.Time    // Currently displayed month
    calendarDay   int           // Selected day (1-31)
    calendarTasks map[string][]*Card  // Tasks grouped by date "2024-11-15"
}
```

### 3. Add Date Parsing Utilities (calendar.go - NEW FILE)

Create `calendar.go`:

```go
package main

import (
    "fmt"
    "time"
)

// ParseDueDate parses the DueDate string from a card
func ParseDueDate(dateStr string) (time.Time, error) {
    if dateStr == "" {
        return time.Time{}, fmt.Errorf("empty date")
    }

    // Try multiple formats
    formats := []string{
        "2006-01-02",           // ISO format from GitHub
        "2006-01-02T15:04:05Z", // Full ISO with time
        "01/02/2006",           // US format
        "Jan 2, 2006",          // Human readable
    }

    for _, format := range formats {
        if t, err := time.Parse(format, dateStr); err == nil {
            return t, nil
        }
    }

    return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
}

// GetTasksByDate groups tasks by their due date
func (m *Model) GetTasksByDate() map[string][]*Card {
    taskMap := make(map[string][]*Card)

    for _, card := range m.board.Cards {
        if card.DueDate == "" {
            continue
        }

        dueDate, err := ParseDueDate(card.DueDate)
        if err != nil {
            continue
        }

        dateKey := dueDate.Format("2006-01-02")
        taskMap[dateKey] = append(taskMap[dateKey], card)
    }

    return taskMap
}

// GetMonthGrid returns a 2D array of calendar days for a given month
func GetMonthGrid(year int, month time.Month) [][]int {
    firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
    lastDay := firstDay.AddDate(0, 1, -1)

    // Get weekday of first day (0=Sunday, 6=Saturday)
    firstWeekday := int(firstDay.Weekday())
    daysInMonth := lastDay.Day()

    // Calculate number of weeks needed
    totalCells := firstWeekday + daysInMonth
    numWeeks := (totalCells + 6) / 7

    // Build grid
    grid := make([][]int, numWeeks)
    day := 1

    for week := 0; week < numWeeks; week++ {
        grid[week] = make([]int, 7)

        for dayOfWeek := 0; dayOfWeek < 7; dayOfWeek++ {
            // First week: skip days before month starts
            if week == 0 && dayOfWeek < firstWeekday {
                grid[week][dayOfWeek] = 0
                continue
            }

            // Past end of month
            if day > daysInMonth {
                grid[week][dayOfWeek] = 0
                continue
            }

            grid[week][dayOfWeek] = day
            day++
        }
    }

    return grid
}
```

### 4. Add Calendar Rendering (view_calendar.go - NEW FILE)

Create `view_calendar.go`:

```go
package main

import (
    "fmt"
    "strings"
    "time"

    "github.com/charmbracelet/lipgloss"
)

// renderCalendarView renders the calendar view
func (m Model) renderCalendarView() string {
    var sections []string

    // Title bar with month/year
    title := m.renderCalendarTitle()
    sections = append(sections, title)

    // Calendar grid
    calendar := m.renderCalendarGrid()
    sections = append(sections, calendar)

    // Status bar
    status := m.renderCalendarStatus()
    sections = append(sections, status)

    return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderCalendarTitle renders the title with month/year
func (m Model) renderCalendarTitle() string {
    monthYear := m.calendarMonth.Format("January 2006")
    title := fmt.Sprintf("📋 tkan - %s", m.board.Name)
    viewLabel := fmt.Sprintf("Calendar View - %s", monthYear)

    titleStyle := styleTitle.Width(m.width)
    padding := m.width - len(title) - len(viewLabel) - 2

    return titleStyle.Render(title + strings.Repeat(" ", padding) + viewLabel)
}

// renderCalendarGrid renders the month grid with tasks
func (m Model) renderCalendarGrid() string {
    // Update task map
    taskMap := m.GetTasksByDate()

    year, month, _ := m.calendarMonth.Date()
    grid := GetMonthGrid(year, month)

    // Day headers
    dayHeaders := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
    headerRow := ""
    cellWidth := (m.width - 8) / 7  // 7 days, minus borders

    for _, day := range dayHeaders {
        headerStyle := lipgloss.NewStyle().
            Width(cellWidth).
            Align(lipgloss.Center).
            Bold(true).
            Foreground(lipgloss.Color("12"))
        headerRow += headerStyle.Render(day)
    }

    var rows []string
    rows = append(rows, headerRow)
    rows = append(rows, strings.Repeat("─", m.width-2))

    // Render each week
    today := time.Now()
    for _, week := range grid {
        row := m.renderCalendarWeek(week, year, month, taskMap, today, cellWidth)
        rows = append(rows, row)
        rows = append(rows, strings.Repeat("─", m.width-2))
    }

    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        Width(m.width - 4).
        Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

// renderCalendarWeek renders a single week row
func (m Model) renderCalendarWeek(week []int, year int, month time.Month,
    taskMap map[string][]*Card, today time.Time, cellWidth int) string {

    var cells []string
    cellHeight := 3  // Lines per cell

    for dayOfWeek, day := range week {
        cell := m.renderCalendarCell(day, year, month, taskMap, today,
            cellWidth, cellHeight, dayOfWeek == 6 || dayOfWeek == 0)
        cells = append(cells, cell)
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, cells...)
}

// renderCalendarCell renders a single day cell
func (m Model) renderCalendarCell(day int, year int, month time.Month,
    taskMap map[string][]*Card, today time.Time,
    width int, height int, isWeekend bool) string {

    if day == 0 {
        // Empty cell (before/after month)
        return lipgloss.NewStyle().
            Width(width).
            Height(height).
            Render("")
    }

    // Check if this is today
    date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
    isToday := date.Year() == today.Year() &&
               date.Month() == today.Month() &&
               date.Day() == today.Day()

    // Check if selected
    isSelected := day == m.calendarDay

    // Get tasks for this date
    dateKey := date.Format("2006-01-02")
    tasks := taskMap[dateKey]

    // Build cell content
    var lines []string

    // Line 1: Day number
    dayStr := fmt.Sprintf("%2d", day)
    if isToday {
        dayStr = lipgloss.NewStyle().
            Foreground(lipgloss.Color("10")).
            Bold(true).
            Render(dayStr)
    }
    lines = append(lines, dayStr)

    // Line 2: Task count or first task
    if len(tasks) == 0 {
        lines = append(lines, "")
    } else if len(tasks) == 1 {
        taskTitle := truncate(tasks[0].Title, width-2)
        lines = append(lines, lipgloss.NewStyle().
            Foreground(lipgloss.Color("14")).
            Render(taskTitle))
    } else {
        lines = append(lines, lipgloss.NewStyle().
            Foreground(lipgloss.Color("11")).
            Render(fmt.Sprintf("[%d tasks]", len(tasks))))
    }

    // Line 3: Status indicator
    statusLine := ""
    if len(tasks) > 0 {
        // Show status dots
        statusCounts := make(map[string]int)
        for _, task := range tasks {
            statusCounts[task.Column]++
        }

        if statusCounts["TODO"] > 0 {
            statusLine += "●"  // Red dot
        }
        if statusCounts["IN_PROGRESS"] > 0 {
            statusLine += "●"  // Yellow dot
        }
        if statusCounts["DONE"] > 0 {
            statusLine += "●"  // Green dot
        }
    }
    lines = append(lines, statusLine)

    content := lipgloss.JoinVertical(lipgloss.Left, lines...)

    // Style the cell
    cellStyle := lipgloss.NewStyle().
        Width(width).
        Height(height).
        Padding(0, 1)

    if isSelected {
        cellStyle = cellStyle.
            Border(lipgloss.RoundedBorder()).
            BorderForeground(lipgloss.Color("12"))
    }

    if isWeekend {
        cellStyle = cellStyle.
            Background(lipgloss.Color("235"))
    }

    if isToday {
        cellStyle = cellStyle.
            Background(lipgloss.Color("237")).
            Bold(true)
    }

    return cellStyle.Render(content)
}

// renderCalendarStatus renders the status bar with keyboard shortcuts
func (m Model) renderCalendarStatus() string {
    help := "← → Next/Prev Month | ↑↓←→ Navigate Days | Enter: View Day Tasks | d: Detail | b: Board View"

    return styleStatus.
        Width(m.width).
        Render(help)
}

// Helper function to truncate strings
func truncate(s string, max int) string {
    if len(s) <= max {
        return s
    }
    return s[:max-3] + "..."
}
```

### 5. Add Keyboard Navigation (update_keyboard.go)

Add to existing `handleKeyPress`:

```go
func (m model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "c":
        // Toggle to calendar view
        m.viewMode = ViewCalendar
        m.calendarMonth = time.Now()
        m.calendarDay = time.Now().Day()
        return m, nil

    case "b":
        // Back to board view
        if m.viewMode == ViewCalendar {
            m.viewMode = ViewBoard
            return m, nil
        }
    }

    // Calendar-specific navigation
    if m.viewMode == ViewCalendar {
        return m.handleCalendarKeyPress(msg)
    }

    // ... existing keyboard handling ...
}

func (m model) handleCalendarKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    switch msg.String() {
    case "left", "h":
        // Previous month
        m.calendarMonth = m.calendarMonth.AddDate(0, -1, 0)
        m.calendarDay = 1
        return m, nil

    case "right", "l":
        // Next month
        m.calendarMonth = m.calendarMonth.AddDate(0, 1, 0)
        m.calendarDay = 1
        return m, nil

    case "up", "k":
        // Previous week (subtract 7 days)
        m.calendarDay = max(1, m.calendarDay-7)
        return m, nil

    case "down", "j":
        // Next week (add 7 days)
        year, month, _ := m.calendarMonth.Date()
        lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
        m.calendarDay = min(lastDay, m.calendarDay+7)
        return m, nil

    case "enter":
        // View tasks for selected day
        m.showCalendarDayDetail = true
        return m, nil

    case "t":
        // Jump to today
        now := time.Now()
        m.calendarMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
        m.calendarDay = now.Day()
        return m, nil
    }

    return m, nil
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

### 6. Update Main View (view.go)

Update the `View()` function:

```go
func (m Model) View() string {
    // ... existing code ...

    // Render based on view mode
    switch m.viewMode {
    case ViewProjectList:
        return m.renderProjectListView()
    case ViewBoard:
        return m.renderBoardView()
    case ViewTable:
        return "Table view (not implemented)"
    case ViewCalendar:  // NEW
        return m.renderCalendarView()
    default:
        return "Unknown view mode"
    }
}
```

## GitHub Projects Integration

### Fetching Date Fields from GitHub

Update `backend_github.go` to fetch date field values:

```go
// In LoadBoard(), add date field extraction
fieldValues(first: 20) {
    nodes {
        // ... existing field types ...

        ... on ProjectV2ItemFieldDateValue {
            date
            field {
                ... on ProjectV2Field {
                    name
                }
            }
        }
    }
}
```

Then map the date field to `Card.DueDate`:

```go
// In processFieldValues()
if dateValue, ok := fieldNode["date"].(string); ok {
    if fieldName == "Target Date" || fieldName == "Due Date" {
        card.DueDate = dateValue
    }
}
```

## Usage

### User Workflow

1. **Open calendar view:**
   ```
   Press 'c' key
   ```

2. **Navigate months:**
   ```
   ← → or h/l for prev/next month
   ```

3. **Navigate days:**
   ```
   ↑↓←→ or hjkl to move between days
   't' to jump to today
   ```

4. **View day details:**
   ```
   Press Enter on a day to see all tasks due that day
   ```

5. **Return to board:**
   ```
   Press 'b'
   ```

## Testing Plan

### Unit Tests

```go
// calendar_test.go
func TestGetMonthGrid(t *testing.T) {
    // Test November 2024
    grid := GetMonthGrid(2024, time.November)

    // Should have 5 weeks
    assert.Equal(t, 5, len(grid))

    // First day (Nov 1) is a Friday (index 5)
    assert.Equal(t, 1, grid[0][5])

    // Last day is 30
    assert.Equal(t, 30, grid[4][6])
}

func TestParseDueDate(t *testing.T) {
    tests := []struct{
        input string
        valid bool
    }{
        {"2024-11-15", true},
        {"2024-11-15T10:30:00Z", true},
        {"11/15/2024", true},
        {"invalid", false},
    }

    for _, tt := range tests {
        _, err := ParseDueDate(tt.input)
        if tt.valid {
            assert.NoError(t, err)
        } else {
            assert.Error(t, err)
        }
    }
}
```

### Manual Testing

1. Create test cards with various due dates
2. Navigate through months
3. Verify task counts are accurate
4. Test edge cases (empty months, many tasks in one day)
5. Verify "today" highlighting works

## Future Enhancements

1. **Week view** - Show just one week with more detail
2. **Drag & drop** - Drag tasks to new dates
3. **Color coding** - Different colors for priority/status
4. **Multi-select** - Select multiple days
5. **Iteration view** - Show sprints/iterations instead of calendar months
6. **Mini calendar** - Show 3-month view (prev, current, next)

## Files to Create/Modify

**New Files:**
- `calendar.go` - Calendar utilities and date parsing
- `view_calendar.go` - Calendar rendering
- `calendar_test.go` - Unit tests

**Modified Files:**
- `types.go` - Add ViewCalendar and calendar state
- `view.go` - Add ViewCalendar case
- `update_keyboard.go` - Add calendar navigation
- `backend_github.go` - Fetch date fields

## Estimated Effort

- **Core implementation:** 4-6 hours
- **Testing & polish:** 2-3 hours
- **GitHub integration:** 1-2 hours
- **Total:** ~8-11 hours

This gives you a fully functional calendar view that integrates with GitHub Projects date fields!
