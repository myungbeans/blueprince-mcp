package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "blueprince-tools",
		Short: "CLI tools for testing Blue Prince MCP server",
		Long: `A command line interface for testing the Blue Prince MCP server tools locally.
This CLI connects to the running MCP server and allows you to test create, read, update, 
and list operations on your notes.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Blue Prince MCP Tools CLI")
			fmt.Println("Use 'blueprince-tools --help' to see available commands")
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().String("config", "cmd/config/local/config.yaml", "Path to config file")
	rootCmd.PersistentFlags().String("host", "", "Override server host")
	rootCmd.PersistentFlags().Int("port", 0, "Override server port")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output")

	// Add subcommands
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newCreateCmd())
	rootCmd.AddCommand(newReadCmd())
	rootCmd.AddCommand(newUpdateCmd())
	rootCmd.AddCommand(newDeleteCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}