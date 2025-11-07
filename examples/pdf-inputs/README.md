# PDF Inputs Example

This example demonstrates how to send PDF files to OpenRouter models using the Go client.

## Features Demonstrated

1. **PDF from URL** - Send publicly accessible PDFs directly
2. **Local PDF files** - Encode and send local PDF files using base64
3. **PDF parsing engines** - Configure different parsing engines for optimal cost/quality
4. **Multiple files** - Send multiple PDFs and files in a single request
5. **File annotations** - Reuse file annotations to avoid re-parsing costs
6. **ContentBuilder** - Build complex messages with PDFs, images, and text

## Setup

Set your OpenRouter API key:

```bash
export OPENROUTER_API_KEY="your-api-key"
```

## Running the Example

```bash
cd examples/pdf-inputs
go run main.go
```

## PDF Processing Engines

OpenRouter provides three PDF processing engines:

1. **`pdf-text`** (Free)
   - Best for well-structured PDFs with clear text content
   - Fast processing
   - Good for digital documents with selectable text

2. **`mistral-ocr`** ($0.0004 per 1,000 pages)
   - Best for scanned documents or PDFs with images
   - Uses OCR to extract text from images
   - Higher quality for complex documents

3. **`native`** (charged as input tokens)
   - Uses the model's native file processing capabilities
   - Only available for models that support file input natively
   - Quality and cost vary by model

If you don't specify an engine, OpenRouter will default to the model's native support (if available), otherwise `pdf-text`.

## Usage Patterns

### Basic PDF from URL

```go
message := openrouter.CreateUserMessageWithPDF(
    "What are the main points in this document?",
    "https://example.com/document.pdf",
    "document.pdf",
)

resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "anthropic/claude-sonnet-4",
    Messages: []openrouter.Message{message},
})
```

### Local PDF File

```go
// Encode the PDF to base64
message, err := openrouter.CreateUserMessageWithBase64PDF(
    "Summarize this document",
    "/path/to/document.pdf",
    "document.pdf",
)

resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "google/gemma-3-27b-it",
    Messages: []openrouter.Message{message},
})
```

### Configure PDF Engine

```go
message := openrouter.CreateUserMessageWithPDF(
    "Extract key concepts",
    "https://example.com/document.pdf",
    "document.pdf",
)

// Use mistral-ocr for scanned documents
plugin := openrouter.CreateFileParserPlugin("mistral-ocr")

resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "google/gemma-3-27b-it",
    Messages: []openrouter.Message{message},
    Plugins:  []openrouter.Plugin{plugin},
})
```

### Multiple Files

```go
files := []openrouter.File{
    {
        Filename: "document1.pdf",
        FileData: "https://example.com/doc1.pdf",
    },
    {
        Filename: "document2.pdf",
        FileData: "https://example.com/doc2.pdf",
    },
}

message := openrouter.CreateUserMessageWithFiles(
    "Compare these documents",
    files,
)
```

### Reuse File Annotations

File annotations allow you to avoid re-parsing the same PDF in follow-up requests:

```go
// First request
resp1, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "google/gemma-3-27b-it",
    Messages: []openrouter.Message{messageWithPDF},
})

// Follow-up request - include the assistant's response with annotations
followUpMessages := []openrouter.Message{
    messageWithPDF,
    resp1.Choices[0].Message, // Contains annotations
    openrouter.CreateUserMessage("Tell me more about that"),
}

resp2, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "google/gemma-3-27b-it",
    Messages: followUpMessages,
})
// PDF is NOT re-parsed - saves processing time and costs!
```

### Using ContentBuilder

For complex messages with multiple content types:

```go
builder := openrouter.NewContentBuilder()
builder.AddText("Analyze this document:")
builder.AddPDF("https://example.com/doc.pdf", "document.pdf")
builder.AddImage("https://example.com/chart.png")
builder.AddText("Focus on the trends shown in the image.")

message := openrouter.Message{
    Role:    "user",
    Content: builder.Build(),
}
```

## Cost Optimization Tips

1. **Use `pdf-text` for digital PDFs** - It's free and works well for most documents
2. **Reuse file annotations** - Include the assistant's response with annotations in follow-up requests
3. **Use `mistral-ocr` only for scanned documents** - It's more expensive but necessary for image-based PDFs
4. **Check model support** - Some models have native file support which may be more cost-effective

## Model Compatibility

PDF input works with **any model** on OpenRouter. However:

- Models with native file support (like Claude, GPT-4, Gemini) process files directly
- Other models use OpenRouter's PDF parsing plugins
- Check the model's documentation for native file support capabilities

## Error Handling

```go
resp, err := client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
    Model:    "anthropic/claude-sonnet-4",
    Messages: []openrouter.Message{messageWithPDF},
})
if err != nil {
    if apiErr, ok := err.(*openrouter.APIError); ok {
        log.Printf("API Error: %s (Code: %s)", apiErr.Message, apiErr.Code)
    } else {
        log.Printf("Error: %v", err)
    }
    return
}
```

## Related Documentation

- [OpenRouter PDF Documentation](https://openrouter.ai/docs/features/multimodal/pdfs)
- [Image Inputs Example](../image-inputs/)
- [API Reference](https://openrouter.ai/docs)
