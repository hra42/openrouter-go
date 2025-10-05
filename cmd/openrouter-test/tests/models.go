package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunModelsTest tests the ListModels endpoint with filtering and validation.
func RunModelsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Models\n")

	// Test 1: List all models
	fmt.Printf("   Testing list all models...\n")
	start := time.Now()
	resp, err := client.ListModels(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list models", err)
		return false
	}

	fmt.Printf("   ✅ Listed all models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total models: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		printError("No models returned", nil)
		return false
	}

	// Display first few models
	if verbose {
		fmt.Printf("\n   First 5 models:\n")
		for i, model := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, model.Name, model.ID)
			if model.ContextLength != nil {
				fmt.Printf("         Context: %.0f tokens\n", *model.ContextLength)
			}
			fmt.Printf("         Pricing: $%s/M prompt, $%s/M completion\n",
				model.Pricing.Prompt, model.Pricing.Completion)
		}
	} else {
		// Show just a couple in non-verbose mode
		for i, model := range resp.Data {
			if i >= 2 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", model.Name, model.ID)
		}
	}

	// Test 2: Validate model structure
	fmt.Printf("\n   Validating model data structure...\n")
	firstModel := resp.Data[0]

	// Check required fields
	if firstModel.ID == "" {
		printError("Model missing ID", nil)
		return false
	}
	if firstModel.Name == "" {
		printError("Model missing Name", nil)
		return false
	}
	if firstModel.Description == "" {
		printError("Model missing Description", nil)
		return false
	}

	// Check architecture
	if len(firstModel.Architecture.InputModalities) == 0 {
		printError("Model missing InputModalities", nil)
		return false
	}
	if len(firstModel.Architecture.OutputModalities) == 0 {
		printError("Model missing OutputModalities", nil)
		return false
	}
	if firstModel.Architecture.Tokenizer == "" {
		printError("Model missing Tokenizer", nil)
		return false
	}

	// Check pricing
	if firstModel.Pricing.Prompt == "" {
		printError("Model missing Prompt pricing", nil)
		return false
	}
	if firstModel.Pricing.Completion == "" {
		printError("Model missing Completion pricing", nil)
		return false
	}

	printSuccess("Model structure validation passed")

	if verbose {
		fmt.Printf("\n   First model details:\n")
		fmt.Printf("      ID: %s\n", firstModel.ID)
		fmt.Printf("      Name: %s\n", firstModel.Name)
		fmt.Printf("      Description: %s\n", truncateString(firstModel.Description, 80))
		if firstModel.ContextLength != nil {
			fmt.Printf("      Context Length: %.0f tokens\n", *firstModel.ContextLength)
		}
		fmt.Printf("      Input Modalities: %v\n", firstModel.Architecture.InputModalities)
		fmt.Printf("      Output Modalities: %v\n", firstModel.Architecture.OutputModalities)
		fmt.Printf("      Tokenizer: %s\n", firstModel.Architecture.Tokenizer)
		if firstModel.Architecture.InstructType != nil {
			fmt.Printf("      Instruct Type: %s\n", *firstModel.Architecture.InstructType)
		}
		fmt.Printf("      Is Moderated: %v\n", firstModel.TopProvider.IsModerated)
		if len(firstModel.SupportedParameters) > 0 {
			fmt.Printf("      Supported Parameters: %v\n", firstModel.SupportedParameters)
		}
	}

	// Test 3: Filter by category
	fmt.Printf("\n   Testing category filter (programming)...\n")
	start = time.Now()
	categoryResp, err := client.ListModels(ctx, &openrouter.ListModelsOptions{
		Category: "programming",
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list models by category", err)
		return false
	}

	fmt.Printf("   ✅ Listed programming models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Programming models: %d\n", len(categoryResp.Data))

	if len(categoryResp.Data) == 0 {
		fmt.Printf("   ⚠️  No programming models found (this might be expected)\n")
	} else if verbose {
		fmt.Printf("\n   Top 3 programming models:\n")
		for i, model := range categoryResp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, model.Name, model.ID)
		}
	}

	// Test 4: Check for specific well-known models
	fmt.Printf("\n   Checking for well-known models...\n")
	wellKnownModels := []string{
		"openai/gpt-4o",
		"anthropic/claude-3.5-sonnet",
		"google/gemini-pro",
		"meta-llama/llama-3.1-8b-instruct",
	}

	foundModels := make(map[string]bool)
	for _, model := range resp.Data {
		for _, knownModel := range wellKnownModels {
			if model.ID == knownModel {
				foundModels[knownModel] = true
			}
		}
	}

	foundCount := len(foundModels)
	fmt.Printf("   Found %d/%d well-known models\n", foundCount, len(wellKnownModels))

	if verbose {
		for _, knownModel := range wellKnownModels {
			status := "❌"
			if foundModels[knownModel] {
				status = "✅"
			}
			fmt.Printf("      %s %s\n", status, knownModel)
		}
	}

	// Test 5: Verify pricing information
	fmt.Printf("\n   Validating pricing information...\n")
	hasPricingInfo := 0
	for _, model := range resp.Data {
		if model.Pricing.Prompt != "" && model.Pricing.Completion != "" {
			hasPricingInfo++
		}
	}

	pricingPercent := (float64(hasPricingInfo) / float64(len(resp.Data))) * 100
	fmt.Printf("   %.1f%% of models have pricing info (%d/%d)\n",
		pricingPercent, hasPricingInfo, len(resp.Data))

	if pricingPercent < 90 {
		fmt.Printf("   ⚠️  Warning: Less than 90%% of models have pricing info\n")
	} else {
		printSuccess("Pricing validation passed")
	}

	printSuccess("List models tests completed")
	return true
}

// RunModelEndpointsTest tests the ListModelEndpoints endpoint for specific models.
func RunModelEndpointsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Model Endpoints\n")

	// Test 1: List endpoints for GPT-4
	fmt.Printf("   Testing endpoints for GPT-4...\n")
	start := time.Now()
	resp, err := client.ListModelEndpoints(ctx, "openai", "gpt-4")
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list model endpoints", err)
		return false
	}

	fmt.Printf("   ✅ Listed GPT-4 endpoints (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s (%s)\n", resp.Data.Name, resp.Data.ID)
	fmt.Printf("      Total endpoints: %d\n", len(resp.Data.Endpoints))

	if len(resp.Data.Endpoints) == 0 {
		printError("No endpoints returned", nil)
		return false
	}

	// Display endpoint details
	if verbose {
		fmt.Printf("\n   Model Details:\n")
		fmt.Printf("      Description: %s\n", truncateString(resp.Data.Description, 100))
		if resp.Data.Architecture.Tokenizer != nil {
			fmt.Printf("      Tokenizer: %s\n", *resp.Data.Architecture.Tokenizer)
		}
		if resp.Data.Architecture.InstructType != nil {
			fmt.Printf("      Instruct Type: %s\n", *resp.Data.Architecture.InstructType)
		}
		fmt.Printf("      Input Modalities: %v\n", resp.Data.Architecture.InputModalities)
		fmt.Printf("      Output Modalities: %v\n", resp.Data.Architecture.OutputModalities)

		fmt.Printf("\n   First 3 endpoints:\n")
		for i, endpoint := range resp.Data.Endpoints {
			if i >= 3 {
				break
			}
			fmt.Printf("      Endpoint %d:\n", i+1)
			fmt.Printf("         Provider: %s\n", endpoint.ProviderName)
			fmt.Printf("         Name: %s\n", endpoint.Name)
			fmt.Printf("         Status: %.0f\n", endpoint.Status)
			fmt.Printf("         Context Length: %.0f tokens\n", endpoint.ContextLength)
			if endpoint.MaxCompletionTokens != nil {
				fmt.Printf("         Max Completion Tokens: %.0f\n", *endpoint.MaxCompletionTokens)
			}
			if endpoint.Quantization != nil && *endpoint.Quantization != "" {
				fmt.Printf("         Quantization: %s\n", *endpoint.Quantization)
			}
			fmt.Printf("         Pricing - Prompt: $%s/M, Completion: $%s/M\n",
				endpoint.Pricing.Prompt, endpoint.Pricing.Completion)
			if endpoint.UptimeLast30m != nil {
				fmt.Printf("         Uptime (30m): %.2f%%\n", *endpoint.UptimeLast30m*100)
			}
			if len(endpoint.SupportedParameters) > 0 {
				fmt.Printf("         Supported Parameters: %d\n", len(endpoint.SupportedParameters))
			}
		}
	} else {
		// Non-verbose: just show a sample
		endpoint := resp.Data.Endpoints[0]
		fmt.Printf("      Example endpoint: %s (Provider: %s)\n",
			endpoint.Name, endpoint.ProviderName)
		fmt.Printf("      Pricing: $%s/M prompt, $%s/M completion\n",
			endpoint.Pricing.Prompt, endpoint.Pricing.Completion)
	}

	// Test 2: Validate endpoint structure
	fmt.Printf("\n   Validating endpoint data structure...\n")
	firstEndpoint := resp.Data.Endpoints[0]

	// Check required fields
	if firstEndpoint.Name == "" {
		printError("Endpoint missing Name", nil)
		return false
	}
	if firstEndpoint.ProviderName == "" {
		printError("Endpoint missing ProviderName", nil)
		return false
	}
	if firstEndpoint.ContextLength == 0 {
		printError("Endpoint missing ContextLength", nil)
		return false
	}
	// Status can be 0, 1, or other numeric values, so we don't validate it's non-zero

	// Check pricing
	if firstEndpoint.Pricing.Prompt == "" {
		printError("Endpoint missing Prompt pricing", nil)
		return false
	}
	if firstEndpoint.Pricing.Completion == "" {
		printError("Endpoint missing Completion pricing", nil)
		return false
	}

	printSuccess("Endpoint structure validation passed")

	// Test 3: List endpoints for Claude
	fmt.Printf("\n   Testing endpoints for Claude-3.5 Sonnet...\n")
	start = time.Now()
	claudeResp, err := client.ListModelEndpoints(ctx, "anthropic", "claude-3.5-sonnet")
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list Claude endpoints", err)
		return false
	}

	fmt.Printf("   ✅ Listed Claude endpoints (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", claudeResp.Data.Name)
	fmt.Printf("      Total endpoints: %d\n", len(claudeResp.Data.Endpoints))

	if verbose && len(claudeResp.Data.Endpoints) > 0 {
		fmt.Printf("      Example endpoint: %s\n", claudeResp.Data.Endpoints[0].ProviderName)
	}

	// Test 4: Test with invalid model (error handling)
	fmt.Printf("\n   Testing error handling with invalid model...\n")
	_, err = client.ListModelEndpoints(ctx, "invalid", "nonexistent-model")

	if err == nil {
		fmt.Printf("   ⚠️  Expected error but got success (model might exist)\n")
	} else {
		printSuccess("Error handling works correctly")
		if verbose {
			fmt.Printf("      Error: %v\n", err)
		}
	}

	// Test 5: Test with empty parameters
	fmt.Printf("\n   Testing validation with empty parameters...\n")

	_, err = client.ListModelEndpoints(ctx, "", "gpt-4")
	if err == nil {
		printError("Should have errored with empty author", nil)
		return false
	}
	printSuccess("Empty author validation passed")

	_, err = client.ListModelEndpoints(ctx, "openai", "")
	if err == nil {
		printError("Should have errored with empty slug", nil)
		return false
	}
	printSuccess("Empty slug validation passed")

	// Test 6: Compare pricing across endpoints
	if verbose && len(resp.Data.Endpoints) > 1 {
		fmt.Printf("\n   Pricing comparison for %s:\n", resp.Data.Name)
		fmt.Printf("   %-30s %-15s %-15s\n", "Provider", "Prompt/M", "Completion/M")
		fmt.Printf("   %s\n", strings.Repeat("-", 60))
		for i, endpoint := range resp.Data.Endpoints {
			if i >= 5 {
				fmt.Printf("   ... and %d more endpoints\n", len(resp.Data.Endpoints)-5)
				break
			}
			fmt.Printf("   %-30s $%-14s $%-14s\n",
				endpoint.ProviderName,
				endpoint.Pricing.Prompt,
				endpoint.Pricing.Completion,
			)
		}
	}

	printSuccess("Model endpoints tests completed")
	return true
}

// RunProvidersTest tests the ListProviders endpoint.
func RunProvidersTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Providers\n")

	// Test: List all providers
	fmt.Printf("   Testing list all providers...\n")
	start := time.Now()
	resp, err := client.ListProviders(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list providers", err)
		return false
	}

	fmt.Printf("   ✅ Listed all providers (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total providers: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		printError("No providers returned", nil)
		return false
	}

	// Display first few providers
	if verbose {
		fmt.Printf("\n   First 5 providers:\n")
		for i, provider := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, provider.Name, provider.Slug)
			if provider.PrivacyPolicyURL != nil {
				fmt.Printf("         Privacy Policy: %s\n", *provider.PrivacyPolicyURL)
			}
			if provider.TermsOfServiceURL != nil {
				fmt.Printf("         Terms of Service: %s\n", *provider.TermsOfServiceURL)
			}
			if provider.StatusPageURL != nil {
				fmt.Printf("         Status Page: %s\n", *provider.StatusPageURL)
			}
		}
	} else {
		// Show just a couple in non-verbose mode
		for i, provider := range resp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", provider.Name, provider.Slug)
		}
	}

	// Validate provider structure
	fmt.Printf("\n   Validating provider data structure...\n")
	firstProvider := resp.Data[0]

	// Check required fields
	if firstProvider.Name == "" {
		printError("Provider missing Name", nil)
		return false
	}
	if firstProvider.Slug == "" {
		printError("Provider missing Slug", nil)
		return false
	}

	printSuccess("Provider structure validation passed")

	if verbose {
		fmt.Printf("\n   First provider details:\n")
		fmt.Printf("      Name: %s\n", firstProvider.Name)
		fmt.Printf("      Slug: %s\n", firstProvider.Slug)
		if firstProvider.PrivacyPolicyURL != nil {
			fmt.Printf("      Privacy Policy URL: %s\n", *firstProvider.PrivacyPolicyURL)
		} else {
			fmt.Printf("      Privacy Policy URL: (not provided)\n")
		}
		if firstProvider.TermsOfServiceURL != nil {
			fmt.Printf("      Terms of Service URL: %s\n", *firstProvider.TermsOfServiceURL)
		} else {
			fmt.Printf("      Terms of Service URL: (not provided)\n")
		}
		if firstProvider.StatusPageURL != nil {
			fmt.Printf("      Status Page URL: %s\n", *firstProvider.StatusPageURL)
		} else {
			fmt.Printf("      Status Page URL: (not provided)\n")
		}
	}

	// Check for well-known providers
	fmt.Printf("\n   Checking for well-known providers...\n")
	wellKnownProviders := []string{
		"openai",
		"anthropic",
		"google",
		"meta",
	}

	foundProviders := make(map[string]bool)
	for _, provider := range resp.Data {
		for _, knownProvider := range wellKnownProviders {
			if provider.Slug == knownProvider {
				foundProviders[knownProvider] = true
			}
		}
	}

	foundCount := len(foundProviders)
	fmt.Printf("   Found %d/%d well-known providers\n", foundCount, len(wellKnownProviders))

	if verbose {
		for _, knownProvider := range wellKnownProviders {
			status := "❌"
			if foundProviders[knownProvider] {
				status = "✅"
			}
			fmt.Printf("      %s %s\n", status, knownProvider)
		}
	}

	// Verify policy URLs
	fmt.Printf("\n   Validating policy URLs...\n")
	hasPrivacyPolicy := 0
	hasTermsOfService := 0
	hasStatusPage := 0

	for _, provider := range resp.Data {
		if provider.PrivacyPolicyURL != nil && *provider.PrivacyPolicyURL != "" {
			hasPrivacyPolicy++
		}
		if provider.TermsOfServiceURL != nil && *provider.TermsOfServiceURL != "" {
			hasTermsOfService++
		}
		if provider.StatusPageURL != nil && *provider.StatusPageURL != "" {
			hasStatusPage++
		}
	}

	privacyPercent := (float64(hasPrivacyPolicy) / float64(len(resp.Data))) * 100
	termsPercent := (float64(hasTermsOfService) / float64(len(resp.Data))) * 100
	statusPercent := (float64(hasStatusPage) / float64(len(resp.Data))) * 100

	fmt.Printf("   %.1f%% have privacy policy (%d/%d)\n", privacyPercent, hasPrivacyPolicy, len(resp.Data))
	fmt.Printf("   %.1f%% have terms of service (%d/%d)\n", termsPercent, hasTermsOfService, len(resp.Data))
	fmt.Printf("   %.1f%% have status page (%d/%d)\n", statusPercent, hasStatusPage, len(resp.Data))

	printSuccess("List providers tests completed")
	return true
}
