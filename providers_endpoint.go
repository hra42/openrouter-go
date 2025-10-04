package openrouter

import (
	"context"
)

// ListProviders retrieves a list of all providers available through the OpenRouter API.
// Returns provider information including name, slug, and policy URLs.
//
// Example:
//
//	ctx := context.Background()
//	providers, err := client.ListProviders(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, provider := range providers.Data {
//	    fmt.Printf("%s (%s)\n", provider.Name, provider.Slug)
//	}
func (c *Client) ListProviders(ctx context.Context) (*ProvidersResponse, error) {
	var response ProvidersResponse
	if err := c.doRequest(ctx, "GET", "/providers", nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
