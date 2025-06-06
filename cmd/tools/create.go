package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
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
		Use:   "create <path>",
		Short: "Create a new note",
		Long:  `Creates a new note using the create_note tool with the specified metadata and content.`,
		Args:  cobra.ExactArgs(1),
		Example: `  # Create a simple note
  blueprince-tools create people/simon.md --title "Simon Jones" --category "people" --content "14 years old, science fair runner-up"
  
  # Create with full metadata
  blueprince-tools create rooms/library.md \
    --title "Library" \
    --category "rooms" \
    --primary-subject "library" \
    --tags "rooms,books,quiet" \
    --confidence "high" \
    --status "confirmed" \
    --content "Large room with many books and a reading area"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			notePath := args[0]

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

			if content == "" {
				content = fmt.Sprintf("# %s\n\nContent for %s", title, title)
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

			// Call the create_note tool
			arguments := map[string]interface{}{
				"path":     notePath,
				"metadata": metadata,
				"content":  content,
			}

			resp, err := client.CallTool("create_note", arguments)
			if err != nil {
				return fmt.Errorf("failed to call create_note: %w", err)
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
	cmd.Flags().StringVar(&content, "content", "", "Note content (markdown)")

	return cmd
}