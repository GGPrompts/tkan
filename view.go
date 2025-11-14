package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the entire application UI (required by Bubbletea)
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Check minimum size
	if m.width < 80 || m.height < 20 {
		return fmt.Sprintf("Terminal too small (%d×%d). Need at least 80×20.", m.width, m.height)
	}

	// Render based on view mode
	switch m.viewMode {
	case ViewProjectList:
		return m.renderProjectListView()
	case ViewBoard:
		return m.renderBoardView()
	case ViewTable:
		return "Table view (not implemented in Phase 1)"
	case ViewHelp:
		return m.renderHelpView()
	default:
		return "Unknown view mode"
	}
}

// renderBoardView renders the Kanban board view
func (m Model) renderBoardView() string {
	var sections []string

	// Title bar
	title := m.renderTitle()
	sections = append(sections, title)

	// Main content area (board + optional detail panel)
	mainContent := m.renderMainContent()
	sections = append(sections, mainContent)

	// Status bar
	status := m.renderStatus()
	sections = append(sections, status)

	boardView := lipgloss.JoinVertical(lipgloss.Left, sections...)

	// Render form overlay if form is open
	if m.formMode != FormNone {
		return m.renderFormOverlay(boardView)
	}

	return boardView
}

// renderTitle renders the title bar
func (m Model) renderTitle() string {
	boardName := m.board.Name
	viewLabel := "Board View"

	title := fmt.Sprintf("📋 tkan - %s", boardName)
	titleStyle := styleTitle.Width(m.boardWidth)

	if m.showDetails {
		// Title spans both board and detail areas
		titleStyle = styleTitle.Width(m.width)
	}

	return titleStyle.Render(title + strings.Repeat(" ", m.width-len(title)-10) + viewLabel)
}

// renderMainContent renders the board and optional detail panel side by side
func (m Model) renderMainContent() string {
	boardContent := m.renderBoard()

	if m.showDetails {
		detailContent := m.renderDetailPanel()

		// Join board and detail panel horizontally (detail panel has its own border)
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			boardContent,
			detailContent,
		)
	}

	return boardContent
}

// renderBoard renders the Kanban board columns and cards
func (m Model) renderBoard() string {
	contentHeight := m.getContentHeight()

	// Column headers
	headers := m.renderColumnHeaders()

	// Column contents (cards stacked vertically)
	columns := m.renderColumns(contentHeight)

	// Join headers and columns
	board := lipgloss.JoinVertical(lipgloss.Left, headers, columns)

	return lipgloss.NewStyle().
		Width(m.boardWidth).
		Height(contentHeight + 1). // +1 for headers
		Render(board)
}

// renderColumnHeaders renders the column headers with counts
func (m Model) renderColumnHeaders() string {
	var headers []string
	visibleColumns := m.getVisibleColumns()

	for i, col := range visibleColumns {
		count := len(col.Cards)
		label := fmt.Sprintf("%s (%d)", col.Name, count)

		// Use selected style if this column is selected
		style := styleColumnHeader
		if i == m.selectedColumn {
			style = styleColumnHeaderSelected
		}

		// Each column gets equal width
		colWidth := m.boardWidth / len(visibleColumns)
		headers = append(headers, style.Width(colWidth).Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, headers...)
}

// renderColumns renders all columns with their cards
func (m Model) renderColumns(contentHeight int) string {
	var columns []string
	visibleColumns := m.getVisibleColumns()
	colWidth := m.boardWidth / len(visibleColumns)

	for i, col := range visibleColumns {
		columnContent := m.renderColumn(col, i, contentHeight, colWidth)
		columns = append(columns, columnContent)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

// renderColumn renders a single column with its cards using Solitaire-style stacking
func (m Model) renderColumn(col Column, colIndex int, contentHeight int, colWidth int) string {
	var columnContent strings.Builder

	// Check if we should show drop indicator in this column
	showDropIndicator := m.draggingCard != nil && m.dropTargetColumn == colIndex

	// Empty column
	if len(col.Cards) == 0 {
		if showDropIndicator && m.dropTargetIndex == 0 {
			// Show drop indicator at top of empty column
			dropLine := strings.Repeat("─", cardWidth)
			columnContent.WriteString(styleDropIndicator.Render(dropLine) + "\n")
		}

		return lipgloss.NewStyle().
			Width(colWidth).
			Height(contentHeight).
			Align(lipgloss.Center).
			Render(columnContent.String())
	}

	// Calculate how many cards we can show with stacking
	maxStackedCards := (contentHeight - cardHeight) / 2
	if maxStackedCards < 0 {
		maxStackedCards = 0
	}

	cardsToShow := len(col.Cards)
	if cardsToShow > maxStackedCards+1 {
		cardsToShow = maxStackedCards + 1
	}

	startIndex := len(col.Cards) - cardsToShow
	if startIndex < 0 {
		startIndex = 0
	}

	for i := startIndex; i < len(col.Cards); i++ {
		// Show drop indicator before this card if needed
		if showDropIndicator && m.dropTargetIndex == i {
			dropLine := strings.Repeat("─", cardWidth)
			columnContent.WriteString(styleDropIndicator.Render(dropLine) + "\n")
		}

		card := col.Cards[i]
		isLast := i == len(col.Cards)-1
		isSelected := colIndex == m.selectedColumn && i == m.selectedCard

		// Check if this is the card being dragged
		isDragging := m.draggingCard != nil && m.dragFromColumn == colIndex && i == m.dragFromIndex

		if isLast {
			// Last card - show full card
			if isDragging {
				columnContent.WriteString(renderCardGhost(card.Title))
			} else {
				columnContent.WriteString(renderCard(card.Title, isSelected))
			}
		} else {
			// Stacked card - show only top 2 lines
			if isDragging {
				columnContent.WriteString(renderCardTopLinesGhost(card.Title))
			} else {
				columnContent.WriteString(renderCardTopLines(card.Title, isSelected))
			}
			columnContent.WriteString("\n")
		}
	}

	// Show drop indicator at end if needed
	if showDropIndicator && m.dropTargetIndex == len(col.Cards) {
		dropLine := strings.Repeat("─", cardWidth)
		columnContent.WriteString("\n" + styleDropIndicator.Render(dropLine))
	}

	// If we couldn't show all cards, indicate how many are hidden
	if startIndex > 0 {
		hidden := fmt.Sprintf("(%d more above)", startIndex)
		columnContent.WriteString("\n" + styleSubdued.Render(hidden))
	}

	return lipgloss.NewStyle().
		Width(colWidth).
		Height(contentHeight).
		Align(lipgloss.Center).
		Render(columnContent.String())
}

// renderDetailPanel renders the card detail panel
func (m Model) renderDetailPanel() string {
	card := m.getCurrentCard()

	if card == nil {
		return styleDetailPanel.
			Width(m.detailWidth).
			Height(m.getContentHeight()).
			Render("No card selected")
	}

	var details []string

	// Card title
	details = append(details, styleDetailTitle.Render(card.Title))
	details = append(details, "")

	// Description
	if card.Description != "" {
		details = append(details, styleDetailLabel.Render("Description:"))
		details = append(details, wrapText(card.Description, m.detailWidth-4))
		details = append(details, "")
	}

	// Tags
	if len(card.Tags) > 0 {
		details = append(details, styleDetailLabel.Render("Tags:"))
		tagStr := ""
		for _, tag := range card.Tags {
			tagStr += styleTag.Render("#"+tag)
		}
		details = append(details, tagStr)
		details = append(details, "")
	}

	// Assignee
	if card.Assignee != "" {
		details = append(details, styleDetailLabel.Render("Assigned: ")+styleDetailValue.Render(card.Assignee))
	}

	// Due date
	if card.DueDate != "" {
		details = append(details, styleDetailLabel.Render("Due: ")+styleDetailValue.Render(card.DueDate))
	}

	// Timestamps
	details = append(details, "")
	details = append(details, styleDetailLabel.Render("Created: ")+styleDetailValue.Render(card.CreatedAt.Format("Jan 2, 2006")))
	details = append(details, styleDetailLabel.Render("Modified: ")+styleDetailValue.Render(card.ModifiedAt.Format("Jan 2, 2006")))

	content := strings.Join(details, "\n")

	return styleDetailPanel.
		Width(m.detailWidth).
		Height(m.getContentHeight()).
		Render(content)
}

// renderStatus renders the status bar
func (m Model) renderStatus() string {
	var help string

	switch m.viewMode {
	case ViewProjectList:
		help = "↑/↓: Navigate | Enter: Open project | q: Quit"
	case ViewBoard:
		archiveStatus := "hidden"
		if m.showArchive {
			archiveStatus = "visible"
		}
		projectsDebug := fmt.Sprintf("[%d projects]", len(m.projects))
		help = fmt.Sprintf("←/→: Columns | ↑/↓: Cards | Tab: Details | a: Archive (%s) | p: Projects %s | q: Quit", archiveStatus, projectsDebug)
	default:
		help = "q: Quit"
	}

	return styleStatus.
		Width(m.width).
		Render(help)
}

// renderProjectListView renders the project selection list
func (m Model) renderProjectListView() string {
	var sections []string

	// Title
	title := styleTitle.Width(m.width).Render("📋 tkan - Select Project")
	sections = append(sections, title)

	// Project list
	var projectLines []string
	projectLines = append(projectLines, "")
	projectLines = append(projectLines, styleDetailTitle.Render("Available Projects:"))
	projectLines = append(projectLines, "")

	cwd, _ := os.Getwd()
	for i, project := range m.projects {
		relPath := GetProjectRelativePath(project, cwd)

		// Format: [*] Project Name (path)
		prefix := "   "
		var style lipgloss.Style
		if i == m.selectedProject {
			prefix = " ▶ "
			// Selected project style - no width constraint
			style = lipgloss.NewStyle().
				Foreground(colorSelected).
				Bold(true)
		} else {
			style = styleDetailValue
		}

		line := fmt.Sprintf("%s%s", prefix, project.Name)
		if relPath != "" {
			line += styleSubdued.Render(fmt.Sprintf(" (%s)", relPath))
		}

		projectLines = append(projectLines, style.Render(line))
	}

	projectLines = append(projectLines, "")
	projectLines = append(projectLines, styleSubdued.Render("Press Enter to open the selected project"))

	content := strings.Join(projectLines, "\n")

	// Center the content
	contentStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height - 4).
		Padding(2, 4)

	sections = append(sections, contentStyle.Render(content))

	// Status bar
	sections = append(sections, m.renderStatus())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// wrapText wraps text to fit within the given width
func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}

// renderHelpView renders the help screen
func (m Model) renderHelpView() string {
	var sections []string

	// Title
	title := styleTitle.Width(m.width).Render("📚 tkan - Help")
	sections = append(sections, title)

	// Help content
	helpContent := `
NAVIGATION
  ←/→ or h/l     Move between columns
  ↑/↓ or k/j     Move between cards in a column
  Home or g      Jump to first column
  End or G       Jump to last column

ACTIONS
  n              Create new card
  e              Edit selected card
  d              Delete selected card
  m              Move card to different column
  Mouse drag     Drag & drop cards between columns

VIEWS
  Tab            Toggle detail panel
  v              Switch to table view (not yet implemented)
  a              Toggle archive column visibility
  p              Back to project list (if multiple projects)
  ?              Toggle this help screen

SEARCH & FILTER
  /              Search/filter cards (not yet implemented)

GENERAL
  q or Ctrl+C    Quit tkan

MOUSE SUPPORT
  Click          Select card
  Double-click   Edit card (not yet implemented)
  Drag & drop    Move cards between columns
                 (150ms hold delay to prevent accidental drags)

Press ? or Esc or Enter or Space to close this help screen`

	contentStyle := lipgloss.NewStyle().
		Width(m.width - 8).
		Padding(2, 4)

	sections = append(sections, contentStyle.Render(helpContent))

	// Status bar
	status := styleStatus.Width(m.width).Render("Press ? or Esc or Enter or Space to close")
	sections = append(sections, status)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderFormOverlay renders the card creation/editing form as an overlay
func (m Model) renderFormOverlay(background string) string {
	// Determine form title
	formTitle := "Create New Card"
	if m.formMode == FormEditCard {
		formTitle = "Edit Card"
	}

	// Build form content
	var formLines []string
	formLines = append(formLines, styleDetailTitle.Render(formTitle))
	formLines = append(formLines, "")

	// Render form inputs
	formLines = append(formLines, styleDetailLabel.Render("Title:"))
	if len(m.formInputs) > 0 {
		formLines = append(formLines, m.formInputs[0].View())
	}
	formLines = append(formLines, "")

	formLines = append(formLines, styleDetailLabel.Render("Description:"))
	if len(m.formInputs) > 1 {
		formLines = append(formLines, m.formInputs[1].View())
	}
	formLines = append(formLines, "")

	// Instructions
	formLines = append(formLines, styleSubdued.Render("Tab/↑/↓: Navigate fields"))
	formLines = append(formLines, styleSubdued.Render("Enter: Next field / Save (on last field)"))
	formLines = append(formLines, styleSubdued.Render("Ctrl+S / Ctrl+Enter: Save"))
	formLines = append(formLines, styleSubdued.Render("Esc: Cancel"))

	formContent := strings.Join(formLines, "\n")

	// Style the form as a centered modal
	formBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		Width(70).
		Render(formContent)

	// Center the form
	centeredForm := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		formBox,
	)

	// Layer the form over the background (simple approach - just replace)
	return centeredForm
}
