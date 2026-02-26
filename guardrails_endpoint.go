package openrouter

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// buildPaginationQuery builds query parameters for offset and limit pagination.
func buildPaginationQuery(endpoint string, options *ListGuardrailsOptions) string {
	if options == nil {
		return endpoint
	}

	params := url.Values{}

	if options.Offset != nil {
		params.Add("offset", strconv.Itoa(*options.Offset))
	}

	if options.Limit != nil {
		params.Add("limit", strconv.Itoa(*options.Limit))
	}

	if len(params) > 0 {
		return fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	return endpoint
}

// ListGuardrails returns a list of all guardrails for the organization.
// Requires a Provisioning API key.
func (c *Client) ListGuardrails(ctx context.Context, options *ListGuardrailsOptions) (*ListGuardrailsResponse, error) {
	endpoint := buildPaginationQuery("/guardrails", options)

	var response ListGuardrailsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateGuardrail creates a new guardrail with the specified configuration.
// Requires a Provisioning API key.
func (c *Client) CreateGuardrail(ctx context.Context, request *CreateGuardrailRequest) (*Guardrail, error) {
	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	if request.Name == "" {
		return nil, &ValidationError{Message: "name is required"}
	}

	var response Guardrail
	if err := c.doRequest(ctx, "POST", "/guardrails", request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetGuardrail retrieves details about a specific guardrail by its ID.
// Requires a Provisioning API key.
func (c *Client) GetGuardrail(ctx context.Context, id string) (*Guardrail, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s", id)

	var response Guardrail
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateGuardrail updates an existing guardrail by its ID.
// Requires a Provisioning API key.
// All fields in the request are optional - only include the fields you want to update.
func (c *Client) UpdateGuardrail(ctx context.Context, id string, request *UpdateGuardrailRequest) (*Guardrail, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s", id)

	var response Guardrail
	if err := c.doRequest(ctx, "PATCH", endpoint, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteGuardrail deletes a guardrail by its ID.
// Requires a Provisioning API key.
// WARNING: This operation is irreversible!
func (c *Client) DeleteGuardrail(ctx context.Context, id string) (*DeleteGuardrailResponse, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s", id)

	var response DeleteGuardrailResponse
	if err := c.doRequest(ctx, "DELETE", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListAllKeyAssignments returns all key assignments across all guardrails.
// Requires a Provisioning API key.
func (c *Client) ListAllKeyAssignments(ctx context.Context, options *ListGuardrailsOptions) (*ListGuardrailKeyAssignmentsResponse, error) {
	endpoint := buildPaginationQuery("/guardrails/key-assignments", options)

	var response ListGuardrailKeyAssignmentsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListGuardrailKeyAssignments returns key assignments for a specific guardrail.
// Requires a Provisioning API key.
func (c *Client) ListGuardrailKeyAssignments(ctx context.Context, id string, options *ListGuardrailsOptions) (*ListGuardrailKeyAssignmentsResponse, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	endpoint := buildPaginationQuery(fmt.Sprintf("/guardrails/%s/key-assignments", id), options)

	var response ListGuardrailKeyAssignmentsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// AssignKeysToGuardrail assigns API keys to a guardrail by their hashes.
// Requires a Provisioning API key.
func (c *Client) AssignKeysToGuardrail(ctx context.Context, id string, request *AssignKeysRequest) (*AssignKeysResponse, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	if len(request.KeyHashes) == 0 {
		return nil, &ValidationError{Message: "key_hashes cannot be empty"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s/key-assignments", id)

	var response AssignKeysResponse
	if err := c.doRequest(ctx, "POST", endpoint, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UnassignKeysFromGuardrail removes API key assignments from a guardrail.
// Requires a Provisioning API key.
func (c *Client) UnassignKeysFromGuardrail(ctx context.Context, id string, request *AssignKeysRequest) error {
	if id == "" {
		return &ValidationError{Message: "id is required"}
	}

	if request == nil {
		return &ValidationError{Message: "request cannot be nil"}
	}

	if len(request.KeyHashes) == 0 {
		return &ValidationError{Message: "key_hashes cannot be empty"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s/key-assignments", id)

	if err := c.doRequest(ctx, "DELETE", endpoint, request, nil); err != nil {
		return err
	}

	return nil
}

// ListAllMemberAssignments returns all member assignments across all guardrails.
// Requires a Provisioning API key.
func (c *Client) ListAllMemberAssignments(ctx context.Context, options *ListGuardrailsOptions) (*ListGuardrailMemberAssignmentsResponse, error) {
	endpoint := buildPaginationQuery("/guardrails/member-assignments", options)

	var response ListGuardrailMemberAssignmentsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ListGuardrailMemberAssignments returns member assignments for a specific guardrail.
// Requires a Provisioning API key.
func (c *Client) ListGuardrailMemberAssignments(ctx context.Context, id string, options *ListGuardrailsOptions) (*ListGuardrailMemberAssignmentsResponse, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	endpoint := buildPaginationQuery(fmt.Sprintf("/guardrails/%s/member-assignments", id), options)

	var response ListGuardrailMemberAssignmentsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// AssignMembersToGuardrail assigns members to a guardrail by their user IDs.
// Requires a Provisioning API key.
func (c *Client) AssignMembersToGuardrail(ctx context.Context, id string, request *AssignMembersRequest) (*AssignMembersResponse, error) {
	if id == "" {
		return nil, &ValidationError{Message: "id is required"}
	}

	if request == nil {
		return nil, &ValidationError{Message: "request cannot be nil"}
	}

	if len(request.MemberUserIDs) == 0 {
		return nil, &ValidationError{Message: "member_user_ids cannot be empty"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s/member-assignments", id)

	var response AssignMembersResponse
	if err := c.doRequest(ctx, "POST", endpoint, request, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UnassignMembersFromGuardrail removes member assignments from a guardrail.
// Requires a Provisioning API key.
func (c *Client) UnassignMembersFromGuardrail(ctx context.Context, id string, request *AssignMembersRequest) error {
	if id == "" {
		return &ValidationError{Message: "id is required"}
	}

	if request == nil {
		return &ValidationError{Message: "request cannot be nil"}
	}

	if len(request.MemberUserIDs) == 0 {
		return &ValidationError{Message: "member_user_ids cannot be empty"}
	}

	endpoint := fmt.Sprintf("/guardrails/%s/member-assignments", id)

	if err := c.doRequest(ctx, "DELETE", endpoint, request, nil); err != nil {
		return err
	}

	return nil
}
