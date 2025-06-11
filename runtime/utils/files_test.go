package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFiles(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "listfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files and directories
	testFiles := []string{
		"file1.txt",
		"file2.md",
		"subdir/file3.txt",
		"subdir/file4.md",
		"subdir/nested/file5.txt",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file, err)
		}
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Create hidden files/directories (should be skipped)
	hiddenFiles := []string{
		".hidden_file",
		".hidden_dir/file.txt",
	}

	for _, file := range hiddenFiles {
		fullPath := filepath.Join(tempDir, file)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file, err)
		}
		err = os.WriteFile(fullPath, []byte("hidden content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create hidden file %s: %v", file, err)
		}
	}

	// Test ListFiles
	files, err := ListFiles(tempDir)
	if err != nil {
		t.Fatalf("ListFiles() failed: %v", err)
	}

	// Verify correct number of files (should exclude hidden files)
	expectedCount := len(testFiles)
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d: %v", expectedCount, len(files), files)
	}

	// Verify all expected files are present
	for _, expectedFile := range testFiles {
		found := false
		for _, actualFile := range files {
			if actualFile == expectedFile {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected file %s not found in results: %v", expectedFile, files)
		}
	}

	// Verify hidden files are not included
	for _, hiddenFile := range hiddenFiles {
		for _, actualFile := range files {
			if actualFile == hiddenFile {
				t.Errorf("Hidden file %s should not be included in results", hiddenFile)
			}
		}
	}
}

func TestValidateDir(t *testing.T) {
	// Test with valid directory
	tempDir, err := os.MkdirTemp("", "validatedir_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = ValidateDir(tempDir)
	if err != nil {
		t.Errorf("ValidateDir() should succeed for valid directory, got: %v", err)
	}

	// Test with non-existent path
	nonExistentPath := filepath.Join(tempDir, "nonexistent")
	err = ValidateDir(nonExistentPath)
	if err == nil {
		t.Error("ValidateDir() should fail for non-existent path")
	}

	// Test with file instead of directory
	testFile := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = ValidateDir(testFile)
	if err == nil {
		t.Error("ValidateDir() should fail for file instead of directory")
	}
}

func TestExpandTilde(t *testing.T) {
	// Test path without tilde
	normalPath := "/usr/local/bin"
	result, err := ExpandTilde(normalPath)
	if err != nil {
		t.Errorf("ExpandTilde() failed for normal path: %v", err)
	}
	if result != normalPath {
		t.Errorf("ExpandTilde() should return unchanged path, got: %s", result)
	}

	// Test path with tilde
	tilePath := "~/Documents"
	result, err = ExpandTilde(tilePath)
	if err != nil {
		t.Errorf("ExpandTilde() failed for tilde path: %v", err)
	}
	if strings.HasPrefix(result, "~") {
		t.Errorf("ExpandTilde() should expand tilde, got: %s", result)
	}
	if !strings.HasSuffix(result, "Documents") {
		t.Errorf("ExpandTilde() should preserve path after tilde, got: %s", result)
	}

	// Test just tilde
	justTilde := "~"
	result, err = ExpandTilde(justTilde)
	if err != nil {
		t.Errorf("ExpandTilde() failed for just tilde: %v", err)
	}
	if result == "~" {
		t.Errorf("ExpandTilde() should expand just tilde, got: %s", result)
	}
}

func TestResolveAndCleanPath(t *testing.T) {
	// Test with normal path
	normalPath := "/usr/local/bin"
	result, err := ResolveAndCleanPath(normalPath)
	if err != nil {
		t.Errorf("ResolveAndCleanPath() failed for normal path: %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ResolveAndCleanPath() should return absolute path, got: %s", result)
	}

	// Test with path containing tilde
	tilePath := "~/Documents"
	result, err = ResolveAndCleanPath(tilePath)
	if err != nil {
		t.Errorf("ResolveAndCleanPath() failed for tilde path: %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ResolveAndCleanPath() should return absolute path, got: %s", result)
	}
	if strings.Contains(result, "~") {
		t.Errorf("ResolveAndCleanPath() should not contain tilde, got: %s", result)
	}

	// Test with relative path
	relativePath := "./test"
	result, err = ResolveAndCleanPath(relativePath)
	if err != nil {
		t.Errorf("ResolveAndCleanPath() failed for relative path: %v", err)
	}
	if !filepath.IsAbs(result) {
		t.Errorf("ResolveAndCleanPath() should return absolute path, got: %s", result)
	}
}

func TestValidateNotePath(t *testing.T) {
	// Test valid relative path
	validPath := "notes/test.md"
	result, err := ValidateNotePath(validPath)
	if err != nil {
		t.Errorf("ValidateNotePath() should succeed for valid path, got: %v", err)
	}
	if result != validPath {
		t.Errorf("ValidateNotePath() should return cleaned path, got: %s", result)
	}

	// Test empty path
	_, err = ValidateNotePath("")
	if err == nil {
		t.Error("ValidateNotePath() should fail for empty path")
	}

	// Test path traversal attempt
	traversalPath := "../../../etc/passwd"
	_, err = ValidateNotePath(traversalPath)
	if err == nil {
		t.Error("ValidateNotePath() should fail for path traversal attempt")
	}

	// Test absolute path
	absolutePath := "/etc/passwd"
	_, err = ValidateNotePath(absolutePath)
	if err == nil {
		t.Error("ValidateNotePath() should fail for absolute path")
	}

	// Test path with dot-dot in middle
	dotDotPath := "notes/../config/file.txt"
	result, err = ValidateNotePath(dotDotPath)
	if err != nil {
		t.Errorf("ValidateNotePath() should handle dot-dot in middle: %v", err)
	}
	if strings.Contains(result, "..") {
		t.Errorf("ValidateNotePath() should clean dot-dot, got: %s", result)
	}
}

func TestBuildSecureNotePath(t *testing.T) {
	// Create temporary vault directory
	tempVault, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatalf("Failed to create temp vault: %v", err)
	}
	defer os.RemoveAll(tempVault)

	notesDir := filepath.Join(tempVault, "notes")
	err = os.MkdirAll(notesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create notes directory: %v", err)
	}

	// Test valid note path
	validNotePath := "category/test.md"
	result, err := BuildSecureNotePath(tempVault, validNotePath)
	if err != nil {
		t.Errorf("BuildSecureNotePath() should succeed for valid path, got: %v", err)
	}

	expectedPath := filepath.Join(tempVault, "notes", validNotePath)
	if result != expectedPath {
		t.Errorf("BuildSecureNotePath() should return correct path, got: %s, expected: %s", result, expectedPath)
	}

	// Verify the path is within notes directory
	if !strings.HasPrefix(result, notesDir) {
		t.Errorf("BuildSecureNotePath() should return path within notes directory")
	}
}

func TestExtractStringParam(t *testing.T) {
	// Test valid parameter
	params := map[string]any{
		"test_param": "test_value",
		"other_param": 123,
	}

	result, err := ExtractStringParam(params, "test_param")
	if err != nil {
		t.Errorf("ExtractStringParam() should succeed for valid string param, got: %v", err)
	}
	if result != "test_value" {
		t.Errorf("ExtractStringParam() should return correct value, got: %s", result)
	}

	// Test missing parameter
	_, err = ExtractStringParam(params, "missing_param")
	if err == nil {
		t.Error("ExtractStringParam() should fail for missing parameter")
	}

	// Test non-string parameter
	_, err = ExtractStringParam(params, "other_param")
	if err == nil {
		t.Error("ExtractStringParam() should fail for non-string parameter")
	}

	// Test nil params
	_, err = ExtractStringParam(nil, "test_param")
	if err == nil {
		t.Error("ExtractStringParam() should fail for nil params")
	}
}

func TestShouldSkipPath(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected bool
	}{
		{"regular file", "test.txt", false},
		{"regular directory", "testdir", false},
		{"hidden file", ".hidden", true},
		{"hidden directory", ".hiddendir", true},
		{"dot file", ".gitignore", true},
		{"double dot", "..", true},
		{"single dot", ".", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock DirEntry
			info := &mockDirEntry{name: tc.filename}
			result := ShouldSkipPath("/test/path/"+tc.filename, info)
			if result != tc.expected {
				t.Errorf("ShouldSkipPath(%s) = %v, expected %v", tc.filename, result, tc.expected)
			}
		})
	}
}

func TestEnsureDirExists(t *testing.T) {
	// Test creating new directory
	tempParent, err := os.MkdirTemp("", "ensuredir_test")
	if err != nil {
		t.Fatalf("Failed to create temp parent dir: %v", err)
	}
	defer os.RemoveAll(tempParent)

	newDir := filepath.Join(tempParent, "newdir")
	err = EnsureDirExists(newDir, 0755)
	if err != nil {
		t.Errorf("EnsureDirExists() should succeed for new directory, got: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("EnsureDirExists() should create directory")
	}

	// Test with existing directory (should succeed)
	err = EnsureDirExists(newDir, 0755)
	if err != nil {
		t.Errorf("EnsureDirExists() should succeed for existing directory, got: %v", err)
	}

	// Test with file instead of directory (should fail)
	testFile := filepath.Join(tempParent, "testfile.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	err = EnsureDirExists(testFile, 0755)
	if err == nil {
		t.Error("EnsureDirExists() should fail when path exists as file")
	}
}

// Mock DirEntry for testing
type mockDirEntry struct {
	name string
}

func (m *mockDirEntry) Name() string               { return m.name }
func (m *mockDirEntry) IsDir() bool                { return false }
func (m *mockDirEntry) Type() os.FileMode          { return 0 }
func (m *mockDirEntry) Info() (os.FileInfo, error) { return nil, nil }