# Guardrails

Guardrails cap spending, restrict which models/providers a key can call, and can enforce Zero-Data-Retention. They are configured once and then attached to API keys or organization members.

All guardrail endpoints require a **Provisioning (Management) API key**.

## Create and list

```go
client := openrouter.NewClient(openrouter.WithAPIKey(mgmtKey))

limit := 25.0
interval := openrouter.ResetIntervalMonthly
enforceZDR := true

g, err := client.CreateGuardrail(ctx, &openrouter.CreateGuardrailRequest{
    Name:             "contractor-sandbox",
    LimitUSD:         &limit,
    ResetInterval:    &interval,
    AllowedProviders: []string{"openai", "anthropic"},
    AllowedModels:    []string{"openai/gpt-4o-mini"},
    EnforceZDR:       &enforceZDR,
})
if err != nil { return err }

list, _ := client.ListGuardrails(ctx, nil)
for _, g := range list.Data {
    fmt.Printf("%s ($%.2f / %v)\n", g.Name, *g.LimitUSD, *g.ResetInterval)
}
```

Leave any field nil/empty to disable it: an empty `AllowedModels` means "no model restriction", `LimitUSD: nil` means "no spend cap", and so on.

## Update and delete

```go
newLimit := 100.0
_, err = client.UpdateGuardrail(ctx, g.ID, &openrouter.UpdateGuardrailRequest{
    LimitUSD: &newLimit,
})

_, err = client.DeleteGuardrail(ctx, g.ID) // irreversible
```

## Assign to keys and members

```go
_, err := client.AssignKeysToGuardrail(ctx, g.ID, &openrouter.AssignKeysRequest{
    KeyHashes: []string{"sk-or-v1-hash-1", "sk-or-v1-hash-2"},
})

_, err = client.AssignMembersToGuardrail(ctx, g.ID, &openrouter.AssignMembersRequest{
    MemberUserIDs: []string{"user_abc123"},
})
```

`UnassignKeysFromGuardrail` / `UnassignMembersFromGuardrail` take the same request shape. `ListAllKeyAssignments` and `ListAllMemberAssignments` return the full cross-guardrail picture; `ListGuardrailKeyAssignments(id, ...)` and `ListGuardrailMemberAssignments(id, ...)` scope to a single guardrail.

## Reset intervals

| Constant | JSON value |
|---|---|
| `ResetIntervalDaily` | `daily` |
| `ResetIntervalWeekly` | `weekly` |
| `ResetIntervalMonthly` | `monthly` |

`ResetInterval` only matters when `LimitUSD` is set.

See the end-to-end flow in `cmd/openrouter-test/tests/guardrails.go` (run `go run cmd/openrouter-test/main.go -test guardrails`).
