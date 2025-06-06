package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <path>",
		Short: "Read a specific note",
		Long:  `Reads the content of a specific note file using the read_note tool.`,
		Args:  cobra.ExactArgs(1),
		Example: `  # Read a specific note
  blueprince-tools read people/simon_jones.md
  
  # Read with verbose output
  blueprince-tools read rooms/library.md --verbose`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := NewClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to create client: %w", err)
			}

			notePath := args[0]

			// Call the read_note tool
			arguments := map[string]interface{}{
				"path": notePath,
			}

			resp, err := client.CallTool("read_note", arguments)
			if err != nil {
				return fmt.Errorf("failed to call read_note: %w", err)
			}

			return client.PrettyPrint(resp)
		},
	}

	return cmd
}