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
func ListFiles(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, blob fs.DirEntry, err error) error {
		if err != nil {
			// Propagate errors encountered during traversal
			return err
		}

		// Skip hidden directories and their contents
		if blob.IsDir() && strings.HasPrefix(blob.Name(), ".") {
			return filepath.SkipDir
		}

		// If it's a file and not hidden, add its path
		if !blob.IsDir() && !strings.HasPrefix(blob.Name(), ".") {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				// Skip if relative path fails
				return nil
			}
			paths = append(paths, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory '%s': %w", root, err)
	}
	return paths, nil
}

// ValidateDir checks if the given path is a directory.
func ValidateDir(path string) error {
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

// ResolveAndCleanPath expands tilde, gets the absolute path, and cleans it.
func ResolveAndCleanPath(path string) (string, error) {
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

// EnsureDirExists checks if a directory exists at the given path, and creates it if it doesn't.
// It uses the provided file mode for creation.
func EnsureDirExists(path string, perm os.FileMode) error {
	_, err := os.Stat(path)
	// create dir if does not exist
	if os.IsNotExist(err) {
		// rerport errors during creation
		if err := os.MkdirAll(path, perm); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", path, err)
		}
		return nil
	}

	// other error types
	if err != nil {
		return fmt.Errorf("failed to check directory '%s': %w", path, err)
	}
	return nil
}
