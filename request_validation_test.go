package openrouter

import (
	"testing"
)

func floatPtr(f float64) *float64 { return &f }
func intPtr(i int) *int           { return &i }

// boolPtr is declared in broadcast_test.go

func TestValidateChatCompletionParams(t *testing.T) {
	tests := []struct {
		name      string
		req       *ChatCompletionRequest
		wantErr   bool
		wantField string
	}{
		{
			name: "valid request",
			req:  &ChatCompletionRequest{},
		},
		{
			name:      "tool_choice without tools",
			req:       &ChatCompletionRequest{ToolChoice: "auto"},
			wantErr:   true,
			wantField: "tool_choice",
		},
		{
			name: "tool_choice with tools",
			req: &ChatCompletionRequest{
				ToolChoice: "auto",
				Tools:      []Tool{{Type: "function", Function: Function{Name: "test"}}},
			},
		},
		{
			name:      "parallel_tool_calls without tools",
			req:       &ChatCompletionRequest{ParallelToolCalls: boolPtr(true)},
			wantErr:   true,
			wantField: "parallel_tool_calls",
		},
		{
			name:      "top_logprobs without logprobs enabled",
			req:       &ChatCompletionRequest{TopLogProbs: intPtr(5)},
			wantErr:   true,
			wantField: "top_logprobs",
		},
		{
			name: "top_logprobs with logprobs enabled",
			req:  &ChatCompletionRequest{TopLogProbs: intPtr(5), LogProbs: boolPtr(true)},
		},
		{
			name:      "negative temperature",
			req:       &ChatCompletionRequest{Temperature: floatPtr(-0.1)},
			wantErr:   true,
			wantField: "temperature",
		},
		{
			name: "zero temperature is valid",
			req:  &ChatCompletionRequest{Temperature: floatPtr(0)},
		},
		{
			name:      "top_p above 1",
			req:       &ChatCompletionRequest{TopP: floatPtr(1.5)},
			wantErr:   true,
			wantField: "top_p",
		},
		{
			name:      "top_p below 0",
			req:       &ChatCompletionRequest{TopP: floatPtr(-0.1)},
			wantErr:   true,
			wantField: "top_p",
		},
		{
			name: "top_p at boundary",
			req:  &ChatCompletionRequest{TopP: floatPtr(1.0)},
		},
		{
			name:      "min_p above 1",
			req:       &ChatCompletionRequest{MinP: floatPtr(1.1)},
			wantErr:   true,
			wantField: "min_p",
		},
		{
			name:      "negative top_k",
			req:       &ChatCompletionRequest{TopK: intPtr(-1)},
			wantErr:   true,
			wantField: "top_k",
		},
		{
			name: "zero top_k is valid",
			req:  &ChatCompletionRequest{TopK: intPtr(0)},
		},
		{
			name:      "zero max_tokens",
			req:       &ChatCompletionRequest{MaxTokens: intPtr(0)},
			wantErr:   true,
			wantField: "max_tokens",
		},
		{
			name:      "negative max_tokens",
			req:       &ChatCompletionRequest{MaxTokens: intPtr(-5)},
			wantErr:   true,
			wantField: "max_tokens",
		},
		{
			name: "positive max_tokens",
			req:  &ChatCompletionRequest{MaxTokens: intPtr(100)},
		},
		{
			name:      "frequency_penalty too high",
			req:       &ChatCompletionRequest{FrequencyPenalty: floatPtr(2.5)},
			wantErr:   true,
			wantField: "frequency_penalty",
		},
		{
			name:      "frequency_penalty too low",
			req:       &ChatCompletionRequest{FrequencyPenalty: floatPtr(-2.5)},
			wantErr:   true,
			wantField: "frequency_penalty",
		},
		{
			name: "frequency_penalty at boundary",
			req:  &ChatCompletionRequest{FrequencyPenalty: floatPtr(2.0)},
		},
		{
			name:      "presence_penalty too high",
			req:       &ChatCompletionRequest{PresencePenalty: floatPtr(3.0)},
			wantErr:   true,
			wantField: "presence_penalty",
		},
		{
			name:      "top_logprobs too high",
			req:       &ChatCompletionRequest{TopLogProbs: intPtr(25), LogProbs: boolPtr(true)},
			wantErr:   true,
			wantField: "top_logprobs",
		},
		{
			name:      "top_logprobs negative",
			req:       &ChatCompletionRequest{TopLogProbs: intPtr(-1), LogProbs: boolPtr(true)},
			wantErr:   true,
			wantField: "top_logprobs",
		},
		{
			name:      "invalid provider sort",
			req:       &ChatCompletionRequest{Provider: &Provider{Sort: "invalid"}},
			wantErr:   true,
			wantField: "provider.sort",
		},
		{
			name: "valid provider sort - price",
			req:  &ChatCompletionRequest{Provider: &Provider{Sort: "price"}},
		},
		{
			name: "valid provider sort - throughput",
			req:  &ChatCompletionRequest{Provider: &Provider{Sort: "throughput"}},
		},
		{
			name: "valid provider sort - latency",
			req:  &ChatCompletionRequest{Provider: &Provider{Sort: "latency"}},
		},
		{
			name: "empty provider sort is valid",
			req:  &ChatCompletionRequest{Provider: &Provider{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatCompletionParams(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateChatCompletionParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				valErr, ok := IsValidationError(err)
				if !ok {
					t.Errorf("expected ValidationError, got %T", err)
					return
				}
				if valErr.Field != tt.wantField {
					t.Errorf("expected field %q, got %q", tt.wantField, valErr.Field)
				}
			}
		})
	}
}

func TestValidateCompletionParams(t *testing.T) {
	tests := []struct {
		name      string
		req       *CompletionRequest
		wantErr   bool
		wantField string
	}{
		{
			name: "valid request",
			req:  &CompletionRequest{},
		},
		{
			name:      "negative temperature",
			req:       &CompletionRequest{Temperature: floatPtr(-1)},
			wantErr:   true,
			wantField: "temperature",
		},
		{
			name:      "top_p out of range",
			req:       &CompletionRequest{TopP: floatPtr(2.0)},
			wantErr:   true,
			wantField: "top_p",
		},
		{
			name:      "min_p out of range",
			req:       &CompletionRequest{MinP: floatPtr(-0.5)},
			wantErr:   true,
			wantField: "min_p",
		},
		{
			name:      "negative top_k",
			req:       &CompletionRequest{TopK: intPtr(-1)},
			wantErr:   true,
			wantField: "top_k",
		},
		{
			name:      "zero max_tokens",
			req:       &CompletionRequest{MaxTokens: intPtr(0)},
			wantErr:   true,
			wantField: "max_tokens",
		},
		{
			name:      "frequency_penalty out of range",
			req:       &CompletionRequest{FrequencyPenalty: floatPtr(3.0)},
			wantErr:   true,
			wantField: "frequency_penalty",
		},
		{
			name:      "presence_penalty out of range",
			req:       &CompletionRequest{PresencePenalty: floatPtr(-3.0)},
			wantErr:   true,
			wantField: "presence_penalty",
		},
		{
			name:      "best_of less than n",
			req:       &CompletionRequest{BestOf: intPtr(2), N: intPtr(5)},
			wantErr:   true,
			wantField: "best_of",
		},
		{
			name: "best_of equal to n",
			req:  &CompletionRequest{BestOf: intPtr(3), N: intPtr(3)},
		},
		{
			name: "best_of greater than n",
			req:  &CompletionRequest{BestOf: intPtr(5), N: intPtr(2)},
		},
		{
			name:      "invalid provider sort",
			req:       &CompletionRequest{Provider: &Provider{Sort: "random"}},
			wantErr:   true,
			wantField: "provider.sort",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCompletionParams(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCompletionParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				valErr, ok := IsValidationError(err)
				if !ok {
					t.Errorf("expected ValidationError, got %T", err)
					return
				}
				if valErr.Field != tt.wantField {
					t.Errorf("expected field %q, got %q", tt.wantField, valErr.Field)
				}
			}
		})
	}
}
