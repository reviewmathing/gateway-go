package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GetRootDir returns the root directory of the application.
// It checks APP_ROOT_DIR environment variable first for explicit control.
// In development environments (go run, IDE), it returns the working directory.
// In production, it returns the directory containing the executable.
func GetRootDir() (string, error) {
	// Allow explicit override via environment variable
	if rootDir := os.Getenv("APP_ROOT_DIR"); rootDir != "" {
		return rootDir, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)

	// Check if running in development environment
	if isDevEnvironment(exePath) {
		return os.Getwd()
	}

	// Production environment - use executable directory
	return exeDir, nil
}

// isDevEnvironment checks if the executable path indicates a development environment.
// Returns true for temporary build directories used by go run, IDEs, and debuggers.
func isDevEnvironment(path string) bool {
	tempPatterns := []string{
		"go-build",         // go run
		"GoLand",           // GoLand
		"Caches/JetBrains", // JetBrains IDEs
		"__debug_bin",      // debuggers
		"/tmp/",            // Linux tmp
		"\\Temp\\",         // Windows tmp
	}

	for _, pattern := range tempPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}
