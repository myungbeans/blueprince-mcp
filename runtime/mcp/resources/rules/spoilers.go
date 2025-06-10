package rules

import (
	"context"

	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

// RegisterSpoilerRules adds a special resource containing the spoiler prevention rules
func RegisterSpoilerRules(ctx context.Context, s *server.MCPServer) error {
	logger := utils.Logger(ctx)

	rulesResource := mcp.NewResource(
		"rules://blue-prince/spoiler-protection",
		"Blue Prince Spoiler Protection Rules",
		mcp.WithResourceDescription("CRITICAL rules for spoiler-free Blue Prince assistance that the assistant must follow at all times"),
		mcp.WithMIMEType("text/markdown; charset=utf-8"),
	)

	rulesHandler := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		logger.Info("Loading spoiler prevention rules", zap.String("uri", req.Params.URI))
		return []mcp.ResourceContents{
			&mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "text/markdown; charset=utf-8",
				Text:     SpoilerRules,
			},
		}, nil
	}

	s.AddResource(rulesResource, rulesHandler)
	logger.Info("Registered spoiler prevention rules resource",
		zap.String("uri", "rules://blue-prince/spoiler-protection"))

	return nil
}
