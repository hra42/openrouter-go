package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProviders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/providers" {
			t.Errorf("expected path /providers, got %s", r.URL.Path)
		}

		// Send response
		response := ProvidersResponse{
			Data: []ProviderInfo{
				{
					Name:               "OpenAI",
					Slug:               "openai",
					PrivacyPolicyURL:   stringPtr("https://openai.com/privacy"),
					TermsOfServiceURL:  stringPtr("https://openai.com/terms"),
					StatusPageURL:      stringPtr("https://status.openai.com"),
				},
				{
					Name:               "Anthropic",
					Slug:               "anthropic",
					PrivacyPolicyURL:   stringPtr("https://anthropic.com/privacy"),
					TermsOfServiceURL:  stringPtr("https://anthropic.com/terms"),
					StatusPageURL:      nil,
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

	resp, err := client.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 providers, got %d", len(resp.Data))
	}

	// Check first provider
	if resp.Data[0].Name != "OpenAI" {
		t.Errorf("expected first provider name 'OpenAI', got %q", resp.Data[0].Name)
	}
	if resp.Data[0].Slug != "openai" {
		t.Errorf("expected first provider slug 'openai', got %q", resp.Data[0].Slug)
	}
	if resp.Data[0].PrivacyPolicyURL == nil || *resp.Data[0].PrivacyPolicyURL != "https://openai.com/privacy" {
		t.Errorf("unexpected privacy policy URL for OpenAI")
	}

	// Check second provider with nil StatusPageURL
	if resp.Data[1].StatusPageURL != nil {
		t.Errorf("expected nil StatusPageURL for Anthropic, got %v", resp.Data[1].StatusPageURL)
	}
}

func TestListProvidersEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ProvidersResponse{
			Data: []ProviderInfo{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.ListProviders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 providers, got %d", len(resp.Data))
	}
}
