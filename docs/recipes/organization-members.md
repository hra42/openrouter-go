# Organization Members

List members of the organization that owns the authenticated Management (Provisioning) key. Useful for building admin dashboards or feeding user IDs into [workspace](./workspaces.md) or [guardrail](./guardrails.md) assignments.

```go
client := openrouter.NewClient(openrouter.WithAPIKey(mgmtKey))

resp, err := client.ListOrganizationMembers(ctx, nil)
if err != nil { return err }

fmt.Printf("total=%d\n", resp.TotalCount)
for _, m := range resp.Data {
    fmt.Printf("  %s (%s) — %s\n", m.Email, m.Role, m.ID)
}
```

Paginate with offset/limit:

```go
offset, limit := 0, 50
_, err := client.ListOrganizationMembers(ctx, &openrouter.ListOrganizationMembersOptions{
    Offset: &offset,
    Limit:  &limit,
})
```

Roles returned are `OrganizationMemberRoleAdmin` (`org:admin`) or `OrganizationMemberRoleMember` (`org:member`). Requires a Provisioning key; inference keys get 401.

See [`examples/list-organization-members/main.go`](../../examples/list-organization-members/main.go).
