package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles all messages and updates the model (required by Bubbletea)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case tea.MouseMsg:
		return m.handleMouseMsg(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.calculateLayout()
		return m, nil

	case boardLoadedMsg:
		if msg.err != nil {
			// Handle error - for now just keep current board
			return m, nil
		}
		m.board = msg.board
		return m, nil
	}

	return m, nil
}
