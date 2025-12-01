package openrouter

import (
	"context"
)

// EmbeddingOption is a functional option for embedding requests.
type EmbeddingOption func(*EmbeddingRequest)

// CreateEmbedding sends an embedding request to the OpenRouter API.
// The input can be a string or []string for text embeddings.
func (c *Client) CreateEmbedding(ctx context.Context, input interface{}, model string, opts ...EmbeddingOption) (*EmbeddingResponse, error) {
	// Validate inputs
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if input == nil {
		return nil, &ValidationError{
			Field:   "input",
			Message: "input is required",
		}
	}

	// Validate empty string input
	if str, ok := input.(string); ok && str == "" {
		return nil, &ValidationError{
			Field:   "input",
			Message: "input cannot be empty",
		}
	}

	// Validate empty slice input
	if arr, ok := input.([]string); ok && len(arr) == 0 {
		return nil, &ValidationError{
			Field:   "input",
			Message: "input cannot be empty",
		}
	}

	if model == "" {
		return nil, ErrNoModel
	}

	// Build request
	req := &EmbeddingRequest{
		Input: input,
		Model: model,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Make request
	var resp EmbeddingResponse
	err := c.doRequest(ctx, "POST", "/embeddings", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateEmbeddings is an alias for CreateEmbedding that emphasizes batch embedding.
// The input should be a []string for batch text embeddings.
func (c *Client) CreateEmbeddings(ctx context.Context, inputs []string, model string, opts ...EmbeddingOption) (*EmbeddingResponse, error) {
	return c.CreateEmbedding(ctx, inputs, model, opts...)
}

// ListEmbeddingsModels retrieves a list of available embedding models.
func (c *Client) ListEmbeddingsModels(ctx context.Context) (*EmbeddingsModelsResponse, error) {
	var response EmbeddingsModelsResponse
	if err := c.doRequest(ctx, "GET", "/embeddings/models", nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// WithEmbeddingEncodingFormat sets the encoding format for the embeddings.
// Valid values are "float" (default) or "base64".
func WithEmbeddingEncodingFormat(format string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		r.EncodingFormat = format
	}
}

// WithEmbeddingDimensions sets the number of dimensions for the output embeddings.
// Only supported by some models.
func WithEmbeddingDimensions(dimensions int) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		r.Dimensions = &dimensions
	}
}

// WithEmbeddingUser sets an optional unique identifier for the end-user.
func WithEmbeddingUser(user string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		r.User = user
	}
}

// WithEmbeddingProvider sets provider-specific routing parameters.
func WithEmbeddingProvider(provider Provider) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		r.Provider = &provider
	}
}

// WithEmbeddingInputType sets the input type for the embeddings.
// Some models distinguish between document and query embeddings.
func WithEmbeddingInputType(inputType string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		r.InputType = inputType
	}
}

// WithEmbeddingOnlyProviders restricts the request to only use specified providers.
func WithEmbeddingOnlyProviders(providers ...string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Only = providers
	}
}

// WithEmbeddingIgnoreProviders specifies providers to skip for this request.
func WithEmbeddingIgnoreProviders(providers ...string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Ignore = providers
	}
}

// WithEmbeddingProviderOrder sets the order of providers to try.
func WithEmbeddingProviderOrder(providers ...string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Order = providers
	}
}

// WithEmbeddingAllowFallbacks controls whether to allow backup providers.
func WithEmbeddingAllowFallbacks(allow bool) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.AllowFallbacks = &allow
	}
}

// WithEmbeddingDataCollection controls whether to use providers that may store data.
// Use "allow" to allow data collection, "deny" to prevent it.
func WithEmbeddingDataCollection(policy string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.DataCollection = policy
	}
}

// WithEmbeddingProviderSort sorts providers by the specified attribute.
// Valid values: "price" (lowest cost), "throughput" (highest), "latency" (lowest)
func WithEmbeddingProviderSort(sort string) EmbeddingOption {
	return func(r *EmbeddingRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Sort = sort
	}
}

// GetEmbeddingVector extracts the embedding vector as []float64 from an EmbeddingData.
// Returns nil if the embedding is not in float format or conversion fails.
func (e *EmbeddingData) GetEmbeddingVector() []float64 {
	if e.Embedding == nil {
		return nil
	}

	// Handle the case where the embedding is already []float64
	if vec, ok := e.Embedding.([]float64); ok {
		return vec
	}

	// Handle the case where JSON unmarshals as []interface{}
	if arr, ok := e.Embedding.([]interface{}); ok {
		result := make([]float64, len(arr))
		for i, v := range arr {
			if f, ok := v.(float64); ok {
				result[i] = f
			} else {
				return nil
			}
		}
		return result
	}

	return nil
}

// GetEmbeddingBase64 extracts the embedding as a base64 string from an EmbeddingData.
// Returns empty string if the embedding is not in base64 format.
func (e *EmbeddingData) GetEmbeddingBase64() string {
	if str, ok := e.Embedding.(string); ok {
		return str
	}
	return ""
}
