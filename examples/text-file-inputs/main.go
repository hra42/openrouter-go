package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))
	ctx := context.Background()

	// Example 1: Single text file
	fmt.Println("Example 1: Analyzing a single code file")
	fmt.Println("----------------------------------------")
	runSingleFileExample(client, ctx)

	fmt.Println()

	// Example 2: Multiple text files
	fmt.Println("Example 2: Analyzing multiple files")
	fmt.Println("----------------------------------------")
	runMultipleFilesExample(client, ctx)

	fmt.Println()

	// Example 3: Using ContentBuilder
	fmt.Println("Example 3: Using ContentBuilder")
	fmt.Println("----------------------------------------")
	runContentBuilderExample(client, ctx)

	fmt.Println()

	// Example 4: Direct text content
	fmt.Println("Example 4: Direct text content")
	fmt.Println("----------------------------------------")
	runDirectContentExample(client, ctx)
}

func runSingleFileExample(client *openrouter.Client, ctx context.Context) {
	fmt.Println("Code pattern for single file:")
	fmt.Println(`
    message, err := openrouter.CreateUserMessageWithTextFile(
        "Review this code:",
        "/path/to/your/code.py",
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
        openrouter.WithModel("anthropic/claude-sonnet-4"),
    )
    `)
}

func runMultipleFilesExample(client *openrouter.Client, ctx context.Context) {
	fmt.Println("Code pattern for multiple files:")
	fmt.Println(`
    message, err := openrouter.CreateUserMessageWithTextFiles(
        "Compare these configuration files and identify differences:",
        "/path/to/config1.yaml",
        "/path/to/config2.yaml",
    )
    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
        openrouter.WithModel("openai/gpt-4"),
    )
    `)
}

func runContentBuilderExample(client *openrouter.Client, ctx context.Context) {
	fmt.Println("Code pattern with ContentBuilder:")
	fmt.Println(`
    builder := openrouter.NewContentBuilder()
    builder.AddText("I have several files to review:")

    builder, err := builder.AddTextFile("/path/to/main.go")
    if err != nil {
        log.Fatal(err)
    }

    builder.AddText("And here's the test file:")
    builder, err = builder.AddTextFile("/path/to/main_test.go")
    if err != nil {
        log.Fatal(err)
    }

    builder.AddText("Do these files work together correctly?")

    message := builder.BuildMessage("user")
    resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
        openrouter.WithModel("anthropic/claude-sonnet-4"),
    )
    `)
}

func runDirectContentExample(client *openrouter.Client, ctx context.Context) {
	// Example with direct content (not from file)
	jsonContent := `{
  "name": "example",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}`

	message := openrouter.CreateUserMessageWithTextContent(
		"Analyze this package.json file:",
		jsonContent,
		"package.json",
	)

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("openai/gpt-4o-mini"),
		openrouter.WithMaxTokens(150),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("Response: %s\n", content)
	}
}
