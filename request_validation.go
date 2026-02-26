package openrouter

// validateChatCompletionParams validates chat completion request parameters
// after options have been applied. Returns a ValidationError if any parameter
// is invalid.
func validateChatCompletionParams(req *ChatCompletionRequest) error {
	// ToolChoice set but no tools
	if req.ToolChoice != nil && len(req.Tools) == 0 {
		return &ValidationError{Field: "tool_choice", Message: "tool_choice is set but no tools are provided"}
	}

	// ParallelToolCalls set but no tools
	if req.ParallelToolCalls != nil && len(req.Tools) == 0 {
		return &ValidationError{Field: "parallel_tool_calls", Message: "parallel_tool_calls is set but no tools are provided"}
	}

	// TopLogProbs > 0 but LogProbs not enabled
	if req.TopLogProbs != nil && *req.TopLogProbs > 0 {
		if req.LogProbs == nil || !*req.LogProbs {
			return &ValidationError{Field: "top_logprobs", Message: "top_logprobs requires logprobs to be enabled"}
		}
	}

	// Temperature range
	if req.Temperature != nil && *req.Temperature < 0 {
		return &ValidationError{Field: "temperature", Message: "temperature must be non-negative"}
	}

	// TopP range [0, 1]
	if req.TopP != nil && (*req.TopP < 0 || *req.TopP > 1) {
		return &ValidationError{Field: "top_p", Message: "top_p must be between 0 and 1"}
	}

	// MinP range [0, 1]
	if req.MinP != nil && (*req.MinP < 0 || *req.MinP > 1) {
		return &ValidationError{Field: "min_p", Message: "min_p must be between 0 and 1"}
	}

	// TopK >= 0
	if req.TopK != nil && *req.TopK < 0 {
		return &ValidationError{Field: "top_k", Message: "top_k must be non-negative"}
	}

	// MaxTokens > 0
	if req.MaxTokens != nil && *req.MaxTokens <= 0 {
		return &ValidationError{Field: "max_tokens", Message: "max_tokens must be positive"}
	}

	// FrequencyPenalty [-2, 2]
	if req.FrequencyPenalty != nil && (*req.FrequencyPenalty < -2 || *req.FrequencyPenalty > 2) {
		return &ValidationError{Field: "frequency_penalty", Message: "frequency_penalty must be between -2 and 2"}
	}

	// PresencePenalty [-2, 2]
	if req.PresencePenalty != nil && (*req.PresencePenalty < -2 || *req.PresencePenalty > 2) {
		return &ValidationError{Field: "presence_penalty", Message: "presence_penalty must be between -2 and 2"}
	}

	// TopLogProbs [0, 20]
	if req.TopLogProbs != nil && (*req.TopLogProbs < 0 || *req.TopLogProbs > 20) {
		return &ValidationError{Field: "top_logprobs", Message: "top_logprobs must be between 0 and 20"}
	}

	// Provider.Sort validation
	if req.Provider != nil && req.Provider.Sort != "" {
		validSorts := map[string]bool{"price": true, "throughput": true, "latency": true}
		if !validSorts[req.Provider.Sort] {
			return &ValidationError{Field: "provider.sort", Message: "sort must be one of: price, throughput, latency"}
		}
	}

	return nil
}

// validateCompletionParams validates legacy completion request parameters
// after options have been applied. Returns a ValidationError if any parameter
// is invalid.
func validateCompletionParams(req *CompletionRequest) error {
	// Temperature range
	if req.Temperature != nil && *req.Temperature < 0 {
		return &ValidationError{Field: "temperature", Message: "temperature must be non-negative"}
	}

	// TopP range [0, 1]
	if req.TopP != nil && (*req.TopP < 0 || *req.TopP > 1) {
		return &ValidationError{Field: "top_p", Message: "top_p must be between 0 and 1"}
	}

	// MinP range [0, 1]
	if req.MinP != nil && (*req.MinP < 0 || *req.MinP > 1) {
		return &ValidationError{Field: "min_p", Message: "min_p must be between 0 and 1"}
	}

	// TopK >= 0
	if req.TopK != nil && *req.TopK < 0 {
		return &ValidationError{Field: "top_k", Message: "top_k must be non-negative"}
	}

	// MaxTokens > 0
	if req.MaxTokens != nil && *req.MaxTokens <= 0 {
		return &ValidationError{Field: "max_tokens", Message: "max_tokens must be positive"}
	}

	// FrequencyPenalty [-2, 2]
	if req.FrequencyPenalty != nil && (*req.FrequencyPenalty < -2 || *req.FrequencyPenalty > 2) {
		return &ValidationError{Field: "frequency_penalty", Message: "frequency_penalty must be between -2 and 2"}
	}

	// PresencePenalty [-2, 2]
	if req.PresencePenalty != nil && (*req.PresencePenalty < -2 || *req.PresencePenalty > 2) {
		return &ValidationError{Field: "presence_penalty", Message: "presence_penalty must be between -2 and 2"}
	}

	// BestOf >= N (if N is set)
	if req.BestOf != nil && req.N != nil && *req.BestOf < *req.N {
		return &ValidationError{Field: "best_of", Message: "best_of must be greater than or equal to n"}
	}

	// Provider.Sort validation
	if req.Provider != nil && req.Provider.Sort != "" {
		validSorts := map[string]bool{"price": true, "throughput": true, "latency": true}
		if !validSorts[req.Provider.Sort] {
			return &ValidationError{Field: "provider.sort", Message: "sort must be one of: price, throughput, latency"}
		}
	}

	return nil
}
