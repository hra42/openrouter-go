package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListOrganizationMembers(t *testing.T) {
	first := "Ada"
	last := "Lovelace"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/organization/members" {
			t.Errorf("expected path /organization/members, got %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("expected empty query string, got %q", r.URL.RawQuery)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-key" {
			t.Errorf("expected Authorization header 'Bearer test-key', got %q", auth)
		}

		response := ListOrganizationMembersResponse{
			Data: []OrganizationMember{
				{
					ID:        "user_123",
					Email:     "ada@example.com",
					FirstName: &first,
					LastName:  &last,
					Role:      OrganizationMemberRoleAdmin,
				},
				{
					ID:        "user_456",
					Email:     "anon@example.com",
					FirstName: nil,
					LastName:  nil,
					Role:      OrganizationMemberRoleMember,
				},
			},
			TotalCount: 2,
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListOrganizationMembers(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.TotalCount != 2 {
		t.Errorf("expected TotalCount 2, got %d", resp.TotalCount)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 members, got %d", len(resp.Data))
	}
	if resp.Data[0].ID != "user_123" {
		t.Errorf("expected ID 'user_123', got %q", resp.Data[0].ID)
	}
	if resp.Data[0].Email != "ada@example.com" {
		t.Errorf("expected Email 'ada@example.com', got %q", resp.Data[0].Email)
	}
	if resp.Data[0].FirstName == nil || *resp.Data[0].FirstName != "Ada" {
		t.Errorf("expected FirstName 'Ada', got %v", resp.Data[0].FirstName)
	}
	if resp.Data[0].Role != OrganizationMemberRoleAdmin {
		t.Errorf("expected Role org:admin, got %q", resp.Data[0].Role)
	}
	if resp.Data[1].FirstName != nil {
		t.Errorf("expected FirstName nil, got %v", *resp.Data[1].FirstName)
	}
	if resp.Data[1].LastName != nil {
		t.Errorf("expected LastName nil, got %v", *resp.Data[1].LastName)
	}
	if resp.Data[1].Role != OrganizationMemberRoleMember {
		t.Errorf("expected Role org:member, got %q", resp.Data[1].Role)
	}
}

func TestListOrganizationMembersWithOffsetAndLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("offset") != "5" {
			t.Errorf("expected offset '5', got %q", query.Get("offset"))
		}
		if query.Get("limit") != "10" {
			t.Errorf("expected limit '10', got %q", query.Get("limit"))
		}

		response := ListOrganizationMembersResponse{Data: []OrganizationMember{}, TotalCount: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	offset := 5
	limit := 10
	if _, err := client.ListOrganizationMembers(context.Background(), &ListOrganizationMembersOptions{
		Offset: &offset,
		Limit:  &limit,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListOrganizationMembersOffsetOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if query.Get("offset") != "3" {
			t.Errorf("expected offset '3', got %q", query.Get("offset"))
		}
		if query.Has("limit") {
			t.Errorf("expected no 'limit' param, got %q", query.Get("limit"))
		}

		response := ListOrganizationMembersResponse{Data: []OrganizationMember{}, TotalCount: 0}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	offset := 3
	if _, err := client.ListOrganizationMembers(context.Background(), &ListOrganizationMembersOptions{
		Offset: &offset,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListOrganizationMembersError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"code":401,"message":"unauthorized"}}`))
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("bad-key"), WithBaseURL(server.URL))

	if _, err := client.ListOrganizationMembers(context.Background(), nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}
