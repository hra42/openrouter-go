# Embeddings & chunking

```go
resp, err := client.CreateEmbedding(ctx,
    "hello world",
    "openai/text-embedding-3-small",
)
if err != nil { return err }
vec := resp.Data[0].Embedding
```

For batch input, use `CreateEmbeddings` with `[]string`:

```go
resp, _ := client.CreateEmbeddings(ctx,
    []string{"doc one", "doc two"},
    "openai/text-embedding-3-small",
)
```

## Long documents

Use the built-in chunker to split oversized inputs before embedding. See [`examples/embedding-chunking/main.go`](../../examples/embedding-chunking/main.go) — it covers chunk sizing, batching, and stitching results back together.

## Discovering embedding models

```go
models, _ := client.ListEmbeddingsModels(ctx)
```
