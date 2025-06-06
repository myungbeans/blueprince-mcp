package tools

import (
	"github.com/myungbeans/blueprince-mcp/cmd/config"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Register(s *server.MCPServer, cfg *config.Config, logger *zap.Logger) {
	listNotesTool := mcp.NewTool("list_notes",
		mcp.WithDescription("Lists all notes."),
		// TODO: Add parameters for filtering later if needed
	)

	// Add Note tool handlers
	s.AddTool(listNotesTool, ListNotesHandler(cfg, logger))
	s.AddTool(createNoteTool(), createNoteHandler(cfg, logger))
}
