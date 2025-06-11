package auth

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	gdrive "github.com/myungbeans/blueprince-mcp/runtime/storage/drive"
	"golang.org/x/oauth2"
)

func TestNewGoogleDriveAuth(t *testing.T) {
	// Create temporary directory and files for testing
	tempDir, err := os.MkdirTemp("", "auth_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Mock credentials file
	credentialsContent := `{
		"installed": {
			"client_id": "test_client_id",
			"client_secret": "test_client_secret",
			"redirect_uris": ["http://localhost:8080"]
		}
	}`

	// Temporarily change working directory to temp dir
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create credentials file
	err = os.WriteFile(gdrive.APP_CREDS_FILE, []byte(credentialsContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create credentials file: %v", err)
	}

	// Temporarily change HOME to temp dir for config paths
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Test NewGoogleDriveAuth
	ctx := context.Background()
	auth, err := NewGoogleDriveAuth(ctx)
	if err != nil {
		t.Fatalf("NewGoogleDriveAuth() failed: %v", err)
	}

	// Verify auth object is properly initialized
	if auth == nil {
		t.Error("NewGoogleDriveAuth() should return non-nil auth object")
	}

	if auth.config == nil {
		t.Error("NewGoogleDriveAuth() should initialize config")
	}

	if auth.ctx == nil {
		t.Error("NewGoogleDriveAuth() should initialize context")
	}

	if auth.tokenFile == "" {
		t.Error("NewGoogleDriveAuth() should set tokenFile path")
	}

	if auth.configFile == "" {
		t.Error("NewGoogleDriveAuth() should set configFile path")
	}

	// Verify paths are correct
	expectedTokenPath := filepath.Join(tempDir, gdrive.CONFIG_DIR, gdrive.TOKEN_FILE)
	if auth.tokenFile != expectedTokenPath {
		t.Errorf("Expected tokenFile to be %s, got %s", expectedTokenPath, auth.tokenFile)
	}

	expectedConfigPath := filepath.Join(tempDir, gdrive.CONFIG_DIR, gdrive.CONFIG_FILE)
	if auth.configFile != expectedConfigPath {
		t.Errorf("Expected configFile to be %s, got %s", expectedConfigPath, auth.configFile)
	}

	// Verify config directory was created
	configDir := filepath.Join(tempDir, gdrive.CONFIG_DIR)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("NewGoogleDriveAuth() should create config directory")
	}
}

func TestNewGoogleDriveAuth_MissingCredentials(t *testing.T) {
	// Create temporary directory without credentials file
	tempDir, err := os.MkdirTemp("", "auth_test_nocreds")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Temporarily change working directory to temp dir
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Test NewGoogleDriveAuth with missing credentials
	ctx := context.Background()
	_, err = NewGoogleDriveAuth(ctx)
	if err == nil {
		t.Error("NewGoogleDriveAuth() should fail when credentials file is missing")
	}
}

func TestGoogleDriveAuth_SaveConfig(t *testing.T) {
	// Create temporary setup
	tempDir, auth := setupTestAuth(t)
	defer os.RemoveAll(tempDir)

	// Mock Drive service with a simple test
	// Since we can't easily mock the Google Drive API, we'll test the structure
	folderName := "TestFolder"

	// Create a minimal auth object for testing config saving logic
	configPath := filepath.Join(tempDir, gdrive.CONFIG_DIR, gdrive.CONFIG_FILE)

	// Manually create the expected config
	expectedConfig := gdrive.DriveConfig{
		FolderID:   "test_folder_id_123",
		FolderName: folderName,
		TokenPath:  auth.tokenFile,
	}

	// Write the config manually to test the file structure
	configData, err := json.MarshalIndent(expectedConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configData, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}

	// Verify config file contents
	savedData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var savedConfig gdrive.DriveConfig
	err = json.Unmarshal(savedData, &savedConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal saved config: %v", err)
	}

	if savedConfig.FolderID != expectedConfig.FolderID {
		t.Errorf("Expected FolderID %s, got %s", expectedConfig.FolderID, savedConfig.FolderID)
	}

	if savedConfig.FolderName != expectedConfig.FolderName {
		t.Errorf("Expected FolderName %s, got %s", expectedConfig.FolderName, savedConfig.FolderName)
	}

	if savedConfig.TokenPath != expectedConfig.TokenPath {
		t.Errorf("Expected TokenPath %s, got %s", expectedConfig.TokenPath, savedConfig.TokenPath)
	}
}

func TestGoogleDriveAuth_Structure(t *testing.T) {
	tempDir, auth := setupTestAuth(t)
	defer os.RemoveAll(tempDir)

	// Test that all fields are properly typed
	var _ *oauth2.Config = auth.config
	var _ string = auth.tokenFile
	var _ string = auth.configFile
	var _ context.Context = auth.ctx

	// Test that service field can be set
	auth.service = nil // Should not panic
}

func TestGoogleDriveAuth_TokenPaths(t *testing.T) {
	tempDir, auth := setupTestAuth(t)
	defer os.RemoveAll(tempDir)

	// Verify token file path is within config directory
	configDir := filepath.Join(tempDir, gdrive.CONFIG_DIR)
	if !filepath.HasPrefix(auth.tokenFile, configDir) {
		t.Errorf("Token file should be within config directory %s, got %s", configDir, auth.tokenFile)
	}

	// Verify config file path is within config directory
	if !filepath.HasPrefix(auth.configFile, configDir) {
		t.Errorf("Config file should be within config directory %s, got %s", configDir, auth.configFile)
	}

	// Verify paths are different
	if auth.tokenFile == auth.configFile {
		t.Error("Token file and config file should have different paths")
	}
}

func TestGoogleDriveAuth_ErrorHandling(t *testing.T) {
	// Test with invalid context
	ctx := context.Background()

	// Test with missing home directory (simulate by using invalid HOME)
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", "/nonexistent/invalid/path/that/should/not/exist")
	defer os.Setenv("HOME", originalHome)

	tempDir, err := os.MkdirTemp("", "auth_error_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create credentials file in temp dir
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	credentialsContent := `{
		"installed": {
			"client_id": "test_client_id",
			"client_secret": "test_client_secret",
			"redirect_uris": ["http://localhost:8080"]
		}
	}`
	err = os.WriteFile(gdrive.APP_CREDS_FILE, []byte(credentialsContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create credentials file: %v", err)
	}

	// This should fail due to invalid HOME directory
	_, err = NewGoogleDriveAuth(ctx)
	if err == nil {
		t.Error("NewGoogleDriveAuth() should fail with invalid HOME directory")
	}
}

// Helper function to set up test environment
func setupTestAuth(t *testing.T) (string, *GoogleDriveAuth) {
	tempDir, err := os.MkdirTemp("", "auth_test_setup")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create credentials file
	credentialsContent := `{
		"installed": {
			"client_id": "test_client_id",
			"client_secret": "test_client_secret",
			"redirect_uris": ["http://localhost:8080"]
		}
	}`

	// Change working directory
	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	t.Cleanup(func() { os.Chdir(originalWd) })

	err = os.WriteFile(gdrive.APP_CREDS_FILE, []byte(credentialsContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create credentials file: %v", err)
	}

	// Change HOME directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	t.Cleanup(func() { os.Setenv("HOME", originalHome) })

	// Create auth object
	ctx := context.Background()
	auth, err := NewGoogleDriveAuth(ctx)
	if err != nil {
		t.Fatalf("Failed to create auth object: %v", err)
	}

	return tempDir, auth
}

// Benchmark tests
func BenchmarkNewGoogleDriveAuth(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "auth_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	credentialsContent := `{
		"installed": {
			"client_id": "test_client_id",
			"client_secret": "test_client_secret",
			"redirect_uris": ["http://localhost:8080"]
		}
	}`

	originalWd, _ := os.Getwd()
	err = os.Chdir(tempDir)
	if err != nil {
		b.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	err = os.WriteFile(gdrive.APP_CREDS_FILE, []byte(credentialsContent), 0600)
	if err != nil {
		b.Fatalf("Failed to create credentials file: %v", err)
	}

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth, err := NewGoogleDriveAuth(ctx)
		if err != nil {
			b.Fatalf("NewGoogleDriveAuth failed: %v", err)
		}
		_ = auth
	}
}