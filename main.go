package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Parse command-line flags
	var (
		githubProject = flag.String("github", "", "Use GitHub Project (format: owner/project-number or owner/repo/project-number)")
		help          = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		fmt.Println("tkan - Terminal Kanban Board")
		fmt.Println("\nUsage:")
		fmt.Println("  tkan                    # Use local .tkan.yaml files")
		fmt.Println("  tkan --github owner/1   # Use GitHub Project #1 from owner")
		fmt.Println("  tkan --github owner/repo/1  # Use GitHub Project #1 from owner/repo")
		fmt.Println("\nExamples:")
		fmt.Println("  tkan --github matt/1")
		fmt.Println("  tkan --github microsoft/vscode/2")
		os.Exit(0)
	}

	var backend Backend
	var board *Board
	var projects []Project

	if *githubProject != "" {
		// Parse GitHub project specification
		parts := strings.Split(*githubProject, "/")
		if len(parts) < 2 || len(parts) > 3 {
			fmt.Println("Invalid GitHub project format. Use: owner/project-number or owner/repo/project-number")
			os.Exit(1)
		}

		owner := parts[0]
		var repoName string
		var projectNumStr string

		if len(parts) == 2 {
			// Format: owner/project-number
			projectNumStr = parts[1]
		} else {
			// Format: owner/repo/project-number
			repoName = parts[1]
			projectNumStr = parts[2]
		}

		projectNum, err := strconv.Atoi(projectNumStr)
		if err != nil {
			fmt.Printf("Invalid project number: %s\n", projectNumStr)
			os.Exit(1)
		}

		// Create GitHub backend
		backend = NewGitHubBackend(owner, projectNum, repoName)
		
		// Load board from GitHub
		board, err = backend.LoadBoard()
		if err != nil {
			fmt.Printf("Error loading GitHub project: %v\n", err)
			fmt.Printf("\nMake sure you have:\n")
			fmt.Printf("1. Authenticated with: gh auth login\n")
			fmt.Printf("2. Added project scope: gh auth refresh -s project\n")
			fmt.Printf("3. Access to the project: %s\n", *githubProject)
			os.Exit(1)
		}

		// Create a fake project entry for the UI
		projects = []Project{{
			Name: fmt.Sprintf("GitHub: %s", board.Name),
			Path: fmt.Sprintf("github:%s", *githubProject),
			Dir:  "GitHub",
		}}
	} else {
		// Use local backend
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current directory: %v\n", err)
			os.Exit(1)
		}

		// Scan for projects (current directory and subdirectories)
		projects, err = ScanProjects(cwd)
		if err != nil {
			fmt.Printf("Error scanning for projects: %v\n", err)
			os.Exit(1)
		}

		// If no projects found, create a default one in current directory
		if len(projects) == 0 {
			fmt.Println("No projects found. Creating default board...")
			board = CreateDefaultBoard()
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

		// Create local backend
		backend = NewLocalBackend(projects[0].Path)

		// Load the first project by default
		board, err = backend.LoadBoard()
		if err != nil {
			fmt.Printf("Error loading board: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize model with backend
	m := NewModelWithBackend(board, projects, backend)

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
