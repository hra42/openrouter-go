---
slug: /
---

# openrouter-go docs

Zero-dependency Go client for the [OpenRouter](https://openrouter.ai) API.

```bash
go get github.com/hra42/openrouter-go
```

## Start here

- [Getting started](./recipes/getting-started.md)
- [All recipes](./recipes/README.md)
- [Godoc](https://pkg.go.dev/github.com/hra42/openrouter-go)

## For AI coding agents

- [`AGENTS.md`](https://github.com/hra42/openrouter-go/blob/main/AGENTS.md) — SDK conventions and pitfalls
- [`llms.txt`](https://github.com/hra42/openrouter-go/blob/main/llms.txt) — task → file index
- [`docs/api-surface.json`](https://github.com/hra42/openrouter-go/blob/main/docs/api-surface.json) — machine-readable public API snapshot

## Status

The chat, completion, streaming, tool-calling, structured-output, multimodal,
embeddings, models, providers, keys, credits, activity, transforms, web-search,
MCP-conversion, and broadcast-webhook APIs are stable.

The [Responses API](./recipes/responses.md) is **beta** and may change without
notice.
