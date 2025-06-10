package drive

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func CredsPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return filepath.Join(wd, APP_CREDS_FILE), nil
}

func TokenPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}
	return filepath.Join(homeDir, CONFIG_DIR, TOKEN_FILE), nil
}

func ConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to get user home directory: %w", err)
	}

	return filepath.Join(homeDir, CONFIG_DIR, CONFIG_FILE), nil
}

func EnsureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get user home directory: %w", err)
	}
	configDir := filepath.Join(homeDir, CONFIG_DIR)
	return utils.EnsureDirExists(configDir, 0700)
}

// loadDriveConfig loads the Google Drive Configuration file into its struct model
func loadDriveConfig() (*DriveConfig, error) {
	configPath, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read drive config file: %w", err)
	}

	var config DriveConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("failed to parse drive config: %w", err)
	}

	return &config, nil
}

// LoadToken loads an OAuth2 token from file
func LoadToken() (*oauth2.Token, error) {
	tokenPath, err := TokenPath()
	if err != nil {
		return nil, err
	}

	tokenData, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal(tokenData, &token)
	return &token, err
}

// SaveToken saves an OAuth2 token to a local file
func SaveToken(tokenPath string, token *oauth2.Token) error {
	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

// LoadCredentials loads Google Drive credentials and creates OAuth config
func LoadCredentials() (*oauth2.Config, error) {
	credsPath, err := CredsPath()
	if err != nil {
		return nil, err
	}

	credentialsData, err := os.ReadFile(credsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	config, err := google.ConfigFromJSON(credentialsData, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return config, nil
}
