package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunOrganizationMembersTest tests the ListOrganizationMembers endpoint.
func RunOrganizationMembersTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: List Organization Members\n")

	fmt.Printf("   Testing list all organization members...\n")
	start := time.Now()
	resp, err := client.ListOrganizationMembers(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  List organization members requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true
			}
			if reqErr.StatusCode == 404 {
				fmt.Printf("   ⚠️  No organization associated with this key (404)\n")
				fmt.Printf("   Skipping test\n")
				return true
			}
		}
		printError("Failed to list organization members", err)
		return false
	}

	fmt.Printf("   ✅ Retrieved organization members (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total members: %d\n", resp.TotalCount)
	fmt.Printf("      Returned in page: %d\n", len(resp.Data))

	// Validate response structure for at least one member
	if len(resp.Data) > 0 {
		fmt.Printf("\n   Validating response structure...\n")
		first := resp.Data[0]
		if first.ID == "" {
			printError("Member missing ID", nil)
			return false
		}
		if first.Email == "" {
			printError("Member missing Email", nil)
			return false
		}
		if first.Role != openrouter.OrganizationMemberRoleAdmin && first.Role != openrouter.OrganizationMemberRoleMember {
			fmt.Printf("   ⚠️  Unexpected role value: %q\n", first.Role)
		}
		printSuccess("Response structure validation passed")

		if verbose {
			fmt.Printf("\n   First member:\n")
			fmt.Printf("      ID: %s\n", first.ID)
			fmt.Printf("      Email: %s\n", first.Email)
			if first.FirstName != nil {
				fmt.Printf("      First name: %s\n", *first.FirstName)
			} else {
				fmt.Printf("      First name: (null)\n")
			}
			if first.LastName != nil {
				fmt.Printf("      Last name: %s\n", *first.LastName)
			} else {
				fmt.Printf("      Last name: (null)\n")
			}
			fmt.Printf("      Role: %s\n", first.Role)
		}
	} else {
		fmt.Printf("   ℹ️  Organization has no members (unusual)\n")
	}

	// Test pagination with limit=1
	if resp.TotalCount > 0 {
		fmt.Printf("\n   Testing pagination with limit=1...\n")
		limit := 1
		start = time.Now()
		pageResp, err := client.ListOrganizationMembers(ctx, &openrouter.ListOrganizationMembersOptions{
			Limit: &limit,
		})
		elapsed = time.Since(start)
		if err != nil {
			printError("Failed to list with limit=1", err)
			return false
		}
		if len(pageResp.Data) > 1 {
			fmt.Printf("   ❌ Expected at most 1 member with limit=1, got %d\n", len(pageResp.Data))
			return false
		}
		fmt.Printf("   ✅ limit=1 returned %d member(s) (%.2fs)\n", len(pageResp.Data), elapsed.Seconds())
	}

	// Custom timeout
	fmt.Printf("\n   Testing with custom timeout...\n")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := client.ListOrganizationMembers(ctxWithTimeout, nil); err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Provisioning key required (expected)\n")
			} else {
				printError("Failed with custom timeout", err)
				return false
			}
		} else if err != context.DeadlineExceeded {
			printError("Failed with custom timeout", err)
			return false
		}
	} else {
		printSuccess("Custom timeout context works")
	}

	printSuccess("List organization members tests completed")
	return true
}
