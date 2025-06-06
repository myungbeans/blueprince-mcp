package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"gopkg.in/yaml.v3"
)

// ServerConfig holds the server-specific configurations.
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Config holds all application configurations.
type Config struct {
	Server            ServerConfig `yaml:"server"`
	ObsidianVaultPath string       `yaml:"obsidian_vault_path"`
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
