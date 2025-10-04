package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Note: The Activity endpoint requires a provisioning key, not a regular inference API key
	// Provisioning keys can be created at: https://openrouter.ai/settings/provisioning-keys
	fmt.Println("Activity Endpoint Example")
	fmt.Println("=========================")
	fmt.Println()
	fmt.Println("Note: This endpoint requires a provisioning key from:")
	fmt.Println("https://openrouter.ai/settings/provisioning-keys")
	fmt.Println()

	// Create client
	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
		openrouter.WithReferer("https://github.com/hra42/openrouter-go"),
		openrouter.WithAppName("Activity Example"),
	)

	ctx := context.Background()

	// Example 1: Get all activity data (last 30 completed UTC days)
	fmt.Println("Example 1: Get All Activity Data")
	fmt.Println("---------------------------------")

	activity, err := client.GetActivity(ctx, nil)
	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("Error: This endpoint requires a provisioning key.\n")
				fmt.Printf("Please create one at: https://openrouter.ai/settings/provisioning-keys\n")
				fmt.Printf("Then set it as your OPENROUTER_API_KEY environment variable.\n\n")
				os.Exit(1)
			}
		}
		log.Fatalf("Error getting activity: %v", err)
	}

	fmt.Printf("Total activity records: %d\n", len(activity.Data))

	if len(activity.Data) == 0 {
		fmt.Println("No activity data found. This is normal for new accounts.")
		return
	}

	// Calculate summary statistics
	totalUsage := 0.0
	totalRequests := 0.0
	uniqueDates := make(map[string]bool)
	uniqueModels := make(map[string]bool)
	modelUsage := make(map[string]float64)

	for _, data := range activity.Data {
		totalUsage += data.Usage
		totalRequests += data.Requests
		uniqueDates[data.Date] = true
		uniqueModels[data.Model] = true
		modelUsage[data.Model] += data.Usage
	}

	fmt.Printf("Unique dates: %d\n", len(uniqueDates))
	fmt.Printf("Unique models: %d\n", len(uniqueModels))
	fmt.Printf("Total usage: $%.4f\n", totalUsage)
	fmt.Printf("Total requests: %.0f\n\n", totalRequests)

	// Show first 5 records
	fmt.Println("First 5 Activity Records:")
	for i, data := range activity.Data {
		if i >= 5 {
			break
		}
		fmt.Printf("\n%d. Date: %s\n", i+1, data.Date)
		fmt.Printf("   Model: %s\n", data.Model)
		fmt.Printf("   Provider: %s\n", data.ProviderName)
		fmt.Printf("   Requests: %.0f\n", data.Requests)
		fmt.Printf("   Usage: $%.4f\n", data.Usage)
		fmt.Printf("   Tokens: %.0f prompt, %.0f completion", data.PromptTokens, data.CompletionTokens)
		if data.ReasoningTokens > 0 {
			fmt.Printf(", %.0f reasoning", data.ReasoningTokens)
		}
		fmt.Println()
		if data.BYOKUsageInference > 0 {
			fmt.Printf("   BYOK Usage (Inference): $%.4f\n", data.BYOKUsageInference)
		}
	}

	// Example 2: Filter by specific date
	fmt.Println("\nExample 2: Filter by Specific Date")
	fmt.Println("-----------------------------------")

	// Get yesterday's date in YYYY-MM-DD format
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	fmt.Printf("Requesting activity for: %s\n\n", yesterday)

	dateActivity, err := client.GetActivity(ctx, &openrouter.ActivityOptions{
		Date: yesterday,
	})
	if err != nil {
		log.Fatalf("Error getting activity by date: %v", err)
	}

	fmt.Printf("Records for %s: %d\n", yesterday, len(dateActivity.Data))

	if len(dateActivity.Data) > 0 {
		// Calculate totals for this date
		dateUsage := 0.0
		dateRequests := 0.0
		modelCount := make(map[string]int)

		for _, data := range dateActivity.Data {
			dateUsage += data.Usage
			dateRequests += data.Requests
			modelCount[data.Model]++
		}

		fmt.Printf("Total usage for this date: $%.4f\n", dateUsage)
		fmt.Printf("Total requests for this date: %.0f\n", dateRequests)
		fmt.Printf("Models used: %d\n\n", len(modelCount))

		// Show breakdown by model
		fmt.Println("Breakdown by Model:")
		for model, count := range modelCount {
			fmt.Printf("  - %s: %d endpoint(s)\n", model, count)
		}
	} else {
		fmt.Printf("No activity found for %s\n", yesterday)
	}

	// Example 3: Usage analysis - find most used models
	fmt.Println("\nExample 3: Most Used Models by Spend")
	fmt.Println("-------------------------------------")

	// Sort models by usage (in a simple way for the example)
	type modelStats struct {
		model    string
		usage    float64
		requests float64
	}

	var stats []modelStats
	for model, usage := range modelUsage {
		requests := 0.0
		for _, data := range activity.Data {
			if data.Model == model {
				requests += data.Requests
			}
		}
		stats = append(stats, modelStats{model, usage, requests})
	}

	// Simple bubble sort by usage (descending)
	for i := 0; i < len(stats); i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[j].usage > stats[i].usage {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}

	// Show top 5 models by spend
	fmt.Println("Top 5 Models by Spend:")
	for i, stat := range stats {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, stat.model)
		fmt.Printf("   Total spend: $%.4f\n", stat.usage)
		fmt.Printf("   Total requests: %.0f\n", stat.requests)
		fmt.Printf("   Avg cost per request: $%.6f\n\n", stat.usage/stat.requests)
	}

	// Example 4: Provider distribution
	fmt.Println("Example 4: Provider Distribution")
	fmt.Println("---------------------------------")

	providerUsage := make(map[string]float64)
	providerRequests := make(map[string]float64)

	for _, data := range activity.Data {
		providerUsage[data.ProviderName] += data.Usage
		providerRequests[data.ProviderName] += data.Requests
	}

	fmt.Printf("Total providers used: %d\n\n", len(providerUsage))

	for provider, usage := range providerUsage {
		fmt.Printf("Provider: %s\n", provider)
		fmt.Printf("  Total spend: $%.4f\n", usage)
		fmt.Printf("  Total requests: %.0f\n", providerRequests[provider])
		fmt.Printf("  %% of total spend: %.2f%%\n\n", (usage/totalUsage)*100)
	}
}
