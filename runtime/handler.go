package runtime

import (
	"context"

	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/resources/files"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/resources/rules"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/tools/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/mcp/tools/screenshots"
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
	s.AddTool(notes.ListTool(), notes.ListHandler(ctx, h.cfg))
	s.AddTool(notes.CreateTool(), notes.CreateHandler(ctx, h.cfg))
	s.AddTool(notes.ReadTool(), notes.ReadHandler(ctx, h.cfg))
	s.AddTool(notes.UpdateTool(), notes.UpdateHandler(ctx, h.cfg))
	s.AddTool(notes.DeleteTool(), notes.DeleteHandler(ctx, h.cfg))
	s.AddTool(screenshots.DownloadTool(), screenshots.DownloadHandler(ctx, h.store))
	s.AddTool(screenshots.ListTool(), screenshots.ListHandler(ctx, h.cfg, h.store))
	// TODO: need to figure out image compression s.AddTool(screenshots.ViewTool(), screenshots.ViewHandler(ctx, h.cfg))
	s.AddTool(screenshots.AnalyzeTool(), screenshots.AnalyzeHandler(ctx, h.cfg))
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
