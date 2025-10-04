package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderRoutingOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Check that provider options are set correctly
		if req.Provider == nil {
			t.Error("Provider should not be nil")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	messages := []Message{CreateUserMessage("Test")}

	tests := []struct {
		name   string
		option ChatCompletionOption
		verify func(t *testing.T, req ChatCompletionRequest)
	}{
		{
			name:   "WithZDR",
			option: WithZDR(true),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.ZDR == nil || !*req.Provider.ZDR {
					t.Error("ZDR not set to true")
				}
			},
		},
		{
			name:   "WithProviderOrder",
			option: WithProviderOrder("openai", "anthropic"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || len(req.Provider.Order) != 2 {
					t.Error("Provider order not set correctly")
				}
			},
		},
		{
			name:   "WithAllowFallbacks",
			option: WithAllowFallbacks(false),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.AllowFallbacks == nil || *req.Provider.AllowFallbacks {
					t.Error("AllowFallbacks not set to false")
				}
			},
		},
		{
			name:   "WithRequireParameters",
			option: WithRequireParameters(true),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.RequireParameters == nil || !*req.Provider.RequireParameters {
					t.Error("RequireParameters not set to true")
				}
			},
		},
		{
			name:   "WithDataCollection",
			option: WithDataCollection("deny"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.DataCollection != "deny" {
					t.Error("DataCollection not set to deny")
				}
			},
		},
		{
			name:   "WithOnlyProviders",
			option: WithOnlyProviders("azure", "openai"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || len(req.Provider.Only) != 2 {
					t.Error("Only providers not set correctly")
				}
			},
		},
		{
			name:   "WithIgnoreProviders",
			option: WithIgnoreProviders("cohere"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || len(req.Provider.Ignore) != 1 || req.Provider.Ignore[0] != "cohere" {
					t.Error("Ignore providers not set correctly")
				}
			},
		},
		{
			name:   "WithQuantizations",
			option: WithQuantizations("fp8", "fp16"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || len(req.Provider.Quantizations) != 2 {
					t.Error("Quantizations not set correctly")
				}
			},
		},
		{
			name:   "WithProviderSort",
			option: WithProviderSort("throughput"),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "throughput" {
					t.Error("Provider sort not set correctly")
				}
			},
		},
		{
			name:   "WithMaxPrice",
			option: WithMaxPrice(MaxPrice{Prompt: 1.0, Completion: 2.0}),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.MaxPrice == nil {
					t.Error("MaxPrice not set")
				}
				if req.Provider.MaxPrice.Prompt != 1.0 || req.Provider.MaxPrice.Completion != 2.0 {
					t.Error("MaxPrice values not set correctly")
				}
			},
		},
		{
			name:   "WithNitro",
			option: WithNitro(),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "throughput" {
					t.Error("Nitro shortcut not working")
				}
			},
		},
		{
			name:   "WithFloorPrice",
			option: WithFloorPrice(),
			verify: func(t *testing.T, req ChatCompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "price" {
					t.Error("FloorPrice shortcut not working")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new test server for each test to capture request
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req ChatCompletionRequest
				json.NewDecoder(r.Body).Decode(&req)
				tt.verify(t, req)
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
			}))
			defer testServer.Close()

			testClient := NewClient(WithAPIKey("test-key"), WithBaseURL(testServer.URL))
			_, err := testClient.ChatComplete(context.Background(), messages,
				WithModel("test-model"),
				tt.option,
			)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCompletionProviderOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	tests := []struct {
		name   string
		option CompletionOption
		verify func(t *testing.T, req CompletionRequest)
	}{
		{
			name:   "WithCompletionZDR",
			option: WithCompletionZDR(true),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.ZDR == nil || !*req.Provider.ZDR {
					t.Error("ZDR not set to true")
				}
			},
		},
		{
			name:   "WithCompletionProviderOrder",
			option: WithCompletionProviderOrder("openai", "anthropic"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || len(req.Provider.Order) != 2 {
					t.Error("Provider order not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionAllowFallbacks",
			option: WithCompletionAllowFallbacks(false),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.AllowFallbacks == nil || *req.Provider.AllowFallbacks {
					t.Error("AllowFallbacks not set to false")
				}
			},
		},
		{
			name:   "WithCompletionRequireParameters",
			option: WithCompletionRequireParameters(true),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.RequireParameters == nil || !*req.Provider.RequireParameters {
					t.Error("RequireParameters not set to true")
				}
			},
		},
		{
			name:   "WithCompletionDataCollection",
			option: WithCompletionDataCollection("allow"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.DataCollection != "allow" {
					t.Error("DataCollection not set to allow")
				}
			},
		},
		{
			name:   "WithCompletionOnlyProviders",
			option: WithCompletionOnlyProviders("azure"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || len(req.Provider.Only) != 1 || req.Provider.Only[0] != "azure" {
					t.Error("Only providers not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionIgnoreProviders",
			option: WithCompletionIgnoreProviders("deepinfra", "together"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || len(req.Provider.Ignore) != 2 {
					t.Error("Ignore providers not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionQuantizations",
			option: WithCompletionQuantizations("int4", "int8"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || len(req.Provider.Quantizations) != 2 {
					t.Error("Quantizations not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionProviderSort",
			option: WithCompletionProviderSort("latency"),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "latency" {
					t.Error("Provider sort not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionMaxPrice",
			option: WithCompletionMaxPrice(MaxPrice{Prompt: 0.5, Completion: 1.0, Request: 0.01}),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.MaxPrice == nil {
					t.Error("MaxPrice not set")
				}
				if req.Provider.MaxPrice.Prompt != 0.5 || req.Provider.MaxPrice.Completion != 1.0 || req.Provider.MaxPrice.Request != 0.01 {
					t.Error("MaxPrice values not set correctly")
				}
			},
		},
		{
			name:   "WithCompletionNitro",
			option: WithCompletionNitro(),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "throughput" {
					t.Error("Nitro shortcut not working")
				}
			},
		},
		{
			name:   "WithCompletionFloorPrice",
			option: WithCompletionFloorPrice(),
			verify: func(t *testing.T, req CompletionRequest) {
				if req.Provider == nil || req.Provider.Sort != "price" {
					t.Error("FloorPrice shortcut not working")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new test server for each test to capture request
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req CompletionRequest
				json.NewDecoder(r.Body).Decode(&req)
				tt.verify(t, req)
				w.WriteHeader(200)
				json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
			}))
			defer testServer.Close()

			testClient := NewClient(WithAPIKey("test-key"), WithBaseURL(testServer.URL))
			_, err := testClient.Complete(context.Background(), "test prompt",
				WithCompletionModel("test-model"),
				tt.option,
			)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestAppAttributionOptions(t *testing.T) {
	// Test that app attribution headers are properly set
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check HTTP-Referer header
		if referer := r.Header.Get("HTTP-Referer"); referer != "https://myapp.com" {
			t.Errorf("HTTP-Referer not set correctly, got: %s", referer)
		}

		// Check X-Title header
		if title := r.Header.Get("X-Title"); title != "My AI Assistant" {
			t.Errorf("X-Title not set correctly, got: %s", title)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithReferer("https://myapp.com"),
		WithAppName("My AI Assistant"),
	)

	messages := []Message{CreateUserMessage("Test")}
	_, err := client.ChatComplete(context.Background(), messages, WithModel("test-model"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAppAttributionOptionsCompletion(t *testing.T) {
	// Test that app attribution headers work with completion endpoints
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check HTTP-Referer header
		if referer := r.Header.Get("HTTP-Referer"); referer != "https://localhost:3000" {
			t.Errorf("HTTP-Referer not set correctly, got: %s", referer)
		}

		// Check X-Title header (especially important for localhost)
		if title := r.Header.Get("X-Title"); title != "Development App" {
			t.Errorf("X-Title not set correctly, got: %s", title)
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithReferer("https://localhost:3000"),
		WithAppName("Development App"),
	)

	_, err := client.Complete(context.Background(), "test prompt", WithCompletionModel("test-model"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMultipleProviderOptions(t *testing.T) {
	// Test that multiple provider options work together correctly
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify multiple options are applied
		if req.Provider == nil {
			t.Error("Provider should not be nil")
			return
		}

		// Check ZDR
		if req.Provider.ZDR == nil || !*req.Provider.ZDR {
			t.Error("ZDR not set")
		}

		// Check Order
		if len(req.Provider.Order) != 2 {
			t.Error("Provider order not set")
		}

		// Check AllowFallbacks
		if req.Provider.AllowFallbacks == nil || *req.Provider.AllowFallbacks {
			t.Error("AllowFallbacks should be false")
		}

		// Check Sort
		if req.Provider.Sort != "throughput" {
			t.Error("Sort not set to throughput")
		}

		// Check MaxPrice
		if req.Provider.MaxPrice == nil || req.Provider.MaxPrice.Prompt != 2.5 {
			t.Error("MaxPrice not set correctly")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []Message{CreateUserMessage("Test")}

	_, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithZDR(true),
		WithProviderOrder("openai", "anthropic"),
		WithAllowFallbacks(false),
		WithProviderSort("throughput"),
		WithMaxPrice(MaxPrice{Prompt: 2.5, Completion: 5.0}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func TestUntestedOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify WithTools
		if req.Tools == nil || len(req.Tools) != 1 {
			t.Error("Tools not set correctly")
		}

		// Verify WithToolChoice
		if req.ToolChoice == nil {
			t.Error("ToolChoice not set")
		}

		// Verify WithParallelToolCalls
		if req.ParallelToolCalls == nil || *req.ParallelToolCalls != false {
			t.Error("ParallelToolCalls not set correctly")
		}

		// Verify WithMessages (should override the messages passed to ChatComplete)
		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages from WithMessages, got %d", len(req.Messages))
		}

		// Verify WithResponseFormat (JSON schema)
		if req.ResponseFormat == nil || req.ResponseFormat.Type != "json_schema" {
			t.Error("ResponseFormat not set correctly")
		}

		// Verify WithTransforms
		if req.Transforms == nil || len(req.Transforms) != 1 || req.Transforms[0] != "middle-out" {
			t.Error("Transforms not set correctly")
		}

		// Verify WithModels
		if req.Models == nil || len(req.Models) != 2 {
			t.Error("Models not set correctly")
		}

		// Verify WithRoute
		if req.Route != "fallback" {
			t.Error("Route not set correctly")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	tools := []Tool{
		{
			Type: "function",
			Function: Function{
				Name: "test_function",
			},
		},
	}

	messages := []Message{CreateUserMessage("Original")}
	overrideMessages := []Message{
		CreateSystemMessage("System"),
		CreateUserMessage("Override"),
	}

	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"field": map[string]interface{}{"type": "string"},
		},
	}

	parallelCalls := false

	_, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithTools(tools...),
		WithToolChoice("auto"),
		WithParallelToolCalls(&parallelCalls),
		WithMessages(overrideMessages),
		WithJSONSchema("test_schema", true, schema),
		WithTransforms("middle-out"),
		WithModels("model1", "model2"),
		WithRoute("fallback"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompletionUntestedOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify WithCompletionProvider
		if req.Provider == nil {
			t.Error("Provider not set")
		}

		// Verify WithCompletionResponseFormat (JSON schema)
		if req.ResponseFormat == nil {
			t.Error("ResponseFormat not set")
		}

		// Verify WithCompletionTransforms
		if req.Transforms == nil || len(req.Transforms) != 1 {
			t.Error("Transforms not set correctly")
		}

		// Verify WithCompletionModels
		if req.Models == nil || len(req.Models) != 2 {
			t.Error("Models not set correctly")
		}

		// Verify WithCompletionRoute
		if req.Route != "fallback" {
			t.Error("Route not set correctly")
		}

		// Verify WithCompletionPlugins
		if req.Plugins == nil || len(req.Plugins) != 1 {
			t.Error("Plugins not set correctly")
		}

		// Verify WithCompletionWebSearchOptions
		if req.WebSearchOptions == nil {
			t.Error("WebSearchOptions not set")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	provider := Provider{
		Order: []string{"openai"},
	}

	schema := map[string]interface{}{
		"type": "object",
	}

	_, err := client.Complete(context.Background(), "Test prompt",
		WithCompletionModel("test-model"),
		WithCompletionProvider(provider),
		WithCompletionJSONSchema("test_schema", true, schema),
		WithCompletionTransforms("middle-out"),
		WithCompletionModels("model1", "model2"),
		WithCompletionRoute("fallback"),
		WithCompletionPlugins(Plugin{ID: "web"}),
		WithCompletionWebSearchOptions(&WebSearchOptions{SearchContextSize: "medium"}),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestJSONModeOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify WithJSONMode sets type to json_object
		if req.ResponseFormat == nil || req.ResponseFormat.Type != "json_object" {
			t.Error("JSON mode not set correctly")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []Message{CreateUserMessage("Test")}

	_, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithJSONMode(),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompletionJSONModeOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify WithCompletionJSONMode sets type to json_object
		if req.ResponseFormat == nil || req.ResponseFormat.Type != "json_object" {
			t.Error("JSON mode not set correctly")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.Complete(context.Background(), "Test",
		WithCompletionModel("test-model"),
		WithCompletionJSONMode(),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTopKAndRepetitionPenalty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		json.NewDecoder(r.Body).Decode(&req)

		// Verify WithTopK
		if req.TopK == nil || *req.TopK != 40 {
			t.Error("TopK not set correctly")
		}

		// Verify WithRepetitionPenalty
		if req.RepetitionPenalty == nil || *req.RepetitionPenalty != 1.1 {
			t.Error("RepetitionPenalty not set correctly")
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	messages := []Message{CreateUserMessage("Test")}

	_, err := client.ChatComplete(context.Background(), messages,
		WithModel("test-model"),
		WithTopK(40),
		WithRepetitionPenalty(1.1),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
