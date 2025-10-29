package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// NewModel creates a new Model with the given board and projects
func NewModel(board *Board, projects []Project) Model {
	// Default to local backend for backward compatibility
	backend := NewLocalBackend(".tkan.yaml")
	if len(projects) > 0 {
		backend = NewLocalBackend(projects[0].Path)
	}
	return NewModelWithBackend(board, projects, backend)
}

// NewModelWithBackend creates a new Model with a specific backend
func NewModelWithBackend(board *Board, projects []Project, backend Backend) Model {
	// If we have multiple projects, start in project list view
	// If only one project, go straight to board view
	viewMode := ViewBoard
	if len(projects) > 1 {
		viewMode = ViewProjectList
	}

	return Model{
		board:            board,
		projects:         projects,
		selectedProject:  0,
		backend:          backend,
		viewMode:         viewMode,
		selectedColumn:   0,
		selectedCard:     0,
		showDetails:      true,  // Start with details panel visible
		showArchive:      false, // Hide archive by default
		width:            0,
		height:           0,
		ready:            false,
		dropTargetColumn: -1, // Initialize drop target as invalid
		dropTargetIndex:  -1,
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

// getColumnAtPosition determines which column is at the given screen position
// Returns columnIndex (or -1 if outside board area)
func (m *Model) getColumnAtPosition(x, y int) int {
	// Check if click is within the board area
	if x < 0 || x >= m.boardWidth {
		return -1
	}

	// Get visible columns and calculate column width
	visibleColumns := m.getVisibleColumns()
	if len(visibleColumns) == 0 {
		return -1
	}

	colWidth := m.boardWidth / len(visibleColumns)

	// Determine which column was clicked
	columnIndex := x / colWidth
	if columnIndex < 0 || columnIndex >= len(visibleColumns) {
		return -1
	}

	return columnIndex
}

// getDropPosition determines where a card would be dropped in a column
// Returns columnIndex, insertIndex where insertIndex is the position to insert
// insertIndex = 0 means insert at start, insertIndex = len(cards) means insert at end
// Returns -1, -1 if outside valid drop area
func (m *Model) getDropPosition(x, y int) (columnIndex, insertIndex int) {
	// Layout calculation (must match renderBoardView exactly):
	// Line 0-1: Title bar (2 lines)
	// Line 2: Column headers (1 line)
	// Line 3+: Card area starts here

	const titleHeight = 2    // Title bar
	const headerHeight = 1   // Column headers
	const cardAreaStartY = 3 // Cards start at Y=3

	// Check if click is in the card area
	if y < cardAreaStartY {
		return -1, -1
	}

	// Get column index
	columnIndex = m.getColumnAtPosition(x, y)
	if columnIndex == -1 {
		return -1, -1
	}

	// Calculate relative Y position within card area
	relY := y - cardAreaStartY

	// Get the actual column from visible columns
	visibleColumns := m.getVisibleColumns()
	col := visibleColumns[columnIndex]

	// Empty column - insert at position 0
	if len(col.Cards) == 0 {
		return columnIndex, 0
	}

	// Calculate insertion position based on Y
	insertIndex = m.getInsertIndexInColumn(col, relY)

	return columnIndex, insertIndex
}

// getCardIndexInColumn determines which card in a column was clicked
// based on the Y position relative to the card area start
func (m *Model) getCardIndexInColumn(col Column, relY int) int {
	numCards := len(col.Cards)
	if numCards == 0 {
		return -1
	}

	// Card rendering logic (from view.go renderColumn):
	// - Each stacked card shows 2 lines
	// - Last card shows full 5 lines
	// - We may not show all cards if column is too long

	const cardHeight = 5    // Full card height
	const stackedHeight = 2 // Visible height of stacked cards

	contentHeight := m.getContentHeight()

	// Calculate how many cards are actually visible
	maxStackedCards := (contentHeight - cardHeight) / stackedHeight
	if maxStackedCards < 0 {
		maxStackedCards = 0
	}

	cardsToShow := numCards
	if cardsToShow > maxStackedCards+1 {
		cardsToShow = maxStackedCards + 1
	}

	startIndex := numCards - cardsToShow
	if startIndex < 0 {
		startIndex = 0
	}

	// Calculate which visible card was clicked
	// Each stacked card except the last takes stackedHeight lines
	numStackedCards := cardsToShow - 1
	stackedAreaHeight := numStackedCards * stackedHeight

	if relY < stackedAreaHeight {
		// Clicked on a stacked card
		clickedStackedIndex := relY / stackedHeight
		return startIndex + clickedStackedIndex
	}

	// Check if clicked on the last (full) card
	lastCardStartY := stackedAreaHeight
	lastCardEndY := lastCardStartY + cardHeight

	if relY >= lastCardStartY && relY < lastCardEndY {
		// Clicked on the last card
		return numCards - 1
	}

	return -1 // Clicked below all cards
}

// getInsertIndexInColumn determines where to insert a card in a column
// Returns the index where the card should be inserted (0 = start, len(cards) = end)
func (m *Model) getInsertIndexInColumn(col Column, relY int) int {
	numCards := len(col.Cards)
	if numCards == 0 {
		return 0
	}

	const cardHeight = 5    // Full card height
	const stackedHeight = 2 // Visible height of stacked cards

	contentHeight := m.getContentHeight()

	// Calculate how many cards are actually visible
	maxStackedCards := (contentHeight - cardHeight) / stackedHeight
	if maxStackedCards < 0 {
		maxStackedCards = 0
	}

	cardsToShow := numCards
	if cardsToShow > maxStackedCards+1 {
		cardsToShow = maxStackedCards + 1
	}

	startIndex := numCards - cardsToShow
	if startIndex < 0 {
		startIndex = 0
	}

	// Calculate which visible card the mouse is over
	numStackedCards := cardsToShow - 1
	stackedAreaHeight := numStackedCards * stackedHeight

	if relY < stackedAreaHeight {
		// Mouse is over a stacked card
		cardIndex := startIndex + (relY / stackedHeight)

		// Determine if mouse is in top half or bottom half of the card segment
		posInCard := relY % stackedHeight
		if posInCard < stackedHeight/2 {
			// Top half - insert before this card
			return cardIndex
		} else {
			// Bottom half - insert after this card
			return cardIndex + 1
		}
	}

	// Mouse is over the last (full) card area
	lastCardStartY := stackedAreaHeight
	lastCardEndY := lastCardStartY + cardHeight

	if relY >= lastCardStartY && relY < lastCardEndY {
		// Determine if mouse is in top half or bottom half of the last card
		posInLastCard := relY - lastCardStartY
		if posInLastCard < cardHeight/2 {
			// Top half - insert before last card
			return numCards - 1
		} else {
			// Bottom half - insert after last card (at end)
			return numCards
		}
	}

	// Below all cards - insert at end
	return numCards
}

// moveCard moves a card from one position to another (within or across columns)
func (m *Model) moveCard(fromColIndex, fromCardIndex, toColIndex, insertIndex int) {
	visibleColumns := m.getVisibleColumns()

	// Validate indices
	if fromColIndex < 0 || fromColIndex >= len(visibleColumns) {
		return
	}
	if toColIndex < 0 || toColIndex >= len(visibleColumns) {
		return
	}

	fromCol := visibleColumns[fromColIndex]
	toCol := visibleColumns[toColIndex]

	if fromCardIndex < 0 || fromCardIndex >= len(fromCol.Cards) {
		return
	}

	// Get the card to move
	card := fromCol.Cards[fromCardIndex]

	// Find the actual columns in the board (not just visible columns)
	var fromColPtr, toColPtr *Column
	for i := range m.board.Columns {
		if m.board.Columns[i].Name == fromCol.Name {
			fromColPtr = &m.board.Columns[i]
		}
		if m.board.Columns[i].Name == toCol.Name {
			toColPtr = &m.board.Columns[i]
		}
	}

	if fromColPtr == nil || toColPtr == nil {
		return
	}

	// Handle reordering within the same column
	if fromColIndex == toColIndex {
		// Check if actually moving to a different position
		if fromCardIndex == insertIndex || fromCardIndex+1 == insertIndex {
			return // No effective move
		}

		// Remove card from source position
		fromColPtr.Cards = append(fromColPtr.Cards[:fromCardIndex], fromColPtr.Cards[fromCardIndex+1:]...)

		// Adjust insert index if needed (if we removed a card before the insert position)
		adjustedInsertIndex := insertIndex
		if fromCardIndex < insertIndex {
			adjustedInsertIndex--
		}

		// Insert at new position
		if adjustedInsertIndex >= len(fromColPtr.Cards) {
			fromColPtr.Cards = append(fromColPtr.Cards, card)
		} else {
			fromColPtr.Cards = append(fromColPtr.Cards[:adjustedInsertIndex], append([]*Card{card}, fromColPtr.Cards[adjustedInsertIndex:]...)...)
		}

		m.selectedColumn = toColIndex
		m.selectedCard = adjustedInsertIndex
	} else {
		// Moving to a different column

		// Remove card from source column
		fromColPtr.Cards = append(fromColPtr.Cards[:fromCardIndex], fromColPtr.Cards[fromCardIndex+1:]...)

		// Insert into target column at specified position
		if insertIndex >= len(toColPtr.Cards) {
			toColPtr.Cards = append(toColPtr.Cards, card)
			m.selectedCard = len(toColPtr.Cards) - 1
		} else {
			toColPtr.Cards = append(toColPtr.Cards[:insertIndex], append([]*Card{card}, toColPtr.Cards[insertIndex:]...)...)
			m.selectedCard = insertIndex
		}

		// Update card's column field
		card.Column = toCol.Name
		m.selectedColumn = toColIndex
	}

	// Update modification time
	card.ModifiedAt = time.Now()

	// Save changes using backend
	if m.backend != nil {
		// For GitHub backend, update the card's column
		m.backend.MoveCard(card.ID, toCol.Name)
		// For local backend, save the entire board
		m.backend.SaveBoard(m.board)
	}
}
