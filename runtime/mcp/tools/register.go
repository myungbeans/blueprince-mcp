package tools

import (
	"context"

	"github.com/myungbeans/blueprince-mcp/cmd/config"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func Register(ctx context.Context, s *server.MCPServer, cfg *config.Config) {
	listNotesTool := mcp.NewTool("list_notes",
		mcp.WithDescription("Lists all notes."),
		// TODO: Add parameters for filtering later if needed
	)

	// Add Note tool handlers
	s.AddTool(listNotesTool, ListNotesHandler(ctx, cfg))
	s.AddTool(createNoteTool(), createNoteHandler(ctx, cfg))
}
