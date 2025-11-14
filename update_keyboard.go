package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyMsg handles keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle form input first if form is open
	if m.formMode != FormNone {
		return m.handleFormKeyMsg(msg)
	}

	// Handle delete confirmation
	if m.confirmingDelete {
		return m.handleDeleteConfirmation(msg)
	}

	// Global shortcuts
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "?":
		// Toggle help screen
		if m.viewMode == ViewHelp {
			// Return to previous view
			m.viewMode = m.previousView
		} else {
			// Show help
			m.previousView = m.viewMode
			m.viewMode = ViewHelp
		}
		return m, nil

	case "tab":
		m.toggleDetails()
		return m, nil

	case "v", "V":
		// Toggle view mode (board <-> table)
		if m.viewMode == ViewBoard {
			m.viewMode = ViewTable
			m.buildTable() // Build table when switching to table view
		} else if m.viewMode == ViewTable {
			m.viewMode = ViewBoard
		}
		return m, nil
	}

	// View-specific shortcuts
	switch m.viewMode {
	case ViewProjectList:
		return m.handleProjectListKeyMsg(msg)
	case ViewBoard:
		return m.handleBoardKeyMsg(msg)
	case ViewTable:
		return m.handleTableKeyMsg(msg)
	case ViewHelp:
		return m.handleHelpKeyMsg(msg)
	case ViewProjectSource:
		return m.handleProjectSourceKeyMsg(msg)
	}

	return m, nil
}

// handleProjectListKeyMsg handles keyboard input for project list view
func (m Model) handleProjectListKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.selectedProject > 0 {
			m.selectedProject--
		}
		return m, nil

	case "down", "j":
		if m.selectedProject < len(m.projects)-1 {
			m.selectedProject++
		}
		return m, nil

	case "enter":
		// Load the selected project
		if err := m.loadSelectedProject(); err != nil {
			// TODO: Show error message
			return m, nil
		}
		return m, nil
	}

	return m, nil
}

// handleBoardKeyMsg handles keyboard input for board view
func (m Model) handleBoardKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Navigation
	case "left", "h":
		m.moveSelectionLeft()
		return m, nil

	case "right", "l":
		m.moveSelectionRight()
		return m, nil

	case "up", "k":
		m.moveSelectionUp()
		return m, nil

	case "down", "j":
		m.moveSelectionDown()
		return m, nil

	// Jump to first/last column
	case "home", "g":
		m.selectedColumn = 0
		m.selectedCard = 0
		return m, nil

	case "end", "G":
		visibleColumns := m.getVisibleColumns()
		m.selectedColumn = len(visibleColumns) - 1
		m.selectedCard = 0
		return m, nil

	// Toggle archive column
	case "a", "A":
		m.toggleArchive()
		// Adjust selected column if needed
		visibleColumns := m.getVisibleColumns()
		if m.selectedColumn >= len(visibleColumns) {
			m.selectedColumn = len(visibleColumns) - 1
			m.selectedCard = 0
		}
		return m, nil

	// Back to project list or show source selector
	case "p", "P":
		if len(m.projects) > 1 {
			// Multiple projects - show project list
			m.viewMode = ViewProjectList
		} else {
			// 0 or 1 project - show source selector
			m.selectedSourceOpt = 0
			m.viewMode = ViewProjectSource
		}
		return m, nil

	// Card creation and editing
	case "n":
		// New card
		m.openCreateCardForm()
		return m, nil

	case "e":
		// Edit card
		m.openEditCardForm()
		return m, nil

	case "d":
		// Show delete confirmation
		card := m.getCurrentCard()
		if card != nil {
			m.confirmingDelete = true
			m.deletingCardID = card.ID
		}
		return m, nil

	case "m":
		// Move card
		return m, nil

	case "/":
		// Search/filter
		return m, nil
	}

	return m, nil
}

// handleTableKeyMsg handles keyboard input for table view
func (m Model) handleTableKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Make sure table exists
	if m.table == nil {
		return m, nil
	}

	switch msg.String() {
	// Navigation
	case "up", "k":
		m.table.CursorUp()
		return m, nil

	case "down", "j":
		m.table.CursorDown()
		return m, nil

	case "left", "h":
		m.table.CursorLeft()
		return m, nil

	case "right", "l":
		m.table.CursorRight()
		return m, nil

	// Sorting
	case "ctrl+s":
		x, _ := m.table.GetCursorLocation()
		_, order := m.table.GetOrder()
		if order == 1 { // Ascending
			m.table.OrderByDesc(x)
		} else {
			m.table.OrderByAsc(x)
		}
		return m, nil

	// Toggle archive
	case "a":
		m.toggleArchive()
		m.buildTable() // Rebuild table with new filter
		return m, nil

	// Edit selected card
	case "e":
		m.openEditCardForm()
		return m, nil

	// Delete selected card - show confirmation
	case "d":
		card := m.getSelectedCardInTable()
		if card != nil {
			m.confirmingDelete = true
			m.deletingCardID = card.ID
		}
		return m, nil

	// Filter (typing alphanumerics)
	default:
		if len(msg.String()) == 1 {
			r := msg.Runes[0]
			// Check if it's alphanumeric
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
				m.handleTableFilter(msg.String())
				return m, nil
			}
		}

		// Backspace for filter
		if msg.String() == "backspace" {
			m.handleTableFilter(msg.String())
			return m, nil
		}
	}

	return m, nil
}

// handleHelpKeyMsg handles keyboard input for help view
func (m Model) handleHelpKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter", " ":
		// Return to previous view
		m.viewMode = m.previousView
		return m, nil
	}

	return m, nil
}

// handleFormKeyMsg handles keyboard input when card form is open
func (m Model) handleFormKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Check if this is GitHub owner input mode
	if m.formMode == FormMode(100) {
		switch msg.String() {
		case "esc":
			// Cancel - go back to source selector
			m.formMode = FormNone
			m.formInputs = nil
			m.viewMode = ViewProjectSource
			return m, nil

		case "enter":
			// Load GitHub projects for entered owner
			owner := ""
			if len(m.formInputs) > 0 {
				owner = m.formInputs[0].Value()
			}
			m.formMode = FormNone
			m.formInputs = nil
			if owner != "" {
				return m.loadGitHubProjects(owner)
			}
			m.viewMode = ViewProjectSource
			return m, nil
		}

		// Update the text input
		if len(m.formInputs) > 0 {
			m.formInputs[0], cmd = m.formInputs[0].Update(msg)
		}
		return m, cmd
	}

	// Regular card form handling
	switch msg.String() {
	case "esc":
		// Cancel form
		m.closeCardForm()
		return m, nil

	case "ctrl+s", "ctrl+enter":
		// Save form
		m.saveCardForm()
		return m, nil

	case "tab", "shift+tab", "up", "down":
		// Navigate between form fields
		if msg.String() == "tab" || msg.String() == "down" {
			m.formFocusIndex++
			if m.formFocusIndex >= len(m.formInputs) {
				m.formFocusIndex = 0
			}
		} else {
			m.formFocusIndex--
			if m.formFocusIndex < 0 {
				m.formFocusIndex = len(m.formInputs) - 1
			}
		}

		// Update focus
		for i := range m.formInputs {
			if i == m.formFocusIndex {
				m.formInputs[i].Focus()
			} else {
				m.formInputs[i].Blur()
			}
		}

		return m, nil

	case "enter":
		// Enter on last field saves the form
		if m.formFocusIndex == len(m.formInputs)-1 {
			m.saveCardForm()
			return m, nil
		}
		// Otherwise move to next field
		m.formFocusIndex++
		if m.formFocusIndex >= len(m.formInputs) {
			m.formFocusIndex = 0
		}
		for i := range m.formInputs {
			if i == m.formFocusIndex {
				m.formInputs[i].Focus()
			} else {
				m.formInputs[i].Blur()
			}
		}
		return m, nil
	}

	// Update the focused text input
	if m.formFocusIndex >= 0 && m.formFocusIndex < len(m.formInputs) {
		m.formInputs[m.formFocusIndex], cmd = m.formInputs[m.formFocusIndex].Update(msg)
	}

	return m, cmd
}

// handleProjectSourceKeyMsg handles keyboard input for project source selection view
func (m Model) handleProjectSourceKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel - go back to board
		m.viewMode = ViewBoard
		return m, nil

	case "up", "k":
		if m.selectedSourceOpt > 0 {
			m.selectedSourceOpt--
		}
		return m, nil

	case "down", "j":
		if m.selectedSourceOpt < 3 { // 4 options (0-3)
			m.selectedSourceOpt++
		}
		return m, nil

	case "enter":
		return m.handleProjectSourceSelection()
	}

	return m, nil
}

// handleTableFilter handles filtering in table view
func (m *Model) handleTableFilter(key string) {
	i, s := m.table.GetFilter()
	x, _ := m.table.GetCursorLocation()

	// If we're on a different column, start fresh
	if x != i && key != "backspace" {
		m.table.SetFilter(x, key)
		return
	}

	// Handle backspace
	if key == "backspace" {
		if len(s) == 1 {
			m.table.UnsetFilter()
			return
		} else if len(s) > 1 {
			s = s[0 : len(s)-1]
		} else {
			return
		}
	} else {
		// Append character to filter
		s = s + key
	}

	m.table.SetFilter(i, s)
}

// handleDeleteConfirmation handles keyboard input when showing delete confirmation
func (m Model) handleDeleteConfirmation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm delete
		m.confirmDelete()
		m.confirmingDelete = false
		m.deletingCardID = ""
		return m, nil

	case "n", "N", "esc":
		// Cancel delete
		m.confirmingDelete = false
		m.deletingCardID = ""
		return m, nil
	}

	return m, nil
}

// confirmDelete actually deletes the card after confirmation
func (m *Model) confirmDelete() {
	if m.deletingCardID == "" {
		return
	}

	// Find and delete the card from the board
	for i, c := range m.board.Cards {
		if c.ID == m.deletingCardID {
			m.board.Cards = append(m.board.Cards[:i], m.board.Cards[i+1:]...)
			break
		}
	}

	// Remove from columns
	for i := range m.board.Columns {
		for j, c := range m.board.Columns[i].Cards {
			if c.ID == m.deletingCardID {
				m.board.Columns[i].Cards = append(m.board.Columns[i].Cards[:j], m.board.Columns[i].Cards[j+1:]...)
				break
			}
		}
	}

	// Save changes
	if m.backend != nil {
		m.backend.SaveBoard(m.board)
	}

	// Rebuild table if in table view
	if m.viewMode == ViewTable {
		m.buildTable()
	}

	// Adjust selection if in board view
	if m.viewMode == ViewBoard {
		col := m.getCurrentColumn()
		if col != nil && m.selectedCard >= len(col.Cards) && m.selectedCard > 0 {
			m.selectedCard--
		}
	}
}
