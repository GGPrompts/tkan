package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

// handleMouseMsg handles mouse input
func (m Model) handleMouseMsg(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Mouse handling will be implemented in Phase 2 (drag & drop)
	// For now, just return the model unchanged
	return m, nil
}
