package openrouter

import (
	"math"
	"testing"
)

func TestParsePricing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"valid price", "0.0025", 0.0025},
		{"zero", "0", 0},
		{"whole number", "5", 5},
		{"empty string", "", 0},
		{"negative value", "-1", 0},
		{"invalid string", "not-a-number", 0},
		{"scientific notation", "2.5e-4", 0.00025},
		{"large value", "1000000", 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParsePricing(tt.input)
			if math.Abs(result-tt.expected) > 1e-12 {
				t.Errorf("ParsePricing(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParsePricingPtr(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		result := ParsePricingPtr(nil)
		if result != 0 {
			t.Errorf("ParsePricingPtr(nil) = %v, want 0", result)
		}
	})

	t.Run("valid pointer", func(t *testing.T) {
		s := "0.005"
		result := ParsePricingPtr(&s)
		if math.Abs(result-0.005) > 1e-12 {
			t.Errorf("ParsePricingPtr(&%q) = %v, want 0.005", s, result)
		}
	})

	t.Run("empty string pointer", func(t *testing.T) {
		s := ""
		result := ParsePricingPtr(&s)
		if result != 0 {
			t.Errorf("ParsePricingPtr(&%q) = %v, want 0", s, result)
		}
	})
}

func TestEstimateCost(t *testing.T) {
	tests := []struct {
		name           string
		pricing        ModelPricing
		usage          Usage
		wantPrompt     float64
		wantCompletion float64
		wantTotal      float64
	}{
		{
			name: "typical GPT-4o pricing",
			pricing: ModelPricing{
				Prompt:     "2.5",
				Completion: "10",
			},
			usage: Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
			},
			wantPrompt:     0.0025,
			wantCompletion: 0.005,
			wantTotal:      0.0075,
		},
		{
			name: "zero tokens",
			pricing: ModelPricing{
				Prompt:     "2.5",
				Completion: "10",
			},
			usage: Usage{
				PromptTokens:     0,
				CompletionTokens: 0,
			},
			wantPrompt:     0,
			wantCompletion: 0,
			wantTotal:      0,
		},
		{
			name: "free model",
			pricing: ModelPricing{
				Prompt:     "0",
				Completion: "0",
			},
			usage: Usage{
				PromptTokens:     10000,
				CompletionTokens: 5000,
			},
			wantPrompt:     0,
			wantCompletion: 0,
			wantTotal:      0,
		},
		{
			name: "empty pricing strings",
			pricing: ModelPricing{
				Prompt:     "",
				Completion: "",
			},
			usage: Usage{
				PromptTokens:     1000,
				CompletionTokens: 500,
			},
			wantPrompt:     0,
			wantCompletion: 0,
			wantTotal:      0,
		},
		{
			name: "one million tokens",
			pricing: ModelPricing{
				Prompt:     "1",
				Completion: "2",
			},
			usage: Usage{
				PromptTokens:     1000000,
				CompletionTokens: 1000000,
			},
			wantPrompt:     1,
			wantCompletion: 2,
			wantTotal:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateCost(tt.pricing, tt.usage)
			if math.Abs(result.PromptCost-tt.wantPrompt) > 1e-10 {
				t.Errorf("PromptCost = %v, want %v", result.PromptCost, tt.wantPrompt)
			}
			if math.Abs(result.CompletionCost-tt.wantCompletion) > 1e-10 {
				t.Errorf("CompletionCost = %v, want %v", result.CompletionCost, tt.wantCompletion)
			}
			if math.Abs(result.TotalCost-tt.wantTotal) > 1e-10 {
				t.Errorf("TotalCost = %v, want %v", result.TotalCost, tt.wantTotal)
			}
		})
	}
}

func TestEstimateCostFromTokens(t *testing.T) {
	result := EstimateCostFromTokens("5", "15", 2000, 1000)
	wantPrompt := 0.01
	wantCompletion := 0.015
	wantTotal := 0.025

	if math.Abs(result.PromptCost-wantPrompt) > 1e-10 {
		t.Errorf("PromptCost = %v, want %v", result.PromptCost, wantPrompt)
	}
	if math.Abs(result.CompletionCost-wantCompletion) > 1e-10 {
		t.Errorf("CompletionCost = %v, want %v", result.CompletionCost, wantCompletion)
	}
	if math.Abs(result.TotalCost-wantTotal) > 1e-10 {
		t.Errorf("TotalCost = %v, want %v", result.TotalCost, wantTotal)
	}
}

func TestEstimateCostWithCaching(t *testing.T) {
	cacheRead := "0.5"
	cacheWrite := "3"
	pricing := ModelPricing{
		Prompt:          "2.5",
		Completion:      "10",
		InputCacheRead:  &cacheRead,
		InputCacheWrite: &cacheWrite,
	}
	usage := Usage{
		PromptTokens:     1000,
		CompletionTokens: 500,
	}

	result := EstimateCostWithCaching(pricing, usage, 5000, 2000)

	// Base prompt: 1000 * 2.5 / 1M = 0.0025
	// Cache read: 5000 * 0.5 / 1M = 0.0025
	// Cache write: 2000 * 3 / 1M = 0.006
	// Total prompt: 0.0025 + 0.0025 + 0.006 = 0.011
	// Completion: 500 * 10 / 1M = 0.005
	// Total: 0.016

	wantPrompt := 0.011
	wantCompletion := 0.005
	wantTotal := 0.016

	if math.Abs(result.PromptCost-wantPrompt) > 1e-10 {
		t.Errorf("PromptCost = %v, want %v", result.PromptCost, wantPrompt)
	}
	if math.Abs(result.CompletionCost-wantCompletion) > 1e-10 {
		t.Errorf("CompletionCost = %v, want %v", result.CompletionCost, wantCompletion)
	}
	if math.Abs(result.TotalCost-wantTotal) > 1e-10 {
		t.Errorf("TotalCost = %v, want %v", result.TotalCost, wantTotal)
	}
}

func TestEstimateCostWithCachingNilPrices(t *testing.T) {
	pricing := ModelPricing{
		Prompt:          "2.5",
		Completion:      "10",
		InputCacheRead:  nil,
		InputCacheWrite: nil,
	}
	usage := Usage{
		PromptTokens:     1000,
		CompletionTokens: 500,
	}

	result := EstimateCostWithCaching(pricing, usage, 5000, 2000)

	// Cache prices are nil, so only base cost applies
	wantPrompt := 0.0025
	wantCompletion := 0.005
	wantTotal := 0.0075

	if math.Abs(result.PromptCost-wantPrompt) > 1e-10 {
		t.Errorf("PromptCost = %v, want %v", result.PromptCost, wantPrompt)
	}
	if math.Abs(result.CompletionCost-wantCompletion) > 1e-10 {
		t.Errorf("CompletionCost = %v, want %v", result.CompletionCost, wantCompletion)
	}
	if math.Abs(result.TotalCost-wantTotal) > 1e-10 {
		t.Errorf("TotalCost = %v, want %v", result.TotalCost, wantTotal)
	}
}
