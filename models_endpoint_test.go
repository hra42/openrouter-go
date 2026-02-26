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
					CanonicalSlug: new("openai-gpt-4-turbo"),
					Created:       1234567890.0,
					Description:   "GPT-4 Turbo model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"text"},
						Tokenizer:        "cl100k_base",
						InstructType:     new("chat"),
					},
					TopProvider: ModelTopProvider{
						ContextLength:       &contextLength,
						MaxCompletionTokens: &maxCompTokens,
						IsModerated:         true,
					},
					PerRequestLimits:    nil,
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
		_ = json.NewEncoder(w).Encode(response)
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
					CanonicalSlug: new("anthropic-claude-3.5-sonnet"),
					Created:       1234567890.0,
					Description:   "Claude 3.5 Sonnet model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"text"},
						Tokenizer:        "claude",
						InstructType:     new("chat"),
					},
					TopProvider: ModelTopProvider{
						ContextLength:       &contextLength,
						MaxCompletionTokens: &maxCompTokens,
						IsModerated:         false,
					},
					PerRequestLimits:    nil,
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
		_ = json.NewEncoder(w).Encode(response)
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

func TestListModelsWithSupportedParameters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/models" {
			t.Errorf("expected path /models, got %s", r.URL.Path)
		}

		// Verify query parameter
		supportedParams := r.URL.Query().Get("supported_parameters")
		if supportedParams != "tools" {
			t.Errorf("expected supported_parameters 'tools', got %q", supportedParams)
		}

		// Send empty response
		response := ModelsResponse{Data: []Model{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListModels(context.Background(), &ListModelsOptions{
		SupportedParameters: "tools",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListModelsWithMultipleOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		// Verify query parameters
		query := r.URL.Query()
		if query.Get("category") != "programming" {
			t.Errorf("expected category 'programming', got %q", query.Get("category"))
		}
		if query.Get("supported_parameters") != "tools" {
			t.Errorf("expected supported_parameters 'tools', got %q", query.Get("supported_parameters"))
		}

		// Send empty response
		response := ModelsResponse{Data: []Model{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListModels(context.Background(), &ListModelsOptions{
		Category:            "programming",
		SupportedParameters: "tools",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListModelsUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/models/user" {
			t.Errorf("expected path /models/user, got %s", r.URL.Path)
		}

		// Verify auth header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization 'Bearer test-key', got %q", auth)
		}

		// Send response with user-filtered models
		temperature := 0.7
		contextLength := 128000.0
		maxCompTokens := 16384.0
		expirationDate := "2026-06-01"
		promptTokens := 100000.0
		completionTokens := 16384.0
		modality := "text->text"

		response := ModelsResponse{
			Data: []Model{
				{
					ID:            "openai/gpt-4o",
					Name:          "GPT-4o",
					CanonicalSlug: new("openai-gpt-4o"),
					Created:       1234567890.0,
					Description:   "GPT-4o model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text", "image"},
						OutputModalities: []string{"text"},
						Tokenizer:        "GPT",
						InstructType:     new("chatml"),
						Modality:         &modality,
					},
					TopProvider: ModelTopProvider{
						ContextLength:       &contextLength,
						MaxCompletionTokens: &maxCompTokens,
						IsModerated:         true,
					},
					PerRequestLimits: &ModelPerRequestLimits{
						PromptTokens:     &promptTokens,
						CompletionTokens: &completionTokens,
					},
					SupportedParameters: []string{"temperature", "top_p", "max_tokens", "tools"},
					DefaultParameters: &ModelDefaultParameters{
						Temperature: &temperature,
					},
					Pricing: ModelPricing{
						Prompt:     "0.005",
						Completion: "0.015",
						Image:      "0.003",
						Request:    "0",
					},
					ExpirationDate: &expirationDate,
				},
				{
					ID:            "anthropic/claude-3.5-sonnet",
					Name:          "Claude 3.5 Sonnet",
					CanonicalSlug: new("anthropic-claude-3.5-sonnet"),
					Created:       1234567890.0,
					Description:   "Claude 3.5 Sonnet",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"text"},
						Tokenizer:        "Claude",
					},
					TopProvider: ModelTopProvider{
						IsModerated: false,
					},
					Pricing: ModelPricing{
						Prompt:     "0.003",
						Completion: "0.015",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.ListModelsUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 models, got %d", len(resp.Data))
	}

	// Verify first model
	model := resp.Data[0]
	if model.ID != "openai/gpt-4o" {
		t.Errorf("expected model ID 'openai/gpt-4o', got %q", model.ID)
	}
	if model.Name != "GPT-4o" {
		t.Errorf("expected model name 'GPT-4o', got %q", model.Name)
	}
	if model.TopProvider.IsModerated != true {
		t.Error("expected IsModerated to be true")
	}

	// Verify new fields
	if model.ExpirationDate == nil || *model.ExpirationDate != "2026-06-01" {
		t.Error("expected expiration date '2026-06-01'")
	}
	if model.Architecture.Modality == nil || *model.Architecture.Modality != "text->text" {
		t.Error("expected modality 'text->text'")
	}
	if model.PerRequestLimits == nil {
		t.Fatal("expected per_request_limits to be non-nil")
	}
	if model.PerRequestLimits.PromptTokens == nil || *model.PerRequestLimits.PromptTokens != 100000.0 {
		t.Error("expected prompt_tokens 100000")
	}
	if model.PerRequestLimits.CompletionTokens == nil || *model.PerRequestLimits.CompletionTokens != 16384.0 {
		t.Error("expected completion_tokens 16384")
	}

	// Verify second model has nil optional fields
	model2 := resp.Data[1]
	if model2.ExpirationDate != nil {
		t.Errorf("expected nil expiration date, got %v", model2.ExpirationDate)
	}
}

func TestListModelsUserEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/models/user" {
			t.Errorf("expected path /models/user, got %s", r.URL.Path)
		}

		response := ModelsResponse{Data: []Model{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.ListModelsUser(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 models, got %d", len(resp.Data))
	}
}

func TestListModelsUserUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Invalid API key",
				Type:    "authentication_error",
				Code:    "401",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListModelsUser(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	reqErr, ok := IsRequestError(err)
	if !ok {
		t.Errorf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code 401, got %d", reqErr.StatusCode)
	}
}

//go:fix inline
func stringPtr(s string) *string {
	return new(s)
}
