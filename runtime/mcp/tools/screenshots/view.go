package screenshots

import (
	"context"
	"fmt"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// TODO: Images exceed max length allowed in a conversation on Claude Desktop.
// Needs to be batched or streamed. WIP

func ViewTool() mcp.Tool {
	tool := mcp.Tool{
		Name: "view_screenshot",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"file_name": map[string]string{
					"type":        "string",
					"description": "File name of the screenshot file. The file must be directly in the vault's ./screenshots dir",
				},
			},
		},
	}

	tool.Description = `
This Tool displays the contents of a screenshot file.

This Tool is part of a multi-step WORKFLOW that is made up of 
1. download_screenshots
2. view_screenshot
3. analyze_screenshot
4. create_note

WORKFLOW: After getting a successful response from this tool, you should iterate over each downloaded file. For each fileName in the response's comma separated list of fileNames:
- Immediately after the response from view_screenshot, call the analyze_screenshot tool. This tool will ask the MCP Host to analyze the contents of the screenshot and format them in the appropriate format for a create_note tool call.
- Call the create_note tool. This tool will store the outputs of analyze_screenshot into a note containing the analyzed contents of the screenshot.
`
	return tool
}

// ViewHandler creates a handler for viewing the contents of a screenshot file
func ViewHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)

	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := request.GetArguments()
		if params == nil {
			return mcp.NewToolResultError("Missing arguments for view_screenshot"), nil
		}

		// Extract and validate path parameter
		imgName, err := utils.ExtractStringParam(params, "file_name")
		if err != nil || imgName == "" {
			return mcp.NewToolResultError(fmt.Sprintf("Parameter validation failed: %v", err)), nil
		}

		// Build the path, validate it
		cleanFilePath, err := utils.ValidatePath(imgName)
		fullPath, err := utils.BuildSecurePath(cfg.ObsidianVaultPath, vault.SCREENSHOT_DIR, cleanFilePath)
		if err != nil {
			logger.Warn("Invalid screenshot path", zap.String("path", fullPath), zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Toggle aspect ratios as needed
		smImg, err := utils.CompressImage(fullPath, 200, 200, 70)
		if err != nil {
			logger.Warn("Failed to compress the image", zap.String("path", fullPath), zap.Error(err))
			return mcp.NewToolResultError(err.Error()), nil
		}

		mimeType := utils.GetMimeType(fullPath)
		return mcp.NewToolResultImage("Screenshot loaded successfully", smImg.Data, mimeType), nil

		// // Check if file exists and read the content
		// content, err := os.ReadFile(fullPath)
		// if err != nil {
		// 	logger.Error("Failed to read screenshot file", zap.String("filePath", fullPath), zap.Error(err))
		// 	if os.IsNotExist(err) {
		// 		return mcp.NewToolResultError(fmt.Sprintf("Screenshot not found: '%s'", fullPath)), nil
		// 	}
		// 	return mcp.NewToolResultError(fmt.Sprintf("Failed to read screenshot file '%s': %v", fullPath, err)), nil
		// }
		// b64Data := base64.StdEncoding.EncodeToString(content)

		// return mcp.NewToolResultText(fmt.Sprintf("data:image/jpeg;base64,%s", b64Data)), nil
	}
}
