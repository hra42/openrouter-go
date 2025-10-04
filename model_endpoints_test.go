package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListModelEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/models/openai/gpt-4/endpoints" {
			t.Errorf("expected path /models/openai/gpt-4/endpoints, got %s", r.URL.Path)
		}

		// Send response
		contextLength := 8192.0
		maxCompTokens := 4096.0
		maxPromptTokens := 4096.0
		uptime := 0.99
		quantization := "fp16"

		response := ModelEndpointsResponse{
			Data: ModelEndpointsData{
				ID:          "openai/gpt-4",
				Name:        "GPT-4",
				Created:     1234567890.0,
				Description: "GPT-4 model",
				Architecture: ModelEndpointsArchitecture{
					Tokenizer:        stringPtr("cl100k_base"),
					InstructType:     stringPtr("chat"),
					InputModalities:  []string{"text"},
					OutputModalities: []string{"text"},
				},
				Endpoints: []ModelEndpoint{
					{
						Name:                "openai/gpt-4",
						ContextLength:       contextLength,
						ProviderName:        "OpenAI",
						Quantization:        &quantization,
						MaxCompletionTokens: &maxCompTokens,
						MaxPromptTokens:     &maxPromptTokens,
						SupportedParameters: []string{"temperature", "top_p", "max_tokens"},
						Status:              1,
						UptimeLast30m:       &uptime,
						Pricing: ModelEndpointPricing{
							Request:    "0",
							Image:      "0",
							Prompt:     "0.03",
							Completion: "0.06",
						},
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

	resp, err := client.ListModelEndpoints(context.Background(), "openai", "gpt-4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response data
	if resp.Data.ID != "openai/gpt-4" {
		t.Errorf("expected ID 'openai/gpt-4', got %q", resp.Data.ID)
	}
	if resp.Data.Name != "GPT-4" {
		t.Errorf("expected name 'GPT-4', got %q", resp.Data.Name)
	}
	if resp.Data.Description != "GPT-4 model" {
		t.Errorf("expected description 'GPT-4 model', got %q", resp.Data.Description)
	}

	// Verify architecture
	if resp.Data.Architecture.Tokenizer == nil || *resp.Data.Architecture.Tokenizer != "cl100k_base" {
		t.Error("expected tokenizer 'cl100k_base'")
	}
	if len(resp.Data.Architecture.InputModalities) != 1 || resp.Data.Architecture.InputModalities[0] != "text" {
		t.Errorf("expected input modalities ['text'], got %v", resp.Data.Architecture.InputModalities)
	}

	// Verify endpoints
	if len(resp.Data.Endpoints) != 1 {
		t.Errorf("expected 1 endpoint, got %d", len(resp.Data.Endpoints))
	}

	endpoint := resp.Data.Endpoints[0]
	if endpoint.Name != "openai/gpt-4" {
		t.Errorf("expected endpoint name 'openai/gpt-4', got %q", endpoint.Name)
	}
	if endpoint.ProviderName != "OpenAI" {
		t.Errorf("expected provider name 'OpenAI', got %q", endpoint.ProviderName)
	}
	if endpoint.Status != 1 {
		t.Errorf("expected status 1, got %f", endpoint.Status)
	}
	if endpoint.ContextLength != 8192.0 {
		t.Errorf("expected context length 8192.0, got %f", endpoint.ContextLength)
	}

	// Verify pricing
	if endpoint.Pricing.Prompt != "0.03" {
		t.Errorf("expected prompt pricing '0.03', got %q", endpoint.Pricing.Prompt)
	}
	if endpoint.Pricing.Completion != "0.06" {
		t.Errorf("expected completion pricing '0.06', got %q", endpoint.Pricing.Completion)
	}

	// Verify optional fields
	if endpoint.Quantization == nil || *endpoint.Quantization != "fp16" {
		t.Error("expected quantization 'fp16'")
	}
	if endpoint.UptimeLast30m == nil || *endpoint.UptimeLast30m != 0.99 {
		t.Error("expected uptime 0.99")
	}
}

func TestListModelEndpointsEmptyAuthor(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.ListModelEndpoints(context.Background(), "", "gpt-4")
	if err == nil {
		t.Fatal("expected error for empty author, got nil")
	}
	if err.Error() != "author cannot be empty" {
		t.Errorf("expected 'author cannot be empty' error, got %q", err.Error())
	}
}

func TestListModelEndpointsEmptySlug(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.ListModelEndpoints(context.Background(), "openai", "")
	if err == nil {
		t.Fatal("expected error for empty slug, got nil")
	}
	if err.Error() != "slug cannot be empty" {
		t.Errorf("expected 'slug cannot be empty' error, got %q", err.Error())
	}
}

func TestListModelEndpointsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Model not found",
				Type:    "not_found",
				Code:    "404",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListModelEndpoints(context.Background(), "invalid", "nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent model, got nil")
	}

	if !IsRequestError(err) {
		t.Errorf("expected RequestError, got %T", err)
	}

	reqErr := err.(*RequestError)
	if reqErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", reqErr.StatusCode)
	}
}

func TestListModelEndpointsMultipleEndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models/anthropic/claude-3.5-sonnet/endpoints" {
			t.Errorf("expected path /models/anthropic/claude-3.5-sonnet/endpoints, got %s", r.URL.Path)
		}

		contextLength1 := 200000.0
		contextLength2 := 200000.0
		maxCompTokens := 8192.0
		uptime1 := 0.995
		uptime2 := 0.98

		response := ModelEndpointsResponse{
			Data: ModelEndpointsData{
				ID:          "anthropic/claude-3.5-sonnet",
				Name:        "Claude 3.5 Sonnet",
				Created:     1234567890.0,
				Description: "Claude 3.5 Sonnet model",
				Architecture: ModelEndpointsArchitecture{
					Tokenizer:        stringPtr("claude"),
					InstructType:     stringPtr("chat"),
					InputModalities:  []string{"text", "image"},
					OutputModalities: []string{"text"},
				},
				Endpoints: []ModelEndpoint{
					{
						Name:                "anthropic/claude-3.5-sonnet",
						ContextLength:       contextLength1,
						ProviderName:        "Anthropic",
						MaxCompletionTokens: &maxCompTokens,
						SupportedParameters: []string{"temperature", "top_p", "max_tokens"},
						Status:              1,
						UptimeLast30m:       &uptime1,
						Pricing: ModelEndpointPricing{
							Request:    "0",
							Image:      "0.004",
							Prompt:     "0.003",
							Completion: "0.015",
						},
					},
					{
						Name:                "anthropic/claude-3.5-sonnet",
						ContextLength:       contextLength2,
						ProviderName:        "AWS Bedrock",
						MaxCompletionTokens: &maxCompTokens,
						SupportedParameters: []string{"temperature", "top_p"},
						Status:              1,
						UptimeLast30m:       &uptime2,
						Pricing: ModelEndpointPricing{
							Request:    "0",
							Image:      "0.0048",
							Prompt:     "0.0035",
							Completion: "0.018",
						},
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

	resp, err := client.ListModelEndpoints(context.Background(), "anthropic", "claude-3.5-sonnet")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify we got multiple endpoints
	if len(resp.Data.Endpoints) != 2 {
		t.Errorf("expected 2 endpoints, got %d", len(resp.Data.Endpoints))
	}

	// Verify both providers are present
	providers := make(map[string]bool)
	for _, endpoint := range resp.Data.Endpoints {
		providers[endpoint.ProviderName] = true
	}

	if !providers["Anthropic"] {
		t.Error("expected Anthropic provider")
	}
	if !providers["AWS Bedrock"] {
		t.Error("expected AWS Bedrock provider")
	}

	// Verify multimodal support
	if len(resp.Data.Architecture.InputModalities) != 2 {
		t.Errorf("expected 2 input modalities, got %d", len(resp.Data.Architecture.InputModalities))
	}
}
