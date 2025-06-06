package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	var (
		title          string
		category       string
		primarySubject string
		tags           []string
		confidence     string
		status         string
		content        string
	)

	cmd := &cobra.Command{
		Use:   "update <path>",
		Short: "Update an existing note",
		Long: `Updates an existing note using the update_note tool. 
Note: This completely replaces the note content. The MCP client should handle merging logic.`,
		Args: cobra.ExactArgs(1),
		Example: `  # Update note content
  blueprince-tools update people/simon.md --content "14 years old, science fair runner-up, found a key"
  
  # Update with new metadata
  blueprince-tools update rooms/library.md \
    --title "Grand Library" \
    --tags "rooms,books,quiet,secret_passage" \
    --status "confirmed" \
    --content "Large room with books. Found secret passage behind bookshelf."`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			notePath := args[0]

			// For update, we need to provide complete metadata and content
			// In a real scenario, the MCP client would read the existing note first
			// and merge the changes, but for testing we'll require explicit values

			if content == "" {
				return fmt.Errorf("content is required for update. Use --content flag")
			}

			// Set defaults if not provided
			if title == "" {
				// Extract title from path
				pathParts := strings.Split(notePath, "/")
				if len(pathParts) > 0 {
					filename := pathParts[len(pathParts)-1]
					title = strings.TrimSuffix(filename, ".md")
					title = strings.ReplaceAll(title, "_", " ")
					title = strings.Title(strings.ToLower(title))
				}
			}

			if category == "" {
				// Extract category from path
				pathParts := strings.Split(notePath, "/")
				if len(pathParts) > 1 {
					category = pathParts[0]
				} else {
					category = "general"
				}
			}

			if primarySubject == "" {
				// Use filename as primary subject
				pathParts := strings.Split(notePath, "/")
				if len(pathParts) > 0 {
					filename := pathParts[len(pathParts)-1]
					primarySubject = strings.TrimSuffix(filename, ".md")
				}
			}

			if len(tags) == 0 {
				// Default tags from category and primary subject
				tags = []string{category, primarySubject}
			}

			if confidence == "" {
				confidence = "medium"
			}

			if status == "" {
				status = "needs_investigation"
			}

			// Build metadata
			metadata := map[string]interface{}{
				"title":           title,
				"category":        category,
				"primary_subject": primarySubject,
				"tags":            tags,
				"confidence":      confidence,
				"status":          status,
			}

			// Call the update_note tool
			arguments := map[string]interface{}{
				"path":     notePath,
				"metadata": metadata,
				"content":  content,
			}

			resp, err := client.CallTool("update_note", arguments)
			if err != nil {
				return fmt.Errorf("failed to call update_note: %w", err)
			}

			return client.PrettyPrint(resp)
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&title, "title", "t", "", "Note title")
	cmd.Flags().StringVarP(&category, "category", "c", "", "Note category (people, rooms, items, lore, puzzles, general)")
	cmd.Flags().StringVarP(&primarySubject, "primary-subject", "s", "", "Primary subject")
	cmd.Flags().StringSliceVar(&tags, "tags", nil, "Comma-separated list of tags")
	cmd.Flags().StringVar(&confidence, "confidence", "", "Confidence level (high, medium, low)")
	cmd.Flags().StringVar(&status, "status", "", "Status (complete, needs_investigation, active_investigation, theory, confirmed)")
	cmd.Flags().StringVar(&content, "content", "", "Note content (markdown) - REQUIRED")

	// Mark content as required
	cmd.MarkFlagRequired("content")

	return cmd
}