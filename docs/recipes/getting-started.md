# Getting started

Install:

```bash
go get github.com/hra42/openrouter-go
```

Minimum viable chat completion:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/hra42/openrouter-go"
)

func main() {
    client := openrouter.NewClient(openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")))

    resp, err := client.ChatComplete(context.Background(),
        []openrouter.Message{openrouter.CreateUserMessage("Hello, world")},
        openrouter.WithModel("openai/gpt-4o-mini"),
    )
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(resp.Choices[0].Message.Content)
}
```

## Message helpers

```go
openrouter.CreateSystemMessage("You are a helpful assistant.")
openrouter.CreateUserMessage("Hello!")
openrouter.CreateAssistantMessage("Hi there! How can I help?")
openrouter.CreateToolMessage("Function result", "tool-call-id")
openrouter.CreateMultiModalMessage("user", "What's in this image?", "https://example.com/image.jpg")
```

## Client configuration

```go
client := openrouter.NewClient(
    openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
    openrouter.WithTimeout(60 * time.Second),
    openrouter.WithRetry(3, time.Second),
    openrouter.WithDefaultModel("openai/gpt-4o-mini"),
    openrouter.WithReferer("https://myapp.com"),
    openrouter.WithAppName("MyApp"),
    openrouter.WithHeader("X-Custom", "value"),
)
```

See [`examples/basic/main.go`](../../examples/basic/main.go).
