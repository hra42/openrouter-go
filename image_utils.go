package openrouter

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EncodeImageToBase64 reads an image file and encodes it to a base64 data URL.
// It automatically detects the image format based on the file extension.
// Supported formats: png, jpg, jpeg, webp, gif
func EncodeImageToBase64(imagePath string) (string, error) {
	// Read the image file
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Detect content type from file extension
	ext := strings.ToLower(filepath.Ext(imagePath))
	var contentType string
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	default:
		return "", fmt.Errorf("unsupported image format: %s (supported: png, jpg, jpeg, webp, gif)", ext)
	}

	// Encode to base64
	base64Image := base64.StdEncoding.EncodeToString(imageData)

	// Return as data URL
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Image), nil
}

// EncodeImageBytesToBase64 encodes image bytes to a base64 data URL.
// The contentType should be one of: "image/png", "image/jpeg", "image/webp", "image/gif"
func EncodeImageBytesToBase64(imageData []byte, contentType string) string {
	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Image)
}

// CreateUserMessageWithBase64Image creates a user message with a base64-encoded image.
// This is a convenience function that combines EncodeImageToBase64 and CreateUserMessageWithImage.
func CreateUserMessageWithBase64Image(text string, imagePath string) (Message, error) {
	base64Image, err := EncodeImageToBase64(imagePath)
	if err != nil {
		return Message{}, err
	}
	return CreateUserMessageWithImage(text, base64Image), nil
}

// CreateUserMessageWithBase64Images creates a user message with multiple base64-encoded images.
// This is a convenience function that combines EncodeImageToBase64 and CreateUserMessageWithImages.
func CreateUserMessageWithBase64Images(text string, imagePaths ...string) (Message, error) {
	imageURLs := make([]string, 0, len(imagePaths))
	for _, path := range imagePaths {
		base64Image, err := EncodeImageToBase64(path)
		if err != nil {
			return Message{}, fmt.Errorf("failed to encode %s: %w", path, err)
		}
		imageURLs = append(imageURLs, base64Image)
	}
	return CreateUserMessageWithImages(text, imageURLs...), nil
}
