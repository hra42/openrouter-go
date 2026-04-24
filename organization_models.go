package openrouter

// ListOrganizationMembersOptions represents options for listing organization members.
type ListOrganizationMembersOptions struct {
	Offset *int
	Limit  *int
}

// OrganizationMemberRole is the role of a member within the organization.
type OrganizationMemberRole string

const (
	OrganizationMemberRoleAdmin  OrganizationMemberRole = "org:admin"
	OrganizationMemberRoleMember OrganizationMemberRole = "org:member"
)

// OrganizationMember represents a single member of an organization.
type OrganizationMember struct {
	ID        string                 `json:"id"`
	Email     string                 `json:"email"`
	FirstName *string                `json:"first_name"`
	LastName  *string                `json:"last_name"`
	Role      OrganizationMemberRole `json:"role"`
}

// ListOrganizationMembersResponse is the response from listing organization members.
type ListOrganizationMembersResponse struct {
	Data       []OrganizationMember `json:"data"`
	TotalCount int                  `json:"total_count"`
}
