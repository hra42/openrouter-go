package openrouter

import (
	"context"
	"fmt"
)

// Complete sends a legacy completion request to the OpenRouter API.
func (c *Client) Complete(ctx context.Context, prompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	// Validate inputs
	if err := c.validateCompletionRequest(prompt); err != nil {
		return nil, err
	}

	// Build request
	req := &CompletionRequest{
		Model:  c.defaultModel,
		Prompt: prompt,
		Stream: false,
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
	var resp CompletionResponse
	err := c.doRequest(ctx, "POST", "/completions", req, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// CompleteStream sends a streaming completion request to the OpenRouter API.
// This method returns a stream that can be used to receive events as they arrive.
func (c *Client) CompleteStream(ctx context.Context, prompt string, opts ...CompletionOption) (*CompletionStream, error) {
	// Validate inputs
	if err := c.validateCompletionRequest(prompt); err != nil {
		return nil, err
	}

	// Build request
	req := &CompletionRequest{
		Model:  c.defaultModel,
		Prompt: prompt,
		Stream: true,
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
	stream, err := c.createStream(ctx, "/completions", req)
	if err != nil {
		return nil, err
	}

	return &CompletionStream{
		stream: stream,
	}, nil
}

// validateCompletionRequest validates the completion request parameters.
func (c *Client) validateCompletionRequest(prompt string) error {
	if c.apiKey == "" {
		return ErrNoAPIKey
	}

	if prompt == "" {
		return ErrNoPrompt
	}

	return nil
}

// CompleteWithContext is a convenience method that combines prompt completion with context.
func (c *Client) CompleteWithContext(ctx context.Context, contextPrompt, userPrompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	fullPrompt := fmt.Sprintf("%s\n\n%s", contextPrompt, userPrompt)
	return c.Complete(ctx, fullPrompt, opts...)
}

// CompleteWithExamples is a convenience method for few-shot prompting.
func (c *Client) CompleteWithExamples(ctx context.Context, instruction string, examples []string, prompt string, opts ...CompletionOption) (*CompletionResponse, error) {
	fullPrompt := instruction

	if len(examples) > 0 {
		fullPrompt += "\n\nExamples:\n"
		for i, example := range examples {
			fullPrompt += fmt.Sprintf("%d. %s\n", i+1, example)
		}
	}

	fullPrompt += fmt.Sprintf("\n\nNow: %s", prompt)

	return c.Complete(ctx, fullPrompt, opts...)
}
