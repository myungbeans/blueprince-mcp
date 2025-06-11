package screenshots

import (
	"context"
	"strings"

	"github.com/myungbeans/blueprince-mcp/runtime/models/storage"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ListTool() mcp.Tool {
	return mcp.Tool{
		Name:        "list_screenshots",
		Description: "Lists all screenshots within the pre-configured Google Drive folder. A successful response includes a comma separated list of file names",
		InputSchema: mcp.ToolInputSchema{
			Type:       "object",
			Properties: map[string]any{},
		},
	}
}

// ListHandler creates a handler for listing files from Goolge Drive
func ListHandler(ctx context.Context, store storage.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		files, err := store.ListFiles()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), err
		}

		return mcp.NewToolResultText(strings.Join(files, ",")), nil
	}
}
