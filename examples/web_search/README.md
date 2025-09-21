# Web Search Examples

This directory contains examples of using OpenRouter's web search feature to enhance model responses with real-time web data.

## Features Demonstrated

1. **Basic Web Search** - Using the `:online` model shortcut
2. **Web Plugin with Defaults** - Using the web plugin with automatic settings
3. **Custom Web Plugin** - Customizing search parameters and prompts
4. **Native Search with Context** - Using provider-native search with context size settings
5. **Forcing Exa Search** - Explicitly using Exa search engine
6. **Parsing Annotations** - Extracting URL citations from responses

## Running the Examples

Set your OpenRouter API key:
```bash
export OPENROUTER_API_KEY="your-api-key"
```

Run the examples:
```bash
go run main.go
```

## Web Search Options

### Model Shortcut
Append `:online` to any model to enable web search:
```go
openrouter.WithModel("openai/gpt-4o:online")
```

### Plugin Configuration
Configure the web plugin with custom settings:
```go
plugin := openrouter.Plugin{
    ID:         "web",
    Engine:     "exa",      // "native", "exa", or "" for auto
    MaxResults: 5,          // Number of search results
    SearchPrompt: "...",    // Custom prompt for attaching results
}
```

### Search Context Size
For native search, specify the context size:
```go
openrouter.WithWebSearchOptions(&openrouter.WebSearchOptions{
    SearchContextSize: "high", // "low", "medium", or "high"
})
```

## Pricing

- **Exa Search**: $4 per 1000 results (default 5 results = $0.02 per request)
- **Native Search**: Varies by provider and context size (see documentation)

## Notes

- Native search is available for OpenAI and Anthropic models
- Exa search works with all models
- Search results are included in the response annotations
- Citations should be formatted as markdown links with domain names