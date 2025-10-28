package main

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadBoard loads a board from a YAML file
func LoadBoard(filename string) (*Board, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read board file: %w", err)
	}

	var board Board
	if err := yaml.Unmarshal(data, &board); err != nil {
		return nil, fmt.Errorf("failed to parse board YAML: %w", err)
	}

	// Populate column cards from board cards
	board.PopulateColumnCards()

	return &board, nil
}

// SaveBoard saves a board to a YAML file
func SaveBoard(filename string, board *Board) error {
	board.ModifiedAt = time.Now()

	data, err := yaml.Marshal(board)
	if err != nil {
		return fmt.Errorf("failed to marshal board to YAML: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write board file: %w", err)
	}

	return nil
}

// PopulateColumnCards populates column cards from board cards
func (b *Board) PopulateColumnCards() {
	// Clear existing cards in columns
	for i := range b.Columns {
		b.Columns[i].Cards = nil
	}

	// Distribute cards to columns based on Card.Column field
	for _, card := range b.Cards {
		for i := range b.Columns {
			if b.Columns[i].Name == card.Column {
				b.Columns[i].Cards = append(b.Columns[i].Cards, card)
				break
			}
		}
	}
}

// CreateDefaultBoard creates a default board with sample data
func CreateDefaultBoard() *Board {
	now := time.Now()

	board := &Board{
		Name:        "My Project",
		Description: "A sample Kanban board",
		Columns: []Column{
			{Name: "BACKLOG"},
			{Name: "TODO"},
			{Name: "PROGRESS"},
			{Name: "REVIEW"},
			{Name: "DONE"},
			{Name: "ARCHIVE"},
		},
		Cards: []*Card{
			{
				ID:          "1",
				Title:       "Fix login flow",
				Description: "Users can't authenticate via OAuth. Error 401 on token refresh.",
				Tags:        []string{"bug", "p1"},
				Assignee:    "@alice",
				DueDate:     "2025-01-15",
				CreatedAt:   now.AddDate(0, 0, -10),
				ModifiedAt:  now.AddDate(0, 0, -1),
				Column:      "TODO",
			},
			{
				ID:          "0",
				Title:       "New feature idea",
				Description: "Consider adding dark mode support.",
				Tags:        []string{"enhancement"},
				Assignee:    "",
				DueDate:     "",
				CreatedAt:   now.AddDate(0, 0, -20),
				ModifiedAt:  now.AddDate(0, 0, -15),
				Column:      "BACKLOG",
			},
			{
				ID:          "2",
				Title:       "Add OAuth support",
				Description: "Implement OAuth 2.0 authentication flow with Google and GitHub providers.",
				Tags:        []string{"feature"},
				Assignee:    "@bob",
				DueDate:     "2025-01-20",
				CreatedAt:   now.AddDate(0, 0, -8),
				ModifiedAt:  now.AddDate(0, 0, -1),
				Column:      "PROGRESS",
			},
			{
				ID:          "3",
				Title:       "Review PR #42",
				Description: "Code review for authentication refactor pull request.",
				Tags:        []string{"code-review"},
				Assignee:    "@charlie",
				DueDate:     "2025-01-18",
				CreatedAt:   now.AddDate(0, 0, -5),
				ModifiedAt:  now.AddDate(0, 0, -1),
				Column:      "REVIEW",
			},
			{
				ID:          "4",
				Title:       "Setup database",
				Description: "Configure PostgreSQL database and run migrations.",
				Tags:        []string{"infra", "done"},
				Assignee:    "@dave",
				DueDate:     "",
				CreatedAt:   now.AddDate(0, 0, -15),
				ModifiedAt:  now.AddDate(0, 0, -10),
				Column:      "DONE",
			},
			{
				ID:          "5",
				Title:       "Write tests",
				Description: "Add unit tests for authentication module.",
				Tags:        []string{"test"},
				Assignee:    "@alice",
				DueDate:     "2025-01-22",
				CreatedAt:   now.AddDate(0, 0, -9),
				ModifiedAt:  now.AddDate(0, 0, -2),
				Column:      "TODO",
			},
		},
		CreatedAt:  now.AddDate(0, 0, -20),
		ModifiedAt: now,
	}

	board.PopulateColumnCards()
	return board
}
