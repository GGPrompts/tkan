package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Scan for projects (current directory and subdirectories)
	projects, err := ScanProjects(cwd)
	if err != nil {
		fmt.Printf("Error scanning for projects: %v\n", err)
		os.Exit(1)
	}

	// If no projects found, create a default one in current directory
	if len(projects) == 0 {
		fmt.Println("No projects found. Creating default board...")
		board := CreateDefaultBoard()
		if err := SaveBoard(".tkan.yaml", board); err != nil {
			fmt.Printf("Error saving board: %v\n", err)
			os.Exit(1)
		}
		projects = append(projects, Project{
			Name: board.Name,
			Path: ".tkan.yaml",
			Dir:  cwd,
		})
	}

	// Load the first project by default
	board, err := LoadBoard(projects[0].Path)
	if err != nil {
		fmt.Printf("Error loading board: %v\n", err)
		os.Exit(1)
	}

	// Initialize model
	m := NewModel(board, projects)

	// Create Bubbletea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
