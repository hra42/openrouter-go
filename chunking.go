package openrouter

import (
	"context"
	"math"
	"regexp"
	"strings"
	"unicode"
)

// ChunkStrategy defines the strategy for chunking text.
// Strategies are ordered by semantic quality, with semantic strategies (sections,
// paragraphs, sentences) preferred over size-based strategies (tokens, characters).
type ChunkStrategy string

const (
	// Semantic strategies (recommended) - preserve meaning at natural boundaries

	// ChunkBySections splits at markdown headers or document structure.
	// Best for structured documents with clear sections.
	ChunkBySections ChunkStrategy = "sections"

	// ChunkByParagraphs splits at paragraph boundaries (double newlines).
	// Good for prose and articles.
	ChunkByParagraphs ChunkStrategy = "paragraphs"

	// ChunkBySentences splits at sentence boundaries.
	// Fine-grained semantic chunking.
	ChunkBySentences ChunkStrategy = "sentences"

	// Size-based strategies (fallback) - when finer control is needed

	// ChunkByTokens splits by estimated token count.
	ChunkByTokens ChunkStrategy = "tokens"

	// ChunkByCharacters splits by character count.
	ChunkByCharacters ChunkStrategy = "characters"
)

// ChunkConfig configures text chunking behavior.
type ChunkConfig struct {
	// Strategy determines how text is split.
	Strategy ChunkStrategy

	// ChunkSize is the target size per chunk.
	// For semantic strategies: target token count (units are merged to reach this).
	// For size-based strategies: exact limit (tokens or characters).
	ChunkSize int

	// Overlap specifies units to overlap between chunks for context preservation.
	// For semantic strategies: number of units (sentences/paragraphs) to repeat.
	// For size-based strategies: number of tokens/characters to overlap.
	Overlap int

	// TrimWhitespace removes leading/trailing whitespace from chunks.
	TrimWhitespace bool

	// PreserveWords prevents splitting mid-word (for tokens/characters strategies).
	PreserveWords bool

	// SectionHeaders defines custom section delimiters for sections strategy.
	// Default: ["#", "##", "###", "####", "#####", "######"] for markdown.
	SectionHeaders []string
}

// TextChunk represents a single chunk of text with metadata.
type TextChunk struct {
	// Text is the chunk content.
	Text string

	// Index is the zero-based index of this chunk.
	Index int

	// StartPos is the byte offset in the original text where this chunk starts.
	StartPos int

	// EndPos is the byte offset in the original text where this chunk ends.
	EndPos int
}

// ChunkedEmbeddingResult contains embeddings for all chunks of a document.
type ChunkedEmbeddingResult struct {
	// Chunks contains the embedding for each chunk.
	Chunks []ChunkEmbedding

	// TotalTokensUsed is the sum of tokens used across all chunk embeddings.
	TotalTokensUsed int

	// Model is the embedding model used.
	Model string
}

// ChunkEmbedding contains the embedding for a single chunk.
type ChunkEmbedding struct {
	// Chunk is the original text chunk with metadata.
	Chunk TextChunk

	// Embedding is the vector representation.
	Embedding []float64
}

// DefaultChunkConfig returns a sensible default configuration.
func DefaultChunkConfig() ChunkConfig {
	return ChunkConfig{
		Strategy:       ChunkByParagraphs,
		ChunkSize:      512,
		Overlap:        1,
		TrimWhitespace: true,
		PreserveWords:  true,
	}
}

// ChunkText splits text into chunks according to the configuration.
func ChunkText(text string, config ChunkConfig) ([]TextChunk, error) {
	if text == "" {
		return []TextChunk{}, nil
	}

	// Apply defaults
	if config.ChunkSize <= 0 {
		config.ChunkSize = 512
	}

	var chunks []TextChunk

	switch config.Strategy {
	case ChunkBySections:
		chunks = chunkBySections(text, config)
	case ChunkByParagraphs:
		chunks = chunkByParagraphs(text, config)
	case ChunkBySentences:
		chunks = chunkBySentences(text, config)
	case ChunkByTokens:
		chunks = chunkByTokens(text, config)
	case ChunkByCharacters:
		chunks = chunkByCharacters(text, config)
	default:
		// Default to paragraphs
		chunks = chunkByParagraphs(text, config)
	}

	// Apply trimming if requested
	if config.TrimWhitespace {
		for i := range chunks {
			chunks[i].Text = strings.TrimSpace(chunks[i].Text)
		}
	}

	// Filter out empty chunks
	var filtered []TextChunk
	for _, chunk := range chunks {
		if chunk.Text != "" {
			filtered = append(filtered, chunk)
		}
	}

	// Re-index after filtering
	for i := range filtered {
		filtered[i].Index = i
	}

	return filtered, nil
}

// ChunkTexts splits multiple texts into chunks.
func ChunkTexts(texts []string, config ChunkConfig) ([][]TextChunk, error) {
	results := make([][]TextChunk, len(texts))
	for i, text := range texts {
		chunks, err := ChunkText(text, config)
		if err != nil {
			return nil, err
		}
		results[i] = chunks
	}
	return results, nil
}

// EstimateTokens estimates the token count for text using character-based heuristics.
// Uses approximately 4 characters per token for English text.
func EstimateTokens(text string) int {
	if text == "" {
		return 0
	}
	// Approximately 4 characters per token for English
	return (len(text) + 3) / 4
}

// EstimateTokensFromWords estimates token count based on word count.
// Uses approximately 1.3 tokens per word for English text.
func EstimateTokensFromWords(text string) int {
	if text == "" {
		return 0
	}
	words := len(strings.Fields(text))
	return int(math.Ceil(float64(words) * 1.3))
}

// CreateChunkedEmbedding creates embeddings for a long text by chunking it first.
func (c *Client) CreateChunkedEmbedding(ctx context.Context, text string, model string, config ChunkConfig, opts ...EmbeddingOption) (*ChunkedEmbeddingResult, error) {
	chunks, err := ChunkText(text, config)
	if err != nil {
		return nil, err
	}

	if len(chunks) == 0 {
		return &ChunkedEmbeddingResult{
			Chunks: []ChunkEmbedding{},
			Model:  model,
		}, nil
	}

	// Extract text from chunks for batch embedding
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Text
	}

	// Create embeddings for all chunks
	resp, err := c.CreateEmbeddings(ctx, texts, model, opts...)
	if err != nil {
		return nil, err
	}

	// Build result
	result := &ChunkedEmbeddingResult{
		Chunks: make([]ChunkEmbedding, len(chunks)),
		Model:  resp.Model,
	}

	for i, chunk := range chunks {
		if i < len(resp.Data) {
			result.Chunks[i] = ChunkEmbedding{
				Chunk:     chunk,
				Embedding: resp.Data[i].GetEmbeddingVector(),
			}
		}
	}

	if resp.Usage != nil {
		result.TotalTokensUsed = resp.Usage.TotalTokens
	}

	return result, nil
}

// CreateChunkedEmbeddings creates embeddings for multiple long texts by chunking each.
func (c *Client) CreateChunkedEmbeddings(ctx context.Context, texts []string, model string, config ChunkConfig, opts ...EmbeddingOption) ([]*ChunkedEmbeddingResult, error) {
	results := make([]*ChunkedEmbeddingResult, len(texts))

	for i, text := range texts {
		result, err := c.CreateChunkedEmbedding(ctx, text, model, config, opts...)
		if err != nil {
			return nil, err
		}
		results[i] = result
	}

	return results, nil
}

// AverageEmbeddings computes the element-wise average of multiple embeddings.
// Useful for combining chunk embeddings into a single document embedding.
func AverageEmbeddings(embeddings [][]float64) []float64 {
	if len(embeddings) == 0 {
		return nil
	}

	dim := len(embeddings[0])
	if dim == 0 {
		return nil
	}

	result := make([]float64, dim)

	for _, emb := range embeddings {
		if len(emb) != dim {
			// Skip embeddings with mismatched dimensions
			continue
		}
		for i, v := range emb {
			result[i] += v
		}
	}

	count := float64(len(embeddings))
	for i := range result {
		result[i] /= count
	}

	return result
}

// WeightedAverageEmbeddings computes a weighted average of embeddings.
// weights should have the same length as embeddings.
func WeightedAverageEmbeddings(embeddings [][]float64, weights []float64) []float64 {
	if len(embeddings) == 0 || len(weights) == 0 {
		return nil
	}

	if len(embeddings) != len(weights) {
		return nil
	}

	dim := len(embeddings[0])
	if dim == 0 {
		return nil
	}

	result := make([]float64, dim)
	totalWeight := 0.0

	for i, emb := range embeddings {
		if len(emb) != dim {
			continue
		}
		weight := weights[i]
		totalWeight += weight
		for j, v := range emb {
			result[j] += v * weight
		}
	}

	if totalWeight == 0 {
		return nil
	}

	for i := range result {
		result[i] /= totalWeight
	}

	return result
}

// CosineSimilarity computes the cosine similarity between two embedding vectors.
// Returns a value between -1 and 1, where 1 means identical direction.
func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
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

// --- Internal chunking implementations ---

// chunkBySections splits text at markdown headers or custom section delimiters.
func chunkBySections(text string, config ChunkConfig) []TextChunk {
	headers := config.SectionHeaders
	if len(headers) == 0 {
		// Default markdown headers
		headers = []string{"######", "#####", "####", "###", "##", "#"}
	}

	// Build regex pattern for headers
	patterns := make([]string, len(headers))
	for i, h := range headers {
		patterns[i] = regexp.QuoteMeta(h)
	}
	pattern := `(?m)^(` + strings.Join(patterns, "|") + `)[ \t]+`
	headerRegex := regexp.MustCompile(pattern)

	// Find all header positions
	matches := headerRegex.FindAllStringIndex(text, -1)

	if len(matches) == 0 {
		// No headers found, fall back to paragraphs
		return chunkByParagraphs(text, config)
	}

	// Split into sections
	var sections []string
	var positions []int

	for i, match := range matches {
		start := match[0]
		var end int
		if i+1 < len(matches) {
			end = matches[i+1][0]
		} else {
			end = len(text)
		}

		section := text[start:end]
		sections = append(sections, section)
		positions = append(positions, start)
	}

	// Handle content before first header
	if len(matches) > 0 && matches[0][0] > 0 {
		preContent := strings.TrimSpace(text[:matches[0][0]])
		if preContent != "" {
			sections = append([]string{preContent}, sections...)
			positions = append([]int{0}, positions...)
		}
	}

	// Merge sections to reach target size
	return mergeUnitsToSize(sections, positions, config)
}

// chunkByParagraphs splits text at double newlines.
func chunkByParagraphs(text string, config ChunkConfig) []TextChunk {
	// Normalize line endings
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// Split on double newlines
	paragraphSplitter := regexp.MustCompile(`\n\s*\n`)
	parts := paragraphSplitter.Split(text, -1)

	// Track positions
	var paragraphs []string
	var positions []int
	pos := 0

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			// Find actual position in original text
			idx := strings.Index(text[pos:], part)
			if idx >= 0 {
				positions = append(positions, pos+idx)
			} else {
				positions = append(positions, pos)
			}
			paragraphs = append(paragraphs, trimmed)
			pos = positions[len(positions)-1] + len(part)
		}
	}

	return mergeUnitsToSize(paragraphs, positions, config)
}

// chunkBySentences splits text at sentence boundaries.
func chunkBySentences(text string, config ChunkConfig) []TextChunk {
	sentences := splitSentences(text)

	// Track positions
	var positions []int
	pos := 0

	for _, sentence := range sentences {
		idx := strings.Index(text[pos:], sentence)
		if idx >= 0 {
			positions = append(positions, pos+idx)
			pos = positions[len(positions)-1] + len(sentence)
		} else {
			positions = append(positions, pos)
		}
	}

	return mergeUnitsToSize(sentences, positions, config)
}

// splitSentences splits text into sentences, handling common abbreviations.
func splitSentences(text string) []string {
	if text == "" {
		return nil
	}

	// Common abbreviations that shouldn't end sentences
	abbreviations := []string{
		"Mr.", "Mrs.", "Ms.", "Dr.", "Prof.", "Sr.", "Jr.",
		"vs.", "etc.", "e.g.", "i.e.", "a.m.", "p.m.",
		"Inc.", "Ltd.", "Corp.", "Co.",
		"Jan.", "Feb.", "Mar.", "Apr.", "Jun.", "Jul.", "Aug.", "Sep.", "Oct.", "Nov.", "Dec.",
		"St.", "Ave.", "Blvd.", "Rd.",
	}

	// Placeholder for abbreviations
	const placeholder = "\x00ABBR\x00"
	processed := text

	// Temporarily replace abbreviations
	for _, abbr := range abbreviations {
		processed = strings.ReplaceAll(processed, abbr, strings.ReplaceAll(abbr, ".", placeholder))
	}

	// Also protect numbers with decimals (e.g., 3.14, $1.00)
	numberRegex := regexp.MustCompile(`(\d)\.(\d)`)
	processed = numberRegex.ReplaceAllString(processed, "$1"+placeholder+"$2")

	// Protect ellipsis
	processed = strings.ReplaceAll(processed, "...", placeholder+placeholder+placeholder)

	// Split on sentence-ending punctuation
	sentenceRegex := regexp.MustCompile(`([.!?]+)(\s+|$)`)
	parts := sentenceRegex.Split(processed, -1)
	delimiters := sentenceRegex.FindAllString(processed, -1)

	var sentences []string
	for i, part := range parts {
		if part == "" {
			continue
		}
		// Restore placeholders
		sentence := strings.ReplaceAll(part, placeholder, ".")

		// Add delimiter back
		if i < len(delimiters) {
			delimiter := strings.ReplaceAll(delimiters[i], placeholder, ".")
			sentence += strings.TrimRight(delimiter, " \t\n")
		}

		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			sentences = append(sentences, sentence)
		}
	}

	return sentences
}

// mergeUnitsToSize merges semantic units (sentences, paragraphs) to reach target token size.
func mergeUnitsToSize(units []string, positions []int, config ChunkConfig) []TextChunk {
	if len(units) == 0 {
		return nil
	}

	var chunks []TextChunk
	var currentTexts []string
	var currentStart int
	currentTokens := 0

	for i, unit := range units {
		unitTokens := EstimateTokens(unit)

		// Check if adding this unit would exceed the target
		if len(currentTexts) > 0 && currentTokens+unitTokens > config.ChunkSize {
			// Create chunk from accumulated units
			chunkText := strings.Join(currentTexts, "\n\n")
			chunks = append(chunks, TextChunk{
				Text:     chunkText,
				Index:    len(chunks),
				StartPos: currentStart,
				EndPos:   currentStart + len(chunkText),
			})

			// Start new chunk, potentially with overlap
			if config.Overlap > 0 && len(currentTexts) > 0 {
				overlapStart := max(len(currentTexts)-config.Overlap, 0)
				currentTexts = currentTexts[overlapStart:]
				currentTokens = 0
				for _, t := range currentTexts {
					currentTokens += EstimateTokens(t)
				}
				// Update start position for overlapped content
				currentStart = positions[i-len(currentTexts)]
			} else {
				currentTexts = nil
				currentTokens = 0
				if i < len(positions) {
					currentStart = positions[i]
				}
			}
		}

		if len(currentTexts) == 0 && i < len(positions) {
			currentStart = positions[i]
		}

		currentTexts = append(currentTexts, unit)
		currentTokens += unitTokens
	}

	// Add remaining content
	if len(currentTexts) > 0 {
		chunkText := strings.Join(currentTexts, "\n\n")
		chunks = append(chunks, TextChunk{
			Text:     chunkText,
			Index:    len(chunks),
			StartPos: currentStart,
			EndPos:   currentStart + len(chunkText),
		})
	}

	return chunks
}

// chunkByTokens splits text by estimated token count.
func chunkByTokens(text string, config ChunkConfig) []TextChunk {
	if config.PreserveWords {
		return chunkByTokensPreserveWords(text, config)
	}

	// Simple character-based splitting (4 chars per token)
	charsPerChunk := config.ChunkSize * 4
	overlapChars := config.Overlap * 4

	var chunks []TextChunk
	pos := 0

	for pos < len(text) {
		end := min(pos+charsPerChunk, len(text))

		chunk := text[pos:end]
		chunks = append(chunks, TextChunk{
			Text:     chunk,
			Index:    len(chunks),
			StartPos: pos,
			EndPos:   end,
		})

		pos = end - overlapChars
		if pos < 0 || overlapChars == 0 {
			pos = end
		}
	}

	return chunks
}

// chunkByTokensPreserveWords splits by tokens while preserving word boundaries.
func chunkByTokensPreserveWords(text string, config ChunkConfig) []TextChunk {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	// Estimate tokens per word (~1.3)
	wordsPerChunk := max(int(float64(config.ChunkSize)/1.3), 1)

	overlapWords := int(float64(config.Overlap) / 1.3)

	var chunks []TextChunk
	pos := 0
	textPos := 0

	for pos < len(words) {
		end := min(pos+wordsPerChunk, len(words))

		chunkWords := words[pos:end]
		chunkText := strings.Join(chunkWords, " ")

		// Find actual position in text (safely)
		startPos := textPos
		if textPos < len(text) {
			idx := strings.Index(text[textPos:], chunkWords[0])
			if idx >= 0 {
				startPos = textPos + idx
			}
		}

		chunks = append(chunks, TextChunk{
			Text:     chunkText,
			Index:    len(chunks),
			StartPos: startPos,
			EndPos:   startPos + len(chunkText),
		})

		// Update textPos safely
		textPos = min(startPos+len(chunkText), len(text))

		// Calculate next position with overlap
		nextPos := end - overlapWords
		if nextPos <= pos || overlapWords == 0 {
			nextPos = end
		}
		pos = nextPos
	}

	return chunks
}

// chunkByCharacters splits text by character count.
func chunkByCharacters(text string, config ChunkConfig) []TextChunk {
	if config.PreserveWords {
		return chunkByCharactersPreserveWords(text, config)
	}

	chunkSize := config.ChunkSize
	overlap := config.Overlap

	var chunks []TextChunk
	pos := 0

	for pos < len(text) {
		end := min(pos+chunkSize, len(text))

		chunks = append(chunks, TextChunk{
			Text:     text[pos:end],
			Index:    len(chunks),
			StartPos: pos,
			EndPos:   end,
		})

		pos = end - overlap
		if pos < 0 || overlap == 0 {
			pos = end
		}
	}

	return chunks
}

// chunkByCharactersPreserveWords splits by characters while preserving word boundaries.
func chunkByCharactersPreserveWords(text string, config ChunkConfig) []TextChunk {
	chunkSize := config.ChunkSize
	overlap := config.Overlap

	var chunks []TextChunk
	pos := 0

	for pos < len(text) {
		end := min(pos+chunkSize, len(text))

		// Find last word boundary before end
		if end < len(text) {
			for end > pos && !unicode.IsSpace(rune(text[end])) {
				end--
			}
			// If no space found, just use the original end
			if end == pos {
				end = min(pos+chunkSize, len(text))
			}
		}

		chunkText := strings.TrimSpace(text[pos:end])
		if chunkText != "" {
			chunks = append(chunks, TextChunk{
				Text:     chunkText,
				Index:    len(chunks),
				StartPos: pos,
				EndPos:   end,
			})
		}

		// Calculate overlap position
		overlapPos := end - overlap
		if overlapPos <= pos || overlap == 0 {
			pos = end
		} else {
			// Find word boundary for overlap
			for overlapPos < end && !unicode.IsSpace(rune(text[overlapPos])) {
				overlapPos++
			}
			pos = overlapPos
		}

		// Skip leading whitespace
		for pos < len(text) && unicode.IsSpace(rune(text[pos])) {
			pos++
		}
	}

	return chunks
}
