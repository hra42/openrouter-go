package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunChunkingTest tests the text chunking utilities.
func RunChunkingTest(ctx context.Context, client *openrouter.Client, embeddingModel string, verbose bool) bool {
	fmt.Printf("🔄 Test: Text Chunking\n")

	// Test 1: Paragraph chunking
	fmt.Printf("   Testing paragraph-based chunking...\n")

	longText := `Artificial intelligence has transformed how we interact with technology.
Machine learning models can now understand and generate human language with remarkable accuracy.

Natural language processing enables computers to analyze text, extract meaning, and respond appropriately.
This has led to significant advances in chatbots, virtual assistants, and automated customer service.

Deep learning architectures like transformers have revolutionized the field.
Models like GPT and BERT have achieved state-of-the-art results on many benchmarks.

The future of AI promises even more exciting developments.
As models become more capable, we can expect new applications across industries.`

	config := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkByParagraphs,
		ChunkSize:      100,
		Overlap:        0,
		TrimWhitespace: true,
	}

	chunks, err := openrouter.ChunkText(longText, config)
	if err != nil {
		printError("Failed to chunk text by paragraphs", err)
		return false
	}

	fmt.Printf("   ✅ Created %d chunks from paragraph-based chunking\n", len(chunks))

	if len(chunks) == 0 {
		printError("No chunks created", nil)
		return false
	}

	if verbose {
		for i, chunk := range chunks {
			preview := chunk.Text
			if len(preview) > 60 {
				preview = preview[:60] + "..."
			}
			fmt.Printf("      Chunk %d: %q\n", i, preview)
		}
	}

	// Test 2: Sentence chunking
	fmt.Printf("   Testing sentence-based chunking...\n")

	config2 := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySentences,
		ChunkSize:      50,
		TrimWhitespace: true,
	}

	chunks2, err := openrouter.ChunkText(longText, config2)
	if err != nil {
		printError("Failed to chunk text by sentences", err)
		return false
	}

	fmt.Printf("   ✅ Created %d chunks from sentence-based chunking\n", len(chunks2))

	// Test 3: Section-based chunking with markdown
	fmt.Printf("   Testing section-based chunking (markdown)...\n")

	markdownText := `# Introduction
This is the introduction to our documentation.

## Getting Started
To get started, install the package using your package manager.

### Prerequisites
You will need Go 1.18 or later installed.

### Installation
Run the following command to install.

## Usage
Here is how to use the library in your projects.

### Basic Example
Create a client and make your first request.

## API Reference
Complete API documentation follows.`

	config3 := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySections,
		ChunkSize:      100,
		TrimWhitespace: true,
	}

	chunks3, err := openrouter.ChunkText(markdownText, config3)
	if err != nil {
		printError("Failed to chunk text by sections", err)
		return false
	}

	fmt.Printf("   ✅ Created %d chunks from section-based chunking\n", len(chunks3))

	if verbose {
		for i, chunk := range chunks3 {
			// Show first line (header)
			lines := strings.SplitN(chunk.Text, "\n", 2)
			fmt.Printf("      Chunk %d: %q\n", i, lines[0])
		}
	}

	// Test 4: Token estimation
	fmt.Printf("   Testing token estimation...\n")

	testText := "This is a test sentence with multiple words for token estimation."
	tokens := openrouter.EstimateTokens(testText)
	tokensFromWords := openrouter.EstimateTokensFromWords(testText)

	fmt.Printf("   ✅ Token estimates: %d (char-based), %d (word-based)\n", tokens, tokensFromWords)

	// Test 5: Character-based chunking with overlap
	fmt.Printf("   Testing character-based chunking with overlap...\n")

	config4 := openrouter.ChunkConfig{
		Strategy:      openrouter.ChunkByCharacters,
		ChunkSize:     100,
		Overlap:       20,
		PreserveWords: true,
	}

	chunks4, err := openrouter.ChunkText(longText, config4)
	if err != nil {
		printError("Failed to chunk text by characters", err)
		return false
	}

	fmt.Printf("   ✅ Created %d chunks with character-based chunking (overlap: 20)\n", len(chunks4))

	// Test 6: Embedding utilities
	fmt.Printf("   Testing embedding utilities...\n")

	emb1 := []float64{1.0, 0.0, 0.0}
	emb2 := []float64{0.0, 1.0, 0.0}
	emb3 := []float64{0.0, 0.0, 1.0}

	avg := openrouter.AverageEmbeddings([][]float64{emb1, emb2, emb3})
	if avg == nil || len(avg) != 3 {
		printError("AverageEmbeddings failed", nil)
		return false
	}
	fmt.Printf("   ✅ AverageEmbeddings computed: [%.2f, %.2f, %.2f]\n", avg[0], avg[1], avg[2])

	similarity := openrouter.CosineSimilarity(emb1, emb1)
	if similarity != 1.0 {
		printError(fmt.Sprintf("Expected cosine similarity 1.0, got %.4f", similarity), nil)
		return false
	}
	fmt.Printf("   ✅ CosineSimilarity (identical): %.4f\n", similarity)

	similarity2 := openrouter.CosineSimilarity(emb1, emb2)
	if similarity2 != 0.0 {
		printError(fmt.Sprintf("Expected cosine similarity 0.0 for orthogonal, got %.4f", similarity2), nil)
		return false
	}
	fmt.Printf("   ✅ CosineSimilarity (orthogonal): %.4f\n", similarity2)

	printSuccess("Text chunking tests passed")
	return true
}

// RunChunkedEmbeddingTest tests chunked embeddings with the API.
func RunChunkedEmbeddingTest(ctx context.Context, client *openrouter.Client, embeddingModel string, verbose bool) bool {
	fmt.Printf("🔄 Test: Chunked Embeddings\n")

	// Create a long document for chunking
	longDocument := `Artificial Intelligence and Machine Learning: A Comprehensive Overview

Introduction
Artificial intelligence (AI) represents one of the most transformative technologies of our time.
It encompasses a broad range of techniques that enable machines to mimic human intelligence.
From simple rule-based systems to complex neural networks, AI has evolved significantly over decades.

Machine Learning Fundamentals
Machine learning is a subset of AI that focuses on learning from data.
Rather than being explicitly programmed, ML systems improve through experience.
Supervised learning uses labeled data to train models for classification and regression tasks.
Unsupervised learning discovers patterns in unlabeled data through clustering and dimensionality reduction.
Reinforcement learning enables agents to learn optimal behaviors through trial and error.

Deep Learning Revolution
Deep learning has revolutionized AI through the use of neural networks with many layers.
Convolutional neural networks excel at image recognition and computer vision tasks.
Recurrent neural networks and transformers have advanced natural language processing.
Large language models can now generate human-like text and engage in complex conversations.

Applications Across Industries
Healthcare uses AI for diagnosis, drug discovery, and personalized treatment plans.
Finance applies ML for fraud detection, algorithmic trading, and risk assessment.
Transportation benefits from autonomous vehicles and optimized logistics.
Manufacturing leverages AI for quality control, predictive maintenance, and robotics.

Future Directions
The field continues to advance with new architectures and training techniques.
Multimodal AI systems can process and generate multiple types of content.
Efficient AI aims to reduce computational requirements while maintaining performance.
Ethical AI focuses on fairness, transparency, and responsible development.`

	// Test 1: Create chunked embedding with paragraph strategy
	fmt.Printf("   Testing chunked embedding with paragraph strategy...\n")

	config := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkByParagraphs,
		ChunkSize:      200,
		Overlap:        1,
		TrimWhitespace: true,
	}

	start := time.Now()
	result, err := client.CreateChunkedEmbedding(ctx, longDocument, embeddingModel, config)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to create chunked embedding", err)
		return false
	}

	fmt.Printf("   ✅ Created chunked embedding (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Model: %s\n", result.Model)
	fmt.Printf("      Chunks: %d\n", len(result.Chunks))
	fmt.Printf("      Total tokens: %d\n", result.TotalTokensUsed)

	if len(result.Chunks) == 0 {
		printError("No chunk embeddings returned", nil)
		return false
	}

	// Verify each chunk has an embedding
	for i, chunkEmb := range result.Chunks {
		if len(chunkEmb.Embedding) == 0 {
			printError(fmt.Sprintf("Chunk %d has empty embedding", i), nil)
			return false
		}
		if verbose {
			preview := chunkEmb.Chunk.Text
			if len(preview) > 50 {
				preview = preview[:50] + "..."
			}
			fmt.Printf("      Chunk %d: %d dims, %q\n", i, len(chunkEmb.Embedding), preview)
		}
	}

	// Test 2: Compute document-level embedding by averaging
	fmt.Printf("   Computing average document embedding...\n")

	var embeddings [][]float64
	for _, chunkEmb := range result.Chunks {
		embeddings = append(embeddings, chunkEmb.Embedding)
	}

	docEmbedding := openrouter.AverageEmbeddings(embeddings)
	if docEmbedding == nil {
		printError("Failed to compute average embedding", nil)
		return false
	}

	fmt.Printf("   ✅ Document embedding dimensions: %d\n", len(docEmbedding))

	// Test 3: Compare similarity between chunks
	if len(result.Chunks) >= 2 {
		fmt.Printf("   Testing similarity between chunks...\n")

		sim := openrouter.CosineSimilarity(result.Chunks[0].Embedding, result.Chunks[1].Embedding)
		fmt.Printf("   ✅ Similarity between chunk 0 and 1: %.4f\n", sim)

		// Similar document chunks should have positive similarity
		if sim < -0.5 {
			fmt.Printf("   ⚠️  Unexpectedly low similarity between adjacent chunks\n")
		}
	}

	// Test 4: Test with section-based chunking
	fmt.Printf("   Testing chunked embedding with section strategy...\n")

	config2 := openrouter.ChunkConfig{
		Strategy:       openrouter.ChunkBySections,
		ChunkSize:      150,
		TrimWhitespace: true,
	}

	result2, err := client.CreateChunkedEmbedding(ctx, longDocument, embeddingModel, config2)
	if err != nil {
		printError("Failed to create section-based chunked embedding", err)
		return false
	}

	fmt.Printf("   ✅ Section-based chunking: %d chunks\n", len(result2.Chunks))

	// Test 5: Weighted average (weight by chunk length)
	fmt.Printf("   Testing weighted average embedding...\n")

	var weights []float64
	for _, chunkEmb := range result.Chunks {
		// Weight by text length as a simple example
		weights = append(weights, float64(len(chunkEmb.Chunk.Text)))
	}

	weightedEmb := openrouter.WeightedAverageEmbeddings(embeddings, weights)
	if weightedEmb == nil {
		printError("Failed to compute weighted average embedding", nil)
		return false
	}

	fmt.Printf("   ✅ Weighted embedding dimensions: %d\n", len(weightedEmb))

	// Compare weighted vs unweighted
	simWeighted := openrouter.CosineSimilarity(docEmbedding, weightedEmb)
	fmt.Printf("   ✅ Similarity (weighted vs unweighted): %.4f\n", simWeighted)

	printSuccess("Chunked embedding tests passed")
	return true
}
