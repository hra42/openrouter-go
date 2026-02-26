package openrouter

import (
	"context"
	"fmt"
)

// ChatComplete sends a chat completion request to the OpenRouter API.
func (c *Client) ChatComplete(ctx context.Context, messages []Message, opts ...ChatCompletionOption) (*ChatCompletionResponse, error) {
	// Validate inputs
	if err := c.validateChatRequest(messages); err != nil {
		return nil, err
	}

	// Build request
	req := &ChatCompletionRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   false,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate response format if present
	if req.ResponseFormat != nil {
		if err := validateResponseFormat(req.ResponseFormat); err != nil {
			return nil, err
		}
	}

	// Validate request parameters
	if err := validateChatCompletionParams(req); err != nil {
		return nil, err
	}

	// Handle model suffixes
	req.Model = processModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Make request
	var resp ChatCompletionResponse
	err := c.doRequest(ctx, "POST", "/chat/completions", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// ChatCompleteStream sends a streaming chat completion request to the OpenRouter API.
// This method returns a stream that can be used to receive events as they arrive.
func (c *Client) ChatCompleteStream(ctx context.Context, messages []Message, opts ...ChatCompletionOption) (*ChatStream, error) {
	// Validate inputs
	if err := c.validateChatRequest(messages); err != nil {
		return nil, err
	}

	// Build request
	req := &ChatCompletionRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   true,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate response format if present
	if req.ResponseFormat != nil {
		if err := validateResponseFormat(req.ResponseFormat); err != nil {
			return nil, err
		}
	}

	// Validate request parameters
	if err := validateChatCompletionParams(req); err != nil {
		return nil, err
	}

	// Handle model suffixes
	req.Model = processModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Apply per-request timeout to connection setup only
	connectCtx := ctx
	if req.requestTimeout != nil {
		var cancel context.CancelFunc
		connectCtx, cancel = context.WithTimeout(ctx, *req.requestTimeout)
		defer cancel()
	}

	// Create stream (timeout applies to connection, not stream lifetime)
	stream, err := c.createStream(connectCtx, "/chat/completions", req)
	if err != nil {
		return nil, err
	}

	return &ChatStream{
		stream: stream,
	}, nil
}

// validateChatRequest validates the chat completion request parameters.
func (c *Client) validateChatRequest(messages []Message) error {
	if c.apiKey == "" {
		return ErrNoAPIKey
	}

	if len(messages) == 0 {
		return ErrNoMessages
	}

	// Validate message roles
	validRoles := map[string]bool{
		"system":    true,
		"user":      true,
		"assistant": true,
		"tool":      true,
	}

	for i, msg := range messages {
		if msg.Role == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("messages[%d].role", i),
				Message: "role is required",
			}
		}

		if !validRoles[msg.Role] {
			return &ValidationError{
				Field:   fmt.Sprintf("messages[%d].role", i),
				Message: fmt.Sprintf("invalid role '%s', must be one of: system, user, assistant, tool", msg.Role),
			}
		}

		if msg.Content == nil && msg.Role != "assistant" {
			return &ValidationError{
				Field:   fmt.Sprintf("messages[%d].content", i),
				Message: "content is required for non-assistant messages",
			}
		}
	}

	return nil
}

// CreateChatMessage is a helper function to create a chat message.
func CreateChatMessage(role string, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

// CreateSystemMessage creates a system message.
func CreateSystemMessage(content string) Message {
	return CreateChatMessage("system", content)
}

// CreateUserMessage creates a user message.
func CreateUserMessage(content string) Message {
	return CreateChatMessage("user", content)
}

// CreateAssistantMessage creates an assistant message.
func CreateAssistantMessage(content string) Message {
	return CreateChatMessage("assistant", content)
}

// CreateToolMessage creates a tool message.
func CreateToolMessage(content string, toolCallID string) Message {
	return Message{
		Role:       "tool",
		Content:    content,
		ToolCallID: toolCallID,
	}
}

// CreateMultiModalMessage creates a message with text and image content.
func CreateMultiModalMessage(role string, text string, imageURL string) Message {
	return Message{
		Role: role,
		Content: []ContentPart{
			{Type: "text", Text: text},
			{Type: "image_url", ImageURL: &ImageURL{URL: imageURL}},
		},
	}
}

// CreateUserMessageWithImage creates a user message with text and a single image.
func CreateUserMessageWithImage(text string, imageURL string) Message {
	return CreateMultiModalMessage("user", text, imageURL)
}

// CreateUserMessageWithImages creates a user message with text and multiple images.
func CreateUserMessageWithImages(text string, imageURLs ...string) Message {
	parts := []ContentPart{
		{Type: "text", Text: text},
	}
	for _, url := range imageURLs {
		parts = append(parts, ContentPart{
			Type:     "image_url",
			ImageURL: &ImageURL{URL: url},
		})
	}
	return Message{
		Role:    "user",
		Content: parts,
	}
}

// CreateUserMessageWithImageDetail creates a user message with an image and detail level.
// The detail parameter can be "auto", "low", or "high".
// - "low" is faster and cheaper, suitable for understanding general content
// - "high" provides more detailed analysis at higher cost
// - "auto" lets the model decide based on image size (default)
func CreateUserMessageWithImageDetail(text string, imageURL string, detail string) Message {
	return Message{
		Role: "user",
		Content: []ContentPart{
			{Type: "text", Text: text},
			{Type: "image_url", ImageURL: &ImageURL{URL: imageURL, Detail: detail}},
		},
	}
}

// CreateContentParts is a helper to build custom content part arrays.
type ContentBuilder struct {
	parts []ContentPart
}

// NewContentBuilder creates a new content builder.
func NewContentBuilder() *ContentBuilder {
	return &ContentBuilder{
		parts: []ContentPart{},
	}
}

// AddText adds a text content part.
func (cb *ContentBuilder) AddText(text string) *ContentBuilder {
	cb.parts = append(cb.parts, ContentPart{
		Type: "text",
		Text: text,
	})
	return cb
}

// AddImage adds an image content part.
func (cb *ContentBuilder) AddImage(imageURL string) *ContentBuilder {
	cb.parts = append(cb.parts, ContentPart{
		Type:     "image_url",
		ImageURL: &ImageURL{URL: imageURL},
	})
	return cb
}

// AddImageWithDetail adds an image content part with a detail level.
func (cb *ContentBuilder) AddImageWithDetail(imageURL string, detail string) *ContentBuilder {
	cb.parts = append(cb.parts, ContentPart{
		Type:     "image_url",
		ImageURL: &ImageURL{URL: imageURL, Detail: detail},
	})
	return cb
}

// AddTextFile reads a text file and adds its content to the message.
// The content is formatted with a filename header for context.
func (cb *ContentBuilder) AddTextFile(filePath string) (*ContentBuilder, error) {
	content, filename, err := ReadTextFileWithFilename(filePath)
	if err != nil {
		return cb, err
	}

	cb.parts = append(cb.parts, ContentPart{
		Type: "text",
		Text: fmt.Sprintf("=== %s ===\n%s", filename, content),
	})
	return cb, nil
}

// AddTextFileWithLabel reads a text file and adds its content with a custom label.
func (cb *ContentBuilder) AddTextFileWithLabel(filePath string, label string) (*ContentBuilder, error) {
	content, err := ReadTextFile(filePath)
	if err != nil {
		return cb, err
	}

	cb.parts = append(cb.parts, ContentPart{
		Type: "text",
		Text: fmt.Sprintf("=== %s ===\n%s", label, content),
	})
	return cb, nil
}

// AddTextContent adds raw text content with an optional label.
func (cb *ContentBuilder) AddTextContent(content string, label string) *ContentBuilder {
	var text string
	if label != "" {
		text = fmt.Sprintf("=== %s ===\n%s", label, content)
	} else {
		text = content
	}

	cb.parts = append(cb.parts, ContentPart{
		Type: "text",
		Text: text,
	})
	return cb
}

// Build returns the content parts array.
func (cb *ContentBuilder) Build() []ContentPart {
	return cb.parts
}

// BuildMessage creates a message with the builder's content.
func (cb *ContentBuilder) BuildMessage(role string) Message {
	return Message{
		Role:    role,
		Content: cb.parts,
	}
}
