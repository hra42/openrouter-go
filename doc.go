// Package openrouter is a zero-dependency Go client for the OpenRouter API
// (https://openrouter.ai). It provides complete bindings for chat completions,
// legacy completions, streaming via Server-Sent Events, tool calling,
// structured outputs, multimodal inputs (image, audio, PDF, text file),
// embeddings, the Responses API (beta), the Anthropic-compatible Messages
// endpoint, broadcast webhook parsing, and the full account/admin surface
// (models, providers, keys, activity, credits, guardrails, ZDR).
//
// The entire package uses only the Go standard library.
//
// # Quick start
//
//	client := openrouter.NewClient(openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")))
//
//	resp, err := client.ChatComplete(ctx,
//	    []openrouter.Message{openrouter.CreateUserMessage("Hello")},
//	    openrouter.WithModel("openai/gpt-4o-mini"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(resp.Choices[0].Message.Content)
//
// # Design
//
// The package follows a small set of conventions used uniformly across every
// endpoint:
//
//   - Functional options. Both [NewClient] and every request method take a
//     variadic list of options (for example [WithAPIKey], [WithModel],
//     [WithTools], [WithResponseFormat]). New behavior is added via new
//     option functions rather than by growing config structs.
//   - context.Context on every call, for cancellation and deadlines.
//   - Streaming via typed iterators ([ChatStream], [CompletionStream],
//     [ResponsesStream]) that must be closed by the caller. The canonical
//     loop is:
//
//	stream, err := client.ChatCompleteStream(ctx, msgs, openrouter.WithModel(m))
//	if err != nil { return err }
//	defer stream.Close()
//	for event := range stream.Events() {
//	    // accumulate event.Choices[0].Delta.Content, etc.
//	}
//	if err := stream.Err(); err != nil { return err }
//
//   - Errors unwrap to [*RequestError], which exposes helpers such as
//     IsRateLimitError, IsAuthenticationError, IsContextLengthError, and
//     IsModerationError for structured handling.
//   - Automatic retry with exponential backoff on transient failures,
//     configurable via [WithMaxRetries] and [WithRetryDelay].
//   - Thread-safe: a single [Client] is safe for concurrent use by multiple
//     goroutines.
//
// # Stability
//
// The chat, completion, streaming, tool-calling, structured-output,
// multimodal, embeddings, models, providers, keys, credits, activity,
// transforms, web-search, MCP-conversion, and broadcast-webhook APIs are
// stable.
//
// The Responses API (see [Client.CreateResponse] and [Client.CreateResponseStream])
// is in beta and may have breaking changes; avoid using it in production
// workloads.
//
// # Further reading
//
// Runnable examples live under examples/ in the repository. Agent-oriented
// conventions live in AGENTS.md; build and test commands live in CLAUDE.md;
// a task-to-file index lives in llms.txt.
package openrouter
