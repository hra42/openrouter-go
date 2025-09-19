package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client.apiKey != apiKey {
		t.Errorf("expected apiKey %q, got %q", apiKey, client.apiKey)
	}

	if client.baseURL != defaultBaseURL {
		t.Errorf("expected baseURL %q, got %q", defaultBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}

	if client.httpClient.Timeout != defaultTimeout {
		t.Errorf("expected timeout %v, got %v", defaultTimeout, client.httpClient.Timeout)
	}
}

func TestClientOptions(t *testing.T) {
	apiKey := "test-api-key"
	customBaseURL := "https://custom.api.com"
	customTimeout := 60 * time.Second
	customModel := "custom-model"
	referer := "https://example.com"
	appName := "TestApp"

	client := NewClient(apiKey,
		WithBaseURL(customBaseURL),
		WithTimeout(customTimeout),
		WithDefaultModel(customModel),
		WithReferer(referer),
		WithAppName(appName),
		WithHeader("X-Custom", "value"),
		WithRetry(5, 2*time.Second),
	)

	if client.baseURL != customBaseURL {
		t.Errorf("expected baseURL %q, got %q", customBaseURL, client.baseURL)
	}

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}

	if client.defaultModel != customModel {
		t.Errorf("expected defaultModel %q, got %q", customModel, client.defaultModel)
	}

	if client.referer != referer {
		t.Errorf("expected referer %q, got %q", referer, client.referer)
	}

	if client.appName != appName {
		t.Errorf("expected appName %q, got %q", appName, client.appName)
	}

	if client.customHeaders["X-Custom"] != "value" {
		t.Errorf("expected custom header X-Custom=value, got %q", client.customHeaders["X-Custom"])
	}

	if client.maxRetries != 5 {
		t.Errorf("expected maxRetries 5, got %d", client.maxRetries)
	}

	if client.retryDelay != 2*time.Second {
		t.Errorf("expected retryDelay 2s, got %v", client.retryDelay)
	}
}

func TestClientWithCustomHTTPClient(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 120 * time.Second,
	}

	client := NewClient("api-key", WithHTTPClient(customHTTPClient))

	if client.httpClient != customHTTPClient {
		t.Error("expected custom HTTP client to be used")
	}

	if client.httpClient.Timeout != 120*time.Second {
		t.Errorf("expected custom timeout 120s, got %v", client.httpClient.Timeout)
	}
}

func TestDoRequest(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", authHeader)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", contentType)
		}

		// Check custom headers
		if referer := r.Header.Get("HTTP-Referer"); referer != "https://test.com" {
			t.Errorf("expected HTTP-Referer header 'https://test.com', got %q", referer)
		}

		if appName := r.Header.Get("X-Title"); appName != "TestApp" {
			t.Errorf("expected X-Title header 'TestApp', got %q", appName)
		}

		// Return successful response
		response := ChatCompletionResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   "test-model",
			Choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    "assistant",
						Content: "Test response",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-key",
		WithBaseURL(server.URL),
		WithReferer("https://test.com"),
		WithAppName("TestApp"),
	)

	req := ChatCompletionRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	var resp ChatCompletionResponse
	err := client.doRequest(context.Background(), "POST", "/chat/completions", req, &resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "test-id" {
		t.Errorf("expected response ID 'test-id', got %q", resp.ID)
	}

	if len(resp.Choices) != 1 {
		t.Errorf("expected 1 choice, got %d", len(resp.Choices))
	}
}

func TestDoRequestErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedError  string
		isRequestError bool
	}{
		{
			name:       "400 Bad Request",
			statusCode: 400,
			responseBody: `{
				"error": {
					"message": "Invalid request",
					"type": "invalid_request_error"
				}
			}`,
			expectedError:  "Invalid request",
			isRequestError: true,
		},
		{
			name:       "401 Unauthorized",
			statusCode: 401,
			responseBody: `{
				"error": {
					"message": "Invalid API key",
					"type": "authentication_error"
				}
			}`,
			expectedError:  "Invalid API key",
			isRequestError: true,
		},
		{
			name:       "429 Rate Limit",
			statusCode: 429,
			responseBody: `{
				"error": {
					"message": "Rate limit exceeded",
					"type": "rate_limit_error"
				}
			}`,
			expectedError:  "Rate limit exceeded",
			isRequestError: true,
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     500,
			responseBody:   "Internal Server Error",
			expectedError:  "Internal Server Error",
			isRequestError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewClient("test-key",
				WithBaseURL(server.URL),
				WithRetry(0, 0), // Disable retries for testing
			)

			req := ChatCompletionRequest{
				Model: "test-model",
				Messages: []Message{
					{Role: "user", Content: "Hello"},
				},
			}

			var resp ChatCompletionResponse
			err := client.doRequest(context.Background(), "POST", "/chat/completions", req, &resp)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.isRequestError {
				reqErr, ok := err.(*RequestError)
				if !ok {
					t.Fatalf("expected RequestError, got %T", err)
				}

				if reqErr.StatusCode != tt.statusCode {
					t.Errorf("expected status code %d, got %d", tt.statusCode, reqErr.StatusCode)
				}

				if reqErr.Message != tt.expectedError {
					t.Errorf("expected error message %q, got %q", tt.expectedError, reqErr.Message)
				}
			}
		})
	}
}

func TestDoRequestRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++

		// Fail first 2 attempts with 500, succeed on third
		if attempts < 3 {
			w.WriteHeader(500)
			w.Write([]byte("Server Error"))
			return
		}

		// Success on third attempt
		response := ChatCompletionResponse{
			ID:    "success",
			Model: "test-model",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test-key",
		WithBaseURL(server.URL),
		WithRetry(3, 10*time.Millisecond), // Short delay for testing
	)

	req := ChatCompletionRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	var resp ChatCompletionResponse
	err := client.doRequestWithRetry(context.Background(), "POST", "/chat/completions", req, &resp)

	if err != nil {
		t.Fatalf("unexpected error after retries: %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}

	if resp.ID != "success" {
		t.Errorf("expected response ID 'success', got %q", resp.ID)
	}
}

func TestContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := NewClient("test-key", WithBaseURL(server.URL))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	req := ChatCompletionRequest{
		Model: "test-model",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	var resp ChatCompletionResponse
	err := client.doRequest(ctx, "POST", "/chat/completions", req, &resp)

	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded error, got %v", ctx.Err())
	}
}