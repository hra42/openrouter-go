// Package main demonstrates the Video Generation API endpoints.
//
// This example shows how to:
// - Submit a video generation request
// - Poll the job status until it reaches a terminal state
// - Download the generated video bytes to disk
// - List the available video generation models
//
// Usage:
//
//	export OPENROUTER_API_KEY="your-api-key"
//	go run examples/videos/main.go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	model := flag.String("model", "alibaba/wan-2.6", "Video generation model to use")
	prompt := flag.String("prompt", "A serene mountain landscape at sunset, cinematic", "Text prompt")
	output := flag.String("o", "video.mp4", "Output file path")
	pollInterval := flag.Duration("interval", 5*time.Second, "Polling interval")
	timeout := flag.Duration("timeout", 10*time.Minute, "Overall timeout")
	flag.Parse()

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	fmt.Println("=== Listing available video models ===")
	listModels(ctx, client)

	fmt.Println()
	fmt.Println("=== Submitting video generation job ===")
	job, err := client.CreateVideo(ctx, *model, *prompt,
		openrouter.WithVideoAspectRatio(openrouter.VideoAspectRatio16x9),
		openrouter.WithVideoResolution(openrouter.VideoResolution480p),
	)
	if err != nil {
		log.Fatalf("Failed to submit video: %v", err)
	}
	fmt.Printf("Submitted. Job ID: %s (status: %s)\n", job.ID, job.Status)

	fmt.Println()
	fmt.Println("=== Polling for completion ===")
	completed, err := pollUntilDone(ctx, client, job.ID, *pollInterval)
	if err != nil {
		log.Fatalf("Polling failed: %v", err)
	}

	if completed.Status != openrouter.VideoStatusCompleted {
		log.Fatalf("Job ended with status %q: %s", completed.Status, completed.Error)
	}

	if completed.Usage != nil && completed.Usage.Cost != nil {
		fmt.Printf("Cost: $%.4f (BYOK: %t)\n", *completed.Usage.Cost, completed.Usage.IsBYOK)
	}

	fmt.Println()
	fmt.Println("=== Downloading video content ===")
	content, err := client.GetVideoContent(ctx, job.ID, 0)
	if err != nil {
		log.Fatalf("Failed to download video: %v", err)
	}
	if err := os.WriteFile(*output, content.Content, 0o644); err != nil {
		log.Fatalf("Failed to write %s: %v", *output, err)
	}
	fmt.Printf("Wrote %d bytes to %s (Content-Type: %s)\n", len(content.Content), *output, content.ContentType)
}

func listModels(ctx context.Context, client *openrouter.Client) {
	resp, err := client.ListVideoModels(ctx)
	if err != nil {
		log.Printf("Failed to list models: %v", err)
		return
	}
	fmt.Printf("Found %d video models:\n", len(resp.Data))
	for _, m := range resp.Data {
		fmt.Printf("  - %s (%s)\n", m.ID, m.Name)
	}
}

func pollUntilDone(ctx context.Context, client *openrouter.Client, jobID string, interval time.Duration) (*openrouter.VideoGenerationResponse, error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		resp, err := client.GetVideo(ctx, jobID)
		if err != nil {
			return nil, err
		}

		fmt.Printf("  status: %s\n", resp.Status)

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
