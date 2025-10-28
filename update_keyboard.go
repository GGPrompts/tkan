package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyMsg handles keyboard input
func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global shortcuts
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

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

	// Future: These will be implemented in later phases
	case "n":
		// New card
		return m, nil

	case "e":
		// Edit card
		return m, nil

	case "d":
		// Delete card
		return m, nil

	case "m":
		// Move card
		return m, nil

	case "/":
		// Search/filter
		return m, nil

	case "?":
		// Help
		return m, nil
	}

	return m, nil
}

// handleTableKeyMsg handles keyboard input for table view (Phase 2)
func (m Model) handleTableKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Not implemented in Phase 1
	return m, nil
}
