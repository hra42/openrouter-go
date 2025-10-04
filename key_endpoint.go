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

// CreateKey creates a new API key with the specified name and optional limit.
// Requires a Provisioning API key (not a regular inference API key).
//
// IMPORTANT: The response contains the actual API key value in the Key field.
// This is the ONLY time the key value will be returned. Store it securely!
//
// Example:
//
//	ctx := context.Background()
//	limit := 100.0
//	includeBYOK := true
//	keyResp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
//	    Name:               "Production API Key",
//	    Limit:              &limit,
//	    IncludeBYOKInLimit: &includeBYOK,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("New API Key: %s\n", keyResp.Key) // SAVE THIS!
//	fmt.Printf("Label: %s\n", keyResp.Data.Label)
//	fmt.Printf("Limit: $%.2f\n", keyResp.Data.Limit)
func (c *Client) CreateKey(ctx context.Context, request *CreateKeyRequest) (*CreateKeyResponse, error) {
	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	if request.Name == "" {
		return nil, &ValidationError{Message: "name is required"}
	}

	var response CreateKeyResponse
	if err := c.doRequest(ctx, "POST", "/keys", request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetKeyByHash retrieves details about a specific API key by its hash.
// Requires a Provisioning API key (not a regular inference API key).
//
// Example:
//
//	ctx := context.Background()
//	// Get the hash from ListKeys or from key creation
//	hash := "abc123hash"
//	keyDetails, err := client.GetKeyByHash(ctx, hash)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Label: %s\n", keyDetails.Data.Label)
//	fmt.Printf("Limit: $%.2f\n", keyDetails.Data.Limit)
//	fmt.Printf("Disabled: %v\n", keyDetails.Data.Disabled)
func (c *Client) GetKeyByHash(ctx context.Context, hash string) (*GetKeyByHashResponse, error) {
	if hash == "" {
		return nil, &ValidationError{Message: "hash is required"}
	}

	endpoint := fmt.Sprintf("/keys/%s", hash)

	var response GetKeyByHashResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteKey deletes an API key by its hash.
// Requires a Provisioning API key (not a regular inference API key).
//
// WARNING: This operation is irreversible! The API key will be permanently deleted.
//
// Example:
//
//	ctx := context.Background()
//	// Get the hash from ListKeys or from key creation
//	hash := "abc123hash"
//	result, err := client.DeleteKey(ctx, hash)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	if result.Data.Success {
//	    fmt.Println("API key successfully deleted")
//	}
func (c *Client) DeleteKey(ctx context.Context, hash string) (*DeleteKeyResponse, error) {
	if hash == "" {
		return nil, &ValidationError{Message: "hash is required"}
	}

	endpoint := fmt.Sprintf("/keys/%s", hash)

	var response DeleteKeyResponse
	if err := c.doRequest(ctx, "DELETE", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
