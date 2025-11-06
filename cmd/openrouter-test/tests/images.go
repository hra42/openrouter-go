package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunImageURLTest tests image input with URL
func RunImageURLTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Image Input (URL)\n")

	// Use the test image hosted on hra42.com (same as test-image.png)
	imageURL := "https://hra42.com/test-image.png"

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImage(
			"What's in this image? Keep your response brief.",
			imageURL,
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
		// Some models might not support vision
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support vision\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunMultipleImagesTest tests multiple images in a single request
func RunMultipleImagesTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Multiple Images\n")

	// Use the test image and a nature scene
	imageURL1 := "https://hra42.com/test-image.png"
	imageURL2 := "https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/1280px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg"

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImages(
			"I'm showing you two images. What do you see in each? Keep your response brief.",
			imageURL1,
			imageURL2,
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
		// Some models might not support vision or multiple images
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support multiple images\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunImageDetailTest tests image input with detail level
func RunImageDetailTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Image with Detail Level\n")

	imageURL := "https://hra42.com/test-image.png"

	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImageDetail(
			"Describe this image in detail. What colors do you see?",
			imageURL,
			"high", // Request high detail analysis
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
		// Some models might not support vision or detail levels
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support detail levels\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunContentBuilderTest tests the ContentBuilder API
func RunContentBuilderTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: ContentBuilder\n")

	imageURL := "https://hra42.com/test-image.png"

	// Build a complex message with interleaved text and images
	content := openrouter.NewContentBuilder().
		AddText("Here's a laboratory image:").
		AddImage(imageURL).
		AddText("What scientific equipment or materials do you see? Keep your response brief.")

	messages := []openrouter.Message{
		content.BuildMessage("user"),
	}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some models might not support vision
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support vision\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
	}

	return true
}

// RunBase64ImageTest tests base64-encoded local image
func RunBase64ImageTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Base64-Encoded Local Image\n")

	// Use the test image in the cmd/openrouter-test directory
	imagePath := "test-image.png"

	// Create message with base64-encoded image
	message, err := openrouter.CreateUserMessageWithBase64Image(
		"What's in this image? Keep your response brief.",
		imagePath,
	)
	if err != nil {
		printError("Failed to encode image", err)
		return false
	}

	messages := []openrouter.Message{message}

	start := time.Now()
	resp, err := client.ChatComplete(ctx, messages,
		openrouter.WithModel(model),
		openrouter.WithMaxTokens(maxTokens),
		openrouter.WithTemperature(0.7),
	)
	elapsed := time.Since(start)

	if err != nil {
		// Some models might not support vision
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support vision\n", model)
				return true // Don't fail the test
			}
		}
		printError("Failed", err)
		return false
	}

	printSuccess(fmt.Sprintf("Success! (%.2fs)", elapsed.Seconds()))

	if verbose || true {
		response := resp.Choices[0].Message.Content.(string)
		if len(response) > 150 {
			response = response[:150] + "..."
		}
		fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		fmt.Printf("   Model: %s\n", resp.Model)
		fmt.Printf("   Tokens: %d prompt, %d completion, %d total\n",
			resp.Usage.PromptTokens, resp.Usage.CompletionTokens, resp.Usage.TotalTokens)
		fmt.Printf("   Image: %s (base64-encoded)\n", imagePath)
	}

	return true
}
