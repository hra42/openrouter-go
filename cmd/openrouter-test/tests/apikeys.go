package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunKeyTest tests the GetKey endpoint to retrieve current API key information.
func RunKeyTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Get API Key Info\n")

	// Test: Get current API key information
	fmt.Printf("   Testing get API key info...\n")
	start := time.Now()
	resp, err := client.GetKey(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to get API key info", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved API key info (%.2fs)\n", elapsed.Seconds())

	// Display key information
	fmt.Printf("      Label: %s\n", resp.Data.Label)
	if resp.Data.Limit != nil {
		fmt.Printf("      Limit: $%.2f\n", *resp.Data.Limit)
	} else {
		fmt.Printf("      Limit: Unlimited\n")
	}
	fmt.Printf("      Usage: $%.2f\n", resp.Data.Usage)
	if resp.Data.LimitRemaining != nil {
		fmt.Printf("      Remaining: $%.2f\n", *resp.Data.LimitRemaining)
	} else {
		fmt.Printf("      Remaining: N/A\n")
	}
	fmt.Printf("      Free Tier: %v\n", resp.Data.IsFreeTier)
	fmt.Printf("      Provisioning Key: %v\n", resp.Data.IsProvisioningKey)

	if resp.Data.RateLimit != nil {
		fmt.Printf("      Rate Limit: %.0f requests per %s\n", resp.Data.RateLimit.Requests, resp.Data.RateLimit.Interval)
	}

	// Validate response structure
	fmt.Printf("\n   Validating response structure...\n")

	// Check required fields
	if resp.Data.Label == "" {
		printError("API key missing Label", nil)
		return false
	}

	// Usage should be non-negative
	if resp.Data.Usage < 0 {
		fmt.Printf("   ❌ Invalid Usage value: %.2f (should be >= 0)\n", resp.Data.Usage)
		return false
	}

	// If limit is set, it should be non-negative
	if resp.Data.Limit != nil && *resp.Data.Limit < 0 {
		fmt.Printf("   ❌ Invalid Limit value: %.2f (should be >= 0)\n", *resp.Data.Limit)
		return false
	}

	// If limit remaining is set, validate it matches calculation
	if resp.Data.Limit != nil && resp.Data.LimitRemaining != nil {
		expectedRemaining := *resp.Data.Limit - resp.Data.Usage
		if *resp.Data.LimitRemaining != expectedRemaining {
			fmt.Printf("   ⚠️  LimitRemaining (%.2f) doesn't match calculation (%.2f - %.2f = %.2f)\n",
				*resp.Data.LimitRemaining, *resp.Data.Limit, resp.Data.Usage, expectedRemaining)
		}
	}

	printSuccess("Response structure validation passed")

	if verbose {
		fmt.Printf("\n   API Key details:\n")
		fmt.Printf("      Label: %s\n", resp.Data.Label)
		if resp.Data.Limit != nil {
			fmt.Printf("      Limit: $%.4f\n", *resp.Data.Limit)
		} else {
			fmt.Printf("      Limit: nil (unlimited)\n")
		}
		fmt.Printf("      Usage: $%.4f\n", resp.Data.Usage)
		if resp.Data.LimitRemaining != nil {
			fmt.Printf("      Remaining: $%.4f\n", *resp.Data.LimitRemaining)
		} else {
			fmt.Printf("      Remaining: nil\n")
		}
		fmt.Printf("      Is Free Tier: %v\n", resp.Data.IsFreeTier)
		fmt.Printf("      Is Provisioning Key: %v\n", resp.Data.IsProvisioningKey)

		if resp.Data.RateLimit != nil {
			fmt.Printf("      Rate Limit:\n")
			fmt.Printf("         Interval: %s\n", resp.Data.RateLimit.Interval)
			fmt.Printf("         Requests: %.0f\n", resp.Data.RateLimit.Requests)
		} else {
			fmt.Printf("      Rate Limit: nil\n")
		}

		// Calculate usage percentage if limit exists
		if resp.Data.Limit != nil && *resp.Data.Limit > 0 {
			usagePercent := (resp.Data.Usage / *resp.Data.Limit) * 100
			fmt.Printf("      Usage: %.2f%%\n", usagePercent)

			if usagePercent > 80 {
				fmt.Printf("      ⚠️  Warning: Usage is above 80%%\n")
			}
		}
	}

	// Test with custom timeout
	fmt.Printf("\n   Testing with custom timeout...\n")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.GetKey(ctxWithTimeout)
	if err != nil {
		printError("Failed with custom timeout", err)
		return false
	}
	printSuccess("Custom timeout context works")

	// Informational checks
	if resp.Data.IsFreeTier {
		fmt.Printf("\n   ℹ️  This is a free tier API key\n")
	}

	if resp.Data.IsProvisioningKey {
		fmt.Printf("\n   ℹ️  This is a provisioning key (for account management)\n")
	} else {
		fmt.Printf("\n   ℹ️  This is an inference key (for API calls)\n")
	}

	if resp.Data.Limit != nil && resp.Data.LimitRemaining != nil {
		if *resp.Data.LimitRemaining <= 0 {
			fmt.Printf("\n   ⚠️  Warning: No credits remaining!\n")
		} else if *resp.Data.LimitRemaining < 1.0 {
			fmt.Printf("\n   ⚠️  Warning: Less than $1 remaining\n")
		}
	}

	printSuccess("Get API key info tests completed")
	return true
}

// RunListKeysTest tests the ListKeys endpoint to retrieve all API keys for the account.
func RunListKeysTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List API Keys\n")

	// Test: List all API keys
	fmt.Printf("   Testing list all API keys...\n")
	start := time.Now()
	resp, err := client.ListKeys(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  List keys endpoint requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true // Don't fail the test - this is expected with regular API keys
			}
		}
		printError("Failed to list API keys", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved API keys list (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total API keys: %d\n", len(resp.Data))

	// Display API keys information
	if len(resp.Data) > 0 {
		// Calculate some statistics
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

		fmt.Printf("      Active keys: %d\n", activeKeys)
		fmt.Printf("      Disabled keys: %d\n", disabledKeys)
		if totalLimit > 0 {
			fmt.Printf("      Total limit across all keys: $%.2f\n", totalLimit)
		}

		if verbose {
			fmt.Printf("\n   First 5 API keys:\n")
			for i, key := range resp.Data {
				if i >= 5 {
					break
				}
				status := "Active"
				if key.Disabled {
					status = "Disabled"
				}
				fmt.Printf("      %d. %s (%s)\n", i+1, key.Label, status)
				fmt.Printf("         Name: %s\n", key.Name)
				fmt.Printf("         Limit: $%.2f\n", key.Limit)
				fmt.Printf("         Created: %s\n", key.CreatedAt)
				fmt.Printf("         Updated: %s\n", key.UpdatedAt)
				fmt.Printf("         Hash: %s\n", key.Hash)
			}
		} else if len(resp.Data) > 0 {
			// Show just one example in non-verbose mode
			example := resp.Data[0]
			status := "Active"
			if example.Disabled {
				status = "Disabled"
			}
			fmt.Printf("      Example: %s (%s, $%.2f limit)\n", example.Label, status, example.Limit)
		}
	} else {
		fmt.Printf("   ℹ️  No API keys found (this might be unusual)\n")
	}

	// Test 2: Filter with options (if we have keys)
	if len(resp.Data) > 0 {
		fmt.Printf("\n   Testing with include_disabled option...\n")
		includeDisabled := true
		start = time.Now()
		filteredResp, err := client.ListKeys(ctx, &openrouter.ListKeysOptions{
			IncludeDisabled: &includeDisabled,
		})
		elapsed = time.Since(start)

		if err != nil {
			printError("Failed to list keys with options", err)
			return false
		}

		fmt.Printf("   ✅ Retrieved keys with include_disabled=true (%.2fs)\n", elapsed.Seconds())
		fmt.Printf("      Keys returned: %d\n", len(filteredResp.Data))

		// Test 3: Test pagination with offset
		if len(resp.Data) > 1 {
			fmt.Printf("\n   Testing pagination with offset...\n")
			offset := 1
			start = time.Now()
			paginatedResp, err := client.ListKeys(ctx, &openrouter.ListKeysOptions{
				Offset: &offset,
			})
			elapsed = time.Since(start)

			if err != nil {
				printError("Failed to list keys with offset", err)
				return false
			}

			fmt.Printf("   ✅ Retrieved keys with offset=1 (%.2fs)\n", elapsed.Seconds())
			fmt.Printf("      Keys returned: %d\n", len(paginatedResp.Data))

			if len(paginatedResp.Data) > 0 && verbose {
				fmt.Printf("      First key after offset: %s\n", paginatedResp.Data[0].Label)
			}
		}
	}

	// Test 4: Validate response structure
	if len(resp.Data) > 0 {
		fmt.Printf("\n   Validating response structure...\n")
		firstKey := resp.Data[0]

		// Check required fields
		if firstKey.Name == "" {
			printError("API key missing Name", nil)
			return false
		}
		if firstKey.Label == "" {
			printError("API key missing Label", nil)
			return false
		}
		if firstKey.Hash == "" {
			printError("API key missing Hash", nil)
			return false
		}
		if firstKey.CreatedAt == "" {
			printError("API key missing CreatedAt", nil)
			return false
		}
		// Note: UpdatedAt is optional and may be empty for some keys
		if firstKey.UpdatedAt == "" {
			fmt.Printf("   ⚠️  API key UpdatedAt is empty (may be normal for some keys)\n")
		}

		// Limit should be non-negative
		if firstKey.Limit < 0 {
			fmt.Printf("   ❌ Invalid Limit value: %.2f (should be >= 0)\n", firstKey.Limit)
			return false
		}

		printSuccess("Response structure validation passed")

		if verbose {
			fmt.Printf("\n   First key details:\n")
			fmt.Printf("      Name: %s\n", firstKey.Name)
			fmt.Printf("      Label: %s\n", firstKey.Label)
			fmt.Printf("      Limit: $%.4f\n", firstKey.Limit)
			fmt.Printf("      Disabled: %v\n", firstKey.Disabled)
			fmt.Printf("      Created At: %s\n", firstKey.CreatedAt)
			fmt.Printf("      Updated At: %s\n", firstKey.UpdatedAt)
			fmt.Printf("      Hash: %s\n", firstKey.Hash)
		}
	}

	// Test 5: Test with custom timeout
	fmt.Printf("\n   Testing with custom timeout...\n")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.ListKeys(ctxWithTimeout, nil)
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

	// Informational summary
	if len(resp.Data) > 0 {
		activeCount := 0
		for _, key := range resp.Data {
			if !key.Disabled {
				activeCount++
			}
		}

		if activeCount == 0 {
			fmt.Printf("\n   ⚠️  Warning: No active API keys found!\n")
		} else if activeCount == len(resp.Data) {
			fmt.Printf("\n   ℹ️  All %d API keys are active\n", activeCount)
		} else {
			fmt.Printf("\n   ℹ️  %d active, %d disabled API keys\n", activeCount, len(resp.Data)-activeCount)
		}

		// Test GetKeyByHash with the first key
		if len(resp.Data) > 0 {
			firstHash := resp.Data[0].Hash
			fmt.Printf("\n   Testing GetKeyByHash with hash: %s\n", firstHash)
			start = time.Now()
			keyDetails, err := client.GetKeyByHash(ctx, firstHash)
			elapsed = time.Since(start)

			if err != nil {
				printError("Failed to get key by hash", err)
				return false
			}

			fmt.Printf("   ✅ Retrieved key details by hash (%.2fs)\n", elapsed.Seconds())

			// Validate that the details match
			if keyDetails.Data.Hash != firstHash {
				fmt.Printf("   ❌ Hash mismatch: expected %s, got %s\n", firstHash, keyDetails.Data.Hash)
				return false
			}
			if keyDetails.Data.Label != resp.Data[0].Label {
				fmt.Printf("   ❌ Label mismatch: expected %s, got %s\n", resp.Data[0].Label, keyDetails.Data.Label)
				return false
			}

			printSuccess("GetKeyByHash validation passed")

			if verbose {
				fmt.Printf("\n   Key details retrieved:\n")
				fmt.Printf("      Hash: %s\n", keyDetails.Data.Hash)
				fmt.Printf("      Label: %s\n", keyDetails.Data.Label)
				fmt.Printf("      Name: %s\n", keyDetails.Data.Name)
				fmt.Printf("      Limit: $%.2f\n", keyDetails.Data.Limit)
				fmt.Printf("      Disabled: %v\n", keyDetails.Data.Disabled)
			}

			// Test with empty hash (should fail)
			fmt.Printf("\n   Testing GetKeyByHash validation...\n")
			_, err = client.GetKeyByHash(ctx, "")
			if err == nil {
				printError("Should have failed with empty hash", nil)
				return false
			}
			if _, ok := openrouter.IsValidationError(err); !ok {
				fmt.Printf("   ❌ Expected ValidationError for empty hash, got %T\n", err)
				return false
			}
			printSuccess("Empty hash validation works")
		}
	}

	printSuccess("List API keys tests completed")
	return true
}

// RunCreateKeyTest tests the CreateKey endpoint to create a new API key.
func RunCreateKeyTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Create API Key\n")

	// Test: Create API key
	fmt.Printf("\n   Testing create API key...\n")
	start := time.Now()

	// Create a key with a timestamp to make it unique and identifiable
	keyName := fmt.Sprintf("Test Key (Created by openrouter-go test suite at %s)", time.Now().Format("2006-01-02 15:04:05"))
	limit := 1.0 // $1 limit for testing

	resp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
		Name:  keyName,
		Limit: &limit,
	})
	elapsed := time.Since(start)

	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Create key endpoint requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true // Don't fail the test - this is expected with regular API keys
			}
		}
		printError("Failed to create API key", err)
		return false
	}

	fmt.Printf("   ✅ Created API key (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Name: %s\n", resp.Data.Name)
	fmt.Printf("      Label: %s\n", resp.Data.Label)
	fmt.Printf("      Limit: $%.2f\n", resp.Data.Limit)

	// Validate response structure
	fmt.Printf("\n   Validating response structure...\n")

	// Check that the key was returned
	if resp.Key == "" {
		printError("No API key value returned", nil)
		return false
	}

	fmt.Printf("   ✅ API key value returned (length: %d characters)\n", len(resp.Key))

	// Verify it starts with the expected prefix
	if !strings.HasPrefix(resp.Key, "sk-or-v1-") {
		fmt.Printf("   ⚠️  API key doesn't start with expected prefix 'sk-or-v1-'\n")
	}

	// Check required fields
	if resp.Data.Name != keyName {
		fmt.Printf("   ❌ API key name mismatch: expected %q, got %q\n", keyName, resp.Data.Name)
		return false
	}
	if resp.Data.Label == "" {
		printError("API key missing Label", nil)
		return false
	}
	if resp.Data.Hash == "" {
		printError("API key missing Hash", nil)
		return false
	}
	if resp.Data.CreatedAt == "" {
		printError("API key missing CreatedAt", nil)
		return false
	}
	// Note: UpdatedAt is optional and may be empty for newly created keys
	if resp.Data.UpdatedAt == "" {
		fmt.Printf("   ⚠️  API key UpdatedAt is empty (may be normal for new keys)\n")
	}

	// Validate limit
	if resp.Data.Limit != limit {
		fmt.Printf("   ❌ Limit mismatch: expected %.2f, got %.2f\n", limit, resp.Data.Limit)
		return false
	}

	// Should not be disabled on creation
	if resp.Data.Disabled {
		fmt.Printf("   ⚠️  Newly created key is disabled\n")
	}

	printSuccess("Response structure validation passed")

	if verbose {
		fmt.Printf("\n   Created key details:\n")
		fmt.Printf("      Name: %s\n", resp.Data.Name)
		fmt.Printf("      Label: %s\n", resp.Data.Label)
		fmt.Printf("      Limit: $%.4f\n", resp.Data.Limit)
		fmt.Printf("      Disabled: %v\n", resp.Data.Disabled)
		fmt.Printf("      Created At: %s\n", resp.Data.CreatedAt)
		fmt.Printf("      Updated At: %s\n", resp.Data.UpdatedAt)
		fmt.Printf("      Hash: %s\n", resp.Data.Hash)
		fmt.Printf("      Key (first 20 chars): %s...\n", resp.Key[:min(20, len(resp.Key))])
	}

	// Important security reminder
	fmt.Printf("\n   ⚠️  IMPORTANT SECURITY REMINDERS:\n")
	fmt.Printf("      1. The full API key value is: %s\n", resp.Key)
	fmt.Printf("      2. This is the ONLY time this value will be shown!\n")
	fmt.Printf("      3. Store it securely or delete it if you don't need it\n")
	fmt.Printf("      4. You can delete this test key at: https://openrouter.ai/settings/keys\n")

	// Test validation
	fmt.Printf("\n   Testing input validation...\n")

	// Test with empty name (should fail)
	_, err = client.CreateKey(ctx, &openrouter.CreateKeyRequest{
		Name: "",
	})
	if err == nil {
		printError("Should have failed with empty name", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty name, got %T\n", err)
		return false
	}
	printSuccess("Empty name validation works")

	// Test with nil request (should fail)
	_, err = client.CreateKey(ctx, nil)
	if err == nil {
		printError("Should have failed with nil request", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for nil request, got %T\n", err)
		return false
	}
	printSuccess("Nil request validation works")

	printSuccess("Create API key tests completed")

	// Clean up: Delete the test key
	fmt.Printf("\n   Cleaning up: Deleting test key...\n")
	deleteResp, err := client.DeleteKey(ctx, resp.Data.Hash)
	if err != nil {
		fmt.Printf("   ⚠️  Warning: Failed to delete test key: %v\n", err)
		fmt.Printf("   You may need to manually delete the test key at: https://openrouter.ai/settings/keys\n")
		fmt.Printf("   Look for: %s\n", keyName)
	} else if deleteResp.Data.Success {
		printSuccess("Test key deleted successfully")
	}

	return true
}

// RunUpdateKeyTest tests the UpdateKey endpoint to modify API key properties.
func RunUpdateKeyTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Update API Key\n")

	// Create a test key specifically for update testing
	fmt.Printf("\n   Creating a temporary key for update testing...\n")
	keyName := fmt.Sprintf("UPDATE TEST - Created at %s", time.Now().Format("2006-01-02 15:04:05"))
	initialLimit := 1.0

	createResp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
		Name:  keyName,
		Limit: &initialLimit,
	})

	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Update key endpoint requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true // Don't fail the test - this is expected with regular API keys
			}
		}
		printError("Failed to create temporary key for update testing", err)
		return false
	}

	testKeyHash := createResp.Data.Hash
	fmt.Printf("   ✅ Created temporary key: %s (hash: %s)\n", createResp.Data.Label, testKeyHash)

	// Test 1: Update just the name
	fmt.Printf("\n   Testing update key name...\n")
	newName := fmt.Sprintf("Updated at %s", time.Now().Format("15:04:05"))
	start := time.Now()
	updateResp, err := client.UpdateKey(ctx, testKeyHash, &openrouter.UpdateKeyRequest{
		Name: &newName,
	})
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to update key name", err)
		return false
	}

	fmt.Printf("   ✅ Updated key name (%.2fs)\n", elapsed.Seconds())

	if updateResp.Data.Name != newName {
		fmt.Printf("   ❌ Name not updated: expected %q, got %q\n", newName, updateResp.Data.Name)
		return false
	}
	fmt.Printf("   ✅ Name update verified: %s\n", updateResp.Data.Name)

	// Test 2: Update the limit
	fmt.Printf("\n   Testing update key limit...\n")
	newLimit := 2.0
	start = time.Now()
	updateResp, err = client.UpdateKey(ctx, testKeyHash, &openrouter.UpdateKeyRequest{
		Limit: &newLimit,
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to update key limit", err)
		return false
	}

	fmt.Printf("   ✅ Updated key limit (%.2fs)\n", elapsed.Seconds())

	if updateResp.Data.Limit != newLimit {
		fmt.Printf("   ❌ Limit not updated: expected %.2f, got %.2f\n", newLimit, updateResp.Data.Limit)
		return false
	}
	fmt.Printf("   ✅ Limit update verified: $%.2f\n", updateResp.Data.Limit)

	// Test 3: Update disabled status
	fmt.Printf("\n   Testing update key disabled status...\n")
	newDisabled := true
	start = time.Now()
	updateResp, err = client.UpdateKey(ctx, testKeyHash, &openrouter.UpdateKeyRequest{
		Disabled: &newDisabled,
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to update key disabled status", err)
		return false
	}

	fmt.Printf("   ✅ Updated key disabled status (%.2fs)\n", elapsed.Seconds())

	if updateResp.Data.Disabled != newDisabled {
		fmt.Printf("   ❌ Disabled status not updated: expected %v, got %v\n", newDisabled, updateResp.Data.Disabled)
		return false
	}
	fmt.Printf("   ✅ Disabled status update verified: %v\n", updateResp.Data.Disabled)

	// Test 4: Update multiple fields at once
	fmt.Printf("\n   Testing update multiple fields...\n")
	multiName := "Multi-field update test"
	multiLimit := 3.0
	reenableKey := false
	start = time.Now()
	updateResp, err = client.UpdateKey(ctx, testKeyHash, &openrouter.UpdateKeyRequest{
		Name:     &multiName,
		Limit:    &multiLimit,
		Disabled: &reenableKey,
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to update multiple fields", err)
		return false
	}

	fmt.Printf("   ✅ Updated multiple fields (%.2fs)\n", elapsed.Seconds())

	if updateResp.Data.Name != multiName {
		printError("Name not updated in multi-field update", nil)
		return false
	}
	if updateResp.Data.Limit != multiLimit {
		printError("Limit not updated in multi-field update", nil)
		return false
	}
	printSuccess("Multiple fields update verified")

	// Test validation
	fmt.Printf("\n   Testing input validation...\n")

	// Test with empty hash (should fail)
	_, err = client.UpdateKey(ctx, "", &openrouter.UpdateKeyRequest{
		Name: &newName,
	})
	if err == nil {
		printError("Should have failed with empty hash", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty hash, got %T\n", err)
		return false
	}
	printSuccess("Empty hash validation works")

	// Test with nil request (should fail)
	_, err = client.UpdateKey(ctx, testKeyHash, nil)
	if err == nil {
		printError("Should have failed with nil request", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for nil request, got %T\n", err)
		return false
	}
	printSuccess("Nil request validation works")

	// Test with non-existent hash (should fail)
	_, err = client.UpdateKey(ctx, "nonexistent-hash-12345", &openrouter.UpdateKeyRequest{
		Name: &newName,
	})
	if err == nil {
		printError("Should have failed with non-existent hash", nil)
		return false
	}
	if reqErr, ok := err.(*openrouter.RequestError); ok {
		if reqErr.StatusCode != 404 {
			fmt.Printf("   ⚠️  Expected 404 for non-existent hash, got %d\n", reqErr.StatusCode)
		} else {
			printSuccess("Non-existent hash validation works")
		}
	}

	printSuccess("Update API key tests completed")

	// Clean up: Delete the test key
	fmt.Printf("\n   Cleaning up: Deleting test key...\n")
	deleteResp, err := client.DeleteKey(ctx, testKeyHash)
	if err != nil {
		fmt.Printf("   ⚠️  Warning: Failed to delete test key: %v\n", err)
		fmt.Printf("   You may need to manually delete the test key at: https://openrouter.ai/settings/keys\n")
		fmt.Printf("   Hash: %s\n", testKeyHash)
	} else if deleteResp.Data.Success {
		printSuccess("Test key deleted successfully")
	}

	return true
}

// RunDeleteKeyTest tests the DeleteKey endpoint to remove API keys.
func RunDeleteKeyTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Delete API Key\n")

	// First create a key specifically for deletion testing
	fmt.Printf("\n   Creating a temporary key for deletion testing...\n")
	keyName := fmt.Sprintf("DELETE TEST - Created at %s (safe to delete)", time.Now().Format("2006-01-02 15:04:05"))
	limit := 0.01 // Minimal limit

	createResp, err := client.CreateKey(ctx, &openrouter.CreateKeyRequest{
		Name:  keyName,
		Limit: &limit,
	})

	if err != nil {
		// Check if it's a provisioning key error
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Delete key endpoint requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true // Don't fail the test - this is expected with regular API keys
			}
		}
		printError("Failed to create temporary key for deletion", err)
		return false
	}

	keyHash := createResp.Data.Hash
	fmt.Printf("   ✅ Created temporary key: %s (hash: %s)\n", createResp.Data.Label, keyHash)

	// Test: Delete the key
	fmt.Printf("\n   Testing delete API key...\n")
	start := time.Now()
	deleteResp, err := client.DeleteKey(ctx, keyHash)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to delete API key", err)
		fmt.Printf("   You may need to manually delete the test key at: https://openrouter.ai/settings/keys\n")
		fmt.Printf("   Look for: %s\n", keyName)
		return false
	}

	fmt.Printf("   ✅ Deleted API key (%.2fs)\n", elapsed.Seconds())

	// Validate response
	if verbose {
		fmt.Printf("   Debug: Delete response success = %v\n", deleteResp.Data.Success)
	}

	// Note: Some APIs return a 200 OK with no body on successful deletion
	// We'll verify deletion by checking if the key still exists rather than relying on the response
	if deleteResp.Data.Success {
		printSuccess("Delete operation confirmed successful via response")
	} else {
		fmt.Printf("   ℹ️  Delete response success field is false/empty (checking actual deletion status...)\n")
	}

	// Verify the key is actually gone
	fmt.Printf("\n   Verifying key was deleted...\n")
	_, err = client.GetKeyByHash(ctx, keyHash)
	if err == nil {
		printError("Key still exists after deletion!", nil)
		return false
	}

	if reqErr, ok := err.(*openrouter.RequestError); ok {
		if reqErr.StatusCode == 404 {
			printSuccess("Confirmed key no longer exists (404) - deletion successful")
		} else {
			fmt.Printf("   ⚠️  Unexpected status code when verifying deletion: %d\n", reqErr.StatusCode)
			// Still consider it a failure if we get a different error code
			return false
		}
	} else {
		// Got an error but not a RequestError - could be network issue
		fmt.Printf("   ⚠️  Got error verifying deletion: %v\n", err)
		// We'll be lenient here since the delete call itself succeeded
	}

	// Test validation
	fmt.Printf("\n   Testing input validation...\n")

	// Test with empty hash (should fail)
	_, err = client.DeleteKey(ctx, "")
	if err == nil {
		printError("Should have failed with empty hash", nil)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty hash, got %T\n", err)
		return false
	}
	printSuccess("Empty hash validation works")

	// Test with non-existent hash (should fail with 404)
	_, err = client.DeleteKey(ctx, "nonexistent-hash-12345")
	if err == nil {
		printError("Should have failed with non-existent hash", nil)
		return false
	}
	if reqErr, ok := err.(*openrouter.RequestError); ok {
		if reqErr.StatusCode != 404 {
			fmt.Printf("   ⚠️  Expected 404 for non-existent hash, got %d\n", reqErr.StatusCode)
		} else {
			printSuccess("Non-existent hash validation works")
		}
	}

	// Test double deletion (should fail with 404)
	fmt.Printf("\n   Testing double deletion...\n")
	_, err = client.DeleteKey(ctx, keyHash)
	if err == nil {
		printError("Should have failed deleting already-deleted key", nil)
		return false
	}
	if reqErr, ok := err.(*openrouter.RequestError); ok {
		if reqErr.StatusCode == 404 {
			printSuccess("Double deletion properly fails with 404")
		} else {
			fmt.Printf("   ⚠️  Expected 404 for double deletion, got %d\n", reqErr.StatusCode)
		}
	}

	printSuccess("Delete API key tests completed")
	return true
}
