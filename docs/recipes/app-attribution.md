# App attribution

OpenRouter uses two headers to attribute traffic to your app (visible on the leaderboard and in some provider analytics):

```go
client := openrouter.NewClient(
    openrouter.WithAPIKey(key),
    openrouter.WithReferer("https://myapp.com"), // HTTP-Referer
    openrouter.WithAppName("MyApp"),             // X-Title
)
```

Set both in production apps. See [`examples/app-attribution/main.go`](../../examples/app-attribution/main.go).
