// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

// AnthropicOption is a functional option for configuring AnthropicMessagesRequest.
type AnthropicOption = RequestOption[*AnthropicMessagesRequest]

// WithAnthropicModel sets the model for the Anthropic Messages API request.
func WithAnthropicModel(model string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Model = model
	}
}

// WithAnthropicMaxTokens sets the maximum number of tokens to generate.
func WithAnthropicMaxTokens(maxTokens int) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.MaxTokens = maxTokens
	}
}

// WithAnthropicMessages sets the messages for the request.
func WithAnthropicMessages(messages []AnthropicMessage) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Messages = messages
	}
}

// WithAnthropicSystemString sets a string system prompt.
func WithAnthropicSystemString(system string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.System = system
	}
}

// WithAnthropicSystemBlocks sets a structured system prompt using text blocks.
func WithAnthropicSystemBlocks(blocks []AnthropicTextBlock) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.System = blocks
	}
}

// WithAnthropicTemperature sets the temperature for sampling.
func WithAnthropicTemperature(temperature float64) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Temperature = &temperature
	}
}

// WithAnthropicTopP sets the top_p for nucleus sampling.
func WithAnthropicTopP(topP float64) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.TopP = &topP
	}
}

// WithAnthropicTopK sets the top_k for sampling.
func WithAnthropicTopK(topK int) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.TopK = &topK
	}
}

// WithAnthropicStopSequences sets custom stop sequences.
func WithAnthropicStopSequences(sequences ...string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.StopSequences = sequences
	}
}

// WithAnthropicTools sets the available tools.
func WithAnthropicTools(tools ...AnthropicTool) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Tools = tools
	}
}

// WithAnthropicToolChoiceAuto sets tool choice to "auto".
func WithAnthropicToolChoiceAuto() AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.ToolChoice = &AnthropicToolChoice{Type: AnthropicToolChoiceAuto}
	}
}

// WithAnthropicToolChoiceAny sets tool choice to "any" (force tool use).
func WithAnthropicToolChoiceAny() AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.ToolChoice = &AnthropicToolChoice{Type: AnthropicToolChoiceAny}
	}
}

// WithAnthropicToolChoiceNone sets tool choice to "none".
func WithAnthropicToolChoiceNone() AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.ToolChoice = &AnthropicToolChoice{Type: AnthropicToolChoiceNone}
	}
}

// WithAnthropicToolChoiceSpecific forces the use of a specific tool.
func WithAnthropicToolChoiceSpecific(name string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.ToolChoice = &AnthropicToolChoice{Type: AnthropicToolChoiceTool, Name: name}
	}
}

// WithAnthropicThinkingEnabled enables extended thinking with a token budget.
func WithAnthropicThinkingEnabled(budgetTokens int) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Thinking = &AnthropicThinkingConfig{
			Type:         AnthropicThinkingEnabled,
			BudgetTokens: budgetTokens,
		}
	}
}

// WithAnthropicThinkingDisabled explicitly disables extended thinking.
func WithAnthropicThinkingDisabled() AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Thinking = &AnthropicThinkingConfig{
			Type: AnthropicThinkingDisabled,
		}
	}
}

// WithAnthropicProvider sets provider-specific parameters.
func WithAnthropicProvider(provider Provider) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Provider = &provider
	}
}

// WithAnthropicPlugins sets plugin configurations.
func WithAnthropicPlugins(plugins ...Plugin) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Plugins = plugins
	}
}

// WithAnthropicModels sets fallback models.
func WithAnthropicModels(models ...string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Models = models
	}
}

// WithAnthropicUser sets a unique identifier for the end-user.
func WithAnthropicUser(user string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.User = user
	}
}

// WithAnthropicSessionID sets a session identifier for multi-turn conversations.
func WithAnthropicSessionID(sessionID string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.SessionID = sessionID
	}
}

// WithAnthropicServiceTier sets the service tier.
func WithAnthropicServiceTier(tier string) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.ServiceTier = tier
	}
}

// WithAnthropicRequestMetadata sets the request body metadata.
func WithAnthropicRequestMetadata(metadata AnthropicRequestMetadata) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.Metadata = &metadata
	}
}

// WithAnthropicHeaderMetadata sets metadata for X-* header injection.
func WithAnthropicHeaderMetadata(metadata map[string]any) AnthropicOption {
	return func(r *AnthropicMessagesRequest) {
		r.HeaderMetadata = metadata
	}
}
