package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all notes in the vault",
		Long:  `Lists all notes in the Blue Prince vault using the list_notes tool.`,
		Example: `  # List all notes
  blueprince-tools list
  
  # List with verbose output
  blueprince-tools list --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			// Call the list_notes tool with no arguments
			resp, err := client.CallTool("list_notes", map[string]interface{}{})
			if err != nil {
				return fmt.Errorf("failed to call list_notes: %w", err)
			}

			return client.PrettyPrint(resp)
		},
	}

	return cmd
}