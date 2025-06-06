package resources

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register scans the rootDir's children and registers all found files as MCP Resources.
func Register(ctx context.Context, s *server.MCPServer, rootDir string) error {
	logger := utils.Logger(ctx)

	// First, register the spoiler prevention rules as a special resource
	if err := registerSpoilerPreventionRules(ctx, s); err != nil {
		logger.Error("Failed to register spoiler prevention rules", zap.Error(err))
		return err
	}

	absRootDir, err := utils.ResolveAndCleanPath(rootDir)
	if err != nil {
		return err
	}

	logger.Info("Scanning for meta file resources in " + absRootDir)
	// TODO: refactor to use common util for traversing files?
	err = filepath.WalkDir(absRootDir, func(fullFilePath string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			logger.Error("Error accessing path during resource scan",
				zap.String("path", fullFilePath),
				zap.Error(walkErr),
			)
			// Skip problematic files, only error on dir errors to prevent further attempts
			if d == nil || !d.IsDir() {
				return nil
			}
			return walkErr
		}

		// Skip unwanted paths (.obsidian, hidden files, etc.)
		if utils.ShouldSkipPath(fullFilePath, d) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Process only files
		if d.IsDir() {
			return nil
		}

		relativePath, err := filepath.Rel(absRootDir, fullFilePath)
		if err != nil {
			logger.Error("Could not get relative path",
				zap.String("fullFilePath", fullFilePath),
				zap.String("absRootDir", absRootDir),
				zap.Error(err),
			)
			return nil // Skip this file
		}
		relativePath = filepath.ToSlash(relativePath) // Ensure URI uses forward slashes

		resourceURI := "file:///" + relativePath
		resourceName := d.Name()
		resourceDescription := "Meta File resource: " + relativePath

		mimeType := getMimeType(resourceName)

		resource := mcp.NewResource(
			resourceURI,
			resourceName,
			mcp.WithResourceDescription(resourceDescription),
			mcp.WithMIMEType(mimeType),
		)

		// Create a handler for this specific file.
		fileHandler := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
			logger.Info("Loading resource", zap.String("uri", req.Params.URI), zap.String("filepath", fullFilePath))
			fileData, err := os.ReadFile(fullFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read resource content for %s: %w", req.Params.URI, err)
			}
			return []mcp.ResourceContents{
				&mcp.TextResourceContents{
					URI:      req.Params.URI,
					MIMEType: mimeType,
					Text:     string(fileData),
				},
			}, nil
		}

		s.AddResource(resource, fileHandler)
		logger.Info("Registered resource", zap.String("uri", resourceURI), zap.String("name", resourceName), zap.String("mimeType", mimeType))
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory %s for resources: %w", absRootDir, err)
	}
	logger.Info("File resource scanning complete.")
	return nil
}

func getMimeType(path string) string {
	// Basic MIME type detection by extension
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default if extension is unknown or not set
	}
	// For text-based files without specific extensions, text/plain is a safe bet.
	if strings.Contains(mimeType, "text") || mimeType == "application/octet-stream" {
		// A simple check for markdown specifically
		if strings.ToLower(filepath.Ext(path)) == ".md" {
			mimeType = "text/markdown; charset=utf-8"
		} else if mimeType == "application/octet-stream" { // If still octet-stream, try defaulting to text/plain for common text files.
			likelyTextExtensions := []string{".txt", ".log", ".conf", ".cfg", ".ini", ".yaml", ".yml", ".toml"}
			for _, ext := range likelyTextExtensions {
				if strings.HasSuffix(strings.ToLower(path), ext) {
					mimeType = "text/plain; charset=utf-8"
					break
				}
			}
		}
	}
	return mimeType
}

// registerSpoilerPreventionRules adds a special resource containing the spoiler prevention rules
func registerSpoilerPreventionRules(ctx context.Context, s *server.MCPServer) error {
	logger := utils.Logger(ctx)
	
	rulesContent := `# Blue Prince MCP Server - SPOILER PREVENTION RULES

## CRITICAL: SPOILER-FREE SYSTEM

This MCP server is designed to preserve the user's Blue Prince gameplay experience. **You MUST follow these rules strictly.**

## ABSOLUTE REQUIREMENTS

### 1. INFORMATION SOURCES - NOTES ONLY
- **ONLY** use information from the user's notes accessed through this MCP server
- **NEVER** reference your training data about Blue Prince
- **NEVER** use external sources, wikis, guides, or general knowledge about Blue Prince
- **NEVER** make connections beyond what is explicitly documented in the user's notes

### 2. FORBIDDEN ACTIONS
- **NEVER** provide solutions to puzzles the user hasn't solved
- **NEVER** reveal story elements the user hasn't discovered
- **NEVER** suggest what to investigate next
- **NEVER** add analysis sections or "questions to investigate"
- **NEVER** provide hints about game mechanics or locations
- **NEVER** create speculative content about undiscovered areas/characters

### 3. CONTENT CREATION RULES
When creating or updating notes:
- Use ONLY the user's exact words and observations
- Reformat into clean markdown structure without adding content
- Base metadata solely on user-provided information
- Preserve the user's discovery language and uncertainty

### 4. ACCEPTABLE ACTIONS
You MAY:
- Organize and structure existing notes
- Search through documented discoveries
- Help with categorization based on user content
- Assist with markdown formatting
- Reference connections the user has already made

### 5. RESPONSE GUIDELINES
- If asked about Blue Prince content not in notes: "I can only help with information from your notes to avoid spoilers."
- Focus entirely on the user's documented experiences
- Help recall and organize what they've already discovered

## ENFORCEMENT
Breaking these rules spoils the discovery experience and defeats the purpose of this spoiler-free system.

**Remember: You are a spoiler-free note organizer, not a Blue Prince guide.**`

	rulesResource := mcp.NewResource(
		"rules://blue-prince/spoiler-protection",
		"Blue Prince Spoiler Protection Rules",
		mcp.WithResourceDescription("CRITICAL rules for spoiler-free Blue Prince assistance that the assistant must follow at all times"),
		mcp.WithMIMEType("text/markdown; charset=utf-8"),
	)

	rulesHandler := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		logger.Info("Loading spoiler prevention rules", zap.String("uri", req.Params.URI))
		return []mcp.ResourceContents{
			&mcp.TextResourceContents{
				URI:      req.Params.URI,
				MIMEType: "text/markdown; charset=utf-8",
				Text:     rulesContent,
			},
		}, nil
	}

	s.AddResource(rulesResource, rulesHandler)
	logger.Info("Registered spoiler prevention rules resource", 
		zap.String("uri", "rules://blue-prince/spoiler-protection"))
	
	return nil
}
