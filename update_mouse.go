package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseMsg handles mouse input
func (m Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Only handle mouse in board view
	if m.viewMode != ViewBoard {
		return m, nil
	}

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionPress {
			return m.handleMousePress(msg)
		} else if msg.Action == tea.MouseActionRelease {
			return m.handleMouseRelease(msg)
		} else if msg.Action == tea.MouseActionMotion && m.draggingCard != nil {
			// Motion while left button held (dragging)
			return m.handleMouseMotion(msg)
		}
	case tea.MouseButtonNone:
		// Mouse motion without button - update drop target if dragging
		if msg.Action == tea.MouseActionMotion && m.draggingCard != nil {
			return m.handleMouseMotion(msg)
		}
	}

	return m, nil
}

// handleMouseMotion updates the drop target during drag
func (m Model) handleMouseMotion(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Update drop target for visual feedback
	colIndex, insertIndex := m.getDropPosition(msg.X, msg.Y)
	m.dropTargetColumn = colIndex
	m.dropTargetIndex = insertIndex

	return m, nil
}

// handleMousePress handles mouse button press (start potential drag)
func (m Model) handleMousePress(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Store press position to detect click vs drag later
	m.mousePressX = msg.X
	m.mousePressY = msg.Y

	// Get column
	colIndex := m.getColumnAtPosition(msg.X, msg.Y)
	if colIndex == -1 {
		return m, nil
	}

	// Get the visible columns
	visibleColumns := m.getVisibleColumns()
	if colIndex >= len(visibleColumns) {
		return m, nil
	}

	col := visibleColumns[colIndex]
	if len(col.Cards) == 0 {
		return m, nil // Can't drag from empty column
	}

	// Calculate which card was clicked
	const cardAreaStartY = 3
	relY := msg.Y - cardAreaStartY
	cardIndex := m.getCardIndexInColumn(col, relY)

	if cardIndex < 0 || cardIndex >= len(col.Cards) {
		return m, nil
	}

	// Prepare for potential drag
	m.draggingCard = col.Cards[cardIndex]
	m.dragFromColumn = colIndex
	m.dragFromIndex = cardIndex

	// Initialize drop target to current position (will update on motion)
	m.dropTargetColumn = colIndex
	m.dropTargetIndex = cardIndex

	return m, nil
}

// handleMouseRelease handles mouse button release (drop or click)
func (m Model) handleMouseRelease(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Not dragging anything
	if m.draggingCard == nil {
		return m, nil
	}

	// Calculate distance moved since press
	dx := msg.X - m.mousePressX
	dy := msg.Y - m.mousePressY
	distanceMoved := dx*dx + dy*dy

	// Threshold: less than 4 pixels² = click, otherwise = drag
	if distanceMoved < 4 {
		// This was a click (select card)
		m.selectedColumn = m.dragFromColumn
		m.selectedCard = m.dragFromIndex

		// Update last click time/position for double-click detection
		// (future enhancement: double-click to auto-move to done)
		// m.lastClickTime = time.Now()
		// m.lastClickX = msg.X
		// m.lastClickY = msg.Y
	} else {
		// This was a drag - get drop position
		toColIndex, insertIndex := m.getDropPosition(msg.X, msg.Y)

		if toColIndex != -1 {
			// Move card to the target position
			m.moveCard(m.dragFromColumn, m.dragFromIndex, toColIndex, insertIndex)
		}
	}

	// Clear drag state and drop target
	m.draggingCard = nil
	m.dragFromColumn = -1
	m.dragFromIndex = -1
	m.dropTargetColumn = -1
	m.dropTargetIndex = -1

	return m, nil
}
