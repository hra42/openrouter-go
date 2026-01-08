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
	LogProbs          *bool                  `json:"logprobs,omitempty"`
	TopLogProbs       *int                   `json:"top_logprobs,omitempty"`
	ResponseFormat    *ResponseFormat        `json:"response_format,omitempty"`
	Tools             []Tool                 `json:"tools,omitempty"`
	ToolChoice        interface{}            `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool                  `json:"parallel_tool_calls,omitempty"`
	Provider          *Provider              `json:"provider,omitempty"`
	Transforms        []string               `json:"transforms,omitempty"`
	Models            []string               `json:"models,omitempty"`
	Route             string                 `json:"route,omitempty"`
	Plugins           []Plugin               `json:"plugins,omitempty"`
	WebSearchOptions  *WebSearchOptions      `json:"web_search_options,omitempty"`
	User              string                 `json:"user,omitempty"`
	Metadata          map[string]interface{} `json:"-"` // Used for headers
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
	ResponseFormat    *ResponseFormat        `json:"response_format,omitempty"`
	Provider          *Provider              `json:"provider,omitempty"`
	Transforms        []string               `json:"transforms,omitempty"`
	Models            []string               `json:"models,omitempty"`
	Route             string                 `json:"route,omitempty"`
	Plugins           []Plugin               `json:"plugins,omitempty"`
	WebSearchOptions  *WebSearchOptions      `json:"web_search_options,omitempty"`
	User              string                 `json:"user,omitempty"`
	Metadata          map[string]interface{} `json:"-"` // Used for headers
}

// Message represents a message in the chat completion request.
type Message struct {
	Role        string         `json:"role"`
	Content     MessageContent `json:"content"`
	Name        string         `json:"name,omitempty"`
	ToolCalls   []ToolCall     `json:"tool_calls,omitempty"`
	ToolCallID  string         `json:"tool_call_id,omitempty"`
	Annotations []Annotation   `json:"annotations,omitempty"`
}

// MessageContent can be either a string or an array of content parts.
type MessageContent interface{}

// ContentPart represents a part of message content (text, image, file, or audio).
type ContentPart struct {
	Type       string      `json:"type"`
	Text       string      `json:"text,omitempty"`
	ImageURL   *ImageURL   `json:"image_url,omitempty"`
	File       *File       `json:"file,omitempty"`
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

// ImageURL represents an image URL in the message content.
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// File represents a file attachment in the message content.
// Supports both URLs and base64-encoded data URLs.
type File struct {
	// Filename is the name of the file (e.g., "document.pdf")
	Filename string `json:"filename"`
	// FileData can be either a URL or a base64-encoded data URL
	FileData string `json:"file_data"`
}

// InputAudio represents an audio input in the message content.
// Audio files must be base64-encoded - direct URLs are not supported.
type InputAudio struct {
	// Data is the base64-encoded audio data
	Data string `json:"data"`
	// Format specifies the audio format (e.g., "wav", "mp3")
	Format string `json:"format"`
}

// ResponseFormat specifies the format of the response.
type ResponseFormat struct {
	Type       string      `json:"type"`
	JSONSchema *JSONSchema `json:"json_schema,omitempty"`
}

// JSONSchema defines the structure for structured output format.
type JSONSchema struct {
	Name   string                 `json:"name"`
	Strict bool                   `json:"strict"`
	Schema map[string]interface{} `json:"schema"`
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
	Order []string `json:"order,omitempty"`
	// RequireParameters only uses providers that support all parameters in the request
	RequireParameters *bool `json:"require_parameters,omitempty"`
	// DataCollection controls whether to use providers that may store data ("allow" or "deny")
	DataCollection string `json:"data_collection,omitempty"`
	// AllowFallbacks allows backup providers when the primary is unavailable
	AllowFallbacks *bool `json:"allow_fallbacks,omitempty"`
	// Ignore specifies provider slugs to skip for this request
	Ignore []string `json:"ignore,omitempty"`
	// Quantizations filters providers by quantization levels (e.g. ["int4", "int8"])
	Quantizations []string `json:"quantizations,omitempty"`
	// ZDR restricts routing to only Zero Data Retention endpoints
	ZDR *bool `json:"zdr,omitempty"`
	// Only specifies provider slugs to allow for this request
	Only []string `json:"only,omitempty"`
	// Sort providers by "price", "throughput", or "latency"
	Sort string `json:"sort,omitempty"`
	// MaxPrice specifies maximum pricing constraints for the request
	MaxPrice *MaxPrice `json:"max_price,omitempty"`

	// Deprecated: Use Ignore instead
	IgnoreProviders []string `json:"-"`
	// Deprecated: Use Quantizations instead
	QuantizationFallback map[string]string `json:"-"`
	// Internal provider parameters
	ProviderParams map[string]interface{} `json:"-"`
}

// ChatCompletionResponse represents a chat completion response from the OpenRouter API.
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
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
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	FinishReason string    `json:"finish_reason"`
	LogProbs     *LogProbs `json:"logprobs,omitempty"`
	Delta        *Message  `json:"delta,omitempty"` // For streaming
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
	Token       string       `json:"token"`
	LogProb     float64      `json:"logprob"`
	Bytes       []int        `json:"bytes,omitempty"`
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
	ID    string
	Event string
	Data  string
	Retry *time.Duration
}

// ErrorResponse represents an error response from the OpenRouter API.
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// APIError represents the error details in an error response.
type APIError struct {
	Message  string                 `json:"message"`
	Type     string                 `json:"type"`
	Code     string                 `json:"code,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Plugin represents a plugin configuration for enhancing model responses.
type Plugin struct {
	// ID is the plugin identifier (e.g., "web" for web search, "file-parser" for file parsing)
	ID string `json:"id"`
	// Engine specifies which search engine to use ("native", "exa", or undefined for auto)
	Engine string `json:"engine,omitempty"`
	// MaxResults specifies the maximum number of search results (defaults to 5)
	MaxResults int `json:"max_results,omitempty"`
	// SearchPrompt customizes the prompt used to attach search results
	SearchPrompt string `json:"search_prompt,omitempty"`
	// PDF contains configuration for PDF file parsing (when ID is "file-parser")
	PDF *PDFParserConfig `json:"pdf,omitempty"`
}

// PDFParserConfig configures PDF file parsing behavior.
type PDFParserConfig struct {
	// Engine specifies the PDF parsing engine to use.
	// Options:
	// - "pdf-text": Best for well-structured PDFs with clear text content (Free)
	// - "mistral-ocr": Best for scanned documents or PDFs with images ($0.0004 per 1,000 pages)
	// - "native": Only for models with native file support (charged as input tokens)
	// If not specified, defaults to model's native support, then "pdf-text"
	Engine string `json:"engine,omitempty"`
}

// WebSearchOptions configures native web search behavior for supported models.
type WebSearchOptions struct {
	// SearchContextSize determines the amount of search context ("low", "medium", or "high")
	SearchContextSize string `json:"search_context_size,omitempty"`
}

// Annotation represents an annotation in a message response.
type Annotation struct {
	// Type of annotation (e.g., "url_citation", "file")
	Type string `json:"type"`
	// URLCitation contains details for URL citation annotations
	URLCitation *URLCitation `json:"url_citation,omitempty"`
	// FileAnnotation contains details for file annotations (parsed file metadata)
	FileAnnotation *FileAnnotation `json:"file,omitempty"`
}

// URLCitation represents a URL citation in a message annotation.
type URLCitation struct {
	// URL of the cited source
	URL string `json:"url"`
	// Title of the web search result
	Title string `json:"title"`
	// Content of the web search result
	Content string `json:"content,omitempty"`
	// StartIndex is the index of the first character of the citation in the message
	StartIndex int `json:"start_index"`
	// EndIndex is the index of the last character of the citation in the message
	EndIndex int `json:"end_index"`
}

// FileAnnotation represents metadata about a parsed file.
// Can be sent back in subsequent requests to avoid re-parsing costs.
type FileAnnotation struct {
	// Filename is the name of the file
	Filename string `json:"filename"`
	// ContentType is the MIME type of the file (e.g., "application/pdf")
	ContentType string `json:"content_type,omitempty"`
	// ParsedContent contains the extracted text/data from the file
	ParsedContent string `json:"parsed_content,omitempty"`
	// Metadata contains additional file information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ModelsResponse represents the response from the list models endpoint.
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// Model represents a model available on OpenRouter.
type Model struct {
	ID                  string                  `json:"id"`
	Name                string                  `json:"name"`
	CanonicalSlug       *string                 `json:"canonical_slug"`
	Created             float64                 `json:"created"`
	Description         string                  `json:"description"`
	ContextLength       *float64                `json:"context_length"`
	HuggingFaceID       *string                 `json:"hugging_face_id"`
	Architecture        ModelArchitecture       `json:"architecture"`
	TopProvider         ModelTopProvider        `json:"top_provider"`
	PerRequestLimits    *ModelPerRequestLimits  `json:"per_request_limits"`
	SupportedParameters []string                `json:"supported_parameters,omitempty"`
	DefaultParameters   *ModelDefaultParameters `json:"default_parameters"`
	Pricing             ModelPricing            `json:"pricing"`
}

// ModelArchitecture contains information about a model's architecture.
type ModelArchitecture struct {
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer        string   `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
}

// ModelTopProvider contains information about the top provider for a model.
type ModelTopProvider struct {
	ContextLength       *float64 `json:"context_length"`
	MaxCompletionTokens *float64 `json:"max_completion_tokens"`
	IsModerated         bool     `json:"is_moderated"`
}

// ModelPerRequestLimits contains per-request limits for a model.
type ModelPerRequestLimits struct {
	// Currently empty but may be extended in the future
}

// ModelDefaultParameters contains default generation parameters for a model.
type ModelDefaultParameters struct {
	Temperature      *float64 `json:"temperature"`
	TopP             *float64 `json:"top_p"`
	FrequencyPenalty *float64 `json:"frequency_penalty"`
}

// ModelPricing contains pricing information for a model.
type ModelPricing struct {
	Prompt            string  `json:"prompt"`
	Completion        string  `json:"completion"`
	Image             string  `json:"image"`
	Request           string  `json:"request"`
	InputCacheRead    *string `json:"input_cache_read"`
	InputCacheWrite   *string `json:"input_cache_write"`
	WebSearch         string  `json:"web_search"`
	InternalReasoning string  `json:"internal_reasoning"`
}

// ModelEndpointsResponse represents the response from the model endpoints endpoint.
type ModelEndpointsResponse struct {
	Data ModelEndpointsData `json:"data"`
}

// ModelEndpointsData contains details about a model and its endpoints.
type ModelEndpointsData struct {
	ID           string                     `json:"id"`
	Name         string                     `json:"name"`
	Created      float64                    `json:"created"`
	Description  string                     `json:"description"`
	Architecture ModelEndpointsArchitecture `json:"architecture"`
	Endpoints    []ModelEndpoint            `json:"endpoints"`
}

// ModelEndpointsArchitecture contains architecture information for a model's endpoints.
type ModelEndpointsArchitecture struct {
	Tokenizer        *string  `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
}

// ModelEndpoint represents a single endpoint for a model.
type ModelEndpoint struct {
	Name                string               `json:"name"`
	ContextLength       float64              `json:"context_length"`
	Pricing             ModelEndpointPricing `json:"pricing"`
	ProviderName        string               `json:"provider_name"`
	Quantization        *string              `json:"quantization"`
	MaxCompletionTokens *float64             `json:"max_completion_tokens"`
	MaxPromptTokens     *float64             `json:"max_prompt_tokens"`
	SupportedParameters []string             `json:"supported_parameters"`
	Status              float64              `json:"status"`
	UptimeLast30m       *float64             `json:"uptime_last_30m"`
}

// ModelEndpointPricing contains pricing information for a specific endpoint.
type ModelEndpointPricing struct {
	Request    string `json:"request"`
	Image      string `json:"image"`
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

// ProvidersResponse represents the response from the list providers endpoint.
type ProvidersResponse struct {
	Data []ProviderInfo `json:"data"`
}

// ProviderInfo represents information about a provider available on OpenRouter.
type ProviderInfo struct {
	Name              string  `json:"name"`
	Slug              string  `json:"slug"`
	PrivacyPolicyURL  *string `json:"privacy_policy_url"`
	TermsOfServiceURL *string `json:"terms_of_service_url"`
	StatusPageURL     *string `json:"status_page_url"`
}

// CreditsResponse represents the response from the credits endpoint.
type CreditsResponse struct {
	Data CreditsData `json:"data"`
}

// CreditsData contains credit balance information for the authenticated user.
type CreditsData struct {
	TotalCredits float64 `json:"total_credits"`
	TotalUsage   float64 `json:"total_usage"`
}

// ActivityResponse represents the response from the activity endpoint.
type ActivityResponse struct {
	Data []ActivityData `json:"data"`
}

// ActivityData represents daily user activity data grouped by model endpoint.
type ActivityData struct {
	Date               string  `json:"date"`
	Model              string  `json:"model"`
	ModelPermaslug     string  `json:"model_permaslug"`
	EndpointID         string  `json:"endpoint_id"`
	ProviderName       string  `json:"provider_name"`
	Usage              float64 `json:"usage"`
	BYOKUsageInference float64 `json:"byok_usage_inference"`
	Requests           float64 `json:"requests"`
	PromptTokens       float64 `json:"prompt_tokens"`
	CompletionTokens   float64 `json:"completion_tokens"`
	ReasoningTokens    float64 `json:"reasoning_tokens"`
}

// KeyResponse represents the response from the get current API key endpoint.
type KeyResponse struct {
	Data KeyData `json:"data"`
}

// KeyData contains information about an API key.
type KeyData struct {
	Label             string        `json:"label"`
	Limit             *float64      `json:"limit"`
	Usage             float64       `json:"usage"`
	IsFreeTier        bool          `json:"is_free_tier"`
	LimitRemaining    *float64      `json:"limit_remaining"`
	IsProvisioningKey bool          `json:"is_provisioning_key"`
	RateLimit         *KeyRateLimit `json:"rate_limit,omitempty"`
}

// KeyRateLimit contains rate limit information for an API key.
type KeyRateLimit struct {
	Interval string  `json:"interval"`
	Requests float64 `json:"requests"`
}

// ListKeysResponse represents the response from the list API keys endpoint.
type ListKeysResponse struct {
	Data []APIKey `json:"data"`
}

// APIKey represents an API key in the list keys response.
type APIKey struct {
	Name      string  `json:"name"`
	Label     string  `json:"label"`
	Limit     float64 `json:"limit"`
	Disabled  bool    `json:"disabled"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Hash      string  `json:"hash"`
}

// ListKeysOptions contains optional parameters for listing API keys.
type ListKeysOptions struct {
	// Offset for pagination
	Offset *int `json:"offset,omitempty"`
	// IncludeDisabled controls whether to include disabled API keys
	IncludeDisabled *bool `json:"include_disabled,omitempty"`
}

// CreateKeyRequest represents a request to create a new API key.
type CreateKeyRequest struct {
	// Name is the required name/label for the API key
	Name string `json:"name"`
	// Limit is the optional credit limit for the API key
	Limit *float64 `json:"limit,omitempty"`
	// IncludeBYOKInLimit controls whether BYOK usage counts toward the limit
	IncludeBYOKInLimit *bool `json:"include_byok_in_limit,omitempty"`
}

// CreateKeyResponse represents the response from creating a new API key.
type CreateKeyResponse struct {
	Data APIKey `json:"data"`
	// Key contains the actual API key value (only returned on creation)
	Key string `json:"key,omitempty"`
}

// GetKeyByHashResponse represents the response from getting an API key by hash.
type GetKeyByHashResponse struct {
	Data APIKey `json:"data"`
}

// DeleteKeyResponse represents the response from deleting an API key.
type DeleteKeyResponse struct {
	Data DeleteKeyData `json:"data"`
}

// DeleteKeyData contains the result of a delete operation.
type DeleteKeyData struct {
	Success bool `json:"success"`
}

// UpdateKeyRequest represents a request to update an existing API key.
// All fields are optional - only include fields you want to update.
type UpdateKeyRequest struct {
	// Name is the new name/label for the API key
	Name *string `json:"name,omitempty"`
	// Disabled controls whether the API key is disabled
	Disabled *bool `json:"disabled,omitempty"`
	// Limit is the new credit limit for the API key
	Limit *float64 `json:"limit,omitempty"`
	// IncludeBYOKInLimit controls whether BYOK usage counts toward the limit
	IncludeBYOKInLimit *bool `json:"include_byok_in_limit,omitempty"`
}

// UpdateKeyResponse represents the response from updating an API key.
type UpdateKeyResponse struct {
	Data APIKey `json:"data"`
}

// EmbeddingRequest represents a request to create embeddings.
type EmbeddingRequest struct {
	// Input is the text or texts to embed. Can be a string, array of strings,
	// array of integers, array of integer arrays, or array of multimodal content parts.
	Input interface{} `json:"input"`
	// Model specifies the embedding model to use (required).
	Model string `json:"model"`
	// EncodingFormat specifies the format for the returned embeddings ("float" or "base64").
	EncodingFormat string `json:"encoding_format,omitempty"`
	// Dimensions specifies the number of dimensions for the output embeddings.
	Dimensions *int `json:"dimensions,omitempty"`
	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`
	// Provider contains provider-specific routing parameters.
	Provider *Provider `json:"provider,omitempty"`
	// InputType specifies the type of input (e.g., for search vs. document embeddings).
	InputType string `json:"input_type,omitempty"`
}

// EmbeddingResponse represents the response from an embeddings request.
type EmbeddingResponse struct {
	// ID is the unique identifier for this embeddings request.
	ID string `json:"id,omitempty"`
	// Object is the object type, always "list".
	Object string `json:"object"`
	// Data contains the list of embedding objects.
	Data []EmbeddingData `json:"data"`
	// Model is the model used to generate the embeddings.
	Model string `json:"model"`
	// Usage contains token usage information for the request.
	Usage *EmbeddingUsage `json:"usage,omitempty"`
}

// EmbeddingData represents a single embedding in the response.
type EmbeddingData struct {
	// Object is the object type, always "embedding".
	Object string `json:"object"`
	// Embedding is the embedding vector. Can be []float64 or base64 string depending on encoding_format.
	Embedding interface{} `json:"embedding"`
	// Index is the index of this embedding in the input list.
	Index int `json:"index"`
}

// EmbeddingUsage contains token usage information for an embeddings request.
type EmbeddingUsage struct {
	// PromptTokens is the number of tokens in the input.
	PromptTokens int `json:"prompt_tokens"`
	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
	// Cost is the cost of the request (if available).
	Cost *float64 `json:"cost,omitempty"`
}

// EmbeddingContentPart represents a content part for multimodal embedding input.
type EmbeddingContentPart struct {
	// Type specifies the type of content ("text" or "image_url").
	Type string `json:"type"`
	// Text is the text content (when Type is "text").
	Text string `json:"text,omitempty"`
	// ImageURL contains the image URL (when Type is "image_url").
	ImageURL *EmbeddingImageURL `json:"image_url,omitempty"`
}

// EmbeddingImageURL represents an image URL for multimodal embedding input.
type EmbeddingImageURL struct {
	// URL is the image URL.
	URL string `json:"url"`
}

// EmbeddingMultimodalInput represents multimodal input for embeddings.
type EmbeddingMultimodalInput struct {
	// Content is the array of content parts.
	Content []EmbeddingContentPart `json:"content"`
}

// EmbeddingsModelsResponse represents the response from the list embeddings models endpoint.
type EmbeddingsModelsResponse struct {
	Data []Model `json:"data"`
}
