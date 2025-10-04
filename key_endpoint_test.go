package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetKey(t *testing.T) {
	limit := 100.0
	limitRemaining := 75.5

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/key" {
			t.Errorf("expected path /key, got %s", r.URL.Path)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		// Send response
		response := KeyResponse{
			Data: KeyData{
				Label:             "My API Key",
				Limit:             &limit,
				Usage:             24.5,
				IsFreeTier:        false,
				LimitRemaining:    &limitRemaining,
				IsProvisioningKey: false,
				RateLimit: &KeyRateLimit{
					Interval: "10s",
					Requests: 200,
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

	resp, err := client.GetKey(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.Label != "My API Key" {
		t.Errorf("expected Label 'My API Key', got %q", resp.Data.Label)
	}
	if resp.Data.Limit == nil || *resp.Data.Limit != 100.0 {
		t.Errorf("expected Limit 100.0, got %v", resp.Data.Limit)
	}
	if resp.Data.Usage != 24.5 {
		t.Errorf("expected Usage 24.5, got %f", resp.Data.Usage)
	}
	if resp.Data.IsFreeTier != false {
		t.Errorf("expected IsFreeTier false, got %t", resp.Data.IsFreeTier)
	}
	if resp.Data.LimitRemaining == nil || *resp.Data.LimitRemaining != 75.5 {
		t.Errorf("expected LimitRemaining 75.5, got %v", resp.Data.LimitRemaining)
	}
	if resp.Data.IsProvisioningKey != false {
		t.Errorf("expected IsProvisioningKey false, got %t", resp.Data.IsProvisioningKey)
	}
	if resp.Data.RateLimit == nil {
		t.Fatal("expected RateLimit to be present, got nil")
	}
	if resp.Data.RateLimit.Interval != "10s" {
		t.Errorf("expected RateLimit.Interval '10s', got %q", resp.Data.RateLimit.Interval)
	}
	if resp.Data.RateLimit.Requests != 200 {
		t.Errorf("expected RateLimit.Requests 200, got %f", resp.Data.RateLimit.Requests)
	}
}

func TestGetKeyFreeTier(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := KeyResponse{
			Data: KeyData{
				Label:             "Free Tier Key",
				Limit:             nil,
				Usage:             5.25,
				IsFreeTier:        true,
				LimitRemaining:    nil,
				IsProvisioningKey: false,
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

	resp, err := client.GetKey(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.Label != "Free Tier Key" {
		t.Errorf("expected Label 'Free Tier Key', got %q", resp.Data.Label)
	}
	if resp.Data.Limit != nil {
		t.Errorf("expected Limit nil, got %v", resp.Data.Limit)
	}
	if resp.Data.IsFreeTier != true {
		t.Errorf("expected IsFreeTier true, got %t", resp.Data.IsFreeTier)
	}
	if resp.Data.LimitRemaining != nil {
		t.Errorf("expected LimitRemaining nil, got %v", resp.Data.LimitRemaining)
	}
	if resp.Data.RateLimit != nil {
		t.Errorf("expected RateLimit nil, got %v", resp.Data.RateLimit)
	}
}

func TestGetKeyProvisioningKey(t *testing.T) {
	limit := 1000.0
	limitRemaining := 1000.0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := KeyResponse{
			Data: KeyData{
				Label:             "Provisioning Key",
				Limit:             &limit,
				Usage:             0.0,
				IsFreeTier:        false,
				LimitRemaining:    &limitRemaining,
				IsProvisioningKey: true,
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

	resp, err := client.GetKey(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.IsProvisioningKey != true {
		t.Errorf("expected IsProvisioningKey true, got %t", resp.Data.IsProvisioningKey)
	}
	if resp.Data.Usage != 0.0 {
		t.Errorf("expected Usage 0.0, got %f", resp.Data.Usage)
	}
}

func TestGetKeyError(t *testing.T) {
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

	_, err := client.GetKey(context.Background())
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
