package drive

import (
	"fmt"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/cmd/setup/drive/auth"
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
The folder cannot be the root of your Google Drive - it must be a specific folder.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetFolder := args[0]

		// Validate folder name
		if targetFolder == "" || targetFolder == "/" || targetFolder == "." {
			return fmt.Errorf("folder name cannot be empty or root directory")
		}

		// Get credentials path
		credentialsPath := filepath.Join("cmd", "setup", "drive", "auth", ".credentials.json")

		// Initialize OAuth flow
		authenticator, err := auth.NewGoogleDriveAuth(credentialsPath)
		if err != nil {
			return fmt.Errorf("failed to initialize Google Drive authentication: %w", err)
		}

		// Perform OAuth flow
		if err := authenticator.Authenticate(); err != nil {
			return fmt.Errorf("failed to authenticate with Google Drive: %w", err)
		}

		// Test permissions and folder access
		if err := authenticator.TestPermissions(targetFolder); err != nil {
			return fmt.Errorf("failed to verify Google Drive permissions: %w", err)
		}

		// Save configuration
		if err := authenticator.SaveConfig(targetFolder); err != nil {
			return fmt.Errorf("failed to save Google Drive configuration: %w", err)
		}

		// Update config files with the folder name
		if err := updateConfigFiles(targetFolder); err != nil {
			return fmt.Errorf("failed to update configuration files: %w", err)
		}

		fmt.Printf("✅ Google Drive setup completed successfully!\n")
		fmt.Printf("📁 Configured folder: %s\n", targetFolder)
		fmt.Printf("📝 Updated configuration files with folder settings\n")

		return nil
	},
}

// updateConfigFiles updates both config.yaml and claude_desktop/config.json with the Google Drive folder
func updateConfigFiles(folderName string) error {
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
