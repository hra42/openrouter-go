package openrouter

import "context"

// ListZDREndpoints retrieves the list of endpoints compatible with Zero Data Retention.
// This endpoint previews the impact of ZDR on available endpoints across all models.
func (c *Client) ListZDREndpoints(ctx context.Context) (*ZDREndpointsResponse, error) {
	var response ZDREndpointsResponse
	if err := c.doRequest(ctx, "GET", "/endpoints/zdr", nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
