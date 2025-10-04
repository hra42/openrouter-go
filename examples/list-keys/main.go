package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	fmt.Println("=== Listing All API Keys ===")
	fmt.Println("NOTE: This endpoint requires a Provisioning API key")
	fmt.Println("Create one at: https://openrouter.ai/settings/provisioning-keys")
	fmt.Println()

	// List all API keys
	resp, err := client.ListKeys(context.Background(), nil)
	if err != nil {
		log.Fatalf("Error listing API keys: %v\n", err)
	}

	fmt.Printf("Total API keys: %d\n\n", len(resp.Data))

	// Display key information
	if len(resp.Data) == 0 {
		fmt.Println("No API keys found.")
		return
	}

	// Calculate statistics
	activeKeys := 0
	disabledKeys := 0
	totalLimit := 0.0

	for _, key := range resp.Data {
		if key.Disabled {
			disabledKeys++
		} else {
			activeKeys++
		}
		totalLimit += key.Limit
	}

	fmt.Printf("Statistics:\n")
	fmt.Printf("  Active keys: %d\n", activeKeys)
	fmt.Printf("  Disabled keys: %d\n", disabledKeys)
	fmt.Printf("  Total limit across all keys: $%.2f\n\n", totalLimit)

	// Display each key
	fmt.Println("API Keys:")
	fmt.Println(strings.Repeat("-", 80))

	for i, key := range resp.Data {
		status := "✅ Active"
		if key.Disabled {
			status = "❌ Disabled"
		}

		fmt.Printf("\n%d. %s [%s]\n", i+1, key.Label, status)
		fmt.Printf("   Name: %s\n", key.Name)
		fmt.Printf("   Limit: $%.2f\n", key.Limit)
		fmt.Printf("   Created: %s\n", key.CreatedAt)
		fmt.Printf("   Updated: %s\n", key.UpdatedAt)
		fmt.Printf("   Hash: %s\n", key.Hash)
	}

	fmt.Println()

	// Example 2: List with pagination (offset)
	if len(resp.Data) > 1 {
		fmt.Println("\n=== Example: Using Pagination (Offset) ===")
		offset := 1
		paginatedResp, err := client.ListKeys(context.Background(), &openrouter.ListKeysOptions{
			Offset: &offset,
		})
		if err != nil {
			log.Fatalf("Error listing keys with offset: %v\n", err)
		}

		fmt.Printf("Keys starting from offset %d: %d keys\n", offset, len(paginatedResp.Data))
		if len(paginatedResp.Data) > 0 {
			fmt.Printf("First key: %s\n", paginatedResp.Data[0].Label)
		}
	}

	// Example 3: Include disabled keys
	fmt.Println("\n=== Example: Including Disabled Keys ===")
	includeDisabled := true
	allKeysResp, err := client.ListKeys(context.Background(), &openrouter.ListKeysOptions{
		IncludeDisabled: &includeDisabled,
	})
	if err != nil {
		log.Fatalf("Error listing keys with include_disabled: %v\n", err)
	}

	fmt.Printf("Total keys (including disabled): %d\n", len(allKeysResp.Data))

	disabledCount := 0
	for _, key := range allKeysResp.Data {
		if key.Disabled {
			disabledCount++
		}
	}
	fmt.Printf("Disabled keys found: %d\n", disabledCount)

	// Warning if all keys are disabled
	if activeKeys == 0 && len(resp.Data) > 0 {
		fmt.Println("\n⚠️  WARNING: All API keys are disabled!")
	} else if activeKeys > 0 {
		fmt.Printf("\n✅ You have %d active API key(s)\n", activeKeys)
	}
}
