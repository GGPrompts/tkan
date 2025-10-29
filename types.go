package main

import (
	"time"
)

// Card represents a single task card
type Card struct {
	ID          string    `yaml:"id"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Tags        []string  `yaml:"tags,omitempty"`
	Assignee    string    `yaml:"assignee,omitempty"`
	DueDate     string    `yaml:"due_date,omitempty"`
	CreatedAt   time.Time `yaml:"created_at"`
	ModifiedAt  time.Time `yaml:"modified_at"`
	Column      string    `yaml:"column"` // Which column this card belongs to
}

// Column represents a column in the Kanban board
type Column struct {
	Name  string  `yaml:"name"`
	Cards []*Card `yaml:"-"` // Populated at runtime from Board.Cards
}

// Board represents the entire Kanban board
type Board struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description,omitempty"`
	Columns     []Column  `yaml:"columns"`
	Cards       []*Card   `yaml:"cards"`
	CreatedAt   time.Time `yaml:"created_at"`
	ModifiedAt  time.Time `yaml:"modified_at"`
}

// ViewMode represents the current view (project list, board, or table)
type ViewMode int

const (
	ViewProjectList ViewMode = iota
	ViewBoard
	ViewTable
)

// Model is the Bubbletea model for the entire application
type Model struct {
	// Data
	board          *Board
	projects       []Project // List of available projects
	selectedProject int      // Which project is selected in project list
	backend        Backend   // Backend for persistence

	// UI State
	viewMode       ViewMode
	selectedColumn int  // Which column is selected (0-4)
	selectedCard   int  // Which card in the column is selected
	showDetails    bool // Show detail panel
	showArchive    bool // Show archive column
	width          int
	height         int

	// Layout (calculated from width/height)
	boardWidth  int // Width of the board area (67% when details shown, 100% when hidden)
	detailWidth int // Width of detail panel (33% when shown, 0 when hidden)

	// Keyboard/Mouse state
	ready bool

	// Mouse drag state (Solitaire pattern)
	draggingCard   *Card     // Card currently being dragged
	dragFromColumn int       // Which column the drag started from
	dragFromIndex  int       // Card index in the source column
	mousePressX    int       // X position where mouse was pressed
	mousePressY    int       // Y position where mouse was pressed
	mouseHeldDown  bool      // Whether mouse button is currently held
	potentialDrag  bool      // Whether we're waiting to see if this becomes a drag
	dragStartTime  time.Time // When the mouse was pressed (for drag delay)

	// Drop target tracking (for visual feedback)
	dropTargetColumn int // Column where card would be dropped (-1 if none)
	dropTargetIndex  int // Position where card would be inserted

	// Double-click detection
	lastClickTime time.Time
	lastClickX    int
	lastClickY    int
}

// Project represents a discovered project with a .tkan.yaml file
type Project struct {
	Name string // Display name (from board or directory name)
	Path string // Full path to .tkan.yaml file
	Dir  string // Directory containing the project
}

// Layout represents the calculated layout dimensions
type Layout struct {
	BoardWidth  int
	DetailWidth int
	ShowDetails bool
}

// Msg types for Bubbletea
type boardLoadedMsg struct {
	board *Board
	err   error
}
