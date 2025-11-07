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

	// Example 1: Send a PDF via URL
	fmt.Println("Example 1: Analyzing a PDF from URL")
	fmt.Println("----------------------------------------")
	runPDFURLExample(client, ctx)

	fmt.Println()

	// Example 2: Send a local PDF file (base64-encoded)
	fmt.Println("Example 2: Analyzing a local PDF file")
	fmt.Println("----------------------------------------")
	runLocalPDFExample(client, ctx)

	fmt.Println()

	// Example 3: Configure PDF parsing engine
	fmt.Println("Example 3: Using different PDF parsing engines")
	fmt.Println("----------------------------------------")
	runPDFEngineExample(client, ctx)

	fmt.Println()

	// Example 4: Multiple files in one request
	fmt.Println("Example 4: Sending multiple files")
	fmt.Println("----------------------------------------")
	runMultipleFilesExample(client, ctx)

	fmt.Println()

	// Example 5: Reusing file annotations
	fmt.Println("Example 5: Reusing file annotations to save costs")
	fmt.Println("----------------------------------------")
	runFileAnnotationsExample(client, ctx)

	fmt.Println()

	// Example 6: Using ContentBuilder with PDFs
	fmt.Println("Example 6: Using ContentBuilder")
	fmt.Println("----------------------------------------")
	runContentBuilderExample(client, ctx)
}

func runPDFURLExample(client *openrouter.Client, ctx context.Context) {
	// Using the Bitcoin whitepaper as a publicly accessible PDF
	message := openrouter.CreateUserMessageWithPDF(
		"What are the main points in this document?",
		"https://bitcoin.org/bitcoin.pdf",
		"bitcoin.pdf",
	)

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("anthropic/claude-sonnet-4"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("Response: %s\n", content)
	}
}

func runLocalPDFExample(client *openrouter.Client, ctx context.Context) {
	// For this example, we'll show the code but won't actually run it
	// since we don't have a local PDF file in the repository
	fmt.Println("Code example (requires a local PDF file):")
	fmt.Println(`
	message, err := openrouter.CreateUserMessageWithBase64PDF(
		"Summarize this document",
		"/path/to/your/document.pdf",
		"document.pdf",
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("google/gemma-3-27b-it"),
	)
	`)
	fmt.Println("Note: Use EncodePDFToBase64() to encode a PDF file to base64")
}

func runPDFEngineExample(client *openrouter.Client, ctx context.Context) {
	// Configure PDF parsing with a specific engine
	message := openrouter.CreateUserMessageWithPDF(
		"Extract the key concepts from this document",
		"https://bitcoin.org/bitcoin.pdf",
		"bitcoin.pdf",
	)

	// Use the pdf-text engine (free, good for well-structured PDFs)
	pdfPlugin := openrouter.CreateFileParserPlugin("pdf-text")

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("google/gemma-3-27b-it"),
		openrouter.WithPlugins(pdfPlugin),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("Response: %s\n", content)
	}

	fmt.Println("\nAvailable PDF parsing engines:")
	fmt.Println("- 'pdf-text': Free, best for well-structured PDFs")
	fmt.Println("- 'mistral-ocr': $0.0004 per 1K pages, best for scanned PDFs")
	fmt.Println("- 'native': Uses model's native file support (charged as tokens)")
	fmt.Println("- '' (empty): Auto-selects native or pdf-text")
}

func runMultipleFilesExample(client *openrouter.Client, ctx context.Context) {
	// You can send multiple files (PDFs, images, etc.) in one request
	files := []openrouter.File{
		{
			Filename: "bitcoin.pdf",
			FileData: "https://bitcoin.org/bitcoin.pdf",
		},
		// You could add more files here, including images:
		// {
		//     Filename: "chart.png",
		//     FileData: "https://example.com/chart.png",
		// },
	}

	message := openrouter.CreateUserMessageWithFiles(
		"Compare and analyze these documents",
		files,
	)

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("anthropic/claude-sonnet-4"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("Response: %s\n", content)
	}
}

func runFileAnnotationsExample(client *openrouter.Client, ctx context.Context) {
	// First request with PDF
	message := openrouter.CreateUserMessageWithPDF(
		"What are the main concepts in this paper?",
		"https://bitcoin.org/bitcoin.pdf",
		"bitcoin.pdf",
	)

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("google/gemma-3-27b-it"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("First response: %s\n", content)
	}

	// Check if we got file annotations
	if len(resp.Choices[0].Message.Annotations) > 0 {
		fmt.Println("\nReceived file annotations - can reuse in follow-up!")

		// Follow-up request using the annotations
		// This avoids re-parsing the PDF and saves costs
		followUpMessages := []openrouter.Message{
			message,
			resp.Choices[0].Message, // Include the assistant's response with annotations
			openrouter.CreateUserMessage("Can you elaborate on the first point?"),
		}

		followUpResp, err := client.ChatComplete(ctx, followUpMessages,
			openrouter.WithModel("google/gemma-3-27b-it"),
		)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		if content, ok := followUpResp.Choices[0].Message.Content.(string); ok {
			fmt.Printf("\nFollow-up response: %s\n", content)
		}
		fmt.Println("\n✓ PDF was NOT re-parsed - saved processing costs!")
	}
}

func runContentBuilderExample(client *openrouter.Client, ctx context.Context) {
	// Build complex messages with PDFs and other content types
	builder := openrouter.NewContentBuilder()
	builder.AddText("Please analyze this document and the image:")
	builder.AddPDF("https://bitcoin.org/bitcoin.pdf", "bitcoin.pdf")

	// You can also add images to the same message
	builder.AddImage("https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Bitcoin.svg/400px-Bitcoin.svg.png")

	message := openrouter.Message{
		Role:    "user",
		Content: builder.Build(),
	}

	resp, err := client.ChatComplete(ctx, []openrouter.Message{message},
		openrouter.WithModel("anthropic/claude-sonnet-4"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Model: %s\n", resp.Model)
	if content, ok := resp.Choices[0].Message.Content.(string); ok {
		fmt.Printf("Response: %s\n", content)
	}
}
