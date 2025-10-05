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

	// Handle model suffixes
	req.Model = processModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Create stream
	stream, err := c.createStream(ctx, "/chat/completions", req)
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
