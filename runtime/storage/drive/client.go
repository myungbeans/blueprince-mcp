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
	secretsPath    string
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

func NewStore(ctx context.Context, svc *drive.Service, vaultPath, secretsPath, folderID string) *GoogleDrive {
	screenshotsDir := filepath.Join(vaultPath, vault.SCREENSHOT_DIR)
	return &GoogleDrive{
		VaultPath:      vaultPath,
		secretsPath:    secretsPath,
		FolderID:       folderID,
		ScreenshotsDir: screenshotsDir,
		Client:         svc,
	}
}

func GetSvc(ctx context.Context, secretsPath, credsPath string) (*drive.Service, error) {
	// Load the token file
	token, err := LoadToken(secretsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	// Load client credentials
	creds, err := LoadCredentials(credsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load credentials: %w", err)
	}

	// Create HTTP client with token
	client := creds.Client(ctx, token)

	// Create Google Drive service
	service, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create Drive service: %w", err)
	}
	return service, nil
}
