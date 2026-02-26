package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListGuardrails(t *testing.T) {
	desc := "Test guardrail"
	limitUSD := 100.0
	interval := ResetIntervalMonthly
	enforceZDR := true

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails" {
			t.Errorf("expected path /guardrails, got %s", r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		response := ListGuardrailsResponse{
			Data: []Guardrail{
				{
					ID:               "gr_123",
					Name:             "Spending Limit",
					Description:      &desc,
					LimitUSD:         &limitUSD,
					ResetInterval:    &interval,
					AllowedProviders: []string{"openai", "anthropic"},
					AllowedModels:    []string{"openai/gpt-4o"},
					EnforceZDR:       &enforceZDR,
					CreatedAt:        "2024-01-01T00:00:00Z",
				},
			},
			TotalCount: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListGuardrails(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalCount != 1 {
		t.Errorf("expected TotalCount 1, got %d", resp.TotalCount)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 guardrail, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != "gr_123" {
		t.Errorf("expected ID 'gr_123', got %q", resp.Data[0].ID)
	}
	if resp.Data[0].Name != "Spending Limit" {
		t.Errorf("expected Name 'Spending Limit', got %q", resp.Data[0].Name)
	}
	if resp.Data[0].LimitUSD == nil || *resp.Data[0].LimitUSD != 100.0 {
		t.Errorf("expected LimitUSD 100.0, got %v", resp.Data[0].LimitUSD)
	}
	if resp.Data[0].EnforceZDR == nil || *resp.Data[0].EnforceZDR != true {
		t.Errorf("expected EnforceZDR true, got %v", resp.Data[0].EnforceZDR)
	}
	if len(resp.Data[0].AllowedProviders) != 2 {
		t.Errorf("expected 2 AllowedProviders, got %d", len(resp.Data[0].AllowedProviders))
	}
}

func TestListGuardrailsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("offset") != "5" {
			t.Errorf("expected offset '5', got %q", query.Get("offset"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("expected limit '10', got %q", query.Get("limit"))
		}

		response := ListGuardrailsResponse{Data: []Guardrail{}, TotalCount: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	offset := 5
	limit := 10
	_, err := client.ListGuardrails(context.Background(), &ListGuardrailsOptions{
		Offset: &offset,
		Limit:  &limit,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListGuardrailsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ListGuardrailsResponse{Data: []Guardrail{}, TotalCount: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListGuardrails(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected 0 guardrails, got %d", len(resp.Data))
	}
}

func TestListGuardrailsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Provisioning key required",
				Type:    "authentication_error",
				Code:    "invalid_key_type",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("invalid-key"), WithBaseURL(server.URL))

	_, err := client.ListGuardrails(context.Background(), nil)
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
}

func TestCreateGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails" {
			t.Errorf("expected path /guardrails, got %s", r.URL.Path)
		}

		var reqBody CreateGuardrailRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Name != "Test Guardrail" {
			t.Errorf("expected Name 'Test Guardrail', got %q", reqBody.Name)
		}
		if reqBody.LimitUSD == nil || *reqBody.LimitUSD != 50.0 {
			t.Errorf("expected LimitUSD 50.0, got %v", reqBody.LimitUSD)
		}

		response := Guardrail{
			ID:        "gr_new",
			Name:      "Test Guardrail",
			LimitUSD:  reqBody.LimitUSD,
			CreatedAt: "2024-01-10T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	limitUSD := 50.0
	resp, err := client.CreateGuardrail(context.Background(), &CreateGuardrailRequest{
		Name:     "Test Guardrail",
		LimitUSD: &limitUSD,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "gr_new" {
		t.Errorf("expected ID 'gr_new', got %q", resp.ID)
	}
	if resp.Name != "Test Guardrail" {
		t.Errorf("expected Name 'Test Guardrail', got %q", resp.Name)
	}
}

func TestCreateGuardrailMinimal(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody CreateGuardrailRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Name != "Minimal Guardrail" {
			t.Errorf("expected Name 'Minimal Guardrail', got %q", reqBody.Name)
		}
		if reqBody.LimitUSD != nil {
			t.Errorf("expected LimitUSD nil, got %v", reqBody.LimitUSD)
		}

		response := Guardrail{
			ID:        "gr_min",
			Name:      "Minimal Guardrail",
			CreatedAt: "2024-01-10T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateGuardrail(context.Background(), &CreateGuardrailRequest{
		Name: "Minimal Guardrail",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Name != "Minimal Guardrail" {
		t.Errorf("expected Name 'Minimal Guardrail', got %q", resp.Name)
	}
}

func TestCreateGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Test nil request
	_, err := client.CreateGuardrail(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Test empty name
	_, err = client.CreateGuardrail(context.Background(), &CreateGuardrailRequest{
		Name: "",
	})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestCreateGuardrailError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Provisioning key required",
				Type:    "authentication_error",
				Code:    "invalid_key_type",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("invalid-key"), WithBaseURL(server.URL))

	_, err := client.CreateGuardrail(context.Background(), &CreateGuardrailRequest{
		Name: "Test",
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
}

func TestGetGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		expectedPath := "/guardrails/gr_123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		response := Guardrail{
			ID:        "gr_123",
			Name:      "Test Guardrail",
			CreatedAt: "2024-01-01T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.GetGuardrail(context.Background(), "gr_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "gr_123" {
		t.Errorf("expected ID 'gr_123', got %q", resp.ID)
	}
	if resp.Name != "Test Guardrail" {
		t.Errorf("expected Name 'Test Guardrail', got %q", resp.Name)
	}
}

func TestGetGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.GetGuardrail(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestGetGuardrailNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Guardrail not found",
				Type:    "not_found_error",
				Code:    "guardrail_not_found",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.GetGuardrail(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, reqErr.StatusCode)
	}
}

func TestUpdateGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH request, got %s", r.Method)
		}
		expectedPath := "/guardrails/gr_123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		var reqBody UpdateGuardrailRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if reqBody.Name == nil || *reqBody.Name != "Updated Name" {
			t.Errorf("expected Name 'Updated Name', got %v", reqBody.Name)
		}

		updatedAt := "2024-01-10T00:00:00Z"
		response := Guardrail{
			ID:        "gr_123",
			Name:      "Updated Name",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: &updatedAt,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	newName := "Updated Name"
	resp, err := client.UpdateGuardrail(context.Background(), "gr_123", &UpdateGuardrailRequest{
		Name: &newName,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Name != "Updated Name" {
		t.Errorf("expected Name 'Updated Name', got %q", resp.Name)
	}
	if resp.UpdatedAt == nil {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestUpdateGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Test empty id
	newName := "Test"
	_, err := client.UpdateGuardrail(context.Background(), "", &UpdateGuardrailRequest{
		Name: &newName,
	})
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Test nil request
	_, err = client.UpdateGuardrail(context.Background(), "gr_123", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestUpdateGuardrailNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Guardrail not found",
				Type:    "not_found_error",
				Code:    "guardrail_not_found",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	newName := "Test"
	_, err := client.UpdateGuardrail(context.Background(), "nonexistent", &UpdateGuardrailRequest{
		Name: &newName,
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, reqErr.StatusCode)
	}
}

func TestDeleteGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/guardrails/gr_123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		response := DeleteGuardrailResponse{Deleted: true}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.DeleteGuardrail(context.Background(), "gr_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Deleted {
		t.Errorf("expected Deleted true, got %t", resp.Deleted)
	}
}

func TestDeleteGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.DeleteGuardrail(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestDeleteGuardrailNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "Guardrail not found",
				Type:    "not_found_error",
				Code:    "guardrail_not_found",
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.DeleteGuardrail(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	reqErr, ok := err.(*RequestError)
	if !ok {
		t.Fatalf("expected RequestError, got %T", err)
	}
	if reqErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, reqErr.StatusCode)
	}
}

func TestListAllKeyAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/key-assignments" {
			t.Errorf("expected path /guardrails/key-assignments, got %s", r.URL.Path)
		}

		assignedBy := "user_admin"
		response := ListGuardrailKeyAssignmentsResponse{
			Data: []GuardrailKeyAssignment{
				{
					ID:             "ka_1",
					KeyHash:        "hash_abc",
					OrganizationID: "org_1",
					GuardrailID:    "gr_123",
					AssignedBy:     &assignedBy,
					CreatedAt:      "2024-01-01T00:00:00Z",
				},
			},
			TotalCount: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListAllKeyAssignments(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalCount != 1 {
		t.Errorf("expected TotalCount 1, got %d", resp.TotalCount)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(resp.Data))
	}
	if resp.Data[0].KeyHash != "hash_abc" {
		t.Errorf("expected KeyHash 'hash_abc', got %q", resp.Data[0].KeyHash)
	}
}

func TestListGuardrailKeyAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/guardrails/gr_123/key-assignments" {
			t.Errorf("expected path /guardrails/gr_123/key-assignments, got %s", r.URL.Path)
		}

		response := ListGuardrailKeyAssignmentsResponse{
			Data:       []GuardrailKeyAssignment{},
			TotalCount: 0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListGuardrailKeyAssignments(context.Background(), "gr_123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 0 {
		t.Errorf("expected TotalCount 0, got %d", resp.TotalCount)
	}
}

func TestListGuardrailKeyAssignmentsValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.ListGuardrailKeyAssignments(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestAssignKeysToGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/gr_123/key-assignments" {
			t.Errorf("expected path /guardrails/gr_123/key-assignments, got %s", r.URL.Path)
		}

		var reqBody AssignKeysRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if len(reqBody.KeyHashes) != 2 {
			t.Errorf("expected 2 key hashes, got %d", len(reqBody.KeyHashes))
		}

		response := AssignKeysResponse{AssignedCount: 2}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.AssignKeysToGuardrail(context.Background(), "gr_123", &AssignKeysRequest{
		KeyHashes: []string{"hash_1", "hash_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AssignedCount != 2 {
		t.Errorf("expected AssignedCount 2, got %d", resp.AssignedCount)
	}
}

func TestAssignKeysToGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Empty id
	_, err := client.AssignKeysToGuardrail(context.Background(), "", &AssignKeysRequest{
		KeyHashes: []string{"hash_1"},
	})
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Nil request
	_, err = client.AssignKeysToGuardrail(context.Background(), "gr_123", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Empty key hashes
	_, err = client.AssignKeysToGuardrail(context.Background(), "gr_123", &AssignKeysRequest{
		KeyHashes: []string{},
	})
	if err == nil {
		t.Fatal("expected error for empty key hashes, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestUnassignKeysFromGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/gr_123/key-assignments" {
			t.Errorf("expected path /guardrails/gr_123/key-assignments, got %s", r.URL.Path)
		}

		var reqBody AssignKeysRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if len(reqBody.KeyHashes) != 1 {
			t.Errorf("expected 1 key hash, got %d", len(reqBody.KeyHashes))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	err := client.UnassignKeysFromGuardrail(context.Background(), "gr_123", &AssignKeysRequest{
		KeyHashes: []string{"hash_1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnassignKeysFromGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Empty id
	err := client.UnassignKeysFromGuardrail(context.Background(), "", &AssignKeysRequest{
		KeyHashes: []string{"hash_1"},
	})
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Nil request
	err = client.UnassignKeysFromGuardrail(context.Background(), "gr_123", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Empty key hashes
	err = client.UnassignKeysFromGuardrail(context.Background(), "gr_123", &AssignKeysRequest{
		KeyHashes: []string{},
	})
	if err == nil {
		t.Fatal("expected error for empty key hashes, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestListAllMemberAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/member-assignments" {
			t.Errorf("expected path /guardrails/member-assignments, got %s", r.URL.Path)
		}

		assignedBy := "user_admin"
		response := ListGuardrailMemberAssignmentsResponse{
			Data: []GuardrailMemberAssignment{
				{
					ID:             "ma_1",
					UserID:         "user_123",
					OrganizationID: "org_1",
					GuardrailID:    "gr_123",
					AssignedBy:     &assignedBy,
					CreatedAt:      "2024-01-01T00:00:00Z",
				},
			},
			TotalCount: 1,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListAllMemberAssignments(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalCount != 1 {
		t.Errorf("expected TotalCount 1, got %d", resp.TotalCount)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 assignment, got %d", len(resp.Data))
	}
	if resp.Data[0].UserID != "user_123" {
		t.Errorf("expected UserID 'user_123', got %q", resp.Data[0].UserID)
	}
}

func TestListGuardrailMemberAssignments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/guardrails/gr_123/member-assignments" {
			t.Errorf("expected path /guardrails/gr_123/member-assignments, got %s", r.URL.Path)
		}

		response := ListGuardrailMemberAssignmentsResponse{
			Data:       []GuardrailMemberAssignment{},
			TotalCount: 0,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListGuardrailMemberAssignments(context.Background(), "gr_123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 0 {
		t.Errorf("expected TotalCount 0, got %d", resp.TotalCount)
	}
}

func TestListGuardrailMemberAssignmentsValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	_, err := client.ListGuardrailMemberAssignments(context.Background(), "", nil)
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestAssignMembersToGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/gr_123/member-assignments" {
			t.Errorf("expected path /guardrails/gr_123/member-assignments, got %s", r.URL.Path)
		}

		var reqBody AssignMembersRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if len(reqBody.MemberUserIDs) != 2 {
			t.Errorf("expected 2 member user IDs, got %d", len(reqBody.MemberUserIDs))
		}

		response := AssignMembersResponse{AssignedCount: 2}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.AssignMembersToGuardrail(context.Background(), "gr_123", &AssignMembersRequest{
		MemberUserIDs: []string{"user_1", "user_2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AssignedCount != 2 {
		t.Errorf("expected AssignedCount 2, got %d", resp.AssignedCount)
	}
}

func TestAssignMembersToGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Empty id
	_, err := client.AssignMembersToGuardrail(context.Background(), "", &AssignMembersRequest{
		MemberUserIDs: []string{"user_1"},
	})
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Nil request
	_, err = client.AssignMembersToGuardrail(context.Background(), "gr_123", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Empty member user IDs
	_, err = client.AssignMembersToGuardrail(context.Background(), "gr_123", &AssignMembersRequest{
		MemberUserIDs: []string{},
	})
	if err == nil {
		t.Fatal("expected error for empty member user IDs, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}

func TestUnassignMembersFromGuardrail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/guardrails/gr_123/member-assignments" {
			t.Errorf("expected path /guardrails/gr_123/member-assignments, got %s", r.URL.Path)
		}

		var reqBody AssignMembersRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		if len(reqBody.MemberUserIDs) != 1 {
			t.Errorf("expected 1 member user ID, got %d", len(reqBody.MemberUserIDs))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	err := client.UnassignMembersFromGuardrail(context.Background(), "gr_123", &AssignMembersRequest{
		MemberUserIDs: []string{"user_1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnassignMembersFromGuardrailValidation(t *testing.T) {
	client := NewClient(WithAPIKey("test-key"))

	// Empty id
	err := client.UnassignMembersFromGuardrail(context.Background(), "", &AssignMembersRequest{
		MemberUserIDs: []string{"user_1"},
	})
	if err == nil {
		t.Fatal("expected error for empty id, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Nil request
	err = client.UnassignMembersFromGuardrail(context.Background(), "gr_123", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}

	// Empty member user IDs
	err = client.UnassignMembersFromGuardrail(context.Background(), "gr_123", &AssignMembersRequest{
		MemberUserIDs: []string{},
	})
	if err == nil {
		t.Fatal("expected error for empty member user IDs, got nil")
	}
	if _, ok := IsValidationError(err); !ok {
		t.Errorf("expected ValidationError, got %T", err)
	}
}
