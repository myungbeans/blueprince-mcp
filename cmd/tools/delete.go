package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <path>",
		Short: "Delete a note",
		Long:  `Deletes a note file using the delete_note tool. The note will be permanently removed from the vault.`,
		Args:  cobra.ExactArgs(1),
		Example: `  # Delete a note
  blueprince-tools delete people/simon.md
  
  # Delete a room note
  blueprince-tools delete rooms/library.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			notePath := args[0]

			// Call the delete_note tool
			arguments := map[string]interface{}{
				"path": notePath,
			}

			resp, err := client.CallTool("delete_note", arguments)
			if err != nil {
				return fmt.Errorf("failed to call delete_note: %w", err)
			}

			return client.PrettyPrint(resp)
		},
	}

	return cmd
}