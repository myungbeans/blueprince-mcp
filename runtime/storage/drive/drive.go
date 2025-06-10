package drive

import (
	"fmt"

	"google.golang.org/api/drive/v3"
)

// TODO: What's the right format to represent the img file?
func (g *GoogleDrive) GetFile(filename, dirname string) (any, error) {
	return nil, nil
}

func (g *GoogleDrive) ListFiles(dirname string) ([]string, error) {
	return []string{}, nil
}

func (g *GoogleDrive) MoveFile(filename, destination string) error {
	return nil
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

// ListFolderContents lists contents of a folder by folder ID
func (g *GoogleDrive) ListFolderContents(folderID string, pageSize int64) ([]*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents", folderID)
	result, err := g.Client.Files.List().Q(query).PageSize(pageSize).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to list folder contents: %w", err)
	}
	return result.Files, nil
}
