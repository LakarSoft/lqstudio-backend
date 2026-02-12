package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SaveFile saves an uploaded file to the storage path and returns the URL path
func SaveFile(file multipart.File, filename string, storagePath string) (string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Build full file path
	fullPath := filepath.Join(storagePath, filename)

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return public URL path (e.g., /uploads/payment-screenshots/filename.jpg)
	urlPath := filepath.Join("/uploads", "payment-screenshots", filename)
	// Convert Windows backslashes to forward slashes for URL
	urlPath = filepath.ToSlash(urlPath)

	return urlPath, nil
}

// ValidateFile validates file size and MIME type
func ValidateFile(fileHeader *multipart.FileHeader, maxSize int64, allowedTypes []string) error {
	// Check file size
	if fileHeader.Size > maxSize {
		return fmt.Errorf("file size exceeds %dMB limit", maxSize/(1024*1024))
	}

	// Open file to read content for MIME type validation
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Detect content type from file content
	contentType := http.DetectContentType(buffer[:n])

	// Check if content type is in allowed types
	allowed := false
	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("invalid file type: %s. Only JPG, JPEG, and PNG files are allowed", contentType)
	}

	// Reset file pointer to beginning for subsequent reads
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	return nil
}

// GenerateUniqueFilename generates a unique filename with timestamp and random suffix
func GenerateUniqueFilename(originalFilename string) string {
	// Get file extension (sanitized)
	ext := strings.ToLower(filepath.Ext(originalFilename))

	// Sanitize extension - only allow jpg, jpeg, png
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		ext = ".jpg" // Default to .jpg if extension is weird
	}

	// Generate timestamp
	timestamp := time.Now().Unix()

	// Generate random bytes
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomHex := hex.EncodeToString(randomBytes)

	// Format: payment_{timestamp}_{random}.ext
	return fmt.Sprintf("payment_%d_%s%s", timestamp, randomHex, ext)
}

// SanitizeFilename removes potentially dangerous characters from filename
func SanitizeFilename(filename string) string {
	// Remove any path components to prevent directory traversal
	filename = filepath.Base(filename)

	// Replace potentially dangerous characters
	dangerous := []string{"..", "/", "\\", "\x00"}
	for _, char := range dangerous {
		filename = strings.ReplaceAll(filename, char, "")
	}

	return filename
}
