package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: Base64-encoded audio from local file
	fmt.Println("=== Example 1: Base64-encoded Audio from Local File ===")
	base64EncodedAudio(client)

	// Example 2: Using ContentBuilder for complex messages
	fmt.Println("\n=== Example 2: Using ContentBuilder ===")
	contentBuilderExample(client)

	// Example 3: Manual audio encoding
	fmt.Println("\n=== Example 3: Manual Audio Encoding ===")
	manualEncoding(client)
}

func base64EncodedAudio(client *openrouter.Client) {
	// Note: This example requires a local audio file to work
	// Replace "path/to/your/audio.wav" with an actual audio file path
	audioPath := "path/to/your/audio.wav"

	// Check if file exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		log.Printf("Audio file not found: %s", audioPath)
		log.Println("Please update the audioPath variable with a valid audio file path")
		log.Println("Supported formats: wav, mp3")
		return
	}

	// Create a message with a base64-encoded audio file
	message, err := openrouter.CreateUserMessageWithBase64Audio(
		"Please transcribe this audio file.",
		audioPath,
	)
	if err != nil {
		log.Printf("Error encoding audio: %v", err)
		return
	}

	messages := []openrouter.Message{message}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.5-flash"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}

func contentBuilderExample(client *openrouter.Client) {
	// Note: This example requires a local audio file to work
	audioPath := "path/to/your/audio.wav"

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		log.Printf("Audio file not found: %s", audioPath)
		log.Println("Skipping this example - please provide a valid audio file")
		return
	}

	// Use ContentBuilder for more complex message construction
	content := openrouter.NewContentBuilder().
		AddText("I have an audio file to transcribe:")

	// Add the audio file
	content, err := content.AddBase64Audio(audioPath)
	if err != nil {
		log.Printf("Error encoding audio: %v", err)
		return
	}

	content.AddText("Please transcribe the audio and provide a summary.")

	messages := []openrouter.Message{
		content.BuildMessage("user"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.5-flash"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func manualEncoding(client *openrouter.Client) {
	// Example showing manual control over encoding
	audioPath := "path/to/your/audio.mp3"

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		log.Printf("Audio file not found: %s", audioPath)
		log.Println("Skipping this example - please provide a valid audio file")
		return
	}

	// Manually encode the audio file
	base64Audio, format, err := openrouter.EncodeAudioToBase64(audioPath)
	if err != nil {
		log.Printf("Error encoding audio: %v", err)
		return
	}

	fmt.Printf("Detected audio format: %s\n", format)
	fmt.Printf("Base64 encoded audio length: %d characters\n", len(base64Audio))

	// Create a message with the manually encoded audio
	message := openrouter.CreateUserMessageWithAudio(
		"What is said in this audio?",
		base64Audio,
		format,
	)

	messages := []openrouter.Message{message}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.5-flash"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
