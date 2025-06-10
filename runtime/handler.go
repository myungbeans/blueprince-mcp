package runtime

import (
	"context"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/resources/files"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/resources/rules"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/tools"
	"github.com/myungbeans/blueprince-mcp/runtime/models/storage"

	"github.com/mark3labs/mcp-go/server"
)

type Handler struct {
	cfg   *config.Config
	store storage.Store
}

func NewHandler(cfg *config.Config, store storage.Store) *Handler {
	return &Handler{
		cfg:   cfg,
		store: store,
	}
}

func (h *Handler) RegisterTools(ctx context.Context, s *server.MCPServer) {
	// Register Tools
	s.AddTool(tools.ListNotesTool(), tools.ListNotesHandler(ctx, h.cfg))
	s.AddTool(tools.CreateNoteTool(), tools.CreateNoteHandler(ctx, h.cfg))
	s.AddTool(tools.ReadNoteTool(), tools.ReadNoteHandler(ctx, h.cfg))
	s.AddTool(tools.UpdateNoteTool(), tools.UpdateNoteHandler(ctx, h.cfg))
	s.AddTool(tools.DeleteNoteTool(), tools.DeleteNoteHandler(ctx, h.cfg))
	// s.AddTool(tools.GetScreenshotsTool(h.cfg.ObsidianVaultPath), tools.ProcessScreenshotsHandler(ctx, h.cfg))
}

func (h *Handler) RegisterResources(ctx context.Context, s *server.MCPServer) error {
	if err := rules.RegisterSpoilerRules(ctx, s); err != nil {
		return err
	}

	if err := files.RegisterVault(ctx, s, h.cfg.ObsidianVaultPath); err != nil {
		return err
	}

	return nil
}
