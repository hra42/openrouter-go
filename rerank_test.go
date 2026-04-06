package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRerank(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/rerank" {
			t.Errorf("expected path /rerank, got %s", r.URL.Path)
		}

		// Verify request body
		var reqBody RerankRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		if reqBody.Model != "cohere/rerank-v3.5" {
			t.Errorf("expected model 'cohere/rerank-v3.5', got %q", reqBody.Model)
		}
		if reqBody.Query != "What is the capital of France?" {
			t.Errorf("expected query 'What is the capital of France?', got %q", reqBody.Query)
		}
		if len(reqBody.Documents) != 2 {
			t.Errorf("expected 2 documents, got %d", len(reqBody.Documents))
		}

		// Send response
		cost := 0.000001
		response := RerankResponse{
			ID:       "rerank-123",
			Model:    "cohere/rerank-v3.5",
			Provider: "Cohere",
			Results: []RerankResult{
				{
					Index:          0,
					RelevanceScore: 0.95,
					Document:       RerankDocument{Text: "Paris is the capital of France."},
				},
				{
					Index:          1,
					RelevanceScore: 0.12,
					Document:       RerankDocument{Text: "Berlin is the capital of Germany."},
				},
			},
			Usage: &RerankUsage{
				TotalTokens: 20,
				SearchUnits: 1,
				Cost:        &cost,
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

	documents := []string{
		"Paris is the capital of France.",
		"Berlin is the capital of Germany.",
	}

	resp, err := client.Rerank(context.Background(), "What is the capital of France?", documents, "cohere/rerank-v3.5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "rerank-123" {
		t.Errorf("expected id 'rerank-123', got %q", resp.ID)
	}
	if resp.Model != "cohere/rerank-v3.5" {
		t.Errorf("expected model 'cohere/rerank-v3.5', got %q", resp.Model)
	}
	if resp.Provider != "Cohere" {
		t.Errorf("expected provider 'Cohere', got %q", resp.Provider)
	}
	if len(resp.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(resp.Results))
	}
	if resp.Results[0].Index != 0 {
		t.Errorf("expected first result index 0, got %d", resp.Results[0].Index)
	}
	if resp.Results[0].RelevanceScore != 0.95 {
		t.Errorf("expected first result relevance_score 0.95, got %f", resp.Results[0].RelevanceScore)
	}
	if resp.Results[0].Document.Text != "Paris is the capital of France." {
		t.Errorf("expected first result document text, got %q", resp.Results[0].Document.Text)
	}
	if resp.Usage == nil {
		t.Fatal("expected usage to be set")
	}
	if resp.Usage.TotalTokens != 20 {
		t.Errorf("expected 20 total tokens, got %d", resp.Usage.TotalTokens)
	}
	if resp.Usage.SearchUnits != 1 {
		t.Errorf("expected 1 search unit, got %d", resp.Usage.SearchUnits)
	}
}

func TestRerankWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody RerankRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		// Verify options were applied
		if reqBody.TopN == nil || *reqBody.TopN != 2 {
			t.Errorf("expected top_n 2, got %v", reqBody.TopN)
		}
		if reqBody.Provider == nil {
			t.Error("expected provider to be set")
		} else {
			if len(reqBody.Provider.Only) != 1 || reqBody.Provider.Only[0] != "Cohere" {
				t.Errorf("expected only=['Cohere'], got %v", reqBody.Provider.Only)
			}
		}

		// Return only top 2 results
		response := RerankResponse{
			Model: "cohere/rerank-v3.5",
			Results: []RerankResult{
				{
					Index:          0,
					RelevanceScore: 0.95,
					Document:       RerankDocument{Text: "Paris is the capital of France."},
				},
				{
					Index:          2,
					RelevanceScore: 0.45,
					Document:       RerankDocument{Text: "France is a country in Europe."},
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

	documents := []string{
		"Paris is the capital of France.",
		"Berlin is the capital of Germany.",
		"France is a country in Europe.",
	}

	resp, err := client.Rerank(
		context.Background(),
		"What is the capital of France?",
		documents,
		"cohere/rerank-v3.5",
		WithRerankTopN(2),
		WithRerankOnlyProviders("Cohere"),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(resp.Results))
	}
}

func TestRerankValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Test missing query
	_, err := client.Rerank(context.Background(), "", []string{"doc"}, "cohere/rerank-v3.5")
	if err == nil {
		t.Error("expected error for empty query")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "query" {
		t.Errorf("expected validation error for query field, got %v", err)
	}

	// Test empty documents
	_, err = client.Rerank(context.Background(), "query", []string{}, "cohere/rerank-v3.5")
	if err == nil {
		t.Error("expected error for empty documents")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "documents" {
		t.Errorf("expected validation error for documents field, got %v", err)
	}

	// Test nil documents
	_, err = client.Rerank(context.Background(), "query", nil, "cohere/rerank-v3.5")
	if err == nil {
		t.Error("expected error for nil documents")
	}
	if validErr, ok := err.(*ValidationError); !ok || validErr.Field != "documents" {
		t.Errorf("expected validation error for documents field, got %v", err)
	}

	// Test missing model
	_, err = client.Rerank(context.Background(), "query", []string{"doc"}, "")
	if err == nil {
		t.Error("expected error for empty model")
	}
	if err != ErrNoModel {
		t.Errorf("expected ErrNoModel, got %v", err)
	}

	// Test missing API key
	clientNoKey := NewClient()
	_, err = clientNoKey.Rerank(context.Background(), "query", []string{"doc"}, "model")
	if err == nil {
		t.Error("expected error for missing API key")
	}
	if err != ErrNoAPIKey {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}

func TestRerankProviderOptions(t *testing.T) {
	tests := []struct {
		name   string
		opts   []RerankOption
		check  func(*RerankRequest) bool
		errMsg string
	}{
		{
			name: "WithRerankProviderOrder",
			opts: []RerankOption{WithRerankProviderOrder("Cohere", "Google")},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					len(r.Provider.Order) == 2 &&
					r.Provider.Order[0] == "Cohere" &&
					r.Provider.Order[1] == "Google"
			},
			errMsg: "provider order not set correctly",
		},
		{
			name: "WithRerankIgnoreProviders",
			opts: []RerankOption{WithRerankIgnoreProviders("expensive-provider")},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					len(r.Provider.Ignore) == 1 &&
					r.Provider.Ignore[0] == "expensive-provider"
			},
			errMsg: "provider ignore not set correctly",
		},
		{
			name: "WithRerankAllowFallbacks",
			opts: []RerankOption{WithRerankAllowFallbacks(false)},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					r.Provider.AllowFallbacks != nil &&
					*r.Provider.AllowFallbacks == false
			},
			errMsg: "allow fallbacks not set correctly",
		},
		{
			name: "WithRerankDataCollection",
			opts: []RerankOption{WithRerankDataCollection("deny")},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					r.Provider.DataCollection == "deny"
			},
			errMsg: "data collection not set correctly",
		},
		{
			name: "WithRerankProviderSort",
			opts: []RerankOption{WithRerankProviderSort("price")},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					r.Provider.Sort == "price"
			},
			errMsg: "provider sort not set correctly",
		},
		{
			name: "WithRerankTopN",
			opts: []RerankOption{WithRerankTopN(5)},
			check: func(r *RerankRequest) bool {
				return r.TopN != nil && *r.TopN == 5
			},
			errMsg: "top_n not set correctly",
		},
		{
			name: "WithRerankProvider",
			opts: []RerankOption{WithRerankProvider(Provider{
				Order: []string{"Cohere"},
				Sort:  "price",
			})},
			check: func(r *RerankRequest) bool {
				return r.Provider != nil &&
					len(r.Provider.Order) == 1 &&
					r.Provider.Order[0] == "Cohere" &&
					r.Provider.Sort == "price"
			},
			errMsg: "provider not set correctly",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := &RerankRequest{
				Model:     "test-model",
				Query:     "test query",
				Documents: []string{"doc1"},
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

func TestRerankErrorResponse(t *testing.T) {
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

	_, err := client.Rerank(context.Background(), "query", []string{"doc"}, "invalid-model")
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

func BenchmarkRerankRoundTrip(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := RerankResponse{
			ID:    "rerank-bench",
			Model: "test-model",
			Results: []RerankResult{
				{
					Index:          0,
					RelevanceScore: 0.95,
					Document:       RerankDocument{Text: "doc1"},
				},
			},
			Usage: &RerankUsage{TotalTokens: 10, SearchUnits: 1},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	ctx := context.Background()
	documents := []string{"doc1", "doc2"}
	b.ResetTimer()
	for b.Loop() {
		_, _ = client.Rerank(ctx, "query", documents, "test-model")
	}
}
