# Account & keys

Credits, activity, and API key management.

```go
// Credits remaining
credits, _ := client.GetCredits(ctx)

// Usage analytics
activity, _ := client.GetActivity(ctx)

// Current key metadata
key, _ := client.GetKey(ctx)

// List / create / update / delete keys (requires a provisioning key)
keys, _ := client.ListKeys(ctx)
newKey, _ := client.CreateKey(ctx, openrouter.CreateKeyRequest{Name: "staging"})
```

Examples:

- [`examples/get-credits/`](../../examples/get-credits/main.go)
- [`examples/activity/`](../../examples/activity/main.go)
- [`examples/key/`](../../examples/key/main.go)
- [`examples/list-keys/`](../../examples/list-keys/main.go)
- [`examples/create-key/`](../../examples/create-key/main.go)

See `account_models.go` and `keys_endpoint.go` for full field definitions.
