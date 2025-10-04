package openrouter

import (
	"context"
	"fmt"
	"net/url"
)

// ListModelsOptions contains optional parameters for listing models.
type ListModelsOptions struct {
	// Category filters models by category (e.g. "programming"). Sorted from most to least used.
	Category string
}

// ListModels retrieves a list of models available through the OpenRouter API.
// Note: supported_parameters is a union of all parameters supported by all providers for each model.
// There may not be a single provider which offers all of the listed parameters for a model.
func (c *Client) ListModels(ctx context.Context, opts *ListModelsOptions) (*ModelsResponse, error) {
	endpoint := "/models"

	// Add query parameters if options are provided
	if opts != nil && opts.Category != "" {
		params := url.Values{}
		params.Add("category", opts.Category)
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	var response ModelsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
