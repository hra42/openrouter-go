// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

import (
	"context"
	"fmt"
)

// Predefined validation errors for Anthropic Messages API.
var (
	// ErrNoAnthropicMessages is returned when no messages are provided.
	ErrNoAnthropicMessages = &ValidationError{
		Field:   "messages",
		Message: "at least one message is required",
	}

	// ErrAnthropicMaxTokensRequired is returned when max_tokens is not set.
	ErrAnthropicMaxTokensRequired = &ValidationError{
		Field:   "max_tokens",
		Message: "max_tokens is required and must be greater than 0",
	}
)

// CreateAnthropicMessage sends a request to the Anthropic Messages API via OpenRouter.
//
// Example:
//
//	messages := []openrouter.AnthropicMessage{
//	    openrouter.CreateAnthropicUserMessage("Hello!"),
//	}
//	resp, err := client.CreateAnthropicMessage(ctx, messages,
//	    openrouter.WithAnthropicModel("anthropic/claude-sonnet-4"),
//	    openrouter.WithAnthropicMaxTokens(1024),
//	)
func (c *Client) CreateAnthropicMessage(ctx context.Context, messages []AnthropicMessage, opts ...AnthropicOption) (*AnthropicMessagesResponse, error) {
	// Validate messages
	if err := validateAnthropicMessages(messages); err != nil {
		return nil, err
	}

	// Check API key
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	// Build request
	req := &AnthropicMessagesRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   false,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate request
	if err := validateAnthropicRequest(req); err != nil {
		return nil, err
	}

	// Make request
	var resp AnthropicMessagesResponse
	err := c.doRequest(ctx, "POST", "/messages", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CreateAnthropicMessageStream sends a streaming request to the Anthropic Messages API via OpenRouter.
//
// Example:
//
//	messages := []openrouter.AnthropicMessage{
//	    openrouter.CreateAnthropicUserMessage("Tell me a story"),
//	}
//	stream, err := client.CreateAnthropicMessageStream(ctx, messages,
//	    openrouter.WithAnthropicModel("anthropic/claude-sonnet-4"),
//	    openrouter.WithAnthropicMaxTokens(1024),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer stream.Close()
//
//	for event := range stream.Events() {
//	    fmt.Print(event.GetTextDelta())
//	}
func (c *Client) CreateAnthropicMessageStream(ctx context.Context, messages []AnthropicMessage, opts ...AnthropicOption) (*AnthropicStream, error) {
	// Validate messages
	if err := validateAnthropicMessages(messages); err != nil {
		return nil, err
	}

	// Check API key
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	// Build request
	req := &AnthropicMessagesRequest{
		Model:    c.defaultModel,
		Messages: messages,
		Stream:   true,
	}

	// Apply options
	for _, opt := range opts {
		opt(req)
	}

	// Validate request
	if err := validateAnthropicRequest(req); err != nil {
		return nil, err
	}

	// Create stream
	stream, err := c.createStream(ctx, "/messages", req)
	if err != nil {
		return nil, err
	}

	return &AnthropicStream{
		stream: stream,
	}, nil
}

// validateAnthropicMessages validates the messages array.
func validateAnthropicMessages(messages []AnthropicMessage) error {
	if len(messages) == 0 {
		return ErrNoAnthropicMessages
	}

	validRoles := map[string]bool{
		"user":      true,
		"assistant": true,
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
				Message: fmt.Sprintf("invalid role '%s', must be one of: user, assistant", msg.Role),
			}
		}

		if msg.Content == nil {
			return &ValidationError{
				Field:   fmt.Sprintf("messages[%d].content", i),
				Message: "content is required",
			}
		}
	}

	return nil
}

// validateAnthropicRequest validates the full Anthropic request.
func validateAnthropicRequest(req *AnthropicMessagesRequest) error {
	if req.Model == "" {
		return ErrNoModel
	}

	if req.MaxTokens <= 0 {
		return ErrAnthropicMaxTokensRequired
	}

	// Validate thinking config if present
	if req.Thinking != nil {
		validTypes := map[string]bool{
			AnthropicThinkingEnabled:  true,
			AnthropicThinkingDisabled: true,
		}
		if !validTypes[req.Thinking.Type] {
			return &ValidationError{
				Field:   "thinking.type",
				Message: fmt.Sprintf("invalid thinking type '%s', must be one of: enabled, disabled", req.Thinking.Type),
			}
		}
		if req.Thinking.Type == AnthropicThinkingEnabled && req.Thinking.BudgetTokens <= 0 {
			return &ValidationError{
				Field:   "thinking.budget_tokens",
				Message: "budget_tokens must be greater than 0 when thinking is enabled",
			}
		}
	}

	// Validate tools if present
	for i, tool := range req.Tools {
		// Built-in tools don't need a name in the same way
		if tool.Type == "" && tool.Name == "" {
			return &ValidationError{
				Field:   fmt.Sprintf("tools[%d].name", i),
				Message: "tool name is required for custom tools",
			}
		}
	}

	return nil
}
