package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	BackupDirName     string       `yaml:"backup_dir_name"`
}

// LoadConfig reads the configuration from the given YAML file path.
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

	// Validate that the ObsidianVaultPath exists and is a directory
	info, err := os.Stat(cfg.ObsidianVaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config error: obsidian_vault_path '%s' does not exist", cfg.ObsidianVaultPath)
		}
		return nil, fmt.Errorf("config error: failed to stat obsidian_vault_path '%s': %w", cfg.ObsidianVaultPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("config error: obsidian_vault_path '%s' is not a directory", cfg.ObsidianVaultPath)
	}

	// Validate required subdirectories
	if err := validateVaultStructure(cfg.ObsidianVaultPath); err != nil {
		return nil, fmt.Errorf("config error in vault '%s': %w", cfg.ObsidianVaultPath, err)
	}

	return &cfg, nil
}

// validateVaultStructure checks for the presence of required subdirectories within the vault
func validateVaultStructure(vaultPath string) error {
	requiredSubdirs := []string{"people", "puzzles", "rooms", "items", "lore", "general"}

	for _, subdir := range requiredSubdirs {
		subdirPath := filepath.Join(vaultPath, subdir)
		info, err := os.Stat(subdirPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("required subdirectory '%s' not found", subdir)
			}
			return fmt.Errorf("failed to check subdirectory '%s': %w", subdir, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path '%s' exists but is not a directory, expected a subdirectory", subdir)
		}
	}
	return nil
}
