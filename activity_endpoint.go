package openrouter

import (
	"context"
	"fmt"
	"net/url"
)

// ActivityOptions contains optional parameters for retrieving activity data.
type ActivityOptions struct {
	// Date filters by a single UTC date in the last 30 days (YYYY-MM-DD format).
	// Example: "2024-01-15"
	// Note: The API returns dates with timestamps (e.g., "2024-01-15 00:00:00") but expects
	// the filter parameter in YYYY-MM-DD format.
	Date string
}

// GetActivity retrieves daily user activity data grouped by model endpoint for the last 30 (completed) UTC days.
//
// If ingesting on a schedule, it is recommended to wait for ~30 minutes after the UTC boundary to
// request the previous day, because events are aggregated by request start time, and some reasoning
// models may take a few minutes to complete.
//
// Note: A provisioning key is required to access this endpoint, to ensure that your historic usage
// is not accessible to just anyone in your org with an inference API key.
//
// Example:
//
//	ctx := context.Background()
//	// Get all activity
//	activity, err := client.GetActivity(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Filter by date
//	activity, err := client.GetActivity(ctx, &openrouter.ActivityOptions{
//	    Date: "2024-01-15",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, data := range activity.Data {
//	    fmt.Printf("%s - %s: %.2f requests, $%.4f usage\n",
//	        data.Date, data.Model, data.Requests, data.Usage)
//	}
func (c *Client) GetActivity(ctx context.Context, opts *ActivityOptions) (*ActivityResponse, error) {
	endpoint := "/activity"

	// Add query parameters if options are provided
	if opts != nil && opts.Date != "" {
		params := url.Values{}
		params.Add("date", opts.Date)
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	var response ActivityResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
