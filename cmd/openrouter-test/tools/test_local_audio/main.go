package main

import (
	"fmt"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Try multiple possible paths for the test audio
	audioPaths := []string{
		"test-mp3.mp3",                     // Run from cmd/openrouter-test
		"cmd/openrouter-test/test-mp3.mp3", // Run from repo root
		"../../test-mp3.mp3",               // Run from nested directory
	}

	var audioPath string
	var base64Audio string
	var format string
	var err error

	// Try each path until we find the audio file
	for _, path := range audioPaths {
		base64Audio, format, err = openrouter.EncodeAudioToBase64(path)
		if err == nil {
			audioPath = path
			break
		}
	}

	if err != nil {
		fmt.Printf("❌ Error: Could not find or encode test audio file\n")
		fmt.Printf("   Tried: %v\n", audioPaths)
		fmt.Printf("   Last error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully encoded audio file!\n")
	fmt.Printf("   Path: %s\n", audioPath)
	fmt.Printf("   Format: %s\n", format)
	fmt.Printf("   Base64 length: %d characters\n", len(base64Audio))
	fmt.Printf("   First 100 chars: %s...\n", base64Audio[:100])

	// Test creating a message
	message, err := openrouter.CreateUserMessageWithBase64Audio(
		"Please transcribe this audio.",
		audioPath,
	)
	if err != nil {
		fmt.Printf("❌ Error creating message: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully created message with audio!\n")
	fmt.Printf("   Message role: %s\n", message.Role)

	parts, ok := message.Content.([]openrouter.ContentPart)
	if !ok {
		fmt.Printf("❌ Error: Message content is not ContentPart array\n")
		os.Exit(1)
	}

	fmt.Printf("   Content parts: %d\n", len(parts))
	for i, part := range parts {
		fmt.Printf("   Part %d type: %s\n", i+1, part.Type)
		if part.Type == "input_audio" && part.InputAudio != nil {
			fmt.Printf("      Audio format: %s\n", part.InputAudio.Format)
			fmt.Printf("      Audio data length: %d\n", len(part.InputAudio.Data))
		}
	}
}
