// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

import (
	"time"
)

// ChatCompletionRequest represents a chat completion request to the OpenRouter API.
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`

	// Optional parameters
	Temperature      *float64               `json:"temperature,omitempty"`
	TopP             *float64               `json:"top_p,omitempty"`
	TopK             *int                   `json:"top_k,omitempty"`
	FrequencyPenalty *float64               `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64               `json:"presence_penalty,omitempty"`
	RepetitionPenalty *float64              `json:"repetition_penalty,omitempty"`
	MaxTokens        *int                   `json:"max_tokens,omitempty"`
	MinP             *float64               `json:"min_p,omitempty"`
	TopA             *float64               `json:"top_a,omitempty"`
	Seed             *int                   `json:"seed,omitempty"`
	Stop             []string               `json:"stop,omitempty"`
	Stream           bool                   `json:"stream,omitempty"`
	LogProbs         *bool                  `json:"logprobs,omitempty"`
	TopLogProbs      *int                   `json:"top_logprobs,omitempty"`
	ResponseFormat   *ResponseFormat        `json:"response_format,omitempty"`
	Tools            []Tool                 `json:"tools,omitempty"`
	ToolChoice       interface{}            `json:"tool_choice,omitempty"`
	Provider         *Provider              `json:"provider,omitempty"`
	Transforms       []string               `json:"transforms,omitempty"`
	Models           []string               `json:"models,omitempty"`
	Route            string                 `json:"route,omitempty"`
	Metadata         map[string]interface{} `json:"-"` // Used for headers
}

// CompletionRequest represents a legacy completion request to the OpenRouter API.
type CompletionRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`

	// Optional parameters
	Temperature       *float64               `json:"temperature,omitempty"`
	TopP              *float64               `json:"top_p,omitempty"`
	TopK              *int                   `json:"top_k,omitempty"`
	FrequencyPenalty  *float64               `json:"frequency_penalty,omitempty"`
	PresencePenalty   *float64               `json:"presence_penalty,omitempty"`
	RepetitionPenalty *float64               `json:"repetition_penalty,omitempty"`
	MaxTokens         *int                   `json:"max_tokens,omitempty"`
	MinP              *float64               `json:"min_p,omitempty"`
	TopA              *float64               `json:"top_a,omitempty"`
	Seed              *int                   `json:"seed,omitempty"`
	Stop              []string               `json:"stop,omitempty"`
	Stream            bool                   `json:"stream,omitempty"`
	LogProbs          *int                   `json:"logprobs,omitempty"`
	Echo              *bool                  `json:"echo,omitempty"`
	N                 *int                   `json:"n,omitempty"`
	BestOf            *int                   `json:"best_of,omitempty"`
	Suffix            *string                `json:"suffix,omitempty"`
	Provider          *Provider              `json:"provider,omitempty"`
	Transforms        []string               `json:"transforms,omitempty"`
	Models            []string               `json:"models,omitempty"`
	Route             string                 `json:"route,omitempty"`
	Metadata          map[string]interface{} `json:"-"` // Used for headers
}

// Message represents a message in the chat completion request.
type Message struct {
	Role       string       `json:"role"`
	Content    MessageContent `json:"content"`
	Name       string       `json:"name,omitempty"`
	ToolCalls  []ToolCall   `json:"tool_calls,omitempty"`
	ToolCallID string       `json:"tool_call_id,omitempty"`
}

// MessageContent can be either a string or an array of content parts.
type MessageContent interface{}

// ContentPart represents a part of message content (text or image).
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL represents an image URL in the message content.
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// ResponseFormat specifies the format of the response.
type ResponseFormat struct {
	Type       string                 `json:"type"`
	JSONSchema map[string]interface{} `json:"json_schema,omitempty"`
}

// Tool represents a tool/function that can be called.
type Tool struct {
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// Function represents a callable function.
type Function struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ToolCall represents a tool call made by the model.
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call made by the model.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// Provider represents provider-specific parameters for routing requests.
type Provider struct {
	// Order specifies provider slugs to try in order (e.g. ["anthropic", "openai"])
	Order            []string               `json:"order,omitempty"`
	// RequireParameters only uses providers that support all parameters in the request
	RequireParameters *bool                 `json:"require_parameters,omitempty"`
	// DataCollection controls whether to use providers that may store data ("allow" or "deny")
	DataCollection   string                 `json:"data_collection,omitempty"`
	// AllowFallbacks allows backup providers when the primary is unavailable
	AllowFallbacks   *bool                  `json:"allow_fallbacks,omitempty"`
	// Ignore specifies provider slugs to skip for this request
	Ignore           []string               `json:"ignore,omitempty"`
	// Quantizations filters providers by quantization levels (e.g. ["int4", "int8"])
	Quantizations    []string               `json:"quantizations,omitempty"`
	// ZDR restricts routing to only Zero Data Retention endpoints
	ZDR              *bool                  `json:"zdr,omitempty"`
	// Only specifies provider slugs to allow for this request
	Only             []string               `json:"only,omitempty"`
	// Sort providers by "price", "throughput", or "latency"
	Sort             string                 `json:"sort,omitempty"`
	// MaxPrice specifies maximum pricing constraints for the request
	MaxPrice         *MaxPrice              `json:"max_price,omitempty"`

	// Deprecated: Use Ignore instead
	IgnoreProviders  []string               `json:"-"`
	// Deprecated: Use Quantizations instead
	QuantizationFallback map[string]string  `json:"-"`
	// Internal provider parameters
	ProviderParams   map[string]interface{} `json:"-"`
}

// ChatCompletionResponse represents a chat completion response from the OpenRouter API.
type ChatCompletionResponse struct {
	ID                string    `json:"id"`
	Object            string    `json:"object"`
	Created           int64     `json:"created"`
	Model             string    `json:"model"`
	Choices           []Choice  `json:"choices"`
	Usage             Usage     `json:"usage"`
	SystemFingerprint string    `json:"system_fingerprint,omitempty"`
}

// CompletionResponse represents a legacy completion response from the OpenRouter API.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// Choice represents a choice in the chat completion response.
type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`
	FinishReason string       `json:"finish_reason"`
	LogProbs     *LogProbs    `json:"logprobs,omitempty"`
	Delta        *Message     `json:"delta,omitempty"` // For streaming
}

// CompletionChoice represents a choice in the legacy completion response.
type CompletionChoice struct {
	Index        int       `json:"index"`
	Text         string    `json:"text"`
	FinishReason string    `json:"finish_reason"`
	LogProbs     *LogProbs `json:"logprobs,omitempty"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LogProbs represents log probability information.
type LogProbs struct {
	Content []LogProbContent `json:"content,omitempty"`
}

// LogProbContent represents log probability content.
type LogProbContent struct {
	Token       string     `json:"token"`
	LogProb     float64    `json:"logprob"`
	Bytes       []int      `json:"bytes,omitempty"`
	TopLogProbs []TopLogProb `json:"top_logprobs,omitempty"`
}

// TopLogProb represents top log probability information.
type TopLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

// MaxPrice represents maximum pricing constraints for a request.
type MaxPrice struct {
	// Prompt specifies max price per million prompt tokens
	Prompt float64 `json:"prompt,omitempty"`
	// Completion specifies max price per million completion tokens
	Completion float64 `json:"completion,omitempty"`
	// Request specifies max price per request (for providers with per-request pricing)
	Request float64 `json:"request,omitempty"`
	// Image specifies max price per image
	Image float64 `json:"image,omitempty"`
}

// StreamEvent represents a server-sent event for streaming responses.
type StreamEvent struct {
	ID      string
	Event   string
	Data    string
	Retry   *time.Duration
}

// ErrorResponse represents an error response from the OpenRouter API.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents the error details in an error response.
type APIError struct {
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	Code    string                 `json:"code,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}