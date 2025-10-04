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
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: Get credit balance
	fmt.Println("=== Example 1: Get Credit Balance ===")
	getCredits(client)

	// Example 2: Display remaining credits
	fmt.Println("\n=== Example 2: Display Remaining Credits ===")
	displayRemainingCredits(client)

	// Example 3: Check if credits are running low
	fmt.Println("\n=== Example 3: Check Credit Status ===")
	checkCreditStatus(client, 5.0) // Alert if less than $5 remaining
}

func getCredits(client *openrouter.Client) {
	resp, err := client.GetCredits(context.Background())
	if err != nil {
		log.Printf("Error getting credits: %v", err)
		return
	}

	fmt.Printf("Total Credits: $%.2f\n", resp.Data.TotalCredits)
	fmt.Printf("Total Usage: $%.2f\n", resp.Data.TotalUsage)
}

func displayRemainingCredits(client *openrouter.Client) {
	resp, err := client.GetCredits(context.Background())
	if err != nil {
		log.Printf("Error getting credits: %v", err)
		return
	}

	remaining := resp.Data.TotalCredits - resp.Data.TotalUsage
	fmt.Printf("Remaining Credits: $%.2f\n", remaining)

	if resp.Data.TotalCredits > 0 {
		usagePercent := (resp.Data.TotalUsage / resp.Data.TotalCredits) * 100
		fmt.Printf("Usage: %.2f%%\n", usagePercent)
	}

	if remaining < 0 {
		fmt.Printf("⚠️  Warning: Usage exceeds purchased credits\n")
	}
}

func checkCreditStatus(client *openrouter.Client, threshold float64) {
	resp, err := client.GetCredits(context.Background())
	if err != nil {
		log.Printf("Error getting credits: %v", err)
		return
	}

	remaining := resp.Data.TotalCredits - resp.Data.TotalUsage

	fmt.Printf("Current Balance: $%.2f\n", remaining)
	fmt.Printf("Alert Threshold: $%.2f\n", threshold)

	if remaining < 0 {
		fmt.Printf("🔴 Status: Overdrawn\n")
		fmt.Printf("   You have exceeded your purchased credits by $%.2f\n", -remaining)
	} else if remaining < threshold {
		fmt.Printf("🟡 Status: Low Balance\n")
		fmt.Printf("   Your balance is below $%.2f. Consider purchasing more credits.\n", threshold)
	} else {
		fmt.Printf("🟢 Status: Sufficient Balance\n")
		fmt.Printf("   You have $%.2f remaining.\n", remaining)
	}
}
