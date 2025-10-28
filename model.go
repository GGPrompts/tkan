package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// NewModel creates a new Model with the given board and projects
func NewModel(board *Board, projects []Project) Model {
	// If we have multiple projects, start in project list view
	// If only one project, go straight to board view
	viewMode := ViewBoard
	if len(projects) > 1 {
		viewMode = ViewProjectList
	}

	return Model{
		board:           board,
		projects:        projects,
		selectedProject: 0,
		viewMode:        viewMode,
		selectedColumn:  0,
		selectedCard:    0,
		showDetails:     true,  // Start with details panel visible
		showArchive:     false, // Hide archive by default
		width:           0,
		height:          0,
		ready:           false,
	}
}

// Init initializes the model (required by Bubbletea)
func (m Model) Init() tea.Cmd {
	return nil
}

// setSize updates the model dimensions and recalculates layout
func (m *Model) setSize(width, height int) {
	m.width = width
	m.height = height
	m.calculateLayout()
}

// calculateLayout computes the board and detail panel widths
// Phase 1: Simple toggle between 67%/33% (details shown) or 100%/0% (details hidden)
func (m *Model) calculateLayout() {
	if m.showDetails {
		// Detail panel visible: Board gets 67%, detail gets 33%
		m.detailWidth = m.width / 3
		m.boardWidth = m.width - m.detailWidth - 1 // -1 for divider
	} else {
		// Detail panel hidden: Board gets 100%
		m.boardWidth = m.width
		m.detailWidth = 0
	}
}

// getContentHeight returns the height available for content (excluding title and status bars)
func (m Model) getContentHeight() int {
	contentHeight := m.height
	contentHeight -= 3 // Title bar (1) + separator (1) + column headers (1)
	contentHeight -= 2 // Status bar (1) + bottom border (1)
	return contentHeight
}

// getCurrentColumn returns the currently selected column
func (m Model) getCurrentColumn() *Column {
	if m.selectedColumn >= 0 && m.selectedColumn < len(m.board.Columns) {
		return &m.board.Columns[m.selectedColumn]
	}
	return nil
}

// getCurrentCard returns the currently selected card
func (m Model) getCurrentCard() *Card {
	col := m.getCurrentColumn()
	if col == nil {
		return nil
	}

	if m.selectedCard >= 0 && m.selectedCard < len(col.Cards) {
		return col.Cards[m.selectedCard]
	}
	return nil
}

// moveSelectionLeft moves the selection to the previous column
func (m *Model) moveSelectionLeft() {
	if m.selectedColumn > 0 {
		m.selectedColumn--
		m.selectedCard = 0 // Reset card selection in new column
	}
}

// moveSelectionRight moves the selection to the next column
func (m *Model) moveSelectionRight() {
	if m.selectedColumn < len(m.board.Columns)-1 {
		m.selectedColumn++
		m.selectedCard = 0 // Reset card selection in new column
	}
}

// moveSelectionUp moves the selection to the previous card in the column
func (m *Model) moveSelectionUp() {
	if m.selectedCard > 0 {
		m.selectedCard--
	}
}

// moveSelectionDown moves the selection to the next card in the column
func (m *Model) moveSelectionDown() {
	col := m.getCurrentColumn()
	if col != nil && m.selectedCard < len(col.Cards)-1 {
		m.selectedCard++
	}
}

// toggleDetails toggles the visibility of the detail panel
func (m *Model) toggleDetails() {
	m.showDetails = !m.showDetails
	m.calculateLayout()
}

// toggleArchive toggles the visibility of the archive column
func (m *Model) toggleArchive() {
	m.showArchive = !m.showArchive
}

// getVisibleColumns returns columns to display (excludes ARCHIVE if showArchive is false)
func (m Model) getVisibleColumns() []Column {
	if m.showArchive {
		return m.board.Columns
	}

	// Filter out ARCHIVE column
	var visible []Column
	for _, col := range m.board.Columns {
		if col.Name != "ARCHIVE" {
			visible = append(visible, col)
		}
	}
	return visible
}

// loadSelectedProject loads the currently selected project
func (m *Model) loadSelectedProject() error {
	if m.selectedProject < 0 || m.selectedProject >= len(m.projects) {
		return nil
	}

	project := m.projects[m.selectedProject]
	board, err := LoadBoard(project.Path)
	if err != nil {
		return err
	}

	m.board = board
	m.viewMode = ViewBoard
	m.selectedColumn = 0
	m.selectedCard = 0
	return nil
}
