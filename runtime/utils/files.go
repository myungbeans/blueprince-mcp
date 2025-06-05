package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ListFiles lists all non-hidden files within the specified directory and its subdirectories.
// It skips directories and files/directories starting with a ".".
// Returns paths relative to the starting dirPath.
func ListFiles(dirPath string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Propagate errors encountered during traversal
			return err
		}

		// Skip hidden directories and their contents
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}

		// If it's a file and not hidden, add its path
		if !d.IsDir() && !strings.HasPrefix(d.Name(), ".") {
			paths = append(paths, path) // Fallback to full path
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory '%s': %w", dirPath, err)
	}

	return paths, nil
}

// CreateDir creates (if needed) a directory at the given path
func CreateDir(path string, perm os.FileMode) error {
	_, err := os.Stat(path)
	// dir already exists
	if err == nil {
		return nil
	}

	// Failure to chceck the dir
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check directory '%s': %w", path, err)
	}

	// Create the dir
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", path, err)
	}

	return nil
}

// IsDir checks if the given path is a directory.
func IsDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path '%s' does not exist", path)
		}
		return fmt.Errorf("failed to stat path '%s': %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path '%s' is not a directory", path)
	}
	return nil
}

// ExpandTilde expands a leading tilde in a path to the user's home directory.
func ExpandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	return strings.Replace(path, "~", homeDir, 1), nil
}

// ResolvePath expands tilde, gets the absolute path, and cleans it.
func ResolvePath(path string) (string, error) {
	expandedPath, err := ExpandTilde(path)
	if err != nil {
		return "", fmt.Errorf("failed to expand path: %w", err)
	}

	absPath, err := filepath.Abs(expandedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for '%s': %w", expandedPath, err)
	}

	return filepath.Clean(absPath), nil
}
