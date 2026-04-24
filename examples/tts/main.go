// Package main demonstrates the Create Speech (TTS) API endpoint.
//
// This example shows how to:
// - Synthesize audio from text using OpenRouter's TTS endpoint
// - Write the returned bytes to disk as an audio file
// - Control the output format and speed via functional options
//
// Usage:
//
//	export OPENROUTER_API_KEY="your-api-key"
//	go run examples/tts/main.go
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

	client := openrouter.NewClient(
		openrouter.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	fmt.Println("=== Example 1: Basic TTS (pcm) ===")
	basicSpeech(ctx, client)

	fmt.Println()
	fmt.Println("=== Example 2: MP3 output at 1.25x speed ===")
	mp3Speech(ctx, client)
}

func basicSpeech(ctx context.Context, client *openrouter.Client) {
	resp, err := client.CreateSpeech(ctx,
		"Hello from OpenRouter's Go SDK.",
		"hexgrad/kokoro-82m",
		"af_bella",
	)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	out := "speech.pcm"
	if err := os.WriteFile(out, resp.Audio, 0o644); err != nil {
		log.Fatalf("Failed to write %s: %v", out, err)
	}

	fmt.Printf("Wrote %d bytes to %s (Content-Type: %s, Format: %s)\n",
		len(resp.Audio), out, resp.ContentType, resp.Format)
}

func mp3Speech(ctx context.Context, client *openrouter.Client) {
	resp, err := client.CreateSpeech(ctx,
		"This clip is in MP3 format, slightly sped up.",
		"hexgrad/kokoro-82m",
		"af_bella",
		openrouter.WithSpeechResponseFormat(openrouter.SpeechFormatMP3),
		openrouter.WithSpeechSpeed(1.25),
	)
	if err != nil {
		log.Fatalf("Failed to synthesize speech: %v", err)
	}

	out := "speech.mp3"
	if err := os.WriteFile(out, resp.Audio, 0o644); err != nil {
		log.Fatalf("Failed to write %s: %v", out, err)
	}

	fmt.Printf("Wrote %d bytes to %s (Content-Type: %s, Format: %s)\n",
		len(resp.Audio), out, resp.ContentType, resp.Format)
}
