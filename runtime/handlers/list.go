package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/server/config"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// ListNotesHandler creates a handler for listing notes, with access to the application config.
func ListNotesHandler(cfg *config.Config, logger *zap.Logger) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileNames, err := utils.ListFiles(cfg.ObsidianVaultPath)
		if err != nil {
			// You could log the detailed error server-side here if desired
			logger.Error("Error listing notes", zap.String("vaultPath", cfg.ObsidianVaultPath), zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Error listing notes: %s. Please ensure 'obsidian_vault_path' is correct.", err.Error())), nil
		}

		if len(fileNames) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No notes found in the vault: %s", cfg.ObsidianVaultPath)), nil
		}

		// Format the output as a newline-separated list of file names.
		resultText := "Files in vault (" + cfg.ObsidianVaultPath + "):\n" + strings.Join(fileNames, "\n")
		return mcp.NewToolResultText(resultText), nil
	}
}
