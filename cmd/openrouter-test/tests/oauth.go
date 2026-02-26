package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunOAuthPKCETest tests the OAuth PKCE helper utilities and auth code API endpoints.
func RunOAuthPKCETest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: OAuth PKCE\n")

	// Section 1: Test PKCE helpers (no network required)
	fmt.Printf("   Testing PKCE helper utilities...\n")

	// Test GenerateCodeVerifier
	start := time.Now()
	verifier, err := openrouter.GenerateCodeVerifier()
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to generate code verifier", err)
		return false
	}

	if verifier == "" {
		printError("Code verifier is empty", nil)
		return false
	}

	if len(verifier) != 43 {
		fmt.Printf("   ❌ Expected verifier length 43, got %d\n", len(verifier))
		return false
	}

	fmt.Printf("   ✅ Generated code verifier (%.2fs)\n", elapsed.Seconds())
	printVerbose(verbose, "Verifier: %s", verifier)

	// Test uniqueness
	verifier2, err := openrouter.GenerateCodeVerifier()
	if err != nil {
		printError("Failed to generate second code verifier", err)
		return false
	}
	if verifier == verifier2 {
		printError("Generated identical verifiers", nil)
		return false
	}
	printSuccess("Code verifiers are unique")

	// Test CreateS256CodeChallenge with RFC 7636 test vector
	rfcVerifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expectedChallenge := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	challenge := openrouter.CreateS256CodeChallenge(rfcVerifier)
	if challenge != expectedChallenge {
		fmt.Printf("   ❌ RFC 7636 test vector mismatch: expected %q, got %q\n", expectedChallenge, challenge)
		return false
	}
	printSuccess("S256 code challenge matches RFC 7636 test vector")

	// Test with our generated verifier
	challenge = openrouter.CreateS256CodeChallenge(verifier)
	if challenge == "" {
		printError("Code challenge is empty", nil)
		return false
	}
	printVerbose(verbose, "Challenge: %s", challenge)
	printSuccess("Created S256 code challenge from generated verifier")

	// Test BuildAuthURL
	authURL, err := openrouter.BuildAuthURL("https://openrouter.ai/auth", openrouter.AuthURLParams{
		CallbackURL:         "https://myapp.com/callback",
		CodeChallenge:       challenge,
		CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
	})
	if err != nil {
		printError("Failed to build auth URL", err)
		return false
	}

	if !strings.HasPrefix(authURL, "https://openrouter.ai/auth?") {
		fmt.Printf("   ❌ Auth URL has unexpected prefix: %s\n", authURL)
		return false
	}
	if !strings.Contains(authURL, "callback_url=") {
		printError("Auth URL missing callback_url", nil)
		return false
	}
	if !strings.Contains(authURL, "code_challenge=") {
		printError("Auth URL missing code_challenge", nil)
		return false
	}
	if !strings.Contains(authURL, "code_challenge_method=S256") {
		printError("Auth URL missing code_challenge_method", nil)
		return false
	}

	printSuccess("Built auth URL with all parameters")
	printVerbose(verbose, "Auth URL: %s", authURL)

	// Test BuildAuthURL validation
	_, err = openrouter.BuildAuthURL("https://openrouter.ai/auth", openrouter.AuthURLParams{})
	if err == nil {
		printError("Expected error for empty callback_url", nil)
		return false
	}
	printSuccess("BuildAuthURL validates empty callback_url")

	// Section 2: Test CreateAuthCode API call
	fmt.Printf("\n   Testing CreateAuthCode API endpoint...\n")
	start = time.Now()
	_, err = client.CreateAuthCode(ctx, &openrouter.CreateAuthCodeRequest{
		CallbackURL:         "https://myapp.com/callback",
		CodeChallenge:       &challenge,
		CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
	})
	elapsed = time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  OAuth API requires appropriate permissions: %v\n", reqErr.Message)
				fmt.Printf("   Skipping API test (this is expected with regular inference keys)\n")
				printSuccess("OAuth PKCE helper tests completed (API test skipped)")
				return true
			}
		}
		printError("Failed to create auth code", err)
		return false
	}

	fmt.Printf("   ✅ Created auth code (%.2fs)\n", elapsed.Seconds())

	// Section 3: Test input validation on client methods
	fmt.Printf("\n   Testing input validation...\n")

	_, err = client.ExchangeAuthCode(ctx, nil)
	if err == nil {
		printError("Expected error for nil ExchangeAuthCode request", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for nil request, got %T\n", err)
		return false
	}
	printSuccess("ExchangeAuthCode nil request validation works")

	_, err = client.ExchangeAuthCode(ctx, &openrouter.ExchangeAuthCodeRequest{})
	if err == nil {
		printError("Expected error for empty code", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty code, got %T\n", err)
		return false
	}
	printSuccess("ExchangeAuthCode empty code validation works")

	_, err = client.CreateAuthCode(ctx, nil)
	if err == nil {
		printError("Expected error for nil CreateAuthCode request", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for nil request, got %T\n", err)
		return false
	}
	printSuccess("CreateAuthCode nil request validation works")

	_, err = client.CreateAuthCode(ctx, &openrouter.CreateAuthCodeRequest{})
	if err == nil {
		printError("Expected error for empty callback_url", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty callback_url, got %T\n", err)
		return false
	}
	printSuccess("CreateAuthCode empty callback_url validation works")

	printSuccess("OAuth PKCE tests completed")
	return true
}
