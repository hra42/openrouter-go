package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunTTSTest tests the Create Speech (TTS) endpoint.
func RunTTSTest(ctx context.Context, client *openrouter.Client, model, voice string, verbose bool) bool {
	fmt.Printf("🔄 Test: TTS (Create Speech)\n")

	input := "Hello, world. This is a test of the OpenRouter text-to-speech endpoint."

	fmt.Printf("   Model: %s\n", model)
	fmt.Printf("   Voice: %s\n", voice)
	fmt.Printf("   Input: %q\n", input)

	start := time.Now()
	resp, err := client.CreateSpeech(ctx, input, model, voice,
		openrouter.WithSpeechResponseFormat(openrouter.SpeechFormatMP3),
	)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to synthesize speech", err)
		return false
	}

	if len(resp.Audio) == 0 {
		printError("No audio bytes returned", nil)
		return false
	}
	if resp.ContentType == "" {
		printError("No Content-Type returned", nil)
		return false
	}

	fmt.Printf("   ✅ Synthesized audio (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Bytes: %d\n", len(resp.Audio))
	fmt.Printf("      Content-Type: %s\n", resp.ContentType)
	fmt.Printf("      Format: %s\n", resp.Format)

	if verbose {
		sample := resp.Audio
		if len(sample) > 16 {
			sample = sample[:16]
		}
		fmt.Printf("      First bytes: % x\n", sample)
	}

	return true
}
