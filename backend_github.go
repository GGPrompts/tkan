package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// GitHubBackend implements a backend that uses GitHub Projects
type GitHubBackend struct {
	owner       string // GitHub owner (user or org)
	projectNum  int    // Project number
	repoName    string // Repository name for the project
}

// GitHubProjectItem represents an item from GitHub Projects
type GitHubProjectItem struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Content   map[string]interface{} `json:"content"`
	FieldValues map[string]interface{} `json:"fieldValues"`
}

// NewGitHubBackend creates a new GitHub Projects backend
func NewGitHubBackend(owner string, projectNum int, repoName string) *GitHubBackend {
	return &GitHubBackend{
		owner:      owner,
		projectNum: projectNum,
		repoName:   repoName,
	}
}

// GitHubProjectInfo represents a GitHub project from the list
type GitHubProjectInfo struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
	Owner  string `json:"owner"`
}

// ListGitHubProjects lists all GitHub projects accessible to the current user
func ListGitHubProjects(owner string) ([]GitHubProjectInfo, error) {
	// Use gh CLI to list projects
	var cmd *exec.Cmd
	if owner != "" {
		cmd = exec.Command("gh", "project", "list", "--owner", owner, "--format", "json", "--limit", "100")
	} else {
		cmd = exec.Command("gh", "project", "list", "--owner", "@me", "--format", "json", "--limit", "100")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %v (output: %s)", err, string(output))
	}

	// Parse JSON output
	var rawProjects struct {
		Projects []struct {
			Number int    `json:"number"`
			Title  string `json:"title"`
			Owner  struct {
				Login string `json:"login"`
			} `json:"owner"`
		} `json:"projects"`
	}

	if err := json.Unmarshal(output, &rawProjects); err != nil {
		return nil, fmt.Errorf("failed to parse project list: %v", err)
	}

	// Convert to our format
	projects := make([]GitHubProjectInfo, 0, len(rawProjects.Projects))
	for _, p := range rawProjects.Projects {
		// Use the actual owner login from the response
		actualOwner := p.Owner.Login
		if actualOwner == "" {
			actualOwner = owner // Fallback to provided owner if not in response
		}
		projects = append(projects, GitHubProjectInfo{
			Number: p.Number,
			Title:  p.Title,
			Owner:  actualOwner,
		})
	}

	return projects, nil
}

// LoadBoard fetches the project from GitHub and converts to our Board format
func (g *GitHubBackend) LoadBoard() (*Board, error) {
	// First, get project info
	projectInfo, err := g.getProjectInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get project info: %v", err)
	}

	// Get all items in the project
	items, err := g.getProjectItems()
	if err != nil {
		return nil, fmt.Errorf("failed to get project items: %v", err)
	}

	// Create board with standard columns
	boardName := "GitHub Project"
	if title, ok := projectInfo["title"].(string); ok {
		boardName = title
	}
	boardDesc := ""
	if desc, ok := projectInfo["shortDescription"].(string); ok {
		boardDesc = desc
	}

	// Construct GitHub project URL
	// Format: https://github.com/users/OWNER/projects/NUM (for users)
	//     or: https://github.com/orgs/OWNER/projects/NUM (for orgs)
	// We'll default to users, but both formats work
	boardURL := fmt.Sprintf("https://github.com/users/%s/projects/%d", g.owner, g.projectNum)

	board := &Board{
		Name:        boardName,
		Description: boardDesc,
		URL:         boardURL,
		Columns: []Column{
			{Name: "BACKLOG"},
			{Name: "TODO"},
			{Name: "PROGRESS"},
			{Name: "REVIEW"},
			{Name: "DONE"},
			{Name: "ARCHIVE"},
		},
		Cards:      []*Card{},
		CreatedAt:  time.Now(), // GitHub doesn't expose project creation time easily
		ModifiedAt: time.Now(),
	}

	// Convert GitHub items to our cards
	for _, item := range items {
		card := g.itemToCard(item)
		if card != nil {
			board.Cards = append(board.Cards, card)
		}
	}

	// Populate cards into columns
	for _, card := range board.Cards {
		for i := range board.Columns {
			if board.Columns[i].Name == card.Column {
				board.Columns[i].Cards = append(board.Columns[i].Cards, card)
				break
			}
		}
	}

	return board, nil
}

// getProjectInfo fetches basic project information
func (g *GitHubBackend) getProjectInfo() (map[string]interface{}, error) {
	query := fmt.Sprintf(`
	{
		user(login: "%s") {
			projectV2(number: %d) {
				title
				shortDescription
				id
			}
		}
	}`, g.owner, g.projectNum)

	cmd := exec.Command("gh", "api", "graphql", "-f", fmt.Sprintf("query=%s", query))
	output, err := cmd.Output()
	if err != nil {
		// Try organization if user fails
		query = fmt.Sprintf(`
		{
			organization(login: "%s") {
				projectV2(number: %d) {
					title
					shortDescription
					id
				}
			}
		}`, g.owner, g.projectNum)
		
		cmd = exec.Command("gh", "api", "graphql", "-f", fmt.Sprintf("query=%s", query))
		output, err = cmd.Output()
		if err != nil {
			return nil, err
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	// Extract project data from nested structure
	data := result["data"].(map[string]interface{})
	var project map[string]interface{}
	
	if userOrOrg, ok := data["user"]; ok && userOrOrg != nil {
		project = userOrOrg.(map[string]interface{})["projectV2"].(map[string]interface{})
	} else if userOrOrg, ok := data["organization"]; ok && userOrOrg != nil {
		project = userOrOrg.(map[string]interface{})["projectV2"].(map[string]interface{})
	}

	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	return project, nil
}

// getProjectItems fetches all items in the project
func (g *GitHubBackend) getProjectItems() ([]GitHubProjectItem, error) {
	// Use gh CLI to list items
	cmd := exec.Command("gh", "project", "item-list", 
		fmt.Sprintf("%d", g.projectNum),
		"--owner", g.owner,
		"--format", "json",
		"--limit", "100")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list project items: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	// Parse items from the response
	itemsRaw, ok := result["items"].([]interface{})
	if !ok {
		return []GitHubProjectItem{}, nil
	}

	items := []GitHubProjectItem{}
	for _, itemRaw := range itemsRaw {
		item, ok := itemRaw.(map[string]interface{})
		if !ok {
			continue
		}
		
		// Safely extract ID
		id, _ := item["id"].(string)
		if id == "" {
			continue
		}
		
		ghItem := GitHubProjectItem{
			ID:          id,
			Type:        "DraftIssue", // Default type
			Content:     map[string]interface{}{},
			FieldValues: map[string]interface{}{},
		}

		// Extract content
		if content, ok := item["content"].(map[string]interface{}); ok {
			ghItem.Content = content
			// Get the actual type from content
			if itemType, ok := content["type"].(string); ok {
				ghItem.Type = itemType
			}
		}

		// Extract field values (if present)
		if fields, ok := item["fieldValues"].(map[string]interface{}); ok {
			ghItem.FieldValues = fields
		}

		items = append(items, ghItem)
	}

	return items, nil
}

// itemToCard converts a GitHub project item to our Card format
func (g *GitHubBackend) itemToCard(item GitHubProjectItem) *Card {
	card := &Card{
		ID:         item.ID,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	// Extract title and body from content
	if title, ok := item.Content["title"].(string); ok {
		card.Title = title
	}
	if body, ok := item.Content["body"].(string); ok {
		card.Description = body
	}

	// Extract URL (for issues/PRs linked to the item)
	if url, ok := item.Content["url"].(string); ok {
		card.URL = url
	}

	// Map Status field to our Column
	if status, ok := item.FieldValues["Status"].(string); ok {
		card.Column = g.mapStatusToColumn(status)
	} else {
		card.Column = "BACKLOG" // Default column
	}

	// Extract other fields if they exist
	if assignee, ok := item.FieldValues["Assignees"].(string); ok {
		card.Assignee = assignee
	}
	if dueDate, ok := item.FieldValues["Target Date"].(string); ok {
		card.DueDate = dueDate
	}

	// Extract labels as tags
	if labels, ok := item.Content["labels"].([]interface{}); ok {
		for _, label := range labels {
			if labelMap, ok := label.(map[string]interface{}); ok {
				if name, ok := labelMap["name"].(string); ok {
					card.Tags = append(card.Tags, name)
				}
			}
		}
	}

	return card
}

// mapStatusToColumn maps GitHub Project Status to our column names
func (g *GitHubBackend) mapStatusToColumn(status string) string {
	// Common GitHub Project status mappings
	statusMap := map[string]string{
		"Backlog":     "BACKLOG",
		"Todo":        "TODO",
		"To Do":       "TODO",
		"In Progress": "PROGRESS",
		"In Review":   "REVIEW",
		"Review":      "REVIEW",
		"Done":        "DONE",
		"Closed":      "ARCHIVE",
		"Archive":     "ARCHIVE",
	}

	if column, ok := statusMap[status]; ok {
		return column
	}
	
	// Default mapping for unknown statuses
	return "BACKLOG"
}

// mapColumnToStatus maps our column names to GitHub Project Status
func (g *GitHubBackend) mapColumnToStatus(column string) string {
	statusMap := map[string]string{
		"BACKLOG": "Backlog",
		"TODO":    "Todo",
		"PROGRESS": "In Progress",
		"REVIEW":   "In Review",
		"DONE":     "Done",
		"ARCHIVE":  "Archive",
	}

	if status, ok := statusMap[column]; ok {
		return status
	}
	
	return "Backlog"
}

// SaveBoard is a no-op for GitHub backend (changes are immediate)
func (g *GitHubBackend) SaveBoard(board *Board) error {
	// GitHub changes are immediate, no need to save
	// This could be used for batch updates in the future
	return nil
}

// MoveCard moves a card to a different column in GitHub
func (g *GitHubBackend) MoveCard(cardID string, toColumn string) error {
	status := g.mapColumnToStatus(toColumn)
	
	// Build GraphQL mutation to update the Status field
	query := fmt.Sprintf(`
	mutation {
		updateProjectV2ItemFieldValue(input: {
			projectId: "%s"
			itemId: "%s"
			fieldId: "%s"
			value: {
				singleSelectOptionId: "%s"
			}
		}) {
			projectV2Item {
				id
			}
		}
	}`, g.getProjectID(), cardID, g.getStatusFieldID(), g.getStatusOptionID(status))

	cmd := exec.Command("gh", "api", "graphql", "-f", fmt.Sprintf("query=%s", query))
	_, err := cmd.Output()
	return err
}

// UpdateCard updates a card's details in GitHub
func (g *GitHubBackend) UpdateCard(card *Card) error {
	// For draft issues, we can update title and body
	query := fmt.Sprintf(`
	mutation {
		updateProjectV2DraftIssue(input: {
			draftIssueId: "%s"
			title: "%s"
			body: "%s"
		}) {
			draftIssue {
				id
			}
		}
	}`, card.ID, strings.ReplaceAll(card.Title, "\"", "\\\""), 
		strings.ReplaceAll(card.Description, "\"", "\\\""))

	cmd := exec.Command("gh", "api", "graphql", "-f", fmt.Sprintf("query=%s", query))
	_, err := cmd.Output()
	return err
}

// CreateCard creates a new draft issue in the project
func (g *GitHubBackend) CreateCard(title, description, column string) (*Card, error) {
	// Create draft issue using gh CLI
	cmd := exec.Command("gh", "project", "item-create",
		fmt.Sprintf("%d", g.projectNum),
		"--owner", g.owner,
		"--title", title,
		"--body", description,
		"--format", "json")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, err
	}

	card := &Card{
		ID:          result["id"].(string),
		Title:       title,
		Description: description,
		Column:      column,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
	}

	// Set the initial column
	if column != "BACKLOG" {
		g.MoveCard(card.ID, column)
	}

	return card, nil
}

// DeleteCard archives a card in GitHub (moves to Archive column)
func (g *GitHubBackend) DeleteCard(cardID string) error {
	return g.MoveCard(cardID, "ARCHIVE")
}

// Helper methods for field IDs (would need to be fetched from project schema)
func (g *GitHubBackend) getProjectID() string {
	// This would need to be fetched and cached
	return ""
}

func (g *GitHubBackend) getStatusFieldID() string {
	// This would need to be fetched from project fields
	return ""
}

func (g *GitHubBackend) getStatusOptionID(status string) string {
	// This would need to be fetched from field options
	return ""
}