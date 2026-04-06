package openrouter

import (
	"context"
)

// RerankOption is a functional option for rerank requests.
type RerankOption func(*RerankRequest)

// Rerank sends a rerank request to the OpenRouter API.
// It reranks the given documents by relevance to the query using the specified model.
func (c *Client) Rerank(ctx context.Context, query string, documents []string, model string, opts ...RerankOption) (*RerankResponse, error) {
	// Validate inputs
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if query == "" {
		return nil, &ValidationError{
			Field:   "query",
			Message: "query is required",
		}
	}

	if len(documents) == 0 {
		return nil, &ValidationError{
			Field:   "documents",
			Message: "at least one document is required",
		}
	}

	if model == "" {
		return nil, ErrNoModel
	}

	// Build request
	req := &RerankRequest{
		Model:     model,
		Query:     query,
		Documents: documents,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Make request
	var resp RerankResponse
	err := c.doRequest(ctx, "POST", "/rerank", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// WithRerankTopN sets the number of most relevant documents to return.
func WithRerankTopN(topN int) RerankOption {
	return func(r *RerankRequest) {
		r.TopN = &topN
	}
}

// WithRerankProvider sets provider-specific routing parameters.
func WithRerankProvider(provider Provider) RerankOption {
	return func(r *RerankRequest) {
		r.Provider = &provider
	}
}

// WithRerankOnlyProviders restricts the request to only use specified providers.
func WithRerankOnlyProviders(providers ...string) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Only = providers
	}
}

// WithRerankIgnoreProviders specifies providers to skip for this request.
func WithRerankIgnoreProviders(providers ...string) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Ignore = providers
	}
}

// WithRerankProviderOrder sets the order of providers to try.
func WithRerankProviderOrder(providers ...string) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Order = providers
	}
}

// WithRerankAllowFallbacks controls whether to allow backup providers.
func WithRerankAllowFallbacks(allow bool) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.AllowFallbacks = &allow
	}
}

// WithRerankDataCollection controls whether to use providers that may store data.
// Use "allow" to allow data collection, "deny" to prevent it.
func WithRerankDataCollection(policy string) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.DataCollection = policy
	}
}

// WithRerankProviderSort sorts providers by the specified attribute.
// Valid values: "price" (lowest cost), "throughput" (highest), "latency" (lowest)
func WithRerankProviderSort(sort string) RerankOption {
	return func(r *RerankRequest) {
		if r.Provider == nil {
			r.Provider = &Provider{}
		}
		r.Provider.Sort = sort
	}
}
