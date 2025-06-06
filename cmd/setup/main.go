package main

import (
	"encoding/json"
	"fmt" // Ensure fmt is imported for printing errors
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// --- Configuration ---
const (
	defaultVaultBaseDir = "~" // used ONLY to place ~/Documents
	defaultVaultName    = "Documents/blueprince_mcp"
	yamlConfigFile      = "cmd/config/local/config.yaml"          // Relative to project root
	jsonConfigFile      = "cmd/config/claude_desktop/config.json" // Relative to project root
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
	err = updateYamlConfigVaultPath(yamlConfigFile, absVaultPath)
	if err != nil {
		return fmt.Errorf("failed to update config file '%s': %w", yamlConfigFile, err)
	}
	// Update claude_desktop_config.json
	err = updateClaudeDesktopConfig(jsonConfigFile, absVaultPath)
	if err != nil {
		return fmt.Errorf("failed to update claude desktop config file '%s': %w", jsonConfigFile, err)
	}

	logger.Info("Setup complete!")
	logger.Info("Your Obsidian vault for Blue Prince MCP is configured", zap.String("vaultPath", absVaultPath))
	logger.Info("Please ensure your Obsidian application is pointed to this vault if you intend to use it directly.")
	return nil
}

func main() {
	// Validate that the setup script is run from the project root by checking for a sentinel file.
	// We'll check for the YAML config file at its new location.
	if _, err := os.Stat(yamlConfigFile); os.IsNotExist(err) {
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

	// Kick off here!
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to execute root command", zap.Error(err))
		os.Exit(1)
	}
	return
}

// updateYamlConfigVaultPath reads the YAML config file, updates the OBSIDIAN_VAULT_PATH, and writes it back.
func updateYamlConfigVaultPath(configPath, vaultPath string) error {
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

	configMap["obsidian_vault_path"] = vaultPath // Use lowercase key to match YAML

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

// updateClaudeDesktopConfig modifies the sample config.json for Claude Desktop,
// and updates the OBSIDIAN_VAULT_PATH env var
func updateClaudeDesktopConfig(configPath, vaultPath string) error {
	logger.Info("Updating Claude Desktop JSON configuration file",
		zap.String("configFile", configPath),
		zap.String("newVaultPath", vaultPath),
	)

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

	// Navigate and update the path
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
		logger.Info("Claude desktop config: 'env' block not found, created a new one.")
	}

	env["OBSIDIAN_VAULT_PATH"] = vaultPath

	// Marshal back to JSON with indentation for readability
	updatedJSON, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated claude desktop config: %w", err)
	}

	// Create a backup before writing (optional, but good practice)
	backupPath := configPath + ".bak"
	if err := os.WriteFile(backupPath, jsonFile, 0644); err != nil {
		logger.Warn("Failed to create backup of claude desktop config file", zap.String("backupPath", backupPath), zap.Error(err))
	} else {
		logger.Info("Created backup of claude desktop config file", zap.String("backupPath", backupPath))
	}

	// Write the updated JSON config file
	if err := os.WriteFile(configPath, updatedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write updated claude desktop config file '%s': %w", configPath, err)
	}
	logger.Info("Successfully updated Claude Desktop JSON configuration file", zap.String("configFile", configPath))
	logger.Info("To complete setup in Claude Desktop, see https://modelcontextprotocol.io/quickstart/server#testing-your-server-with-claude-for-desktop", zap.String("configFile", configPath))
	return nil
}
