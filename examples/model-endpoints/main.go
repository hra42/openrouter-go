package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable (optional for listing model endpoints)
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("Note: OPENROUTER_API_KEY not set. Using unauthenticated access.")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: List endpoints for a specific model
	fmt.Println("=== Example 1: List Endpoints for GPT-4 ===")
	listModelEndpoints(client, "openai", "gpt-4")

	// Example 2: List endpoints for Claude
	fmt.Println("\n=== Example 2: List Endpoints for Claude-3.5 Sonnet ===")
	listModelEndpoints(client, "anthropic", "claude-3.5-sonnet")

	// Example 3: Compare endpoints by pricing
	fmt.Println("\n=== Example 3: Compare Endpoints by Pricing ===")
	compareEndpointsByPricing(client, "anthropic", "claude-3.5-sonnet")
}

func listModelEndpoints(client *openrouter.Client, author, slug string) {
	resp, err := client.ListModelEndpoints(context.Background(), author, slug)
	if err != nil {
		log.Printf("Error listing model endpoints: %v", err)
		return
	}

	data := resp.Data
	fmt.Printf("Model: %s (%s)\n", data.Name, data.ID)
	fmt.Printf("Description: %s\n", data.Description)

	fmt.Printf("\nArchitecture:\n")
	if data.Architecture.Tokenizer != nil {
		fmt.Printf("  Tokenizer: %s\n", *data.Architecture.Tokenizer)
	}
	if data.Architecture.InstructType != nil {
		fmt.Printf("  Instruct Type: %s\n", *data.Architecture.InstructType)
	}
	fmt.Printf("  Input Modalities: %s\n", strings.Join(data.Architecture.InputModalities, ", "))
	fmt.Printf("  Output Modalities: %s\n", strings.Join(data.Architecture.OutputModalities, ", "))

	fmt.Printf("\nAvailable Endpoints: %d\n\n", len(data.Endpoints))

	for i, endpoint := range data.Endpoints {
		fmt.Printf("Endpoint %d:\n", i+1)
		fmt.Printf("  Provider: %s\n", endpoint.ProviderName)
		fmt.Printf("  Name: %s\n", endpoint.Name)
		fmt.Printf("  Status: %.0f\n", endpoint.Status)
		fmt.Printf("  Context Length: %.0f tokens\n", endpoint.ContextLength)

		if endpoint.MaxCompletionTokens != nil {
			fmt.Printf("  Max Completion Tokens: %.0f\n", *endpoint.MaxCompletionTokens)
		}
		if endpoint.MaxPromptTokens != nil {
			fmt.Printf("  Max Prompt Tokens: %.0f\n", *endpoint.MaxPromptTokens)
		}

		if endpoint.Quantization != nil && *endpoint.Quantization != "" {
			fmt.Printf("  Quantization: %s\n", *endpoint.Quantization)
		}

		fmt.Printf("  Pricing:\n")
		fmt.Printf("    Prompt: $%s/M tokens\n", endpoint.Pricing.Prompt)
		fmt.Printf("    Completion: $%s/M tokens\n", endpoint.Pricing.Completion)
		fmt.Printf("    Request: $%s\n", endpoint.Pricing.Request)
		fmt.Printf("    Image: $%s\n", endpoint.Pricing.Image)

		if endpoint.UptimeLast30m != nil {
			fmt.Printf("  Uptime (Last 30m): %.2f%%\n", *endpoint.UptimeLast30m*100)
		}

		if len(endpoint.SupportedParameters) > 0 {
			fmt.Printf("  Supported Parameters: %s\n", strings.Join(endpoint.SupportedParameters, ", "))
		}

		fmt.Println()
	}
}

func compareEndpointsByPricing(client *openrouter.Client, author, slug string) {
	resp, err := client.ListModelEndpoints(context.Background(), author, slug)
	if err != nil {
		log.Printf("Error listing model endpoints: %v", err)
		return
	}

	data := resp.Data
	fmt.Printf("Model: %s\n\n", data.Name)
	fmt.Printf("%-30s %-15s %-15s %-15s\n", "Provider", "Prompt/M", "Completion/M", "Request")
	fmt.Printf("%s\n", strings.Repeat("-", 75))

	for _, endpoint := range data.Endpoints {
		fmt.Printf("%-30s $%-14s $%-14s $%-14s\n",
			endpoint.ProviderName,
			endpoint.Pricing.Prompt,
			endpoint.Pricing.Completion,
			endpoint.Pricing.Request,
		)
	}
}
