package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// ReadNoteTool returns the configured mcp.Tool for reading notes
func ReadNoteTool() mcp.Tool {
	return mcp.Tool{
		Name:        "read_note",
		Description: "Reads the content of a specific note by its path. Use this to retrieve the full content of a note file including metadata and content.",
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

// ReadNoteHandler creates a handler for reading the content of a specific note
func ReadNoteHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for read_note"), nil
		}

		// Extract and validate path parameter
		notePath, err := utils.ExtractStringParam(params, "path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		// Validate and clean the path for security
		cleanPath, err := utils.ValidateNotePath(notePath)
		if err != nil {
			logger.Warn("Invalid note path", zap.String("originalPath", notePath), zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Build secure full path
		fullPath, err := utils.BuildSecureNotePath(cfg.ObsidianVaultPath, cleanPath)
		if err != nil {
			logger.Warn("Security validation failed for note path",
				zap.String("notePath", notePath),
				zap.String("cleanPath", cleanPath),
				zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Check if file exists and read the content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Failed to read note file", zap.String("filePath", fullPath), zap.Error(err))
			if os.IsNotExist(err) {
				return mcp.NewToolResultError(fmt.Sprintf("Note not found: '%s'", notePath)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read note file '%s': %v", notePath, err)), nil
		}

		logger.Info("Note read successfully", zap.String("path", notePath))
		return mcp.NewToolResultText(string(content)), nil
	}
}
