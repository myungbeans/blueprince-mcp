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
		logger.Info("Using Obsidian vault path from env variable", zap.String("variable", envVaultPath), zap.String("path", vaultPath))
		cfg = &config.Config{
			ObsidianVaultPath: vaultPath,
		}
	} else {
		cfg, err = config.LoadConfig(defaultConfigFilePath)
		if err != nil {
			logger.Fatal("Failed to load configuration", zap.Error(err))
		}
	}

	store, err := drive.NewStore(ctx, vaultPath)
	if err != nil {
		logger.Fatal("Failed to initialize store", zap.Error(err))
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
