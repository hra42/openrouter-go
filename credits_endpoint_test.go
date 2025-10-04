package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCredits(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/credits" {
			t.Errorf("expected path /credits, got %s", r.URL.Path)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		// Send response
		response := CreditsResponse{
			Data: CreditsData{
				TotalCredits: 100.50,
				TotalUsage:   25.75,
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

	resp, err := client.GetCredits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.TotalCredits != 100.50 {
		t.Errorf("expected TotalCredits 100.50, got %f", resp.Data.TotalCredits)
	}
	if resp.Data.TotalUsage != 25.75 {
		t.Errorf("expected TotalUsage 25.75, got %f", resp.Data.TotalUsage)
	}

	// Verify remaining calculation
	remaining := resp.Data.TotalCredits - resp.Data.TotalUsage
	expectedRemaining := 74.75
	if remaining != expectedRemaining {
		t.Errorf("expected remaining %.2f, got %.2f", expectedRemaining, remaining)
	}
}

func TestGetCreditsZeroBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := CreditsResponse{
			Data: CreditsData{
				TotalCredits: 0.0,
				TotalUsage:   0.0,
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

	resp, err := client.GetCredits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.TotalCredits != 0.0 {
		t.Errorf("expected TotalCredits 0.0, got %f", resp.Data.TotalCredits)
	}
	if resp.Data.TotalUsage != 0.0 {
		t.Errorf("expected TotalUsage 0.0, got %f", resp.Data.TotalUsage)
	}
}

func TestGetCreditsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Invalid API key",
				Type:    "authentication_error",
				Code:    "invalid_api_key",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.GetCredits(context.Background())
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
	if reqErr.Message != "Invalid API key" {
		t.Errorf("expected error message 'Invalid API key', got %q", reqErr.Message)
	}
}
