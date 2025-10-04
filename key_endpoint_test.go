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

func TestListKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/keys" {
			t.Errorf("expected path /keys, got %s", r.URL.Path)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		// Send response
		response := ListKeysResponse{
			Data: []APIKey{
				{
					Name:      "sk-or-v1-abc123",
					Label:     "Production Key",
					Limit:     100.0,
					Disabled:  false,
					CreatedAt: "2024-01-01T00:00:00Z",
					UpdatedAt: "2024-01-02T00:00:00Z",
					Hash:      "abc123hash",
				},
				{
					Name:      "sk-or-v1-def456",
					Label:     "Development Key",
					Limit:     50.0,
					Disabled:  false,
					CreatedAt: "2024-01-03T00:00:00Z",
					UpdatedAt: "2024-01-04T00:00:00Z",
					Hash:      "def456hash",
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

	resp, err := client.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 keys, got %d", len(resp.Data))
	}

	// Validate first key
	if resp.Data[0].Label != "Production Key" {
		t.Errorf("expected Label 'Production Key', got %q", resp.Data[0].Label)
	}
	if resp.Data[0].Limit != 100.0 {
		t.Errorf("expected Limit 100.0, got %f", resp.Data[0].Limit)
	}
	if resp.Data[0].Disabled != false {
		t.Errorf("expected Disabled false, got %t", resp.Data[0].Disabled)
	}
	if resp.Data[0].Hash != "abc123hash" {
		t.Errorf("expected Hash 'abc123hash', got %q", resp.Data[0].Hash)
	}
}

func TestListKeysWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()

		offset := query.Get("offset")
		if offset != "10" {
			t.Errorf("expected offset '10', got %q", offset)
		}

		includeDisabled := query.Get("include_disabled")
		if includeDisabled != "true" {
			t.Errorf("expected include_disabled 'true', got %q", includeDisabled)
		}

		response := ListKeysResponse{
			Data: []APIKey{
				{
					Name:      "sk-or-v1-disabled",
					Label:     "Disabled Key",
					Limit:     25.0,
					Disabled:  true,
					CreatedAt: "2024-01-05T00:00:00Z",
					UpdatedAt: "2024-01-06T00:00:00Z",
					Hash:      "disabledhash",
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

	offset := 10
	includeDisabled := true

	resp, err := client.ListKeys(context.Background(), &ListKeysOptions{
		Offset:          &offset,
		IncludeDisabled: &includeDisabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("expected 1 key, got %d", len(resp.Data))
	}

	if resp.Data[0].Disabled != true {
		t.Errorf("expected Disabled true, got %t", resp.Data[0].Disabled)
	}
}

func TestListKeysEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ListKeysResponse{
			Data: []APIKey{},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.ListKeys(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 keys, got %d", len(resp.Data))
	}
}

func TestListKeysError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Provisioning key required",
				Type:    "authentication_error",
				Code:    "invalid_key_type",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListKeys(context.Background(), nil)
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

func TestCreateKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/keys" {
			t.Errorf("expected path /keys, got %s", r.URL.Path)
		}

		// Verify Authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		// Read and verify request body
		var reqBody CreateKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Name != "Test API Key" {
			t.Errorf("expected Name 'Test API Key', got %q", reqBody.Name)
		}
		if reqBody.Limit == nil || *reqBody.Limit != 50.0 {
			t.Errorf("expected Limit 50.0, got %v", reqBody.Limit)
		}

		// Send response
		response := CreateKeyResponse{
			Data: APIKey{
				Name:      "sk-or-v1-newkey123",
				Label:     "Test API Key",
				Limit:     50.0,
				Disabled:  false,
				CreatedAt: "2024-01-10T00:00:00Z",
				UpdatedAt: "2024-01-10T00:00:00Z",
				Hash:      "newkeyhash",
			},
			Key: "sk-or-v1-newkey123-actual-secret-value",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	limit := 50.0
	resp, err := client.CreateKey(context.Background(), &CreateKeyRequest{
		Name:  "Test API Key",
		Limit: &limit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate response
	if resp.Data.Label != "Test API Key" {
		t.Errorf("expected Label 'Test API Key', got %q", resp.Data.Label)
	}
	if resp.Data.Limit != 50.0 {
		t.Errorf("expected Limit 50.0, got %f", resp.Data.Limit)
	}
	if resp.Key != "sk-or-v1-newkey123-actual-secret-value" {
		t.Errorf("expected Key 'sk-or-v1-newkey123-actual-secret-value', got %q", resp.Key)
	}
	if resp.Data.Disabled != false {
		t.Errorf("expected Disabled false, got %t", resp.Data.Disabled)
	}
}

func TestCreateKeyMinimal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read and verify request body
		var reqBody CreateKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Name != "Minimal Key" {
			t.Errorf("expected Name 'Minimal Key', got %q", reqBody.Name)
		}
		if reqBody.Limit != nil {
			t.Errorf("expected Limit to be nil, got %v", reqBody.Limit)
		}
		if reqBody.IncludeBYOKInLimit != nil {
			t.Errorf("expected IncludeBYOKInLimit to be nil, got %v", reqBody.IncludeBYOKInLimit)
		}

		response := CreateKeyResponse{
			Data: APIKey{
				Name:      "sk-or-v1-minimal123",
				Label:     "Minimal Key",
				Limit:     0,
				Disabled:  false,
				CreatedAt: "2024-01-10T00:00:00Z",
				UpdatedAt: "2024-01-10T00:00:00Z",
				Hash:      "minimalhash",
			},
			Key: "sk-or-v1-minimal123-secret",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.CreateKey(context.Background(), &CreateKeyRequest{
		Name: "Minimal Key",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.Label != "Minimal Key" {
		t.Errorf("expected Label 'Minimal Key', got %q", resp.Data.Label)
	}
	if resp.Key == "" {
		t.Error("expected Key to be set")
	}
}

func TestCreateKeyWithBYOK(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateKeyRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.IncludeBYOKInLimit == nil || *reqBody.IncludeBYOKInLimit != true {
			t.Errorf("expected IncludeBYOKInLimit true, got %v", reqBody.IncludeBYOKInLimit)
		}

		response := CreateKeyResponse{
			Data: APIKey{
				Name:      "sk-or-v1-byok123",
				Label:     "BYOK Key",
				Limit:     100.0,
				Disabled:  false,
				CreatedAt: "2024-01-10T00:00:00Z",
				UpdatedAt: "2024-01-10T00:00:00Z",
				Hash:      "byokhash",
			},
			Key: "sk-or-v1-byok123-secret",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	limit := 100.0
	includeBYOK := true
	resp, err := client.CreateKey(context.Background(), &CreateKeyRequest{
		Name:               "BYOK Key",
		Limit:              &limit,
		IncludeBYOKInLimit: &includeBYOK,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.Label != "BYOK Key" {
		t.Errorf("expected Label 'BYOK Key', got %q", resp.Data.Label)
	}
}

func TestCreateKeyValidation(t *testing.T) {
	client := NewClient(
		WithAPIKey("test-key"),
	)

	// Test nil request
	_, err := client.CreateKey(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Test empty name
	_, err = client.CreateKey(context.Background(), &CreateKeyRequest{
		Name: "",
	})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if !IsValidationError(err) {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestCreateKeyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Provisioning key required",
				Type:    "authentication_error",
				Code:    "invalid_key_type",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.CreateKey(context.Background(), &CreateKeyRequest{
		Name: "Test Key",
	})
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
