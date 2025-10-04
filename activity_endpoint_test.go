package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetActivity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/activity" {
			t.Errorf("expected path /activity, got %s", r.URL.Path)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		// Send response
		response := ActivityResponse{
			Data: []ActivityData{
				{
					Date:               "2024-01-15",
					Model:              "openai/gpt-4",
					ModelPermaslug:     "openai/gpt-4",
					EndpointID:         "endpoint-1",
					ProviderName:       "OpenAI",
					Usage:              12.50,
					BYOKUsageInference: 0.0,
					Requests:           100,
					PromptTokens:       5000,
					CompletionTokens:   2500,
					ReasoningTokens:    0,
				},
				{
					Date:               "2024-01-15",
					Model:              "anthropic/claude-3-opus",
					ModelPermaslug:     "anthropic/claude-3-opus",
					EndpointID:         "endpoint-2",
					ProviderName:       "Anthropic",
					Usage:              25.75,
					BYOKUsageInference: 5.0,
					Requests:           50,
					PromptTokens:       10000,
					CompletionTokens:   5000,
					ReasoningTokens:    1000,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.GetActivity(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 activity records, got %d", len(resp.Data))
	}

	// Verify first record
	first := resp.Data[0]
	if first.Date != "2024-01-15" {
		t.Errorf("expected Date '2024-01-15', got %q", first.Date)
	}
	if first.Model != "openai/gpt-4" {
		t.Errorf("expected Model 'openai/gpt-4', got %q", first.Model)
	}
	if first.Usage != 12.50 {
		t.Errorf("expected Usage 12.50, got %f", first.Usage)
	}
	if first.Requests != 100 {
		t.Errorf("expected Requests 100, got %f", first.Requests)
	}

	// Verify second record
	second := resp.Data[1]
	if second.Model != "anthropic/claude-3-opus" {
		t.Errorf("expected Model 'anthropic/claude-3-opus', got %q", second.Model)
	}
	if second.ReasoningTokens != 1000 {
		t.Errorf("expected ReasoningTokens 1000, got %f", second.ReasoningTokens)
	}
}

func TestGetActivityWithDate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameter
		date := r.URL.Query().Get("date")
		if date != "2024-01-15" {
			t.Errorf("expected date query parameter '2024-01-15', got %q", date)
		}

		response := ActivityResponse{
			Data: []ActivityData{
				{
					Date:               "2024-01-15",
					Model:              "openai/gpt-4",
					ModelPermaslug:     "openai/gpt-4",
					EndpointID:         "endpoint-1",
					ProviderName:       "OpenAI",
					Usage:              10.0,
					BYOKUsageInference: 0.0,
					Requests:           50,
					PromptTokens:       2500,
					CompletionTokens:   1250,
					ReasoningTokens:    0,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.GetActivity(context.Background(), &ActivityOptions{
		Date: "2024-01-15",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 activity record, got %d", len(resp.Data))
	}

	if resp.Data[0].Date != "2024-01-15" {
		t.Errorf("expected Date '2024-01-15', got %q", resp.Data[0].Date)
	}
}

func TestGetActivityEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ActivityResponse{
			Data: []ActivityData{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.GetActivity(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 activity records, got %d", len(resp.Data))
	}
}

func TestGetActivityError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Provisioning key required",
				Type:    "authentication_error",
				Code:    "provisioning_key_required",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.GetActivity(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, reqErr.StatusCode)
	}
	if reqErr.Message != "Provisioning key required" {
		t.Errorf("expected error message 'Provisioning key required', got %q", reqErr.Message)
	}
}
