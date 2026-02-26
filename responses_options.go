// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

// RequestOption is a generic functional option type for configuring requests.
type RequestOption[T any] func(T)

// ResponsesOption is a functional option for configuring ResponsesRequest.
type ResponsesOption = RequestOption[*ResponsesRequest]

// WithResponsesModel sets the model for the Responses API request.
func WithResponsesModel(model string) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Model = model
	}
}

// WithResponsesInput sets the input for the Responses API request.
// Input can be a string for simple text, or []ResponsesInputItem for structured messages.
func WithResponsesInput(input any) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Input = input
	}
}

// WithResponsesMaxOutputTokens sets the maximum number of output tokens.
func WithResponsesMaxOutputTokens(tokens int) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.MaxOutputTokens = &tokens
	}
}

// WithResponsesTemperature sets the temperature for sampling (0-2).
func WithResponsesTemperature(temperature float64) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Temperature = &temperature
	}
}

// WithResponsesTopP sets the top_p value for nucleus sampling (0-1).
func WithResponsesTopP(topP float64) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.TopP = &topP
	}
}

// WithResponsesReasoning sets the reasoning configuration.
// If ptr is nil, r.Reasoning is left nil. Otherwise, a defensive copy is made.
func WithResponsesReasoning(ptr *ResponsesReasoning) ResponsesOption {
	return func(r *ResponsesRequest) {
		if ptr == nil {
			r.Reasoning = nil
			return
		}
		newReasoning := &ResponsesReasoning{}
		*newReasoning = *ptr
		r.Reasoning = newReasoning
	}
}

// WithResponsesReasoningEffort sets the reasoning effort level.
// Valid values: "minimal", "low", "medium", "high"
func WithResponsesReasoningEffort(effort string) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Reasoning = &ResponsesReasoning{Effort: effort}
	}
}

// WithResponsesTools sets the available tools/functions for the request.
// Uses ResponsesTool which has a flat structure expected by the Responses API.
func WithResponsesTools(tools ...ResponsesTool) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Tools = tools
	}
}

// WithResponsesToolChoice sets the tool choice strategy for ResponsesRequest.ToolChoice.
// Can be:
//   - "auto" - let the model decide whether to call tools
//   - "none" - disable tool calling
//   - A specific tool: map[string]any{"type": "function", "function": map[string]any{"name": "functionName"}}
func WithResponsesToolChoice(toolChoice any) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.ToolChoice = toolChoice
	}
}

// WithResponsesPlugins appends plugins to the request.
// This appends to any existing plugins to avoid overwriting plugins added by
// other options like WithResponsesWebSearch.
func WithResponsesPlugins(plugins ...Plugin) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Plugins = append(r.Plugins, plugins...)
	}
}

// WithResponsesWebSearch enables web search with the specified maximum results.
// This is a convenience function that adds the web search plugin.
func WithResponsesWebSearch(maxResults int) ResponsesOption {
	return func(r *ResponsesRequest) {
		webPlugin := Plugin{
			ID:         "web",
			MaxResults: maxResults,
		}
		r.Plugins = append(r.Plugins, webPlugin)
	}
}

// WithResponsesMetadata sets metadata headers for the request.
func WithResponsesMetadata(metadata map[string]any) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Metadata = metadata
	}
}

// Reasoning effort level constants for the Responses API.
const (
	// ReasoningEffortMinimal indicates minimal reasoning effort for simple queries.
	ReasoningEffortMinimal = "minimal"
	// ReasoningEffortLow indicates low reasoning effort for straightforward tasks.
	ReasoningEffortLow = "low"
	// ReasoningEffortMedium indicates medium reasoning effort for moderate complexity.
	ReasoningEffortMedium = "medium"
	// ReasoningEffortHigh indicates high reasoning effort for complex problems.
	ReasoningEffortHigh = "high"
)
