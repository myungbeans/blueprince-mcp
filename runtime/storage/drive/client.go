package drive

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	VaultPath      string
	FolderID       string
	FolderName     string
	ScreenshotsDir string
	Client         *drive.Service
}

// DriveConfig represents the Google Drive configuration file generated during OAuth flow
type DriveConfig struct {
	FolderID   string `json:"folder_id"`
	FolderName string `json:"folder_name"`
	TokenPath  string `json:"token_path"`
}

func NewStore(ctx context.Context, vaultPath string) (*GoogleDrive, error) {
	// Load Google Drive configuration - this is where the user's token lives
	driveConfig, err := loadDriveConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Google Drive config: %w", err)
	}

	// Initialize Google Drive service
	driveService, err := initDriveSvc(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Google Drive service: %w", err)
	}

	screenshotsDir := filepath.Join(vaultPath, vault.SCREENSHOT_DIR)

	return &GoogleDrive{
		VaultPath:      vaultPath,
		FolderID:       driveConfig.FolderID,
		FolderName:     driveConfig.FolderName,
		ScreenshotsDir: screenshotsDir,
		Client:         driveService,
	}, nil
}

// init initializes the Google Drive client using the token stored on the Google Drive config file
func initDriveSvc(ctx context.Context) (*drive.Service, error) {
	// Load the token file
	token, err := LoadToken()
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	// Load client credentials
	config, err := LoadCredentials()
	if err != nil {
		return nil, err
	}

	// Create HTTP client with token
	client := config.Client(ctx, token)

	// Create Drive service
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %w", err)
	}

	return service, nil
}
