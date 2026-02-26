package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunModelsUserTest tests the ListModelsUser endpoint.
func RunModelsUserTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List User Models\n")

	// Test 1: List user-filtered models
	fmt.Printf("   Testing list user models...\n")
	start := time.Now()
	resp, err := client.ListModelsUser(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list user models", err)
		return false
	}

	if resp == nil {
		printError("Response is nil", nil)
		return false
	}

	fmt.Printf("   ✅ Listed user models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total models: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		fmt.Printf("   ⚠️  No user models returned (this might be expected depending on account settings)\n")
		printSuccess("User models test completed")
		return true
	}

	// Display first few models
	if verbose {
		fmt.Printf("\n   First 5 user models:\n")
		for i, m := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, m.ID, m.Name)
			if m.ContextLength != nil {
				fmt.Printf("         Context Length: %.0f tokens\n", *m.ContextLength)
			}
			fmt.Printf("         Pricing: $%s/M prompt, $%s/M completion\n",
				m.Pricing.Prompt, m.Pricing.Completion)
			if m.ExpirationDate != nil {
				fmt.Printf("         Expiration: %s\n", *m.ExpirationDate)
			}
			if m.Architecture.Modality != nil {
				fmt.Printf("         Modality: %s\n", *m.Architecture.Modality)
			}
		}
	} else {
		for i, m := range resp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", m.ID, m.Name)
		}
	}

	// Test 2: Validate model structure
	fmt.Printf("\n   Validating model data structure...\n")
	firstModel := resp.Data[0]

	if firstModel.ID == "" {
		printError("Model missing ID", nil)
		return false
	}
	if firstModel.Name == "" {
		printError("Model missing Name", nil)
		return false
	}
	if firstModel.Pricing.Prompt == "" {
		printError("Model missing Prompt pricing", nil)
		return false
	}
	if firstModel.Pricing.Completion == "" {
		printError("Model missing Completion pricing", nil)
		return false
	}

	printSuccess("Model structure validation passed")

	printSuccess("User models test completed")
	return true
}
