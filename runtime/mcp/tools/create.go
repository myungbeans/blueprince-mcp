package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/myungbeans/blueprince-mcp/cmd/config"
	"github.com/myungbeans/blueprince-mcp/runtime/models/notes"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	// "github.com/mark3labs/mcp-go/server"
	// "github.com/myungbeans/blueprince-mcp/cmd/config"
	// "go.uber.org/zap"
)

// createNoteTool returns the configured mcp.Tool for creating notes
func createNoteTool() mcp.Tool {
	createNoteTool := mcp.Tool{
		Name:        "create_note",
		Description: "Creates a structured Blue Prince note with intelligent organization based on user input. Relies on MCP Client to intelligently extract some parameters for the request, as well as format the request appropriately. VERY IMPORTANT: the MCP Client should avoid any modifications that could be spoilers for the game -- rely heavily on what the user inputted and only on information available from existing notes",
		InputSchema: notes.GetMCPSchema(),
	}

	// Add detailed examples
	createNoteTool.Description += `

EXAMPLES:

User Input: "room: nook, paintings of tiger and a cupcake stand? Weird"
Expected Call:
{
  "path": "rooms/nook_tiger_paintings.md",
  "metadata": {
    "title": "Nook - Tiger Paintings and Cupcake Stand",
    "category": "rooms",
    "primary_subject": "paintings",
    "tags": ["rooms", "nook", "paintings", "tiger", "cupcake_stand", "weird_elements"],
    "confidence": "medium",
    "status": "needs_investigation"
  },
  "content": "# Nook - Tiger Paintings and Cupcake Stand\n\n## Initial Observation\nPaintings of tiger and a cupcake stand? Weird\n\n## Related Elements\n*Connections will be added as more notes are discovered*"
}

User Input: "Simon is 14, won science fair runner-up, inherits mansion"
Expected Call:
{
  "path": "people/simon_jones.md",
  "metadata": {
    "title": "Simon P. Jones - Protagonist",
    "category": "person",
    "primary_subject": "simon_jones",
    "tags": ["people", "simon_jones", "protagonist", "14_years_old", "science_fair", "inheritance"],
    "confidence": "high",
    "status": "confirmed"
  },
  "content": "# Simon P. Jones - Protagonist\n\n## Basic Info\n- Age: 14\n- Science fair runner-up\n- Inherits Mt. Holly Estate\n\n## Relationships\n- Relationships will be added as more notes are discovered"
}`

	return createNoteTool
}

func createNoteHandler(ctx context.Context, cfg *config.Config) server.ToolHandlerFunc {
	logger := utils.Logger(ctx)
	// TODO: START IMPLEMENTATION FROM HERE
	return nil
}
