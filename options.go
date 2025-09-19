package openrouter

import (
	"net/http"
	"time"
)

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL for the API.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithDefaultModel sets a default model to use for requests.
func WithDefaultModel(model string) ClientOption {
	return func(c *Client) {
		c.defaultModel = model
	}
}

// WithReferer sets the HTTP-Referer header for requests.
func WithReferer(referer string) ClientOption {
	return func(c *Client) {
		c.referer = referer
	}
}

// WithAppName sets the X-Title header for requests.
func WithAppName(appName string) ClientOption {
	return func(c *Client) {
		c.appName = appName
	}
}

// WithRetry configures retry behavior.
func WithRetry(maxRetries int, retryDelay time.Duration) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
		c.retryDelay = retryDelay
	}
}

// WithHeader adds a custom header to all requests.
func WithHeader(key, value string) ClientOption {
	return func(c *Client) {
		c.customHeaders[key] = value
	}
}

// ChatCompletionOption is a functional option for chat completion requests.
type ChatCompletionOption func(*ChatCompletionRequest)

// WithModel sets the model for the request.
func WithModel(model string) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Model = model
	}
}

// WithTemperature sets the temperature parameter.
func WithTemperature(temperature float64) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Temperature = &temperature
	}
}

// WithTopP sets the top_p parameter.
func WithTopP(topP float64) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.TopP = &topP
	}
}

// WithTopK sets the top_k parameter.
func WithTopK(topK int) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.TopK = &topK
	}
}

// WithMaxTokens sets the max_tokens parameter.
func WithMaxTokens(maxTokens int) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.MaxTokens = &maxTokens
	}
}

// WithStop sets the stop sequences.
func WithStop(stop ...string) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Stop = stop
	}
}

// WithFrequencyPenalty sets the frequency penalty.
func WithFrequencyPenalty(penalty float64) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.FrequencyPenalty = &penalty
	}
}

// WithPresencePenalty sets the presence penalty.
func WithPresencePenalty(penalty float64) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.PresencePenalty = &penalty
	}
}

// WithRepetitionPenalty sets the repetition penalty.
func WithRepetitionPenalty(penalty float64) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.RepetitionPenalty = &penalty
	}
}

// WithSeed sets the random seed.
func WithSeed(seed int) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Seed = &seed
	}
}

// WithTools sets the available tools/functions.
func WithTools(tools ...Tool) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Tools = tools
	}
}

// WithToolChoice sets the tool choice strategy.
func WithToolChoice(toolChoice interface{}) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.ToolChoice = toolChoice
	}
}

// WithResponseFormat sets the response format.
func WithResponseFormat(format ResponseFormat) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.ResponseFormat = &format
	}
}

// WithLogProbs enables log probabilities in the response.
func WithLogProbs(topLogProbs int) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		enabled := true
		r.LogProbs = &enabled
		if topLogProbs > 0 {
			r.TopLogProbs = &topLogProbs
		}
	}
}

// WithProvider sets provider-specific parameters.
func WithProvider(provider Provider) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Provider = &provider
	}
}

// WithTransforms sets the transforms to apply.
func WithTransforms(transforms ...string) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Transforms = transforms
	}
}

// WithModels sets the models for fallback.
func WithModels(models ...string) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Models = models
	}
}

// WithRoute sets the routing preference.
func WithRoute(route string) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Route = route
	}
}

// WithMetadata sets metadata headers for the request.
func WithMetadata(metadata map[string]interface{}) ChatCompletionOption {
	return func(r *ChatCompletionRequest) {
		r.Metadata = metadata
	}
}

// CompletionOption is a functional option for completion requests.
type CompletionOption func(*CompletionRequest)

// WithCompletionModel sets the model for the completion request.
func WithCompletionModel(model string) CompletionOption {
	return func(r *CompletionRequest) {
		r.Model = model
	}
}

// WithCompletionTemperature sets the temperature for completion.
func WithCompletionTemperature(temperature float64) CompletionOption {
	return func(r *CompletionRequest) {
		r.Temperature = &temperature
	}
}

// WithCompletionTopP sets the top_p for completion.
func WithCompletionTopP(topP float64) CompletionOption {
	return func(r *CompletionRequest) {
		r.TopP = &topP
	}
}

// WithCompletionMaxTokens sets the max_tokens for completion.
func WithCompletionMaxTokens(maxTokens int) CompletionOption {
	return func(r *CompletionRequest) {
		r.MaxTokens = &maxTokens
	}
}

// WithCompletionStop sets stop sequences for completion.
func WithCompletionStop(stop ...string) CompletionOption {
	return func(r *CompletionRequest) {
		r.Stop = stop
	}
}

// WithCompletionLogProbs sets the number of log probabilities to return.
func WithCompletionLogProbs(logProbs int) CompletionOption {
	return func(r *CompletionRequest) {
		r.LogProbs = &logProbs
	}
}

// WithCompletionEcho enables echoing the prompt in the response.
func WithCompletionEcho(echo bool) CompletionOption {
	return func(r *CompletionRequest) {
		r.Echo = &echo
	}
}

// WithCompletionN sets the number of completions to generate.
func WithCompletionN(n int) CompletionOption {
	return func(r *CompletionRequest) {
		r.N = &n
	}
}

// WithCompletionBestOf sets the number of completions to generate server-side.
func WithCompletionBestOf(bestOf int) CompletionOption {
	return func(r *CompletionRequest) {
		r.BestOf = &bestOf
	}
}

// WithCompletionSuffix sets the suffix to append after the completion.
func WithCompletionSuffix(suffix string) CompletionOption {
	return func(r *CompletionRequest) {
		r.Suffix = &suffix
	}
}

// WithCompletionProvider sets provider-specific parameters for completion.
func WithCompletionProvider(provider Provider) CompletionOption {
	return func(r *CompletionRequest) {
		r.Provider = &provider
	}
}

// WithCompletionMetadata sets metadata headers for the completion request.
func WithCompletionMetadata(metadata map[string]interface{}) CompletionOption {
	return func(r *CompletionRequest) {
		r.Metadata = metadata
	}
}