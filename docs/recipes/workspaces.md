# Workspaces

Workspaces group API keys, default models, and observability settings under a single tenant inside your organization. The workspace endpoints are part of the Management API and require a **Provisioning (Management) API key** — a regular inference key will return 401.

```go
client := openrouter.NewClient(openrouter.WithAPIKey(mgmtKey))

list, err := client.ListWorkspaces(ctx, nil)
if err != nil { return err }
for _, ws := range list.Data {
    fmt.Printf("%s (%s)\n", ws.Name, ws.Slug)
}
```

Create, update, and delete:

```go
desc := "Team sandbox"
created, err := client.CreateWorkspace(ctx, &openrouter.CreateWorkspaceRequest{
    Name:        "Sandbox",
    Slug:        "sandbox",
    Description: &desc,
})

newName := "Sandbox (renamed)"
_, err = client.UpdateWorkspace(ctx, created.Data.ID, &openrouter.UpdateWorkspaceRequest{
    Name: &newName,
})

_, err = client.DeleteWorkspace(ctx, created.Data.ID)
```

`GetWorkspace`, `UpdateWorkspace`, and `DeleteWorkspace` all accept either the workspace UUID or its slug.

## Members

Add or remove organization members in bulk. Members inherit their organization role inside the workspace.

```go
added, err := client.AddWorkspaceMembers(ctx, "sandbox", []string{"user_abc123"})
fmt.Printf("added=%d\n", added.AddedCount)

_, err = client.RemoveWorkspaceMembers(ctx, "sandbox", []string{"user_abc123"})
```

Workspaces with active API keys cannot be deleted, and the default workspace is protected. Both constraints surface as regular `*APIError` values — check `errors.As` and act on the HTTP status.

See [`examples/workspaces/main.go`](../../examples/workspaces/main.go) for an end-to-end flow.
