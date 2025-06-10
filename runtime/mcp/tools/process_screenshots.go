package tools

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/mark3labs/mcp-go/mcp"
// 	"github.com/myungbeans/blueprince-mcp/cmd/config"
// 	"github.com/myungbeans/blueprince-mcp/runtime/utils"
// 	"go.uber.org/zap"
// )

// // ProcessScreenshotsRequest represents the input for processing screenshots from Google Drive
// type ProcessScreenshotsRequest struct {
// 	MaxScreenshots int  `json:"max_screenshots,omitempty"` // Optional limit on number of screenshots to process
// 	ForceReprocess bool `json:"force_reprocess,omitempty"` // Optional flag to reprocess already downloaded screenshots
// }

// // ProcessScreenshotsResponse represents the output of screenshot processing
// type ProcessScreenshotsResponse struct {
// 	ProcessedCount int      `json:"processed_count"`
// 	SkippedCount   int      `json:"skipped_count"`
// 	ErrorCount     int      `json:"error_count"`
// 	ProcessedFiles []string `json:"processed_files"`
// 	Errors         []string `json:"errors,omitempty"`
// }

// // GetScreenshotsTool creates the MCP tool for fetching screenshots from a screenshots Store
// func GetScreenshotsTool() mcp.Tool {
// 	return mcp.Tool{
// 		Name:        "get_screenshots",
// 		Description: "Downloads screenshots from the configured Store (e.g. Google Drive folder) and store them to the local vault under `./screenshots/`. The force_reprocess arg allows repeat downloads of files already fetched.",
// 		InputSchema: mcp.ToolInputSchema{
// 			Type: "object",
// 			Properties: map[string]any{
// 				"max_screenshots": map[string]any{
// 					"type":        "number",
// 					"description": "Maximum number of screenshots to process (default: 10)",
// 				},
// 				"force_reprocess": map[string]any{
// 					"type":        "boolean",
// 					"description": "Force reprocessing of already downloaded screenshots (default: false)",
// 				},
// 			},
// 		},
// 	}
// }

// // GetScreenshotHandler creates the MCP handler for getting and downloading screenshots from a Store
// func GetScreenshotHandler(ctx context.Context, cfg *config.Config) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
// 		logger := utils.Logger(ctx)
// 		logger.Info("Processing screenshots tool called", zap.String("tool", req.Params.Name))

// 		args, ok := req.Params.Arguments.(map[string]interface{})
// 		if !ok {
// 			args = make(map[string]interface{})
// 		}
// 		result, err := ProcessScreenshots(ctx, cfg.ObsidianVaultPath, args)
// 		if err != nil {
// 			logger.Error("Screenshot processing failed", zap.Error(err))
// 			return mcp.NewToolResultError(fmt.Sprintf("Screenshot processing failed: %v", err)), nil
// 		}

// 		return mcp.NewToolResultText(fmt.Sprintf(
// 			"Screenshot processing completed:\n- Processed: %d\n- Skipped: %d\n- Errors: %d\n- Files: %v",
// 			result.ProcessedCount,
// 			result.SkippedCount,
// 			result.ErrorCount,
// 			result.ProcessedFiles,
// 		)), nil
// 	}
// }

// // ProcessScreenshots handles the screenshot processing workflow
// func ProcessScreenshots(ctx context.Context, vaultPath string, args map[string]interface{}) (*ProcessScreenshotsResponse, error) {
// 	logger := utils.Logger(ctx)
// 	logger.Info("Starting screenshot processing from Google Drive")

// 	// Parse arguments
// 	var req ProcessScreenshotsRequest
// 	if argsJSON, err := json.Marshal(args); err == nil {
// 		json.Unmarshal(argsJSON, &req)
// 	}

// 	// Set defaults
// 	if req.MaxScreenshots <= 0 {
// 		req.MaxScreenshots = 10
// 	}

// 	logger.Info("Processing screenshots",
// 		zap.Int("maxScreenshots", req.MaxScreenshots),
// 		zap.Bool("forceReprocess", req.ForceReprocess),
// 	)

// 	// Initialize screenshot processor
// 	processor, err := NewScreenshotProcessor(ctx, vaultPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize screenshot processor: %w", err)
// 	}

// 	// Process screenshots
// 	result, err := processor.ProcessScreenshots(ctx, req.MaxScreenshots, req.ForceReprocess)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to process screenshots: %w", err)
// 	}

// 	logger.Info("Screenshot processing completed",
// 		zap.Int("processed", result.ProcessedCount),
// 		zap.Int("skipped", result.SkippedCount),
// 		zap.Int("errors", result.ErrorCount),
// 	)

// 	return result, nil
// }
