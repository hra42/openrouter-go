package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunZDREndpointsTest tests the ListZDREndpoints endpoint.
func RunZDREndpointsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List ZDR Endpoints\n")

	// Test 1: List all ZDR endpoints
	fmt.Printf("   Testing list ZDR endpoints...\n")
	start := time.Now()
	resp, err := client.ListZDREndpoints(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list ZDR endpoints", err)
		return false
	}

	if resp == nil {
		printError("Response is nil", nil)
		return false
	}

	fmt.Printf("   ✅ Listed ZDR endpoints (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total endpoints: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		fmt.Printf("   ⚠️  No ZDR endpoints returned (this might be expected depending on account)\n")
		printSuccess("ZDR endpoints test completed")
		return true
	}

	// Display first few endpoints
	if verbose {
		fmt.Printf("\n   First 5 ZDR endpoints:\n")
		for i, ep := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (Model: %s)\n", i+1, ep.Name, ep.ModelName)
			fmt.Printf("         Model ID: %s\n", ep.ModelID)
			fmt.Printf("         Provider: %s\n", ep.ProviderName)
			fmt.Printf("         Context Length: %.0f tokens\n", ep.ContextLength)
			fmt.Printf("         Pricing: $%s/M prompt, $%s/M completion\n",
				ep.Pricing.Prompt, ep.Pricing.Completion)
			if ep.Tag != nil {
				fmt.Printf("         Tag: %s\n", *ep.Tag)
			}
			if ep.SupportsImplicitCaching != nil {
				fmt.Printf("         Supports Implicit Caching: %v\n", *ep.SupportsImplicitCaching)
			}
			if ep.LatencyLast30m != nil {
				fmt.Printf("         Latency P50: %.1fms, P99: %.1fms\n",
					ep.LatencyLast30m.P50, ep.LatencyLast30m.P99)
			}
			if ep.ThroughputLast30m != nil {
				fmt.Printf("         Throughput P50: %.1f, P99: %.1f\n",
					ep.ThroughputLast30m.P50, ep.ThroughputLast30m.P99)
			}
		}
	} else {
		for i, ep := range resp.Data {
			if i >= 3 {
				break
			}
			fmt.Printf("      Example: %s (%s)\n", ep.Name, ep.ProviderName)
		}
	}

	// Test 2: Validate endpoint structure
	fmt.Printf("\n   Validating endpoint data structure...\n")
	firstEndpoint := resp.Data[0]

	if firstEndpoint.Name == "" {
		printError("Endpoint missing Name", nil)
		return false
	}
	if firstEndpoint.ModelID == "" {
		printError("Endpoint missing ModelID", nil)
		return false
	}
	if firstEndpoint.ProviderName == "" {
		printError("Endpoint missing ProviderName", nil)
		return false
	}
	if firstEndpoint.Pricing.Prompt == "" {
		printError("Endpoint missing Prompt pricing", nil)
		return false
	}
	if firstEndpoint.Pricing.Completion == "" {
		printError("Endpoint missing Completion pricing", nil)
		return false
	}

	printSuccess("Endpoint structure validation passed")

	printSuccess("ZDR endpoints test completed")
	return true
}
