package main

import (
	"time"

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
	// Only update drop target if actually dragging (not just potential drag)
	if m.draggingCard != nil {
		// Update drop target for visual feedback
		colIndex, insertIndex := m.getDropPosition(msg.X, msg.Y)
		m.dropTargetColumn = colIndex
		m.dropTargetIndex = insertIndex
	}

	// If mouse moved significantly while in potential drag, cancel the selection
	// and start dragging immediately (without waiting for timer)
	if m.potentialDrag && m.mouseHeldDown {
		dx := msg.X - m.mousePressX
		dy := msg.Y - m.mousePressY
		distanceMoved := dx*dx + dy*dy

		// If moved more than 4 pixels, start dragging immediately
		if distanceMoved > 4 {
			// Get the card to drag
			visibleColumns := m.getVisibleColumns()
			if m.dragFromColumn < len(visibleColumns) {
				col := visibleColumns[m.dragFromColumn]
				if m.dragFromIndex < len(col.Cards) {
					m.draggingCard = col.Cards[m.dragFromIndex]
					m.dropTargetColumn = m.dragFromColumn
					m.dropTargetIndex = m.dragFromIndex
					m.potentialDrag = false
				}
			}
		}
	}

	return m, nil
}

// handleMousePress handles mouse button press (start potential drag)
func (m Model) handleMousePress(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Store press position and time
	m.mousePressX = msg.X
	m.mousePressY = msg.Y
	m.dragStartTime = time.Now()
	m.mouseHeldDown = true

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

	// Immediately select the card (visual feedback)
	m.selectedColumn = colIndex
	m.selectedCard = cardIndex

	// Store potential drag info but don't start dragging yet
	m.potentialDrag = true
	m.dragFromColumn = colIndex
	m.dragFromIndex = cardIndex

	// Start a timer to initiate drag after delay
	return m, tickCmd()
}

// handleMouseRelease handles mouse button release (drop or click)
func (m Model) handleMouseRelease(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Clear mouse held state
	m.mouseHeldDown = false

	// If we were actually dragging, handle the drop
	if m.draggingCard != nil {
		// Get drop position
		toColIndex, insertIndex := m.getDropPosition(msg.X, msg.Y)

		if toColIndex != -1 {
			// Move card to the target position
			m.moveCard(m.dragFromColumn, m.dragFromIndex, toColIndex, insertIndex)
		}

		// Clear drag state
		m.draggingCard = nil
		m.dropTargetColumn = -1
		m.dropTargetIndex = -1
	}

	// Clear potential drag state
	m.potentialDrag = false
	m.dragFromColumn = -1
	m.dragFromIndex = -1

	return m, nil
}
