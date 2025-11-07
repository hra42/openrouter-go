package main

import (
	"fmt"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	dataURL, err := openrouter.EncodeImageToBase64("test-image.png")
	if err != nil {
		fmt.Printf("❌ Failed to encode: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Image encoded successfully!\n")
	fmt.Printf("   Data URL length: %d bytes\n", len(dataURL))
	fmt.Printf("   Prefix: %s...\n", dataURL[:50])
	fmt.Println("\nThe test image is ready for use in e2e tests!")
}
