package openrouter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListOrganizationMembers returns a paginated list of members in the organization
// associated with the authenticated management key.
// Requires a Provisioning API key.
//
// Example:
//
//	ctx := context.Background()
//	resp, err := client.ListOrganizationMembers(ctx, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, m := range resp.Data {
//	    fmt.Printf("%s (%s)\n", m.Email, m.Role)
//	}
func (c *Client) ListOrganizationMembers(ctx context.Context, options *ListOrganizationMembersOptions) (*ListOrganizationMembersResponse, error) {
	endpoint := "/organization/members"
	if options != nil {
		params := url.Values{}
		if options.Offset != nil {
			params.Add("offset", strconv.Itoa(*options.Offset))
		}
		if options.Limit != nil {
			params.Add("limit", strconv.Itoa(*options.Limit))
		}
		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	var response ListOrganizationMembersResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
