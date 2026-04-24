package openrouter

// Audio response formats for the Create Speech endpoint.
const (
	SpeechFormatMP3 = "mp3"
	SpeechFormatPCM = "pcm"
)

// SpeechRequest represents a text-to-speech synthesis request.
type SpeechRequest struct {
	// Input is the text to synthesize (required).
	Input string `json:"input"`
	// Model is the TTS model identifier (required).
	Model string `json:"model"`
	// Voice is the provider-specific voice identifier (required).
	Voice string `json:"voice"`
	// ResponseFormat selects the audio output format ("mp3" or "pcm"). Defaults to "pcm" upstream.
	ResponseFormat string `json:"response_format,omitempty"`
	// Speed is the playback speed multiplier. Only used by providers that support it (e.g. OpenAI TTS).
	Speed *float64 `json:"speed,omitempty"`
	// Provider contains provider-specific passthrough configuration.
	Provider *SpeechProvider `json:"provider,omitempty"`
}

// SpeechProvider contains provider-specific passthrough configuration for speech requests.
type SpeechProvider struct {
	// Options is a map keyed by provider slug. The map for the matched provider
	// is spread into the upstream request body.
	Options map[string]map[string]any `json:"options,omitempty"`
}

// SpeechResponse represents the result of a text-to-speech request.
type SpeechResponse struct {
	// Audio contains the raw audio bytes (mp3 or pcm, depending on Format).
	Audio []byte
	// ContentType is the Content-Type header returned by the server.
	ContentType string
	// Format echoes the requested response format ("mp3" or "pcm").
	Format string
}
