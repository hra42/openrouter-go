package openrouter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateUserMessageWithImage(t *testing.T) {
	msg := CreateUserMessageWithImage(
		"What's in this image?",
		"https://example.com/image.jpg",
	)

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatal("expected content to be []ContentPart")
	}

	if len(parts) != 2 {
		t.Fatalf("expected 2 content parts, got %d", len(parts))
	}

	// Check text part
	if parts[0].Type != "text" {
		t.Errorf("expected first part type 'text', got %q", parts[0].Type)
	}
	if parts[0].Text != "What's in this image?" {
		t.Errorf("unexpected text: %q", parts[0].Text)
	}

	// Check image part
	if parts[1].Type != "image_url" {
		t.Errorf("expected second part type 'image_url', got %q", parts[1].Type)
	}
	if parts[1].ImageURL == nil {
		t.Fatal("expected ImageURL to be non-nil")
	}
	if parts[1].ImageURL.URL != "https://example.com/image.jpg" {
		t.Errorf("unexpected image URL: %q", parts[1].ImageURL.URL)
	}
}

func TestCreateUserMessageWithImages(t *testing.T) {
	msg := CreateUserMessageWithImages(
		"Compare these images",
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
	)

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatal("expected content to be []ContentPart")
	}

	// 1 text part + 3 image parts = 4 total
	if len(parts) != 4 {
		t.Fatalf("expected 4 content parts, got %d", len(parts))
	}

	// Check text part
	if parts[0].Type != "text" {
		t.Errorf("expected first part type 'text', got %q", parts[0].Type)
	}

	// Check image parts
	expectedURLs := []string{
		"https://example.com/image1.jpg",
		"https://example.com/image2.jpg",
		"https://example.com/image3.jpg",
	}
	for i, expectedURL := range expectedURLs {
		partIndex := i + 1
		if parts[partIndex].Type != "image_url" {
			t.Errorf("expected part %d type 'image_url', got %q", partIndex, parts[partIndex].Type)
		}
		if parts[partIndex].ImageURL == nil {
			t.Fatalf("expected part %d ImageURL to be non-nil", partIndex)
		}
		if parts[partIndex].ImageURL.URL != expectedURL {
			t.Errorf("part %d: expected URL %q, got %q", partIndex, expectedURL, parts[partIndex].ImageURL.URL)
		}
	}
}

func TestCreateUserMessageWithImageDetail(t *testing.T) {
	msg := CreateUserMessageWithImageDetail(
		"Describe this image",
		"https://example.com/image.jpg",
		"high",
	)

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatal("expected content to be []ContentPart")
	}

	if len(parts) != 2 {
		t.Fatalf("expected 2 content parts, got %d", len(parts))
	}

	// Check image part has detail
	if parts[1].ImageURL == nil {
		t.Fatal("expected ImageURL to be non-nil")
	}
	if parts[1].ImageURL.Detail != "high" {
		t.Errorf("expected detail 'high', got %q", parts[1].ImageURL.Detail)
	}
}

func TestContentBuilder(t *testing.T) {
	builder := NewContentBuilder()

	content := builder.
		AddText("First text").
		AddImage("https://example.com/image1.jpg").
		AddText("Second text").
		AddImageWithDetail("https://example.com/image2.jpg", "low").
		Build()

	if len(content) != 4 {
		t.Fatalf("expected 4 content parts, got %d", len(content))
	}

	// Check first text
	if content[0].Type != "text" || content[0].Text != "First text" {
		t.Error("unexpected first content part")
	}

	// Check first image
	if content[1].Type != "image_url" || content[1].ImageURL.URL != "https://example.com/image1.jpg" {
		t.Error("unexpected second content part")
	}

	// Check second text
	if content[2].Type != "text" || content[2].Text != "Second text" {
		t.Error("unexpected third content part")
	}

	// Check second image with detail
	if content[3].Type != "image_url" ||
		content[3].ImageURL.URL != "https://example.com/image2.jpg" ||
		content[3].ImageURL.Detail != "low" {
		t.Error("unexpected fourth content part")
	}
}

func TestContentBuilderBuildMessage(t *testing.T) {
	builder := NewContentBuilder()

	msg := builder.
		AddText("Hello").
		AddImage("https://example.com/image.jpg").
		BuildMessage("user")

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatal("expected content to be []ContentPart")
	}

	if len(parts) != 2 {
		t.Fatalf("expected 2 content parts, got %d", len(parts))
	}
}

func TestImageMessageSerialization(t *testing.T) {
	msg := CreateUserMessageWithImage(
		"What's in this image?",
		"https://example.com/image.jpg",
	)

	// Serialize to JSON
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}

	// Deserialize back
	var decoded Message
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	// Verify structure
	if decoded.Role != "user" {
		t.Errorf("expected role 'user', got %q", decoded.Role)
	}

	// Content should be decoded as []interface{} (JSON arrays)
	// We need to check it's properly serialized
	t.Logf("Serialized JSON: %s", string(data))

	// Verify the JSON contains expected fields
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	content, ok := jsonMap["content"].([]interface{})
	if !ok {
		t.Fatal("expected content to be an array")
	}

	if len(content) != 2 {
		t.Fatalf("expected 2 content parts, got %d", len(content))
	}
}

func TestChatCompleteWithImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		// Verify the message structure
		if len(req.Messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(req.Messages))
		}

		// Content should be array for multimodal
		contentArray, ok := req.Messages[0].Content.([]interface{})
		if !ok {
			t.Errorf("expected content to be array, got %T", req.Messages[0].Content)
		}

		if len(contentArray) != 2 {
			t.Errorf("expected 2 content parts, got %d", len(contentArray))
		}

		// Send response
		response := ChatCompletionResponse{
			ID:      "chat-123",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "gpt-4-vision",
			Choices: []Choice{
				{
					Index: 0,
					Message: Message{
						Role:    "assistant",
						Content: "I see a beautiful nature scene.",
					},
					FinishReason: "stop",
				},
			},
			Usage: Usage{
				PromptTokens:     100,
				CompletionTokens: 10,
				TotalTokens:      110,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	messages := []Message{
		CreateUserMessageWithImage(
			"What's in this image?",
			"https://example.com/image.jpg",
		),
	}

	resp, err := client.ChatComplete(context.Background(), messages,
		WithModel("gpt-4-vision"),
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Choices[0].Message.Content != "I see a beautiful nature scene." {
		t.Errorf("unexpected response: %q", resp.Choices[0].Message.Content)
	}
}

func TestEncodeImageToBase64(t *testing.T) {
	// Create a temporary test image
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.png")

	// Write a simple 1x1 PNG image
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xDE, 0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, // IDAT chunk
		0x54, 0x08, 0xD7, 0x63, 0xF8, 0xFF, 0xFF, 0x3F,
		0x00, 0x05, 0xFE, 0x02, 0xFE, 0xDC, 0xCC, 0x59,
		0xE7, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, // IEND chunk
		0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	if err := os.WriteFile(imagePath, pngData, 0644); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	// Test encoding
	dataURL, err := EncodeImageToBase64(imagePath)
	if err != nil {
		t.Fatalf("failed to encode image: %v", err)
	}

	// Verify it starts with the correct data URL prefix
	expectedPrefix := "data:image/png;base64,"
	if len(dataURL) < len(expectedPrefix) {
		t.Fatal("data URL too short")
	}
	if dataURL[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, dataURL[:len(expectedPrefix)])
	}

	// Verify it contains base64 data
	if len(dataURL) <= len(expectedPrefix) {
		t.Error("data URL has no base64 data")
	}
}

func TestEncodeImageToBase64UnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(imagePath, []byte("not an image"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := EncodeImageToBase64(imagePath)
	if err == nil {
		t.Error("expected error for unsupported format, got nil")
	}
}

func TestEncodeImageBytesToBase64(t *testing.T) {
	imageData := []byte("fake image data")
	dataURL := EncodeImageBytesToBase64(imageData, "image/jpeg")

	expectedPrefix := "data:image/jpeg;base64,"
	if dataURL[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected prefix %q, got %q", expectedPrefix, dataURL[:len(expectedPrefix)])
	}
}

func TestCreateUserMessageWithBase64Image(t *testing.T) {
	tmpDir := t.TempDir()
	imagePath := filepath.Join(tmpDir, "test.jpg")

	// Write a minimal JPEG
	jpegData := []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46,
		0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x01,
		0x00, 0x01, 0x00, 0x00, 0xFF, 0xD9,
	}

	if err := os.WriteFile(imagePath, jpegData, 0644); err != nil {
		t.Fatalf("failed to write test image: %v", err)
	}

	msg, err := CreateUserMessageWithBase64Image(
		"What's in this image?",
		imagePath,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatal("expected content to be []ContentPart")
	}

	if len(parts) != 2 {
		t.Fatalf("expected 2 content parts, got %d", len(parts))
	}

	// Verify the image URL is a data URL
	if parts[1].ImageURL == nil {
		t.Fatal("expected ImageURL to be non-nil")
	}
	expectedPrefix := "data:image/jpeg;base64,"
	if len(parts[1].ImageURL.URL) < len(expectedPrefix) {
		t.Fatal("image URL too short")
	}
	if parts[1].ImageURL.URL[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("expected data URL prefix, got %q", parts[1].ImageURL.URL[:len(expectedPrefix)])
	}
}
