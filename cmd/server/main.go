package main

import (
	"context"
	"fmt"
	"os"

	"github.com/myungbeans/blueprince-mcp/cmd/server/config"
	"github.com/myungbeans/blueprince-mcp/runtime/handlers"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Assuming config.yaml is at the project root and the binary is run from the project root.
	// e.g., go run ./cmd/server/main.go

	// Initialize logger early
	logger, err := zap.NewDevelopment() // Use NewDevelopment for more human-friendly output during dev
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync() // Flushes buffer, if any

	// Validate that the server is run from the project root by checking for a sentinel file.
	// config.yaml is a good candidate.
	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		logger.Fatal("Configuration file not found. Please run the server from the project root directory.", zap.String("configFile", "config.yaml"), zap.Error(err))
	}

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		return
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Blue Prince Architect Notes",
		"0.0.0",
		server.WithToolCapabilities(false),
	)

	// TODO: remove generated hello_world
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	// Register Tools
	// createNoteTool := mcp.NewTool("create_note",
	// 	mcp.WithDescription("Creates a new note."),
	// 	mcp.WithString("title",
	// 		mcp.Required(),
	// 		mcp.Description("The title of the note."),
	// 	),
	// 	mcp.WithString("content",
	// 		mcp.Required(),
	// 		mcp.Description("The content of the note."),
	// 	),
	// )

	listNotesTool := mcp.NewTool("list_notes",
		mcp.WithDescription("Lists all notes."),
		// Add parameters for pagination/filtering later if needed
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
	s.AddTool(tool, helloHandler)
	s.AddTool(listNotesTool, handlers.ListNotesHandler(cfg, logger)) // Pass logger to the factory
	// s.AddTool(createNoteTool, createNoteHandler)
	// s.AddTool(readNoteTool, readNoteHandler)
	// s.AddTool(updateNoteTool, updateNoteHandler)
	// s.AddTool(deleteNoteTool, deleteNoteHandler)
	// s.AddTool(searchNotesTool, searchNotesHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
