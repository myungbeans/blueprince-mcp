package screenshots

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func AnalyzeTool() mcp.Tool {
	tool := mcp.Tool{
		Name: "analyze_screenshot",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"file_name": map[string]string{
					"type":        "string",
					"description": "Specific filename to look for in local vault",
				},
			},
		},
	}

	tool.Description = `
This Tool responds with b64 encoded data from the screenshot image file requesteed. 

This Tool is part of a multi-step WORKFLOW that is made up of 
1. download_screenshots
2. analyze_screenshot
3. create_note

WORKFLOW: 
After getting a successful response from this tool, you should:
1. Analyze the contents of the file. When analyzing:
- ONLY include descriptions of what is in the image
- NEVER make any logical leaps or connections to other notes
- NEVER add follow up questions or things to investigate
- ALWAYS be careful to avoid the possibility of spoilers
- ALWAYS ignore HUD elements in the top left and top right portions of the screen with counters for resources (such as steps, dice, keys, gems, and coins).
- ALWAYS annotate any text found in the image EXACTLY as it is on the screen.
- ALWAYS includendicate which room the screenshot is from (the room is indicated in the bottom right corner of the screenshot)
- If there is ANY part of your analysis that you are unsure of, indicate so by flagging it with "==NEEDS USER VERIFICATION==" in the note, being specific about which specific parts of your analysis you are unsure of.

2. Call the create_note tool with your analysis.
`
	return tool
}

// AnalyzeHandler creates a handler for analyzing screenshots
func AnalyzeHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
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

		// Check if file exists and read the content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			logger.Error("Failed to read screenshot file", zap.String("filePath", fullPath), zap.Error(err))
			if os.IsNotExist(err) {
				return mcp.NewToolResultError(fmt.Sprintf("Screenshot not found: '%s'", fullPath)), nil
			}
			return mcp.NewToolResultError(fmt.Sprintf("Failed to read screenshot file '%s': %v", fullPath, err)), nil
		}
		b64Data := base64.StdEncoding.EncodeToString(content)

		return mcp.NewToolResultText(fmt.Sprintf("data:image/jpeg;base64,%s", b64Data)), nil
	}
}
