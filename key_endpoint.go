package openrouter

import (
	"context"
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
