package notes

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"gopkg.in/yaml.v3"
)

// Metadata represents the YAML frontmatter structure
type Metadata struct {
	Title          string   `json:"title" yaml:"title"`
	Category       string   `json:"category" yaml:"category"`
	PrimarySubject string   `json:"primary_subject" yaml:"primary_subject"`
	Tags           []string `json:"tags" yaml:"tags"`
	Confidence     string   `json:"confidence" yaml:"confidence"`
	Status         string   `json:"status" yaml:"status"`
	CreatedAt      string   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
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
				"description": "Raw markdown content containing ONLY the user's observations and input. NEVER add investigation prompts, analysis questions, or additional sections not explicitly provided by the user. Preserve exactly what the player observed without speculation or gameplay hints.",
			},
		},
		Required: []string{"path", "metadata", "content"},
	}
}

// ParseMetadata converts map[string]any to Metadata struct
func ParseMetadata(metadataMap map[string]any) (*Metadata, error) {
	metadata := &Metadata{}

	// Parse title
	title, ok := metadataMap["title"].(string)
	if !ok {
		return nil, fmt.Errorf("title must be a string")

	}
	metadata.Title = title

	// Parse category
	category, ok := metadataMap["category"].(string)
	if !ok {
		return nil, fmt.Errorf("category must be a string")
	}
	metadata.Category = category

	if !isValidCategory(metadata.Category) {
		return nil, fmt.Errorf("Invalid category '%s'. Must be one of: %v", metadata.Category, Categories)
	}

	// Parse primary_subject (not stored in metadata struct but validated)
	primarySubject, ok := metadataMap["category"].(string)
	if !ok {
		return nil, fmt.Errorf("primary subject must be a string")
	}
	metadata.PrimarySubject = primarySubject

	// Parse tags
	rawVal, ok := metadataMap["tags"]
	if !ok {
		return nil, fmt.Errorf("tags is required")
	}
	rawTags, ok := rawVal.([]any)
	if !ok {
		return nil, fmt.Errorf("tags must be an array")
	}
	metadata.Tags = make([]string, len(rawTags))
	for i, tag := range rawTags {
		tagStr, ok := tag.(string)
		if !ok {
			return nil, fmt.Errorf("all tags must be strings")
		}
		metadata.Tags[i] = tagStr
	}

	// Parse confidence
	confidence, ok := metadataMap["confidence"].(string)
	if !ok {
		return nil, fmt.Errorf("confidence must be a string")
	}

	if !isValidConfidence(confidence) {
		return nil, fmt.Errorf("confidence must be one of: high, medium, low")
	}
	metadata.Confidence = confidence

	// Parse status
	status, ok := metadataMap["status"].(string)
	if !ok {
		return nil, fmt.Errorf("status must be a string")
	}

	if !isValidStatus(status) {
		return nil, fmt.Errorf("status must be one of: complete, needs_investigation, active_investigation, theory, confirmed")
	}
	metadata.Status = status

	return metadata, nil
}

// isValidCategory checks if category is in the allowed list
func isValidCategory(category string) bool {
	for _, validCategory := range Categories {
		if category == validCategory {
			return true
		}
	}
	return false
}

// isValidConfidence checks if confidence level is valid
func isValidConfidence(confidence string) bool {
	for _, valid := range ConfidenceLevels {
		if confidence == valid {
			return true
		}
	}
	return false
}

// isValidStatus checks if status is valid
func isValidStatus(status string) bool {
	for _, valid := range Statuses {
		if status == valid {
			return true
		}
	}
	return false
}

// CreateContent generates the full file content with YAML frontmatter
func CreateContent(metadata *Metadata, content string) (string, error) {
	// Marshal metadata to YAML
	yamlBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata to YAML: %w", err)
	}

	// Create file content with frontmatter
	fileContent := fmt.Sprintf("---\n%s---\n\n%s", string(yamlBytes), content)
	return fileContent, nil
}
