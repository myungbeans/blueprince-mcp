package screenshots

import (
	"context"
	"fmt"
	"strings"

	"github.com/myungbeans/blueprince-mcp/runtime/models/storage"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func DownloadTool() mcp.Tool {
	tool := mcp.Tool{
		Name: "download_screenshots",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"file_name": map[string]string{
					"type":        "string",
					"description": "Specific filename to look for in Google Drive (scoped to the folder configured during setup)",
				},
			},
		},
	}

	tool.Description = `
This tool downloads screenshots from the Google Drive folder configured during setup. 
The tool can be configured to batch download multiple files or a single file. 
If the param "file_name" is an empty string, all files directly in the pre-configured Google Drive folder will be downloaded.
If the param "file_name" is not an empty string, the tool will attempt to download only the specified file.
Under the hood, this tool moves successfully downloaded files into an archived folder in Drive.	

This Tool is part of a multi-step WORKFLOW that is made up of 
1. download_screenshots
2. analyze_screenshot
3. create_note

WORKFLOW: After getting a successful response from this tool, you should:
1. Call the analyze_screenshot tool for each fileName in the response of this tool. fileNames are comma separated.
`
	return tool
}

// DownloadHandler creates a handler for downloading files from Goolge Drive
func DownloadHandler(ctx context.Context, store storage.Store) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for read_note"), nil
		}

		// Extract and validate file_name parameter
		fileName := ""
		fileNameRaw, ok := params["file_name"]
		if ok {
			fileName, ok = fileNameRaw.(string)
			if !ok {
				errMsg := "parameter 'file_name' must be a string"
				return mcp.NewToolResultError(errMsg), fmt.Errorf("bad parameter error: %s", errMsg)
			}
		}

		files, err := store.GetFiles(fileName)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), err
		}

		return mcp.NewToolResultText(strings.Join(files, ",")), nil
	}
}
