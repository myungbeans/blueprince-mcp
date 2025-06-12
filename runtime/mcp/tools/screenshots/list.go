package screenshots

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/storage"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	LOCAL_SRC        = "local"
	GOOGLE_DRIVE_SRC = "google drive"
)

var sources = []string{LOCAL_SRC, GOOGLE_DRIVE_SRC}

func ListTool() mcp.Tool {
	return mcp.Tool{
		Name:        "list_screenshots",
		Description: "Lists all screenshots within the pre-configured Google Drive folder. A successful response includes a comma separated list of file names",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"source": map[string]any{
					"type":        "string",
					"description": "Source to lookup screenshots from.",
					"enum":        sources,
				},
			},
		},
	}
}

// ListHandler creates a handler for listing files from Goolge Drive
func ListHandler(ctx context.Context, cfg *config.Config, store storage.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for view_screenshot"), nil
		}

		// Extract and validate path parameter
		source, err := utils.ExtractStringParam(params, "source")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		switch source {
		case LOCAL_SRC:
			return listLocalScreenshots(cfg)
		case GOOGLE_DRIVE_SRC:
			files, err := store.ListFiles()
			if err != nil {
				return mcp.NewToolResultError(err.Error()), err
			}
			return mcp.NewToolResultText(strings.Join(files, ",")), nil
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Invalid source: %s", source)), nil
		}
	}
}

func listLocalScreenshots(cfg *config.Config) (*mcp.CallToolResult, error) {
	imgDir := filepath.Join(cfg.ObsidianVaultPath, vault.SCREENSHOT_DIR)
	relativeFilePaths, err := utils.ListFiles(imgDir)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error listing notes from '%s': %v", imgDir, err)), err
	}

	if len(relativeFilePaths) == 0 {
		return mcp.NewToolResultText(fmt.Sprintf("No notes found in: %s", imgDir)), nil
	}

	resultText := "Notes:\n" + strings.Join(relativeFilePaths, "\n")
	return mcp.NewToolResultText(resultText), nil
}
