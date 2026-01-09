// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

// ResponsesRequest represents a request to the OpenRouter Responses API.
//
// WARNING: BETA API - EXPECT BREAKING CHANGES
//
// This API is in beta and may have breaking changes at any time.
// Do not rely on this for production workloads. The API structure,
// parameters, and behavior may change without notice.
//
// For stable production use, consider using ChatCompletionRequest instead.
//
// Features: reasoning, tool calling, web search integration, and streaming.
type ResponsesRequest struct {
	// Model is the model identifier (required).
	// Example: "openai/o4-mini"
	Model string `json:"model"`

	// Input is the input to the model (required).
	// Can be a string for simple text input, or []ResponsesInputItem for structured messages.
	Input interface{} `json:"input"`

	// Stream enables streaming responses via Server-Sent Events.
	Stream bool `json:"stream,omitempty"`

	// MaxOutputTokens limits the number of tokens in the response.
	MaxOutputTokens *int `json:"max_output_tokens,omitempty"`

	// Temperature controls randomness in sampling (0-2).
	Temperature *float64 `json:"temperature,omitempty"`

	// TopP controls nucleus sampling (0-1).
	TopP *float64 `json:"top_p,omitempty"`

	// Reasoning configures reasoning capabilities.
	Reasoning *ResponsesReasoning `json:"reasoning,omitempty"`

	// Tools defines available functions the model can call.
	Tools []Tool `json:"tools,omitempty"`

	// ToolChoice controls tool invocation behavior.
	// Can be "auto", "none", or a specific tool specification.
	ToolChoice interface{} `json:"tool_choice,omitempty"`

	// Plugins configures plugins like web search.
	Plugins []Plugin `json:"plugins,omitempty"`

	// Metadata is used for custom headers (not serialized to JSON).
	Metadata map[string]interface{} `json:"-"`
}

// GetMetadata returns the metadata map for header generation.
func (r *ResponsesRequest) GetMetadata() map[string]interface{} {
	return r.Metadata
}

// ResponsesReasoning configures reasoning capabilities for the Responses API.
type ResponsesReasoning struct {
	// Effort controls the computational intensity of reasoning.
	// Valid values: "minimal", "low", "medium", "high"
	Effort string `json:"effort"`
}

// ResponsesInputItem represents an item in the structured input array.
// Used for multi-turn conversations and function call outputs.
type ResponsesInputItem struct {
	// Type specifies the item type: "message" or "function_call_output"
	Type string `json:"type"`

	// ID is the unique identifier for this item (used in multi-turn conversations).
	ID string `json:"id,omitempty"`

	// Status indicates the completion status (e.g., "completed").
	Status string `json:"status,omitempty"`

	// Role is the message role: "user", "assistant", "system"
	// Only used when Type is "message".
	Role string `json:"role,omitempty"`

	// Content is the message content array.
	// Only used when Type is "message".
	Content []ResponsesInputContent `json:"content,omitempty"`

	// CallID is the ID of the function call being responded to.
	// Only used when Type is "function_call_output".
	CallID string `json:"call_id,omitempty"`

	// Output is the result of the function call.
	// Only used when Type is "function_call_output".
	Output string `json:"output,omitempty"`
}

// ResponsesInputContent represents a content item in a message.
type ResponsesInputContent struct {
	// Type specifies the content type: "input_text"
	Type string `json:"type"`

	// Text is the text content.
	Text string `json:"text"`
}

// ResponsesResponse represents a response from the OpenRouter Responses API.
type ResponsesResponse struct {
	// ID is the unique identifier for this response.
	ID string `json:"id"`

	// Object is the object type, always "response".
	Object string `json:"object"`

	// CreatedAt is the Unix timestamp when the response was created.
	CreatedAt int64 `json:"created_at"`

	// Model is the model that generated the response.
	Model string `json:"model"`

	// Output is the array of output items (messages, function calls, etc.).
	Output []ResponsesOutput `json:"output"`

	// Usage contains token usage information.
	Usage ResponsesUsage `json:"usage"`

	// Status indicates the completion status (e.g., "completed").
	Status string `json:"status"`

	// Metadata contains additional response metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ResponsesOutput represents an output item in the response.
type ResponsesOutput struct {
	// Type specifies the output type: "message" or "function_call"
	Type string `json:"type"`

	// ID is the unique identifier for this output item.
	ID string `json:"id"`

	// Status indicates the completion status.
	Status string `json:"status,omitempty"`

	// Role is the message role (for "message" type).
	Role string `json:"role,omitempty"`

	// Content is the message content array (for "message" type).
	Content []ResponsesOutputContent `json:"content,omitempty"`

	// CallID is the function call ID (for "function_call" type).
	CallID string `json:"call_id,omitempty"`

	// Name is the function name (for "function_call" type).
	Name string `json:"name,omitempty"`

	// Arguments is the JSON-encoded function arguments (for "function_call" type).
	Arguments string `json:"arguments,omitempty"`
}

// ResponsesOutputContent represents a content item in the output.
type ResponsesOutputContent struct {
	// Type specifies the content type: "output_text" or "reasoning"
	Type string `json:"type"`

	// Text is the text content (for "output_text" type).
	Text string `json:"text,omitempty"`

	// Annotations contains URL citations and other annotations.
	Annotations []ResponsesAnnotation `json:"annotations,omitempty"`

	// EncryptedContent contains the encrypted reasoning chain (for "reasoning" type).
	EncryptedContent string `json:"encrypted_content,omitempty"`

	// Summary contains key reasoning steps as text (for "reasoning" type).
	Summary []string `json:"summary,omitempty"`
}

// ResponsesAnnotation represents an annotation in the output content.
type ResponsesAnnotation struct {
	// Type specifies the annotation type (e.g., "url_citation").
	Type string `json:"type"`

	// URL is the cited URL (for "url_citation" type).
	URL string `json:"url,omitempty"`

	// StartIndex is the start position of the citation in the text.
	StartIndex int `json:"start_index,omitempty"`

	// EndIndex is the end position of the citation in the text.
	EndIndex int `json:"end_index,omitempty"`
}

// ResponsesUsage contains token usage information for the Responses API.
type ResponsesUsage struct {
	// InputTokens is the number of tokens in the input.
	InputTokens int `json:"input_tokens"`

	// OutputTokens is the number of tokens in the output.
	OutputTokens int `json:"output_tokens"`

	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`

	// ReasoningTokens is the number of tokens used for reasoning (if reasoning is enabled).
	ReasoningTokens int `json:"reasoning_tokens,omitempty"`
}

// Helper functions for creating input items

// CreateResponsesUserMessage creates a user message input item.
func CreateResponsesUserMessage(text string) ResponsesInputItem {
	return ResponsesInputItem{
		Type: "message",
		Role: "user",
		Content: []ResponsesInputContent{
			{Type: "input_text", Text: text},
		},
	}
}

// CreateResponsesAssistantMessage creates an assistant message input item.
func CreateResponsesAssistantMessage(text string) ResponsesInputItem {
	return ResponsesInputItem{
		Type: "message",
		Role: "assistant",
		Content: []ResponsesInputContent{
			{Type: "input_text", Text: text},
		},
	}
}

// CreateResponsesSystemMessage creates a system message input item.
func CreateResponsesSystemMessage(text string) ResponsesInputItem {
	return ResponsesInputItem{
		Type: "message",
		Role: "system",
		Content: []ResponsesInputContent{
			{Type: "input_text", Text: text},
		},
	}
}

// CreateResponsesMessage creates a message input item with the specified role and text.
func CreateResponsesMessage(role, text string) ResponsesInputItem {
	return ResponsesInputItem{
		Type: "message",
		Role: role,
		Content: []ResponsesInputContent{
			{Type: "input_text", Text: text},
		},
	}
}

// CreateResponsesFunctionOutput creates a function call output item.
func CreateResponsesFunctionOutput(callID, output string) ResponsesInputItem {
	return ResponsesInputItem{
		Type:   "function_call_output",
		CallID: callID,
		Output: output,
	}
}

// GetTextContent extracts the text content from a ResponsesResponse.
// Returns the first text content found in the output, or an empty string if none found.
func (r *ResponsesResponse) GetTextContent() string {
	for _, output := range r.Output {
		if output.Type == "message" {
			for _, content := range output.Content {
				if content.Type == "output_text" && content.Text != "" {
					return content.Text
				}
			}
		}
	}
	return ""
}

// GetFunctionCalls returns all function calls from the response.
func (r *ResponsesResponse) GetFunctionCalls() []ResponsesOutput {
	var calls []ResponsesOutput
	for _, output := range r.Output {
		if output.Type == "function_call" {
			calls = append(calls, output)
		}
	}
	return calls
}

// GetAnnotations returns all annotations from the response.
func (r *ResponsesResponse) GetAnnotations() []ResponsesAnnotation {
	var annotations []ResponsesAnnotation
	for _, output := range r.Output {
		for _, content := range output.Content {
			annotations = append(annotations, content.Annotations...)
		}
	}
	return annotations
}

// GetReasoningSummary returns the reasoning summary if present.
func (r *ResponsesResponse) GetReasoningSummary() []string {
	for _, output := range r.Output {
		for _, content := range output.Content {
			if content.Type == "reasoning" && len(content.Summary) > 0 {
				return content.Summary
			}
		}
	}
	return nil
}
