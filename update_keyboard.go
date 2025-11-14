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
		} else {
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

	// Back to project list (if multiple projects)
	case "p", "P":
		if len(m.projects) > 1 {
			m.viewMode = ViewProjectList
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
		// Delete card
		m.deleteCard()
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

// handleTableKeyMsg handles keyboard input for table view (Phase 2)
func (m Model) handleTableKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Not implemented in Phase 1
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
