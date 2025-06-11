package drive

import (
	"fmt"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/cmd/setup/drive/auth"
	"github.com/myungbeans/blueprince-mcp/runtime/storage/drive"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"github.com/spf13/cobra"
)

var DriveCmd = &cobra.Command{
	Use:   "drive <folder-name>",
	Short: "Configure Google Drive integration with OAuth authentication",
	Long: `Set up Google Drive integration for the Blue Prince MCP.
This command configures OAuth authentication and sets permissions for:
- View files and directories
- List files and directories  
- Edit files & directories (rename, move)
- Create directories
- Download files

Upon successful authentication, this command will:
- Save OAuth tokens locally for future use
- Update config.yaml with the google_drive_screenshot_folder setting
- Update claude_desktop/config.json with the GOOGLE_DRIVE_SCREENSHOT_FOLDER environment variable

You must specify a folder name to use as the root Google Drive folder for the integration.
The folder cannot be the root of your Google Drive - it must be a specific folder.

This should be run on the pre-built binary. 
If run locally via CLI, you must download the app's .credentials.json from SecretManager first and save it locally in the root dir of this project.
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		targetFolder := args[0]

		// Validate folder name
		if targetFolder == "" || targetFolder == "/" || targetFolder == "." {
			return fmt.Errorf("folder name cannot be empty or root directory")
		}

		// Initialize OAuth flow
		authenticator, err := auth.NewAuthenticator(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Google Drive authentication: %w", err)
		}
		// LoadToken will await from web response
		token, err := authenticator.LoadToken()
		if err != nil {
			return fmt.Errorf("failed to retrieve token: %w", err)
		}

		// Use token to set up the Google Drive svc
		authenticator.SetSvc(ctx, token)

		// Test permissions and folder access
		if err := authenticator.TestPermissions(targetFolder); err != nil {
			return fmt.Errorf("failed to verify Google Drive permissions: %w", err)
		}

		// Save configuration
		if err := authenticator.SaveConfig(targetFolder); err != nil {
			return fmt.Errorf("failed to save Google Drive configuration: %w", err)
		}

		// Update config files with the folder name
		if err := updateScreenshotFolderConfig(targetFolder); err != nil {
			return fmt.Errorf("failed to update configuration files: %w", err)
		}
		// Follow OAuth convention for the local location of secrets
		path, err := utils.ResolveAndCleanPath(filepath.Join("~", drive.CONFIG_DIR))
		if err != nil {
			return fmt.Errorf("failed to determine path to secrets: %w", err)
		}
		if err := updateSecretsLocation(path); err != nil {
			return fmt.Errorf("failed to update configuration files with path to secrets: %w", err)
		}
		pwd, err := utils.ResolveAndCleanPath(".")
		if err != nil {
			return fmt.Errorf("failed to determine pwd: %w", err)
		}
		if err := updateRoot(pwd); err != nil {
			return fmt.Errorf("failed to update configuration files with project root: %w", err)
		}

		fmt.Printf("✅ Google Drive setup completed successfully!\n")
		fmt.Printf("📁 Configured folder: %s\n", targetFolder)
		fmt.Printf("📝 Updated configuration files with folder settings\n")

		return nil
	},
}

// updateScreenshotFolderConfig updates both config.yaml and claude_desktop/config.json with the Google Drive folder
func updateScreenshotFolderConfig(folderName string) error {
	// Update config.yaml
	if err := config.UpdateYamlField(config.YamlConfigFile, config.GoogleDriveScreenshotFolderField, folderName); err != nil {
		return fmt.Errorf("failed to update YAML config: %w", err)
	}
	fmt.Printf("📝 Updated %s with folder: %s\n", config.YamlConfigFile, folderName)

	// Update claude_desktop/config.json
	if err := config.UpdateClaudeDesktopEnvVar(config.JsonConfigFile, config.GoogleDriveScreenshotFolderEnv, folderName); err != nil {
		return fmt.Errorf("failed to update Claude Desktop config: %w", err)
	}
	fmt.Printf("📝 Updated %s with folder: %s\n", config.JsonConfigFile, folderName)

	return nil
}

// updateSecretsLocation updates both config.yaml and claude_desktop/config.json with the local location of Google Drive secrets
func updateSecretsLocation(path string) error {
	// Update config.yaml
	if err := config.UpdateYamlField(config.YamlConfigFile, config.GoogleDriveSecretsField, path); err != nil {
		return fmt.Errorf("failed to update YAML config: %w", err)
	}
	fmt.Printf("📝 Updated %s with path: %s\n", config.YamlConfigFile, path)

	// Update claude_desktop/config.json
	if err := config.UpdateClaudeDesktopEnvVar(config.JsonConfigFile, config.GoogleDriveSecretsEnv, path); err != nil {
		return fmt.Errorf("failed to update Claude Desktop config: %w", err)
	}
	fmt.Printf("📝 Updated %s with path: %s\n", config.JsonConfigFile, path)

	return nil
}

// updateRoot updates both config.yaml and claude_desktop/config.json with the local location of Google Drive secrets
func updateRoot(path string) error {
	// Update config.yaml
	if err := config.UpdateYamlField(config.YamlConfigFile, config.RootField, path); err != nil {
		return fmt.Errorf("failed to update YAML config: %w", err)
	}
	fmt.Printf("📝 Updated %s with root: %s\n", config.YamlConfigFile, path)

	// Update claude_desktop/config.json
	if err := config.UpdateClaudeDesktopEnvVar(config.JsonConfigFile, config.RootEnv, path); err != nil {
		return fmt.Errorf("failed to update Claude Desktop config: %w", err)
	}
	fmt.Printf("📝 Updated %s with root: %s\n", config.JsonConfigFile, path)

	return nil
}
