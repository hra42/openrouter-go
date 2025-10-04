package openrouter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// GetKey retrieves information about the current API key including usage, limits,
// and rate limit information for the authenticated user.
//
// Example:
//
//	ctx := context.Background()
//	keyInfo, err := client.GetKey(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Label: %s\n", keyInfo.Data.Label)
//	if keyInfo.Data.Limit != nil {
//	    fmt.Printf("Limit: $%.2f\n", *keyInfo.Data.Limit)
//	}
//	fmt.Printf("Usage: $%.4f\n", keyInfo.Data.Usage)
//	if keyInfo.Data.LimitRemaining != nil {
//	    fmt.Printf("Remaining: $%.4f\n", *keyInfo.Data.LimitRemaining)
//	}
func (c *Client) GetKey(ctx context.Context) (*KeyResponse, error) {
	var response KeyResponse
	if err := c.doRequest(ctx, "GET", "/key", nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListKeys returns a list of all API keys associated with the account.
// Requires a Provisioning API key (not a regular inference API key).
//
// Example:
//
//	ctx := context.Background()
//	keys, err := client.ListKeys(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, key := range keys.Data {
//	    fmt.Printf("Label: %s, Disabled: %v\n", key.Label, key.Disabled)
//	}
//
// With options:
//
//	offset := 10
//	includeDisabled := true
//	keys, err := client.ListKeys(ctx, &openrouter.ListKeysOptions{
//	    Offset: &offset,
//	    IncludeDisabled: &includeDisabled,
//	})
func (c *Client) ListKeys(ctx context.Context, options *ListKeysOptions) (*ListKeysResponse, error) {
	endpoint := "/keys"

	// Add query parameters if options are provided
	if options != nil {
		params := url.Values{}

		if options.Offset != nil {
			params.Add("offset", strconv.Itoa(*options.Offset))
		}

		if options.IncludeDisabled != nil {
			params.Add("include_disabled", strconv.FormatBool(*options.IncludeDisabled))
		}

		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	var response ListKeysResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
