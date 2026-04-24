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
		log.Fatal("OPENROUTER_API_KEY environment variable is required (must be a Provisioning/Management key)")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	fmt.Println("=== Example 1: List all organization members ===")
	listAll(client)

	fmt.Println("\n=== Example 2: First page (limit 5) ===")
	listPage(client, 0, 5)
}

func listAll(client *openrouter.Client) {
	resp, err := client.ListOrganizationMembers(context.Background(), nil)
	if err != nil {
		log.Printf("Error listing organization members: %v", err)
		return
	}

	fmt.Printf("Total members: %d\n", resp.TotalCount)
	for _, m := range resp.Data {
		fmt.Printf("  - %s  %s  (%s)\n", m.Email, formatName(m), m.Role)
	}
}

func listPage(client *openrouter.Client, offset, limit int) {
	resp, err := client.ListOrganizationMembers(context.Background(), &openrouter.ListOrganizationMembersOptions{
		Offset: &offset,
		Limit:  &limit,
	})
	if err != nil {
		log.Printf("Error listing organization members: %v", err)
		return
	}

	fmt.Printf("Page offset=%d limit=%d (total=%d)\n", offset, limit, resp.TotalCount)
	for _, m := range resp.Data {
		fmt.Printf("  - %s (%s)\n", m.Email, m.Role)
	}
}

func formatName(m openrouter.OrganizationMember) string {
	first, last := "", ""
	if m.FirstName != nil {
		first = *m.FirstName
	}
	if m.LastName != nil {
		last = *m.LastName
	}
	if first == "" && last == "" {
		return "(no name)"
	}
	return fmt.Sprintf("%s %s", first, last)
}
