package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewWebPlugin(t *testing.T) {
	plugin := NewWebPlugin()

	if plugin.ID != "web" {
		t.Errorf("expected plugin ID 'web', got %q", plugin.ID)
	}

	if plugin.MaxResults != 5 {
		t.Errorf("expected MaxResults 5, got %d", plugin.MaxResults)
	}

	if plugin.Engine != "" {
		t.Errorf("expected empty Engine (auto), got %q", plugin.Engine)
	}
}

func TestNewWebPluginWithOptions(t *testing.T) {
	searchPrompt := "Custom search prompt"
	plugin := NewWebPluginWithOptions(WebSearchEngineExa, 10, searchPrompt)

	if plugin.ID != "web" {
		t.Errorf("expected plugin ID 'web', got %q", plugin.ID)
	}

	if plugin.Engine != string(WebSearchEngineExa) {
		t.Errorf("expected engine %q, got %q", WebSearchEngineExa, plugin.Engine)
	}

	if plugin.MaxResults != 10 {
		t.Errorf("expected MaxResults 10, got %d", plugin.MaxResults)
	}

	if plugin.SearchPrompt != searchPrompt {
		t.Errorf("expected SearchPrompt %q, got %q", searchPrompt, plugin.SearchPrompt)
	}
}

func TestWithOnlineModel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"openai/gpt-4o", "openai/gpt-4o:online"},
		{"anthropic/claude-3-sonnet", "anthropic/claude-3-sonnet:online"},
		{"meta-llama/llama-3.1-8b-instruct", "meta-llama/llama-3.1-8b-instruct:online"},
	}

	for _, tt := range tests {
		result := WithOnlineModel(tt.input)
		if result != tt.expected {
			t.Errorf("WithOnlineModel(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDefaultWebSearchPrompt(t *testing.T) {
	date := "2024-01-15"
	result := DefaultWebSearchPrompt(date)
	expected := "A web search was conducted on `2024-01-15`. Incorporate the following web search results into your response.\n\n" +
		"IMPORTANT: Cite them using markdown links named using the domain of the source.\n" +
		"Example: [nytimes.com](https://nytimes.com/some-page)."

	if result != expected {
		t.Errorf("DefaultWebSearchPrompt() returned unexpected result.\nGot: %q\nWant: %q", result, expected)
	}
}

func TestParseAnnotations(t *testing.T) {
	annotations := []Annotation{
		{
			Type: "url_citation",
			URLCitation: &URLCitation{
				URL:        "https://example.com/article",
				Title:      "Example Article",
				Content:    "Article content excerpt",
				StartIndex: 10,
				EndIndex:   50,
			},
		},
		{
			Type: "url_citation",
			URLCitation: &URLCitation{
				URL:        "https://test.com/page",
				Title:      "Test Page",
				Content:    "Test content",
				StartIndex: 60,
				EndIndex:   100,
			},
		},
		{
			Type:        "other_type",
			URLCitation: nil,
		},
	}

	citations := ParseAnnotations(annotations)

	if len(citations) != 2 {
		t.Fatalf("expected 2 citations, got %d", len(citations))
	}

	if citations[0].URL != "https://example.com/article" {
		t.Errorf("expected first citation URL 'https://example.com/article', got %q", citations[0].URL)
	}

	if citations[1].Title != "Test Page" {
		t.Errorf("expected second citation title 'Test Page', got %q", citations[1].Title)
	}
}

func TestWithPlugins(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		if len(req.Plugins) != 2 {
			t.Errorf("expected 2 plugins, got %d", len(req.Plugins))
		}

		if req.Plugins[0].ID != "web" {
			t.Errorf("expected first plugin ID 'web', got %q", req.Plugins[0].ID)
		}

		if req.Plugins[1].ID != "custom" {
			t.Errorf("expected second plugin ID 'custom', got %q", req.Plugins[1].ID)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatCompletionResponse{
			ID:      "test-123",
			Choices: []Choice{{Message: Message{Content: "Response with plugins"}}},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	plugin1 := NewWebPlugin()
	plugin2 := Plugin{ID: "custom", MaxResults: 3}

	messages := []Message{CreateUserMessage("Test with plugins")}

	resp, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithPlugins(plugin1, plugin2),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-123" {
		t.Errorf("expected response ID 'test-123', got %q", resp.ID)
	}
}

func TestWithWebSearchOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.WebSearchOptions == nil {
			t.Error("expected WebSearchOptions to be set")
		} else if req.WebSearchOptions.SearchContextSize != "high" {
			t.Errorf("expected SearchContextSize 'high', got %q", req.WebSearchOptions.SearchContextSize)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatCompletionResponse{
			ID:      "test-456",
			Choices: []Choice{{Message: Message{Content: "Response with search options"}}},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	messages := []Message{CreateUserMessage("Test with search options")}

	resp, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithWebSearchOptions(&WebSearchOptions{
			SearchContextSize: string(WebSearchContextHigh),
		}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-456" {
		t.Errorf("expected response ID 'test-456', got %q", resp.ID)
	}
}

func TestChatCompleteWithOnlineSuffix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		if req.Model != "openai/gpt-4o:online" {
			t.Errorf("expected model 'openai/gpt-4o:online', got %q", req.Model)
		}

		// Send response with annotations
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatCompletionResponse{
			ID: "test-789",
			Choices: []Choice{{
				Message: Message{
					Content: "Response with web search",
					Annotations: []Annotation{
						{
							Type: "url_citation",
							URLCitation: &URLCitation{
								URL:        "https://source.com/article",
								Title:      "Source Article",
								Content:    "Article excerpt",
								StartIndex: 0,
								EndIndex:   20,
							},
						},
					},
				},
			}},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	messages := []Message{CreateUserMessage("Search for information")}

	resp, err := client.ChatComplete(context.Background(), messages,
		WithModel("openai/gpt-4o:online"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Choices[0].Message.Annotations) != 1 {
		t.Errorf("expected 1 annotation, got %d", len(resp.Choices[0].Message.Annotations))
	}

	citations := ParseAnnotations(resp.Choices[0].Message.Annotations)
	if len(citations) != 1 || citations[0].URL != "https://source.com/article" {
		t.Errorf("expected citation URL 'https://source.com/article', got %+v", citations)
	}
}

func TestWebSearchEngineConstants(t *testing.T) {
	if string(WebSearchEngineNative) != "native" {
		t.Errorf("expected WebSearchEngineNative to be 'native', got %q", WebSearchEngineNative)
	}

	if string(WebSearchEngineExa) != "exa" {
		t.Errorf("expected WebSearchEngineExa to be 'exa', got %q", WebSearchEngineExa)
	}

	if string(WebSearchEngineAuto) != "" {
		t.Errorf("expected WebSearchEngineAuto to be empty string, got %q", WebSearchEngineAuto)
	}
}

func TestWebSearchContextSizeConstants(t *testing.T) {
	if string(WebSearchContextLow) != "low" {
		t.Errorf("expected WebSearchContextLow to be 'low', got %q", WebSearchContextLow)
	}

	if string(WebSearchContextMedium) != "medium" {
		t.Errorf("expected WebSearchContextMedium to be 'medium', got %q", WebSearchContextMedium)
	}

	if string(WebSearchContextHigh) != "high" {
		t.Errorf("expected WebSearchContextHigh to be 'high', got %q", WebSearchContextHigh)
	}
}

func TestPluginSerialization(t *testing.T) {
	plugin := Plugin{
		ID:           "web",
		Engine:       "exa",
		MaxResults:   7,
		SearchPrompt: "Custom prompt",
	}

	data, err := json.Marshal(plugin)
	if err != nil {
		t.Fatalf("failed to marshal plugin: %v", err)
	}

	var decoded Plugin
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal plugin: %v", err)
	}

	if !reflect.DeepEqual(plugin, decoded) {
		t.Errorf("plugin serialization mismatch.\nOriginal: %+v\nDecoded: %+v", plugin, decoded)
	}
}

func TestWebSearchOptionsSerialization(t *testing.T) {
	options := WebSearchOptions{
		SearchContextSize: "medium",
	}

	data, err := json.Marshal(options)
	if err != nil {
		t.Fatalf("failed to marshal WebSearchOptions: %v", err)
	}

	var decoded WebSearchOptions
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal WebSearchOptions: %v", err)
	}

	if options.SearchContextSize != decoded.SearchContextSize {
		t.Errorf("WebSearchOptions serialization mismatch.\nOriginal: %+v\nDecoded: %+v", options, decoded)
	}
}