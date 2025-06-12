package main

import (
	"context"
	"fmt"
	"os"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime"
	"github.com/myungbeans/blueprince-mcp/runtime/storage/drive"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/server"
)

const (
	server_version        = "0.0.1"
	defaultConfigFilePath = "cmd/config/local/config.yaml"
	envVaultPath          = "OBSIDIAN_VAULT_PATH"
	envRoot               = "ROOT"
	envGoogleDriveFolder  = "GOOGLE_DRIVE_SCREENSHOT_FOLDER"
	envGoogleDriveSecrets = "GOOGLE_DRIVE_SECRETS_DIR"
	loggerKey             = "logger"
)

func main() {
	// Initialize logger early
	// TODO: using NewDevelopment for more human-friendly output during dev
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger. Driving blind is dangerous. Abort!\nError: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flushes buffer, if any

	ctx := context.WithValue(context.Background(), loggerKey, logger)
	var cfg *config.Config
	vaultPath := os.Getenv(envVaultPath)
	if vaultPath != "" {
		logger.Info("Using env vars for configs")
		cfg = &config.Config{
			ObsidianVaultPath:  vaultPath,
			Root:               os.Getenv(envRoot),
			GoogleDriveFolder:  os.Getenv(envGoogleDriveFolder),
			GoogleDriveSecrets: os.Getenv(envGoogleDriveSecrets),
		}
	} else {
		cfg, err = config.LoadConfig(defaultConfigFilePath)
		if err != nil {
			logger.Fatal("Failed to load configuration", zap.Error(err))
		}
	}

	var store *drive.GoogleDrive
	if cfg.GoogleDriveSecrets != "" {
		svc, err := drive.GetSvc(ctx, cfg.GoogleDriveSecrets, cfg.Root)
		if err != nil {
			logger.Fatal("Failed to get Google Drive client", zap.Error(err))
		}

		// Load Google Drive configuration - this is where the user's token lives
		driveConfig, err := drive.LoadDriveConfig(cfg.GoogleDriveSecrets)
		if err != nil {
			logger.Fatal("Failed to load Google Drive config", zap.Error(err))
		}

		store = drive.NewStore(ctx, svc, vaultPath, cfg.GoogleDriveSecrets, driveConfig.FolderID)
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Blue Prince Architect Notes - SPOILER-FREE Note Taking",
		server_version,
	)

	rtime := runtime.NewHandler(cfg, store)
	err = rtime.RegisterResources(ctx, s)
	if err != nil {
		logger.Fatal("Failed to register resources", zap.Error(err))
	}

	rtime.RegisterTools(ctx, s)

	// Start the stdio server
	logger.Info("Initializing server...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}
