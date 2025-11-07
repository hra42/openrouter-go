package tests

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunAudioBase64Test tests audio input with base64-encoded audio from local file
func RunAudioBase64Test(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Audio Input (Base64 from Local File)\n")

	// Try multiple possible paths for the test audio
	audioPaths := []string{
		"test-mp3.mp3",                     // Run from cmd/openrouter-test
		"cmd/openrouter-test/test-mp3.mp3", // Run from repo root
		"../../test-mp3.mp3",               // Run from nested directory
	}

	var audioPath string
	var err error
	var message openrouter.Message

	// Try each path until we find the audio file
	for _, path := range audioPaths {
		message, err = openrouter.CreateUserMessageWithBase64Audio(
			"Please transcribe this audio or describe what you hear. Keep your response brief.",
			path,
		)
		if err == nil {
			audioPath = path
			break
		}
	}

	if err != nil {
		printError("Failed to encode audio (tried: test-mp3.mp3, cmd/openrouter-test/test-mp3.mp3)", err)
		return false
	}

	if verbose {
		fmt.Printf("   Using audio file: %s\n", audioPath)
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
		// Some models might not support audio
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support audio\n", model)
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

// RunAudioContentBuilderTest tests audio input using ContentBuilder with local file
func RunAudioContentBuilderTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Audio Input with ContentBuilder\n")

	// Try multiple possible paths for the test audio
	audioPaths := []string{
		"test-mp3.mp3",                     // Run from cmd/openrouter-test
		"cmd/openrouter-test/test-mp3.mp3", // Run from repo root
		"../../test-mp3.mp3",               // Run from nested directory
	}

	content := openrouter.NewContentBuilder().
		AddText("I'm providing an audio file for you to transcribe or analyze:")

	var audioPath string
	var err error

	// Try each path until we find the audio file
	for _, path := range audioPaths {
		content, err = content.AddBase64Audio(path)
		if err == nil {
			audioPath = path
			break
		}
	}

	if err != nil {
		printError("Failed to encode audio (tried: test-mp3.mp3, cmd/openrouter-test/test-mp3.mp3)", err)
		return false
	}

	if verbose {
		fmt.Printf("   Using audio file: %s\n", audioPath)
	}

	content.AddText("Please provide a brief transcription or description.")

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
		// Some models might not support audio
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
				fmt.Printf("⚠️  Skipped: Model %s not available or doesn't support audio with ContentBuilder\n", model)
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

// RunAudioFormatsTest tests different audio formats (WAV and MP3)
func RunAudioFormatsTest(ctx context.Context, client *openrouter.Client, model string, maxTokens int, verbose bool) bool {
	fmt.Printf("🔄 Test: Audio Formats (WAV and MP3)\n")

	formats := []struct {
		name   string
		format string
		data   string
	}{
		{
			name:   "WAV",
			format: "wav",
			data:   "UklGRiQAAABXQVZFZm10IBAAAAABAAEAQB8AAAB9AAACABAAZGF0YQAAAAA=", // Minimal WAV header
		},
		{
			name:   "MP3",
			format: "mp3",
			data:   "//uQxAAAAAAAAAAAAAAAAAAAAAAAWGluZwAAAA8AAAACAAADhAC", // Minimal MP3 header
		},
	}

	for _, format := range formats {
		fmt.Printf("   Testing %s format...\n", format.name)

		message := openrouter.CreateUserMessageWithAudio(
			fmt.Sprintf("This is a %s audio file test. Please acknowledge.", format.name),
			format.data,
			format.format,
		)

		messages := []openrouter.Message{message}

		start := time.Now()
		resp, err := client.ChatComplete(ctx, messages,
			openrouter.WithModel(model),
			openrouter.WithMaxTokens(maxTokens),
			openrouter.WithTemperature(0.7),
		)
		elapsed := time.Since(start)

		if err != nil {
			// Some models might not support audio or specific formats
			if reqErr, ok := err.(*openrouter.RequestError); ok {
				if reqErr.IsNotFoundError() || reqErr.StatusCode == 403 || reqErr.StatusCode == 400 {
					fmt.Printf("   ⚠️  Skipped %s: Model doesn't support this format\n", format.name)
					continue // Don't fail the test for unsupported formats
				}
			}
			printError(fmt.Sprintf("Failed for %s", format.name), err)
			return false
		}

		printSuccess(fmt.Sprintf("%s format success! (%.2fs)", format.name, elapsed.Seconds()))

		if verbose {
			response := resp.Choices[0].Message.Content.(string)
			if len(response) > 100 {
				response = response[:100] + "..."
			}
			fmt.Printf("   Response: %s\n", strings.TrimSpace(response))
		}
	}

	return true
}
