package tests

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunTextFileTest tests text file input
func RunTextFileTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Text File Input\n")

	// Create a temporary test file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_code.py")
	testCode := `def fibonacci(n):
    """Calculate the nth Fibonacci number."""
    if n <= 1:
        return n
    return fibonacci(n-1) + fibonacci(n-2)

print(fibonacci(10))
`

	if err := os.WriteFile(tmpFile, []byte(testCode), 0644); err != nil {
		printError("Failed to create test file", err)
		return false
	}
	defer func() { _ = os.Remove(tmpFile) }()

	message, err := openrouter.CreateUserMessageWithTextFile(
		"Review this Python code. What does it do and is there any issue? Keep your response brief.",
		tmpFile,
	)
	if err != nil {
		printError("Failed to create message", err)
		return false
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("⚠️  Skipped: Model %s not available\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned"))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunMultipleTextFilesTest tests multiple text files in one request
func RunMultipleTextFilesTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Multiple Text Files\n")

	tmpDir := os.TempDir()

	// Create test files
	file1 := filepath.Join(tmpDir, "config.yaml")
	yaml := `database:
  host: localhost
  port: 5432
  name: mydb
`

	file2 := filepath.Join(tmpDir, "handler.go")
	goCode := `package main

import "fmt"

func HandleRequest() {
    fmt.Println("Hello")
}
`

	if err := os.WriteFile(file1, []byte(yaml), 0644); err != nil { panic(err) }
	if err := os.WriteFile(file2, []byte(goCode), 0644); err != nil { panic(err) }
	defer func() { _ = os.Remove(file1) }()
	defer func() { _ = os.Remove(file2) }()

	message, err := openrouter.CreateUserMessageWithTextFiles(
		"Review these configuration and code files. Do they work together correctly? Keep your response brief.",
		file1, file2,
	)
	if err != nil {
		printError("Failed to create message", err)
		return false
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("⚠️  Skipped: Model %s not available\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned"))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Tokens: %d prompt, %d completion\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
	}

	return true
}

// RunTextContentBuilderTest tests ContentBuilder with text files
func RunTextContentBuilderTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: ContentBuilder with Text Files\n")

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "README.md")
	markdown := `# My Project

This is a sample project.

## Features
- Feature 1
- Feature 2
`

	if err := os.WriteFile(tmpFile, []byte(markdown), 0644); err != nil { panic(err) }
	defer func() { _ = os.Remove(tmpFile) }()

	builder := openrouter.NewContentBuilder()
	builder.AddText("I have a documentation file:")
	builder, err := builder.AddTextFile(tmpFile)
	if err != nil {
		printError("Failed to add text file", err)
		return false
	}
	builder.AddText("Can you summarize it? Keep your response brief.")

	message := builder.BuildMessage("user")

	start := time.Now()
	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 {
				fmt.Printf("⚠️  Skipped: Model %s not available\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned"))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
	}

	return true
}

// RunTextFormatValidationTest tests format validation
func RunTextFormatValidationTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Text Format Validation\n")

	tmpDir := os.TempDir()

	// Test supported formats
	supportedFormats := []string{".txt", ".md", ".json", ".py", ".go", ".js", ".yaml", ".toml"}

	for _, ext := range supportedFormats {
		tmpFile := filepath.Join(tmpDir, "test"+ext)
		if err := os.WriteFile(tmpFile, []byte("test content"), 0644); err != nil { panic(err) }

		_, err := openrouter.ReadTextFile(tmpFile)
		_ = os.Remove(tmpFile)

		if err != nil {
			printError(fmt.Sprintf("Format %s should be supported", ext), err)
			return false
		}
	}

	// Test unsupported format
	tmpFile := filepath.Join(tmpDir, "test.exe")
	if err := os.WriteFile(tmpFile, []byte("binary"), 0644); err != nil { panic(err) }
	_, err := openrouter.ReadTextFile(tmpFile)
	_ = os.Remove(tmpFile)

	if err == nil {
		printError("Unsupported format should fail", fmt.Errorf("expected error for .exe file"))
		return false
	}

	printSuccess("Format validation works correctly")
	return true
}
