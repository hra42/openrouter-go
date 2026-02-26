package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExchangeAuthCode(t *testing.T) {
	verifier := "test-verifier-value"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/auth/keys" {
			t.Errorf("expected path /auth/keys, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		var req ExchangeAuthCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.Code != "test-auth-code" {
			t.Errorf("expected code 'test-auth-code', got %q", req.Code)
		}
		if req.CodeVerifier == nil || *req.CodeVerifier != verifier {
			t.Errorf("expected code_verifier %q, got %v", verifier, req.CodeVerifier)
		}
		if req.CodeChallengeMethod != CodeChallengeMethodS256 {
			t.Errorf("expected code_challenge_method 'S256', got %q", req.CodeChallengeMethod)
		}

		userID := "user-123"
		response := ExchangeAuthCodeResponse{
			Key:    "sk-or-v1-test-key-value",
			UserID: &userID,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ExchangeAuthCode(context.Background(), &ExchangeAuthCodeRequest{
		Code:                "test-auth-code",
		CodeVerifier:        &verifier,
		CodeChallengeMethod: CodeChallengeMethodS256,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Key != "sk-or-v1-test-key-value" {
		t.Errorf("expected key 'sk-or-v1-test-key-value', got %q", resp.Key)
	}
	if resp.UserID == nil || *resp.UserID != "user-123" {
		t.Errorf("expected user_id 'user-123', got %v", resp.UserID)
	}
}

func TestExchangeAuthCodeWithoutPKCE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ExchangeAuthCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.CodeVerifier != nil {
			t.Errorf("expected nil code_verifier, got %v", req.CodeVerifier)
		}
		if req.CodeChallengeMethod != "" {
			t.Errorf("expected empty code_challenge_method, got %q", req.CodeChallengeMethod)
		}

		response := ExchangeAuthCodeResponse{
			Key: "sk-or-v1-no-pkce-key",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ExchangeAuthCode(context.Background(), &ExchangeAuthCodeRequest{
		Code: "test-auth-code",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Key != "sk-or-v1-no-pkce-key" {
		t.Errorf("expected key 'sk-or-v1-no-pkce-key', got %q", resp.Key)
	}
	if resp.UserID != nil {
		t.Errorf("expected nil user_id, got %v", resp.UserID)
	}
}

func TestExchangeAuthCodeValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Nil request
	_, err := client.ExchangeAuthCode(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}

	// Empty code
	_, err = client.ExchangeAuthCode(context.Background(), &ExchangeAuthCodeRequest{})
	if err == nil {
		t.Fatal("expected error for empty code")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}
}

func TestExchangeAuthCodeErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "bad request",
			statusCode: 400,
			body:       `{"error":{"message":"Invalid authorization code","type":"bad_request"}}`,
		},
		{
			name:       "forbidden",
			statusCode: 403,
			body:       `{"error":{"message":"Access denied","type":"forbidden"}}`,
		},
		{
			name:       "server error",
			statusCode: 500,
			body:       `{"error":{"message":"Internal server error","type":"server_error"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

			_, err := client.ExchangeAuthCode(context.Background(), &ExchangeAuthCodeRequest{
				Code: "test-code",
			})
			if err == nil {
				t.Fatal("expected error")
			}

			reqErr, ok := err.(*RequestError)
			if !ok {
				t.Fatalf("expected RequestError, got %T: %v", err, err)
			}
			if reqErr.StatusCode != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, reqErr.StatusCode)
			}
		})
	}
}

func TestCreateAuthCode(t *testing.T) {
	challenge := "test-challenge-value"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/auth/keys/code" {
			t.Errorf("expected path /auth/keys/code, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		var req CreateAuthCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if req.CallbackURL != "https://myapp.com/callback" {
			t.Errorf("expected callback_url 'https://myapp.com/callback', got %q", req.CallbackURL)
		}
		if req.CodeChallenge == nil || *req.CodeChallenge != challenge {
			t.Errorf("expected code_challenge %q, got %v", challenge, req.CodeChallenge)
		}
		if req.CodeChallengeMethod != CodeChallengeMethodS256 {
			t.Errorf("expected code_challenge_method 'S256', got %q", req.CodeChallengeMethod)
		}

		response := CreateAuthCodeResponse{
			Data: CreateAuthCodeData{
				ID:        "code-abc123",
				AppID:     42,
				CreatedAt: "2026-01-01T00:00:00Z",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateAuthCode(context.Background(), &CreateAuthCodeRequest{
		CallbackURL:         "https://myapp.com/callback",
		CodeChallenge:       &challenge,
		CodeChallengeMethod: CodeChallengeMethodS256,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.ID != "code-abc123" {
		t.Errorf("expected ID 'code-abc123', got %q", resp.Data.ID)
	}
	if resp.Data.AppID != 42 {
		t.Errorf("expected AppID 42, got %f", resp.Data.AppID)
	}
	if resp.Data.CreatedAt != "2026-01-01T00:00:00Z" {
		t.Errorf("expected CreatedAt '2026-01-01T00:00:00Z', got %q", resp.Data.CreatedAt)
	}
}

func TestCreateAuthCodeWithAllOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req CreateAuthCodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if req.CallbackURL != "https://myapp.com/callback" {
			t.Errorf("expected callback_url 'https://myapp.com/callback', got %q", req.CallbackURL)
		}
		if req.Limit == nil || *req.Limit != 10.0 {
			t.Errorf("expected limit 10.0, got %v", req.Limit)
		}
		if req.ExpiresAt == nil || *req.ExpiresAt != "2026-12-31T23:59:59Z" {
			t.Errorf("expected expires_at '2026-12-31T23:59:59Z', got %v", req.ExpiresAt)
		}
		if req.KeyLabel == nil || *req.KeyLabel != "My App Key" {
			t.Errorf("expected key_label 'My App Key', got %v", req.KeyLabel)
		}
		if req.UsageLimitType != UsageLimitMonthly {
			t.Errorf("expected usage_limit_type 'monthly', got %q", req.UsageLimitType)
		}

		response := CreateAuthCodeResponse{
			Data: CreateAuthCodeData{
				ID:        "code-full-opts",
				AppID:     99,
				CreatedAt: "2026-01-01T00:00:00Z",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	challenge := "test-challenge"
	limit := 10.0
	expiresAt := "2026-12-31T23:59:59Z"
	keyLabel := "My App Key"

	resp, err := client.CreateAuthCode(context.Background(), &CreateAuthCodeRequest{
		CallbackURL:         "https://myapp.com/callback",
		CodeChallenge:       &challenge,
		CodeChallengeMethod: CodeChallengeMethodS256,
		Limit:               &limit,
		ExpiresAt:           &expiresAt,
		KeyLabel:            &keyLabel,
		UsageLimitType:      UsageLimitMonthly,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Data.ID != "code-full-opts" {
		t.Errorf("expected ID 'code-full-opts', got %q", resp.Data.ID)
	}
}

func TestCreateAuthCodeValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Nil request
	_, err := client.CreateAuthCode(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}

	// Empty callback_url
	_, err = client.CreateAuthCode(context.Background(), &CreateAuthCodeRequest{})
	if err == nil {
		t.Fatal("expected error for empty callback_url")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T: %v", err, err)
	}
}

func TestCreateAuthCodeErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "bad request",
			statusCode: 400,
			body:       `{"error":{"message":"Invalid callback URL","type":"bad_request"}}`,
		},
		{
			name:       "unauthorized",
			statusCode: 401,
			body:       `{"error":{"message":"Unauthorized","type":"unauthorized"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

			_, err := client.CreateAuthCode(context.Background(), &CreateAuthCodeRequest{
				CallbackURL: "https://myapp.com/callback",
			})
			if err == nil {
				t.Fatal("expected error")
			}

			reqErr, ok := err.(*RequestError)
			if !ok {
				t.Fatalf("expected RequestError, got %T: %v", err, err)
			}
			if reqErr.StatusCode != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, reqErr.StatusCode)
			}
		})
	}
}
