package handlers

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/server/config"
	"go.uber.org/zap"
)

// CreateNoteHandler creates a handler for creating new notes.
func CreateNoteHandler(cfg *config.Config, logger *zap.Logger) server.ToolHandlerFunc {
	return nil
}
