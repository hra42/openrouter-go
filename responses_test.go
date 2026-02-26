package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateResponse_BasicText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/responses" {
			t.Errorf("expected path /responses, got %s", r.URL.Path)
		}

		// Parse request body
		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify request fields
		if req.Model != "openai/o4-mini" {
			t.Errorf("expected model 'openai/o4-mini', got %q", req.Model)
		}
		if input, ok := req.Input.(string); !ok || input != "Hello, world!" {
			t.Errorf("expected input 'Hello, world!', got %v", req.Input)
		}
		if req.Stream {
			t.Error("expected stream to be false")
		}

		// Send response
		response := ResponsesResponse{
			ID:        "resp_123",
			Object:    "response",
			CreatedAt: 1234567890,
			Model:     "openai/o4-mini",
			Output: []ResponsesOutput{
				{
					Type:   "message",
					ID:     "msg_123",
					Status: "completed",
					Role:   "assistant",
					Content: []ResponsesOutputContent{
						{
							Type: "output_text",
							Text: "Hello! How can I help you today?",
						},
					},
				},
			},
			Usage: ResponsesUsage{
				InputTokens:  5,
				OutputTokens: 8,
				TotalTokens:  13,
			},
			Status: "completed",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateResponse(context.Background(), "Hello, world!",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesMaxOutputTokens(100),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "resp_123" {
		t.Errorf("expected ID 'resp_123', got %q", resp.ID)
	}

	if len(resp.Output) != 1 {
		t.Fatalf("expected 1 output, got %d", len(resp.Output))
	}

	if resp.GetTextContent() != "Hello! How can I help you today?" {
		t.Errorf("unexpected response content: %q", resp.GetTextContent())
	}

	if resp.Usage.TotalTokens != 13 {
		t.Errorf("expected 13 total tokens, got %d", resp.Usage.TotalTokens)
	}
}

func TestCreateResponse_StructuredInput(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify structured input (it arrives as a generic slice/map due to JSON)
		input, ok := req.Input.([]interface{})
		if !ok {
			t.Fatalf("expected input to be array, got %T", req.Input)
		}
		if len(input) != 2 {
			t.Errorf("expected 2 input items, got %d", len(input))
		}

		response := ResponsesResponse{
			ID:        "resp_456",
			Object:    "response",
			CreatedAt: 1234567890,
			Model:     "openai/o4-mini",
			Output: []ResponsesOutput{
				{
					Type:   "message",
					ID:     "msg_456",
					Status: "completed",
					Role:   "assistant",
					Content: []ResponsesOutputContent{
						{
							Type: "output_text",
							Text: "Your name is Alice. Nice to meet you!",
						},
					},
				},
			},
			Usage: ResponsesUsage{
				InputTokens:  20,
				OutputTokens: 12,
				TotalTokens:  32,
			},
			Status: "completed",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	input := []ResponsesInputItem{
		CreateResponsesUserMessage("My name is Alice"),
		CreateResponsesUserMessage("What is my name?"),
	}

	resp, err := client.CreateResponse(context.Background(), input,
		WithResponsesModel("openai/o4-mini"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "resp_456" {
		t.Errorf("expected ID 'resp_456', got %q", resp.ID)
	}

	expectedText := "Your name is Alice. Nice to meet you!"
	if resp.GetTextContent() != expectedText {
		t.Errorf("expected %q, got %q", expectedText, resp.GetTextContent())
	}
}

func TestCreateResponse_WithReasoning(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify reasoning is set
		if req.Reasoning == nil {
			t.Error("expected reasoning to be set")
		} else if req.Reasoning.Effort != "high" {
			t.Errorf("expected reasoning effort 'high', got %q", req.Reasoning.Effort)
		}

		response := ResponsesResponse{
			ID:     "resp_789",
			Object: "response",
			Model:  "openai/o4-mini",
			Output: []ResponsesOutput{
				{
					Type: "message",
					ID:   "msg_789",
					Role: "assistant",
					Content: []ResponsesOutputContent{
						{
							Type:             "reasoning",
							EncryptedContent: "encrypted-reasoning-chain",
							Summary:          []string{"Step 1: Analyze input", "Step 2: Calculate result"},
						},
						{
							Type: "output_text",
							Text: "The answer is 42.",
						},
					},
				},
			},
			Usage: ResponsesUsage{
				InputTokens:     10,
				OutputTokens:    20,
				TotalTokens:     30,
				ReasoningTokens: 15,
			},
			Status: "completed",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateResponse(context.Background(), "What is the meaning of life?",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesReasoningEffort("high"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	summary := resp.GetReasoningSummary()
	if len(summary) != 2 {
		t.Errorf("expected 2 reasoning summary items, got %d", len(summary))
	}

	if resp.Usage.ReasoningTokens != 15 {
		t.Errorf("expected 15 reasoning tokens, got %d", resp.Usage.ReasoningTokens)
	}
}

func TestCreateResponse_WithTools(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify tools are set (ResponsesTool has flat structure with Name at top level)
		if len(req.Tools) != 1 {
			t.Errorf("expected 1 tool, got %d", len(req.Tools))
		} else {
			if req.Tools[0].Name != "get_weather" {
				t.Errorf("expected tool name 'get_weather', got %q", req.Tools[0].Name)
			}
		}

		response := ResponsesResponse{
			ID:     "resp_tools",
			Object: "response",
			Model:  "openai/o4-mini",
			Output: []ResponsesOutput{
				{
					Type:      "function_call",
					ID:        "fc_123",
					CallID:    "call_abc",
					Name:      "get_weather",
					Arguments: `{"location":"San Francisco"}`,
				},
			},
			Usage: ResponsesUsage{
				InputTokens:  15,
				OutputTokens: 10,
				TotalTokens:  25,
			},
			Status: "completed",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	// Responses API uses a flat tool structure (name, description, parameters at top level)
	weatherTool := CreateResponsesTool(
		"get_weather",
		"Get weather for a location",
		map[string]any{
			"type": "object",
			"properties": map[string]any{
				"location": map[string]any{
					"type":        "string",
					"description": "The city name",
				},
			},
			"required": []string{"location"},
		},
	)

	resp, err := client.CreateResponse(context.Background(), "What's the weather in San Francisco?",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesTools(weatherTool),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	calls := resp.GetFunctionCalls()
	if len(calls) != 1 {
		t.Fatalf("expected 1 function call, got %d", len(calls))
	}

	if calls[0].Name != "get_weather" {
		t.Errorf("expected function name 'get_weather', got %q", calls[0].Name)
	}

	if calls[0].Arguments != `{"location":"San Francisco"}` {
		t.Errorf("unexpected arguments: %q", calls[0].Arguments)
	}
}

func TestCreateResponse_WithWebSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify plugins are set
		if len(req.Plugins) != 1 {
			t.Errorf("expected 1 plugin, got %d", len(req.Plugins))
		} else {
			if req.Plugins[0].ID != "web" {
				t.Errorf("expected plugin id 'web', got %q", req.Plugins[0].ID)
			}
			if req.Plugins[0].MaxResults != 3 {
				t.Errorf("expected max_results 3, got %d", req.Plugins[0].MaxResults)
			}
		}

		response := ResponsesResponse{
			ID:     "resp_search",
			Object: "response",
			Model:  "openai/o4-mini",
			Output: []ResponsesOutput{
				{
					Type: "message",
					ID:   "msg_search",
					Role: "assistant",
					Content: []ResponsesOutputContent{
						{
							Type: "output_text",
							Text: "According to recent sources, the answer is...",
							Annotations: []ResponsesAnnotation{
								{
									Type:       "url_citation",
									URL:        "https://example.com/source",
									StartIndex: 0,
									EndIndex:   30,
								},
							},
						},
					},
				},
			},
			Usage: ResponsesUsage{
				InputTokens:  10,
				OutputTokens: 15,
				TotalTokens:  25,
			},
			Status: "completed",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateResponse(context.Background(), "What is the latest news?",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesWebSearch(3),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	annotations := resp.GetAnnotations()
	if len(annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(annotations))
	}

	if annotations[0].Type != "url_citation" {
		t.Errorf("expected annotation type 'url_citation', got %q", annotations[0].Type)
	}

	if annotations[0].URL != "https://example.com/source" {
		t.Errorf("unexpected annotation URL: %q", annotations[0].URL)
	}
}

func TestCreateResponse_Validation(t *testing.T) {
	tests := []struct {
		name          string
		apiKey        string
		input         interface{}
		model         string
		expectedError error
	}{
		{
			name:          "missing API key",
			apiKey:        "",
			input:         "Hello",
			model:         "openai/o4-mini",
			expectedError: ErrNoAPIKey,
		},
		{
			name:          "nil input",
			apiKey:        "test-key",
			input:         nil,
			model:         "openai/o4-mini",
			expectedError: nil, // will be ValidationError
		},
		{
			name:          "empty string input",
			apiKey:        "test-key",
			input:         "",
			model:         "openai/o4-mini",
			expectedError: nil, // will be ValidationError
		},
		{
			name:          "empty array input",
			apiKey:        "test-key",
			input:         []ResponsesInputItem{},
			model:         "openai/o4-mini",
			expectedError: nil, // will be ValidationError
		},
		{
			name:          "missing model",
			apiKey:        "test-key",
			input:         "Hello",
			model:         "",
			expectedError: ErrNoModel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(WithAPIKey(tt.apiKey))

			var opts []ResponsesOption
			if tt.model != "" {
				opts = append(opts, WithResponsesModel(tt.model))
			}

			_, err := client.CreateResponse(context.Background(), tt.input, opts...)

			if err == nil {
				t.Error("expected error but got nil")
				return
			}

			if tt.expectedError != nil && err != tt.expectedError {
				// Check if it's the right type of validation error
				if _, ok := err.(*ValidationError); !ok {
					t.Errorf("expected ValidationError or %v, got %v", tt.expectedError, err)
				}
			}
		})
	}
}

func TestCreateResponse_InvalidReasoningEffort(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.CreateResponse(context.Background(), "Hello",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesReasoningEffort("invalid"),
	)

	if err == nil {
		t.Error("expected error for invalid reasoning effort")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if valErr.Field != "reasoning.effort" {
		t.Errorf("expected field 'reasoning.effort', got %q", valErr.Field)
	}
}

func TestCreateResponse_InvalidInputRole(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	input := []ResponsesInputItem{
		{
			Type: "message",
			Role: "invalid_role",
			Content: []ResponsesInputContent{
				{Type: "input_text", Text: "Hello"},
			},
		},
	}

	_, err := client.CreateResponse(context.Background(), input,
		WithResponsesModel("openai/o4-mini"),
	)

	if err == nil {
		t.Error("expected error for invalid role")
		return
	}

	valErr, ok := err.(*ValidationError)
	if !ok {
		t.Errorf("expected ValidationError, got %T", err)
		return
	}

	if valErr.Field != "input[0].role" {
		t.Errorf("expected field 'input[0].role', got %q", valErr.Field)
	}
}

func TestCreateResponse_Options(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ResponsesRequest
		_ = json.NewDecoder(r.Body).Decode(&req)

		// Verify all options were applied
		if req.Temperature == nil || *req.Temperature != 0.8 {
			t.Error("temperature not set correctly")
		}
		if req.TopP == nil || *req.TopP != 0.9 {
			t.Error("top_p not set correctly")
		}
		if req.MaxOutputTokens == nil || *req.MaxOutputTokens != 200 {
			t.Error("max_output_tokens not set correctly")
		}

		// Metadata should be in headers
		if r.Header.Get("X-custom-field") != "custom-value" {
			t.Error("metadata not set in headers")
		}

		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(ResponsesResponse{ID: "test", Status: "completed"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.CreateResponse(context.Background(), "Test",
		WithResponsesModel("openai/o4-mini"),
		WithResponsesTemperature(0.8),
		WithResponsesTopP(0.9),
		WithResponsesMaxOutputTokens(200),
		WithResponsesMetadata(map[string]interface{}{
			"custom-field": "custom-value",
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateResponseStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/responses" {
			t.Errorf("expected path /responses, got %s", r.URL.Path)
		}

		var req ResponsesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify stream is enabled
		if !req.Stream {
			t.Error("expected stream to be true")
		}

		// Set headers for SSE
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected http.ResponseWriter to be an http.Flusher")
		}

		// Send streaming events
		events := []string{
			`data: {"id":"resp_stream","object":"response","model":"openai/o4-mini","output":[{"type":"message","id":"msg_1","role":"assistant","content":[{"type":"output_text","text":"Hello"}]}],"usage":{"input_tokens":5,"output_tokens":1,"total_tokens":6},"status":"in_progress"}`,
			`data: {"id":"resp_stream","object":"response","model":"openai/o4-mini","output":[{"type":"message","id":"msg_1","role":"assistant","content":[{"type":"output_text","text":"Hello world"}]}],"usage":{"input_tokens":5,"output_tokens":2,"total_tokens":7},"status":"in_progress"}`,
			`data: {"id":"resp_stream","object":"response","model":"openai/o4-mini","output":[{"type":"message","id":"msg_1","role":"assistant","content":[{"type":"output_text","text":"Hello world!"}]}],"usage":{"input_tokens":5,"output_tokens":3,"total_tokens":8},"status":"completed"}`,
			`data: [DONE]`,
		}

		for _, event := range events {
			_, _ = w.Write([]byte(event + "\n\n"))
			flusher.Flush()
		}
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	stream, err := client.CreateResponseStream(context.Background(), "Hello",
		WithResponsesModel("openai/o4-mini"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer stream.Close() //nolint:errcheck

	eventCount := 0
	var lastEvent ResponsesResponse
	for event := range stream.Events() {
		eventCount++
		lastEvent = event
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if eventCount == 0 {
		t.Error("expected to receive at least one event")
	}

	// Verify final event
	if lastEvent.Status != "completed" {
		t.Errorf("expected final status 'completed', got %q", lastEvent.Status)
	}

	if lastEvent.GetTextContent() != "Hello world!" {
		t.Errorf("expected text 'Hello world!', got %q", lastEvent.GetTextContent())
	}
}

func TestResponsesInputHelpers(t *testing.T) {
	// Test CreateResponsesUserMessage
	msg := CreateResponsesUserMessage("Hello")
	if msg.Type != "message" || msg.Role != "user" {
		t.Error("CreateResponsesUserMessage failed")
	}
	if len(msg.Content) != 1 || msg.Content[0].Text != "Hello" {
		t.Error("CreateResponsesUserMessage content incorrect")
	}

	// Test CreateResponsesAssistantMessage
	msg = CreateResponsesAssistantMessage("Hi there")
	if msg.Type != "message" || msg.Role != "assistant" {
		t.Error("CreateResponsesAssistantMessage failed")
	}

	// Test CreateResponsesSystemMessage
	msg = CreateResponsesSystemMessage("You are helpful")
	if msg.Type != "message" || msg.Role != "system" {
		t.Error("CreateResponsesSystemMessage failed")
	}

	// Test CreateResponsesMessage
	msg = CreateResponsesMessage("user", "Test")
	if msg.Type != "message" || msg.Role != "user" || msg.Content[0].Text != "Test" {
		t.Error("CreateResponsesMessage failed")
	}

	// Test CreateResponsesFunctionOutput
	funcOutput := CreateResponsesFunctionOutput("call_123", `{"result":"success"}`)
	if funcOutput.Type != "function_call_output" {
		t.Error("CreateResponsesFunctionOutput type incorrect")
	}
	if funcOutput.CallID != "call_123" {
		t.Error("CreateResponsesFunctionOutput call_id incorrect")
	}
	if funcOutput.Output != `{"result":"success"}` {
		t.Error("CreateResponsesFunctionOutput output incorrect")
	}
}

func TestResponsesResponse_Helpers(t *testing.T) {
	resp := ResponsesResponse{
		Output: []ResponsesOutput{
			{
				Type: "message",
				Role: "assistant",
				Content: []ResponsesOutputContent{
					{
						Type:    "reasoning",
						Summary: []string{"Step 1", "Step 2"},
					},
					{
						Type: "output_text",
						Text: "The answer",
						Annotations: []ResponsesAnnotation{
							{Type: "url_citation", URL: "https://example.com"},
						},
					},
				},
			},
			{
				Type:      "function_call",
				Name:      "my_function",
				CallID:    "call_1",
				Arguments: `{"arg":"value"}`,
			},
		},
	}

	// Test GetTextContent
	if resp.GetTextContent() != "The answer" {
		t.Errorf("GetTextContent failed: %q", resp.GetTextContent())
	}

	// Test GetFunctionCalls
	calls := resp.GetFunctionCalls()
	if len(calls) != 1 || calls[0].Name != "my_function" {
		t.Error("GetFunctionCalls failed")
	}

	// Test GetAnnotations
	annotations := resp.GetAnnotations()
	if len(annotations) != 1 || annotations[0].URL != "https://example.com" {
		t.Error("GetAnnotations failed")
	}

	// Test GetReasoningSummary
	summary := resp.GetReasoningSummary()
	if len(summary) != 2 || summary[0] != "Step 1" {
		t.Error("GetReasoningSummary failed")
	}
}

func TestReasoningEffortConstants(t *testing.T) {
	if ReasoningEffortMinimal != "minimal" {
		t.Errorf("expected 'minimal', got %q", ReasoningEffortMinimal)
	}
	if ReasoningEffortLow != "low" {
		t.Errorf("expected 'low', got %q", ReasoningEffortLow)
	}
	if ReasoningEffortMedium != "medium" {
		t.Errorf("expected 'medium', got %q", ReasoningEffortMedium)
	}
	if ReasoningEffortHigh != "high" {
		t.Errorf("expected 'high', got %q", ReasoningEffortHigh)
	}
}
