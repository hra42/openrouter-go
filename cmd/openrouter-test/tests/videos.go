package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunVideoModelsTest tests the List Video Models endpoint.
func RunVideoModelsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Video Models\n")

	start := time.Now()
	resp, err := client.ListVideoModels(ctx)
	elapsed := time.Since(start)

	if err != nil {
		printError("Failed to list video models", err)
		return false
	}

	if len(resp.Data) == 0 {
		printError("No video models returned", nil)
		return false
	}

	fmt.Printf("   ✅ Received %d video models (%.2fs)\n", len(resp.Data), elapsed.Seconds())
	if verbose {
		for _, m := range resp.Data {
			fmt.Printf("      - %s (%s)\n", m.ID, m.Name)
			fmt.Printf("        aspect ratios: %v\n", m.SupportedAspectRatios)
			fmt.Printf("        resolutions:   %v\n", m.SupportedResolutions)
			fmt.Printf("        durations:     %v\n", m.SupportedDurations)
		}
	}

	return true
}

// RunVideoGenerationTest tests the full video generation flow: submit, poll, download.
// Video generation can take several minutes; a 15-minute timeout is applied.
func RunVideoGenerationTest(ctx context.Context, client *openrouter.Client, model string, verbose bool) bool {
	fmt.Printf("🔄 Test: Video Generation (full flow)\n")
	fmt.Printf("   Model: %s\n", model)

	prompt := "A serene mountain landscape at sunset, cinematic"
	fmt.Printf("   Prompt: %q\n", prompt)

	// Submit
	submitStart := time.Now()
	job, err := client.CreateVideo(ctx, model, prompt,
		openrouter.WithVideoAspectRatio(openrouter.VideoAspectRatio16x9),
		openrouter.WithVideoResolution(openrouter.VideoResolution480p),
	)
	if err != nil {
		printError("Failed to submit video generation", err)
		return false
	}
	fmt.Printf("   ✅ Submitted in %.2fs (id=%s, status=%s)\n", time.Since(submitStart).Seconds(), job.ID, job.Status)

	// Poll
	pollCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	final, err := pollVideoUntilDone(pollCtx, client, job.ID, 5*time.Second, verbose)
	if err != nil {
		printError("Polling failed", err)
		return false
	}

	if final.Status != openrouter.VideoStatusCompleted {
		printError(fmt.Sprintf("Job did not complete (status=%s)", final.Status), nil)
		if final.Error != "" {
			fmt.Printf("      Error: %s\n", final.Error)
		}
		return false
	}

	fmt.Printf("   ✅ Job completed\n")
	if final.GenerationID != "" {
		fmt.Printf("      Generation ID: %s\n", final.GenerationID)
	}
	if final.Usage != nil && final.Usage.Cost != nil {
		fmt.Printf("      Cost: $%.4f (BYOK: %t)\n", *final.Usage.Cost, final.Usage.IsBYOK)
	}

	// Download
	downloadStart := time.Now()
	content, err := client.GetVideoContent(pollCtx, job.ID, 0)
	if err != nil {
		printError("Failed to download video content", err)
		return false
	}
	if len(content.Content) == 0 {
		printError("No video bytes returned", nil)
		return false
	}

	fmt.Printf("   ✅ Downloaded %d bytes in %.2fs (Content-Type: %s)\n",
		len(content.Content), time.Since(downloadStart).Seconds(), content.ContentType)

	return true
}

func pollVideoUntilDone(ctx context.Context, client *openrouter.Client, jobID string, interval time.Duration, verbose bool) (*openrouter.VideoGenerationResponse, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		resp, err := client.GetVideo(ctx, jobID)
		if err != nil {
			return nil, err
		}

		printVerbose(verbose, "poll status: %s", resp.Status)

		switch resp.Status {
		case openrouter.VideoStatusCompleted,
			openrouter.VideoStatusFailed,
			openrouter.VideoStatusCancelled,
			openrouter.VideoStatusExpired:
			return resp, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}
	}
}
