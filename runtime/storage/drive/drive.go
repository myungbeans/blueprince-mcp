package drive

import (
	"fmt"
	"io"
	"os"

	"github.com/myungbeans/blueprince-mcp/runtime/models/vault"
	"github.com/myungbeans/blueprince-mcp/runtime/utils"
	"google.golang.org/api/drive/v3"
)

// GetFile downloads a file from Google Drive to local storage
func (g *GoogleDrive) GetFiles(filename string) ([]string, error) {
	query := fmt.Sprintf("'%s' in parents and trashed=false", g.FolderID)
	if filename != "" {
		// Find the file in Google Drive folder
		query = fmt.Sprintf("name='%s' and %s", filename, query)
	}

	result, err := g.Client.Files.List().Q(query).PageSize(max_page_size).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search for file '%s': %w", filename, err)
	}

	if len(result.Files) == 0 {
		return nil, fmt.Errorf("file '%s' not found in Google Drive folder", filename)
	}

	// Ensure local dest is ready
	if err := utils.EnsureDirExists(g.ScreenshotsDir, 0755); err != nil {
		return nil, err
	}

	files := make([]string, len(result.Files))
	for i, file := range result.Files {
		// Download file content
		response, err := g.Client.Files.Get(file.Id).Download()
		if err != nil {
			return nil, fmt.Errorf("failed to download file '%s': %w", filename, err)
		}
		defer response.Body.Close()

		// Create local file
		// This needs to be the vaultpath + screenshotsdir +
		fullPath, err := utils.BuildSecurePath(g.VaultPath, vault.SCREENSHOT_DIR, filename)
		if err != nil {
			return nil, fmt.Errorf("Security validation failed for img path %q: %w", filename, err)
		}

		newImgFile, err := os.Create(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create local file '%s': %w", fullPath, err)
		}
		defer newImgFile.Close()

		// Copy content
		_, err = io.Copy(newImgFile, response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file content: %w", err)
		}

		files[i] = filename

		g.MoveFile(file.Id, archive_dir)
	}

	return files, nil
}

// ListFiles lists files in the Google Drive folder
func (g *GoogleDrive) ListFiles() ([]string, error) {
	// Build query to list files in the configured folder
	query := fmt.Sprintf("'%s' in parents and trashed=false", g.FolderID)

	// Add filter to exclude folders from results (only return files)
	query += " and mimeType != 'application/vnd.google-apps.folder'"
	result, err := g.Client.Files.List().Q(query).PageSize(max_page_size).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files in Google Drive: %w", err)
	}

	filenames := make([]string, len(result.Files))
	for i, file := range result.Files {
		filenames[i] = file.Name
	}

	return filenames, nil
}

// MoveFile moves a file from current location to a new directory within Google Drive
// filename: name of the file to move
// destination: destination directory name (will be created if it doesn't exist)
func (g *GoogleDrive) MoveFile(fileId, destination string) error {
	// Find or create destination folder
	destFolderID, err := g.findOrCreateSubfolder(destination)
	if err != nil {
		return fmt.Errorf("failed to find or create destination folder '%s': %w", destination, err)
	}

	// Move the file by updating its parents
	// Remove from current parent and add to new parent
	_, err = g.Client.Files.Update(fileId, &drive.File{}).
		AddParents(destFolderID).
		RemoveParents(g.FolderID).
		Do()

	if err != nil {
		return fmt.Errorf("failed to move file '%s' to '%s': %w", fileId, destination, err)
	}

	return nil
}

// findOrCreateSubfolder finds or creates a subfolder within the configured Google Drive folder
func (g *GoogleDrive) findOrCreateSubfolder(folderName string) (string, error) {
	// Search for existing subfolder
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and '%s' in parents and trashed=false", folderName, g.FolderID)
	result, err := g.Client.Files.List().Q(query).Do()
	if err != nil {
		return "", fmt.Errorf("failed to search for subfolder '%s': %w", folderName, err)
	}

	if len(result.Files) > 0 {
		// Subfolder exists
		return result.Files[0].Id, nil
	}

	// Create new subfolder
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{g.FolderID},
	}

	createdFolder, err := g.Client.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create subfolder '%s': %w", folderName, err)
	}

	return createdFolder.Id, nil
}

// FindOrCreateFolder finds an existing Google Drive folder or creates a new one
func (g *GoogleDrive) FindOrCreateFolder(folderName string) (string, error) {
	// Search for existing folder
	query := fmt.Sprintf("name='%s' and mimeType='application/vnd.google-apps.folder' and trashed=false", folderName)
	r, err := g.Client.Files.List().Q(query).Do()
	if err != nil {
		return "", fmt.Errorf("unable to search for folder: %w", err)
	}

	if len(r.Files) > 0 {
		// Folder exists
		return r.Files[0].Id, nil
	}

	// Create new folder
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	file, err := g.Client.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder: %w", err)
	}

	return file.Id, nil
}
