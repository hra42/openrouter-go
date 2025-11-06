package openrouter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncodePDFToBase64(t *testing.T) {
	// Create a temporary test PDF file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.pdf")

	// Write some PDF-like content (minimal valid PDF structure)
	pdfContent := []byte("%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj 2 0 obj<</Type/Pages/Count 1/Kids[3 0 R]>>endobj 3 0 obj<</Type/Page/MediaBox[0 0 612 792]/Parent 2 0 R>>endobj\nxref\n0 4\n0000000000 65535 f\n0000000009 00000 n\n0000000056 00000 n\n0000000115 00000 n\ntrailer<</Size 4/Root 1 0 R>>\nstartxref\n190\n%%EOF")
	err := os.WriteFile(tmpFile, pdfContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test PDF: %v", err)
	}

	// Test encoding
	dataURL, err := EncodePDFToBase64(tmpFile)
	if err != nil {
		t.Fatalf("EncodePDFToBase64 failed: %v", err)
	}

	// Verify result
	if !strings.HasPrefix(dataURL, "data:application/pdf;base64,") {
		t.Errorf("Expected data URL to start with 'data:application/pdf;base64,', got: %s", dataURL[:50])
	}

	// Verify it's not empty
	if len(dataURL) <= len("data:application/pdf;base64,") {
		t.Error("Encoded data URL appears to be empty")
	}
}

func TestEncodePDFToBase64_NonExistentFile(t *testing.T) {
	_, err := EncodePDFToBase64("/nonexistent/file.pdf")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestEncodePDFBytesToBase64(t *testing.T) {
	testData := []byte("test pdf data")
	dataURL := EncodePDFBytesToBase64(testData)

	if !strings.HasPrefix(dataURL, "data:application/pdf;base64,") {
		t.Errorf("Expected data URL to start with 'data:application/pdf;base64,', got: %s", dataURL)
	}

	// Verify the base64 part contains the data
	if !strings.Contains(dataURL, "dGVzdCBwZGYgZGF0YQ==") {
		t.Error("Base64 encoding does not match expected value")
	}
}

func TestCreateUserMessageWithPDF(t *testing.T) {
	text := "Analyze this PDF"
	pdfURL := "https://example.com/document.pdf"
	filename := "document.pdf"

	msg := CreateUserMessageWithPDF(text, pdfURL, filename)

	// Verify message structure
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	// Verify content is an array
	content, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatalf("Expected content to be []ContentPart, got %T", msg.Content)
	}

	// Verify we have text and file parts
	if len(content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(content))
	}

	// Check text part
	if content[0].Type != "text" || content[0].Text != text {
		t.Errorf("Expected text part with '%s', got type='%s' text='%s'", text, content[0].Type, content[0].Text)
	}

	// Check file part
	if content[1].Type != "file" {
		t.Errorf("Expected file part, got type='%s'", content[1].Type)
	}
	if content[1].File == nil {
		t.Fatal("File is nil")
	}
	if content[1].File.Filename != filename {
		t.Errorf("Expected filename '%s', got '%s'", filename, content[1].File.Filename)
	}
	if content[1].File.FileData != pdfURL {
		t.Errorf("Expected file data '%s', got '%s'", pdfURL, content[1].File.FileData)
	}
}

func TestCreateUserMessageWithBase64PDF(t *testing.T) {
	// Create a temporary test PDF file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.pdf")

	pdfContent := []byte("%PDF-1.4\ntest content")
	err := os.WriteFile(tmpFile, pdfContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test PDF: %v", err)
	}

	// Test encoding and message creation
	text := "Analyze this PDF"
	filename := "test.pdf"

	msg, err := CreateUserMessageWithBase64PDF(text, tmpFile, filename)
	if err != nil {
		t.Fatalf("CreateUserMessageWithBase64PDF failed: %v", err)
	}

	// Verify message structure
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	content, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatalf("Expected content to be []ContentPart, got %T", msg.Content)
	}

	if len(content) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(content))
	}

	// Verify file part has base64 data
	if content[1].File == nil {
		t.Fatal("File is nil")
	}
	if !strings.HasPrefix(content[1].File.FileData, "data:application/pdf;base64,") {
		t.Error("File data is not a base64 data URL")
	}
}

func TestCreateUserMessageWithFiles(t *testing.T) {
	text := "Analyze these files"
	files := []File{
		{
			Filename: "doc1.pdf",
			FileData: "https://example.com/doc1.pdf",
		},
		{
			Filename: "doc2.pdf",
			FileData: "https://example.com/doc2.pdf",
		},
	}

	msg := CreateUserMessageWithFiles(text, files)

	// Verify message structure
	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	content, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatalf("Expected content to be []ContentPart, got %T", msg.Content)
	}

	// Should have 1 text part + 2 file parts
	if len(content) != 3 {
		t.Fatalf("Expected 3 content parts, got %d", len(content))
	}

	// Verify text part
	if content[0].Type != "text" {
		t.Errorf("Expected first part to be text, got '%s'", content[0].Type)
	}

	// Verify file parts
	for i, file := range files {
		partIdx := i + 1
		if content[partIdx].Type != "file" {
			t.Errorf("Expected part %d to be file, got '%s'", partIdx, content[partIdx].Type)
		}
		if content[partIdx].File == nil {
			t.Fatalf("File %d is nil", partIdx)
		}
		if content[partIdx].File.Filename != file.Filename {
			t.Errorf("Expected filename '%s', got '%s'", file.Filename, content[partIdx].File.Filename)
		}
	}
}

func TestCreateFileParserPlugin(t *testing.T) {
	tests := []struct {
		name   string
		engine string
		hasConfig bool
	}{
		{"Empty engine", "", false},
		{"pdf-text engine", "pdf-text", true},
		{"mistral-ocr engine", "mistral-ocr", true},
		{"native engine", "native", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plugin := CreateFileParserPlugin(tt.engine)

			if plugin.ID != "file-parser" {
				t.Errorf("Expected ID 'file-parser', got '%s'", plugin.ID)
			}

			if tt.hasConfig {
				if plugin.PDF == nil {
					t.Error("Expected PDF config to be set, got nil")
				} else if plugin.PDF.Engine != tt.engine {
					t.Errorf("Expected engine '%s', got '%s'", tt.engine, plugin.PDF.Engine)
				}
			} else {
				if plugin.PDF != nil {
					t.Error("Expected PDF config to be nil for empty engine")
				}
			}
		})
	}
}

func TestContentBuilderWithPDF(t *testing.T) {
	builder := NewContentBuilder()
	builder.AddText("First text")
	builder.AddPDF("https://example.com/doc.pdf", "doc.pdf")
	builder.AddText("Second text")

	parts := builder.Build()

	if len(parts) != 3 {
		t.Fatalf("Expected 3 parts, got %d", len(parts))
	}

	// Check first text
	if parts[0].Type != "text" || parts[0].Text != "First text" {
		t.Error("First part should be 'First text'")
	}

	// Check PDF file
	if parts[1].Type != "file" {
		t.Errorf("Second part should be file, got '%s'", parts[1].Type)
	}
	if parts[1].File == nil {
		t.Fatal("File is nil")
	}
	if parts[1].File.Filename != "doc.pdf" {
		t.Errorf("Expected filename 'doc.pdf', got '%s'", parts[1].File.Filename)
	}

	// Check second text
	if parts[2].Type != "text" || parts[2].Text != "Second text" {
		t.Error("Third part should be 'Second text'")
	}
}

func TestContentBuilderAddFile(t *testing.T) {
	builder := NewContentBuilder()
	builder.AddFile("document.pdf", "https://example.com/doc.pdf")

	parts := builder.Build()

	if len(parts) != 1 {
		t.Fatalf("Expected 1 part, got %d", len(parts))
	}

	if parts[0].Type != "file" {
		t.Errorf("Expected file type, got '%s'", parts[0].Type)
	}
	if parts[0].File == nil {
		t.Fatal("File is nil")
	}
	if parts[0].File.Filename != "document.pdf" {
		t.Errorf("Expected filename 'document.pdf', got '%s'", parts[0].File.Filename)
	}
}
