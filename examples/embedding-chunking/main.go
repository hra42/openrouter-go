// Package main demonstrates text chunking and chunked embeddings.
//
// This example shows how to:
// - Chunk long documents using different strategies
// - Create embeddings for document chunks
// - Combine chunk embeddings for document similarity
//
// Usage:
//
//	export OPENROUTER_API_KEY="your-api-key"
//	go run examples/embedding-chunking/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

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

	// Example 1: Basic text chunking with different strategies
	fmt.Println("=== Example 1: Text Chunking Strategies ===")
	fmt.Println()
	demonstrateChunking()

	// Example 2: Chunked embeddings for long documents
	fmt.Println()
	fmt.Println("=== Example 2: Chunked Embeddings ===")
	fmt.Println()
	demonstrateChunkedEmbeddings(ctx, client)

	// Example 3: Document similarity using chunked embeddings
	fmt.Println()
	fmt.Println("=== Example 3: Document Similarity ===")
	fmt.Println()
	demonstrateDocumentSimilarity(ctx, client)
}

func demonstrateChunking() {
	// Sample long text (article about AI)
	article := `Artificial Intelligence: A Brief Overview

Artificial intelligence (AI) is transforming industries across the globe. From healthcare to finance,
AI systems are automating tasks, providing insights, and enabling new capabilities that were previously impossible.

Machine Learning Fundamentals
Machine learning is a subset of AI that focuses on training systems to learn from data.
Supervised learning uses labeled datasets to train models for classification and regression.
Unsupervised learning discovers hidden patterns in data without explicit labels.
Reinforcement learning enables agents to learn optimal behaviors through trial and error.

Deep Learning Revolution
Deep learning uses neural networks with many layers to learn hierarchical representations.
Convolutional neural networks excel at image and video analysis tasks.
Recurrent networks and transformers have revolutionized natural language processing.
Large language models can now generate coherent text and engage in complex reasoning.

Future Directions
The future of AI includes more efficient models, better interpretability, and responsible development.
Researchers are working on multimodal AI that can process text, images, and audio together.
Edge AI brings intelligence to devices, enabling real-time processing without cloud connectivity.`

	// Strategy 1: Paragraph-based chunking (recommended for articles)
	fmt.Println("1. Paragraph-based Chunking:")
	paragraphConfig := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkByParagraphs,
		ChunkSize:      200, // Target tokens per chunk
		Overlap:        1,   // Overlap 1 paragraph for context
		TrimWhitespace: true,
	}

	chunks, _ := openrouter.ChunkText(article, paragraphConfig)
	fmt.Printf("   Created %d chunks\n", len(chunks))
	for i, chunk := range chunks {
		preview := chunk.Text
		if len(preview) > 60 {
			preview = preview[:60] + "..."
		}
		fmt.Printf("   Chunk %d: %q\n", i, preview)
	}

	// Strategy 2: Section-based chunking (for markdown documents)
	fmt.Println("\n2. Section-based Chunking:")
	sectionConfig := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySections,
		ChunkSize:      300,
		TrimWhitespace: true,
	}

	chunks2, _ := openrouter.ChunkText(article, sectionConfig)
	fmt.Printf("   Created %d chunks\n", len(chunks2))
	for i, chunk := range chunks2 {
		lines := strings.SplitN(chunk.Text, "\n", 2)
		fmt.Printf("   Chunk %d: %q\n", i, lines[0])
	}

	// Strategy 3: Sentence-based chunking (fine-grained)
	fmt.Println("\n3. Sentence-based Chunking:")
	sentenceConfig := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySentences,
		ChunkSize:      100,
		Overlap:        2, // Overlap 2 sentences
		TrimWhitespace: true,
	}

	chunks3, _ := openrouter.ChunkText(article, sentenceConfig)
	fmt.Printf("   Created %d chunks\n", len(chunks3))

	// Token estimation
	fmt.Println("\n4. Token Estimation:")
	tokens := openrouter.EstimateTokens(article)
	tokensFromWords := openrouter.EstimateTokensFromWords(article)
	fmt.Printf("   Article length: %d characters\n", len(article))
	fmt.Printf("   Estimated tokens (char-based): %d\n", tokens)
	fmt.Printf("   Estimated tokens (word-based): %d\n", tokensFromWords)
}

func demonstrateChunkedEmbeddings(ctx context.Context, client *openrouter.Client) {
	document := `The history of computing spans several decades of innovation and discovery.
Early computers filled entire rooms and used vacuum tubes for processing.
The invention of the transistor revolutionized electronics and made computers smaller.
Integrated circuits further miniaturized computing power onto tiny chips.

The personal computer revolution brought computing to homes and offices.
Companies like Apple, IBM, and Microsoft made computers accessible to everyone.
Graphical user interfaces replaced command-line interactions.
The mouse became a standard input device alongside the keyboard.

The internet connected computers worldwide and transformed communication.
Email replaced letters for many forms of correspondence.
The World Wide Web made information accessible to billions of people.
Social media platforms changed how people interact and share content.

Mobile computing put powerful computers in everyone's pockets.
Smartphones combine communication, computing, and entertainment.
Apps provide specialized functionality for countless use cases.
Cloud computing enables processing and storage beyond device limitations.`

	config := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkByParagraphs,
		ChunkSize:      100,
		Overlap:        1,
		TrimWhitespace: true,
	}

	fmt.Println("Creating chunked embeddings...")
	result, err := client.CreateChunkedEmbedding(
		ctx,
		document,
		"openai/text-embedding-3-small",
		config,
	)
	if err != nil {
		log.Printf("Failed to create chunked embedding: %v", err)
		return
	}

	fmt.Printf("Model: %s\n", result.Model)
	fmt.Printf("Chunks: %d\n", len(result.Chunks))
	fmt.Printf("Total tokens: %d\n", result.TotalTokensUsed)

	// Show each chunk's embedding info
	for i, chunkEmb := range result.Chunks {
		preview := chunkEmb.Chunk.Text
		if len(preview) > 50 {
			preview = preview[:50] + "..."
		}
		fmt.Printf("\nChunk %d: %d dimensions\n", i, len(chunkEmb.Embedding))
		fmt.Printf("  Text: %q\n", preview)
		fmt.Printf("  First 3 values: [%.4f, %.4f, %.4f, ...]\n",
			chunkEmb.Embedding[0], chunkEmb.Embedding[1], chunkEmb.Embedding[2])
	}

	// Create document-level embedding by averaging
	var embeddings [][]float64
	for _, chunkEmb := range result.Chunks {
		embeddings = append(embeddings, chunkEmb.Embedding)
	}

	docEmbedding := openrouter.AverageEmbeddings(embeddings)
	fmt.Printf("\nDocument embedding: %d dimensions\n", len(docEmbedding))
	fmt.Printf("First 3 values: [%.4f, %.4f, %.4f, ...]\n",
		docEmbedding[0], docEmbedding[1], docEmbedding[2])
}

func demonstrateDocumentSimilarity(ctx context.Context, client *openrouter.Client) {
	// Two documents about related topics
	doc1 := `Machine learning is revolutionizing how we analyze data.
Neural networks can learn complex patterns from examples.
Deep learning has enabled breakthroughs in image recognition.
Natural language processing allows computers to understand text.`

	doc2 := `Artificial intelligence mimics human cognitive functions.
Deep neural networks process information in hierarchical layers.
Computer vision systems can identify objects in images.
Language models generate human-like text responses.`

	// An unrelated document
	doc3 := `The culinary arts involve preparing and presenting food.
Chefs use various techniques like sauteing, roasting, and grilling.
Fresh ingredients are essential for quality dishes.
Presentation affects the dining experience significantly.`

	config := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySentences,
		ChunkSize:      50,
		TrimWhitespace: true,
	}

	// Create embeddings for all documents
	fmt.Println("Creating embeddings for 3 documents...")

	result1, err := client.CreateChunkedEmbedding(ctx, doc1, "openai/text-embedding-3-small", config)
	if err != nil {
		log.Printf("Failed to create embedding for doc1: %v", err)
		return
	}

	result2, err := client.CreateChunkedEmbedding(ctx, doc2, "openai/text-embedding-3-small", config)
	if err != nil {
		log.Printf("Failed to create embedding for doc2: %v", err)
		return
	}

	result3, err := client.CreateChunkedEmbedding(ctx, doc3, "openai/text-embedding-3-small", config)
	if err != nil {
		log.Printf("Failed to create embedding for doc3: %v", err)
		return
	}

	// Average embeddings for each document
	emb1 := averageChunkEmbeddings(result1)
	emb2 := averageChunkEmbeddings(result2)
	emb3 := averageChunkEmbeddings(result3)

	// Calculate similarities
	sim12 := openrouter.CosineSimilarity(emb1, emb2)
	sim13 := openrouter.CosineSimilarity(emb1, emb3)
	sim23 := openrouter.CosineSimilarity(emb2, emb3)

	fmt.Printf("\nDocument 1: AI/ML topic (%d chunks)\n", len(result1.Chunks))
	fmt.Printf("Document 2: AI/ML topic (%d chunks)\n", len(result2.Chunks))
	fmt.Printf("Document 3: Cooking topic (%d chunks)\n", len(result3.Chunks))

	fmt.Println("\nSimilarity Matrix:")
	fmt.Printf("  Doc1 vs Doc2 (both AI): %.4f\n", sim12)
	fmt.Printf("  Doc1 vs Doc3 (AI vs Cooking): %.4f\n", sim13)
	fmt.Printf("  Doc2 vs Doc3 (AI vs Cooking): %.4f\n", sim23)

	fmt.Println("\nInterpretation:")
	if sim12 > sim13 && sim12 > sim23 {
		fmt.Println("  As expected, the two AI documents are most similar to each other.")
	}
	if sim13 < 0.5 && sim23 < 0.5 {
		fmt.Println("  The cooking document has low similarity with both AI documents.")
	}
}

func averageChunkEmbeddings(result *openrouter.ChunkedEmbeddingResult) []float64 {
	var embeddings [][]float64
	for _, chunkEmb := range result.Chunks {
		embeddings = append(embeddings, chunkEmb.Embedding)
	}
	return openrouter.AverageEmbeddings(embeddings)
}
