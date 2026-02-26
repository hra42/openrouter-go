package openrouter

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
	ExpirationDate      *string                 `json:"expiration_date,omitempty"`
}

// ModelArchitecture contains information about a model's architecture.
type ModelArchitecture struct {
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer        string   `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
	Modality         *string  `json:"modality"`
}

// ModelTopProvider contains information about the top provider for a model.
type ModelTopProvider struct {
	ContextLength       *float64 `json:"context_length"`
	MaxCompletionTokens *float64 `json:"max_completion_tokens"`
	IsModerated         bool     `json:"is_moderated"`
}

// ModelPerRequestLimits contains per-request limits for a model.
type ModelPerRequestLimits struct {
	PromptTokens     *float64 `json:"prompt_tokens"`
	CompletionTokens *float64 `json:"completion_tokens"`
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

// PercentileStats contains latency or throughput percentile statistics.
type PercentileStats struct {
	P50 float64 `json:"p50"`
	P75 float64 `json:"p75"`
	P90 float64 `json:"p90"`
	P99 float64 `json:"p99"`
}

// PublicEndpointPricing contains pricing information for a public endpoint.
type PublicEndpointPricing struct {
	Prompt            string `json:"prompt"`
	Completion        string `json:"completion"`
	Request           string `json:"request,omitempty"`
	Image             string `json:"image,omitempty"`
	ImageToken        string `json:"image_token,omitempty"`
	ImageOutput       string `json:"image_output,omitempty"`
	Audio             string `json:"audio,omitempty"`
	AudioOutput       string `json:"audio_output,omitempty"`
	InputAudioCache   string `json:"input_audio_cache,omitempty"`
	WebSearch         string `json:"web_search,omitempty"`
	InternalReasoning string `json:"internal_reasoning,omitempty"`
	InputCacheRead    string `json:"input_cache_read,omitempty"`
	InputCacheWrite   string `json:"input_cache_write,omitempty"`
	Discount          string `json:"discount,omitempty"`
}

// PublicEndpoint represents a single endpoint from the ZDR endpoints listing.
type PublicEndpoint struct {
	Name                    string                `json:"name"`
	ModelID                 string                `json:"model_id"`
	ModelName               string                `json:"model_name"`
	ContextLength           float64               `json:"context_length"`
	Pricing                 PublicEndpointPricing `json:"pricing"`
	ProviderName            string                `json:"provider_name"`
	Tag                     *string               `json:"tag"`
	Quantization            *string               `json:"quantization"`
	MaxCompletionTokens     *float64              `json:"max_completion_tokens"`
	MaxPromptTokens         *float64              `json:"max_prompt_tokens"`
	SupportedParameters     []string              `json:"supported_parameters"`
	Status                  float64               `json:"status"`
	UptimeLast30m           *float64              `json:"uptime_last_30m"`
	SupportsImplicitCaching *bool                 `json:"supports_implicit_caching"`
	LatencyLast30m          *PercentileStats      `json:"latency_last_30m"`
	ThroughputLast30m       *PercentileStats      `json:"throughput_last_30m"`
}

// ZDREndpointsResponse represents the response from the ZDR endpoints listing.
type ZDREndpointsResponse struct {
	Data []PublicEndpoint `json:"data"`
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
