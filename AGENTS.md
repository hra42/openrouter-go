# AGENTS.md

Conventions for AI coding agents (Claude Code, Cursor, etc.) writing code against `github.com/hra42/openrouter-go`. Read this once per session before generating SDK code.

For build, test, and e2e commands see [CLAUDE.md](CLAUDE.md). For a task → file index see [llms.txt](llms.txt).

---

## Before you write code

- **Import**: `import "github.com/hra42/openrouter-go"` — package name is `openrouter`.
- **Go version**: 1.26+. No external dependencies are used or welcome — the SDK is standard-library only.
- **API key**: read from env (`OPENROUTER_API_KEY`). Never hardcode.
- **Context**: every API call takes `context.Context` as the first arg. Plumb it through; don't `context.Background()` inside library-style code.

Minimum viable call:

```go
client := openrouter.NewClient(openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")))
resp, err := client.ChatComplete(ctx,
    []openrouter.Message{openrouter.CreateUserMessage("hi")},
    openrouter.WithModel("openai/gpt-4o-mini"),
)
if err != nil { /* see Errors below */ }
fmt.Println(resp.Choices[0].Message.Content)
```

---

## The one pattern: functional options

Every configurable call uses functional options. **Do not invent config structs.** If you need to set something, look for a `With*` function.

- Client-level: `ClientOption` — applied once at `NewClient(...)`. Examples: `WithAPIKey`, `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithReferer`, `WithAppName`, `WithMaxRetries`, `WithRetryDelay`, `WithCustomHeader`.
- Request-level: distinct option types per endpoint (`ChatCompletionOption`, `CompletionOption`, `ResponsesOption`, …) but same shape. Examples: `WithModel`, `WithMessages`, `WithTemperature`, `WithMaxTokens`, `WithTopP`, `WithStop`, `WithTools`, `WithToolChoice`, `WithResponseFormat`, `WithTransforms`, `WithReasoning`, `WithProvider`, `WithUsage`.

When adding a feature, search `options.go` (and `*_options.go`) first — the option almost certainly already exists.

---

## Streaming contract

Every stream type (`ChatStream`, `CompletionStream`, `ResponsesStream`) follows the same protocol. Violating this leaks the underlying HTTP connection.

```go
stream, err := client.ChatCompleteStream(ctx, messages, openrouter.WithModel(m))
if err != nil { return err }
defer func() { _ = stream.Close() }()   // always

for event := range stream.Events() {
    for _, choice := range event.Choices {
        if choice.Delta != nil {
            if s, ok := choice.Delta.Content.(string); ok {
                fmt.Print(s)
            }
        }
    }
}
if err := stream.Err(); err != nil {     // check AFTER the channel drains
    return err
}
```

Rules:
1. `defer stream.Close()` immediately after the nil-error check.
2. Drain `stream.Events()` to completion OR cancel the context — don't break early without cancel.
3. `stream.Err()` is checked **after** the `Events()` channel closes.
4. Tool calls stream as deltas — accumulate `Choices[0].Delta.ToolCalls` by index before acting.

SSE parsing lives in `internal/sse/` — do not reach into it from user code.

---

## Errors

All API errors unwrap to `*openrouter.RequestError`:

```go
var rerr *openrouter.RequestError
if errors.As(err, &rerr) {
    switch {
    case rerr.IsRateLimitError():      // 429 — respect rerr.RetryAfter
    case rerr.IsAuthenticationError(): // 401/403
    case rerr.IsModerationError():     // provider moderation
    case rerr.IsContextLengthError():  // shrink messages or use WithTransforms
    }
}
```

The client already retries transient failures with exponential backoff (configurable via `WithMaxRetries` / `WithRetryDelay`). Don't layer your own retry on top unless you have a reason.

---

## Which endpoint for which task

| You want to… | Use | Why |
|---|---|---|
| Have a chat-style conversation | `ChatComplete` / `ChatCompleteStream` | Modern, supports tools, structured output, multimodal |
| Use a raw prompt (legacy / base models) | `Complete` / `CompleteStream` | Prompt-in, text-out |
| Use OpenAI's Responses API shape | `CreateResponse` / `CreateResponseStream` | **Beta — unstable**. Prefer ChatComplete for production |
| Talk to Claude via Anthropic shape | `CreateAnthropicMessage` / `Stream` | Drop-in for code already using Anthropic SDK types |
| Get vector embeddings | `CreateEmbedding` / `CreateEmbeddings` | Pair with `Chunk*` helpers in `chunking.go` for long docs |
| Discover models / pricing | `ListModels`, `ListModelEndpoints`, `ListProviders` | Metadata, no inference calls |
| Manage account / keys | `GetCredits`, `GetActivity`, `ListKeys`, `CreateKey`, … | Admin surface |
| Manage workspaces | `ListWorkspaces`, `CreateWorkspace`, `UpdateWorkspace`, `DeleteWorkspace`, `AddWorkspaceMembers`, `RemoveWorkspaceMembers` | Management (Provisioning) key required |
| List organization members | `ListOrganizationMembers` | Management key required |
| Configure spend/model guardrails | `CreateGuardrail`, `ListGuardrails`, `AssignKeysToGuardrail`, `AssignMembersToGuardrail`, … | Management key required. Caps USD spend, restricts models/providers, enforces ZDR |
| Generate a video | `CreateVideo` → `GetVideo` (poll) → `GetVideoContent` | Async: submit, poll until terminal status, then download bytes |
| Synthesize speech | `CreateSpeech` | Returns raw `mp3` or `pcm` bytes |
| Rerank documents | `Rerank` | Second-stage scoring after vector retrieval |
| Parse an OpenRouter Broadcast webhook | `broadcast.ParseTrace*` functions | No `Client` needed — standalone utilities |
| Convert an MCP tool to OpenAI shape | `ConvertMCPTool` | No `Client` needed — standalone utility |
| Complete an OAuth PKCE flow | `ExchangeAuthCode` | Converts a PKCE auth code into an API key |

---

## Multimodal inputs

Messages hold `Content` as either a plain `string` or `[]ContentPart`. Use helpers to build parts:

- Image: `openrouter.NewImageContentPart(url)` or base64 data URL. See `examples/image-inputs/`.
- Audio: `openrouter.NewAudioContentPart(base64Data, format)` where format is `"wav"` or `"mp3"`. Helpers in `audio_utils.go`. See `examples/audio-inputs/`.
- PDF: `openrouter.NewFileContentPart(filename, base64Data)` + `WithPDFEngine(...)` option. See `examples/pdf-inputs/`.
- Text file: same `NewFileContentPart`. See `examples/text-file-inputs/`.

Model must support the modality — check `ListModelEndpoints` or OpenRouter docs before assuming.

---

## Things that will bite you

- **Responses API is beta.** `responses.go` may have breaking changes without notice. Don't put it behind a public stable interface.
- **Web search pricing is not uniform.** `WithWebSearch` can use native provider search or Exa; they price differently. Callers should know which engine they asked for.
- **Tool calling requires a tool-capable model.** If a model doesn't support tools, the API returns an error — don't paper over it with a fallback to `ChatComplete` without tools; surface it.
- **`ChatCompleteStream` returns a `ChatStream`, not a channel.** Don't wrap it in a goroutine that publishes to a chan unless you also propagate `stream.Err()`.
- **`model+":nitro"` / `":floor"` / `":online"` are suffix modifiers.** Concatenate with the base model string — they aren't separate options.
- **Context cancellation interrupts streams cleanly.** Prefer `ctx` cancellation over closing the stream prematurely from a different goroutine.
- **`WithTransforms("middle-out")` silently drops content** when the prompt exceeds the context window. Useful, but don't enable it by default in code that needs determinism.
- **App attribution headers (`WithReferer`, `WithAppName`)** affect OpenRouter's app leaderboard and some provider analytics. Set them in production apps.
- **Two kinds of cost — estimated vs. actual.** `cost.go` (`EstimateCost`, `EstimateCostFromTokens`) computes a *client-side estimate* from model pricing × tokens, useful *before* a call. To get the *actual* cost OpenRouter charged, pass `WithUsage(true)` (chat) or `WithCompletionUsage(true)` (legacy) and read `resp.Usage.Cost` (a `*float64` in credits, nil unless usage accounting was enabled). For streaming, the cost arrives on the final chunk's `Usage`, so accumulate it while ranging over `stream.Events()`.

---

## When adding a new endpoint (SDK contributors)

1. Add an e2e test in `cmd/openrouter-test/main.go` — see the CRITICAL E2E Test Guidelines in [CLAUDE.md](CLAUDE.md). Test functions that hit chat/completion endpoints MUST accept `model string`; never hardcode model names.
2. Follow the functional-options pattern — don't introduce new config shapes.
3. If the endpoint streams, reuse `internal/sse/`.
4. Errors go through `RequestError` — don't invent new error types for the same failure modes.
5. Add an example under `examples/<feature>/main.go` runnable via `go run`.

---

## Reference

- Public godoc: https://pkg.go.dev/github.com/hra42/openrouter-go
- Examples: `examples/` (27 subdirs, each `go run`-able)
- E2E harness: `cmd/openrouter-test/`
- Upstream API docs: https://openrouter.ai/docs
