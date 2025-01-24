// pkg/fileutils/fileutils.go

package fileutils

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"
)

const (
    UploadDir     = "uploads"
    ProfileImages = "profile_images"
)

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