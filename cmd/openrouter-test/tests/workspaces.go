package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunWorkspacesTest exercises the Workspaces endpoints end-to-end:
// list → create → get → update → delete, plus empty add/remove members.
// Requires a Management (Provisioning) API key. With a regular inference key,
// the test skips cleanly on 401/403.
func RunWorkspacesTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Workspaces\n")

	// Step 1: list
	fmt.Printf("   Listing workspaces...\n")
	start := time.Now()
	listResp, err := client.ListWorkspaces(ctx, nil)
	elapsed := time.Since(start)
	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Workspaces endpoint requires a management key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (management keys are separate from inference API keys)\n")
				fmt.Printf("   Create one at: https://openrouter.ai/settings/provisioning-keys\n")
				return true
			}
		}
		printError("Failed to list workspaces", err)
		return false
	}
	fmt.Printf("   ✅ Listed %d workspaces (total=%d, %.2fs)\n", len(listResp.Data), listResp.TotalCount, elapsed.Seconds())
	if verbose && len(listResp.Data) > 0 {
		for _, ws := range listResp.Data {
			fmt.Printf("      - %s (%s) id=%s\n", ws.Name, ws.Slug, ws.ID)
		}
	}

	// Step 2: create throwaway workspace with unique slug
	slug := fmt.Sprintf("sdk-e2e-%d", time.Now().Unix())
	fmt.Printf("\n   Creating workspace %q...\n", slug)
	desc := "Temporary workspace created by openrouter-go e2e tests"
	createResp, err := client.CreateWorkspace(ctx, &openrouter.CreateWorkspaceRequest{
		Name:        "SDK E2E Test",
		Slug:        slug,
		Description: &desc,
	})
	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok && (reqErr.StatusCode == 401 || reqErr.StatusCode == 403) {
			fmt.Printf("   ⚠️  Create workspace forbidden (key lacks permission). Skipping.\n")
			return true
		}
		printError("Failed to create workspace", err)
		return false
	}
	workspaceID := createResp.Data.ID
	fmt.Printf("   ✅ Created workspace id=%s slug=%s\n", workspaceID, createResp.Data.Slug)

	// Ensure cleanup even if later steps fail
	defer func() {
		if _, err := client.DeleteWorkspace(ctx, workspaceID); err != nil {
			fmt.Printf("   ⚠️  Cleanup: failed to delete workspace %s: %v\n", workspaceID, err)
		}
	}()

	// Step 3: get by slug
	fmt.Printf("\n   Getting workspace by slug...\n")
	getResp, err := client.GetWorkspace(ctx, slug)
	if err != nil {
		printError("Failed to get workspace", err)
		return false
	}
	if getResp.Data.ID != workspaceID {
		fmt.Printf("   ❌ Get returned mismatched ID: %s vs %s\n", getResp.Data.ID, workspaceID)
		return false
	}
	fmt.Printf("   ✅ Got workspace id=%s\n", getResp.Data.ID)

	// Step 4: update name
	fmt.Printf("\n   Updating workspace name...\n")
	newName := "SDK E2E Test (updated)"
	updateResp, err := client.UpdateWorkspace(ctx, workspaceID, &openrouter.UpdateWorkspaceRequest{
		Name: &newName,
	})
	if err != nil {
		printError("Failed to update workspace", err)
		return false
	}
	if updateResp.Data.Name != newName {
		fmt.Printf("   ❌ Update did not persist name: got %q\n", updateResp.Data.Name)
		return false
	}
	fmt.Printf("   ✅ Updated name to %q\n", updateResp.Data.Name)

	// Step 5: add/remove members with empty list — either succeeds (0 added) or
	// returns 400. Both are acceptable outcomes for this smoke check.
	fmt.Printf("\n   Adding members (empty list)...\n")
	if addResp, err := client.AddWorkspaceMembers(ctx, workspaceID, []string{}); err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok && reqErr.StatusCode == 400 {
			fmt.Printf("   ℹ️  Empty add list rejected with 400 (acceptable)\n")
		} else {
			printError("Failed to add members", err)
			return false
		}
	} else {
		fmt.Printf("   ✅ AddedCount=%d\n", addResp.AddedCount)
	}

	fmt.Printf("\n   Removing members (empty list)...\n")
	if removeResp, err := client.RemoveWorkspaceMembers(ctx, workspaceID, []string{}); err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok && reqErr.StatusCode == 400 {
			fmt.Printf("   ℹ️  Empty remove list rejected with 400 (acceptable)\n")
		} else {
			printError("Failed to remove members", err)
			return false
		}
	} else {
		fmt.Printf("   ✅ RemovedCount=%d\n", removeResp.RemovedCount)
	}

	printSuccess("Workspaces tests completed")
	return true
}
