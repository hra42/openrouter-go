package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Create a new client with your API key
	client := openrouter.NewClient(
		openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
	)

	// Example 1: Using the :online shortcut
	fmt.Println("=== Example 1: Using :online shortcut ===")
	basicWebSearch(client)

	// Example 2: Using the web plugin with default settings
	fmt.Println("\n=== Example 2: Using web plugin with defaults ===")
	webPluginDefaults(client)

	// Example 3: Customizing the web plugin
	fmt.Println("\n=== Example 3: Customizing web plugin ===")
	customWebPlugin(client)

	// Example 4: Using native search with context size
	fmt.Println("\n=== Example 4: Native search with context size ===")
	nativeSearchWithContext(client)

	// Example 5: Forcing Exa search
	fmt.Println("\n=== Example 5: Forcing Exa search ===")
	forceExaSearch(client)

	// Example 6: Parsing web search annotations
	fmt.Println("\n=== Example 6: Parsing annotations ===")
	parseAnnotationsExample(client)
}

func basicWebSearch(client *openrouter.Client) {
	ctx := context.Background()

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What are the latest developments in quantum computing as of 2024?"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o:online"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func webPluginDefaults(client *openrouter.Client) {
	ctx := context.Background()

	// Create a web plugin with defaults (auto engine, 5 results)
	webPlugin := openrouter.NewWebPlugin()

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What are the current stock prices for major tech companies?"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(webPlugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func customWebPlugin(client *openrouter.Client) {
	ctx := context.Background()

	// Customize the web plugin
	customPrompt := fmt.Sprintf(
		"Web search conducted on %s. Here are relevant results to consider:",
		time.Now().Format("2006-01-02"),
	)

	webPlugin := openrouter.NewWebPluginWithOptions(
		openrouter.WebSearchEngineAuto,
		3, // Only get 3 results
		customPrompt,
	)

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What's the weather forecast for New York City this week?"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("anthropic/claude-3.5-sonnet"),
		openrouter.WithPlugins(webPlugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func nativeSearchWithContext(client *openrouter.Client) {
	ctx := context.Background()

	// Use native search with high context for detailed research
	webPlugin := openrouter.Plugin{
		ID:     "web",
		Engine: string(openrouter.WebSearchEngineNative),
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Research the latest breakthroughs in renewable energy technology"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(webPlugin),
		openrouter.WithWebSearchOptions(&openrouter.WebSearchOptions{
			SearchContextSize: string(openrouter.WebSearchContextHigh),
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func forceExaSearch(client *openrouter.Client) {
	ctx := context.Background()

	// Force Exa search even for models that support native search
	webPlugin := openrouter.Plugin{
		ID:         "web",
		Engine:     string(openrouter.WebSearchEngineExa),
		MaxResults: 10, // Get more results
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("Find recent academic papers on machine learning optimization"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o"),
		openrouter.WithPlugins(webPlugin),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func parseAnnotationsExample(client *openrouter.Client) {
	ctx := context.Background()

	messages := []openrouter.Message{
		openrouter.CreateUserMessage("What are the latest AI regulations in the EU?"),
	}

	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel("openai/gpt-4o:online"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Parse and display annotations
	message := resp.Choices[0].Message
	fmt.Printf("Response: %s\n\n", message.Content)

	// Extract URL citations from annotations
	citations := openrouter.ParseAnnotations(message.Annotations)
	if len(citations) > 0 {
		fmt.Println("Sources cited:")
		for i, citation := range citations {
			fmt.Printf("%d. %s\n   URL: %s\n   Excerpt: %s\n",
				i+1, citation.Title, citation.URL, citation.Content)
		}
	} else {
		fmt.Println("No citations found in response")
	}
}
