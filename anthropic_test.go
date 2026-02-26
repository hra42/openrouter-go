package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateAnthropicMessage(t *testing.T) {
	stopReason := "end_turn"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/messages" {
			t.Errorf("expected path /messages, got %s", r.URL.Path)
		}

		var req AnthropicMessagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Model != "anthropic/claude-sonnet-4" {
			t.Errorf("expected model 'anthropic/claude-sonnet-4', got %q", req.Model)
		}
		if req.MaxTokens != 1024 {
			t.Errorf("expected max_tokens 1024, got %d", req.MaxTokens)
		}
		if len(req.Messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(req.Messages))
		}
		if req.Stream {
			t.Error("expected stream to be false")
		}

		response := AnthropicMessagesResponse{
			ID:   "msg_123",
			Type: "message",
			Role: "assistant",
			Content: []AnthropicResponseContentBlock{
				{
					Type: "text",
					Text: "Hello! How can I help you?",
				},
			},
			Model:      "anthropic/claude-sonnet-4",
			StopReason: &stopReason,
			Usage: AnthropicUsage{
				InputTokens:  10,
				OutputTokens: 8,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	messages := []AnthropicMessage{
		CreateAnthropicUserMessage("Hello!"),
	}

	resp, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(1024),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "msg_123" {
		t.Errorf("expected ID 'msg_123', got %q", resp.ID)
	}

	if resp.Type != "message" {
		t.Errorf("expected type 'message', got %q", resp.Type)
	}

	content := resp.GetTextContent()
	if content != "Hello! How can I help you?" {
		t.Errorf("unexpected content: %q", content)
	}

	if resp.GetStopReason() != "end_turn" {
		t.Errorf("expected stop_reason 'end_turn', got %q", resp.GetStopReason())
	}

	if resp.Usage.InputTokens != 10 {
		t.Errorf("expected 10 input tokens, got %d", resp.Usage.InputTokens)
	}
}

func TestCreateAnthropicMessageWithSystemPrompt(t *testing.T) {
	tests := []struct {
		name   string
		system interface{}
	}{
		{
			name:   "string system",
			system: "You are a helpful assistant.",
		},
		{
			name: "blocks system",
			system: []AnthropicTextBlock{
				{Type: "text", Text: "You are a helpful assistant."},
				{Type: "text", Text: "Be concise."},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("failed to decode request: %v", err)
				}

				if req["system"] == nil {
					t.Error("expected system prompt in request")
				}

				stopReason := "end_turn"
				response := AnthropicMessagesResponse{
					ID:   "msg_sys",
					Type: "message",
					Role: "assistant",
					Content: []AnthropicResponseContentBlock{
						{Type: "text", Text: "OK"},
					},
					StopReason: &stopReason,
					Usage:      AnthropicUsage{InputTokens: 5, OutputTokens: 1},
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
			messages := []AnthropicMessage{CreateAnthropicUserMessage("Hi")}

			var opts []AnthropicOption
			opts = append(opts,
				WithAnthropicModel("anthropic/claude-sonnet-4"),
				WithAnthropicMaxTokens(100),
			)

			switch v := tc.system.(type) {
			case string:
				opts = append(opts, WithAnthropicSystemString(v))
			case []AnthropicTextBlock:
				opts = append(opts, WithAnthropicSystemBlocks(v))
			}

			resp, err := client.CreateAnthropicMessage(context.Background(), messages, opts...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.GetTextContent() != "OK" {
				t.Errorf("unexpected content: %q", resp.GetTextContent())
			}
		})
	}
}

func TestCreateAnthropicMessageWithTools(t *testing.T) {
	stopReason := "tool_use"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req AnthropicMessagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Tools) != 1 {
			t.Errorf("expected 1 tool, got %d", len(req.Tools))
		}
		if req.Tools[0].Name != "get_weather" {
			t.Errorf("expected tool name 'get_weather', got %q", req.Tools[0].Name)
		}

		response := AnthropicMessagesResponse{
			ID:   "msg_tools",
			Type: "message",
			Role: "assistant",
			Content: []AnthropicResponseContentBlock{
				{
					Type:  "tool_use",
					ID:    "toolu_123",
					Name:  "get_weather",
					Input: json.RawMessage(`{"location":"San Francisco"}`),
				},
			},
			StopReason: &stopReason,
			Usage:      AnthropicUsage{InputTokens: 20, OutputTokens: 15},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	tool := CreateAnthropicCustomTool("get_weather", "Get current weather", map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "City name",
			},
		},
		"required": []string{"location"},
	})

	messages := []AnthropicMessage{
		CreateAnthropicUserMessage("What's the weather in San Francisco?"),
	}

	resp, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(1024),
		WithAnthropicTools(tool),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp.IsToolUse() {
		t.Error("expected IsToolUse() to return true")
	}

	toolUses := resp.GetToolUseBlocks()
	if len(toolUses) != 1 {
		t.Fatalf("expected 1 tool use block, got %d", len(toolUses))
	}

	if toolUses[0].Name != "get_weather" {
		t.Errorf("expected tool name 'get_weather', got %q", toolUses[0].Name)
	}
}

func TestCreateAnthropicMessageWithThinking(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req AnthropicMessagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Thinking == nil {
			t.Error("expected thinking config in request")
		} else {
			if req.Thinking.Type != "enabled" {
				t.Errorf("expected thinking type 'enabled', got %q", req.Thinking.Type)
			}
			if req.Thinking.BudgetTokens != 5000 {
				t.Errorf("expected budget_tokens 5000, got %d", req.Thinking.BudgetTokens)
			}
		}

		stopReason := "end_turn"
		response := AnthropicMessagesResponse{
			ID:   "msg_think",
			Type: "message",
			Role: "assistant",
			Content: []AnthropicResponseContentBlock{
				{
					Type:     "thinking",
					Thinking: "Let me think about this...",
				},
				{
					Type: "text",
					Text: "The answer is 42.",
				},
			},
			StopReason: &stopReason,
			Usage:      AnthropicUsage{InputTokens: 10, OutputTokens: 50},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("What is the meaning of life?")}

	resp, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(8192),
		WithAnthropicThinkingEnabled(5000),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	thinking := resp.GetThinkingContent()
	if thinking != "Let me think about this..." {
		t.Errorf("unexpected thinking content: %q", thinking)
	}

	text := resp.GetTextContent()
	if text != "The answer is 42." {
		t.Errorf("unexpected text content: %q", text)
	}
}

func TestCreateAnthropicMessageStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" {
			t.Errorf("expected path /messages, got %s", r.URL.Path)
		}

		var req AnthropicMessagesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if !req.Stream {
			t.Error("expected stream to be true")
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		flusher, _ := w.(http.Flusher)

		// Send events
		events := []string{
			`data: {"type":"message_start","message":{"id":"msg_stream","type":"message","role":"assistant","content":[],"model":"anthropic/claude-sonnet-4","stop_reason":null,"usage":{"input_tokens":10,"output_tokens":0}}}`,
			`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
			`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}}`,
			`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":" world"}}`,
			`data: {"type":"content_block_stop","index":0}`,
			`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`,
			`data: {"type":"message_stop"}`,
			`data: [DONE]`,
		}

		for _, event := range events {
			_, _ = w.Write([]byte(event + "\n\n"))
			flusher.Flush()
		}
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hello!")}

	stream, err := client.CreateAnthropicMessageStream(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
	)

	if err != nil {
		t.Fatalf("unexpected error creating stream: %v", err)
	}
	defer func() { _ = stream.Close() }()

	var text strings.Builder
	eventCount := 0
	for event := range stream.Events() {
		eventCount++
		text.WriteString(event.GetTextDelta())
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if eventCount == 0 {
		t.Error("expected at least one event")
	}

	result := text.String()
	if result != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", result)
	}
}

func TestAnthropicValidationNoMessages(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.CreateAnthropicMessage(context.Background(), nil,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
	)

	if err != ErrNoAnthropicMessages {
		t.Errorf("expected ErrNoAnthropicMessages, got %v", err)
	}
}

func TestAnthropicValidationEmptyMessages(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.CreateAnthropicMessage(context.Background(), []AnthropicMessage{},
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
	)

	if err != ErrNoAnthropicMessages {
		t.Errorf("expected ErrNoAnthropicMessages, got %v", err)
	}
}

func TestAnthropicValidationInvalidRole(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	messages := []AnthropicMessage{
		{Role: "system", Content: "test"},
	}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
	)

	if err == nil {
		t.Error("expected error for invalid role")
	}
	valErr, ok := IsValidationError(err)
	if !ok {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if !strings.Contains(valErr.Message, "invalid role") {
		t.Errorf("expected 'invalid role' in message, got %q", valErr.Message)
	}
}

func TestAnthropicValidationNoMaxTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not have been sent")
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hello")}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
	)

	if err != ErrAnthropicMaxTokensRequired {
		t.Errorf("expected ErrAnthropicMaxTokensRequired, got %v", err)
	}
}

func TestAnthropicValidationInvalidThinkingType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not have been sent")
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hello")}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
		func(r *AnthropicMessagesRequest) {
			r.Thinking = &AnthropicThinkingConfig{Type: "invalid"}
		},
	)

	if err == nil {
		t.Error("expected error for invalid thinking type")
	}
}

func TestAnthropicValidationThinkingNoBudget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request should not have been sent")
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hello")}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
		func(r *AnthropicMessagesRequest) {
			r.Thinking = &AnthropicThinkingConfig{Type: "enabled", BudgetTokens: 0}
		},
	)

	if err == nil {
		t.Error("expected error for thinking enabled without budget")
	}
}

func TestAnthropicOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify tool_choice
		if tc, ok := req["tool_choice"].(map[string]interface{}); ok {
			if tc["type"] != "tool" {
				t.Errorf("expected tool_choice type 'tool', got %v", tc["type"])
			}
			if tc["name"] != "get_weather" {
				t.Errorf("expected tool_choice name 'get_weather', got %v", tc["name"])
			}
		} else {
			t.Error("expected tool_choice in request")
		}

		// Verify metadata
		if meta, ok := req["metadata"].(map[string]interface{}); ok {
			if meta["user_id"] != "user-123" {
				t.Errorf("expected user_id 'user-123', got %v", meta["user_id"])
			}
		} else {
			t.Error("expected metadata in request")
		}

		stopReason := "end_turn"
		response := AnthropicMessagesResponse{
			ID:   "msg_opts",
			Type: "message",
			Role: "assistant",
			Content: []AnthropicResponseContentBlock{
				{Type: "text", Text: "OK"},
			},
			StopReason: &stopReason,
			Usage:      AnthropicUsage{InputTokens: 5, OutputTokens: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hi")}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
		WithAnthropicToolChoiceSpecific("get_weather"),
		WithAnthropicRequestMetadata(AnthropicRequestMetadata{UserID: "user-123"}),
		WithAnthropicTools(CreateAnthropicCustomTool("get_weather", "Get weather", map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		})),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnthropicHeaderMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("expected X-Custom-Header 'custom-value', got %q", r.Header.Get("X-Custom-Header"))
		}

		stopReason := "end_turn"
		response := AnthropicMessagesResponse{
			ID:   "msg_hdr",
			Type: "message",
			Role: "assistant",
			Content: []AnthropicResponseContentBlock{
				{Type: "text", Text: "OK"},
			},
			StopReason: &stopReason,
			Usage:      AnthropicUsage{InputTokens: 5, OutputTokens: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []AnthropicMessage{CreateAnthropicUserMessage("Hi")}

	_, err := client.CreateAnthropicMessage(context.Background(), messages,
		WithAnthropicModel("anthropic/claude-sonnet-4"),
		WithAnthropicMaxTokens(100),
		WithAnthropicHeaderMetadata(map[string]interface{}{
			"Custom-Header": "custom-value",
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAnthropicHelperConstructors(t *testing.T) {
	t.Run("user message", func(t *testing.T) {
		msg := CreateAnthropicUserMessage("hello")
		if msg.Role != "user" {
			t.Errorf("expected role 'user', got %q", msg.Role)
		}
		if msg.Content != "hello" {
			t.Errorf("expected content 'hello', got %v", msg.Content)
		}
	})

	t.Run("assistant message", func(t *testing.T) {
		msg := CreateAnthropicAssistantMessage("hi")
		if msg.Role != "assistant" {
			t.Errorf("expected role 'assistant', got %q", msg.Role)
		}
	})

	t.Run("user message with blocks", func(t *testing.T) {
		blocks := []AnthropicContentBlock{
			CreateAnthropicTextBlock("hello"),
			CreateAnthropicImageURLBlock("https://example.com/image.png"),
		}
		msg := CreateAnthropicUserMessageWithBlocks(blocks)
		if msg.Role != "user" {
			t.Errorf("expected role 'user', got %q", msg.Role)
		}
		contentBlocks, ok := msg.Content.([]AnthropicContentBlock)
		if !ok {
			t.Fatal("expected content to be []AnthropicContentBlock")
		}
		if len(contentBlocks) != 2 {
			t.Errorf("expected 2 content blocks, got %d", len(contentBlocks))
		}
	})

	t.Run("text block", func(t *testing.T) {
		block := CreateAnthropicTextBlock("test")
		if block.Type != "text" || block.Text != "test" {
			t.Errorf("unexpected text block: %+v", block)
		}
	})

	t.Run("image URL block", func(t *testing.T) {
		block := CreateAnthropicImageURLBlock("https://example.com/img.png")
		if block.Type != "image" || block.Source == nil || block.Source.Type != "url" {
			t.Errorf("unexpected image URL block: %+v", block)
		}
	})

	t.Run("image base64 block", func(t *testing.T) {
		block := CreateAnthropicImageBase64Block("image/png", "base64data")
		if block.Type != "image" || block.Source == nil || block.Source.Type != "base64" {
			t.Errorf("unexpected image base64 block: %+v", block)
		}
		if block.Source.MediaType != "image/png" || block.Source.Data != "base64data" {
			t.Errorf("unexpected source fields: %+v", block.Source)
		}
	})

	t.Run("document URL block", func(t *testing.T) {
		block := CreateAnthropicDocumentURLBlock("https://example.com/doc.pdf")
		if block.Type != "document" || block.Source == nil || block.Source.Type != "url" {
			t.Errorf("unexpected document URL block: %+v", block)
		}
	})

	t.Run("custom tool", func(t *testing.T) {
		tool := CreateAnthropicCustomTool("test_tool", "A test tool", map[string]interface{}{"type": "object"})
		if tool.Name != "test_tool" || tool.Description != "A test tool" {
			t.Errorf("unexpected tool: %+v", tool)
		}
	})

	t.Run("bash tool", func(t *testing.T) {
		tool := CreateAnthropicBashTool()
		if tool.Type != "bash_20250124" || tool.Name != "bash" {
			t.Errorf("unexpected bash tool: %+v", tool)
		}
	})

	t.Run("text editor tool", func(t *testing.T) {
		tool := CreateAnthropicTextEditorTool()
		if tool.Type != "text_editor_20250124" || tool.Name != "str_replace_editor" {
			t.Errorf("unexpected text editor tool: %+v", tool)
		}
	})

	t.Run("web search tool", func(t *testing.T) {
		maxUses := 5
		tool := CreateAnthropicWebSearchTool(&maxUses)
		if tool.Type != "web_search_20250305" || tool.Name != "web_search" {
			t.Errorf("unexpected web search tool: %+v", tool)
		}
		if tool.MaxUses == nil || *tool.MaxUses != 5 {
			t.Errorf("expected max_uses 5, got %v", tool.MaxUses)
		}
	})

	t.Run("tool result block", func(t *testing.T) {
		block := CreateAnthropicToolResultBlock("toolu_123", "The weather is sunny")
		if block.Type != "tool_result" || block.ToolUseID != "toolu_123" {
			t.Errorf("unexpected tool result block: %+v", block)
		}
	})
}

func TestAnthropicResponseHelpers(t *testing.T) {
	stopReason := "end_turn"

	t.Run("GetTextBlocks", func(t *testing.T) {
		resp := &AnthropicMessagesResponse{
			Content: []AnthropicResponseContentBlock{
				{Type: "thinking", Thinking: "hmm"},
				{Type: "text", Text: "first"},
				{Type: "text", Text: "second"},
			},
			StopReason: &stopReason,
		}
		blocks := resp.GetTextBlocks()
		if len(blocks) != 2 {
			t.Errorf("expected 2 text blocks, got %d", len(blocks))
		}
	})

	t.Run("GetThinkingContent", func(t *testing.T) {
		resp := &AnthropicMessagesResponse{
			Content: []AnthropicResponseContentBlock{
				{Type: "thinking", Thinking: "part1"},
				{Type: "text", Text: "answer"},
				{Type: "thinking", Thinking: "part2"},
			},
			StopReason: &stopReason,
		}
		thinking := resp.GetThinkingContent()
		if thinking != "part1part2" {
			t.Errorf("expected 'part1part2', got %q", thinking)
		}
	})

	t.Run("nil stop reason", func(t *testing.T) {
		resp := &AnthropicMessagesResponse{}
		if resp.GetStopReason() != "" {
			t.Errorf("expected empty stop reason, got %q", resp.GetStopReason())
		}
		if resp.IsToolUse() {
			t.Error("expected IsToolUse() to be false for nil stop reason")
		}
	})
}
