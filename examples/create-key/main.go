package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	fmt.Println("=== Creating a New API Key ===")
	fmt.Println("NOTE: This endpoint requires a Provisioning API key")
	fmt.Println("Create one at: https://openrouter.ai/settings/provisioning-keys")
	fmt.Println()

	// Example 1: Create a key with a limit
	fmt.Println("Example 1: Creating an API key with a credit limit")
	limit := 10.0 // $10 limit
	resp, err := client.CreateKey(context.Background(), &openrouter.CreateKeyRequest{
		Name:  "My Production API Key",
		Limit: &limit,
	})
	if err != nil {
		log.Fatalf("Error creating API key: %v\n", err)
	}

	fmt.Println("\n✅ API Key Created Successfully!")
	fmt.Println("=====================================")
	fmt.Printf("Label: %s\n", resp.Data.Label)
	fmt.Printf("Limit: $%.2f\n", resp.Data.Limit)
	fmt.Printf("Created: %s\n", resp.Data.CreatedAt)
	fmt.Printf("Hash: %s\n", resp.Data.Hash)
	fmt.Println()

	// ⚠️ IMPORTANT: The key value is ONLY returned once!
	fmt.Println("⚠️  IMPORTANT - SAVE THIS API KEY NOW!")
	fmt.Println("=====================================")
	fmt.Printf("API Key: %s\n", resp.Key)
	fmt.Println()
	fmt.Println("This is the ONLY time the full key value will be displayed.")
	fmt.Println("Store it securely in your password manager or environment variables.")
	fmt.Println()

	// Example 2: Create a key with BYOK limit inclusion
	fmt.Println("\n=== Example 2: Creating a key with BYOK limit ===")
	limit2 := 50.0
	includeBYOK := true
	resp2, err := client.CreateKey(context.Background(), &openrouter.CreateKeyRequest{
		Name:               "Development Key with BYOK",
		Limit:              &limit2,
		IncludeBYOKInLimit: &includeBYOK,
	})
	if err != nil {
		log.Fatalf("Error creating API key with BYOK: %v\n", err)
	}

	fmt.Printf("Created: %s\n", resp2.Data.Label)
	fmt.Printf("Limit: $%.2f (BYOK usage counts toward this limit)\n", resp2.Data.Limit)
	fmt.Printf("API Key: %s\n", resp2.Key)
	fmt.Println()

	// Example 3: Create a key without a limit
	fmt.Println("\n=== Example 3: Creating a key without a credit limit ===")
	resp3, err := client.CreateKey(context.Background(), &openrouter.CreateKeyRequest{
		Name: "Unlimited Key",
	})
	if err != nil {
		log.Fatalf("Error creating unlimited API key: %v\n", err)
	}

	fmt.Printf("Created: %s\n", resp3.Data.Label)
	fmt.Printf("Limit: $%.2f (0 = no specific limit, uses account limit)\n", resp3.Data.Limit)
	fmt.Printf("API Key: %s\n", resp3.Key)
	fmt.Println()

	// Security best practices
	fmt.Println("=== Security Best Practices ===")
	fmt.Println("1. Store API keys in environment variables or secure vaults")
	fmt.Println("2. Never commit API keys to version control")
	fmt.Println("3. Set appropriate credit limits for each key")
	fmt.Println("4. Rotate keys periodically")
	fmt.Println("5. Delete unused keys at: https://openrouter.ai/settings/keys")
	fmt.Println()

	// Display created keys summary
	fmt.Println("=== Summary ===")
	fmt.Printf("Created %d new API keys:\n", 3)
	fmt.Printf("1. %s ($%.2f limit)\n", resp.Data.Label, resp.Data.Limit)
	fmt.Printf("2. %s ($%.2f limit, BYOK included)\n", resp2.Data.Label, resp2.Data.Limit)
	fmt.Printf("3. %s (no limit)\n", resp3.Data.Label)
	fmt.Println()
	fmt.Println("⚠️  Remember to save the API key values shown above!")
	fmt.Println("⚠️  They will not be displayed again.")
}
