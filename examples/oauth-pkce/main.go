// Package main demonstrates the OAuth PKCE authentication flow with OpenRouter.
//
// This example shows how to:
// 1. Generate a PKCE code verifier and challenge
// 2. Build an authorization URL for the user
// 3. Exchange the authorization code for an API key
//
// Usage:
//
//	go run examples/oauth-pkce/main.go
//	go run examples/oauth-pkce/main.go -code <auth-code>
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	code := flag.String("code", "", "Authorization code from callback (if you have one)")
	flag.Parse()

	// Step 1: Generate PKCE code verifier and challenge
	fmt.Println("=== OAuth PKCE Flow Demo ===")
	fmt.Println()

	verifier, err := openrouter.GenerateCodeVerifier()
	if err != nil {
		log.Fatalf("Failed to generate code verifier: %v", err)
	}
	fmt.Printf("Code Verifier:  %s\n", verifier)

	challenge := openrouter.CreateS256CodeChallenge(verifier)
	fmt.Printf("Code Challenge: %s\n", challenge)
	fmt.Println()

	// Step 2: Build the authorization URL
	authURL, err := openrouter.BuildAuthURL("https://openrouter.ai/auth", openrouter.AuthURLParams{
		CallbackURL:         "https://myapp.com/callback",
		CodeChallenge:       challenge,
		CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
	})
	if err != nil {
		log.Fatalf("Failed to build auth URL: %v", err)
	}

	fmt.Println("Direct the user to this URL to authorize:")
	fmt.Printf("  %s\n", authURL)
	fmt.Println()

	// Step 3: Exchange the authorization code (if provided)
	if *code == "" {
		fmt.Println("No authorization code provided. To complete the flow, run:")
		fmt.Printf("  go run examples/oauth-pkce/main.go -code <auth-code>\n")
		fmt.Println()
		fmt.Println("After the user authorizes, they will be redirected to your callback URL")
		fmt.Println("with a 'code' parameter that you can use here.")
		return
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required for code exchange")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))
	ctx := context.Background()

	resp, err := client.ExchangeAuthCode(ctx, &openrouter.ExchangeAuthCodeRequest{
		Code:                *code,
		CodeVerifier:        &verifier,
		CodeChallengeMethod: openrouter.CodeChallengeMethodS256,
	})
	if err != nil {
		log.Fatalf("Failed to exchange auth code: %v", err)
	}

	fmt.Printf("API Key: %s\n", resp.Key)
	if resp.UserID != nil {
		fmt.Printf("User ID: %s\n", *resp.UserID)
	}
}
