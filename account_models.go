package openrouter

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

// CodeChallengeMethod represents the PKCE code challenge method.
type CodeChallengeMethod string

const (
	// CodeChallengeMethodS256 uses SHA-256 hashing per RFC 7636.
	CodeChallengeMethodS256 CodeChallengeMethod = "S256"
	// CodeChallengeMethodPlain uses the plain verifier as the challenge.
	CodeChallengeMethodPlain CodeChallengeMethod = "plain"
)

// ExchangeAuthCodeRequest represents a request to exchange an authorization code for an API key.
type ExchangeAuthCodeRequest struct {
	Code                string              `json:"code"`
	CodeVerifier        *string             `json:"code_verifier,omitempty"`
	CodeChallengeMethod CodeChallengeMethod `json:"code_challenge_method,omitempty"`
}

// ExchangeAuthCodeResponse represents the response from exchanging an authorization code.
type ExchangeAuthCodeResponse struct {
	Key    string  `json:"key"`
	UserID *string `json:"user_id"`
}

