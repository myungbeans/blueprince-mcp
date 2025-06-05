package main

import (
	"fmt"
	"os"

	"github.com/myungbeans/blueprince-mcp/cmd/server/config"
	"github.com/myungbeans/blueprince-mcp/runtime/handlers"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	defaultConfigFilePath = "cmd/server/config/local/config.yaml"
	envVaultPath          = "OBSIDIAN_VAULT_PATH"
)

func main() {
	// Initialize logger early
	logger, err := zap.NewDevelopment() // Use NewDevelopment for more human-friendly output during dev
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flushes buffer, if any

	var cfg *config.Config
	vaultPath := os.Getenv(envVaultPath)
	if vaultPath != "" {
		logger.Info("Using Obsidian vault path from env variable", zap.String("variable", envVaultPath), zap.String("path", vaultPath))
		cfg = &config.Config{
			ObsidianVaultPath: vaultPath,
		}
	} else {
		cfg, err = config.LoadConfig(defaultConfigFilePath)
		if err != nil {
			logger.Fatal("Failed to load configuration", zap.Error(err))
			os.Exit(1)
		}
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Blue Prince Architect Notes",
		"0.0.0",
		server.WithToolCapabilities(false),
	)

	// Register Tools
	registerTools(s, cfg, logger)

	// Start the stdio server
	logger.Info("Initializing server...")
	if err := server.ServeStdio(s); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

func registerTools(s *server.MCPServer, cfg *config.Config, logger *zap.Logger) {
	createNoteTool := mcp.NewTool("create_note",
		mcp.WithDescription("Creates a new note."),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the note."),
		),
		mcp.WithString("tags",
			mcp.Required(),
			mcp.Description("Tags that describe some of the topics included in the note"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("The content of the note."),
		),
	)

	listNotesTool := mcp.NewTool("list_notes",
		mcp.WithDescription("Lists all notes."),
		// TODO: Add parameters for pagination/filtering later if needed
	)

	// readNoteTool := mcp.NewTool("read_note",
	// 	mcp.WithDescription("Reads a specific note by its ID."),
	// 	mcp.WithString("id",
	// 		mcp.Required(),
	// 		mcp.Description("The unique identifier of the note."),
	// 	),
	// )

	// updateNoteTool := mcp.NewTool("update_note",
	// 	mcp.WithDescription("Updates an existing note by its ID."),
	// 	mcp.WithString("id",
	// 		mcp.Required(),
	// 		mcp.Description("The unique identifier of the note to update."),
	// 	),
	// 	mcp.WithString("title",
	// 		mcp.Description("The new title of the note (optional)."),
	// 	),
	// 	mcp.WithString("content",
	// 		mcp.Description("The new content of the note (optional)."),
	// 	),
	// 	// Note: Need logic in handler to ensure at least title or content is provided
	// )

	// deleteNoteTool := mcp.NewTool("delete_note",
	// 	mcp.WithDescription("Deletes a specific note by its ID."),
	// 	mcp.WithString("id",
	// 		mcp.Required(),
	// 		mcp.Description("The unique identifier of the note to delete."),
	// 	),
	// )

	// searchNotesTool := mcp.NewTool("search_notes",
	// 	mcp.WithDescription("Searches notes based on keywords."),
	// 	mcp.WithString("query", mcp.Required(), mcp.Description("The search query or keywords.")),
	// )

	// Add Note tool handlers
	s.AddTool(listNotesTool, handlers.ListNotesHandler(cfg, logger))   // Pass logger to the factory
	s.AddTool(createNoteTool, handlers.CreateNoteHandler(cfg, logger)) // Pass logger to the factory
	// s.AddTool(readNoteTool, readNoteHandler)
	// s.AddTool(updateNoteTool, updateNoteHandler)
	// s.AddTool(deleteNoteTool, deleteNoteHandler)
	// s.AddTool(searchNotesTool, searchNotesHandler)
}
