package openrouter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadTextFile(t *testing.T) {
	tests := []struct {
		name        string
		extension   string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Valid TXT file",
			extension: ".txt",
			content:   "Hello, World!",
			wantErr:   false,
		},
		{
			name:      "Valid Markdown file",
			extension: ".md",
			content:   "# Heading\n\nParagraph",
			wantErr:   false,
		},
		{
			name:      "Valid JSON file",
			extension: ".json",
			content:   `{"key": "value"}`,
			wantErr:   false,
		},
		{
			name:      "Valid Python file",
			extension: ".py",
			content:   "def hello():\n    print('Hello')",
			wantErr:   false,
		},
		{
			name:      "Valid Go file",
			extension: ".go",
			content:   "package main\n\nfunc main() {}",
			wantErr:   false,
		},
		{
			name:      "Valid JavaScript file",
			extension: ".js",
			content:   "console.log('hello');",
			wantErr:   false,
		},
		{
			name:      "Valid YAML file",
			extension: ".yaml",
			content:   "key: value",
			wantErr:   false,
		},
		{
			name:        "Unsupported format .exe",
			extension:   ".exe",
			content:     "binary data",
			wantErr:     true,
			errContains: "unsupported text file format",
		},
		{
			name:        "Unsupported format .bin",
			extension:   ".bin",
			content:     "binary data",
			wantErr:     true,
			errContains: "unsupported text file format",
		},
		{
			name:      "Uppercase extension .TXT",
			extension: ".TXT",
			content:   "Hello",
			wantErr:   false,
		},
		{
			name:      "Uppercase extension .PY",
			extension: ".PY",
			content:   "print('hello')",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test"+tt.extension)
			if err := os.WriteFile(tmpFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			content, err := ReadTextFile(tmpFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTextFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("ReadTextFile() error = %v, should contain %v", err, tt.errContains)
					}
				}
				return
			}

			if content != tt.content {
				t.Errorf("ReadTextFile() content = %v, want %v", content, tt.content)
			}
		})
	}
}

func TestReadTextFile_InvalidUTF8(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	// Write invalid UTF-8 bytes
	invalidUTF8 := []byte{0xff, 0xfe, 0xfd}
	if err := os.WriteFile(tmpFile, invalidUTF8, 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	_, err := ReadTextFile(tmpFile)
	if err == nil {
		t.Error("Expected error for invalid UTF-8, got nil")
	}
	if !strings.Contains(err.Error(), "invalid UTF-8") {
		t.Errorf("Error should mention invalid UTF-8, got: %v", err)
	}
}

func TestReadTextFileWithFilename(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.py")
	expectedContent := "print('hello')"

	if err := os.WriteFile(tmpFile, []byte(expectedContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	content, filename, err := ReadTextFileWithFilename(tmpFile)
	if err != nil {
		t.Fatalf("ReadTextFileWithFilename() error = %v", err)
	}

	if content != expectedContent {
		t.Errorf("content = %v, want %v", content, expectedContent)
	}

	if filename != "test.py" {
		t.Errorf("filename = %v, want test.py", filename)
	}
}

func TestCreateUserMessageWithTextFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "code.js")
	code := "console.log('hello');"

	if err := os.WriteFile(tmpFile, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	msg, err := CreateUserMessageWithTextFile("Review this code:", tmpFile)
	if err != nil {
		t.Fatalf("CreateUserMessageWithTextFile() error = %v", err)
	}

	if msg.Role != "user" {
		t.Errorf("Role = %v, want user", msg.Role)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatalf("Content is not a string, got %T", msg.Content)
	}

	if !strings.Contains(content, "Review this code:") {
		t.Error("Content should contain user text")
	}
	if !strings.Contains(content, "code.js") {
		t.Error("Content should contain filename")
	}
	if !strings.Contains(content, code) {
		t.Error("Content should contain file content")
	}
	if !strings.Contains(content, "===") {
		t.Error("Content should contain filename header markers")
	}
}

func TestCreateUserMessageWithTextFiles(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")

	os.WriteFile(file1, []byte("Content 1"), 0644)
	os.WriteFile(file2, []byte("Content 2"), 0644)

	msg, err := CreateUserMessageWithTextFiles("Review these files:", file1, file2)
	if err != nil {
		t.Fatalf("CreateUserMessageWithTextFiles() error = %v", err)
	}

	content, ok := msg.Content.(string)
	if !ok {
		t.Fatalf("Content is not a string, got %T", msg.Content)
	}

	if !strings.Contains(content, "file1.txt") {
		t.Error("Content should contain first filename")
	}
	if !strings.Contains(content, "file2.txt") {
		t.Error("Content should contain second filename")
	}
	if !strings.Contains(content, "Content 1") {
		t.Error("Content should contain first file content")
	}
	if !strings.Contains(content, "Content 2") {
		t.Error("Content should contain second file content")
	}
}

func TestCreateUserMessageWithTextContent(t *testing.T) {
	content := "SELECT * FROM users;"
	text := "Here's a query:"
	label := "query.sql"

	msg := CreateUserMessageWithTextContent(text, content, label)

	if msg.Role != "user" {
		t.Errorf("Role = %v, want user", msg.Role)
	}

	msgContent, ok := msg.Content.(string)
	if !ok {
		t.Fatalf("Content is not a string, got %T", msg.Content)
	}

	if !strings.Contains(msgContent, text) {
		t.Error("Content should contain text")
	}
	if !strings.Contains(msgContent, label) {
		t.Error("Content should contain label")
	}
	if !strings.Contains(msgContent, content) {
		t.Error("Content should contain query content")
	}
}

func TestCreateUserMessageWithTextContent_NoLabel(t *testing.T) {
	content := "print('hello')"
	text := "Here's code:"

	msg := CreateUserMessageWithTextContent(text, content, "")

	msgContent, ok := msg.Content.(string)
	if !ok {
		t.Fatalf("Content is not a string, got %T", msg.Content)
	}

	if !strings.Contains(msgContent, text) {
		t.Error("Content should contain text")
	}
	if !strings.Contains(msgContent, content) {
		t.Error("Content should contain code content")
	}
}

func TestContentBuilder_AddTextFile(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	goCode := "package main\n\nfunc main() {}"

	if err := os.WriteFile(tmpFile, []byte(goCode), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cb := NewContentBuilder()
	cb.AddText("Here's my code:")
	cb, err := cb.AddTextFile(tmpFile)
	if err != nil {
		t.Fatalf("AddTextFile() error = %v", err)
	}

	parts := cb.Build()
	if len(parts) != 2 {
		t.Fatalf("Expected 2 parts, got %d", len(parts))
	}

	if parts[0].Type != "text" || parts[0].Text != "Here's my code:" {
		t.Error("First part should be user text")
	}

	if parts[1].Type != "text" {
		t.Error("Second part should be text type")
	}
	if !strings.Contains(parts[1].Text, "test.go") {
		t.Error("Second part should contain filename")
	}
	if !strings.Contains(parts[1].Text, goCode) {
		t.Error("Second part should contain code")
	}
}

func TestContentBuilder_AddTextFileWithLabel(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	content := "# Title\n\nContent"

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cb := NewContentBuilder()
	cb, err := cb.AddTextFileWithLabel(tmpFile, "custom_label")
	if err != nil {
		t.Fatalf("AddTextFileWithLabel() error = %v", err)
	}

	parts := cb.Build()
	if len(parts) != 1 {
		t.Fatalf("Expected 1 part, got %d", len(parts))
	}

	if !strings.Contains(parts[0].Text, "custom_label") {
		t.Error("Should contain custom label")
	}
	if !strings.Contains(parts[0].Text, content) {
		t.Error("Should contain file content")
	}
}

func TestContentBuilder_AddTextContent(t *testing.T) {
	cb := NewContentBuilder()
	content := "SELECT * FROM users;"

	cb.AddText("Here's a query:")
	cb.AddTextContent(content, "query.sql")

	parts := cb.Build()
	if len(parts) != 2 {
		t.Fatalf("Expected 2 parts, got %d", len(parts))
	}

	if !strings.Contains(parts[1].Text, "query.sql") {
		t.Error("Should contain label")
	}
	if !strings.Contains(parts[1].Text, content) {
		t.Error("Should contain content")
	}
}

func TestContentBuilder_AddTextContent_NoLabel(t *testing.T) {
	cb := NewContentBuilder()
	content := "print('hello')"

	cb.AddText("Code:")
	cb.AddTextContent(content, "")

	parts := cb.Build()
	if len(parts) != 2 {
		t.Fatalf("Expected 2 parts, got %d", len(parts))
	}

	if !strings.Contains(parts[1].Text, content) {
		t.Error("Should contain content")
	}
}

func TestValidateTextContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "Valid UTF-8",
			content: "Hello, 世界! 🌍",
			wantErr: false,
		},
		{
			name:    "Simple ASCII",
			content: "Hello World",
			wantErr: false,
		},
		{
			name:    "Invalid UTF-8",
			content: string([]byte{0xff, 0xfe}),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTextContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTextContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReadTextFile_FileNotFound(t *testing.T) {
	_, err := ReadTextFile("/nonexistent/path/file.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "failed to read") {
		t.Errorf("Error should mention read failure, got: %v", err)
	}
}

func TestCreateUserMessageWithTextFiles_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	existingFile := filepath.Join(tmpDir, "exists.txt")
	nonexistentFile := filepath.Join(tmpDir, "nonexistent.txt")

	os.WriteFile(existingFile, []byte("content"), 0644)

	_, err := CreateUserMessageWithTextFiles("text", existingFile, nonexistentFile)
	if err == nil {
		t.Error("Expected error for missing file")
	}
}
