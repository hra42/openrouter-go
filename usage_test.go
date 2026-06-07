package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWithUsageOption verifies that WithUsage sets the usage accounting field on
// the chat completion request body.
func TestWithUsageOption(t *testing.T) {
	tests := []struct {
		name        string
		option      ChatCompletionOption
		wantNil     bool
		wantInclude bool
	}{
		{
			name:        "enabled",
			option:      WithUsage(true),
			wantInclude: true,
		},
		{
			name:        "disabled",
			option:      WithUsage(false),
			wantInclude: false,
		},
		{
			name:    "absent",
			option:  nil,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req ChatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					t.Fatalf("decode request: %v", err)
				}

				if tt.wantNil {
					if req.Usage != nil {
						t.Errorf("expected Usage to be nil, got %+v", req.Usage)
					}
				} else {
					if req.Usage == nil {
						t.Fatal("expected Usage to be set, got nil")
					}
					if req.Usage.Include != tt.wantInclude {
						t.Errorf("Usage.Include = %v, want %v", req.Usage.Include, tt.wantInclude)
					}
				}

				w.WriteHeader(200)
				_ = json.NewEncoder(w).Encode(ChatCompletionResponse{ID: "test"})
			}))
			defer server.Close()

			client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
			opts := []ChatCompletionOption{WithModel("test-model")}
			if tt.option != nil {
				opts = append(opts, tt.option)
			}
			if _, err := client.ChatComplete(context.Background(), []Message{CreateUserMessage("Test")}, opts...); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestWithCompletionUsageOption verifies that WithCompletionUsage sets the usage
// accounting field on the legacy completion request body.
func TestWithCompletionUsageOption(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Usage == nil || !req.Usage.Include {
			t.Errorf("expected Usage.Include to be true, got %+v", req.Usage)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(CompletionResponse{ID: "test"})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	if _, err := client.Complete(context.Background(), "Test prompt",
		WithCompletionModel("test-model"),
		WithCompletionUsage(true),
	); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestUsageCostDecoding verifies that the enriched usage accounting fields decode
// correctly from an API response.
func TestUsageCostDecoding(t *testing.T) {
	body := `{
		"id": "gen-123",
		"model": "openai/gpt-4o-mini",
		"choices": [{"index": 0, "message": {"role": "assistant", "content": "hi"}, "finish_reason": "stop"}],
		"usage": {
			"prompt_tokens": 12,
			"completion_tokens": 8,
			"total_tokens": 20,
			"cost": 0.0000123,
			"cost_details": {"upstream_inference_cost": 0.0000111},
			"prompt_tokens_details": {"cached_tokens": 4},
			"completion_tokens_details": {"reasoning_tokens": 3}
		}
	}`

	var resp ChatCompletionResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	u := resp.Usage
	if u.PromptTokens != 12 || u.CompletionTokens != 8 || u.TotalTokens != 20 {
		t.Errorf("token counts wrong: %+v", u)
	}
	if u.Cost == nil || *u.Cost != 0.0000123 {
		t.Errorf("Cost = %v, want 0.0000123", u.Cost)
	}
	if u.CostDetails == nil || u.CostDetails.UpstreamInferenceCost == nil || *u.CostDetails.UpstreamInferenceCost != 0.0000111 {
		t.Errorf("CostDetails wrong: %+v", u.CostDetails)
	}
	if u.PromptTokensDetails == nil || u.PromptTokensDetails.CachedTokens != 4 {
		t.Errorf("PromptTokensDetails wrong: %+v", u.PromptTokensDetails)
	}
	if u.CompletionTokensDetails == nil || u.CompletionTokensDetails.ReasoningTokens != 3 {
		t.Errorf("CompletionTokensDetails wrong: %+v", u.CompletionTokensDetails)
	}
}

// TestUsageCostOmittedWhenAbsent verifies the cost fields stay nil and are omitted
// from JSON when usage accounting is not enabled.
func TestUsageCostOmittedWhenAbsent(t *testing.T) {
	var resp ChatCompletionResponse
	body := `{"id":"x","usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}}`
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Usage.Cost != nil || resp.Usage.CostDetails != nil {
		t.Errorf("cost fields should be nil when absent: %+v", resp.Usage)
	}

	out, err := json.Marshal(resp.Usage)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(out); got != `{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}` {
		t.Errorf("cost fields should be omitted, got %s", got)
	}
}
