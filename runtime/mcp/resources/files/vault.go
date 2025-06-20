package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/myungbeans/blueprince-mcp/runtime/utils"

	"go.uber.org/zap"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterVault scans the rootDir's children and registers all valid files as MCP Resources.
func RegisterVault(ctx context.Context, s *server.MCPServer, rootDir string) error {
	logger := utils.Logger(ctx)

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

		mimeType := utils.GetMimeType(resourceName)

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
