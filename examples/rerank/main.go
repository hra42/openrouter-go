// Package main demonstrates the Rerank API endpoint.
//
// This example shows how to:
// - Rerank documents by relevance to a search query
// - Use the top_n option to limit results
// - Interpret relevance scores
//
// Usage:
//
//	export OPENROUTER_API_KEY="your-api-key"
//	go run examples/rerank/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	// Example 1: Basic reranking
	fmt.Println("=== Example 1: Basic Reranking ===")
	fmt.Println()
	basicRerank(ctx, client)

	// Example 2: Reranking with top_n
	fmt.Println()
	fmt.Println("=== Example 2: Reranking with top_n ===")
	fmt.Println()
	rerankWithTopN(ctx, client)
}

func basicRerank(ctx context.Context, client *openrouter.Client) {
	query := "What programming language is best for web development?"
	documents := []string{
		"Python is a versatile language used for data science, machine learning, and web development with frameworks like Django and Flask.",
		"JavaScript is the most popular language for web development, running in browsers and on servers via Node.js.",
		"Rust is a systems programming language focused on safety, speed, and concurrency.",
		"Go is designed for building scalable network services and cloud infrastructure.",
		"TypeScript adds static typing to JavaScript and is widely used in modern web applications.",
	}

	fmt.Printf("Query: %q\n", query)
	fmt.Printf("Documents: %d\n\n", len(documents))

	resp, err := client.Rerank(ctx, query, documents, "cohere/rerank-v3.5")
	if err != nil {
		log.Fatalf("Failed to rerank: %v", err)
	}

	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Println("Results (sorted by relevance):")
	for i, result := range resp.Results {
		fmt.Printf("  %d. [score: %.4f] %s\n", i+1, result.RelevanceScore, documents[result.Index])
	}

	if resp.Usage != nil {
		fmt.Printf("\nUsage: %d tokens", resp.Usage.TotalTokens)
		if resp.Usage.SearchUnits > 0 {
			fmt.Printf(", %d search units", resp.Usage.SearchUnits)
		}
		fmt.Println()
	}
}

func rerankWithTopN(ctx context.Context, client *openrouter.Client) {
	query := "How do neural networks work?"
	documents := []string{
		"Neural networks are composed of layers of interconnected nodes that process information.",
		"The stock market experienced significant volatility this quarter.",
		"Backpropagation is the algorithm used to train neural networks by adjusting weights.",
		"Mediterranean cuisine features olive oil, fresh vegetables, and seafood.",
		"Deep learning uses multiple layers in neural networks to learn hierarchical features.",
		"The weather forecast predicts rain for the weekend.",
	}

	fmt.Printf("Query: %q\n", query)
	fmt.Printf("Documents: %d (returning top 3)\n\n", len(documents))

	resp, err := client.Rerank(ctx, query, documents, "cohere/rerank-v3.5",
		openrouter.WithRerankTopN(3),
	)
	if err != nil {
		log.Fatalf("Failed to rerank: %v", err)
	}

	fmt.Println("Top 3 most relevant documents:")
	for i, result := range resp.Results {
		fmt.Printf("  %d. [score: %.4f] (original index: %d) %s\n",
			i+1, result.RelevanceScore, result.Index, documents[result.Index])
	}
}
