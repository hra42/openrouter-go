package tests

import (
	"fmt"
	"strings"
)

// printSuccess prints a success message with a checkmark
func printSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// printError prints an error message with an X
func printError(message string, err error) {
	fmt.Printf("❌ %s: %v\n", message, err)
}

// printInfo prints an informational message
func printInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

// printVerbose prints a verbose message if verbose mode is enabled
func printVerbose(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Printf("  "+format+"\n", args...)
	}
}

// truncateString truncates a string to maxLen characters, adding ellipsis if needed
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// sanitizeForDisplay sanitizes a string for display by removing control characters
func sanitizeForDisplay(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || (r >= 32 && r < 127) || r >= 160 {
			return r
		}
		return -1
	}, s)
}
