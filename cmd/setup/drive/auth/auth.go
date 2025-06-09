package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveAuth struct {
	config      *oauth2.Config
	tokenFile   string
	configFile  string
	service     *drive.Service
	ctx         context.Context
}

type DriveConfig struct {
	FolderID   string `json:"folder_id"`
	FolderName string `json:"folder_name"`
	TokenPath  string `json:"token_path"`
}

func NewGoogleDriveAuth(credentialsPath string) (*GoogleDriveAuth, error) {
	ctx := context.Background()
	
	// Read credentials file
	b, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	// Parse credentials and create config
	config, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	// Set up token and config file paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("unable to get user home directory: %v", err)
	}

	tokenFile := filepath.Join(homeDir, ".blueprince_mcp", "drive_token.json")
	configFile := filepath.Join(homeDir, ".blueprince_mcp", "drive_config.json")

	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(tokenFile), 0700); err != nil {
		return nil, fmt.Errorf("unable to create config directory: %v", err)
	}

	return &GoogleDriveAuth{
		config:     config,
		tokenFile:  tokenFile,
		configFile: configFile,
		ctx:        ctx,
	}, nil
}

func (g *GoogleDriveAuth) Authenticate() error {
	// Try to load existing token
	token, err := g.loadToken()
	if err != nil {
		// No existing token, start OAuth flow
		token, err = g.getTokenFromWeb()
		if err != nil {
			return fmt.Errorf("unable to retrieve token from web: %v", err)
		}
		
		// Save token for future use
		if err := g.saveToken(token); err != nil {
			return fmt.Errorf("unable to save token: %v", err)
		}
	}

	// Create Drive service
	client := g.config.Client(g.ctx, token)
	service, err := drive.NewService(g.ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	g.service = service
	return nil
}

func (g *GoogleDriveAuth) TestPermissions(folderName string) error {
	if g.service == nil {
		return fmt.Errorf("not authenticated - call Authenticate() first")
	}

	// Test basic access by listing files
	_, err := g.service.Files.List().PageSize(1).Do()
	if err != nil {
		return fmt.Errorf("unable to access Google Drive: %v", err)
	}

	// Find or create the specified folder
	folderID, err := g.findOrCreateFolder(folderName)
	if err != nil {
		return fmt.Errorf("unable to access folder '%s': %v", folderName, err)
	}

	// Test permissions by trying to list contents of the folder
	_, err = g.service.Files.List().Q(fmt.Sprintf("'%s' in parents", folderID)).PageSize(1).Do()
	if err != nil {
		return fmt.Errorf("unable to list contents of folder '%s': %v", folderName, err)
	}

	fmt.Printf(" Successfully verified access to Google Drive folder: %s\n", folderName)
	return nil
}

func (g *GoogleDriveAuth) SaveConfig(folderName string) error {
	// Find folder ID
	folderID, err := g.findOrCreateFolder(folderName)
	if err != nil {
		return fmt.Errorf("unable to find folder '%s': %v", folderName, err)
	}

	config := DriveConfig{
		FolderID:   folderID,
		FolderName: folderName,
		TokenPath:  g.tokenFile,
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal config: %v", err)
	}

	if err := os.WriteFile(g.configFile, configData, 0600); err != nil {
		return fmt.Errorf("unable to save config file: %v", err)
	}

	fmt.Printf(" Configuration saved to: %s\n", g.configFile)
	return nil
}

func (g *GoogleDriveAuth) findOrCreateFolder(folderName string) (string, error) {
	// Search for existing folder
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", folderName)
	r, err := g.service.Files.List().Q(query).Do()
	if err != nil {
		return "", fmt.Errorf("unable to search for folder: %v", err)
	}

	if len(r.Files) > 0 {
		// Folder exists
		return r.Files[0].Id, nil
	}

	// Create new folder
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	file, err := g.service.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder: %v", err)
	}

	fmt.Printf(" Created new Google Drive folder: %s\n", folderName)
	return file.Id, nil
}

func (g *GoogleDriveAuth) getTokenFromWeb() (*oauth2.Token, error) {
	authURL := g.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	fmt.Print("Enter authorization code: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	token, err := g.config.Exchange(g.ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return token, nil
}

func (g *GoogleDriveAuth) loadToken() (*oauth2.Token, error) {
	f, err := os.Open(g.tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func (g *GoogleDriveAuth) saveToken(token *oauth2.Token) error {
	f, err := os.OpenFile(g.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}