package notes

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// DeleteTool returns the configured mcp.Tool for deleting notes
func DeleteTool() mcp.Tool {
	return mcp.Tool{
		Name:        "delete_note",
		Description: "Deletes a specific note by its path. Use this to permanently remove a note file from the vault.",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"path": map[string]string{
					"type":        "string",
					"description": "Path to the note file relative to the notes directory (e.g., 'people/simon_jones.md', 'rooms/nook_tiger_paintings.md')",
				},
			},
			Required: []string{"path"},
		},
	}
}

// DeleteHandler creates a handler for deleting a specific note
func DeleteHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for delete_note"), nil
		}

		// Extract and validate path parameter
		notePath, err := utils.ExtractStringParam(params, "path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		// Validate and clean the path for security
		cleanPath, err := utils.ValidatePath(notePath)
		if err != nil {
			logger.Warn("Invalid note path", zap.String("originalPath", notePath), zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Build secure full path
		fullPath, err := utils.BuildSecurePath(cfg.ObsidianVaultPath, vault.NOTES_DIR, cleanPath)
		if err != nil {
			logger.Warn("Security validation failed for note path",
				zap.String("notePath", notePath),
				zap.String("cleanPath", cleanPath),
				zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Check if file exists before attempting deletion
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			logger.Warn("Note file not found for deletion", zap.String("path", notePath))
			return mcp.NewToolResultError(fmt.Sprintf("Note not found: '%s'", notePath)), nil
		}

		// Delete the file
		if err := os.Remove(fullPath); err != nil {
			logger.Error("Failed to delete note file", zap.String("filePath", fullPath), zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete note file '%s': %v", notePath, err)), nil
		}

		logger.Info("Note deleted successfully", zap.String("path", notePath))
		return mcp.NewToolResultText(fmt.Sprintf("Successfully deleted note: %s", notePath)), nil
	}
}
