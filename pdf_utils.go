package openrouter

import (
	"encoding/base64"
	"fmt"
	"os"
)

// EncodePDFToBase64 reads a PDF file and encodes it to a base64 data URL.
func EncodePDFToBase64(pdfPath string) (string, error) {
	// Read the PDF file
	pdfData, err := os.ReadFile(pdfPath)
	if err != nil {
		return "", fmt.Errorf("failed to read PDF file: %w", err)
	}

	// Encode to base64
	base64PDF := base64.StdEncoding.EncodeToString(pdfData)

	// Return as data URL
	return fmt.Sprintf("data:application/pdf;base64,%s", base64PDF), nil
}

// EncodePDFBytesToBase64 encodes PDF bytes to a base64 data URL.
func EncodePDFBytesToBase64(pdfData []byte) string {
	base64PDF := base64.StdEncoding.EncodeToString(pdfData)
	return fmt.Sprintf("data:application/pdf;base64,%s", base64PDF)
}

// CreateUserMessageWithPDF creates a user message with text and a PDF file.
// The pdfURL can be either a direct URL to a publicly accessible PDF
// or a base64-encoded data URL (use EncodePDFToBase64 to encode local files).
func CreateUserMessageWithPDF(text string, pdfURL string, filename string) Message {
	return Message{
		Role: "user",
		Content: []ContentPart{
			{Type: "text", Text: text},
			{Type: "file", File: &File{
				Filename: filename,
				FileData: pdfURL,
			}},
		},
	}
}

// CreateUserMessageWithBase64PDF creates a user message with a base64-encoded PDF.
// This is a convenience function that combines EncodePDFToBase64 and CreateUserMessageWithPDF.
func CreateUserMessageWithBase64PDF(text string, pdfPath string, filename string) (Message, error) {
	base64PDF, err := EncodePDFToBase64(pdfPath)
	if err != nil {
		return Message{}, err
	}
	return CreateUserMessageWithPDF(text, base64PDF, filename), nil
}

// CreateUserMessageWithFiles creates a user message with text and multiple files.
// Files can be PDFs, images, or other supported file types.
func CreateUserMessageWithFiles(text string, files []File) Message {
	parts := []ContentPart{
		{Type: "text", Text: text},
	}
	for _, file := range files {
		fileCopy := file // Create a copy to avoid pointer issues
		parts = append(parts, ContentPart{
			Type: "file",
			File: &fileCopy,
		})
	}
	return Message{
		Role:    "user",
		Content: parts,
	}
}

// CreateFileParserPlugin creates a file parser plugin configuration.
// The engine parameter can be:
//   - "pdf-text": Best for well-structured PDFs with clear text content (Free)
//   - "mistral-ocr": Best for scanned documents or PDFs with images ($0.0004 per 1,000 pages)
//   - "native": Only for models with native file support (charged as input tokens)
//   - "" (empty): Defaults to model's native support, then "pdf-text"
func CreateFileParserPlugin(engine string) Plugin {
	plugin := Plugin{
		ID: "file-parser",
	}
	if engine != "" {
		plugin.PDF = &PDFParserConfig{
			Engine: engine,
		}
	}
	return plugin
}

// AddFile adds a file to the content builder.
func (cb *ContentBuilder) AddFile(filename string, fileData string) *ContentBuilder {
	cb.parts = append(cb.parts, ContentPart{
		Type: "file",
		File: &File{
			Filename: filename,
			FileData: fileData,
		},
	})
	return cb
}

// AddPDF adds a PDF file to the content builder.
// The pdfURL can be either a direct URL or a base64-encoded data URL.
func (cb *ContentBuilder) AddPDF(pdfURL string, filename string) *ContentBuilder {
	return cb.AddFile(filename, pdfURL)
}

// AddBase64PDF reads and encodes a PDF file, then adds it to the content builder.
func (cb *ContentBuilder) AddBase64PDF(pdfPath string, filename string) (*ContentBuilder, error) {
	base64PDF, err := EncodePDFToBase64(pdfPath)
	if err != nil {
		return cb, err
	}
	return cb.AddPDF(base64PDF, filename), nil
}
