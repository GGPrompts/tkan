package main

import (
	"fmt"
	"time"
)

// Backend interface defines methods for board persistence
type Backend interface {
	LoadBoard() (*Board, error)
	SaveBoard(*Board) error
	MoveCard(cardID string, toColumn string) error
	UpdateCard(card *Card) error
	CreateCard(title, description, column string) (*Card, error)
	DeleteCard(cardID string) error
}

// LocalBackend implements Backend using local YAML files
type LocalBackend struct {
	filePath string
}

// NewLocalBackend creates a new local file backend
func NewLocalBackend(filePath string) *LocalBackend {
	return &LocalBackend{filePath: filePath}
}

// LoadBoard loads from YAML file
func (l *LocalBackend) LoadBoard() (*Board, error) {
	return LoadBoard(l.filePath)
}

// SaveBoard saves to YAML file
func (l *LocalBackend) SaveBoard(board *Board) error {
	return SaveBoard(l.filePath, board)
}

// MoveCard moves a card to a different column
func (l *LocalBackend) MoveCard(cardID string, toColumn string) error {
	board, err := l.LoadBoard()
	if err != nil {
		return err
	}

	// Find and update the card
	for _, card := range board.Cards {
		if card.ID == cardID {
			card.Column = toColumn
			break
		}
	}

	return l.SaveBoard(board)
}

// UpdateCard updates a card's details
func (l *LocalBackend) UpdateCard(card *Card) error {
	board, err := l.LoadBoard()
	if err != nil {
		return err
	}

	// Find and update the card
	for i, c := range board.Cards {
		if c.ID == card.ID {
			board.Cards[i] = card
			break
		}
	}

	return l.SaveBoard(board)
}

// CreateCard creates a new card
func (l *LocalBackend) CreateCard(title, description, column string) (*Card, error) {
	board, err := l.LoadBoard()
	if err != nil {
		return nil, err
	}

	// Generate new ID
	maxID := 0
	for _, card := range board.Cards {
		var id int
		fmt.Sscanf(card.ID, "%d", &id)
		if id > maxID {
			maxID = id
		}
	}

	newCard := &Card{
		ID:          fmt.Sprintf("%d", maxID+1),
		Title:       title,
		Description: description,
		Column:      column,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	board.Cards = append(board.Cards, newCard)
	
	if err := l.SaveBoard(board); err != nil {
		return nil, err
	}

	return newCard, nil
}

// DeleteCard removes a card (moves to archive)
func (l *LocalBackend) DeleteCard(cardID string) error {
	return l.MoveCard(cardID, "ARCHIVE")
}