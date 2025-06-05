package config

import (
	"fmt"
	"os"
	"path/filepath"

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
		return nil, fmt.Errorf("config error: OBSIDIAN_VAULT_PATH must be set to a valid path in %s", configPath)
	}

	if cfg.ObsidianVaultPath == "/" {
		return nil, fmt.Errorf("config error: OBSIDIAN_VAULT_PATH cannot be the root directory '/' in %s", configPath)
	}

	if err := utils.IsDir(cfg.ObsidianVaultPath); err != nil {
		return nil, fmt.Errorf("config error for OBSIDIAN_VAULT_PATH: %w", err)
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
		if err := utils.IsDir(subdirPath); err != nil {
			return fmt.Errorf("required subdirectory '%s' error: %w", subdir, err)
		}
	}
	return nil
}
