package openrouter

import (
	"strconv"
)

// CostResult holds a cost breakdown for a request.
type CostResult struct {
	// PromptCost is the cost for prompt/input tokens.
	PromptCost float64
	// CompletionCost is the cost for completion/output tokens.
	CompletionCost float64
	// TotalCost is the sum of PromptCost and CompletionCost.
	TotalCost float64
}

// ParsePricing parses a pricing string (e.g., "0.0025") to float64.
// Pricing values represent cost per million tokens.
// Returns 0 on error, empty string, or negative values.
func ParsePricing(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}
	val, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || val < 0 {
		return 0
	}
	return val
}

// ParsePricingPtr parses a *string pricing value.
// Returns 0 if the pointer is nil.
func ParsePricingPtr(priceStr *string) float64 {
	if priceStr == nil {
		return 0
	}
	return ParsePricing(*priceStr)
}

// EstimateCost calculates the estimated cost from ModelPricing and Usage.
// Pricing values are per-million tokens: cost = tokens * price / 1,000,000.
func EstimateCost(pricing ModelPricing, usage Usage) CostResult {
	return EstimateCostFromTokens(
		pricing.Prompt,
		pricing.Completion,
		usage.PromptTokens,
		usage.CompletionTokens,
	)
}

// EstimateCostFromTokens calculates cost from raw pricing strings and token counts.
// This is useful when you don't have the full ModelPricing and Usage structs.
func EstimateCostFromTokens(promptPrice, completionPrice string, promptTokens, completionTokens int) CostResult {
	pp := ParsePricing(promptPrice)
	cp := ParsePricing(completionPrice)

	promptCost := float64(promptTokens) * pp / 1_000_000
	completionCost := float64(completionTokens) * cp / 1_000_000

	return CostResult{
		PromptCost:     promptCost,
		CompletionCost: completionCost,
		TotalCost:      promptCost + completionCost,
	}
}

// EstimateCostWithCaching calculates cost including cache read/write tokens.
// Cache read tokens are typically cheaper than regular prompt tokens,
// and cache write tokens may have a different price as well.
func EstimateCostWithCaching(pricing ModelPricing, usage Usage, cacheReadTokens, cacheWriteTokens int) CostResult {
	// Base cost from non-cached tokens
	baseCost := EstimateCost(pricing, usage)

	// Cache costs
	cacheReadPrice := ParsePricingPtr(pricing.InputCacheRead)
	cacheWritePrice := ParsePricingPtr(pricing.InputCacheWrite)

	cacheReadCost := float64(cacheReadTokens) * cacheReadPrice / 1_000_000
	cacheWriteCost := float64(cacheWriteTokens) * cacheWritePrice / 1_000_000

	return CostResult{
		PromptCost:     baseCost.PromptCost + cacheReadCost + cacheWriteCost,
		CompletionCost: baseCost.CompletionCost,
		TotalCost:      baseCost.TotalCost + cacheReadCost + cacheWriteCost,
	}
}
