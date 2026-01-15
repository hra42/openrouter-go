package openrouter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

// Supported text file extensions
var supportedTextExtensions = map[string]bool{
	// Common text formats
	".txt": true, ".md": true, ".json": true, ".csv": true,

	// Code files
	".js": true, ".jsx": true, ".ts": true, ".tsx": true,
	".py": true, ".go": true, ".java": true, ".rs": true,
	".c": true, ".cpp": true, ".h": true, ".hpp": true,
	".rb": true, ".php": true, ".swift": true, ".kt": true,
	".kts": true, ".scala": true, ".sh": true, ".bash": true,

	// Config files
	".yaml": true, ".yml": true, ".toml": true, ".xml": true,
	".ini": true, ".conf": true, ".config": true, ".env": true,

	// Markup and documentation
	".html": true, ".htm": true, ".css": true, ".scss": true,
	".sass": true, ".less": true, ".sql": true, ".graphql": true,
}

// ReadTextFile reads a text file and returns its content as a string.
// It validates that the file contains valid UTF-8 text.
// Supported formats: .txt, .md, .json, .csv, .js, .py, .go, .java, .rs, .ts,
// .jsx, .tsx, .cpp, .c, .h, .rb, .php, .swift, .kt, .yaml, .yml, .toml, .xml, .ini
func ReadTextFile(filePath string) (string, error) {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if !supportedTextExtensions[ext] {
		return "", fmt.Errorf("unsupported text file format: %s", ext)
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read text file: %w", err)
	}

	// Validate UTF-8 encoding
	if !utf8.Valid(content) {
		return "", fmt.Errorf("file contains invalid UTF-8 encoding: %s", filePath)
	}

	return string(content), nil
}

// ReadTextFileWithFilename reads a text file and returns both the content and filename.
// This is useful when you want to include the filename in your message.
func ReadTextFileWithFilename(filePath string) (content string, filename string, err error) {
	content, err = ReadTextFile(filePath)
	if err != nil {
		return "", "", err
	}

	filename = filepath.Base(filePath)
	return content, filename, nil
}

// ValidateTextContent validates that the given content is valid UTF-8 text.
func ValidateTextContent(content string) error {
	if !utf8.ValidString(content) {
		return fmt.Errorf("content contains invalid UTF-8 encoding")
	}
	return nil
}

// CreateUserMessageWithTextFile creates a user message with text and content from a text file.
// The file content is included inline as text, not base64-encoded.
// Example:
//
//	msg, err := CreateUserMessageWithTextFile(
//	    "Review this code for bugs:",
//	    "/path/to/code.py",
//	)
func CreateUserMessageWithTextFile(text string, filePath string) (Message, error) {
	fileContent, filename, err := ReadTextFileWithFilename(filePath)
	if err != nil {
		return Message{}, err
	}

	// Create a message with the file content
	// Format: user text, then filename header, then file content
	fullText := fmt.Sprintf("%s\n\n=== %s ===\n%s", text, filename, fileContent)

	return Message{
		Role:    "user",
		Content: fullText,
	}, nil
}

// CreateUserMessageWithTextFiles creates a user message with multiple text files.
// Each file's content is included with a header showing the filename.
func CreateUserMessageWithTextFiles(text string, filePaths ...string) (Message, error) {
	var parts []string
	parts = append(parts, text)

	for _, path := range filePaths {
		content, filename, err := ReadTextFileWithFilename(path)
		if err != nil {
			return Message{}, fmt.Errorf("failed to read %s: %w", path, err)
		}
		parts = append(parts, fmt.Sprintf("\n=== %s ===\n%s", filename, content))
	}

	return Message{
		Role:    "user",
		Content: strings.Join(parts, "\n"),
	}, nil
}

// CreateUserMessageWithTextContent creates a user message with inline text content.
// Unlike CreateUserMessageWithTextFile, this doesn't read from a file but formats
// the content with an optional filename for context.
func CreateUserMessageWithTextContent(text string, content string, filename string) Message {
	var fullText string
	if filename != "" {
		fullText = fmt.Sprintf("%s\n\n=== %s ===\n%s", text, filename, content)
	} else {
		fullText = fmt.Sprintf("%s\n\n%s", text, content)
	}

	return Message{
		Role:    "user",
		Content: fullText,
	}
}
