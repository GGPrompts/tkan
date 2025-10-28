package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanProjects scans for .tkan.yaml files starting from the given directory
// It searches the current directory and all subdirectories (up to 3 levels deep)
func ScanProjects(startDir string) ([]Project, error) {
	var projects []Project

	// Check if startDir itself has a .tkan.yaml
	if project, found := checkProjectInDir(startDir); found {
		projects = append(projects, project)
	}

	// Walk subdirectories (max 3 levels deep)
	err := filepath.Walk(startDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip directories we can't access
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Skip if too deep (more than 3 levels)
		relPath, _ := filepath.Rel(startDir, path)
		depth := strings.Count(relPath, string(os.PathSeparator))
		if depth > 3 {
			return filepath.SkipDir
		}

		// Check if this directory has a .tkan.yaml
		if info.IsDir() && path != startDir {
			if project, found := checkProjectInDir(path); found {
				projects = append(projects, project)
				return filepath.SkipDir // Don't scan subdirectories of projects
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return projects, nil
}

// checkProjectInDir checks if a directory contains a .tkan.yaml file
func checkProjectInDir(dir string) (Project, bool) {
	yamlPath := filepath.Join(dir, ".tkan.yaml")
	if _, err := os.Stat(yamlPath); err == nil {
		// Load board to get the name
		board, err := LoadBoard(yamlPath)
		name := filepath.Base(dir) // Default to directory name
		if err == nil && board.Name != "" {
			name = board.Name
		}

		return Project{
			Name: name,
			Path: yamlPath,
			Dir:  dir,
		}, true
	}
	return Project{}, false
}

// GetProjectRelativePath returns a relative path for display
func GetProjectRelativePath(project Project, baseDir string) string {
	relPath, err := filepath.Rel(baseDir, project.Dir)
	if err != nil {
		return project.Dir
	}
	if relPath == "." {
		return "(current directory)"
	}
	return relPath
}
