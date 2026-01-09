// Package openrouter provides Go bindings for the OpenRouter API.
package openrouter

// ResponsesOption is a functional option for configuring ResponsesRequest.
type ResponsesOption func(*ResponsesRequest)

// WithResponsesModel sets the model for the Responses API request.
func WithResponsesModel(model string) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Model = model
	}
}

// WithResponsesInput sets the input for the Responses API request.
// Input can be a string for simple text, or []ResponsesInputItem for structured messages.
func WithResponsesInput(input interface{}) ResponsesOption {
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
func WithResponsesReasoning(reasoning ResponsesReasoning) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Reasoning = &reasoning
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
func WithResponsesTools(tools ...Tool) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Tools = tools
	}
}

// WithResponsesToolChoice sets the tool choice strategy.
// Can be "auto", "none", or a specific tool: map[string]interface{}{"type": "function", "name": "functionName"}
func WithResponsesToolChoice(toolChoice interface{}) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.ToolChoice = toolChoice
	}
}

// WithResponsesPlugins sets the plugins for the request.
func WithResponsesPlugins(plugins ...Plugin) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Plugins = plugins
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
func WithResponsesMetadata(metadata map[string]interface{}) ResponsesOption {
	return func(r *ResponsesRequest) {
		r.Metadata = metadata
	}
}

// Reasoning effort level constants
const (
	ReasoningEffortMinimal = "minimal"
	ReasoningEffortLow     = "low"
	ReasoningEffortMedium  = "medium"
	ReasoningEffortHigh    = "high"
)
