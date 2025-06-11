package notes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// UpdateTool returns the configured mcp.Tool for updating notes
func UpdateTool() mcp.Tool {
	return mcp.Tool{
		Name:        "update_note",
		Description: "Updates an existing note by completely replacing it with new content and metadata. The MCP client should handle reading the existing note, merging user input with existing content, and providing the complete updated note. CRITICAL CONSTRAINTS: (1) NEVER add investigation questions, analysis prompts, or checklists unless explicitly requested. (2) The MCP client should preserve existing user observations and intelligently merge new content. (3) Focus on enhancing existing content rather than adding speculative material.",
		InputSchema: notes.GetMCPSchema(),
	}
}

// UpdateHandler creates a handler for updating existing notes
func UpdateHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for update_note"), nil
		}

		// Extract and validate content parameter
		content, err := utils.ExtractStringParam(params, "content")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		// Check for spoiler-risk content patterns
		if err := spoilerCheck(content); err != nil {
			logger.Warn("Content contains potential spoiler additions", zap.String("reason", err.Error()))
			return mcp.NewToolResultError(fmt.Sprintf("Content validation failed: %v. Please provide only the user's direct observations without additional analysis or investigation prompts.", err)), nil
		}

		// Extract and validate metadata parameter
		metadataRaw, ok := params["metadata"]
		if !ok {
			return mcp.NewToolResultError("Missing required parameter: metadata"), nil
		}
		metadataMap, ok := metadataRaw.(map[string]interface{})
		if !ok {
			return mcp.NewToolResultError("Parameter 'metadata' must be an object"), nil
		}

		metadata, err := notes.ParseMetadata(metadataMap)
		if err != nil {
			logger.Error("Failed to parse metadata", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Invalid metadata: %v", err)), nil
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

		// Check if file exists (this is what makes it an update vs create)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return mcp.NewToolResultError(fmt.Sprintf("Note not found: '%s'. Use create_note to create new notes.", notePath)), nil
		}

		// Ensure directory exists (should already exist, but safety check)
		dir := filepath.Dir(fullPath)
		if err := utils.EnsureDirExists(dir, 0755); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Directory error: %v", err)), nil
		}

		// Update timestamps - preserve created_at, update updated_at
		metadata.UpdatedAt = time.Now().Format(time.RFC3339)
		// Note: We trust the MCP client to preserve created_at from the existing note

		// Create updated file content
		fileContent, err := notes.CreateContent(metadata, content)
		if err != nil {
			logger.Error("Failed to create file content", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create file content: %v", err)), nil
		}

		// Write the updated content (completely overwrite)
		if err := os.WriteFile(fullPath, []byte(fileContent), 0644); err != nil {
			logger.Error("Failed to write updated note file", zap.String("path", fullPath), zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write updated note file: %v", err)), nil
		}

		logger.Info("Note updated successfully", zap.String("path", notePath), zap.String("category", metadata.Category))
		return mcp.NewToolResultText(fmt.Sprintf("Successfully updated note: %s", notePath)), nil
	}
}
