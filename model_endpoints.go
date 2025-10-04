package openrouter

import (
	"context"
	"fmt"
)

// ListModelEndpoints retrieves the list of endpoints for a specific model.
// The model is specified by its author and slug (e.g., author="openai", slug="gpt-4").
// This endpoint provides detailed information about each provider offering the model,
// including pricing, context length, supported parameters, and uptime statistics.
func (c *Client) ListModelEndpoints(ctx context.Context, author, slug string) (*ModelEndpointsResponse, error) {
	if author == "" {
		return nil, fmt.Errorf("author cannot be empty")
	}
	if slug == "" {
		return nil, fmt.Errorf("slug cannot be empty")
	}

	endpoint := fmt.Sprintf("/models/%s/%s/endpoints", author, slug)

	var response ModelEndpointsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
