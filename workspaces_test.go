package openrouter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListWorkspaces(t *testing.T) {
	desc := "prod env"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/workspaces" {
			t.Errorf("expected /workspaces, got %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("expected empty query, got %q", r.URL.RawQuery)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Errorf("auth header = %q", auth)
		}

		resp := ListWorkspacesResponse{
			Data: []Workspace{
				{
					ID:                    "ws-1",
					Name:                  "Production",
					Slug:                  "production",
					Description:           &desc,
					CreatedAt:             "2026-01-01T00:00:00Z",
					IOLoggingSamplingRate: 1.0,
				},
			},
			TotalCount: 1,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.ListWorkspaces(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalCount != 1 {
		t.Errorf("TotalCount = %d", resp.TotalCount)
	}
	if len(resp.Data) != 1 || resp.Data[0].Slug != "production" {
		t.Errorf("unexpected data: %+v", resp.Data)
	}
	if resp.Data[0].Description == nil || *resp.Data[0].Description != "prod env" {
		t.Errorf("description = %v", resp.Data[0].Description)
	}
}

func TestListWorkspacesWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("offset") != "5" {
			t.Errorf("offset = %q", q.Get("offset"))
		}
		if q.Get("limit") != "10" {
			t.Errorf("limit = %q", q.Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ListWorkspacesResponse{Data: []Workspace{}, TotalCount: 0})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	offset, limit := 5, 10
	if _, err := client.ListWorkspaces(context.Background(), &ListWorkspacesOptions{
		Offset: &offset, Limit: &limit,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/workspaces" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req CreateWorkspaceRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad body: %v", err)
		}
		if req.Name != "Prod" {
			t.Errorf("name = %q", req.Name)
		}
		if req.Slug != "prod" {
			t.Errorf("slug = %q", req.Slug)
		}
		// Optional fields must be omitted when nil
		if strings.Contains(string(body), "description") {
			t.Errorf("description should be omitted: %s", string(body))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(CreateWorkspaceResponse{
			Data: Workspace{ID: "ws-new", Name: "Prod", Slug: "prod", CreatedAt: "2026-01-01T00:00:00Z"},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.CreateWorkspace(context.Background(), &CreateWorkspaceRequest{
		Name: "Prod",
		Slug: "prod",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.ID != "ws-new" {
		t.Errorf("ID = %q", resp.Data.ID)
	}
}

func TestGetWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/workspaces/production" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(GetWorkspaceResponse{
			Data: Workspace{ID: "ws-1", Name: "Production", Slug: "production", CreatedAt: "2026-01-01T00:00:00Z"},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.GetWorkspace(context.Background(), "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Slug != "production" {
		t.Errorf("slug = %q", resp.Data.Slug)
	}
}

func TestUpdateWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/workspaces/production" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req UpdateWorkspaceRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad body: %v", err)
		}
		if req.Name == nil || *req.Name != "Updated" {
			t.Errorf("name = %v", req.Name)
		}
		if req.Slug != nil {
			t.Errorf("slug should be nil, got %v", *req.Slug)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(UpdateWorkspaceResponse{
			Data: Workspace{ID: "ws-1", Name: "Updated", Slug: "production", CreatedAt: "2026-01-01T00:00:00Z"},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	name := "Updated"
	resp, err := client.UpdateWorkspace(context.Background(), "production", &UpdateWorkspaceRequest{
		Name: &name,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Data.Name != "Updated" {
		t.Errorf("name = %q", resp.Data.Name)
	}
}

func TestDeleteWorkspace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/workspaces/ws-1" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(DeleteWorkspaceResponse{Deleted: true})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.DeleteWorkspace(context.Background(), "ws-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.Deleted {
		t.Errorf("expected Deleted=true")
	}
}

func TestAddWorkspaceMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/workspaces/production/members/add" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var req bulkWorkspaceMembersRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Fatalf("bad body: %v", err)
		}
		if len(req.UserIDs) != 2 || req.UserIDs[0] != "user_a" || req.UserIDs[1] != "user_b" {
			t.Errorf("user_ids = %v", req.UserIDs)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BulkAddWorkspaceMembersResponse{
			AddedCount: 2,
			Data: []WorkspaceMember{
				{ID: "m1", WorkspaceID: "ws-1", UserID: "user_a", Role: WorkspaceMemberRoleAdmin, CreatedAt: "2026-01-01T00:00:00Z"},
				{ID: "m2", WorkspaceID: "ws-1", UserID: "user_b", Role: WorkspaceMemberRoleMember, CreatedAt: "2026-01-01T00:00:00Z"},
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.AddWorkspaceMembers(context.Background(), "production", []string{"user_a", "user_b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.AddedCount != 2 {
		t.Errorf("AddedCount = %d", resp.AddedCount)
	}
	if len(resp.Data) != 2 {
		t.Errorf("len(Data) = %d", len(resp.Data))
	}
	if resp.Data[0].Role != WorkspaceMemberRoleAdmin {
		t.Errorf("role = %q", resp.Data[0].Role)
	}
}

func TestRemoveWorkspaceMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/workspaces/production/members/remove" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(BulkRemoveWorkspaceMembersResponse{RemovedCount: 1})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))
	resp, err := client.RemoveWorkspaceMembers(context.Background(), "production", []string{"user_a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RemovedCount != 1 {
		t.Errorf("RemovedCount = %d", resp.RemovedCount)
	}
}

func TestWorkspacesUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"code":401,"message":"unauthorized"}}`))
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("bad-key"), WithBaseURL(server.URL), WithRetry(0, 0))
	if _, err := client.ListWorkspaces(context.Background(), nil); err == nil {
		t.Fatal("expected error, got nil")
	}
}
