package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"go.uber.org/zap"
)

// ReadNoteHandler creates a handler for reading the content of a specific note.
func ReadNoteHandler(cfg *config.Config, logger *zap.Logger) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		noteID, err := request.RequireString("id")
		if err != nil {
			logger.Error("Missing required parameter 'id' for read_note", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Missing required parameter 'id': %v", err)), nil
		}

		// Clean the noteID to prevent path traversal issues (e.g. ".." components)
		// and ensure it's treated as a relative path.
		cleanNoteID := filepath.Clean(noteID)
		if strings.HasPrefix(cleanNoteID, "..") || filepath.IsAbs(cleanNoteID) {
			logger.Warn("Invalid note ID (path traversal attempt or absolute path)", zap.String("noteID", noteID), zap.String("cleanNoteID", cleanNoteID))
			return mcp.NewToolResultError(fmt.Sprintf("Invalid note ID: '%s'. Must be a relative path within the vault.", noteID)), nil
		}

		// Construct the full path to the note file
		// NoteID is expected to be relative to the NOTES_DIR
		fullPath := filepath.Join(cfg.ObsidianVaultPath, vault.NOTES_DIR, cleanNoteID)

		// Security check: Ensure the resolved path is still within the ObsidianVaultPath
		absVaultPath, _ := filepath.Abs(filepath.Join(cfg.ObsidianVaultPath, vault.NOTES_DIR)) // Check against notes dir
		absFullPath, err := filepath.Abs(fullPath)
		if err != nil || !strings.HasPrefix(absFullPath, absVaultPath) {
			logger.Warn("Path traversal attempt or invalid path for read_note", zap.String("noteID", noteID), zap.String("fullPath", fullPath), zap.String("absVaultPath", absVaultPath))
			return mcp.NewToolResultError(fmt.Sprintf("Access denied or invalid path for note ID: '%s'", noteID)), nil
		}

		// Read the file content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Failed to read note file", zap.String("filePath", fullPath), zap.Error(err))
			if os.IsNotExist(err) {
				return mcp.NewToolResultError(fmt.Sprintf("Note not found with ID: '%s'", noteID)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read note file with ID '%s': %v", noteID, err)), nil
		}

		logger.Info("Note read successfully", zap.String("filePath", fullPath))
		return mcp.NewToolResultText(string(content)), nil
	}
}
