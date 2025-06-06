package notes

import (
	"github.com/mark3labs/mcp-go/mcp"
)

// Metadata represents the YAML frontmatter structure
type Metadata struct {
	Title      string   `json:"title" yaml:"title"`
	Category   string   `json:"category" yaml:"category"`
	Tags       []string `json:"tags" yaml:"tags"`
	Confidence string   `json:"confidence" yaml:"confidence"`
	Status     string   `json:"status" yaml:"status"`
	CreatedAt  string   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt  string   `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

func GetMCPSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]any{
			"path": map[string]string{
				"type":        "string",
				"description": "MCP Client generated file path following pattern: {category}s/{subject}_{keywords}.md (e.g., 'rooms/nook_tiger_paintings.md')",
			},
			"metadata": map[string]any{
				"type":        "object",
				"description": "YAML frontmatter metadata extracted from user input",
				"properties": map[string]any{
					"title": map[string]string{
						"type":        "string",
						"description": "Human-readable descriptive title",
					},
					"category": map[string]any{
						"type":        "string",
						"description": "Content type category - either user inputted or intelligently extracted from user input by MCP client",
						"enum":        Categories,
					},
					"primary_subject": map[string]string{
						"type":        "string",
						"description": "Primary subject intelligently extracted from user input by MCP client (e.g., 'nook', 'simon_jones', 'coat_of_arms')",
					},
					"tags": map[string]any{
						"type":        "array",
						"description": "Searchable tags useful lookups: [type, subject, ...key_elements, ...descriptive_terms]",
						"items": map[string]string{
							"type": "string",
						},
					},
					"confidence": map[string]any{
						"type":        "string",
						"description": "Information reliability based on user certainty",
						"enum":        []string{"high", "medium", "low"},
					},
					"status": map[string]any{
						"type":        "string",
						"description": "Investigation status",
						"enum":        []string{"complete", "needs_investigation", "active_investigation", "theory", "confirmed"},
					},
				},
				"required": []string{"title", "category", "primary_subject", "tags", "confidence", "status"},
			},
			"content": map[string]string{
				"type":        "string",
				"description": "Markdown content following Blue Prince note templates with structured sections",
			},
		},
		Required: []string{"path", "metadata", "content"},
	}
}
