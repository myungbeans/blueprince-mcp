package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	gdrive "github.com/myungbeans/blueprince-mcp/runtime/storage/drive"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveAuth struct {
	config     *oauth2.Config
	tokenFile  string
	configFile string
	service    *drive.Service
	ctx        context.Context
}

func NewGoogleDriveAuth(ctx context.Context) (*GoogleDriveAuth, error) {
	// Load credentials using shared utility
	config, err := gdrive.LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("unable to load credentials: %v", err)
	}

	// Set up token and config file paths
	tokenFile, err := gdrive.TokenPath()
	if err != nil {
		return nil, fmt.Errorf("unable to get token path: %v", err)
	}

	configFile, err := gdrive.ConfigPath()
	if err != nil {
		return nil, fmt.Errorf("unable to get config path: %v", err)
	}

	// Ensure config directory exists
	if err := gdrive.EnsureConfigDir(); err != nil {
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
	token, err := gdrive.LoadToken()
	if err != nil {
		// No existing token, start OAuth flow
		token, err = g.getTokenFromWeb()
		if err != nil {
			return fmt.Errorf("unable to retrieve token from web: %v", err)
		}

		// Save token for future use
		if err := gdrive.SaveToken(g.tokenFile, token); err != nil {
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
	gd := &gdrive.GoogleDrive{Client: g.service}
	folderID, err := gd.FindOrCreateFolder(folderName)
	if err != nil {
		return fmt.Errorf("unable to access folder '%s': %v", folderName, err)
	}

	// Test permissions by trying to list contents of the folder
	_, err = gd.ListFolderContents(folderID, 1)
	if err != nil {
		return fmt.Errorf("unable to list contents of folder '%s': %v", folderName, err)
	}

	fmt.Printf(" Successfully verified access to Google Drive folder: %s\n", folderName)
	return nil
}

func (g *GoogleDriveAuth) SaveConfig(folderName string) error {
	// Create GoogleDrive instance to use its methods
	gd := &gdrive.GoogleDrive{Client: g.service}

	// Find folder ID
	folderID, err := gd.FindOrCreateFolder(folderName)
	if err != nil {
		return fmt.Errorf("unable to find folder '%s': %v", folderName, err)
	}

	config := gdrive.DriveConfig{
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

func (g *GoogleDriveAuth) getTokenFromWeb() (*oauth2.Token, error) {
	// Create a channel to receive the authorization code
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Start a temporary HTTP server to handle the OAuth callback
	server := &http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Parse the authorization code from the callback URL
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no authorization code received")
			return
		}

		// Send success response to browser
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Blue Prince MCP - Authorization Complete</title>
    <style>
        body { font-family: Arial, sans-serif; text-align: center; padding: 50px; background-color: #f5f5f5; }
        .container { max-width: 500px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .success { color:rgb(39, 75, 174); font-size: 24px; margin-bottom: 20px; }
        .message { color: #2c3e50; font-size: 16px; line-height: 1.6; }
    </style>
</head>
<body>
    <div class="container">
        <div class="success">Authorization Successful!</div>
        <div class="message">
            You have successfully authenticated with Google Drive.<br>
            You can now close this window and return to your terminal.
        </div>
    </div>
</body>
</html>`)

		// Send the code through the channel
		codeCh <- code
	})

	// Start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("failed to start callback server: %v", err)
		}
	}()

	// Update the OAuth config to use the local callback server
	originalRedirectURL := g.config.RedirectURL
	g.config.RedirectURL = "http://localhost:8080"

	// Generate the authorization URL
	authURL := g.config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println("🔐 Starting Google Drive authentication...")
	fmt.Printf("📄 Please visit this URL to authorize the application:\n%s\n\n", authURL)
	fmt.Println("💡 Your browser should open automatically. If not, copy and paste the URL above.")
	fmt.Println("⏳ Waiting for authorization...")

	// Try to open the URL in the user's default browser
	go func() {
		time.Sleep(2 * time.Second) // Give server time to start
		openBrowser(authURL)
	}()

	// Wait for either the authorization code or an error
	var authCode string
	select {
	case code := <-codeCh:
		authCode = code
		fmt.Println("✅ Authorization received successfully!")
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("timeout waiting for authorization (5 minutes)")
	}

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Restore original redirect URL
	g.config.RedirectURL = originalRedirectURL

	// Exchange the authorization code for a token
	token, err := g.config.Exchange(g.ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}

	return token, nil
}

func openBrowser(url string) {
	var err error

	switch {
	case isCommand("open"): // macOS
		err = runCommand("open", url)
	case isCommand("xdg-open"): // Linux
		err = runCommand("xdg-open", url)
	case isCommand("cmd"): // Windows
		err = runCommand("cmd", "/c", "start", url)
	default:
		fmt.Printf("Please manually open this URL in your browser:\n%s\n", url)
		return
	}

	if err != nil {
		fmt.Printf("Could not open browser automatically. Please manually open this URL:\n%s\n", url)
	}
}

func isCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Start()
}
