package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette - soft, professional colors
var (
	// Base colors
	colorBackground = lipgloss.Color("235") // Dark gray
	colorForeground = lipgloss.Color("252") // Light gray
	colorBorder     = lipgloss.Color("240") // Medium gray
	colorDivider    = lipgloss.Color("237") // Slightly lighter than background

	// Accent colors
	colorPrimary   = lipgloss.Color("75")  // Blue
	colorSecondary = lipgloss.Color("140") // Purple
	colorSuccess   = lipgloss.Color("76")  // Green
	colorWarning   = lipgloss.Color("220") // Yellow
	colorDanger    = lipgloss.Color("203") // Red
	colorInfo      = lipgloss.Color("117") // Cyan

	// UI element colors
	colorTitle     = lipgloss.Color("111") // Bright blue
	colorSelected  = lipgloss.Color("212") // Pink/magenta
	colorSubdued   = lipgloss.Color("243") // Subdued gray
	colorHighlight = lipgloss.Color("229") // Bright yellow
)

// Base styles
var (
	// Title bar style
	styleTitle = lipgloss.NewStyle().
			Foreground(colorTitle).
			Bold(true).
			Padding(0, 1)

	// Status bar style
	styleStatus = lipgloss.NewStyle().
			Foreground(colorSubdued).
			Padding(0, 1)

	// Subdued text style (for secondary info)
	styleSubdued = lipgloss.NewStyle().
			Foreground(colorSubdued)

	// Panel border style
	stylePanelBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBorder).
				Padding(0, 1)

	// Divider style (vertical line between panels)
	styleDivider = lipgloss.NewStyle().
			Foreground(colorDivider)
)

// Column styles
var (
	// Column header style (centered, bold)
	styleColumnHeader = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true).
				Align(lipgloss.Center).
				Padding(0, 1)

	// Selected column header style
	styleColumnHeaderSelected = lipgloss.NewStyle().
					Foreground(colorSelected).
					Bold(true).
					Align(lipgloss.Center).
					Padding(0, 1)
)

// Card styles (12 chars wide × 5 lines tall - Solitaire-style)
var (
	cardWidth  = 12
	cardHeight = 5

	// Normal card style
	styleCard = lipgloss.NewStyle().
			Width(cardWidth).
			Height(cardHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(0, 1)

	// Selected card style
	styleCardSelected = lipgloss.NewStyle().
				Width(cardWidth).
				Height(cardHeight).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorSelected).
				Foreground(colorSelected).
				Bold(true).
				Padding(0, 1)

	// Drop indicator style (thin line showing where card will be dropped)
	styleDropIndicator = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)

	// Ghost card style (for card being dragged)
	styleCardGhost = lipgloss.NewStyle().
			Width(cardWidth).
			Height(cardHeight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubdued).
			Foreground(colorSubdued).
			Padding(0, 1)

	// Card content style (for text inside cards)
	styleCardContent = lipgloss.NewStyle().
				Width(cardWidth - 2). // Account for padding
				Foreground(colorForeground)
)

// Detail panel styles
var (
	styleDetailPanel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBorder).
				Padding(1, 2)

	styleDetailTitle = lipgloss.NewStyle().
				Foreground(colorTitle).
				Bold(true).
				Underline(true)

	styleDetailLabel = lipgloss.NewStyle().
				Foreground(colorSubdued).
				Bold(true)

	styleDetailValue = lipgloss.NewStyle().
				Foreground(colorForeground)

	styleTag = lipgloss.NewStyle().
			Foreground(colorInfo).
			Background(lipgloss.Color("238")).
			Padding(0, 1).
			MarginRight(1)
)

// Helper functions for styling

// renderCard renders a card with the given title (wrapped, no tags)
// Card format (12×5):
//   ┌──────────┐
//   │Title     │
//   │wrapped   │
//   │here      │
//   └──────────┘
func renderCard(title string, selected bool) string {
	return renderCardWithStyle(title, selected, false)
}

// renderCardGhost renders a faded ghost card (for dragging)
func renderCardGhost(title string) string {
	return renderCardWithStyle(title, false, true)
}

// renderCardWithStyle renders a card with the given title and style options
func renderCardWithStyle(title string, selected bool, ghost bool) string {
	style := styleCard
	if ghost {
		style = styleCardGhost
	} else if selected {
		style = styleCardSelected
	}

	// Wrap title to fit card width (10 chars with padding)
	maxWidth := cardWidth - 2
	wrappedTitle := wrapCardTitle(title, maxWidth)

	return style.Render(wrappedTitle)
}

// renderCardTopLines renders just the top 2 lines of a card (for stacking)
// This creates the Solitaire-style cascading effect
func renderCardTopLines(title string, selected bool) string {
	// Render full card first
	fullCard := renderCardWithStyle(title, selected, false)

	// Extract just the top 2 lines
	lines := strings.Split(fullCard, "\n")
	if len(lines) >= 2 {
		return lines[0] + "\n" + lines[1]
	}
	return fullCard
}

// renderCardTopLinesGhost renders just the top 2 lines of a ghost card
func renderCardTopLinesGhost(title string) string {
	// Render full ghost card first
	fullCard := renderCardGhost(title)

	// Extract just the top 2 lines
	lines := strings.Split(fullCard, "\n")
	if len(lines) >= 2 {
		return lines[0] + "\n" + lines[1]
	}
	return fullCard
}

// wrapCardTitle wraps a title to fit within the card width
func wrapCardTitle(title string, maxWidth int) string {
	if len(title) <= maxWidth {
		return title
	}

	var lines []string
	words := strings.Fields(title)
	currentLine := ""

	for _, word := range words {
		// If word itself is too long, split it
		if len(word) > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = ""
			}
			// Split long word across multiple lines
			for len(word) > maxWidth {
				lines = append(lines, word[:maxWidth])
				word = word[maxWidth:]
			}
			currentLine = word
			continue
		}

		// Try adding word to current line
		testLine := currentLine
		if currentLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= maxWidth {
			currentLine = testLine
		} else {
			// Word doesn't fit, start new line
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Limit to 3 lines (card has 5 total, 2 for border)
	if len(lines) > 3 {
		lines = lines[:3]
		// Add ellipsis to last line if truncated
		lastLine := lines[2]
		if len(lastLine) > maxWidth-1 {
			lines[2] = lastLine[:maxWidth-1] + "…"
		} else {
			lines[2] = lastLine + "…"
		}
	}

	return strings.Join(lines, "\n")
}
