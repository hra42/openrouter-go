package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEmbedding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/embeddings" {
			t.Errorf("expected path /embeddings, got %s", r.URL.Path)
		}

		// Verify request body
		var reqBody EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if reqBody.Model != "openai/text-embedding-3-small" {
			t.Errorf("expected model 'openai/text-embedding-3-small', got %q", reqBody.Model)
		}

		// Send response
		response := EmbeddingResponse{
			ID:     "emb-123",
			Object: "list",
			Data: []EmbeddingData{
				{
					Object:    "embedding",
					Embedding: []interface{}{0.1, 0.2, 0.3, 0.4, 0.5},
					Index:     0,
				},
			},
			Model: "openai/text-embedding-3-small",
			Usage: &EmbeddingUsage{
				PromptTokens: 5,
				TotalTokens:  5,
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

	resp, err := client.CreateEmbedding(context.Background(), "Hello, world!", "openai/text-embedding-3-small")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Object != "list" {
		t.Errorf("expected object 'list', got %q", resp.Object)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 embedding, got %d", len(resp.Data))
	}
	if resp.Data[0].Object != "embedding" {
		t.Errorf("expected embedding object 'embedding', got %q", resp.Data[0].Object)
	}
	if resp.Data[0].Index != 0 {
		t.Errorf("expected embedding index 0, got %d", resp.Data[0].Index)
	}
	if resp.Usage == nil {
		t.Error("expected usage to be set")
	}
	if resp.Usage.PromptTokens != 5 {
		t.Errorf("expected 5 prompt tokens, got %d", resp.Usage.PromptTokens)
	}
}

func TestCreateEmbeddings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body
		var reqBody EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Verify input is an array
		inputs, ok := reqBody.Input.([]interface{})
		if !ok {
			t.Errorf("expected input to be an array, got %T", reqBody.Input)
		}
		if len(inputs) != 3 {
			t.Errorf("expected 3 inputs, got %d", len(inputs))
		}

		// Send response with multiple embeddings
		response := EmbeddingResponse{
			ID:     "emb-456",
			Object: "list",
			Data: []EmbeddingData{
				{Object: "embedding", Embedding: []interface{}{0.1, 0.2}, Index: 0},
				{Object: "embedding", Embedding: []interface{}{0.3, 0.4}, Index: 1},
				{Object: "embedding", Embedding: []interface{}{0.5, 0.6}, Index: 2},
			},
			Model: "openai/text-embedding-3-small",
			Usage: &EmbeddingUsage{
				PromptTokens: 15,
				TotalTokens:  15,
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

	inputs := []string{"Hello", "World", "Test"}
	resp, err := client.CreateEmbeddings(context.Background(), inputs, "openai/text-embedding-3-small")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 3 {
		t.Errorf("expected 3 embeddings, got %d", len(resp.Data))
	}
}

func TestCreateEmbeddingWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Verify options were applied
		if reqBody.EncodingFormat != "base64" {
			t.Errorf("expected encoding_format 'base64', got %q", reqBody.EncodingFormat)
		}
		if reqBody.Dimensions == nil || *reqBody.Dimensions != 256 {
			t.Errorf("expected dimensions 256, got %v", reqBody.Dimensions)
		}
		if reqBody.User != "test-user" {
			t.Errorf("expected user 'test-user', got %q", reqBody.User)
		}
		if reqBody.InputType != "query" {
			t.Errorf("expected input_type 'query', got %q", reqBody.InputType)
		}

		// Send response
		response := EmbeddingResponse{
			Object: "list",
			Data: []EmbeddingData{
				{
					Object:    "embedding",
					Embedding: "YmFzZTY0ZW5jb2RlZA==",
					Index:     0,
				},
			},
			Model: "openai/text-embedding-3-small",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.CreateEmbedding(
		context.Background(),
		"test query",
		"openai/text-embedding-3-small",
		WithEmbeddingEncodingFormat("base64"),
		WithEmbeddingDimensions(256),
		WithEmbeddingUser("test-user"),
		WithEmbeddingInputType("query"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify base64 response
	base64Str := resp.Data[0].GetEmbeddingBase64()
	if base64Str != "YmFzZTY0ZW5jb2RlZA==" {
		t.Errorf("expected base64 string, got %q", base64Str)
	}
}

func TestCreateEmbeddingWithProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody EmbeddingRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Verify provider options
		if reqBody.Provider == nil {
			t.Error("expected provider to be set")
		} else {
			if len(reqBody.Provider.Only) != 1 || reqBody.Provider.Only[0] != "openai" {
				t.Errorf("expected only=['openai'], got %v", reqBody.Provider.Only)
			}
			if reqBody.Provider.Sort != "price" {
				t.Errorf("expected sort 'price', got %q", reqBody.Provider.Sort)
			}
		}

		response := EmbeddingResponse{
			Object: "list",
			Data: []EmbeddingData{
				{Object: "embedding", Embedding: []interface{}{0.1}, Index: 0},
			},
			Model: "openai/text-embedding-3-small",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.CreateEmbedding(
		context.Background(),
		"test",
		"openai/text-embedding-3-small",
		WithEmbeddingOnlyProviders("openai"),
		WithEmbeddingProviderSort("price"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateEmbeddingValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Test missing input
	_, err := client.CreateEmbedding(context.Background(), nil, "openai/text-embedding-3-small")
	if err == nil {
		t.Error("expected error for nil input")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "input" {
		t.Errorf("expected validation error for input field, got %v", err)
	}

	// Test empty string input
	_, err = client.CreateEmbedding(context.Background(), "", "openai/text-embedding-3-small")
	if err == nil {
		t.Error("expected error for empty string input")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "input" {
		t.Errorf("expected validation error for input field, got %v", err)
	}

	// Test empty slice input
	_, err = client.CreateEmbedding(context.Background(), []string{}, "openai/text-embedding-3-small")
	if err == nil {
		t.Error("expected error for empty slice input")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "input" {
		t.Errorf("expected validation error for input field, got %v", err)
	}

	// Test missing model
	_, err = client.CreateEmbedding(context.Background(), "test", "")
	if err == nil {
		t.Error("expected error for empty model")
	}
	if err != ErrNoModel {
		t.Errorf("expected ErrNoModel, got %v", err)
	}

	// Test missing API key
	clientNoKey := NewClient()
	_, err = clientNoKey.CreateEmbedding(context.Background(), "test", "model")
	if err == nil {
		t.Error("expected error for missing API key")
	}
	if err != ErrNoAPIKey {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestListEmbeddingsModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/embeddings/models" {
			t.Errorf("expected path /embeddings/models, got %s", r.URL.Path)
		}

		// Send response
		contextLength := 8191.0
		response := EmbeddingsModelsResponse{
			Data: []Model{
				{
					ID:            "openai/text-embedding-3-small",
					Name:          "Text Embedding 3 Small",
					Description:   "OpenAI's text embedding model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"embeddings"},
						Tokenizer:        "cl100k_base",
					},
					Pricing: ModelPricing{
						Prompt:     "0.00002",
						Completion: "0",
					},
				},
				{
					ID:            "openai/text-embedding-3-large",
					Name:          "Text Embedding 3 Large",
					Description:   "OpenAI's large text embedding model",
					ContextLength: &contextLength,
					Architecture: ModelArchitecture{
						InputModalities:  []string{"text"},
						OutputModalities: []string{"embeddings"},
						Tokenizer:        "cl100k_base",
					},
					Pricing: ModelPricing{
						Prompt:     "0.00013",
						Completion: "0",
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

	resp, err := client.ListEmbeddingsModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 models, got %d", len(resp.Data))
	}

	model := resp.Data[0]
	if model.ID != "openai/text-embedding-3-small" {
		t.Errorf("expected model ID 'openai/text-embedding-3-small', got %q", model.ID)
	}
	if model.Name != "Text Embedding 3 Small" {
		t.Errorf("expected model name 'Text Embedding 3 Small', got %q", model.Name)
	}
	if len(model.Architecture.OutputModalities) != 1 || model.Architecture.OutputModalities[0] != "embeddings" {
		t.Errorf("expected output modality 'embeddings', got %v", model.Architecture.OutputModalities)
	}
}

func TestGetEmbeddingVector(t *testing.T) {
	// Test with []float64
	data := &EmbeddingData{
		Embedding: []float64{0.1, 0.2, 0.3},
	}
	vec := data.GetEmbeddingVector()
	if len(vec) != 3 || vec[0] != 0.1 {
		t.Errorf("expected [0.1, 0.2, 0.3], got %v", vec)
	}

	// Test with []interface{}
	data2 := &EmbeddingData{
		Embedding: []interface{}{0.4, 0.5, 0.6},
	}
	vec2 := data2.GetEmbeddingVector()
	if len(vec2) != 3 || vec2[0] != 0.4 {
		t.Errorf("expected [0.4, 0.5, 0.6], got %v", vec2)
	}

	// Test with nil
	data3 := &EmbeddingData{
		Embedding: nil,
	}
	vec3 := data3.GetEmbeddingVector()
	if vec3 != nil {
		t.Errorf("expected nil, got %v", vec3)
	}

	// Test with string (base64)
	data4 := &EmbeddingData{
		Embedding: "base64string",
	}
	vec4 := data4.GetEmbeddingVector()
	if vec4 != nil {
		t.Errorf("expected nil for base64 string, got %v", vec4)
	}

	// Test with mixed types in []interface{}
	data5 := &EmbeddingData{
		Embedding: []interface{}{0.1, "string", 0.3},
	}
	vec5 := data5.GetEmbeddingVector()
	if vec5 != nil {
		t.Errorf("expected nil for mixed type array, got %v", vec5)
	}

	// Test with integers in []interface{}
	data6 := &EmbeddingData{
		Embedding: []interface{}{1, 2, 3},
	}
	vec6 := data6.GetEmbeddingVector()
	if vec6 != nil {
		t.Errorf("expected nil for integer array, got %v", vec6)
	}
}

func TestGetEmbeddingBase64(t *testing.T) {
	// Test with string
	data := &EmbeddingData{
		Embedding: "YmFzZTY0ZW5jb2RlZA==",
	}
	str := data.GetEmbeddingBase64()
	if str != "YmFzZTY0ZW5jb2RlZA==" {
		t.Errorf("expected base64 string, got %q", str)
	}

	// Test with []float64
	data2 := &EmbeddingData{
		Embedding: []float64{0.1, 0.2, 0.3},
	}
	str2 := data2.GetEmbeddingBase64()
	if str2 != "" {
		t.Errorf("expected empty string for float array, got %q", str2)
	}
}

func TestEmbeddingProviderOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []EmbeddingOption
		check  func(*EmbeddingRequest) bool
		errMsg string
	}{
		{
			name: "WithEmbeddingProviderOrder",
			opts: []EmbeddingOption{WithEmbeddingProviderOrder("openai", "azure")},
			check: func(r *EmbeddingRequest) bool {
				return r.Provider != nil &&
					len(r.Provider.Order) == 2 &&
					r.Provider.Order[0] == "openai" &&
					r.Provider.Order[1] == "azure"
			},
			errMsg: "provider order not set correctly",
		},
		{
			name: "WithEmbeddingIgnoreProviders",
			opts: []EmbeddingOption{WithEmbeddingIgnoreProviders("expensive-provider")},
			check: func(r *EmbeddingRequest) bool {
				return r.Provider != nil &&
					len(r.Provider.Ignore) == 1 &&
					r.Provider.Ignore[0] == "expensive-provider"
			},
			errMsg: "provider ignore not set correctly",
		},
		{
			name: "WithEmbeddingAllowFallbacks",
			opts: []EmbeddingOption{WithEmbeddingAllowFallbacks(false)},
			check: func(r *EmbeddingRequest) bool {
				return r.Provider != nil &&
					r.Provider.AllowFallbacks != nil &&
					*r.Provider.AllowFallbacks == false
			},
			errMsg: "allow fallbacks not set correctly",
		},
		{
			name: "WithEmbeddingDataCollection",
			opts: []EmbeddingOption{WithEmbeddingDataCollection("deny")},
			check: func(r *EmbeddingRequest) bool {
				return r.Provider != nil &&
					r.Provider.DataCollection == "deny"
			},
			errMsg: "data collection not set correctly",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &EmbeddingRequest{
				Input: "test",
				Model: "test-model",
			}

			for _, opt := range tc.opts {
				opt(req)
			}

			if !tc.check(req) {
				t.Error(tc.errMsg)
			}
		})
	}
}

func TestEmbeddingErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Invalid model",
				Type:    "invalid_request_error",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.CreateEmbedding(context.Background(), "test", "invalid-model")
	if err == nil {
		t.Error("expected error for invalid model")
	}

	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Errorf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", reqErr.StatusCode)
	}
}
