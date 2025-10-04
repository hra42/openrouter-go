package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable (optional for listing providers)
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("Note: OPENROUTER_API_KEY not set. Using unauthenticated access.")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: List all available providers
	fmt.Println("=== Example 1: List All Providers ===")
	listAllProviders(client)

	// Example 2: Display detailed provider information
	fmt.Println("\n=== Example 2: Display Detailed Provider Information ===")
	displayDetailedProviderInfo(client)

	// Example 3: Find providers with complete policy information
	fmt.Println("\n=== Example 3: Providers with Complete Policy Information ===")
	findProvidersWithPolicies(client)
}

func listAllProviders(client *openrouter.Client) {
	resp, err := client.ListProviders(context.Background())
	if err != nil {
		log.Printf("Error listing providers: %v", err)
		return
	}

	fmt.Printf("Total providers available: %d\n\n", len(resp.Data))

	// Display all providers
	for i, provider := range resp.Data {
		fmt.Printf("%d. %s (%s)\n", i+1, provider.Name, provider.Slug)
		if provider.StatusPageURL != nil {
			fmt.Printf("   Status: %s\n", *provider.StatusPageURL)
		}
	}
}

func displayDetailedProviderInfo(client *openrouter.Client) {
	resp, err := client.ListProviders(context.Background())
	if err != nil {
		log.Printf("Error listing providers: %v", err)
		return
	}

	if len(resp.Data) == 0 {
		fmt.Println("No providers available")
		return
	}

	// Display detailed info for the first few providers
	for i, provider := range resp.Data {
		if i >= 3 {
			fmt.Printf("\n... and %d more providers\n", len(resp.Data)-3)
			break
		}

		fmt.Printf("\nProvider: %s\n", provider.Name)
		fmt.Printf("Slug: %s\n", provider.Slug)

		if provider.PrivacyPolicyURL != nil {
			fmt.Printf("Privacy Policy: %s\n", *provider.PrivacyPolicyURL)
		} else {
			fmt.Printf("Privacy Policy: (not provided)\n")
		}

		if provider.TermsOfServiceURL != nil {
			fmt.Printf("Terms of Service: %s\n", *provider.TermsOfServiceURL)
		} else {
			fmt.Printf("Terms of Service: (not provided)\n")
		}

		if provider.StatusPageURL != nil {
			fmt.Printf("Status Page: %s\n", *provider.StatusPageURL)
		} else {
			fmt.Printf("Status Page: (not provided)\n")
		}
	}
}

func findProvidersWithPolicies(client *openrouter.Client) {
	resp, err := client.ListProviders(context.Background())
	if err != nil {
		log.Printf("Error listing providers: %v", err)
		return
	}

	providersWithPolicies := []openrouter.ProviderInfo{}
	for _, provider := range resp.Data {
		if provider.PrivacyPolicyURL != nil && provider.TermsOfServiceURL != nil {
			providersWithPolicies = append(providersWithPolicies, provider)
		}
	}

	fmt.Printf("Providers with both privacy policy and terms of service: %d/%d\n\n",
		len(providersWithPolicies), len(resp.Data))

	// Display first few providers with complete policies
	for i, provider := range providersWithPolicies {
		if i >= 5 {
			fmt.Printf("... and %d more providers with complete policies\n",
				len(providersWithPolicies)-5)
			break
		}
		fmt.Printf("%d. %s\n", i+1, provider.Name)
		fmt.Printf("   Privacy: %s\n", *provider.PrivacyPolicyURL)
		fmt.Printf("   Terms: %s\n", *provider.TermsOfServiceURL)
		if provider.StatusPageURL != nil {
			fmt.Printf("   Status: %s\n", *provider.StatusPageURL)
		}
		fmt.Println()
	}
}
