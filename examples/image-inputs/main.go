package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hra42/openrouter-go"
)

func main() {
	// Get API key from environment variable
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENROUTER_API_KEY environment variable")
	}

	// Create a new client
	client := openrouter.NewClient(openrouter.WithAPIKey(apiKey))

	// Example 1: Single image with URL
	fmt.Println("=== Example 1: Single Image with URL ===")
	singleImageURL(client)

	// Example 2: Multiple images with URLs
	fmt.Println("\n=== Example 2: Multiple Images with URLs ===")
	multipleImagesURL(client)

	// Example 3: Image with detail level
	fmt.Println("\n=== Example 3: Image with Detail Level ===")
	imageWithDetail(client)

	// Example 4: Using ContentBuilder for complex messages
	fmt.Println("\n=== Example 4: Using ContentBuilder ===")
	contentBuilder(client)

	// Example 5: Base64-encoded image from local file
	// Note: This example requires a local image file to work
	// Uncomment to test with your own image
	// fmt.Println("\n=== Example 5: Base64-encoded Image ===")
	// base64EncodedImage(client)
}

func singleImageURL(client *openrouter.Client) {
	// Create a message with text and a single image URL
	// Example: Laboratory test tubes with colorful liquids
	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImage(
			"What's in this image?",
			"https://hra42.com/test-image.png",
		),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
	fmt.Printf("Tokens used: %d\n", resp.Usage.TotalTokens)
}

func multipleImagesURL(client *openrouter.Client) {
	// Create a message with text and multiple images
	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImages(
			"Compare these two images. What are the similarities and differences?",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Placeholder_view_vector.svg/991px-Placeholder_view_vector.svg.png",
		),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func imageWithDetail(client *openrouter.Client) {
	// Create a message with an image and specify detail level
	// "low" - faster and cheaper, suitable for general understanding
	// "high" - more detailed analysis at higher cost
	// "auto" - let the model decide based on image size
	messages := []openrouter.Message{
		openrouter.CreateUserMessageWithImageDetail(
			"Describe this image in detail.",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg",
			"high", // Request high detail analysis
		),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func contentBuilder(client *openrouter.Client) {
	// Use ContentBuilder for more complex message construction
	content := openrouter.NewContentBuilder().
		AddText("I have a few images to show you. First, here's a nature scene:").
		AddImage("https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg").
		AddText("And here's another one:").
		AddImage("https://upload.wikimedia.org/wikipedia/commons/thumb/3/3f/Placeholder_view_vector.svg/991px-Placeholder_view_vector.svg.png").
		AddText("What can you tell me about these images?")

	messages := []openrouter.Message{
		content.BuildMessage("user"),
	}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

func base64EncodedImage(client *openrouter.Client) {
	// Example with a local image file
	// Replace "path/to/your/image.jpg" with an actual image path
	imagePath := "path/to/your/image.jpg"

	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		log.Printf("Image file not found: %s", imagePath)
		log.Println("Please update the imagePath variable with a valid image file path")
		return
	}

	// Create a message with a base64-encoded image
	message, err := openrouter.CreateUserMessageWithBase64Image(
		"What's in this image?",
		imagePath,
	)
	if err != nil {
		log.Printf("Error encoding image: %v", err)
		return
	}

	messages := []openrouter.Message{message}

	resp, err := client.ChatComplete(
		context.Background(),
		messages,
		openrouter.WithModel("google/gemini-2.0-flash-thinking-exp:free"),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
