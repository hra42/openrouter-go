package openrouter

// RerankRequest represents a request to rerank documents against a query.
type RerankRequest struct {
	// Model specifies the rerank model to use (required).
	Model string `json:"model"`
	// Query is the search query to rerank documents against (required).
	Query string `json:"query"`
	// Documents is the list of documents to rerank (required).
	Documents []string `json:"documents"`
	// TopN is the number of most relevant documents to return.
	TopN *int `json:"top_n,omitempty"`
	// Provider contains provider-specific routing parameters.
	Provider *Provider `json:"provider,omitempty"`
}

// RerankResponse represents the response from a rerank request.
type RerankResponse struct {
	// ID is the unique identifier for the rerank response (ORID format).
	ID string `json:"id,omitempty"`
	// Model is the model used for reranking.
	Model string `json:"model"`
	// Provider is the provider that served the rerank request.
	Provider string `json:"provider,omitempty"`
	// Results is the list of rerank results sorted by relevance.
	Results []RerankResult `json:"results"`
	// Usage contains usage statistics for the request.
	Usage *RerankUsage `json:"usage,omitempty"`
}

// RerankResult represents a single reranked document result.
type RerankResult struct {
	// Index is the index of the document in the original input list.
	Index int `json:"index"`
	// RelevanceScore is the relevance score of the document to the query.
	RelevanceScore float64 `json:"relevance_score"`
	// Document contains the original document text.
	Document RerankDocument `json:"document"`
}

// RerankDocument contains the text of a reranked document.
type RerankDocument struct {
	// Text is the document text.
	Text string `json:"text"`
}

// RerankUsage contains usage statistics for a rerank request.
type RerankUsage struct {
	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens,omitempty"`
	// SearchUnits is the number of search units consumed (Cohere billing).
	SearchUnits int `json:"search_units,omitempty"`
	// Cost is the cost of the request in credits.
	Cost *float64 `json:"cost,omitempty"`
}
