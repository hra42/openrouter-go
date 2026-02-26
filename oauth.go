package openrouter

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
)

// AuthURLParams contains the parameters for building an OpenRouter authorization URL.
type AuthURLParams struct {
	// CallbackURL is the URL that OpenRouter will redirect to after authorization (required).
	CallbackURL string
	// CodeChallenge is the PKCE code challenge (optional).
	CodeChallenge string
	// CodeChallengeMethod is the method used to generate the code challenge (optional).
	CodeChallengeMethod CodeChallengeMethod
}

// GenerateCodeVerifier generates a cryptographically random PKCE code verifier.
// The verifier is a 32-byte random value encoded as base64url without padding,
// resulting in a 43-character string per RFC 7636.
func GenerateCodeVerifier() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// CreateS256CodeChallenge creates a PKCE code challenge from a code verifier
// using the S256 method: BASE64URL(SHA256(verifier)) per RFC 7636.
func CreateS256CodeChallenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}

// BuildAuthURL constructs an OpenRouter authorization URL with the given parameters.
// The base URL is typically "https://openrouter.ai/auth".
//
// Example:
//
//	verifier, _ := openrouter.GenerateCodeVerifier()
//	challenge := openrouter.CreateS256CodeChallenge(verifier)
//	authURL, _ := openrouter.BuildAuthURL("https://openrouter.ai/auth", openrouter.AuthURLParams{
//	    CallbackURL:         "https://myapp.com/callback",
//	    CodeChallenge:       challenge,
//	    CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
//	})
//	// Redirect user to authURL
func BuildAuthURL(baseURL string, params AuthURLParams) (string, error) {
	if params.CallbackURL == "" {
		return "", fmt.Errorf("callback_url is required")
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	q := u.Query()
	q.Set("callback_url", params.CallbackURL)

	if params.CodeChallenge != "" {
		q.Set("code_challenge", params.CodeChallenge)
	}

	if params.CodeChallengeMethod != "" {
		q.Set("code_challenge_method", string(params.CodeChallengeMethod))
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
