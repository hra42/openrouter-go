package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable (optional for listing models)
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("Note: OPENROUTER_API_KEY not set. Using unauthenticated access.")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: List all available models
	fmt.Println("=== Example 1: List All Models ===")
	listAllModels(client)

	// Example 2: List models filtered by category
	fmt.Println("\n=== Example 2: List Models by Category (Programming) ===")
	listModelsByCategory(client, "programming")

	// Example 3: Display detailed model information
	fmt.Println("\n=== Example 3: Display Detailed Model Information ===")
	displayDetailedModelInfo(client)
}

func listAllModels(client *openrouter.Client) {
	resp, err := client.ListModels(context.Background(), nil)
	if err != nil {
		log.Printf("Error listing models: %v", err)
		return
	}

	fmt.Printf("Total models available: %d\n\n", len(resp.Data))

	// Display first 10 models
	for i, model := range resp.Data {
		if i >= 10 {
			fmt.Printf("... and %d more models\n", len(resp.Data)-10)
			break
		}
		fmt.Printf("%d. %s (%s)\n", i+1, model.Name, model.ID)
		fmt.Printf("   Description: %s\n", truncate(model.Description, 80))
		if model.ContextLength != nil {
			fmt.Printf("   Context Length: %.0f tokens\n", *model.ContextLength)
		}
		fmt.Printf("   Pricing - Prompt: $%s/M tokens, Completion: $%s/M tokens\n\n",
			model.Pricing.Prompt, model.Pricing.Completion)
	}
}

func listModelsByCategory(client *openrouter.Client, category string) {
	resp, err := client.ListModels(context.Background(), &openrouter.ListModelsOptions{
		Category: category,
	})
	if err != nil {
		log.Printf("Error listing models: %v", err)
		return
	}

	fmt.Printf("Models in category '%s': %d\n\n", category, len(resp.Data))

	// Display first 5 models
	for i, model := range resp.Data {
		if i >= 5 {
			fmt.Printf("... and %d more models in this category\n", len(resp.Data)-5)
			break
		}
		fmt.Printf("%d. %s (%s)\n", i+1, model.Name, model.ID)
		fmt.Printf("   Supported Parameters: %v\n", model.SupportedParameters)
		fmt.Printf("   Moderation: %v\n\n", model.TopProvider.IsModerated)
	}
}

func displayDetailedModelInfo(client *openrouter.Client) {
	resp, err := client.ListModels(context.Background(), nil)
	if err != nil {
		log.Printf("Error listing models: %v", err)
		return
	}

	if len(resp.Data) == 0 {
		fmt.Println("No models available")
		return
	}

	// Display detailed info for the first model
	model := resp.Data[0]
	fmt.Printf("Model: %s\n", model.Name)
	fmt.Printf("ID: %s\n", model.ID)
	fmt.Printf("Description: %s\n", model.Description)
	if model.ContextLength != nil {
		fmt.Printf("Context Length: %.0f tokens\n", *model.ContextLength)
	}
	if model.HuggingFaceID != nil {
		fmt.Printf("Hugging Face ID: %s\n", *model.HuggingFaceID)
	}

	fmt.Printf("\nArchitecture:\n")
	fmt.Printf("  Input Modalities: %v\n", model.Architecture.InputModalities)
	fmt.Printf("  Output Modalities: %v\n", model.Architecture.OutputModalities)
	fmt.Printf("  Tokenizer: %s\n", model.Architecture.Tokenizer)
	if model.Architecture.InstructType != nil {
		fmt.Printf("  Instruct Type: %s\n", *model.Architecture.InstructType)
	}

	fmt.Printf("\nTop Provider:\n")
	if model.TopProvider.ContextLength != nil {
		fmt.Printf("  Context Length: %.0f tokens\n", *model.TopProvider.ContextLength)
	}
	if model.TopProvider.MaxCompletionTokens != nil {
		fmt.Printf("  Max Completion Tokens: %.0f tokens\n", *model.TopProvider.MaxCompletionTokens)
	}
	fmt.Printf("  Is Moderated: %v\n", model.TopProvider.IsModerated)

	if model.DefaultParameters != nil {
		fmt.Printf("\nDefault Parameters:\n")
		if model.DefaultParameters.Temperature != nil {
			fmt.Printf("  Temperature: %.2f\n", *model.DefaultParameters.Temperature)
		}
		if model.DefaultParameters.TopP != nil {
			fmt.Printf("  Top P: %.2f\n", *model.DefaultParameters.TopP)
		}
		if model.DefaultParameters.FrequencyPenalty != nil {
			fmt.Printf("  Frequency Penalty: %.2f\n", *model.DefaultParameters.FrequencyPenalty)
		}
	}

	fmt.Printf("\nPricing:\n")
	fmt.Printf("  Prompt: $%s per million tokens\n", model.Pricing.Prompt)
	fmt.Printf("  Completion: $%s per million tokens\n", model.Pricing.Completion)
	fmt.Printf("  Image: $%s\n", model.Pricing.Image)
	fmt.Printf("  Request: $%s\n", model.Pricing.Request)
	fmt.Printf("  Web Search: $%s\n", model.Pricing.WebSearch)
	fmt.Printf("  Internal Reasoning: $%s\n", model.Pricing.InternalReasoning)
	if model.Pricing.InputCacheRead != nil {
		fmt.Printf("  Input Cache Read: $%s\n", *model.Pricing.InputCacheRead)
	}
	if model.Pricing.InputCacheWrite != nil {
		fmt.Printf("  Input Cache Write: $%s\n", *model.Pricing.InputCacheWrite)
	}

	if len(model.SupportedParameters) > 0 {
		fmt.Printf("\nSupported Parameters: %v\n", model.SupportedParameters)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
