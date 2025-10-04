package openrouter

import (
	"context"
)

// GetCredits retrieves the total credits purchased and used for the authenticated user.
//
// Example:
//
//	ctx := context.Background()
//	credits, err := client.GetCredits(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Credits: %.2f, Usage: %.2f, Remaining: %.2f\n",
//	    credits.Data.TotalCredits,
//	    credits.Data.TotalUsage,
//	    credits.Data.TotalCredits - credits.Data.TotalUsage)
func (c *Client) GetCredits(ctx context.Context) (*CreditsResponse, error) {
	var response CreditsResponse
	if err := c.doRequest(ctx, "GET", "/credits", nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
