package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/models" {
			t.Errorf("expected path /models, got %s", r.URL.Path)
		}

		// Send response
		temperature := 0.7
		topP := 0.9
		freqPenalty := 0.0
		contextLength := 128000.0
		maxCompTokens := 16384.0
		response := ModelsResponse{
			Data: []Model{
				{
					ID:            "openai/gpt-4-turbo",
					Name:          "GPT-4 Turbo",
					CanonicalSlug: stringPtr("openai-gpt-4-turbo"),
					Created:       1234567890.0,
					Description:   "GPT-4 Turbo model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"text"},
						Tokenizer:        "cl100k_base",
						InstructType:     stringPtr("chat"),
					},
					TopProvider: ModelTopProvider{
						ContextLength:       &contextLength,
						MaxCompletionTokens: &maxCompTokens,
						IsModerated:         true,
					},
					PerRequestLimits: nil,
					SupportedParameters: []string{"temperature", "top_p", "max_tokens"},
					DefaultParameters: &ModelDefaultParameters{
						Temperature:      &temperature,
						TopP:             &topP,
						FrequencyPenalty: &freqPenalty,
					},
					Pricing: ModelPricing{
						Prompt:            "0.01",
						Completion:        "0.03",
						Image:             "0",
						Request:           "0",
						InputCacheRead:    nil,
						InputCacheWrite:   nil,
						WebSearch:         "0",
						InternalReasoning: "0",
					},
				},
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

	resp, err := client.ListModels(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("expected 1 model, got %d", len(resp.Data))
	}

	model := resp.Data[0]
	if model.ID != "openai/gpt-4-turbo" {
		t.Errorf("expected model ID 'openai/gpt-4-turbo', got %q", model.ID)
	}
	if model.Name != "GPT-4 Turbo" {
		t.Errorf("expected model name 'GPT-4 Turbo', got %q", model.Name)
	}
	if model.TopProvider.IsModerated != true {
		t.Error("expected IsModerated to be true")
	}
}

func TestListModelsWithCategory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/models" {
			t.Errorf("expected path /models, got %s", r.URL.Path)
		}

		// Verify query parameter
		category := r.URL.Query().Get("category")
		if category != "programming" {
			t.Errorf("expected category 'programming', got %q", category)
		}

		// Send response
		temperature := 0.7
		topP := 0.9
		freqPenalty := 0.0
		contextLength := 128000.0
		maxCompTokens := 16384.0
		response := ModelsResponse{
			Data: []Model{
				{
					ID:            "anthropic/claude-3.5-sonnet",
					Name:          "Claude 3.5 Sonnet",
					CanonicalSlug: stringPtr("anthropic-claude-3.5-sonnet"),
					Created:       1234567890.0,
					Description:   "Claude 3.5 Sonnet model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"text"},
						Tokenizer:        "claude",
						InstructType:     stringPtr("chat"),
					},
					TopProvider: ModelTopProvider{
						ContextLength:       &contextLength,
						MaxCompletionTokens: &maxCompTokens,
						IsModerated:         false,
					},
					PerRequestLimits: nil,
					SupportedParameters: []string{"temperature", "top_p", "max_tokens"},
					DefaultParameters: &ModelDefaultParameters{
						Temperature:      &temperature,
						TopP:             &topP,
						FrequencyPenalty: &freqPenalty,
					},
					Pricing: ModelPricing{
						Prompt:            "0.003",
						Completion:        "0.015",
						Image:             "0",
						Request:           "0",
						InputCacheRead:    nil,
						InputCacheWrite:   nil,
						WebSearch:         "0",
						InternalReasoning: "0",
					},
				},
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

	resp, err := client.ListModels(context.Background(), &ListModelsOptions{
		Category: "programming",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("expected 1 model, got %d", len(resp.Data))
	}

	model := resp.Data[0]
	if model.ID != "anthropic/claude-3.5-sonnet" {
		t.Errorf("expected model ID 'anthropic/claude-3.5-sonnet', got %q", model.ID)
	}
}

func stringPtr(s string) *string {
	return &s
}
