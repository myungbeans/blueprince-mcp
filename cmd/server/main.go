package main

import (
	"fmt"
	"os"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/resources"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/tools"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/server"
)

const (
	server_version        = "0.0.0"
	defaultConfigFilePath = "cmd/config/local/config.yaml"
	envVaultPath          = "OBSIDIAN_VAULT_PATH"
)

func main() {
	// Initialize logger early
	logger, err := zap.NewDevelopment() // Use NewDevelopment for more human-friendly output during dev
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flushes buffer, if any

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
			os.Exit(1)
		}
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Blue Prince Architect Notes",
		server_version,
		server.WithToolCapabilities(false),
	)

	// Register Resources
	if err := resources.Register(s, cfg.ObsidianVaultPath, logger); err != nil {
		logger.Fatal("Failed to register file resources", zap.Error(err))
		os.Exit(1)
	}

	// Register Tools
	tools.Register(s, cfg, logger)

	// Start the stdio server
	logger.Info("Initializing server...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}
