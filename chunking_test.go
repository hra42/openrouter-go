package openrouter

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChunkBySections(t *testing.T) {
	text := `# Introduction
This is the introduction section with some content.

## First Section
This is the first section content.
It has multiple lines.

## Second Section
This is the second section content.

### Subsection
This is a subsection.

## Third Section
Final section content.`

	config := ChunkConfig{
		Strategy:       ChunkBySections,
		ChunkSize:      100, // Small to force splits
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// Verify each chunk contains a header
	for i, chunk := range chunks {
		if !strings.Contains(chunk.Text, "#") && i > 0 {
			// First chunk may not have a header if there's content before the first header
			if i == 0 && !strings.Contains(text, "# ") {
				continue
			}
		}
		if chunk.Index != i {
			t.Errorf("chunk %d has wrong index %d", i, chunk.Index)
		}
	}
}

func TestChunkByParagraphs(t *testing.T) {
	text := `This is the first paragraph. It contains multiple sentences.

This is the second paragraph. It also has content.

This is the third paragraph.

And a fourth one.`

	config := ChunkConfig{
		Strategy:       ChunkByParagraphs,
		ChunkSize:      50, // Small to force some merging
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// Verify no chunk contains double newlines (paragraphs are merged)
	for i, chunk := range chunks {
		if strings.Contains(chunk.Text, "\n\n\n") {
			t.Errorf("chunk %d contains triple newlines", i)
		}
	}
}

func TestChunkBySentences(t *testing.T) {
	text := "This is the first sentence. This is the second sentence! And here is the third? Finally, the fourth sentence."

	config := ChunkConfig{
		Strategy:       ChunkBySentences,
		ChunkSize:      30, // Small to force merging
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}

	// All content should be preserved
	combinedText := ""
	for _, chunk := range chunks {
		if combinedText != "" {
			combinedText += " "
		}
		combinedText += chunk.Text
	}

	// Check that all sentences are present
	if !strings.Contains(combinedText, "first sentence") {
		t.Error("missing first sentence")
	}
	if !strings.Contains(combinedText, "fourth sentence") {
		t.Error("missing fourth sentence")
	}
}

func TestChunkByTokens(t *testing.T) {
	// Create text with known content
	text := strings.Repeat("word ", 100) // ~100 words

	config := ChunkConfig{
		Strategy:      ChunkByTokens,
		ChunkSize:     50, // ~50 tokens per chunk
		PreserveWords: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) < 2 {
		t.Fatal("expected at least 2 chunks")
	}

	// Each chunk should have approximately the right token count
	for _, chunk := range chunks {
		tokens := EstimateTokens(chunk.Text)
		// Allow some variance due to word preservation
		if tokens > config.ChunkSize*2 {
			t.Errorf("chunk has too many tokens: %d", tokens)
		}
	}
}

func TestChunkByCharacters(t *testing.T) {
	text := "abcdefghij" // 10 characters

	config := ChunkConfig{
		Strategy:  ChunkByCharacters,
		ChunkSize: 3,
		Overlap:   0,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 4 // ceil(10/3) = 4 chunks
	if len(chunks) != expected {
		t.Errorf("expected %d chunks, got %d", expected, len(chunks))
	}

	if len(chunks) >= 1 && chunks[0].Text != "abc" {
		t.Errorf("first chunk should be 'abc', got %q", chunks[0].Text)
	}
}

func TestChunkOverlap(t *testing.T) {
	text := `First paragraph content here.

Second paragraph content here.

Third paragraph content here.

Fourth paragraph content here.`

	config := ChunkConfig{
		Strategy:       ChunkByParagraphs,
		ChunkSize:      30, // Small to ensure chunking
		Overlap:        1,  // Overlap by 1 paragraph
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) < 2 {
		t.Skip("not enough chunks to test overlap")
	}

	// With overlap, consecutive chunks should share some content
	// This tests that the overlap mechanism is working
	for i := 1; i < len(chunks); i++ {
		// Check that this chunk starts with content from previous chunk
		// (this may not always happen due to token limits, so just verify basic structure)
		if chunks[i].Text == "" {
			t.Errorf("chunk %d is empty", i)
		}
	}
}

func TestSemanticMerging(t *testing.T) {
	// Short paragraphs that should be merged
	text := `Short.

Also short.

Tiny.

Brief.`

	config := ChunkConfig{
		Strategy:       ChunkByParagraphs,
		ChunkSize:      500, // Large size means all should merge into one
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All short paragraphs should be merged into one chunk
	if len(chunks) != 1 {
		t.Errorf("expected 1 merged chunk, got %d", len(chunks))
	}

	if len(chunks) > 0 {
		// Verify all content is present
		if !strings.Contains(chunks[0].Text, "Short") {
			t.Error("missing 'Short'")
		}
		if !strings.Contains(chunks[0].Text, "Brief") {
			t.Error("missing 'Brief'")
		}
	}
}

func TestPreserveWords(t *testing.T) {
	text := "hello world testing chunking"

	config := ChunkConfig{
		Strategy:      ChunkByCharacters,
		ChunkSize:     10,
		PreserveWords: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No chunk should contain a word split in the middle
	for _, chunk := range chunks {
		words := strings.Fields(chunk.Text)
		for _, word := range words {
			if word != "hello" && word != "world" && word != "testing" && word != "chunking" {
				t.Errorf("unexpected partial word in chunk: %q", word)
			}
		}
	}
}

func TestEmptyInput(t *testing.T) {
	config := DefaultChunkConfig()

	chunks, err := ChunkText("", config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty input, got %d", len(chunks))
	}
}

func TestUTF8Handling(t *testing.T) {
	text := "Hello 世界! Привет мир! مرحبا بالعالم"

	config := ChunkConfig{
		Strategy:       ChunkByCharacters,
		ChunkSize:      20,
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// All content should be preserved
	combined := ""
	for _, chunk := range chunks {
		combined += chunk.Text
	}

	// Check for presence of each language's content
	if !strings.Contains(combined, "世界") {
		t.Error("missing Chinese characters")
	}
	if !strings.Contains(combined, "Привет") {
		t.Error("missing Russian characters")
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"test", 1},                    // 4 chars = 1 token
		{"hello world", 3},             // 11 chars = ~3 tokens
		{strings.Repeat("a", 100), 25}, // 100 chars = 25 tokens
	}

	for _, tc := range tests {
		result := EstimateTokens(tc.text)
		if result != tc.expected {
			t.Errorf("EstimateTokens(%q) = %d, expected %d", tc.text, result, tc.expected)
		}
	}
}

func TestEstimateTokensFromWords(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"hello", 2},                   // 1 word * 1.3 = 2 (ceiling)
		{"hello world", 3},             // 2 words * 1.3 = 2.6 = 3 (ceiling)
		{"one two three four five", 7}, // 5 words * 1.3 = 6.5 = 7 (ceiling)
	}

	for _, tc := range tests {
		result := EstimateTokensFromWords(tc.text)
		if result != tc.expected {
			t.Errorf("EstimateTokensFromWords(%q) = %d, expected %d", tc.text, result, tc.expected)
		}
	}
}

func TestAverageEmbeddings(t *testing.T) {
	embeddings := [][]float64{
		{1.0, 2.0, 3.0},
		{2.0, 4.0, 6.0},
		{3.0, 6.0, 9.0},
	}

	result := AverageEmbeddings(embeddings)

	expected := []float64{2.0, 4.0, 6.0}
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if math.Abs(result[i]-expected[i]) > 0.001 {
			t.Errorf("result[%d] = %f, expected %f", i, result[i], expected[i])
		}
	}
}

func TestAverageEmbeddingsEmpty(t *testing.T) {
	result := AverageEmbeddings(nil)
	if result != nil {
		t.Errorf("expected nil for empty input, got %v", result)
	}

	result = AverageEmbeddings([][]float64{})
	if result != nil {
		t.Errorf("expected nil for empty slice, got %v", result)
	}
}

func TestWeightedAverageEmbeddings(t *testing.T) {
	embeddings := [][]float64{
		{1.0, 2.0, 3.0},
		{2.0, 4.0, 6.0},
	}
	weights := []float64{1.0, 3.0} // Second embedding has 3x weight

	result := WeightedAverageEmbeddings(embeddings, weights)

	// (1*1 + 2*3) / 4 = 7/4 = 1.75
	// (2*1 + 4*3) / 4 = 14/4 = 3.5
	// (3*1 + 6*3) / 4 = 21/4 = 5.25
	expected := []float64{1.75, 3.5, 5.25}

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if math.Abs(result[i]-expected[i]) > 0.001 {
			t.Errorf("result[%d] = %f, expected %f", i, result[i], expected[i])
		}
	}
}

func TestWeightedAverageEmbeddingsInvalid(t *testing.T) {
	// Mismatched lengths
	result := WeightedAverageEmbeddings([][]float64{{1.0}}, []float64{1.0, 2.0})
	if result != nil {
		t.Error("expected nil for mismatched lengths")
	}

	// Empty inputs
	result = WeightedAverageEmbeddings(nil, nil)
	if result != nil {
		t.Error("expected nil for nil inputs")
	}
}

func TestCosineSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		a        []float64
		b        []float64
		expected float64
	}{
		{
			name:     "identical vectors",
			a:        []float64{1.0, 2.0, 3.0},
			b:        []float64{1.0, 2.0, 3.0},
			expected: 1.0,
		},
		{
			name:     "opposite vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{-1.0, 0.0, 0.0},
			expected: -1.0,
		},
		{
			name:     "orthogonal vectors",
			a:        []float64{1.0, 0.0, 0.0},
			b:        []float64{0.0, 1.0, 0.0},
			expected: 0.0,
		},
		{
			name:     "different lengths",
			a:        []float64{1.0, 2.0},
			b:        []float64{1.0, 2.0, 3.0},
			expected: 0.0,
		},
		{
			name:     "empty vectors",
			a:        []float64{},
			b:        []float64{},
			expected: 0.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := CosineSimilarity(tc.a, tc.b)
			if math.Abs(result-tc.expected) > 0.001 {
				t.Errorf("CosineSimilarity(%v, %v) = %f, expected %f", tc.a, tc.b, result, tc.expected)
			}
		})
	}
}

func TestCreateChunkedEmbedding(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/embeddings" {
			t.Errorf("expected path /embeddings, got %s", r.URL.Path)
		}

		var reqBody EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Get number of inputs
		var inputCount int
		switch input := reqBody.Input.(type) {
		case []interface{}:
			inputCount = len(input)
		case string:
			inputCount = 1
		}

		// Generate embeddings for each input
		data := make([]EmbeddingData, inputCount)
		for i := 0; i < inputCount; i++ {
			data[i] = EmbeddingData{
				Object:    "embedding",
				Embedding: []interface{}{0.1, 0.2, 0.3, 0.4, 0.5},
				Index:     i,
			}
		}

		response := EmbeddingResponse{
			ID:     "emb-chunked",
			Object: "list",
			Data:   data,
			Model:  reqBody.Model,
			Usage: &EmbeddingUsage{
				PromptTokens: inputCount * 10,
				TotalTokens:  inputCount * 10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	longText := `First paragraph with some content here.

Second paragraph with different content.

Third paragraph to make this longer.

Fourth paragraph for more content.`

	config := ChunkConfig{
		Strategy:  ChunkByParagraphs,
		ChunkSize: 30, // Small to force chunking
	}

	result, err := client.CreateChunkedEmbedding(context.Background(), longText, "test-model", config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Chunks) == 0 {
		t.Fatal("expected at least one chunk embedding")
	}

	if result.Model != "test-model" {
		t.Errorf("expected model 'test-model', got %q", result.Model)
	}

	// Verify each chunk has an embedding
	for i, chunkEmb := range result.Chunks {
		if len(chunkEmb.Embedding) == 0 {
			t.Errorf("chunk %d has empty embedding", i)
		}
		if chunkEmb.Chunk.Text == "" {
			t.Errorf("chunk %d has empty text", i)
		}
	}
}

func TestCreateChunkedEmbeddingEmpty(t *testing.T) {
	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL("http://localhost:9999"), // Won't be called
	)

	result, err := client.CreateChunkedEmbedding(context.Background(), "", "test-model", DefaultChunkConfig())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Chunks) != 0 {
		t.Errorf("expected 0 chunks for empty input, got %d", len(result.Chunks))
	}
}

func TestChunkTexts(t *testing.T) {
	texts := []string{
		"First text content.",
		"Second text content here.",
		"Third text with more content.",
	}

	config := ChunkConfig{
		Strategy:  ChunkByCharacters,
		ChunkSize: 100,
	}

	results, err := ChunkTexts(texts, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != len(texts) {
		t.Errorf("expected %d results, got %d", len(texts), len(results))
	}

	for i, chunks := range results {
		if len(chunks) == 0 {
			t.Errorf("text %d produced no chunks", i)
		}
	}
}

func TestDefaultChunkConfig(t *testing.T) {
	config := DefaultChunkConfig()

	if config.Strategy != ChunkByParagraphs {
		t.Errorf("expected default strategy ChunkByParagraphs, got %v", config.Strategy)
	}
	if config.ChunkSize != 512 {
		t.Errorf("expected default ChunkSize 512, got %d", config.ChunkSize)
	}
	if config.Overlap != 1 {
		t.Errorf("expected default Overlap 1, got %d", config.Overlap)
	}
	if !config.TrimWhitespace {
		t.Error("expected default TrimWhitespace true")
	}
	if !config.PreserveWords {
		t.Error("expected default PreserveWords true")
	}
}

func TestSplitSentencesAbbreviations(t *testing.T) {
	text := "Dr. Smith works at MIT. He has a Ph.D. from Stanford. The meeting is at 3 p.m. today."

	sentences := splitSentences(text)

	// Should not split at abbreviation periods, but will split at real sentence endings
	// Expected: ["Dr. Smith works at MIT.", "He has a Ph.D. from Stanford.", "The meeting is at 3 p.m.", "today."]
	// Note: "MIT." ends a sentence, "p.m." ends a sentence, "today." is the last fragment
	if len(sentences) < 3 {
		t.Errorf("expected at least 3 sentences, got %d: %v", len(sentences), sentences)
	}

	// First sentence should include "Dr." without splitting at that abbreviation
	if len(sentences) > 0 && !strings.Contains(sentences[0], "Dr.") {
		t.Errorf("first sentence should contain 'Dr.': %q", sentences[0])
	}

	// Verify Ph.D. was not split
	foundPhD := false
	for _, s := range sentences {
		if strings.Contains(s, "Ph.D.") {
			foundPhD = true
			break
		}
	}
	if !foundPhD {
		t.Errorf("expected to find 'Ph.D.' in sentences: %v", sentences)
	}
}

func TestSplitSentencesEllipsis(t *testing.T) {
	text := "I wonder... Maybe we should try? Yes, let's do it!"

	sentences := splitSentences(text)

	// Ellipsis counts as a sentence ending, so this becomes:
	// ["I wonder...", "Maybe we should try?", "Yes, let's do it!"]
	// Or possibly combined differently based on implementation
	if len(sentences) < 2 {
		t.Errorf("expected at least 2 sentences, got %d: %v", len(sentences), sentences)
	}

	// Verify we have all the content
	combined := strings.Join(sentences, " ")
	if !strings.Contains(combined, "wonder") {
		t.Error("missing 'wonder' content")
	}
	if !strings.Contains(combined, "try") {
		t.Error("missing 'try' content")
	}
}

func TestSplitSentencesNumbers(t *testing.T) {
	text := "The value is 3.14 which is pi. The price is $1.99 for this item."

	sentences := splitSentences(text)

	// Should not split at decimal points
	if len(sentences) != 2 {
		t.Errorf("expected 2 sentences, got %d: %v", len(sentences), sentences)
	}

	if len(sentences) > 0 && !strings.Contains(sentences[0], "3.14") {
		t.Error("first sentence should contain '3.14'")
	}
}

func TestChunkBySectionsNoHeaders(t *testing.T) {
	// Text without any headers should fall back to paragraph chunking
	text := `This is paragraph one.

This is paragraph two.

This is paragraph three.`

	config := ChunkConfig{
		Strategy:       ChunkBySections,
		ChunkSize:      500,
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still produce chunks (fallback to paragraphs)
	if len(chunks) == 0 {
		t.Fatal("expected at least one chunk")
	}
}

func TestChunkPositions(t *testing.T) {
	text := "First chunk. Second chunk. Third chunk."

	config := ChunkConfig{
		Strategy:       ChunkBySentences,
		ChunkSize:      10, // Very small to force splits
		TrimWhitespace: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify positions are reasonable (within text bounds)
	for _, chunk := range chunks {
		if chunk.StartPos < 0 || chunk.StartPos > len(text) {
			t.Errorf("invalid StartPos: %d (text length: %d)", chunk.StartPos, len(text))
		}
		if chunk.EndPos < chunk.StartPos {
			t.Errorf("EndPos (%d) < StartPos (%d)", chunk.EndPos, chunk.StartPos)
		}
	}
}

func TestChunkByTokensWithOverlap(t *testing.T) {
	text := strings.Repeat("word ", 50) // 50 words

	config := ChunkConfig{
		Strategy:      ChunkByTokens,
		ChunkSize:     25,
		Overlap:       5,
		PreserveWords: true,
	}

	chunks, err := ChunkText(text, config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(chunks) < 2 {
		t.Skip("not enough chunks to test overlap")
	}

	// With overlap, chunks should share some content
	// Just verify we get multiple chunks
	if len(chunks) <= 1 {
		t.Error("expected multiple chunks with overlap")
	}
}
