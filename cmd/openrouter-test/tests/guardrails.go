package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/hra42/openrouter-go"
)

// RunGuardrailsTest tests the guardrails API endpoints with a full lifecycle test.
func RunGuardrailsTest(ctx context.Context, client *openrouter.Client, verbose bool) bool {
	fmt.Printf("🔄 Test: Guardrails API\n")

	// Test 1: List guardrails
	fmt.Printf("   Testing list guardrails...\n")
	start := time.Now()
	listResp, err := client.ListGuardrails(ctx, nil)
	elapsed := time.Since(start)

	if err != nil {
		if reqErr, ok := err.(*openrouter.RequestError); ok {
			if reqErr.StatusCode == 401 || reqErr.StatusCode == 403 {
				fmt.Printf("   ⚠️  Guardrails API requires a provisioning key: %v\n", reqErr.Message)
				fmt.Printf("   Skipping test (provisioning keys are separate from inference API keys)\n")
				fmt.Printf("   Create a provisioning key at: https://openrouter.ai/settings/provisioning-keys\n")
				return true
			}
		}
		printError("Failed to list guardrails", err)
		return false
	}

	fmt.Printf("   ✅ Listed guardrails (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total guardrails: %d\n", listResp.TotalCount)

	if verbose && len(listResp.Data) > 0 {
		for i, g := range listResp.Data {
			if i >= 5 {
				break
			}
			fmt.Printf("      %d. %s (ID: %s)\n", i+1, g.Name, g.ID)
		}
	}

	// Test 2: Create a guardrail
	fmt.Printf("\n   Testing create guardrail...\n")
	guardrailName := fmt.Sprintf("Test Guardrail (openrouter-go %s)", time.Now().Format("2006-01-02 15:04:05"))
	limitUSD := 1.0
	interval := openrouter.ResetIntervalMonthly
	enforceZDR := true
	desc := "Created by openrouter-go test suite"

	start = time.Now()
	created, err := client.CreateGuardrail(ctx, &openrouter.CreateGuardrailRequest{
		Name:          guardrailName,
		Description:   &desc,
		LimitUSD:      &limitUSD,
		ResetInterval: &interval,
		EnforceZDR:    &enforceZDR,
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to create guardrail", err)
		return false
	}

	fmt.Printf("   ✅ Created guardrail (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      ID: %s\n", created.ID)
	fmt.Printf("      Name: %s\n", created.Name)

	guardrailID := created.ID

	// Test 3: Get guardrail
	fmt.Printf("\n   Testing get guardrail...\n")
	start = time.Now()
	fetched, err := client.GetGuardrail(ctx, guardrailID)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to get guardrail", err)
		// Clean up
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Retrieved guardrail (%.2fs)\n", elapsed.Seconds())

	if fetched.ID != guardrailID {
		fmt.Printf("   ❌ ID mismatch: expected %s, got %s\n", guardrailID, fetched.ID)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}
	if fetched.Name != guardrailName {
		fmt.Printf("   ❌ Name mismatch: expected %s, got %s\n", guardrailName, fetched.Name)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	printSuccess("Get guardrail validation passed")

	if verbose {
		fmt.Printf("      ID: %s\n", fetched.ID)
		fmt.Printf("      Name: %s\n", fetched.Name)
		if fetched.Description != nil {
			fmt.Printf("      Description: %s\n", *fetched.Description)
		}
		if fetched.LimitUSD != nil {
			fmt.Printf("      Limit USD: $%.2f\n", *fetched.LimitUSD)
		}
		if fetched.ResetInterval != nil {
			fmt.Printf("      Reset Interval: %s\n", *fetched.ResetInterval)
		}
		if fetched.EnforceZDR != nil {
			fmt.Printf("      Enforce ZDR: %v\n", *fetched.EnforceZDR)
		}
	}

	// Test 4: Update guardrail
	fmt.Printf("\n   Testing update guardrail...\n")
	newName := fmt.Sprintf("Updated Guardrail (%s)", time.Now().Format("15:04:05"))
	newLimit := 2.0
	start = time.Now()
	updated, err := client.UpdateGuardrail(ctx, guardrailID, &openrouter.UpdateGuardrailRequest{
		Name:     &newName,
		LimitUSD: &newLimit,
	})
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to update guardrail", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Updated guardrail (%.2fs)\n", elapsed.Seconds())

	if updated.Name != newName {
		fmt.Printf("   ❌ Name not updated: expected %q, got %q\n", newName, updated.Name)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	printSuccess("Update guardrail validation passed")

	// Test 5: List key assignments for guardrail
	fmt.Printf("\n   Testing list key assignments for guardrail...\n")
	start = time.Now()
	keyAssignments, err := client.ListGuardrailKeyAssignments(ctx, guardrailID, nil)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list key assignments", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Listed key assignments (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Key assignments: %d\n", keyAssignments.TotalCount)

	// Test 6: List member assignments for guardrail
	fmt.Printf("\n   Testing list member assignments for guardrail...\n")
	start = time.Now()
	memberAssignments, err := client.ListGuardrailMemberAssignments(ctx, guardrailID, nil)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list member assignments", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Listed member assignments (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Member assignments: %d\n", memberAssignments.TotalCount)

	// Test 7: List all key assignments
	fmt.Printf("\n   Testing list all key assignments...\n")
	start = time.Now()
	allKeyAssignments, err := client.ListAllKeyAssignments(ctx, nil)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list all key assignments", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Listed all key assignments (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total key assignments: %d\n", allKeyAssignments.TotalCount)

	// Test 8: List all member assignments
	fmt.Printf("\n   Testing list all member assignments...\n")
	start = time.Now()
	allMemberAssignments, err := client.ListAllMemberAssignments(ctx, nil)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to list all member assignments", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}

	fmt.Printf("   ✅ Listed all member assignments (%.2fs)\n", elapsed.Seconds())
	fmt.Printf("      Total member assignments: %d\n", allMemberAssignments.TotalCount)

	// Test 9: Input validation
	fmt.Printf("\n   Testing input validation...\n")

	_, err = client.CreateGuardrail(ctx, nil)
	if err == nil {
		printError("Should have failed with nil request", nil)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for nil request, got %T\n", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}
	printSuccess("Nil request validation works")

	_, err = client.GetGuardrail(ctx, "")
	if err == nil {
		printError("Should have failed with empty id", nil)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}
	if _, ok := openrouter.IsValidationError(err); !ok {
		fmt.Printf("   ❌ Expected ValidationError for empty id, got %T\n", err)
		_, _ = client.DeleteGuardrail(ctx, guardrailID)
		return false
	}
	printSuccess("Empty id validation works")

	// Test 10: Delete guardrail (cleanup)
	fmt.Printf("\n   Testing delete guardrail...\n")
	start = time.Now()
	deleteResp, err := client.DeleteGuardrail(ctx, guardrailID)
	elapsed = time.Since(start)

	if err != nil {
		printError("Failed to delete guardrail", err)
		return false
	}

	fmt.Printf("   ✅ Deleted guardrail (%.2fs)\n", elapsed.Seconds())

	if deleteResp.Deleted {
		printSuccess("Delete operation confirmed via response")
	} else {
		fmt.Printf("   ℹ️  Delete response 'deleted' field is false (verifying deletion...)\n")
	}

	// Verify deletion
	fmt.Printf("\n   Verifying guardrail was deleted...\n")
	_, err = client.GetGuardrail(ctx, guardrailID)
	if err == nil {
		printError("Guardrail still exists after deletion!", nil)
		return false
	}

	if reqErr, ok := err.(*openrouter.RequestError); ok {
		if reqErr.StatusCode == 404 {
			printSuccess("Confirmed guardrail no longer exists (404)")
		} else {
			fmt.Printf("   ⚠️  Unexpected status code when verifying deletion: %d\n", reqErr.StatusCode)
		}
	}

	printSuccess("Guardrails API tests completed")
	return true
}
