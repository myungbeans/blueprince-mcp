package tools

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	"github.com/myungbeans/blueprince-mcp/runtime/models/storage"
// 	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
// 	"github.com/myungbeans/blueprince-mcp/runtime/utils"
// 	"go.uber.org/zap"
// 	"golang.org/x/oauth2"
// 	"golang.org/x/oauth2/google"
// 	"google.golang.org/api/drive/v3"
// 	"google.golang.org/api/option"
// )

// // ScreenshotProcessor handles Google Drive screenshot processing
// type ScreenshotProcessor struct {
// 	vaultPath      string
// 	storage        *storage.Store
// 	folderID       string
// 	screenshotsDir string
// }

// // DriveConfig represents the Google Drive configuration
// type DriveConfig struct {
// 	FolderID   string `json:"folder_id"`
// 	FolderName string `json:"folder_name"`
// 	TokenPath  string `json:"token_path"`
// }

// // NewScreenshotProcessor creates a new screenshot processor
// func NewScreenshotProcessor(ctx context.Context, vaultPath string) (*ScreenshotProcessor, error) {
// 	logger := utils.Logger(ctx)

// 	// Load Google Drive configuration
// 	driveConfig, err := loadDriveConfig()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to load Google Drive config: %w", err)
// 	}

// 	// Initialize Google Drive service
// 	driveService, err := initializeDriveService(driveConfig.TokenPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize Google Drive service: %w", err)
// 	}

// 	screenshotsDir := filepath.Join(vaultPath, vault.SCREENSHOT_DIR)

// 	return &ScreenshotProcessor{
// 		vaultPath:      vaultPath,
// 		driveService:   driveService,
// 		folderID:       driveConfig.FolderID,
// 		screenshotsDir: screenshotsDir,
// 		logger:         logger,
// 	}, nil
// }

// // ProcessScreenshots processes screenshots from Google Drive
// func (sp *ScreenshotProcessor) ProcessScreenshots(ctx context.Context, maxScreenshots int, forceReprocess bool) (*ProcessScreenshotsResponse, error) {
// 	response := &ProcessScreenshotsResponse{
// 		ProcessedFiles: make([]string, 0),
// 		Errors:         make([]string, 0),
// 	}

// 	// List screenshots in Google Drive folder
// 	screenshots, err := sp.listScreenshots(maxScreenshots)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list screenshots: %w", err)
// 	}

// 	sp.logger.Info("Found screenshots in Google Drive", zap.Int("count", len(screenshots)))

// 	for _, screenshot := range screenshots {
// 		sp.logger.Info("Processing screenshot", zap.String("name", screenshot.Name), zap.String("id", screenshot.Id))

// 		// Check if already processed (unless force reprocess)
// 		localPath := filepath.Join(sp.screenshotsDir, screenshot.Name)
// 		if !forceReprocess && fileExists(localPath) {
// 			sp.logger.Info("Screenshot already exists locally, skipping", zap.String("path", localPath))
// 			response.SkippedCount++
// 			continue
// 		}

// 		// Download screenshot
// 		if err := sp.downloadScreenshot(screenshot, localPath); err != nil {
// 			sp.logger.Error("Failed to download screenshot", zap.String("name", screenshot.Name), zap.Error(err))
// 			response.ErrorCount++
// 			response.Errors = append(response.Errors, fmt.Sprintf("Download failed for %s: %v", screenshot.Name, err))
// 			continue
// 		}

// 		// Analyze screenshot content
// 		content, err := sp.analyzeScreenshot(localPath)
// 		if err != nil {
// 			sp.logger.Error("Failed to analyze screenshot", zap.String("path", localPath), zap.Error(err))
// 			response.ErrorCount++
// 			response.Errors = append(response.Errors, fmt.Sprintf("Analysis failed for %s: %v", screenshot.Name, err))
// 			continue
// 		}

// 		// Create note from analysis
// 		notePath, err := sp.createNoteFromScreenshot(screenshot.Name, content, localPath)
// 		if err != nil {
// 			sp.logger.Error("Failed to create note from screenshot", zap.String("name", screenshot.Name), zap.Error(err))
// 			response.ErrorCount++
// 			response.Errors = append(response.Errors, fmt.Sprintf("Note creation failed for %s: %v", screenshot.Name, err))
// 			continue
// 		}

// 		response.ProcessedCount++
// 		response.ProcessedFiles = append(response.ProcessedFiles, notePath)
// 		sp.logger.Info("Successfully processed screenshot", zap.String("screenshot", screenshot.Name), zap.String("note", notePath))
// 	}

// 	return response, nil
// }

// // listScreenshots lists image files in the configured Google Drive folder
// func (sp *ScreenshotProcessor) listScreenshots(maxResults int) ([]*drive.File, error) {
// 	query := fmt.Sprintf("'%s' in parents and (mimeType contains 'image/' or name contains '.png' or name contains '.jpg' or name contains '.jpeg' or name contains '.gif' or name contains '.bmp' or name contains '.webp') and trashed=false", sp.folderID)

// 	call := sp.driveService.Files.List().Q(query).PageSize(int64(maxResults)).OrderBy("createdTime desc")

// 	result, err := call.Do()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list files from Google Drive: %w", err)
// 	}

// 	return result.Files, nil
// }

// // downloadScreenshot downloads a screenshot from Google Drive to local storage
// func (sp *ScreenshotProcessor) downloadScreenshot(file *drive.File, localPath string) error {
// 	// Ensure directory exists
// 	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
// 		return fmt.Errorf("failed to create directory: %w", err)
// 	}

// 	// Download file content
// 	response, err := sp.driveService.Files.Get(file.Id).Download()
// 	if err != nil {
// 		return fmt.Errorf("failed to download file: %w", err)
// 	}
// 	defer response.Body.Close()

// 	// Create local file
// 	localFile, err := os.Create(localPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create local file: %w", err)
// 	}
// 	defer localFile.Close()

// 	// Copy content
// 	_, err = io.Copy(localFile, response.Body)
// 	if err != nil {
// 		return fmt.Errorf("failed to copy file content: %w", err)
// 	}

// 	sp.logger.Info("Downloaded screenshot", zap.String("source", file.Name), zap.String("destination", localPath))
// 	return nil
// }

// // analyzeScreenshot analyzes the content of a screenshot image
// func (sp *ScreenshotProcessor) analyzeScreenshot(imagePath string) (string, error) {
// 	// For now, this is a placeholder that returns basic information
// 	// In a real implementation, this would use an image analysis service
// 	// or AI vision API to describe the screenshot content

// 	fileInfo, err := os.Stat(imagePath)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get file info: %w", err)
// 	}

// 	filename := filepath.Base(imagePath)
// 	timestamp := fileInfo.ModTime().Format("2006-01-02 15:04:05")

// 	// Basic analysis - in a real implementation, this would be replaced with
// 	// actual image content analysis using vision AI services
// 	content := fmt.Sprintf(`# Screenshot Analysis: %s

// ## Basic Information
// - **Filename**: %s
// - **Timestamp**: %s
// - **File Size**: %d bytes

// ## Visual Content Description
// *This section would contain AI-generated description of the screenshot content*

// **Note**: This is a placeholder. In the full implementation, this would contain:
// - Description of UI elements visible in the screenshot
// - Text content that can be read from the image
// - Visual elements like characters, objects, environments
// - Game state information visible in the screenshot
// - Any other relevant visual information

// ## Analysis Notes
// - Screenshot captured from Blue Prince game
// - Content analysis focuses only on what is directly visible
// - No speculation or external game knowledge added
// - Spoiler-free description of visual elements only`, filename, filename, timestamp, fileInfo.Size())

// 	return content, nil
// }

// // createNoteFromScreenshot creates a spoiler-free note from screenshot analysis
// func (sp *ScreenshotProcessor) createNoteFromScreenshot(screenshotName, content, imagePath string) (string, error) {
// 	// Generate note filename based on screenshot name and timestamp
// 	timestamp := time.Now().Format("2006-01-02-150405")
// 	baseName := strings.TrimSuffix(screenshotName, filepath.Ext(screenshotName))
// 	noteName := fmt.Sprintf("screenshot-%s-%s.md", baseName, timestamp)

// 	// Create note in the general category (screenshots are general observations)
// 	noteCategory := "general"
// 	notesDir := filepath.Join(sp.vaultPath, vault.NOTES_DIR, noteCategory)
// 	notePath := filepath.Join(notesDir, noteName)

// 	// Ensure directory exists
// 	if err := os.MkdirAll(notesDir, 0755); err != nil {
// 		return "", fmt.Errorf("failed to create notes directory: %w", err)
// 	}

// 	// Create note content with metadata
// 	noteContent := fmt.Sprintf(`---
// title: "Screenshot: %s"
// category: %s
// tags: ["screenshot", "visual", "game-capture"]
// screenshot_file: "%s"
// screenshot_path: "%s"
// processed_date: "%s"
// source: "google-drive-import"
// ---

// %s`, baseName, noteCategory, filepath.Base(imagePath), imagePath, time.Now().Format(time.RFC3339), content)

// 	if err := os.WriteFile(notePath, []byte(noteContent), 0644); err != nil {
// 		return "", fmt.Errorf("failed to write note file: %w", err)
// 	}

// 	sp.logger.Info("Created note from screenshot", zap.String("note", notePath), zap.String("screenshot", screenshotName))
// 	return notePath, nil
// }

// // fileExists checks if a file exists
// func fileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	return err == nil
// }
