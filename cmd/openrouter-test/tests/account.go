package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunCreditsTest tests the GetCredits endpoint.
func RunCreditsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Get Credits\n")

	// Test: Get credits for authenticated user
	fmt.Printf("   Testing get credits...\n")
	start := time.Now()
	resp, err := client.GetCredits(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to get credits", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved credits (%.2fs)\n", elapsed.Seconds())

	// Display credits information
	fmt.Printf("      Total Credits: $%.2f\n", resp.Data.TotalCredits)
	fmt.Printf("      Total Usage: $%.2f\n", resp.Data.TotalUsage)

	remaining := resp.Data.TotalCredits - resp.Data.TotalUsage
	fmt.Printf("      Remaining: $%.2f\n", remaining)

	if remaining < 0 {
		fmt.Printf("      ⚠️  Warning: Usage exceeds credits (negative balance)\n")
	}

	// Validate response structure
	fmt.Printf("\n   Validating response structure...\n")

	// Check that values are non-negative (usage can exceed credits, but both should be >= 0)
	if resp.Data.TotalCredits < 0 {
		fmt.Printf("   ❌ Invalid TotalCredits value: %.2f (should be >= 0)\n", resp.Data.TotalCredits)
		return false
	}
	if resp.Data.TotalUsage < 0 {
		fmt.Printf("   ❌ Invalid TotalUsage value: %.2f (should be >= 0)\n", resp.Data.TotalUsage)
		return false
	}

	printSuccess("Response structure validation passed")

	if verbose {
		fmt.Printf("\n   Credit details:\n")
		fmt.Printf("      Total Credits: $%.4f\n", resp.Data.TotalCredits)
		fmt.Printf("      Total Usage: $%.4f\n", resp.Data.TotalUsage)
		fmt.Printf("      Remaining: $%.4f\n", remaining)

		if resp.Data.TotalCredits > 0 {
			usagePercent := (resp.Data.TotalUsage / resp.Data.TotalCredits) * 100
			fmt.Printf("      Usage: %.2f%%\n", usagePercent)
		}
	}

	// Test case variations
	fmt.Printf("\n   Testing with different contexts...\n")

	// Test with custom timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.GetCredits(ctxWithTimeout)
	if err != nil {
		printError("Failed with custom timeout", err)
		return false
	}
	printSuccess("Custom timeout context works")

	printSuccess("Get credits tests completed")
	return true
}

// RunActivityTest tests the GetActivity endpoint.
func RunActivityTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Get Activity\n")

	// Test 1: Get all activity data
	fmt.Printf("   Testing get all activity...\n")
	start := time.Now()
	resp, err := client.GetActivity(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Activity endpoint requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				return true // Don't fail the test - this is expected with regular API keys
			}
		}
		printError("Failed to get activity", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved activity data (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total activity records: %d\n", len(resp.Data))

	// Display activity information
	if len(resp.Data) > 0 {
		// Calculate some statistics
		totalUsage := 0.0
		totalRequests := 0.0
		uniqueDates := make(map[string]bool)
		uniqueModels := make(map[string]bool)

		for _, data := range resp.Data {
			totalUsage += data.Usage
			totalRequests += data.Requests
			uniqueDates[data.Date] = true
			uniqueModels[data.Model] = true
		}

		fmt.Printf("      Unique dates: %d\n", len(uniqueDates))
		fmt.Printf("      Unique models: %d\n", len(uniqueModels))
		fmt.Printf("      Total usage: $%.4f\n", totalUsage)
		fmt.Printf("      Total requests: %.0f\n", totalRequests)

		if verbose {
			fmt.Printf("\n   First 5 activity records:\n")
			for i, data := range resp.Data {
				if i >= 5 {
					break
				}
				fmt.Printf("      %d. %s - %s\n", i+1, data.Date, data.Model)
				fmt.Printf("         Provider: %s\n", data.ProviderName)
				fmt.Printf("         Requests: %.0f\n", data.Requests)
				fmt.Printf("         Usage: $%.4f\n", data.Usage)
				fmt.Printf("         Tokens: %.0f prompt, %.0f completion", data.PromptTokens, data.CompletionTokens)
				if data.ReasoningTokens > 0 {
					fmt.Printf(", %.0f reasoning", data.ReasoningTokens)
				}
				fmt.Printf("\n")
				if data.BYOKUsageInference > 0 {
					fmt.Printf("         BYOK Usage: $%.4f\n", data.BYOKUsageInference)
				}
			}
		} else if len(resp.Data) > 0 {
			// Show just one example in non-verbose mode
			example := resp.Data[0]
			fmt.Printf("      Example: %s - %s (%.0f requests, $%.4f)\n",
				example.Date, example.Model, example.Requests, example.Usage)
		}
	} else {
		fmt.Printf("   ℹ️  No activity data found (this is normal for new accounts)\n")
	}

	// Test 2: Filter by specific date
	if len(resp.Data) > 0 {
		// Get the most recent date from the data
		// The API returns dates like "2025-10-03 00:00:00" but expects "2025-10-03" format
		latestDateRaw := resp.Data[0].Date

		// Parse and extract just the date part (YYYY-MM-DD)
		var latestDate string
		if len(latestDateRaw) >= 10 {
			latestDate = latestDateRaw[:10] // Extract YYYY-MM-DD
		} else {
			latestDate = latestDateRaw
		}

		fmt.Printf("\n   Testing date filter (%s)...\n", latestDate)
		start = time.Now()
		dateResp, err := client.GetActivity(ctx, &openrouter.ActivityOptions{
			Date: latestDate,
		})
		elapsed = time.Since(start)

		if err != nil {
			printError("Failed to get activity with date filter", err)
			return false
		}

		fmt.Printf("   ✅ Retrieved activity for %s (%.2fs)\n", latestDate, elapsed.Seconds())
		fmt.Printf("      Records for this date: %d\n", len(dateResp.Data))

		// Verify all records match the requested date
		// Note: API returns dates with timestamps, so we need to check just the date portion
		allMatch := true
		for _, data := range dateResp.Data {
			// Extract date portion from response (might be "2025-10-03" or "2025-10-03 00:00:00")
			responseDate := data.Date
			if len(responseDate) >= 10 {
				responseDate = responseDate[:10]
			}
			if responseDate != latestDate {
				fmt.Printf("   ❌ Found record with mismatched date: %s (expected %s)\n", data.Date, latestDate)
				allMatch = false
				break
			}
		}

		if allMatch && len(dateResp.Data) > 0 {
			printSuccess("All records match the requested date")
		}

		if verbose && len(dateResp.Data) > 0 {
			// Show activity breakdown by model for this date
			fmt.Printf("\n   Activity breakdown for %s:\n", latestDate)
			modelUsage := make(map[string]float64)
			modelRequests := make(map[string]float64)

			for _, data := range dateResp.Data {
				modelUsage[data.Model] += data.Usage
				modelRequests[data.Model] += data.Requests
			}

			fmt.Printf("   %-40s %-12s %-12s\n", "Model", "Requests", "Usage")
			fmt.Printf("   %s\n", strings.Repeat("-", 70))
			count := 0
			for model, usage := range modelUsage {
				if count >= 5 {
					fmt.Printf("   ... and %d more models\n", len(modelUsage)-5)
					break
				}
				fmt.Printf("   %-40s %-12.0f $%-11.4f\n",
					truncateString(model, 40),
					modelRequests[model],
					usage,
				)
				count++
			}
		}
	}

	// Test 3: Validate response structure
	if len(resp.Data) > 0 {
		fmt.Printf("\n   Validating response structure...\n")
		firstRecord := resp.Data[0]

		// Check required fields
		if firstRecord.Date == "" {
			printError("Activity record missing Date", nil)
			return false
		}
		if firstRecord.Model == "" {
			printError("Activity record missing Model", nil)
			return false
		}
		if firstRecord.ModelPermaslug == "" {
			printError("Activity record missing ModelPermaslug", nil)
			return false
		}
		if firstRecord.EndpointID == "" {
			printError("Activity record missing EndpointID", nil)
			return false
		}
		if firstRecord.ProviderName == "" {
			printError("Activity record missing ProviderName", nil)
			return false
		}

		// Numeric fields should be non-negative
		if firstRecord.Usage < 0 {
			fmt.Printf("   ❌ Invalid Usage value: %.4f (should be >= 0)\n", firstRecord.Usage)
			return false
		}
		if firstRecord.BYOKUsageInference < 0 {
			fmt.Printf("   ❌ Invalid BYOKUsageInference value: %.4f (should be >= 0)\n", firstRecord.BYOKUsageInference)
			return false
		}
		if firstRecord.Requests < 0 {
			fmt.Printf("   ❌ Invalid Requests value: %.0f (should be >= 0)\n", firstRecord.Requests)
			return false
		}
		if firstRecord.PromptTokens < 0 {
			fmt.Printf("   ❌ Invalid PromptTokens value: %.0f (should be >= 0)\n", firstRecord.PromptTokens)
			return false
		}
		if firstRecord.CompletionTokens < 0 {
			fmt.Printf("   ❌ Invalid CompletionTokens value: %.0f (should be >= 0)\n", firstRecord.CompletionTokens)
			return false
		}
		if firstRecord.ReasoningTokens < 0 {
			fmt.Printf("   ❌ Invalid ReasoningTokens value: %.0f (should be >= 0)\n", firstRecord.ReasoningTokens)
			return false
		}

		printSuccess("Response structure validation passed")

		if verbose {
			fmt.Printf("\n   First record details:\n")
			fmt.Printf("      Date: %s\n", firstRecord.Date)
			fmt.Printf("      Model: %s\n", firstRecord.Model)
			fmt.Printf("      Model Permaslug: %s\n", firstRecord.ModelPermaslug)
			fmt.Printf("      Endpoint ID: %s\n", firstRecord.EndpointID)
			fmt.Printf("      Provider: %s\n", firstRecord.ProviderName)
			fmt.Printf("      Usage: $%.4f\n", firstRecord.Usage)
			if firstRecord.BYOKUsageInference > 0 {
				fmt.Printf("      BYOK Usage (Inference): $%.4f\n", firstRecord.BYOKUsageInference)
			}
			fmt.Printf("      Requests: %.0f\n", firstRecord.Requests)
			fmt.Printf("      Prompt Tokens: %.0f\n", firstRecord.PromptTokens)
			fmt.Printf("      Completion Tokens: %.0f\n", firstRecord.CompletionTokens)
			if firstRecord.ReasoningTokens > 0 {
				fmt.Printf("      Reasoning Tokens: %.0f\n", firstRecord.ReasoningTokens)
			}
		}
	}

	// Test 4: Test with invalid date format
	fmt.Printf("\n   Testing error handling with invalid date...\n")
	_, err = client.GetActivity(ctx, &openrouter.ActivityOptions{
		Date: "invalid-date-format",
	})

	if err != nil {
		printSuccess("Error handling works correctly")
		printVerbose(verbose, "Error: %v", err)
	} else {
		fmt.Printf("   ⚠️  No error with invalid date (API may be lenient)\n")
	}

	// Test 5: Test with custom timeout
	fmt.Printf("\n   Testing with custom timeout...\n")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.GetActivity(ctxWithTimeout, nil)
	if err != nil {
		// Only fail if it's not a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Provisioning key required (expected)\n")
			} else {
				printError("Failed with custom timeout", err)
				return false
			}
		} else if err != context.DeadlineExceeded {
			printError("Failed with custom timeout", err)
			return false
		}
	} else {
		printSuccess("Custom timeout context works")
	}

	printSuccess("Get activity tests completed")
	return true
}
