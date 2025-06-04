package main

import (
	"fmt" // Ensure fmt is imported for printing errors
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// --- Configuration ---
const (
	defaultVaultBaseDir = "~" // used ONLY to place ~/Documents
	defaultVaultName    = "Documents/blueprince_mcp"
	configFile          = "config.yaml" // Assumes setup is run from project root
)

var requiredSubdirs = []string{"people", "puzzles", "rooms", "items", "lore", "general"}

var logger *zap.Logger // Global logger instance

var rootCmd = &cobra.Command{
	Use:   "setup [vault-path]",
	Short: "Sets up the Blue Prince MCP Obsidian vault.",
	Long: `This command initializes the Blue Prince MCP Obsidian vault.
It creates the specified vault directory (or uses the default "~/Documents/blueprince_mcp/")
and ensures the required subdirectories (/people, /puzzles, etc.) exist.
These subdirs were chosen to allow LLM clients to effectively lookup relevant info on searches,
while still providing human readabiliy.
Finally, it updates the 'config.yaml' with the correct vault path.`,
	Args: cobra.MaximumNArgs(1), // Allow zero or one argument for vault-path
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Starting Blue Prince MCP setup...")

		// Use default path if no argument is provided
		targetVaultPath := filepath.Join(defaultVaultBaseDir, defaultVaultName)
		if len(args) > 0 {
			targetVaultPath = args[0]
		}

		logger.Info("Using the vault path", zap.String("path", targetVaultPath))

		defer logger.Sync() // Flushes buffer, if any
		return setupVault(targetVaultPath)
	},
}

func setupVault(vaultPath string) error {
	// Clean the provided vaultPath
	absVaultPath, err := utils.ResolvePath(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to resolve and clean vault path: %w", err)
	}
	logger.Debug("Resolved target vault path", zap.String("resolvedPath", absVaultPath), zap.String("initialPath", vaultPath))

	// Create vault directory if it doesn't exist
	logger.Info("Ensuring vault directory exists", zap.String("path", absVaultPath))
	if err := utils.CreateDir(absVaultPath, 0755); err != nil { // 0755 permissions: owner rwx, group rx, others rx
		return fmt.Errorf("failed to ensure vault directory '%s': %w", absVaultPath, err)
	}
	logger.Info("Vault directory is ready", zap.String("path", absVaultPath))

	// Create required subdirectories
	for _, subdir := range requiredSubdirs {
		fullSubdirPath := filepath.Join(absVaultPath, subdir)
		logger.Info("Ensuring subdirectory exists", zap.String("subdir", subdir), zap.String("path", fullSubdirPath))
		if err := utils.CreateDir(fullSubdirPath, 0755); err != nil {
			return fmt.Errorf("failed to ensure subdirectory '%s': %w", fullSubdirPath, err)
		}
	}
	logger.Info("All required subdirectories are ready.")

	// Update config.yaml
	err = updateConfigVaultPath(configFile, absVaultPath)
	if err != nil {
		return fmt.Errorf("failed to update config file '%s': %w", configFile, err)
	}

	logger.Info("Setup complete!")
	logger.Info("Your Obsidian vault for Blue Prince MCP is configured", zap.String("vaultPath", absVaultPath))
	logger.Info("Please ensure your Obsidian application is pointed to this vault if you intend to use it directly.")
	return nil
}

func main() {
	// Validate that the setup script is run from the project root by checking for a sentinel file.
	// config.yaml is a good candidate since this flow writes back to config.yaml.
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Initialize a temporary logger for this specific error, as global logger isn't set yet.
		tempLogger, _ := zap.NewDevelopment()
		defer tempLogger.Sync()
		tempLogger.Fatal("Configuration file not found at root of project directory. Please run the setup command from the project root directory or ensure that config.yaml exists.",
			zap.String("configFile", configFile),
			zap.Error(err),
		)
		return
	}

	// Initialize logger
	var err error
	logger, err = zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flushes buffer, if any

	// Kick off here!
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
	return
}

// updateConfigVaultPath reads the config file, updates the obsidian_vault_path, and writes it back.
func updateConfigVaultPath(configPath, vaultPath string) error {
	logger.Info("Updating configuration file",
		zap.String("configFile", configPath),
		zap.String("newVaultPath", vaultPath),
	)

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

	configMap["obsidian_vault_path"] = vaultPath

	updatedYAML, err := yaml.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	// Create a backup before writing
	backupPath := configPath + ".bak"
	err = os.WriteFile(backupPath, yamlFile, 0644)
	if err != nil {
		logger.Warn("Failed to create backup of config file",
			zap.String("configFile", configPath),
			zap.String("backupPath", backupPath),
			zap.Error(err),
		)
	} else {
		logger.Info("Created backup of config file",
			zap.String("configFile", configPath),
			zap.String("backupPath", backupPath),
		)
	}

	// Write the updated config file even if backup fails (backup is just a nice to have)
	err = os.WriteFile(configPath, updatedYAML, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}
	logger.Info("Successfully updated configuration file", zap.String("configFile", configPath))
	return nil
}
