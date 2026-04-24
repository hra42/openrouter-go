package openrouter

import (
	"context"
)

// SpeechOption is a functional option for speech synthesis requests.
type SpeechOption func(*SpeechRequest)

// CreateSpeech synthesizes audio from the given input text using the specified model and voice.
// It returns the raw audio bytes along with the server's Content-Type and the resolved format.
func (c *Client) CreateSpeech(ctx context.Context, input, model, voice string, opts ...SpeechOption) (*SpeechResponse, error) {
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if input == "" {
		return nil, &ValidationError{
			Field:   "input",
			Message: "input is required",
		}
	}

	if model == "" {
		return nil, ErrNoModel
	}

	if voice == "" {
		return nil, &ValidationError{
			Field:   "voice",
			Message: "voice is required",
		}
	}

	req := &SpeechRequest{
		Input: input,
		Model: model,
		Voice: voice,
	}

	for _, opt := range opts {
		opt(req)
	}

	audio, contentType, err := c.doRequestRaw(ctx, "POST", "/audio/speech", req)
	if err != nil {
		return nil, err
	}

	format := req.ResponseFormat
	if format == "" {
		format = SpeechFormatPCM
	}

	return &SpeechResponse{
		Audio:       audio,
		ContentType: contentType,
		Format:      format,
	}, nil
}

// WithSpeechResponseFormat sets the audio output format ("mp3" or "pcm").
func WithSpeechResponseFormat(format string) SpeechOption {
	return func(r *SpeechRequest) {
		r.ResponseFormat = format
	}
}

// WithSpeechSpeed sets the playback speed multiplier.
// Only used by providers that support it (e.g. OpenAI TTS); ignored otherwise.
func WithSpeechSpeed(speed float64) SpeechOption {
	return func(r *SpeechRequest) {
		r.Speed = &speed
	}
}

// WithSpeechProviderOptions sets provider-specific passthrough options keyed by provider slug.
// The given options map is spread into the upstream request body when the matching provider is used.
func WithSpeechProviderOptions(provider string, options map[string]any) SpeechOption {
	return func(r *SpeechRequest) {
		if r.Provider == nil {
			r.Provider = &SpeechProvider{}
		}
		if r.Provider.Options == nil {
			r.Provider.Options = make(map[string]map[string]any)
		}
		r.Provider.Options[provider] = options
	}
}
