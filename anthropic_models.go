// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

import "encoding/json"

// Anthropic Messages API stop reasons.
const (
	AnthropicStopReasonEndTurn      = "end_turn"
	AnthropicStopReasonMaxTokens    = "max_tokens"
	AnthropicStopReasonStopSequence = "stop_sequence"
	AnthropicStopReasonToolUse      = "tool_use"
)

// Anthropic thinking types.
const (
	AnthropicThinkingEnabled  = "enabled"
	AnthropicThinkingDisabled = "disabled"
)

// Anthropic tool choice types.
const (
	AnthropicToolChoiceAuto = "auto"
	AnthropicToolChoiceAny  = "any"
	AnthropicToolChoiceNone = "none"
	AnthropicToolChoiceTool = "tool"
)

// Anthropic service tier values.
const (
	AnthropicServiceTierAuto     = "auto"
	AnthropicServiceTierStandard = "standard"
)

// Anthropic content block types.
const (
	AnthropicContentTypeText             = "text"
	AnthropicContentTypeImage            = "image"
	AnthropicContentTypeDocument         = "document"
	AnthropicContentTypeToolUse          = "tool_use"
	AnthropicContentTypeToolResult       = "tool_result"
	AnthropicContentTypeThinking         = "thinking"
	AnthropicContentTypeRedactedThinking = "redacted_thinking"
)

// Anthropic source types.
const (
	AnthropicSourceBase64 = "base64"
	AnthropicSourceURL    = "url"
)

// AnthropicMessagesRequest represents a request to the Anthropic Messages API via OpenRouter.
type AnthropicMessagesRequest struct {
	// Model is the model identifier (required).
	Model string `json:"model"`

	// MaxTokens is the maximum number of tokens to generate (required).
	MaxTokens int `json:"max_tokens"`

	// Messages is the array of input messages (required).
	Messages []AnthropicMessage `json:"messages"`

	// System is an optional system prompt. Can be a string or []AnthropicTextBlock.
	System interface{} `json:"system,omitempty"`

	// Temperature controls randomness (0-1).
	Temperature *float64 `json:"temperature,omitempty"`

	// TopP controls nucleus sampling.
	TopP *float64 `json:"top_p,omitempty"`

	// TopK controls top-k sampling.
	TopK *int `json:"top_k,omitempty"`

	// StopSequences specifies custom stop sequences.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// Stream enables streaming responses.
	Stream bool `json:"stream,omitempty"`

	// Tools defines available tools.
	Tools []AnthropicTool `json:"tools,omitempty"`

	// ToolChoice controls tool invocation behavior.
	ToolChoice *AnthropicToolChoice `json:"tool_choice,omitempty"`

	// Thinking configures extended thinking.
	Thinking *AnthropicThinkingConfig `json:"thinking,omitempty"`

	// ServiceTier specifies the service tier ("auto" or "standard").
	ServiceTier string `json:"service_tier,omitempty"`

	// Provider contains provider-specific routing parameters (reuses existing type).
	Provider *Provider `json:"provider,omitempty"`

	// Plugins configures plugins (reuses existing type).
	Plugins []Plugin `json:"plugins,omitempty"`

	// User is an optional unique identifier for the end-user.
	User string `json:"user,omitempty"`

	// SessionID is an optional session identifier for multi-turn conversations.
	SessionID string `json:"session_id,omitempty"`

	// Models specifies fallback models.
	Models []string `json:"models,omitempty"`

	// Metadata is the request body metadata (e.g., user_id).
	Metadata *AnthropicRequestMetadata `json:"metadata,omitempty"`

	// HeaderMetadata is used for X-* header injection (not serialized to JSON).
	HeaderMetadata map[string]interface{} `json:"-"`
}

// GetMetadata returns the header metadata map for header generation.
func (r *AnthropicMessagesRequest) GetMetadata() map[string]interface{} {
	return r.HeaderMetadata
}

// AnthropicMessage represents a message in the Anthropic Messages API.
type AnthropicMessage struct {
	// Role is "user" or "assistant".
	Role string `json:"role"`

	// Content can be a string or []AnthropicContentBlock.
	Content interface{} `json:"content"`
}

// AnthropicContentBlock represents a content block in a message.
// It is a unified struct with a Type discriminator covering text, image, document,
// tool_use, tool_result, thinking, and redacted_thinking.
type AnthropicContentBlock struct {
	// Type is the content block type discriminator.
	Type string `json:"type"`

	// Text is the text content (for "text" type).
	Text string `json:"text,omitempty"`

	// Source is the content source (for "image" and "document" types).
	Source *AnthropicContentSource `json:"source,omitempty"`

	// ID is the tool use ID (for "tool_use" type).
	ID string `json:"id,omitempty"`

	// Name is the tool name (for "tool_use" type).
	Name string `json:"name,omitempty"`

	// Input is the tool input (for "tool_use" type).
	Input json.RawMessage `json:"input,omitempty"`

	// ToolUseID is the ID of the tool_use block being responded to (for "tool_result" type).
	ToolUseID string `json:"tool_use_id,omitempty"`

	// Content is the tool result content (for "tool_result" type). Can be string or []AnthropicContentBlock.
	ToolResultContent interface{} `json:"content,omitempty"`

	// IsError indicates if the tool result is an error (for "tool_result" type).
	IsError *bool `json:"is_error,omitempty"`

	// Thinking is the thinking text (for "thinking" type).
	Thinking string `json:"thinking,omitempty"`

	// Signature is the thinking block signature (for "thinking" type).
	Signature string `json:"signature,omitempty"`

	// Data is redacted data (for "redacted_thinking" type).
	Data string `json:"data,omitempty"`

	// Citations configuration (for "text" type).
	Citations *AnthropicCitationsConfig `json:"citations,omitempty"`
}

// AnthropicContentSource represents a source for image or document content blocks.
type AnthropicContentSource struct {
	// Type is "base64" or "url".
	Type string `json:"type"`

	// MediaType is the MIME type (for base64 sources).
	MediaType string `json:"media_type,omitempty"`

	// Data is the base64-encoded data (for base64 sources).
	Data string `json:"data,omitempty"`

	// URL is the source URL (for URL sources).
	URL string `json:"url,omitempty"`
}

// AnthropicTextBlock represents a text block used in system prompts.
type AnthropicTextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// AnthropicTool represents a tool definition for the Anthropic API.
type AnthropicTool struct {
	// Type is the tool type. For custom tools, omit or set to empty.
	// Built-in types: "bash_20250124", "text_editor_20250124", "web_search_20250305"
	Type string `json:"type,omitempty"`

	// Name is the tool name (required for custom tools).
	Name string `json:"name,omitempty"`

	// Description explains what the tool does.
	Description string `json:"description,omitempty"`

	// InputSchema is the JSON Schema for the tool parameters (custom tools).
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`

	// MaxUses limits the number of times the tool can be used (for web_search).
	MaxUses *int `json:"max_uses,omitempty"`
}

// AnthropicToolChoice controls tool invocation behavior.
type AnthropicToolChoice struct {
	// Type is "auto", "any", "none", or "tool".
	Type string `json:"type"`

	// Name is the tool name (only when Type is "tool").
	Name string `json:"name,omitempty"`

	// DisableParallelToolUse disables parallel tool calling.
	DisableParallelToolUse *bool `json:"disable_parallel_tool_use,omitempty"`
}

// AnthropicThinkingConfig configures extended thinking.
type AnthropicThinkingConfig struct {
	// Type is "enabled" or "disabled".
	Type string `json:"type"`

	// BudgetTokens is the maximum number of thinking tokens (required when Type is "enabled").
	BudgetTokens int `json:"budget_tokens,omitempty"`
}

// AnthropicRequestMetadata contains body-level metadata for the request.
type AnthropicRequestMetadata struct {
	// UserID is an external identifier for the user making the request.
	UserID string `json:"user_id,omitempty"`
}

// AnthropicCitationsConfig controls citation generation.
type AnthropicCitationsConfig struct {
	// Enabled enables citation generation.
	Enabled bool `json:"enabled"`
}

// --- Response Types ---

// AnthropicMessagesResponse represents a response from the Anthropic Messages API.
type AnthropicMessagesResponse struct {
	// ID is the unique message identifier.
	ID string `json:"id"`

	// Type is always "message".
	Type string `json:"type"`

	// Role is always "assistant".
	Role string `json:"role"`

	// Content is the array of content blocks in the response.
	Content []AnthropicResponseContentBlock `json:"content"`

	// Model is the model that generated the response.
	Model string `json:"model"`

	// StopReason indicates why the model stopped generating.
	StopReason *string `json:"stop_reason"`

	// StopSequence is the stop sequence that triggered stopping, if any.
	StopSequence *string `json:"stop_sequence"`

	// Usage contains token usage information.
	Usage AnthropicUsage `json:"usage"`
}

// AnthropicResponseContentBlock represents a content block in the response.
type AnthropicResponseContentBlock struct {
	// Type discriminator: "text", "tool_use", "thinking", "redacted_thinking",
	// "server_tool_use", "web_search_tool_result"
	Type string `json:"type"`

	// Text is the text content (for "text" type).
	Text string `json:"text,omitempty"`

	// Citations for the text block.
	Citations []AnthropicCitation `json:"citations,omitempty"`

	// ID is the tool use ID (for "tool_use" and "server_tool_use" types).
	ID string `json:"id,omitempty"`

	// Name is the tool name (for "tool_use" and "server_tool_use" types).
	Name string `json:"name,omitempty"`

	// Input is the tool input (for "tool_use" and "server_tool_use" types).
	Input json.RawMessage `json:"input,omitempty"`

	// Thinking is the thinking text (for "thinking" type).
	Thinking string `json:"thinking,omitempty"`

	// Signature is the thinking block signature (for "thinking" type).
	Signature string `json:"signature,omitempty"`

	// Data is redacted data (for "redacted_thinking" type).
	Data string `json:"data,omitempty"`

	// Content is the web search result content (for "web_search_tool_result" type).
	Content json.RawMessage `json:"content,omitempty"`
}

// AnthropicCitation represents a citation in a text content block.
type AnthropicCitation struct {
	// Type: "char_location", "page_location", "content_block_location",
	// "web_search_result_location", "search_result_location"
	Type string `json:"type"`

	// CitedText is the text being cited.
	CitedText string `json:"cited_text,omitempty"`

	// DocumentIndex is the index of the source document.
	DocumentIndex *int `json:"document_index,omitempty"`

	// DocumentTitle is the title of the source document.
	DocumentTitle string `json:"document_title,omitempty"`

	// StartCharIndex is the start character index (for "char_location").
	StartCharIndex *int `json:"start_char_index,omitempty"`

	// EndCharIndex is the end character index (for "char_location").
	EndCharIndex *int `json:"end_char_index,omitempty"`

	// StartPageNumber is the start page (for "page_location").
	StartPageNumber *int `json:"start_page_number,omitempty"`

	// EndPageNumber is the end page (for "page_location").
	EndPageNumber *int `json:"end_page_number,omitempty"`

	// StartBlockIndex is the start block index (for "content_block_location").
	StartBlockIndex *int `json:"start_block_index,omitempty"`

	// EndBlockIndex is the end block index (for "content_block_location").
	EndBlockIndex *int `json:"end_block_index,omitempty"`

	// URL is the source URL (for "web_search_result_location").
	URL string `json:"url,omitempty"`

	// Title is the source title (for "web_search_result_location").
	Title string `json:"title,omitempty"`
}

// AnthropicUsage contains token usage information.
type AnthropicUsage struct {
	InputTokens              int    `json:"input_tokens"`
	OutputTokens             int    `json:"output_tokens"`
	CacheCreationInputTokens int    `json:"cache_creation_input_tokens,omitempty"`
	CacheReadInputTokens     int    `json:"cache_read_input_tokens,omitempty"`
	ServiceTier              string `json:"service_tier,omitempty"`
}

// --- Streaming Types ---

// AnthropicStream represents a streaming response from the Anthropic Messages API.
type AnthropicStream = Stream[AnthropicStreamEvent]

// AnthropicStreamEvent represents a single event in the Anthropic streaming response.
type AnthropicStreamEvent struct {
	// Type: "message_start", "content_block_start", "content_block_delta",
	// "content_block_stop", "message_delta", "message_stop", "ping", "error"
	Type string `json:"type"`

	// Message is the full message object (for "message_start").
	Message *AnthropicMessagesResponse `json:"message,omitempty"`

	// Index is the content block index (for content_block_* events).
	Index *int `json:"index,omitempty"`

	// ContentBlock is the content block being started (for "content_block_start").
	ContentBlock *AnthropicResponseContentBlock `json:"content_block,omitempty"`

	// Delta contains incremental content (for "content_block_delta" and "message_delta").
	Delta *AnthropicStreamDelta `json:"delta,omitempty"`

	// Usage contains updated usage info (for "message_delta").
	Usage *AnthropicStreamUsage `json:"usage,omitempty"`

	// Error contains error information (for "error" events).
	Error *AnthropicStreamError `json:"error,omitempty"`
}

// AnthropicStreamDelta represents incremental content in a streaming event.
type AnthropicStreamDelta struct {
	// Type: "text_delta", "input_json_delta", "thinking_delta", "signature_delta"
	Type string `json:"type,omitempty"`

	// Text is the incremental text (for "text_delta").
	Text string `json:"text,omitempty"`

	// PartialJSON is the incremental JSON (for "input_json_delta").
	PartialJSON string `json:"partial_json,omitempty"`

	// Thinking is the incremental thinking text (for "thinking_delta").
	Thinking string `json:"thinking,omitempty"`

	// Signature is the signature (for "signature_delta").
	Signature string `json:"signature,omitempty"`

	// StopReason is set on "message_delta" events.
	StopReason string `json:"stop_reason,omitempty"`

	// StopSequence is set on "message_delta" events.
	StopSequence string `json:"stop_sequence,omitempty"`
}

// AnthropicStreamUsage contains usage information in streaming events.
type AnthropicStreamUsage struct {
	InputTokens  int `json:"input_tokens,omitempty"`
	OutputTokens int `json:"output_tokens,omitempty"`
}

// AnthropicStreamError contains error information in streaming events.
type AnthropicStreamError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// --- Helper Constructors ---

// CreateAnthropicUserMessage creates a user message with text content.
func CreateAnthropicUserMessage(text string) AnthropicMessage {
	return AnthropicMessage{
		Role:    "user",
		Content: text,
	}
}

// CreateAnthropicAssistantMessage creates an assistant message with text content.
func CreateAnthropicAssistantMessage(text string) AnthropicMessage {
	return AnthropicMessage{
		Role:    "assistant",
		Content: text,
	}
}

// CreateAnthropicUserMessageWithBlocks creates a user message with content blocks.
func CreateAnthropicUserMessageWithBlocks(blocks []AnthropicContentBlock) AnthropicMessage {
	return AnthropicMessage{
		Role:    "user",
		Content: blocks,
	}
}

// CreateAnthropicTextBlock creates a text content block.
func CreateAnthropicTextBlock(text string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type: AnthropicContentTypeText,
		Text: text,
	}
}

// CreateAnthropicImageURLBlock creates an image content block from a URL.
func CreateAnthropicImageURLBlock(url string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type: AnthropicContentTypeImage,
		Source: &AnthropicContentSource{
			Type: AnthropicSourceURL,
			URL:  url,
		},
	}
}

// CreateAnthropicImageBase64Block creates an image content block from base64 data.
func CreateAnthropicImageBase64Block(mediaType string, data string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type: AnthropicContentTypeImage,
		Source: &AnthropicContentSource{
			Type:      AnthropicSourceBase64,
			MediaType: mediaType,
			Data:      data,
		},
	}
}

// CreateAnthropicDocumentURLBlock creates a document content block from a URL.
func CreateAnthropicDocumentURLBlock(url string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type: AnthropicContentTypeDocument,
		Source: &AnthropicContentSource{
			Type: AnthropicSourceURL,
			URL:  url,
		},
	}
}

// CreateAnthropicDocumentBase64Block creates a document content block from base64 data.
func CreateAnthropicDocumentBase64Block(mediaType string, data string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type: AnthropicContentTypeDocument,
		Source: &AnthropicContentSource{
			Type:      AnthropicSourceBase64,
			MediaType: mediaType,
			Data:      data,
		},
	}
}

// CreateAnthropicToolUseBlock creates a tool_use content block.
func CreateAnthropicToolUseBlock(id, name string, input json.RawMessage) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type:  AnthropicContentTypeToolUse,
		ID:    id,
		Name:  name,
		Input: input,
	}
}

// CreateAnthropicToolResultBlock creates a tool_result content block.
func CreateAnthropicToolResultBlock(toolUseID string, content string) AnthropicContentBlock {
	return AnthropicContentBlock{
		Type:              AnthropicContentTypeToolResult,
		ToolUseID:         toolUseID,
		ToolResultContent: content,
	}
}

// CreateAnthropicCustomTool creates a custom tool definition.
func CreateAnthropicCustomTool(name, description string, schema map[string]interface{}) AnthropicTool {
	return AnthropicTool{
		Name:        name,
		Description: description,
		InputSchema: schema,
	}
}

// CreateAnthropicBashTool creates a bash tool definition.
func CreateAnthropicBashTool() AnthropicTool {
	return AnthropicTool{
		Type: "bash_20250124",
		Name: "bash",
	}
}

// CreateAnthropicTextEditorTool creates a text editor tool definition.
func CreateAnthropicTextEditorTool() AnthropicTool {
	return AnthropicTool{
		Type: "text_editor_20250124",
		Name: "str_replace_editor",
	}
}

// CreateAnthropicWebSearchTool creates a web search tool definition.
func CreateAnthropicWebSearchTool(maxUses *int) AnthropicTool {
	return AnthropicTool{
		Type:    "web_search_20250305",
		Name:    "web_search",
		MaxUses: maxUses,
	}
}

// --- Response Helper Methods ---

// GetTextContent extracts all text content from the response, concatenated.
func (r *AnthropicMessagesResponse) GetTextContent() string {
	var result string
	for _, block := range r.Content {
		if block.Type == AnthropicContentTypeText {
			result += block.Text
		}
	}
	return result
}

// GetTextBlocks returns all text content blocks from the response.
func (r *AnthropicMessagesResponse) GetTextBlocks() []AnthropicResponseContentBlock {
	var blocks []AnthropicResponseContentBlock
	for _, block := range r.Content {
		if block.Type == AnthropicContentTypeText {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// GetToolUseBlocks returns all tool_use content blocks from the response.
func (r *AnthropicMessagesResponse) GetToolUseBlocks() []AnthropicResponseContentBlock {
	var blocks []AnthropicResponseContentBlock
	for _, block := range r.Content {
		if block.Type == AnthropicContentTypeToolUse {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// GetThinkingContent extracts all thinking text from the response, concatenated.
func (r *AnthropicMessagesResponse) GetThinkingContent() string {
	var result string
	for _, block := range r.Content {
		if block.Type == AnthropicContentTypeThinking {
			result += block.Thinking
		}
	}
	return result
}

// IsToolUse returns true if the response stop reason is "tool_use".
func (r *AnthropicMessagesResponse) IsToolUse() bool {
	return r.GetStopReason() == AnthropicStopReasonToolUse
}

// GetStopReason returns the stop reason string, or empty if nil.
func (r *AnthropicMessagesResponse) GetStopReason() string {
	if r.StopReason != nil {
		return *r.StopReason
	}
	return ""
}

// --- Stream Event Helper Methods ---

// GetTextDelta extracts the text delta from a streaming event.
func (e *AnthropicStreamEvent) GetTextDelta() string {
	if e.Delta != nil && e.Delta.Type == "text_delta" {
		return e.Delta.Text
	}
	return ""
}
