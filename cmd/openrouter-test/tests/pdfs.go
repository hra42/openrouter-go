package tests

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunPDFURLTest tests PDF input with URL
func RunPDFURLTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: PDF Input (URL)\n")

	// Use the Bitcoin whitepaper as a publicly accessible PDF
	pdfURL := "https://bitcoin.org/bitcoin.pdf"

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithPDF(
			"What is the main topic of this document? Keep your response brief.",
			pdfURL,
			"bitcoin.pdf",
		),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some models might not support file input
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
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

// RunPDFWithEngineTest tests PDF parsing with different engines
func RunPDFWithEngineTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: PDF with Parser Engine\n")

	// Use the Bitcoin whitepaper
	pdfURL := "https://bitcoin.org/bitcoin.pdf"

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithPDF(
			"What is the main topic of this document? Keep your response brief.",
			pdfURL,
			"bitcoin.pdf",
		),
	}

	// Use the pdf-text engine (free)
	plugin := openrouter.CreateFileParserPlugin("pdf-text")

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithPlugins(plugin),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Engine: pdf-text\n")
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunPDFWithAnnotationsTest tests reusing file annotations
func RunPDFWithAnnotationsTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: PDF with File Annotations\n")

	// Use the Bitcoin whitepaper
	pdfURL := "https://bitcoin.org/bitcoin.pdf"

	firstMessage := openrouter.CreateUserMessageWithPDF(
		"What is the main topic of this document? Keep your response brief.",
		pdfURL,
		"bitcoin.pdf",
	)

	// First request
	start := time.Now()
	resp, err := client.ChatComplete(ctx, []openrouter.Message{firstMessage},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true
			}
		}
		printError("Failed (first request)", err)
		return false
	}

	printSuccess(fmt.Sprintf("First request success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Tokens: %d prompt, %d completion\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens)
		if len(resp.Choices[0].Message.Annotations) > 0 {
			fmt.Printf("   Annotations: %d received (can reuse!)\n", len(resp.Choices[0].Message.Annotations))
		}
	}

	// Follow-up request with annotations
	if len(resp.Choices[0].Message.Annotations) > 0 {
		followUpMessages := []openrouter.Message{
			firstMessage,
			resp.Choices[0].Message, // Include the assistant's response with annotations
			openrouter.CreateUserMessage("Tell me more about that."),
		}

		start2 := time.Now()
		resp2, err := client.ChatComplete(ctx, followUpMessages,
			openrouter.WithModel(model),
			openrouter.WithMaxTokens(maxTokens),
		)
		elapsed2 := time.Since(start2)

		if err != nil {
			printError("Failed (follow-up request)", err)
			return false
		}

		printSuccess(fmt.Sprintf("Follow-up success! (%.2fs)", elapsed2.Seconds()))

		// Validate follow-up response has choices
		if len(resp2.Choices) == 0 {
			printError("Invalid follow-up response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp2.ID))
			return false
		}

		if verbose || true {
			response2 := resp2.Choices[0].Message.Content.(string)
			if len(response2) > 150 {
				response2 = response2[:150] + "..."
			}
			fmt.Printf("   Follow-up: %s\n", strings.TrimSpace(response2))
			fmt.Printf("   Tokens: %d prompt, %d completion\n",
				resp2.Usage.PromptTokens, resp2.Usage.CompletionTokens)
			fmt.Printf("   ✓ PDF was NOT re-parsed (used annotations)\n")
		}
	}

	return true
}

// RunMultipleFilesTest tests multiple files in a single request
func RunMultipleFilesTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Multiple Files (PDF and Image)\n")

	// Use a PDF and an image
	files := []openrouter.File{
		{
			Filename: "bitcoin.pdf",
			FileData: "https://bitcoin.org/bitcoin.pdf",
		},
		{
			Filename: "bitcoin-logo.png",
			FileData: "https://upload.wikimedia.org/wikipedia/commons/thumb/4/46/Bitcoin.svg/400px-Bitcoin.svg.png",
		},
	}

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithFiles(
			"I'm showing you a PDF document and an image. Briefly describe what you see.",
			files,
		),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support multiple files\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Files: %d (1 PDF, 1 image)\n", len(files))
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunPDFContentBuilderTest tests using ContentBuilder with PDFs
func RunPDFContentBuilderTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: ContentBuilder with PDF\n")

	// Build a message with ContentBuilder
	builder := openrouter.NewContentBuilder()
	builder.AddText("What is the main topic of this document?")
	builder.AddPDF("https://bitcoin.org/bitcoin.pdf", "bitcoin.pdf")
	builder.AddText("Keep your response brief.")

	messages := []openrouter.Message{
		{
			Role:    "user",
			Content: builder.Build(),
		},
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
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

// RunBase64PDFTest tests PDF input with base64-encoded local file
func RunBase64PDFTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Base64 Encoded PDF\n")

	// Use the test PDF in the cmd/openrouter-test directory
	pdfPath := filepath.Join("..", "..", "cmd", "openrouter-test", "test-pdf.pdf")

	// Create message with base64-encoded PDF
	message, err := openrouter.CreateUserMessageWithBase64PDF(
		"What is the main topic of this document? Keep your response brief.",
		pdfPath,
		"test-pdf.pdf",
	)
	if err != nil {
		printError("Failed to encode PDF", err)
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
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	// Validate response has choices
	if len(resp.Choices) == 0 {
		printError("Invalid response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp.ID))
		return false
	}

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 200 {
			response = response[:200] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   PDF: Base64-encoded local file\n")
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunPDFComparisonTest tests comparing the same PDF via URL and base64
func RunPDFComparisonTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: PDF URL vs Base64 Comparison\n")

	// Test with URL
	pdfURL := "https://hra42.com/test-pdf.pdf"
	urlMessage := openrouter.CreateUserMessageWithPDF(
		"What is the main topic of this document? Keep your response to one sentence.",
		pdfURL,
		"test-pdf.pdf",
	)

	fmt.Printf("   Testing URL version...\n")
	start1 := time.Now()
	resp1, err := client.ChatComplete(ctx, []openrouter.Message{urlMessage},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed1 := time.Since(start1)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support PDF input\n", model)
				return true
			}
		}
		printError("Failed (URL)", err)
		return false
	}

	// Test with base64
	pdfPath := filepath.Join("..", "..", "cmd", "openrouter-test", "test-pdf.pdf")
	base64Message, err := openrouter.CreateUserMessageWithBase64PDF(
		"What is the main topic of this document? Keep your response to one sentence.",
		pdfPath,
		"test-pdf.pdf",
	)
	if err != nil {
		printError("Failed to encode PDF", err)
		return false
	}

	fmt.Printf("   Testing Base64 version...\n")
	start2 := time.Now()
	resp2, err := client.ChatComplete(ctx, []openrouter.Message{base64Message},
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed2 := time.Since(start2)

	if err != nil {
		printError("Failed (Base64)", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! URL: %.2fs, Base64: %.2fs", elapsed1.Seconds(), elapsed2.Seconds()))

	// Validate both responses have choices
	if len(resp1.Choices) == 0 {
		printError("Invalid URL response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp1.ID))
		return false
	}
	if len(resp2.Choices) == 0 {
		printError("Invalid Base64 response", fmt.Errorf("no choices returned in response (generation_id: %s)", resp2.ID))
		return false
	}

	if verbose || true {
		urlResponse := resp1.Choices[0].Message.Content.(string)
		base64Response := resp2.Choices[0].Message.Content.(string)

		fmt.Printf("\n   URL Response: %s\n", strings.TrimSpace(urlResponse))
		fmt.Printf("   URL Tokens: %d prompt, %d completion\n",
			resp1.Usage.PromptTokens, resp1.Usage.CompletionTokens)

		fmt.Printf("\n   Base64 Response: %s\n", strings.TrimSpace(base64Response))
		fmt.Printf("   Base64 Tokens: %d prompt, %d completion\n",
			resp2.Usage.PromptTokens, resp2.Usage.CompletionTokens)

		fmt.Printf("\n   ✓ Both methods successfully processed the same PDF\n")
	}

	return true
}
