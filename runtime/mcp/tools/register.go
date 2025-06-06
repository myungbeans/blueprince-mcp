package tools

import (
	"context"

	"github.com/myungbeans/blueprince-mcp/cmd/config"

	"github.com/mark3labs/mcp-go/server"
)

func Register(ctx context.Context, s *server.MCPServer, cfg *config.Config) {
	s.AddTool(listNotesTool(), listNotesHandler(ctx, cfg))
	s.AddTool(createNoteTool(), createNoteHandler(ctx, cfg))
	s.AddTool(readNoteTool(), readNoteHandler(ctx, cfg))
	s.AddTool(updateNoteTool(), updateNoteHandler(ctx, cfg))
	s.AddTool(deleteNoteTool(), deleteNoteHandler(ctx, cfg))
}
