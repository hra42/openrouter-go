package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChatComplete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected path /chat/completions, got %s", r.URL.Path)
		}

		// Parse request body
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify request fields
		if req.Model != "gpt-3.5-turbo" {
			t.Errorf("expected model 'gpt-3.5-turbo', got %q", req.Model)
		}
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(req.Messages))
		}
		if req.Stream != false {
			t.Error("expected stream to be false")
		}

		// Send response
		response := ChatCompletionResponse{
			ID:      "chat-123",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "gpt-3.5-turbo",
			Choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    "assistant",
						Content: "Hello! How can I help you today?",
					},
					FinishReason: "stop",
				},
			},
			Usage: Usage{
				PromptTokens:     10,
				CompletionTokens: 8,
				TotalTokens:      18,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))

	messages := []Message{
		CreateSystemMessage("You are a helpful assistant."),
		CreateUserMessage("Hello!"),
	}

	resp, err := client.ChatComplete(context.Background(), messages,
		WithModel("gpt-3.5-turbo"),
		WithTemperature(0.7),
		WithMaxTokens(100),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "chat-123" {
		t.Errorf("expected ID 'chat-123', got %q", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Hello! How can I help you today?" {
		t.Errorf("unexpected response content: %q", resp.Choices[0].Message.Content)
	}

	if resp.Usage.TotalTokens != 18 {
		t.Errorf("expected 18 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestChatCompleteValidation(t *testing.T) {
	client := NewClient("")

	tests := []struct {
		name          string
		apiKey        string
		messages      []Message
		model         string
		expectedError error
	}{
		{
			name:          "missing API key",
			apiKey:        "",
			messages:      []Message{{Role: "user", Content: "Hello"}},
			model:         "gpt-3.5-turbo",
			expectedError: ErrNoAPIKey,
		},
		{
			name:          "missing messages",
			apiKey:        "test-key",
			messages:      []Message{},
			model:         "gpt-3.5-turbo",
			expectedError: ErrNoMessages,
		},
		{
			name:          "missing model",
			apiKey:        "test-key",
			messages:      []Message{{Role: "user", Content: "Hello"}},
			model:         "",
			expectedError: ErrNoModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.apiKey = tt.apiKey

			var opts []ChatCompletionOption
			if tt.model != "" {
				opts = append(opts, WithModel(tt.model))
			}

			_, err := client.ChatComplete(context.Background(), tt.messages, opts...)

			if err != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestValidateChatRequest(t *testing.T) {
	client := NewClient("test-key")

	tests := []struct {
		name        string
		messages    []Message
		shouldError bool
		errorField  string
	}{
		{
			name: "valid messages",
			messages: []Message{
				{Role: "system", Content: "You are helpful"},
				{Role: "user", Content: "Hello"},
				{Role: "assistant", Content: "Hi there!"},
			},
			shouldError: false,
		},
		{
			name: "missing role",
			messages: []Message{
				{Content: "Hello"},
			},
			shouldError: true,
			errorField:  "messages[0].role",
		},
		{
			name: "invalid role",
			messages: []Message{
				{Role: "invalid", Content: "Hello"},
			},
			shouldError: true,
			errorField:  "messages[0].role",
		},
		{
			name: "missing content for user message",
			messages: []Message{
				{Role: "user"},
			},
			shouldError: true,
			errorField:  "messages[0].content",
		},
		{
			name: "assistant message without content allowed",
			messages: []Message{
				{Role: "assistant", ToolCalls: []ToolCall{{ID: "1", Type: "function", Function: FunctionCall{Name: "test"}}}},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.validateChatRequest(tt.messages)

			if tt.shouldError && err == nil {
				t.Error("expected error but got nil")
			}

			if !tt.shouldError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.shouldError && err != nil {
				if valErr, ok := err.(*ValidationError); ok {
					if valErr.Field != tt.errorField {
						t.Errorf("expected error field %q, got %q", tt.errorField, valErr.Field)
					}
				}
			}
		})
	}
}

func TestChatCompletionOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
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
		if req.FrequencyPenalty == nil || *req.FrequencyPenalty != 0.5 {
			t.Error("frequency_penalty not set correctly")
		}
		if req.PresencePenalty == nil || *req.PresencePenalty != 0.6 {
			t.Error("presence_penalty not set correctly")
		}
		if req.Seed == nil || *req.Seed != 42 {
			t.Error("seed not set correctly")
		}
		if req.LogProbs == nil || !*req.LogProbs {
			t.Error("logprobs not enabled")
		}
		if req.TopLogProbs == nil || *req.TopLogProbs != 5 {
			t.Error("top_logprobs not set correctly")
		}

		// Metadata should be in headers
		if r.Header.Get("X-custom-field") != "custom-value" {
			t.Error("metadata not set in headers")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))

	messages := []Message{CreateUserMessage("Test")}

	_, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithTemperature(0.8),
		WithTopP(0.9),
		WithMaxTokens(200),
		WithStop("\n", "END"),
		WithFrequencyPenalty(0.5),
		WithPresencePenalty(0.6),
		WithSeed(42),
		WithLogProbs(5),
		WithMetadata(map[string]interface{}{
			"custom-field": "custom-value",
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMessageHelpers(t *testing.T) {
	// Test CreateSystemMessage
	msg := CreateSystemMessage("You are helpful")
	if msg.Role != "system" || msg.Content != "You are helpful" {
		t.Error("CreateSystemMessage failed")
	}

	// Test CreateUserMessage
	msg = CreateUserMessage("Hello")
	if msg.Role != "user" || msg.Content != "Hello" {
		t.Error("CreateUserMessage failed")
	}

	// Test CreateAssistantMessage
	msg = CreateAssistantMessage("Hi there")
	if msg.Role != "assistant" || msg.Content != "Hi there" {
		t.Error("CreateAssistantMessage failed")
	}

	// Test CreateToolMessage
	msg = CreateToolMessage("Result", "tool-123")
	if msg.Role != "tool" || msg.Content != "Result" || msg.ToolCallID != "tool-123" {
		t.Error("CreateToolMessage failed")
	}

	// Test CreateMultiModalMessage
	msg = CreateMultiModalMessage("user", "Describe this", "https://example.com/image.jpg")
	if msg.Role != "user" {
		t.Error("CreateMultiModalMessage failed: wrong role")
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok || len(parts) != 2 {
		t.Error("CreateMultiModalMessage failed: wrong content structure")
	}

	if parts[0].Type != "text" || parts[0].Text != "Describe this" {
		t.Error("CreateMultiModalMessage failed: wrong text part")
	}

	if parts[1].Type != "image_url" || parts[1].ImageURL.URL != "https://example.com/image.jpg" {
		t.Error("CreateMultiModalMessage failed: wrong image part")
	}
}

func TestHandleModelSuffix(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name         string
		model        string
		expectedModel string
		expectedSort  string
	}{
		{
			name:         "nitro suffix",
			model:        "meta-llama/llama-3.1-70b-instruct:nitro",
			expectedModel: "meta-llama/llama-3.1-70b-instruct",
			expectedSort:  "throughput",
		},
		{
			name:         "floor suffix",
			model:        "meta-llama/llama-3.1-70b-instruct:floor",
			expectedModel: "meta-llama/llama-3.1-70b-instruct",
			expectedSort:  "price",
		},
		{
			name:         "no suffix",
			model:        "meta-llama/llama-3.1-70b-instruct",
			expectedModel: "meta-llama/llama-3.1-70b-instruct",
			expectedSort:  "",
		},
		{
			name:         "model with colon but not suffix",
			model:        "custom:model:v1",
			expectedModel: "custom:model:v1",
			expectedSort:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &ChatCompletionRequest{
				Model: tt.model,
			}

			actualModel := client.handleModelSuffix(tt.model, req)

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