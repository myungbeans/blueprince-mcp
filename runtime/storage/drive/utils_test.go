package drive

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestCredsPath(t *testing.T) {
	path, err := CredsPath()
	if err != nil {
		t.Fatalf("CredsPath() failed: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("CredsPath() should return absolute path, got: %s", path)
	}

	if !strings.HasSuffix(path, APP_CREDS_FILE) {
		t.Errorf("CredsPath() should end with %s, got: %s", APP_CREDS_FILE, path)
	}
}

func TestTokenPath(t *testing.T) {
	path, err := TokenPath()
	if err != nil {
		t.Fatalf("TokenPath() failed: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("TokenPath() should return absolute path, got: %s", path)
	}

	if !strings.HasSuffix(path, TOKEN_FILE) {
		t.Errorf("TokenPath() should end with %s, got: %s", TOKEN_FILE, path)
	}

	expectedDir := filepath.Join(CONFIG_DIR, TOKEN_FILE)
	if !strings.HasSuffix(path, expectedDir) {
		t.Errorf("TokenPath() should contain %s, got: %s", expectedDir, path)
	}
}

func TestConfigPath(t *testing.T) {
	path, err := ConfigPath()
	if err != nil {
		t.Fatalf("ConfigPath() failed: %v", err)
	}

	if !filepath.IsAbs(path) {
		t.Errorf("ConfigPath() should return absolute path, got: %s", path)
	}

	if !strings.HasSuffix(path, CONFIG_FILE) {
		t.Errorf("ConfigPath() should end with %s, got: %s", CONFIG_FILE, path)
	}

	expectedDir := filepath.Join(CONFIG_DIR, CONFIG_FILE)
	if !strings.HasSuffix(path, expectedDir) {
		t.Errorf("ConfigPath() should contain %s, got: %s", expectedDir, path)
	}
}

func TestEnsureConfigDir(t *testing.T) {
	// Get a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "drive_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Temporarily change home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	err = EnsureConfigDir()
	if err != nil {
		t.Fatalf("EnsureConfigDir() failed: %v", err)
	}

	// Check if directory was created
	configDir := filepath.Join(tempDir, CONFIG_DIR)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("EnsureConfigDir() should create directory %s", configDir)
	}
}

func TestSaveAndLoadToken(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "token_test_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// Test token
	originalToken := &oauth2.Token{
		AccessToken:  "test_access_token",
		RefreshToken: "test_refresh_token",
		TokenType:    "Bearer",
	}

	// Test SaveToken
	err = SaveToken(tempFile.Name(), originalToken)
	if err != nil {
		t.Fatalf("SaveToken() failed: %v", err)
	}

	// Verify file was written
	if _, err := os.Stat(tempFile.Name()); os.IsNotExist(err) {
		t.Errorf("SaveToken() should create file %s", tempFile.Name())
	}

	// Test LoadToken by reading the file directly
	tokenData, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read token file: %v", err)
	}

	var loadedToken oauth2.Token
	err = json.Unmarshal(tokenData, &loadedToken)
	if err != nil {
		t.Fatalf("Failed to unmarshal token: %v", err)
	}

	// Compare tokens
	if loadedToken.AccessToken != originalToken.AccessToken {
		t.Errorf("AccessToken mismatch: got %s, want %s", loadedToken.AccessToken, originalToken.AccessToken)
	}
	if loadedToken.RefreshToken != originalToken.RefreshToken {
		t.Errorf("RefreshToken mismatch: got %s, want %s", loadedToken.RefreshToken, originalToken.RefreshToken)
	}
	if loadedToken.TokenType != originalToken.TokenType {
		t.Errorf("TokenType mismatch: got %s, want %s", loadedToken.TokenType, originalToken.TokenType)
	}
}

func TestLoadDriveConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "drive_config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Temporarily change home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Create config directory
	configDir := filepath.Join(tempDir, CONFIG_DIR)
	err = os.MkdirAll(configDir, 0700)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Test config
	originalConfig := &DriveConfig{
		FolderID:   "test_folder_id",
		FolderName: "test_folder",
		TokenPath:  "/path/to/token",
	}

	// Write config file
	configPath := filepath.Join(configDir, CONFIG_FILE)
	configData, err := json.MarshalIndent(originalConfig, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, configData, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Test loadDriveConfig
	loadedConfig, err := loadDriveConfig()
	if err != nil {
		t.Fatalf("loadDriveConfig() failed: %v", err)
	}

	// Compare configs
	if loadedConfig.FolderID != originalConfig.FolderID {
		t.Errorf("FolderID mismatch: got %s, want %s", loadedConfig.FolderID, originalConfig.FolderID)
	}
	if loadedConfig.FolderName != originalConfig.FolderName {
		t.Errorf("FolderName mismatch: got %s, want %s", loadedConfig.FolderName, originalConfig.FolderName)
	}
	if loadedConfig.TokenPath != originalConfig.TokenPath {
		t.Errorf("TokenPath mismatch: got %s, want %s", loadedConfig.TokenPath, originalConfig.TokenPath)
	}
}

func TestLoadCredentialsFileNotFound(t *testing.T) {
	// Temporarily change working directory to a non-existent path
	originalWd, _ := os.Getwd()
	tempDir, err := os.MkdirTemp("", "creds_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Test LoadCredentials with missing file
	_, err = LoadCredentials()
	if err == nil {
		t.Error("LoadCredentials() should fail when credentials file doesn't exist")
	}
}
