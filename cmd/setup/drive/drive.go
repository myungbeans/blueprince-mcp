package drive

import (
	"fmt"
	"path/filepath"

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

		return nil
	},
}
