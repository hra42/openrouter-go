package openrouter

// ResetInterval represents the reset interval for a guardrail budget.
type ResetInterval string

const (
	// ResetIntervalDaily resets the budget daily.
	ResetIntervalDaily ResetInterval = "daily"
	// ResetIntervalWeekly resets the budget weekly.
	ResetIntervalWeekly ResetInterval = "weekly"
	// ResetIntervalMonthly resets the budget monthly.
	ResetIntervalMonthly ResetInterval = "monthly"
)

// Guardrail represents a guardrail configuration for controlling spending,
// model access, and data policies.
type Guardrail struct {
	ID               string         `json:"id"`
	Name             string         `json:"name"`
	Description      *string        `json:"description"`
	LimitUSD         *float64       `json:"limit_usd"`
	ResetInterval    *ResetInterval `json:"reset_interval"`
	AllowedProviders []string       `json:"allowed_providers"`
	AllowedModels    []string       `json:"allowed_models"`
	EnforceZDR       *bool          `json:"enforce_zdr"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        *string        `json:"updated_at"`
}

// ListGuardrailsResponse represents the response from listing guardrails.
type ListGuardrailsResponse struct {
	Data       []Guardrail `json:"data"`
	TotalCount int         `json:"total_count"`
}

// ListGuardrailsOptions represents options for listing guardrails.
type ListGuardrailsOptions struct {
	Offset *int
	Limit  *int
}

// CreateGuardrailRequest represents a request to create a new guardrail.
type CreateGuardrailRequest struct {
	Name             string         `json:"name"`
	Description      *string        `json:"description,omitempty"`
	LimitUSD         *float64       `json:"limit_usd,omitempty"`
	ResetInterval    *ResetInterval `json:"reset_interval,omitempty"`
	AllowedProviders []string       `json:"allowed_providers,omitempty"`
	AllowedModels    []string       `json:"allowed_models,omitempty"`
	EnforceZDR       *bool          `json:"enforce_zdr,omitempty"`
}

// UpdateGuardrailRequest represents a request to update an existing guardrail.
// All fields are optional - only include fields you want to update.
type UpdateGuardrailRequest struct {
	Name             *string        `json:"name,omitempty"`
	Description      *string        `json:"description,omitempty"`
	LimitUSD         *float64       `json:"limit_usd,omitempty"`
	ResetInterval    *ResetInterval `json:"reset_interval,omitempty"`
	AllowedProviders []string       `json:"allowed_providers,omitempty"`
	AllowedModels    []string       `json:"allowed_models,omitempty"`
	EnforceZDR       *bool          `json:"enforce_zdr,omitempty"`
}

// DeleteGuardrailResponse represents the response from deleting a guardrail.
type DeleteGuardrailResponse struct {
	Deleted bool `json:"deleted"`
}

// GuardrailKeyAssignment represents a key assignment to a guardrail.
type GuardrailKeyAssignment struct {
	ID             string  `json:"id"`
	KeyHash        string  `json:"key_hash"`
	OrganizationID string  `json:"organization_id"`
	GuardrailID    string  `json:"guardrail_id"`
	AssignedBy     *string `json:"assigned_by"`
	CreatedAt      string  `json:"created_at"`
}

// ListGuardrailKeyAssignmentsResponse represents the response from listing guardrail key assignments.
type ListGuardrailKeyAssignmentsResponse struct {
	Data       []GuardrailKeyAssignment `json:"data"`
	TotalCount int                      `json:"total_count"`
}

// AssignKeysRequest represents a request to assign keys to a guardrail.
type AssignKeysRequest struct {
	KeyHashes []string `json:"key_hashes"`
}

// AssignKeysResponse represents the response from assigning keys to a guardrail.
type AssignKeysResponse struct {
	AssignedCount int `json:"assigned_count"`
}

// GuardrailMemberAssignment represents a member assignment to a guardrail.
type GuardrailMemberAssignment struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	OrganizationID string  `json:"organization_id"`
	GuardrailID    string  `json:"guardrail_id"`
	AssignedBy     *string `json:"assigned_by"`
	CreatedAt      string  `json:"created_at"`
}

// ListGuardrailMemberAssignmentsResponse represents the response from listing guardrail member assignments.
type ListGuardrailMemberAssignmentsResponse struct {
	Data       []GuardrailMemberAssignment `json:"data"`
	TotalCount int                         `json:"total_count"`
}

// AssignMembersRequest represents a request to assign members to a guardrail.
type AssignMembersRequest struct {
	MemberUserIDs []string `json:"member_user_ids"`
}

// AssignMembersResponse represents the response from assigning members to a guardrail.
type AssignMembersResponse struct {
	AssignedCount int `json:"assigned_count"`
}
