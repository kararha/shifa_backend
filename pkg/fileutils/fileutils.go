// pkg/fileutils/fileutils.go

package fileutils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	UploadDir     = "uploads"
	ProfileImages = "profile_images" // This will create uploads/profile_images/
	MaxFileSize   = 5 << 20          // 5MB
)

var validImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// SaveFile saves an uploaded file to the specified directory
func SaveFile(file io.Reader, filename string, directory string) (string, error) {
	// Create the upload directory if it doesn't exist
	uploadPath := filepath.Join(UploadDir, directory)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %v", err)
	}

	// Generate a unique filename to prevent overwrites
	ext := filepath.Ext(filename)
	baseFilename := strings.TrimSuffix(filename, ext)
	timestamp := time.Now().Format("20060102150405")
	newFilename := fmt.Sprintf("%s_%s%s", baseFilename, timestamp, ext)

	// Create the full file path
	filePath := filepath.Join(uploadPath, newFilename)

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()

	// Copy the file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	// Return the relative path that will be stored in the database
	return filepath.Join(directory, newFilename), nil
}

// DeleteFile removes a file from the uploads directory
func DeleteFile(relativePath string) error {
	fullPath := filepath.Join(UploadDir, relativePath)
	return os.Remove(fullPath)
}

// GetFullPath returns the full system path for a stored relative path
func GetFullPath(relativePath string) string {
	return filepath.Join(UploadDir, relativePath)
}

// ValidateImage validates the image file
func ValidateImage(file io.Reader, size int64, contentType string) error {
	if size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum limit of %d bytes", MaxFileSize)
	}

	if !validImageTypes[contentType] {
		return fmt.Errorf("invalid file type: %s", contentType)
	}

	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %v", err)
	}

	detectedType := http.DetectContentType(buffer)
	if !validImageTypes[detectedType] {
		return fmt.Errorf("invalid file content type: %s", detectedType)
	}

	return nil
}
