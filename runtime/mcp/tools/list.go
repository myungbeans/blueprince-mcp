package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// ListNotesHandler creates a handler for listing notes, with access to the application config.
func ListNotesHandler(cfg *config.Config, logger *zap.Logger) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		notesDir := filepath.Join(cfg.ObsidianVaultPath, vault.NOTES_DIR)
		relativeFilePaths, err := utils.ListFiles(notesDir)
		if err != nil {
			logger.Error("Error listing notes", zap.String("notesDir", notesDir), zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Error listing notes from '%s': %v", notesDir, err)), nil
		}

		if len(relativeFilePaths) == 0 {
			return mcp.NewToolResultText(fmt.Sprintf("No notes found in: %s", notesDir)), nil
		}

		resultText := "Notes:\n" + strings.Join(relativeFilePaths, "\n")
		return mcp.NewToolResultText(resultText), nil
	}
}
