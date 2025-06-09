package main

import (
	"fmt" // Ensure fmt is imported for printing errors
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/cmd/setup/drive"
	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// --- Configuration ---
const (
	defaultVaultBaseDir = "~" // used ONLY to place ~/Documents
	defaultVaultName    = "Documents/blueprince_mcp"
)

var baseRequiredSubdirs = []string{vault.META_DIR, vault.NOTES_DIR, vault.SCREENSHOT_DIR}

var logger *zap.Logger // Global logger instance

var rootCmd = &cobra.Command{
	Use:   "setup [vault-path]",
	Short: "Sets up the Blue Prince MCP Obsidian vault.",
	Long: `This command initializes the Blue Prince MCP Obsidian vault.
It creates the specified vault directory (or uses the default "~/Documents/blueprince_mcp/")
and ensures the required subdirectories (/notes, /meta, /screenshots) exist.
These subdirs were chosen to allow LLM clients to effectively lookup relevant info on searches,
while still providing human readabiliy.
Finally, it updates the 'config.yaml' and 'claude_desktop_config.json' files with the correct vault path.`,
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
	absVaultPath, err := utils.ResolveAndCleanPath(vaultPath)
	if err != nil {
		return fmt.Errorf("failed to resolve and clean vault path: %w", err)
	}

	// Create vault directory if it doesn't exist
	logger.Info("Ensuring vault directory exists", zap.String("path", absVaultPath))
	if err := utils.EnsureDirExists(absVaultPath, 0755); err != nil { // 0755 permissions: owner rwx, group rx, others rx
		return fmt.Errorf("failed to ensure vault directory '%s': %w", absVaultPath, err)
	}
	logger.Info("Vault directory is ready", zap.String("path", absVaultPath))

	// Create required subdirectories
	for _, subdir := range baseRequiredSubdirs {
		subpath := filepath.Join(absVaultPath, subdir)
		logger.Info("Ensuring subdirectory exists", zap.String("subdir", subdir), zap.String("path", subpath))
		if err := utils.EnsureDirExists(subpath, 0755); err != nil {
			return fmt.Errorf("failed to ensure subdirectory '%s': %w", subpath, err)
		}
	}
	logger.Info("All base required subdirectories are ready.")

	// Create note category subdirectories within NOTES_DIR
	notesDirPath := filepath.Join(absVaultPath, vault.NOTES_DIR)
	for _, categorySubdir := range notes.Categories {
		category := filepath.Join(notesDirPath, categorySubdir)
		logger.Info("Ensuring note category subdirectory exists", zap.String("subdir", categorySubdir), zap.String("path", category))
		if err := utils.EnsureDirExists(category, 0755); err != nil {
			return fmt.Errorf("failed to ensure ./note/{category} subdirectory '%s': %w", category, err)
		}
	}
	logger.Info("All ./note/{category} subdirectories are ready.")

	// Update config.yaml
	err = config.UpdateYamlField(config.YamlConfigFile, config.ObsidianVaultPathField, absVaultPath)
	if err != nil {
		return fmt.Errorf("failed to update config file '%s': %w", config.YamlConfigFile, err)
	}
	logger.Info("Successfully updated configuration file", zap.String("configFile", config.YamlConfigFile))

	// Update claude_desktop_config.json
	err = config.UpdateClaudeDesktopEnvVar(config.JsonConfigFile, config.ObsidianVaultPathEnv, absVaultPath)
	if err != nil {
		return fmt.Errorf("failed to update claude desktop config file '%s': %w", config.JsonConfigFile, err)
	}
	logger.Info("Successfully updated Claude Desktop JSON configuration file", zap.String("configFile", config.JsonConfigFile))
	logger.Info("To complete setup in Claude Desktop, see https://modelcontextprotocol.io/quickstart/server#testing-your-server-with-claude-for-desktop", zap.String("configFile", config.JsonConfigFile))

	logger.Info("Setup complete!")
	logger.Info("Your Obsidian vault for Blue Prince MCP is configured", zap.String("vaultPath", absVaultPath))
	logger.Info("Please ensure your Obsidian application is pointed to this vault if you intend to use it directly.")
	return nil
}

func main() {
	// Validate that the setup script is run from the project root by checking for a sentinel file.
	// We'll check for the YAML config file at its new location.
	if _, err := os.Stat(config.YamlConfigFile); os.IsNotExist(err) {
		// Initialize a temporary logger for this specific error, as global logger isn't set yet.
		tempLogger, _ := zap.NewDevelopment()
		defer tempLogger.Sync()
		tempLogger.Fatal("Please run the setup command from the project root directory.",
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

	// Add subcommands
	rootCmd.AddCommand(drive.DriveCmd)

	// Kick off here!
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
	return
}
