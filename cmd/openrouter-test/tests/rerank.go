package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunRerankTest tests the Rerank endpoint.
func RunRerankTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Rerank\n")

	query := "What is the capital of France?"
	documents := []string{
		"Berlin is the capital of Germany.",
		"Paris is the capital of France.",
		"Madrid is the capital of Spain.",
		"Rome is the capital of Italy.",
		"London is the capital of the United Kingdom.",
	}

	fmt.Printf("   Query: %q\n", query)
	fmt.Printf("   Documents: %d\n", len(documents))

	// Test 1: Basic rerank
	fmt.Printf("   Testing basic rerank...\n")
	start := time.Now()
	resp, err := client.Rerank(ctx, query, documents, model)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to rerank", err)
		return false
	}

	fmt.Printf("   ✅ Reranked documents (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", resp.Model)

	if len(resp.Results) == 0 {
		printError("No results returned", nil)
		return false
	}

	fmt.Printf("      Results: %d\n", len(resp.Results))

	// Verify result structure
	for i, result := range resp.Results {
		if result.Document.Text == "" {
			printError(fmt.Sprintf("Empty document text at result %d", i), nil)
			return false
		}
		if verbose {
			fmt.Printf("      %d. [%.4f] (index=%d) %s\n", i+1, result.RelevanceScore, result.Index, truncateString(result.Document.Text, 60))
		}
	}

	// Display usage if available
	if resp.Usage != nil {
		if resp.Usage.TotalTokens > 0 {
			fmt.Printf("      Total tokens: %d\n", resp.Usage.TotalTokens)
		}
		if resp.Usage.SearchUnits > 0 {
			fmt.Printf("      Search units: %d\n", resp.Usage.SearchUnits)
		}
		if resp.Usage.Cost != nil {
			fmt.Printf("      Cost: $%.6f\n", *resp.Usage.Cost)
		}
	}

	if resp.Provider != "" {
		fmt.Printf("      Provider: %s\n", resp.Provider)
	}

	// Test 2: Rerank with top_n
	fmt.Printf("   Testing rerank with top_n=2...\n")
	start = time.Now()
	resp2, err := client.Rerank(ctx, query, documents, model,
		openrouter.WithRerankTopN(2),
	)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to rerank with top_n", err)
		return false
	}

	fmt.Printf("   ✅ Reranked with top_n=2 (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Results returned: %d\n", len(resp2.Results))

	if len(resp2.Results) > 2 {
		printError(fmt.Sprintf("Expected at most 2 results with top_n=2, got %d", len(resp2.Results)), nil)
		return false
	}

	if verbose {
		for i, result := range resp2.Results {
			fmt.Printf("      %d. [%.4f] %s\n", i+1, result.RelevanceScore, truncateString(result.Document.Text, 60))
		}
	}

	printSuccess("Rerank tests passed")
	return true
}
