package openrouter

import "time"

// StreamEvent represents a server-sent event for streaming responses.
type StreamEvent struct {
	ID    string
	Event string
	Data  string
	Retry *time.Duration
}

// ErrorResponse represents an error response from the OpenRouter API.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents the error details in an error response.
type APIError struct {
	Message  string         `json:"message"`
	Type     string         `json:"type"`
	Code     string         `json:"code,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Plugin represents a plugin configuration for enhancing model responses.
type Plugin struct {
	// ID is the plugin identifier (e.g., "web" for web search, "file-parser" for file parsing)
	ID string `json:"id"`
	// Engine specifies which search engine to use ("native", "exa", or undefined for auto)
	Engine string `json:"engine,omitempty"`
	// MaxResults specifies the maximum number of search results (defaults to 5)
	MaxResults int `json:"max_results,omitempty"`
	// SearchPrompt customizes the prompt used to attach search results
	SearchPrompt string `json:"search_prompt,omitempty"`
	// PDF contains configuration for PDF file parsing (when ID is "file-parser")
	PDF *PDFParserConfig `json:"pdf,omitempty"`
}

// PDFParserConfig configures PDF file parsing behavior.
type PDFParserConfig struct {
	// Engine specifies the PDF parsing engine to use.
	// Options:
	// - "pdf-text": Best for well-structured PDFs with clear text content (Free)
	// - "mistral-ocr": Best for scanned documents or PDFs with images ($0.0004 per 1,000 pages)
	// - "native": Only for models with native file support (charged as input tokens)
	// If not specified, defaults to model's native support, then "pdf-text"
	Engine string `json:"engine,omitempty"`
}

// WebSearchOptions configures native web search behavior for supported models.
type WebSearchOptions struct {
	// SearchContextSize determines the amount of search context ("low", "medium", or "high")
	SearchContextSize string `json:"search_context_size,omitempty"`
}

// Annotation represents an annotation in a message response.
type Annotation struct {
	// Type of annotation (e.g., "url_citation", "file")
	Type string `json:"type"`
	// URLCitation contains details for URL citation annotations
	URLCitation *URLCitation `json:"url_citation,omitempty"`
	// FileAnnotation contains details for file annotations (parsed file metadata)
	FileAnnotation *FileAnnotation `json:"file,omitempty"`
}

// URLCitation represents a URL citation in a message annotation.
type URLCitation struct {
	// URL of the cited source
	URL string `json:"url"`
	// Title of the web search result
	Title string `json:"title"`
	// Content of the web search result
	Content string `json:"content,omitempty"`
	// StartIndex is the index of the first character of the citation in the message
	StartIndex int `json:"start_index"`
	// EndIndex is the index of the last character of the citation in the message
	EndIndex int `json:"end_index"`
}

// FileAnnotation represents metadata about a parsed file.
// Can be sent back in subsequent requests to avoid re-parsing costs.
type FileAnnotation struct {
	// Filename is the name of the file
	Filename string `json:"filename"`
	// ContentType is the MIME type of the file (e.g., "application/pdf")
	ContentType string `json:"content_type,omitempty"`
	// ParsedContent contains the extracted text/data from the file
	ParsedContent string `json:"parsed_content,omitempty"`
	// Metadata contains additional file information
	Metadata map[string]any `json:"metadata,omitempty"`
}
