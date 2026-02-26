package tests

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunEmbeddingTest tests the CreateEmbedding endpoint.
func RunEmbeddingTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Create Embedding\n")

	// Test 1: Single text embedding
	fmt.Printf("   Testing single text embedding...\n")
	start := time.Now()
	resp, err := client.CreateEmbedding(ctx, "Hello, world!", model)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to create embedding", err)
		return false
	}

	fmt.Printf("   ✅ Created embedding (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", resp.Model)
	fmt.Printf("      Object: %s\n", resp.Object)
	fmt.Printf("      Embeddings returned: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		printError("No embeddings returned", nil)
		return false
	}

	// Verify embedding structure
	embedding := resp.Data[0]
	if embedding.Object != "embedding" {
		printError("Expected object 'embedding'", nil)
		return false
	}
	if embedding.Index != 0 {
		printError("Expected index 0", nil)
		return false
	}

	// Get embedding vector
	vec := embedding.GetEmbeddingVector()
	if vec == nil {
		printError("Failed to get embedding vector", nil)
		return false
	}

	fmt.Printf("      Embedding dimensions: %d\n", len(vec))

	if verbose {
		fmt.Printf("      First 5 dimensions: [")
		for i := 0; i < 5 && i < len(vec); i++ {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%.4f", vec[i])
		}
		fmt.Printf("...]\n")
	}

	// Display usage if available
	if resp.Usage != nil {
		fmt.Printf("      Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("      Total tokens: %d\n", resp.Usage.TotalTokens)
		if resp.Usage.Cost != nil {
			fmt.Printf("      Cost: $%.6f\n", *resp.Usage.Cost)
		}
	}

	printSuccess("Single text embedding test passed")
	return true
}

// RunBatchEmbeddingTest tests batch embeddings with CreateEmbeddings.
func RunBatchEmbeddingTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Batch Embeddings\n")

	// Test batch embedding
	inputs := []string{
		"The quick brown fox jumps over the lazy dog.",
		"Machine learning is a subset of artificial intelligence.",
		"OpenRouter provides access to many AI models.",
	}

	fmt.Printf("   Testing batch embedding (%d texts)...\n", len(inputs))
	start := time.Now()
	resp, err := client.CreateEmbeddings(ctx, inputs, model)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to create batch embeddings", err)
		return false
	}

	fmt.Printf("   ✅ Created batch embeddings (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", resp.Model)
	fmt.Printf("      Embeddings returned: %d\n", len(resp.Data))

	if len(resp.Data) != len(inputs) {
		printError(fmt.Sprintf("Expected %d embeddings, got %d", len(inputs), len(resp.Data)), nil)
		return false
	}

	// Verify all embeddings
	for i, embedding := range resp.Data {
		if embedding.Index != i {
			printError(fmt.Sprintf("Expected index %d, got %d", i, embedding.Index), nil)
			return false
		}
		vec := embedding.GetEmbeddingVector()
		if len(vec) == 0 {
			printError(fmt.Sprintf("Empty embedding vector at index %d", i), nil)
			return false
		}
		if verbose {
			fmt.Printf("      Embedding %d: %d dimensions\n", i, len(vec))
		}
	}

	// Display usage
	if resp.Usage != nil {
		fmt.Printf("      Total prompt tokens: %d\n", resp.Usage.PromptTokens)
	}

	printSuccess("Batch embeddings test passed")
	return true
}

// RunEmbeddingWithOptionsTest tests embedding options.
func RunEmbeddingWithOptionsTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Embedding with Options\n")

	// Test with user option
	fmt.Printf("   Testing embedding with user option...\n")
	start := time.Now()
	resp, err := client.CreateEmbedding(
		ctx,
		"Test embedding with options",
		model,
		openrouter.WithEmbeddingUser("test-user-123"),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to create embedding with options", err)
		return false
	}

	fmt.Printf("   ✅ Created embedding with options (%.2fs)\n", elapsed.Seconds())

	if len(resp.Data) == 0 {
		printError("No embeddings returned", nil)
		return false
	}

	vec := resp.Data[0].GetEmbeddingVector()
	if vec == nil {
		printError("Failed to get embedding vector", nil)
		return false
	}

	fmt.Printf("      Embedding dimensions: %d\n", len(vec))

	printSuccess("Embedding with options test passed")
	return true
}

// RunListEmbeddingsModelsTest tests the ListEmbeddingsModels endpoint.
func RunListEmbeddingsModelsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Embeddings Models\n")

	// Test: List all embedding models
	fmt.Printf("   Testing list embedding models...\n")
	start := time.Now()
	resp, err := client.ListEmbeddingsModels(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list embedding models", err)
		return false
	}

	fmt.Printf("   ✅ Listed embedding models (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total models: %d\n", len(resp.Data))

	if len(resp.Data) == 0 {
		printError("No embedding models returned", nil)
		return false
	}

	// Display first few models
	if verbose {
		fmt.Printf("\n   First 5 embedding models:\n")
		for i, model := range resp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (%s)\n", i+1, model.Name, model.ID)
			if model.ContextLength != nil {
				fmt.Printf("         Context: %.0f tokens\n", *model.ContextLength)
			}
			fmt.Printf("         Pricing: $%s/M tokens\n", model.Pricing.Prompt)
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

	// Validate model structure
	fmt.Printf("\n   Validating embedding model data structure...\n")
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

	// Check that output modality includes embeddings
	hasEmbeddingsOutput := slices.Contains(firstModel.Architecture.OutputModalities, "embeddings")
	if !hasEmbeddingsOutput && len(firstModel.Architecture.OutputModalities) > 0 {
		fmt.Printf("   ⚠️  First model doesn't have 'embeddings' output modality: %v\n",
			firstModel.Architecture.OutputModalities)
	}

	printSuccess("Embedding model structure validation passed")

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
		fmt.Printf("      Pricing (prompt): $%s per million tokens\n", firstModel.Pricing.Prompt)
	}

	// Check for well-known embedding models
	fmt.Printf("\n   Checking for well-known embedding models...\n")
	wellKnownModels := []string{
		"openai/text-embedding-3-small",
		"openai/text-embedding-3-large",
		"openai/text-embedding-ada-002",
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
	fmt.Printf("   Found %d/%d well-known embedding models\n", foundCount, len(wellKnownModels))

	if verbose {
		for _, knownModel := range wellKnownModels {
			status := "❌"
			if foundModels[knownModel] {
				status = "✅"
			}
			fmt.Printf("      %s %s\n", status, knownModel)
		}
	}

	printSuccess("List embeddings models tests completed")
	return true
}
