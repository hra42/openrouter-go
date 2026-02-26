package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListZDREndpoints(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/endpoints/zdr" {
			t.Errorf("expected path /endpoints/zdr, got %s", r.URL.Path)
		}

		// Verify auth header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization 'Bearer test-key', got %q", auth)
		}

		// Send response
		contextLength := 128000.0
		maxCompTokens := 16384.0
		maxPromptTokens := 100000.0
		uptime := 0.99
		quantization := "fp16"
		tag := "production"
		supportsCache := true

		response := ZDREndpointsResponse{
			Data: []PublicEndpoint{
				{
					Name:                    "openai/gpt-4o",
					ModelID:                 "openai/gpt-4o",
					ModelName:               "GPT-4o",
					ContextLength:           contextLength,
					ProviderName:            "OpenAI",
					Tag:                     &tag,
					Quantization:            &quantization,
					MaxCompletionTokens:     &maxCompTokens,
					MaxPromptTokens:         &maxPromptTokens,
					SupportedParameters:     []string{"temperature", "top_p", "max_tokens"},
					Status:                  1,
					UptimeLast30m:           &uptime,
					SupportsImplicitCaching: &supportsCache,
					LatencyLast30m: &PercentileStats{
						P50: 150.0,
						P75: 200.0,
						P90: 350.0,
						P99: 500.0,
					},
					ThroughputLast30m: &PercentileStats{
						P50: 50.0,
						P75: 45.0,
						P90: 35.0,
						P99: 20.0,
					},
					Pricing: PublicEndpointPricing{
						Prompt:     "0.005",
						Completion: "0.015",
						Request:    "0",
						Image:      "0.003",
					},
				},
				{
					Name:                "anthropic/claude-3.5-sonnet",
					ModelID:             "anthropic/claude-3.5-sonnet",
					ModelName:           "Claude 3.5 Sonnet",
					ContextLength:       200000.0,
					ProviderName:        "Anthropic",
					SupportedParameters: []string{"temperature", "top_p"},
					Status:              1,
					Pricing: PublicEndpointPricing{
						Prompt:     "0.003",
						Completion: "0.015",
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

	resp, err := client.ListZDREndpoints(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify response data
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 endpoints, got %d", len(resp.Data))
	}

	// Verify first endpoint
	ep := resp.Data[0]
	if ep.Name != "openai/gpt-4o" {
		t.Errorf("expected name 'openai/gpt-4o', got %q", ep.Name)
	}
	if ep.ModelID != "openai/gpt-4o" {
		t.Errorf("expected model ID 'openai/gpt-4o', got %q", ep.ModelID)
	}
	if ep.ModelName != "GPT-4o" {
		t.Errorf("expected model name 'GPT-4o', got %q", ep.ModelName)
	}
	if ep.ProviderName != "OpenAI" {
		t.Errorf("expected provider name 'OpenAI', got %q", ep.ProviderName)
	}
	if ep.ContextLength != 128000.0 {
		t.Errorf("expected context length 128000.0, got %f", ep.ContextLength)
	}
	if ep.Status != 1 {
		t.Errorf("expected status 1, got %f", ep.Status)
	}

	// Verify pricing
	if ep.Pricing.Prompt != "0.005" {
		t.Errorf("expected prompt pricing '0.005', got %q", ep.Pricing.Prompt)
	}
	if ep.Pricing.Completion != "0.015" {
		t.Errorf("expected completion pricing '0.015', got %q", ep.Pricing.Completion)
	}

	// Verify optional fields
	if ep.Tag == nil || *ep.Tag != "production" {
		t.Error("expected tag 'production'")
	}
	if ep.Quantization == nil || *ep.Quantization != "fp16" {
		t.Error("expected quantization 'fp16'")
	}
	if ep.SupportsImplicitCaching == nil || !*ep.SupportsImplicitCaching {
		t.Error("expected supports_implicit_caching to be true")
	}

	// Verify latency stats
	if ep.LatencyLast30m == nil {
		t.Fatal("expected latency stats")
	}
	if ep.LatencyLast30m.P50 != 150.0 {
		t.Errorf("expected latency P50 150.0, got %f", ep.LatencyLast30m.P50)
	}
	if ep.LatencyLast30m.P99 != 500.0 {
		t.Errorf("expected latency P99 500.0, got %f", ep.LatencyLast30m.P99)
	}

	// Verify throughput stats
	if ep.ThroughputLast30m == nil {
		t.Fatal("expected throughput stats")
	}
	if ep.ThroughputLast30m.P50 != 50.0 {
		t.Errorf("expected throughput P50 50.0, got %f", ep.ThroughputLast30m.P50)
	}

	// Verify second endpoint has nil optional fields
	ep2 := resp.Data[1]
	if ep2.ModelName != "Claude 3.5 Sonnet" {
		t.Errorf("expected model name 'Claude 3.5 Sonnet', got %q", ep2.ModelName)
	}
	if ep2.Tag != nil {
		t.Errorf("expected nil tag, got %v", ep2.Tag)
	}
	if ep2.LatencyLast30m != nil {
		t.Errorf("expected nil latency stats, got %v", ep2.LatencyLast30m)
	}
}

func TestListZDREndpointsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/endpoints/zdr" {
			t.Errorf("expected path /endpoints/zdr, got %s", r.URL.Path)
		}

		response := ZDREndpointsResponse{Data: []PublicEndpoint{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)

	resp, err := client.ListZDREndpoints(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("expected 0 endpoints, got %d", len(resp.Data))
	}
}

func TestListZDREndpointsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Invalid API key",
				Type:    "authentication_error",
				Code:    "401",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("invalid-key"),
		WithBaseURL(server.URL),
	)

	_, err := client.ListZDREndpoints(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	reqErr, ok := IsRequestError(err)
	if !ok {
		t.Errorf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status code 401, got %d", reqErr.StatusCode)
	}
}
