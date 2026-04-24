package openrouter

// Workspace represents an OpenRouter workspace.
type Workspace struct {
	ID                              string  `json:"id"`
	Name                            string  `json:"name"`
	Slug                            string  `json:"slug"`
	Description                     *string `json:"description"`
	CreatedBy                       *string `json:"created_by"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	DefaultTextModel                *string `json:"default_text_model"`
	DefaultImageModel               *string `json:"default_image_model"`
	DefaultProviderSort             *string `json:"default_provider_sort"`
	IOLoggingAPIKeyIDs              *[]int  `json:"io_logging_api_key_ids"`
	IOLoggingSamplingRate           float64 `json:"io_logging_sampling_rate"`
	IsDataDiscountLoggingEnabled    bool    `json:"is_data_discount_logging_enabled"`
	IsObservabilityBroadcastEnabled bool    `json:"is_observability_broadcast_enabled"`
	IsObservabilityIOLoggingEnabled bool    `json:"is_observability_io_logging_enabled"`
}

// ListWorkspacesOptions represents options for listing workspaces.
type ListWorkspacesOptions struct {
	Offset *int
	Limit  *int
}

// ListWorkspacesResponse is the response from listing workspaces.
type ListWorkspacesResponse struct {
	Data       []Workspace `json:"data"`
	TotalCount int         `json:"total_count"`
}

// CreateWorkspaceRequest is the payload for creating a workspace.
// Name and Slug are required. All other fields are optional.
type CreateWorkspaceRequest struct {
	Name                            string   `json:"name"`
	Slug                            string   `json:"slug"`
	Description                     *string  `json:"description,omitempty"`
	DefaultTextModel                *string  `json:"default_text_model,omitempty"`
	DefaultImageModel               *string  `json:"default_image_model,omitempty"`
	DefaultProviderSort             *string  `json:"default_provider_sort,omitempty"`
	IOLoggingAPIKeyIDs              *[]int   `json:"io_logging_api_key_ids,omitempty"`
	IOLoggingSamplingRate           *float64 `json:"io_logging_sampling_rate,omitempty"`
	IsDataDiscountLoggingEnabled    *bool    `json:"is_data_discount_logging_enabled,omitempty"`
	IsObservabilityBroadcastEnabled *bool    `json:"is_observability_broadcast_enabled,omitempty"`
	IsObservabilityIOLoggingEnabled *bool    `json:"is_observability_io_logging_enabled,omitempty"`
}

// UpdateWorkspaceRequest is the payload for updating a workspace.
// All fields are optional; only provided fields will be updated.
type UpdateWorkspaceRequest struct {
	Name                            *string  `json:"name,omitempty"`
	Slug                            *string  `json:"slug,omitempty"`
	Description                     *string  `json:"description,omitempty"`
	DefaultTextModel                *string  `json:"default_text_model,omitempty"`
	DefaultImageModel               *string  `json:"default_image_model,omitempty"`
	DefaultProviderSort             *string  `json:"default_provider_sort,omitempty"`
	IOLoggingAPIKeyIDs              *[]int   `json:"io_logging_api_key_ids,omitempty"`
	IOLoggingSamplingRate           *float64 `json:"io_logging_sampling_rate,omitempty"`
	IsDataDiscountLoggingEnabled    *bool    `json:"is_data_discount_logging_enabled,omitempty"`
	IsObservabilityBroadcastEnabled *bool    `json:"is_observability_broadcast_enabled,omitempty"`
	IsObservabilityIOLoggingEnabled *bool    `json:"is_observability_io_logging_enabled,omitempty"`
}

// GetWorkspaceResponse is the response from getting a single workspace.
type GetWorkspaceResponse struct {
	Data Workspace `json:"data"`
}

// CreateWorkspaceResponse is the response from creating a workspace.
type CreateWorkspaceResponse struct {
	Data Workspace `json:"data"`
}

// UpdateWorkspaceResponse is the response from updating a workspace.
type UpdateWorkspaceResponse struct {
	Data Workspace `json:"data"`
}

// DeleteWorkspaceResponse is the response from deleting a workspace.
type DeleteWorkspaceResponse struct {
	Deleted bool `json:"deleted"`
}

// WorkspaceMemberRole is the role of a member within a workspace.
type WorkspaceMemberRole string

const (
	WorkspaceMemberRoleAdmin  WorkspaceMemberRole = "admin"
	WorkspaceMemberRoleMember WorkspaceMemberRole = "member"
)

// WorkspaceMember represents a single workspace membership.
type WorkspaceMember struct {
	ID          string              `json:"id"`
	WorkspaceID string              `json:"workspace_id"`
	UserID      string              `json:"user_id"`
	Role        WorkspaceMemberRole `json:"role"`
	CreatedAt   string              `json:"created_at"`
}

// BulkAddWorkspaceMembersResponse is the response from bulk-adding workspace members.
type BulkAddWorkspaceMembersResponse struct {
	AddedCount int               `json:"added_count"`
	Data       []WorkspaceMember `json:"data"`
}

// BulkRemoveWorkspaceMembersResponse is the response from bulk-removing workspace members.
type BulkRemoveWorkspaceMembersResponse struct {
	RemovedCount int `json:"removed_count"`
}

// bulkWorkspaceMembersRequest is the internal request body for add/remove member endpoints.
type bulkWorkspaceMembersRequest struct {
	UserIDs []string `json:"user_ids"`
}
