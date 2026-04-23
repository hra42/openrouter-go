# Error handling & retries

All API errors unwrap to `*openrouter.RequestError`:

```go
resp, err := client.ChatComplete(ctx, messages, openrouter.WithModel(m))
if err != nil {
    var rerr *openrouter.RequestError
    if errors.As(err, &rerr) {
        switch {
        case rerr.IsRateLimitError():
            // honor rerr.RetryAfter
        case rerr.IsAuthenticationError():
            // bad key / expired
        case rerr.IsContextLengthError():
            // shrink messages or use WithTransforms("middle-out")
        case rerr.IsModerationError():
            // provider moderation
        case rerr.IsServerError():
            // transient — client already retries
        }
    }
    return err
}
```

## Built-in retry

The client retries transient failures with exponential backoff:

```go
client := openrouter.NewClient(
    openrouter.WithAPIKey(key),
    openrouter.WithRetry(5, 2*time.Second), // max attempts, initial delay
)
```

Don't layer your own retry on top unless you have a specific reason.

## Other error types

- `*ValidationError` — local input validation (missing fields, bad types).
- `*StreamError` — streaming-specific failures.

Use the helper `errors.As` rather than type assertions so wrapped errors work correctly.
