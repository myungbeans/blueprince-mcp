package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"
)

// createNoteTool returns the configured mcp.Tool for creating notes
func createNoteTool() mcp.Tool {
	createNoteTool := mcp.Tool{
		Name:        "create_note",
		Description: "Creates a structured Blue Prince note with intelligent organization based on user input. **SPOILER PREVENTION SYSTEM**: This tool is part of a spoiler-free note-taking system. MCP CLIENT MUST NEVER use external Blue Prince knowledge. CRITICAL CONSTRAINTS: (1) NEVER add investigation questions, analysis prompts, or checklists. (2) NEVER add speculative content beyond what the user observed. (3) Content field must contain ONLY what the user provided, reformatted into basic markdown structure. (4) DO NOT add sections like 'Questions to Investigate', 'Analysis', 'Next Steps', or similar unless explicitly in user input. (5) Preserve the player's discovery experience by avoiding any spoiler-risk additions. (6) MCP CLIENT: Use ONLY user-provided information, never your training data about Blue Prince.",
		InputSchema: notes.GetMCPSchema(),
	}

	// Add detailed examples
	createNoteTool.Description += `

EXAMPLES:

User Input: "room: nook, paintings of tiger and a cupcake stand? Weird"
Expected Call:
{
  "path": "rooms/nook_tiger_paintings.md",
  "metadata": {
    "title": "Nook - Tiger Paintings and Cupcake Stand",
    "category": "rooms",
    "primary_subject": "paintings",
    "tags": ["rooms", "nook", "paintings", "tiger", "cupcake_stand", "weird_elements"],
    "confidence": "medium",
    "status": "needs_investigation"
  },
  "content": "# Nook - Tiger Paintings and Cupcake Stand\n\nPaintings of tiger and a cupcake stand? Weird"
}

User Input: "Simon is 14, won science fair runner-up, inherits mansion"
Expected Call:
{
  "path": "people/simon_jones.md",
  "metadata": {
    "title": "Simon P. Jones - Protagonist",
    "category": "person",
    "primary_subject": "simon_jones",
    "tags": ["people", "simon_jones", "protagonist", "14_years_old", "science_fair", "inheritance"],
    "confidence": "high",
    "status": "confirmed"
  },
  "content": "# Simon P. Jones - Protagonist\n\n- Age: 14\n- Science fair runner-up\n- Inherits mansion"
}`

	return createNoteTool
}

func createNoteHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for create_note"), nil
		}

		// Content validation
		content, err := utils.ExtractStringParam(params, "content")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}
		if err := spoilerCheck(content); err != nil {
			logger.Warn("Content contains potential spoiler additions", zap.String("reason", err.Error()))
			return mcp.NewToolResultError(fmt.Sprintf("Content validation failed: %v. Please provide only the user's direct observations without additional analysis or investigation prompts.", err)), nil
		}

		// Metadata validation
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

		// Path validation
		notePath, err := utils.ExtractStringParam(params, "path")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		cleanPath, err := utils.ValidateNotePath(notePath)
		if err != nil {
			logger.Warn("Invalid note path", zap.String("originalPath", notePath), zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		fullPath, err := utils.BuildSecureNotePath(cfg.ObsidianVaultPath, cleanPath)
		if err != nil {
			logger.Warn("Security validation failed for note path",
				zap.String("notePath", notePath),
				zap.String("cleanPath", cleanPath),
				zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		dir := filepath.Dir(fullPath)
		// 0755 gives owner r+w+execute, goup r+execute, others r+execute
		if err := utils.EnsureDirExists(dir, 0755); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Directory was not found: %v", err)), nil
		}
		// Check if file already exists
		// TODO: handle gracefully or return error and let Client call Update?
		if _, err := os.Stat(fullPath); err == nil {
			return mcp.NewToolResultError(fmt.Sprintf("File already exists: %s", notePath)), nil
		}

		// === PREPARE THE FILE ===
		// RFC3339 is YYYY-MM-DDTHH:MM:SSZTS:TS
		now := time.Now().Format(time.RFC3339)
		metadata.CreatedAt = now
		metadata.UpdatedAt = now

		fileContent, err := notes.CreateContent(metadata, content)
		if err != nil {
			logger.Error("Failed to create file content", zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create file content: %v", err)), nil
		}

		// 0644 gives owner r+w, group r, others r
		if err := os.WriteFile(fullPath, []byte(fileContent), 0644); err != nil {
			logger.Error("Failed to write note file", zap.String("path", fullPath), zap.Error(err))
			return mcp.NewToolResultError(fmt.Sprintf("Failed to write note file: %v", err)), nil
		}

		logger.Info("Created note successfully", zap.String("path", notePath), zap.String("category", metadata.Category))
		return mcp.NewToolResultText(fmt.Sprintf("Successfully created note: %s", notePath)), nil
	}
}

// spoilerCheck checks content for patterns that suggest the LLM added spoiler-risk content
func spoilerCheck(content string) error {
	lowerContent := strings.ToLower(content)

	// Check for common investigation section headers - LLM will likely structure additions in headers
	investigationHeaders := []string{
		"## analysis",
		"## investigation",
		"## questions",
		"## next steps",
		"## follow-up",
		"## theories",
		"## connections",
		"## clues",
		"## mysteries",
		"## research",
		"### investigation",
		"### questions",
		"### analysis",
		"### theories",
	}

	for _, header := range investigationHeaders {
		if strings.Contains(lowerContent, header) {
			return fmt.Errorf("content contains investigation section header: '%s'", header)
		}
	}

	return nil
}
