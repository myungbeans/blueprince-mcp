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

		// Skip unwanted paths (hidden files and paths)
		if ShouldSkipPath(path, blob) {
			if blob.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if blob.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			// Skip if relative path fails
			return nil
		}
		paths = append(paths, relPath)
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

// ValidateNotePath validates and cleans a note path for security
// Returns the cleaned path and any validation error
func ValidateNotePath(notePath string) (string, error) {
	if notePath == "" {
		return "", fmt.Errorf("note path cannot be empty")
	}

	// Clean the path to prevent path traversal
	cleanPath := filepath.Clean(notePath)

	// Check for path traversal attempts or absolute paths
	if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return "", fmt.Errorf("invalid note path: '%s'. Must be a relative path within the vault", notePath)
	}

	return cleanPath, nil
}

// BuildSecureNotePath constructs and validates a full path to a note file within the vault
// Returns the full path and any security validation error
func BuildSecureNotePath(vaultPath, cleanNotePath string) (string, error) {
	// Construct the full path to the note file
	fullPath := filepath.Join(vaultPath, "notes", cleanNotePath)

	// Security check: Ensure the resolved path is still within the notes directory
	absNotesPath, err := filepath.Abs(filepath.Join(vaultPath, "notes"))
	if err != nil {
		return "", fmt.Errorf("failed to resolve notes directory path: %w", err)
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve full file path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, absNotesPath) {
		return "", fmt.Errorf("access denied: path '%s' is outside the notes directory", cleanNotePath)
	}

	return fullPath, nil
}

// ExtractStringParam extracts and validates a string parameter from MCP arguments
func ExtractStringParam(params map[string]any, paramName string) (string, error) {
	if params == nil {
		return "", fmt.Errorf("missing arguments")
	}

	paramRaw, ok := params[paramName]
	if !ok {
		return "", fmt.Errorf("missing required parameter: %s", paramName)
	}

	paramValue, ok := paramRaw.(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' must be a string", paramName)
	}

	return paramValue, nil
}

// ShouldSkipPath determines if a path should be skipped during file traversal
func ShouldSkipPath(path string, info os.DirEntry) bool {
	name := info.Name()

	// Skip all hidden files and directories (starting with .)
	if strings.HasPrefix(name, ".") {
		return true
	}

	return false
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
