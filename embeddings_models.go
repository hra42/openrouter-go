package openrouter

// EmbeddingRequest represents a request to create embeddings.
type EmbeddingRequest struct {
	// Input is the text or texts to embed. Can be a string, array of strings,
	// array of integers, array of integer arrays, or array of multimodal content parts.
	Input interface{} `json:"input"`
	// Model specifies the embedding model to use (required).
	Model string `json:"model"`
	// EncodingFormat specifies the format for the returned embeddings ("float" or "base64").
	EncodingFormat string `json:"encoding_format,omitempty"`
	// Dimensions specifies the number of dimensions for the output embeddings.
	Dimensions *int `json:"dimensions,omitempty"`
	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`
	// Provider contains provider-specific routing parameters.
	Provider *Provider `json:"provider,omitempty"`
	// InputType specifies the type of input (e.g., for search vs. document embeddings).
	InputType string `json:"input_type,omitempty"`
}

// EmbeddingResponse represents the response from an embeddings request.
type EmbeddingResponse struct {
	// ID is the unique identifier for this embeddings request.
	ID string `json:"id,omitempty"`
	// Object is the object type, always "list".
	Object string `json:"object"`
	// Data contains the list of embedding objects.
	Data []EmbeddingData `json:"data"`
	// Model is the model used to generate the embeddings.
	Model string `json:"model"`
	// Usage contains token usage information for the request.
	Usage *EmbeddingUsage `json:"usage,omitempty"`
}

// EmbeddingData represents a single embedding in the response.
type EmbeddingData struct {
	// Object is the object type, always "embedding".
	Object string `json:"object"`
	// Embedding is the embedding vector. Can be []float64 or base64 string depending on encoding_format.
	Embedding interface{} `json:"embedding"`
	// Index is the index of this embedding in the input list.
	Index int `json:"index"`
}

// EmbeddingUsage contains token usage information for an embeddings request.
type EmbeddingUsage struct {
	// PromptTokens is the number of tokens in the input.
	PromptTokens int `json:"prompt_tokens"`
	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
	// Cost is the cost of the request (if available).
	Cost *float64 `json:"cost,omitempty"`
}

// EmbeddingContentPart represents a content part for multimodal embedding input.
type EmbeddingContentPart struct {
	// Type specifies the type of content ("text" or "image_url").
	Type string `json:"type"`
	// Text is the text content (when Type is "text").
	Text string `json:"text,omitempty"`
	// ImageURL contains the image URL (when Type is "image_url").
	ImageURL *EmbeddingImageURL `json:"image_url,omitempty"`
}

// EmbeddingImageURL represents an image URL for multimodal embedding input.
type EmbeddingImageURL struct {
	// URL is the image URL.
	URL string `json:"url"`
}

// EmbeddingMultimodalInput represents multimodal input for embeddings.
type EmbeddingMultimodalInput struct {
	// Content is the array of content parts.
	Content []EmbeddingContentPart `json:"content"`
}

// EmbeddingsModelsResponse represents the response from the list embeddings models endpoint.
type EmbeddingsModelsResponse struct {
	Data []Model `json:"data"`
}
