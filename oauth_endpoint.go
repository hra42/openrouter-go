package openrouter

import (
	"context"
)

// ExchangeAuthCode exchanges an authorization code for an API key.
// This is the second step of the OAuth PKCE flow, called after the user
// has authorized the application at OpenRouter and been redirected back
// with an authorization code.
//
// If PKCE was used when creating the auth code, the code_verifier must be
// provided to complete the exchange.
//
// Example:
//
//	ctx := context.Background()
//	verifier := "your-code-verifier"
//	resp, err := client.ExchangeAuthCode(ctx, &openrouter.ExchangeAuthCodeRequest{
//	    Code:                "auth-code-from-callback",
//	    CodeVerifier:        &verifier,
//	    CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("API Key: %s\n", resp.Key)
func (c *Client) ExchangeAuthCode(ctx context.Context, request *ExchangeAuthCodeRequest) (*ExchangeAuthCodeResponse, error) {
	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	if request.Code == "" {
		return nil, &ValidationError{Message: "code is required"}
	}

	var response ExchangeAuthCodeResponse
	if err := c.doRequest(ctx, "POST", "/auth/keys", request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

