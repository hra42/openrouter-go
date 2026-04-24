package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hra42/openrouter-go"
)

func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required (must be a Management/Provisioning key)")
	}

	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))
	ctx := context.Background()

	fmt.Println("=== 1) List workspaces ===")
	listResp, err := client.ListWorkspaces(ctx, nil)
	if err != nil {
		log.Fatalf("list workspaces: %v", err)
	}
	fmt.Printf("Total: %d\n", listResp.TotalCount)
	for _, ws := range listResp.Data {
		fmt.Printf("  - %s (%s) id=%s\n", ws.Name, ws.Slug, ws.ID)
	}

	fmt.Println("\n=== 2) Create workspace ===")
	slug := fmt.Sprintf("example-%d", time.Now().Unix())
	desc := "Workspace created by the openrouter-go example"
	created, err := client.CreateWorkspace(ctx, &openrouter.CreateWorkspaceRequest{
		Name:        "Example Workspace",
		Slug:        slug,
		Description: &desc,
	})
	if err != nil {
		log.Fatalf("create workspace: %v", err)
	}
	fmt.Printf("Created id=%s slug=%s\n", created.Data.ID, created.Data.Slug)

	fmt.Println("\n=== 3) Get workspace by slug ===")
	got, err := client.GetWorkspace(ctx, slug)
	if err != nil {
		log.Fatalf("get workspace: %v", err)
	}
	fmt.Printf("Got: %s (created_at=%s)\n", got.Data.Name, got.Data.CreatedAt)

	fmt.Println("\n=== 4) Update workspace ===")
	newName := "Example Workspace (updated)"
	updated, err := client.UpdateWorkspace(ctx, created.Data.ID, &openrouter.UpdateWorkspaceRequest{
		Name: &newName,
	})
	if err != nil {
		log.Fatalf("update workspace: %v", err)
	}
	fmt.Printf("Updated name: %s\n", updated.Data.Name)

	// Example (disabled): add members. Replace user IDs with real Clerk user IDs.
	//
	// added, err := client.AddWorkspaceMembers(ctx, created.Data.ID, []string{"user_abc123"})
	// if err != nil {
	// 	log.Fatalf("add members: %v", err)
	// }
	// fmt.Printf("Added %d members\n", added.AddedCount)

	fmt.Println("\n=== 5) Delete workspace ===")
	del, err := client.DeleteWorkspace(ctx, created.Data.ID)
	if err != nil {
		log.Fatalf("delete workspace: %v", err)
	}
	fmt.Printf("Deleted: %v\n", del.Deleted)
}
