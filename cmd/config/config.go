package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"gopkg.in/yaml.v3"
)

// Constants for common config paths and field names
const (
	YamlConfigFile = "cmd/config/local/config.yaml"
	JsonConfigFile = "cmd/config/claude_desktop/config.json"

	// YAML field names
	ObsidianVaultPathField           = "obsidian_vault_path"
	GoogleDriveScreenshotFolderField = "google_drive_screenshot_folder"
	GoogleDriveSecretsField          = "google_drive_secrets_dir"
	RootField                        = "root"

	// Environment variable names for Claude Desktop
	ObsidianVaultPathEnv           = "OBSIDIAN_VAULT_PATH"
	GoogleDriveScreenshotFolderEnv = "GOOGLE_DRIVE_SCREENSHOT_FOLDER"
	GoogleDriveSecretsEnv          = "GOOGLE_DRIVE_SECRETS_DIR"
	RootEnv                        = "ROOT"
)

// ServerConfig holds the server-specific configurations.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Config holds all application configurations.
type Config struct {
	Server             ServerConfig `yaml:"server"`
	ObsidianVaultPath  string       `yaml:"obsidian_vault_path"`
	GoogleDriveFolder  string       `yaml:"google_drive_screenshot_folder"`
	GoogleDriveSecrets string       `yaml:"google_drive_secrets_dir"`
	Root               string       `yaml:"root"`
}

// LoadConfig reads the configuration from the given YAML file path and validates it.
func LoadConfig(configPath string) (*Config, error) {
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config error: failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: failed to unmarshal config file %s: %w", configPath, err)
	}

	if cfg.ObsidianVaultPath == "" || cfg.ObsidianVaultPath == "/path/to/your/obsidian/vault" {
		return nil, fmt.Errorf("config error: obsidian_vault_path must be set to a valid path in %s", configPath)
	}

	if cfg.ObsidianVaultPath == "/" {
		return nil, fmt.Errorf("config error: obsidian_vault_path cannot be the root directory '/' in %s", configPath)
	}

	if err := utils.ValidateDir(cfg.ObsidianVaultPath); err != nil {
		return nil, fmt.Errorf("config error for obsidian_vault_path: %w", err)
	}

	// Validate required subdirectories
	if err := validateBaseVaultStructure(cfg.ObsidianVaultPath); err != nil {
		return nil, fmt.Errorf("config error in vault '%s': %w", cfg.ObsidianVaultPath, err)
	}

	// Validate note category subdirectories within the NOTES_DIR
	notesDirPath := filepath.Join(cfg.ObsidianVaultPath, vault.NOTES_DIR)
	if err := validateNoteCategoriesStructure(notesDirPath); err != nil {
		return nil, fmt.Errorf("config error in notes directory structure within '%s': %w", notesDirPath, err)
	}

	return &cfg, nil
}

// validateBaseVaultStructure checks for the presence of top-level required subdirectories within the vault.
func validateBaseVaultStructure(vaultPath string) error {
	requiredSubdirs := []string{vault.NOTES_DIR, vault.META_DIR, vault.SCREENSHOT_DIR}
	for _, subdir := range requiredSubdirs {
		subdirPath := filepath.Join(vaultPath, subdir)
		if err := utils.ValidateDir(subdirPath); err != nil {
			return fmt.Errorf("required subdirectory '%s' error: %w", subdir, err)
		}
	}
	return nil
}

// validateNoteCategoriesStructure checks for the presence of note category subdirectories.
func validateNoteCategoriesStructure(notesBasePath string) error {
	for _, category := range notes.Categories {
		categoryPath := filepath.Join(notesBasePath, category)
		if err := utils.ValidateDir(categoryPath); err != nil {
			return fmt.Errorf("required note category subdirectory '%s' error: %w", category, err)
		}
	}
	return nil
}

// UpdateYamlField updates a specific field in a YAML config file
func UpdateYamlField(configPath, fieldName, value string) error {
	// Read the existing config file
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configMap map[string]any
	err = yaml.Unmarshal(yamlFile, &configMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	// Update the specified field
	configMap[fieldName] = value

	updatedYAML, err := yaml.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	// Create a backup before writing
	if err := createBackup(configPath, yamlFile); err != nil {
		fmt.Printf("⚠️  Warning: Failed to create backup of config file: %v\n", err)
	}

	// Write the updated config file
	err = os.WriteFile(configPath, updatedYAML, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	return nil
}

// UpdateClaudeDesktopEnvVar updates an environment variable in the Claude Desktop config
func UpdateClaudeDesktopEnvVar(configPath, envVarName, value string) error {
	// Read the existing JSON config file
	jsonFile, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read claude desktop config file '%s': %w", configPath, err)
	}

	var configMap map[string]interface{}
	err = json.Unmarshal(jsonFile, &configMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal claude desktop config file '%s': %w", configPath, err)
	}

	// Navigate and update the environment variable
	mcpServers, ok := configMap["mcpServers"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("claude desktop config: 'mcpServers' key not found or not a map")
	}

	blueprinceServer, ok := mcpServers["blueprince_notes_mcp"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("claude desktop config: 'mcpServers.blueprince_notes_mcp' key not found or not a map")
	}

	env, ok := blueprinceServer["env"].(map[string]interface{})
	if !ok {
		// If env block doesn't exist, create it
		env = make(map[string]interface{})
		blueprinceServer["env"] = env
		fmt.Printf("📝 Created 'env' block in Claude Desktop config\n")
	}

	env[envVarName] = value

	// Marshal back to JSON with indentation for readability
	updatedJSON, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated claude desktop config: %w", err)
	}

	// Create a backup before writing
	if err := createBackup(configPath, jsonFile); err != nil {
		fmt.Printf("⚠️  Warning: Failed to create backup of claude desktop config file: %v\n", err)
	}

	// Write the updated JSON config file
	if err := os.WriteFile(configPath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write updated claude desktop config file '%s': %w", configPath, err)
	}

	return nil
}

// createBackup creates a backup file with .bak extension
func createBackup(originalPath string, content []byte) error {
	backupPath := originalPath + ".bak"
	return os.WriteFile(backupPath, content, 0644)
}
