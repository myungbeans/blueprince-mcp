package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// MCPRequest represents a generic MCP request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

// MCPResponse represents a generic MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToolCallParams represents the parameters for calling a tool
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Client represents an MCP client for testing
type Client struct {
	configPath string
	verbose    bool
}

// NewClient creates a new MCP test client
func NewClient(cmd *cobra.Command) (*Client, error) {
	configPath, _ := cmd.Flags().GetString("config")
	verbose, _ := cmd.Flags().GetBool("verbose")

	// Make config path relative to current working directory if not absolute
	if !filepath.IsAbs(configPath) {
		// Try multiple potential paths to find the config
		possiblePaths := []string{
			configPath,                                    // As-is (when run from project root)
			filepath.Join("../..", configPath),          // When run from cmd/tools/
			filepath.Join("..", configPath),             // When run from cmd/
		}
		
		found := false
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				found = true
				break
			}
		}
		
		if !found {
			// Try to find project root by looking for go.mod
			cwd, err := os.Getwd()
			if err == nil {
				for dir := cwd; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
					if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
						// Found project root
						testPath := filepath.Join(dir, configPath)
						if _, err := os.Stat(testPath); err == nil {
							configPath = testPath
							found = true
							break
						}
					}
				}
			}
		}
	}

	// Verify config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", configPath)
	}

	if verbose {
		fmt.Printf("Using config: %s\n", configPath)
	}

	return &Client{
		configPath: configPath,
		verbose:    verbose,
	}, nil
}

// CallTool calls an MCP tool by running the server as a subprocess
func (c *Client) CallTool(toolName string, arguments map[string]interface{}) (*MCPResponse, error) {
	request := MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/call",
		Params: ToolCallParams{
			Name:      toolName,
			Arguments: arguments,
		},
	}

	if c.verbose {
		fmt.Printf("Calling tool: %s\n", toolName)
		if len(arguments) > 0 {
			argBytes, _ := json.MarshalIndent(arguments, "", "  ")
			fmt.Printf("Arguments:\n%s\n", string(argBytes))
		}
	}

	// Marshal request to JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Run the MCP server as a subprocess
	cmd := exec.Command("go", "run", "./cmd/server/main.go")
	
	// Set environment variable for config path
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+c.configPath)
	
	// Set up stdin/stdout pipes
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start server: %w", err)
	}
	defer cmd.Process.Kill()

	// Send the request
	if _, err := stdin.Write(reqBody); err != nil {
		return nil, fmt.Errorf("failed to write request: %w", err)
	}
	if _, err := stdin.Write([]byte("\n")); err != nil {
		return nil, fmt.Errorf("failed to write newline: %w", err)
	}
	stdin.Close()

	// Read the response
	scanner := bufio.NewScanner(stdout)
	var responseLine string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "{") {
			responseLine = line
			break
		}
	}

	if responseLine == "" {
		return nil, fmt.Errorf("no valid JSON response received")
	}

	if c.verbose {
		fmt.Printf("Response:\n%s\n", responseLine)
	}

	// Parse MCP response
	var mcpResp MCPResponse
	if err := json.Unmarshal([]byte(responseLine), &mcpResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &mcpResp, nil
}

// PrettyPrint prints the result in a human-readable format
func (c *Client) PrettyPrint(resp *MCPResponse) error {
	if resp.Error != nil {
		fmt.Printf("❌ Error: %s (code: %d)\n", resp.Error.Message, resp.Error.Code)
		if resp.Error.Data != nil {
			dataBytes, _ := json.MarshalIndent(resp.Error.Data, "", "  ")
			fmt.Printf("Details:\n%s\n", string(dataBytes))
		}
		return fmt.Errorf("MCP error: %s", resp.Error.Message)
	}

	if resp.Result != nil {
		// Try to extract content from result
		if resultMap, ok := resp.Result.(map[string]interface{}); ok {
			if content, ok := resultMap["content"]; ok {
				if contentSlice, ok := content.([]interface{}); ok && len(contentSlice) > 0 {
					if textContent, ok := contentSlice[0].(map[string]interface{}); ok {
						if text, ok := textContent["text"].(string); ok {
							fmt.Printf("✅ Success:\n%s\n", text)
							return nil
						}
					}
				}
			}
		}
		
		// Fallback to pretty print the whole result
		resultBytes, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Printf("✅ Success:\n%s\n", string(resultBytes))
	}

	return nil
}