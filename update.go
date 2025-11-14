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

	case dragStartMsg:
		// Timer expired - start dragging if mouse still held
		if m.mouseHeldDown && m.potentialDrag {
			// Get the card to drag
			visibleColumns := m.getVisibleColumns()
			if m.dragFromColumn < len(visibleColumns) {
				col := visibleColumns[m.dragFromColumn]
				if m.dragFromIndex < len(col.Cards) {
					m.draggingCard = col.Cards[m.dragFromIndex]
					m.dropTargetColumn = m.dragFromColumn
					m.dropTargetIndex = m.dragFromIndex
				}
			}
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.calculateLayout()
		// Rebuild table if in table view
		if m.viewMode == ViewTable && m.table != nil {
			m.buildTable()
		}
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
