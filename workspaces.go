package openrouter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListWorkspaces returns a paginated list of workspaces for the authenticated user.
// Requires a Management (Provisioning) API key.
func (c *Client) ListWorkspaces(ctx context.Context, options *ListWorkspacesOptions) (*ListWorkspacesResponse, error) {
	endpoint := "/workspaces"
	if options != nil {
		params := url.Values{}
		if options.Offset != nil {
			params.Add("offset", strconv.Itoa(*options.Offset))
		}
		if options.Limit != nil {
			params.Add("limit", strconv.Itoa(*options.Limit))
		}
		if len(params) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
		}
	}

	var response ListWorkspacesResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// CreateWorkspace creates a new workspace.
// Requires a Management (Provisioning) API key.
func (c *Client) CreateWorkspace(ctx context.Context, req *CreateWorkspaceRequest) (*CreateWorkspaceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("create workspace request is required")
	}
	var response CreateWorkspaceResponse
	if err := c.doRequest(ctx, "POST", "/workspaces", req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// GetWorkspace retrieves a single workspace by ID (UUID) or slug.
// Requires a Management (Provisioning) API key.
func (c *Client) GetWorkspace(ctx context.Context, idOrSlug string) (*GetWorkspaceResponse, error) {
	endpoint := "/workspaces/" + url.PathEscape(idOrSlug)
	var response GetWorkspaceResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// UpdateWorkspace updates an existing workspace by ID (UUID) or slug.
// Requires a Management (Provisioning) API key.
func (c *Client) UpdateWorkspace(ctx context.Context, idOrSlug string, req *UpdateWorkspaceRequest) (*UpdateWorkspaceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("update workspace request is required")
	}
	endpoint := "/workspaces/" + url.PathEscape(idOrSlug)
	var response UpdateWorkspaceResponse
	if err := c.doRequest(ctx, "PATCH", endpoint, req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteWorkspace deletes a workspace by ID (UUID) or slug.
// The default workspace cannot be deleted. Workspaces with active API keys cannot be deleted.
// Requires a Management (Provisioning) API key.
func (c *Client) DeleteWorkspace(ctx context.Context, idOrSlug string) (*DeleteWorkspaceResponse, error) {
	endpoint := "/workspaces/" + url.PathEscape(idOrSlug)
	var response DeleteWorkspaceResponse
	if err := c.doRequest(ctx, "DELETE", endpoint, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// AddWorkspaceMembers adds multiple organization members to a workspace.
// Members are assigned the same role they hold in the organization.
// Requires a Management (Provisioning) API key.
func (c *Client) AddWorkspaceMembers(ctx context.Context, idOrSlug string, userIDs []string) (*BulkAddWorkspaceMembersResponse, error) {
	endpoint := "/workspaces/" + url.PathEscape(idOrSlug) + "/members/add"
	body := &bulkWorkspaceMembersRequest{UserIDs: userIDs}
	var response BulkAddWorkspaceMembersResponse
	if err := c.doRequest(ctx, "POST", endpoint, body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// RemoveWorkspaceMembers removes multiple members from a workspace.
// Members with active API keys in the workspace cannot be removed.
// Requires a Management (Provisioning) API key.
func (c *Client) RemoveWorkspaceMembers(ctx context.Context, idOrSlug string, userIDs []string) (*BulkRemoveWorkspaceMembersResponse, error) {
	endpoint := "/workspaces/" + url.PathEscape(idOrSlug) + "/members/remove"
	body := &bulkWorkspaceMembersRequest{UserIDs: userIDs}
	var response BulkRemoveWorkspaceMembersResponse
	if err := c.doRequest(ctx, "POST", endpoint, body, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
