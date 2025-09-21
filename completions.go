package openrouter

import (
	"context"
	"fmt"
	"strings"
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
	req.Model = c.handleCompletionModelSuffix(req.Model, req)

	// Ensure model is set
	if req.Model == "" {
		return nil, ErrNoModel
	}

	// Make request
	var resp CompletionResponse
	err := c.doRequestWithRetry(ctx, "POST", "/completions", req, &resp)
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
	req.Model = c.handleCompletionModelSuffix(req.Model, req)

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

// CompletionStream represents a streaming completion response.
type CompletionStream struct {
	stream *eventStream
}

// Events returns a channel that receives streaming events.
func (cs *CompletionStream) Events() <-chan CompletionResponse {
	events := make(chan CompletionResponse)

	go func() {
		defer close(events)

		for event := range cs.stream.Events() {
			// Parse the event data into a CompletionResponse
			var response CompletionResponse
			if err := parseSSEData(event.Data, &response); err != nil {
				cs.stream.setError(err)
				return
			}

			select {
			case events <- response:
			case <-cs.stream.ctx.Done():
				return
			}
		}
	}()

	return events
}

// Err returns any error that occurred during streaming.
func (cs *CompletionStream) Err() error {
	return cs.stream.Err()
}

// Close closes the stream.
func (cs *CompletionStream) Close() error {
	return cs.stream.Close()
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

// handleCompletionModelSuffix processes model suffixes like :nitro and :floor for completion requests
func (c *Client) handleCompletionModelSuffix(model string, req *CompletionRequest) string {
	if strings.HasSuffix(model, ":nitro") {
		// Remove suffix and apply throughput sorting
		model = strings.TrimSuffix(model, ":nitro")
		if req.Provider == nil {
			req.Provider = &Provider{}
		}
		req.Provider.Sort = "throughput"
	} else if strings.HasSuffix(model, ":floor") {
		// Remove suffix and apply price sorting
		model = strings.TrimSuffix(model, ":floor")
		if req.Provider == nil {
			req.Provider = &Provider{}
		}
		req.Provider.Sort = "price"
	}
	return model
}