# OAuth PKCE

For apps that want users to bring their own OpenRouter account. PKCE flow: redirect the user to OpenRouter's authorization page, then exchange the returned code for an API key.

```go
// 1. Create an auth code request (SDK handles the PKCE code challenge).
authReq, _ := client.CreateAuthCode(ctx, &openrouter.CreateAuthCodeRequest{
    CallbackURL: "https://myapp.com/callback",
})
// Redirect the user to authReq's authorization URL; persist the code verifier
// in the user's session so you can use it in step 2.

// 2. Exchange the code returned to your callback for a user API key.
exch, _ := client.ExchangeAuthCode(ctx, &openrouter.ExchangeAuthCodeRequest{
    Code:         codeFromQueryString,
    CodeVerifier: verifier,
})
// exch.Key is the user's API key — store securely.
```

See [`examples/oauth-pkce/main.go`](../../examples/oauth-pkce/main.go) for a runnable HTTP server with the correct field names.
