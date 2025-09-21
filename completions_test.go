package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/completions" {
			t.Errorf("expected path /completions, got %s", r.URL.Path)
		}

		// Parse request body
		var req CompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify request fields
		if req.Model != "gpt-3.5-turbo-instruct" {
			t.Errorf("expected model 'gpt-3.5-turbo-instruct', got %q", req.Model)
		}
		if req.Prompt != "Once upon a time" {
			t.Errorf("expected prompt 'Once upon a time', got %q", req.Prompt)
		}
		if req.Stream != false {
			t.Error("expected stream to be false")
		}

		// Send response
		response := CompletionResponse{
			ID:      "cmpl-123",
			Object:  "text_completion",
			Created: 1234567890,
			Model:   "gpt-3.5-turbo-instruct",
			Choices: []CompletionChoice{
				{
					Index:        0,
					Text:         " in a land far, far away",
					FinishReason: "stop",
				},
			},
			Usage: Usage{
				PromptTokens:     4,
				CompletionTokens: 7,
				TotalTokens:      11,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))

	resp, err := client.Complete(context.Background(), "Once upon a time",
		WithCompletionModel("gpt-3.5-turbo-instruct"),
		WithCompletionTemperature(0.7),
		WithCompletionMaxTokens(100),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "cmpl-123" {
		t.Errorf("expected ID 'cmpl-123', got %q", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Text != " in a land far, far away" {
		t.Errorf("unexpected response text: %q", resp.Choices[0].Text)
	}

	if resp.Usage.TotalTokens != 11 {
		t.Errorf("expected 11 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestCompleteValidation(t *testing.T) {
	tests := []struct {
		name          string
		apiKey        string
		prompt        string
		model         string
		expectedError error
	}{
		{
			name:          "missing API key",
			apiKey:        "",
			prompt:        "Test prompt",
			model:         "gpt-3.5-turbo-instruct",
			expectedError: ErrNoAPIKey,
		},
		{
			name:          "missing prompt",
			apiKey:        "test-key",
			prompt:        "",
			model:         "gpt-3.5-turbo-instruct",
			expectedError: ErrNoPrompt,
		},
		{
			name:          "missing model",
			apiKey:        "test-key",
			prompt:        "Test prompt",
			model:         "",
			expectedError: ErrNoModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey)

			var opts []CompletionOption
			if tt.model != "" {
				opts = append(opts, WithCompletionModel(tt.model))
			}

			_, err := client.Complete(context.Background(), tt.prompt, opts...)

			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestCompletionOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify all options were applied
		if req.Temperature == nil || *req.Temperature != 0.8 {
			t.Error("temperature not set correctly")
		}
		if req.TopP == nil || *req.TopP != 0.9 {
			t.Error("top_p not set correctly")
		}
		if req.MaxTokens == nil || *req.MaxTokens != 200 {
			t.Error("max_tokens not set correctly")
		}
		if len(req.Stop) != 2 || req.Stop[0] != "\n" {
			t.Error("stop sequences not set correctly")
		}
		if req.LogProbs == nil || *req.LogProbs != 5 {
			t.Error("logprobs not set correctly")
		}
		if req.Echo == nil || !*req.Echo {
			t.Error("echo not set correctly")
		}
		if req.N == nil || *req.N != 3 {
			t.Error("n not set correctly")
		}
		if req.BestOf == nil || *req.BestOf != 5 {
			t.Error("best_of not set correctly")
		}
		if req.Suffix == nil || *req.Suffix != " The end." {
			t.Error("suffix not set correctly")
		}

		// Metadata should be in headers
		if r.Header.Get("X-custom-field") != "custom-value" {
			t.Error("metadata not set in headers")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))

	_, err := client.Complete(context.Background(), "Test prompt",
		WithCompletionModel("test-model"),
		WithCompletionTemperature(0.8),
		WithCompletionTopP(0.9),
		WithCompletionMaxTokens(200),
		WithCompletionStop("\n", "END"),
		WithCompletionLogProbs(5),
		WithCompletionEcho(true),
		WithCompletionN(3),
		WithCompletionBestOf(5),
		WithCompletionSuffix(" The end."),
		WithCompletionMetadata(map[string]interface{}{
			"custom-field": "custom-value",
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompleteWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		expectedPrompt := "You are a helpful assistant.\n\nWhat is the capital of France?"
		if req.Prompt != expectedPrompt {
			t.Errorf("expected prompt %q, got %q", expectedPrompt, req.Prompt)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{
			ID: "test",
			Choices: []CompletionChoice{
				{Text: "The capital of France is Paris."},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key",
		WithBaseURL(server.URL),
		WithDefaultModel("test-model"),
	)

	resp, err := client.CompleteWithContext(
		context.Background(),
		"You are a helpful assistant.",
		"What is the capital of France?",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("no choices in response")
	}

	if resp.Choices[0].Text != "The capital of France is Paris." {
		t.Errorf("unexpected response: %q", resp.Choices[0].Text)
	}
}

func TestCompleteWithExamples(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		expectedPrompt := `Translate English to French.

Examples:
1. Hello -> Bonjour
2. Thank you -> Merci


Now: Good morning`

		if req.Prompt != expectedPrompt {
			t.Errorf("prompt mismatch:\nExpected:\n%q\nGot:\n%q", expectedPrompt, req.Prompt)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{
			ID: "test",
			Choices: []CompletionChoice{
				{Text: "Bonjour"},
			},
		})
	}))
	defer server.Close()

	client := NewClient("test-key",
		WithBaseURL(server.URL),
		WithDefaultModel("test-model"),
	)

	examples := []string{
		"Hello -> Bonjour",
		"Thank you -> Merci",
	}

	resp, err := client.CompleteWithExamples(
		context.Background(),
		"Translate English to French.",
		examples,
		"Good morning",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("no choices in response")
	}

	if resp.Choices[0].Text != "Bonjour" {
		t.Errorf("unexpected response: %q", resp.Choices[0].Text)
	}
}

func TestHandleCompletionModelSuffix(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name         string
		model        string
		expectedModel string
		expectedSort  string
	}{
		{
			name:         "nitro suffix",
			model:        "openai/gpt-3.5-turbo-instruct:nitro",
			expectedModel: "openai/gpt-3.5-turbo-instruct",
			expectedSort:  "throughput",
		},
		{
			name:         "floor suffix",
			model:        "openai/gpt-3.5-turbo-instruct:floor",
			expectedModel: "openai/gpt-3.5-turbo-instruct",
			expectedSort:  "price",
		},
		{
			name:         "no suffix",
			model:        "openai/gpt-3.5-turbo-instruct",
			expectedModel: "openai/gpt-3.5-turbo-instruct",
			expectedSort:  "",
		},
		{
			name:         "model with colon in name",
			model:        "org:custom:model",
			expectedModel: "org:custom:model",
			expectedSort:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &CompletionRequest{
				Model: tt.model,
			}

			actualModel := client.handleCompletionModelSuffix(tt.model, req)

			if actualModel != tt.expectedModel {
				t.Errorf("expected model %q, got %q", tt.expectedModel, actualModel)
			}

			if tt.expectedSort != "" {
				if req.Provider == nil {
					t.Error("expected Provider to be set")
				} else if req.Provider.Sort != tt.expectedSort {
					t.Errorf("expected Sort %q, got %q", tt.expectedSort, req.Provider.Sort)
				}
			} else {
				if req.Provider != nil && req.Provider.Sort != "" {
					t.Errorf("expected no Sort, got %q", req.Provider.Sort)
				}
			}
		})
	}
}