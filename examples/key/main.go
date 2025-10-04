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

	// Get API key information
	fmt.Println("=== Getting API Key Information ===")
	resp, err := client.GetKey(context.Background())
	if err != nil {
		log.Fatalf("Error getting API key info: %v", err)
	}

	// Display key information
	fmt.Printf("\nAPI Key Details:\n")
	fmt.Printf("  Label: %s\n", resp.Data.Label)

	if resp.Data.Limit != nil {
		fmt.Printf("  Credit Limit: $%.2f\n", *resp.Data.Limit)
	} else {
		fmt.Printf("  Credit Limit: Unlimited\n")
	}

	fmt.Printf("  Usage: $%.4f\n", resp.Data.Usage)

	if resp.Data.LimitRemaining != nil {
		fmt.Printf("  Remaining: $%.4f\n", *resp.Data.LimitRemaining)

		// Calculate usage percentage
		if resp.Data.Limit != nil && *resp.Data.Limit > 0 {
			usagePercent := (resp.Data.Usage / *resp.Data.Limit) * 100
			fmt.Printf("  Usage: %.2f%%\n", usagePercent)
		}
	} else {
		fmt.Printf("  Remaining: N/A\n")
	}

	fmt.Printf("  Free Tier: %v\n", resp.Data.IsFreeTier)
	fmt.Printf("  Provisioning Key: %v\n", resp.Data.IsProvisioningKey)

	// Display rate limit information if available
	if resp.Data.RateLimit != nil {
		fmt.Printf("\nRate Limit:\n")
		fmt.Printf("  Requests: %.0f per %s\n", resp.Data.RateLimit.Requests, resp.Data.RateLimit.Interval)
	}

	// Display warnings based on usage
	fmt.Println("\nStatus Checks:")

	if resp.Data.IsFreeTier {
		fmt.Println("  ℹ️  This is a free tier API key")
	}

	if resp.Data.IsProvisioningKey {
		fmt.Println("  ℹ️  This is a provisioning key (for account management)")
	} else {
		fmt.Println("  ℹ️  This is an inference key (for API calls)")
	}

	if resp.Data.Limit != nil && resp.Data.LimitRemaining != nil {
		if *resp.Data.LimitRemaining <= 0 {
			fmt.Println("  ⚠️  Warning: No credits remaining!")
		} else if *resp.Data.LimitRemaining < 1.0 {
			fmt.Println("  ⚠️  Warning: Less than $1 remaining")
		} else if resp.Data.Limit != nil && *resp.Data.Limit > 0 {
			usagePercent := (resp.Data.Usage / *resp.Data.Limit) * 100
			if usagePercent > 80 {
				fmt.Printf("  ⚠️  Warning: Usage is above 80%% (%.2f%%)\n", usagePercent)
			} else if usagePercent > 50 {
				fmt.Printf("  ℹ️  Usage is above 50%% (%.2f%%)\n", usagePercent)
			} else {
				fmt.Println("  ✅ Credit usage is healthy")
			}
		}
	}

	fmt.Println("\nAPI key information retrieved successfully!")
}
