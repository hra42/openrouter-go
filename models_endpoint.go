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

	// SupportedParameters filters models by supported parameter (e.g. "tools", "temperature").
	SupportedParameters string

	// UseRSS returns an RSS feed instead of JSON when set to "true".
	UseRSS string

	// UseRSSChatLinks includes chat links in the RSS feed when set to "true".
	UseRSSChatLinks string
}

// ListModels retrieves a list of models available through the OpenRouter API.
// Note: supported_parameters is a union of all parameters supported by all providers for each model.
// There may not be a single provider which offers all of the listed parameters for a model.
func (c *Client) ListModels(ctx context.Context, opts *ListModelsOptions) (*ModelsResponse, error) {
	endpoint := "/models"

	// Add query parameters if options are provided
	if opts != nil {
		params := url.Values{}
		if opts.Category != "" {
			params.Add("category", opts.Category)
		}
		if opts.SupportedParameters != "" {
			params.Add("supported_parameters", opts.SupportedParameters)
		}
		if opts.UseRSS != "" {
			params.Add("use_rss", opts.UseRSS)
		}
		if opts.UseRSSChatLinks != "" {
			params.Add("use_rss_chat_links", opts.UseRSSChatLinks)
		}
		if encoded := params.Encode(); encoded != "" {
			endpoint = fmt.Sprintf("%s?%s", endpoint, encoded)
		}
	}

	var response ModelsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
