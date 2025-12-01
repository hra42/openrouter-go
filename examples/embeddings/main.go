package main

import (
	"context"
	"fmt"
	"log"
	"math"
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

	// Example 1: Single text embedding
	fmt.Println("=== Example 1: Single Text Embedding ===")
	singleEmbedding(client)

	// Example 2: Batch embeddings
	fmt.Println("\n=== Example 2: Batch Embeddings ===")
	batchEmbeddings(client)

	// Example 3: Embedding with options
	fmt.Println("\n=== Example 3: Embedding with Options ===")
	embeddingWithOptions(client)

	// Example 4: Semantic similarity
	fmt.Println("\n=== Example 4: Semantic Similarity ===")
	semanticSimilarity(client)

	// Example 5: List available embedding models
	fmt.Println("\n=== Example 5: List Embedding Models ===")
	listEmbeddingModels(client)
}

func singleEmbedding(client *openrouter.Client) {
	text := "The quick brown fox jumps over the lazy dog."

	resp, err := client.CreateEmbedding(
		context.Background(),
		text,
		"openai/text-embedding-3-small",
	)
	if err != nil {
		log.Printf("Error creating embedding: %v", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Embeddings returned: %d\n", len(resp.Data))

	if len(resp.Data) > 0 {
		vec := resp.Data[0].GetEmbeddingVector()
		if vec != nil {
			fmt.Printf("Embedding dimensions: %d\n", len(vec))
			fmt.Printf("First 5 values: [%.4f, %.4f, %.4f, %.4f, %.4f...]\n",
				vec[0], vec[1], vec[2], vec[3], vec[4])
		}
	}

	if resp.Usage != nil {
		fmt.Printf("Prompt tokens: %d\n", resp.Usage.PromptTokens)
		fmt.Printf("Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}

func batchEmbeddings(client *openrouter.Client) {
	texts := []string{
		"Machine learning is fascinating.",
		"Deep learning uses neural networks.",
		"Natural language processing handles text.",
	}

	resp, err := client.CreateEmbeddings(
		context.Background(),
		texts,
		"openai/text-embedding-3-small",
	)
	if err != nil {
		log.Printf("Error creating batch embeddings: %v", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Embeddings returned: %d\n", len(resp.Data))

	for i, embedding := range resp.Data {
		vec := embedding.GetEmbeddingVector()
		if vec != nil {
			fmt.Printf("  Embedding %d: %d dimensions\n", i, len(vec))
		}
	}

	if resp.Usage != nil {
		fmt.Printf("Total tokens: %d\n", resp.Usage.TotalTokens)
	}
}

func embeddingWithOptions(client *openrouter.Client) {
	text := "This is a query for semantic search."

	resp, err := client.CreateEmbedding(
		context.Background(),
		text,
		"openai/text-embedding-3-small",
		openrouter.WithEmbeddingUser("user-123"),
		openrouter.WithEmbeddingInputType("query"),
	)
	if err != nil {
		log.Printf("Error creating embedding with options: %v", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	if len(resp.Data) > 0 {
		vec := resp.Data[0].GetEmbeddingVector()
		if vec != nil {
			fmt.Printf("Embedding dimensions: %d\n", len(vec))
		}
	}
}

func semanticSimilarity(client *openrouter.Client) {
	texts := []string{
		"I love programming in Go.",
		"Go is my favorite programming language.",
		"The weather is nice today.",
	}

	resp, err := client.CreateEmbeddings(
		context.Background(),
		texts,
		"openai/text-embedding-3-small",
	)
	if err != nil {
		log.Printf("Error creating embeddings: %v", err)
		return
	}

	if len(resp.Data) < 3 {
		log.Printf("Expected 3 embeddings, got %d", len(resp.Data))
		return
	}

	// Get embedding vectors
	vec1 := resp.Data[0].GetEmbeddingVector()
	vec2 := resp.Data[1].GetEmbeddingVector()
	vec3 := resp.Data[2].GetEmbeddingVector()

	if vec1 == nil || vec2 == nil || vec3 == nil {
		log.Printf("Failed to get embedding vectors")
		return
	}

	// Calculate cosine similarities
	sim12 := cosineSimilarity(vec1, vec2)
	sim13 := cosineSimilarity(vec1, vec3)
	sim23 := cosineSimilarity(vec2, vec3)

	fmt.Printf("Text 1: %q\n", texts[0])
	fmt.Printf("Text 2: %q\n", texts[1])
	fmt.Printf("Text 3: %q\n", texts[2])
	fmt.Println()
	fmt.Printf("Similarity (Text 1 vs Text 2): %.4f (related topics)\n", sim12)
	fmt.Printf("Similarity (Text 1 vs Text 3): %.4f (unrelated topics)\n", sim13)
	fmt.Printf("Similarity (Text 2 vs Text 3): %.4f (unrelated topics)\n", sim23)
}

func listEmbeddingModels(client *openrouter.Client) {
	resp, err := client.ListEmbeddingsModels(context.Background())
	if err != nil {
		log.Printf("Error listing embedding models: %v", err)
		return
	}

	fmt.Printf("Available embedding models: %d\n\n", len(resp.Data))

	for i, model := range resp.Data {
		if i >= 10 {
			fmt.Printf("... and %d more models\n", len(resp.Data)-10)
			break
		}
		fmt.Printf("%d. %s (%s)\n", i+1, model.Name, model.ID)
		if model.ContextLength != nil {
			fmt.Printf("   Context: %.0f tokens\n", *model.ContextLength)
		}
		fmt.Printf("   Pricing: $%s per million tokens\n", model.Pricing.Prompt)
	}
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}
