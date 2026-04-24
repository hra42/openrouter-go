package openrouter

import (
	"context"
	"fmt"
)

// VideoGenerationOption is a functional option for video generation requests.
type VideoGenerationOption func(*VideoGenerationRequest)

// CreateVideo submits a video generation job and returns the initial response (including the job ID and polling URL).
// Callers should poll GetVideo with the returned ID until Status is terminal
// (VideoStatusCompleted, VideoStatusFailed, VideoStatusCancelled, or VideoStatusExpired).
func (c *Client) CreateVideo(ctx context.Context, model, prompt string, opts ...VideoGenerationOption) (*VideoGenerationResponse, error) {
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if model == "" {
		return nil, ErrNoModel
	}

	if prompt == "" {
		return nil, &ValidationError{
			Field:   "prompt",
			Message: "prompt is required",
		}
	}

	req := &VideoGenerationRequest{
		Model:  model,
		Prompt: prompt,
	}

	for _, opt := range opts {
		opt(req)
	}

	var resp VideoGenerationResponse
	if err := c.doRequest(ctx, "POST", "/videos", req, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetVideo returns the current status of a video generation job.
func (c *Client) GetVideo(ctx context.Context, jobID string) (*VideoGenerationResponse, error) {
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if jobID == "" {
		return nil, &ValidationError{
			Field:   "jobID",
			Message: "jobID is required",
		}
	}

	var resp VideoGenerationResponse
	if err := c.doRequest(ctx, "GET", "/videos/"+jobID, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetVideoContent downloads the generated video bytes for a completed job.
// Use index to select a specific output when the provider produced multiple videos;
// pass 0 to fetch the default (first) output.
func (c *Client) GetVideoContent(ctx context.Context, jobID string, index int) (*VideoContentResponse, error) {
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	if jobID == "" {
		return nil, &ValidationError{
			Field:   "jobID",
			Message: "jobID is required",
		}
	}

	path := "/videos/" + jobID + "/content"
	if index > 0 {
		path = fmt.Sprintf("%s?index=%d", path, index)
	}

	content, contentType, err := c.doRequestRaw(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	return &VideoContentResponse{
		Content:     content,
		ContentType: contentType,
	}, nil
}

// ListVideoModels returns the list of available video generation models and their capabilities.
func (c *Client) ListVideoModels(ctx context.Context) (*VideoModelsResponse, error) {
	if c.apiKey == "" {
		return nil, ErrNoAPIKey
	}

	var resp VideoModelsResponse
	if err := c.doRequest(ctx, "GET", "/videos/models", nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// WithVideoAspectRatio sets the output aspect ratio.
func WithVideoAspectRatio(ratio VideoAspectRatio) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.AspectRatio = ratio
	}
}

// WithVideoDuration sets the requested video duration in seconds.
func WithVideoDuration(seconds int) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.Duration = &seconds
	}
}

// WithVideoResolution sets the output resolution.
func WithVideoResolution(resolution VideoResolution) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.Resolution = resolution
	}
}

// WithVideoSize sets the exact pixel dimensions (e.g. "1280x720").
// Interchangeable with WithVideoResolution + WithVideoAspectRatio.
func WithVideoSize(size string) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.Size = size
	}
}

// WithVideoSeed sets a deterministic sampling seed. Determinism is not guaranteed for all providers.
func WithVideoSeed(seed int) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.Seed = &seed
	}
}

// WithVideoGenerateAudio enables or disables audio generation alongside the video.
func WithVideoGenerateAudio(enabled bool) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.GenerateAudio = &enabled
	}
}

// WithVideoCallbackURL sets the HTTPS webhook URL to receive a completion notification.
func WithVideoCallbackURL(url string) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.CallbackURL = url
	}
}

// WithVideoFrameImages sets the first- and/or last-frame reference images.
func WithVideoFrameImages(frames ...VideoFrameImage) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.FrameImages = frames
	}
}

// WithVideoInputReferences sets the reference images used to guide generation.
func WithVideoInputReferences(refs ...VideoContentPartImage) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		r.InputReferences = refs
	}
}

// WithVideoProviderOptions sets provider-specific passthrough options keyed by provider slug.
// The given options map is spread into the upstream request body when the matching provider is used.
func WithVideoProviderOptions(provider string, options map[string]any) VideoGenerationOption {
	return func(r *VideoGenerationRequest) {
		if r.Provider == nil {
			r.Provider = &VideoProvider{}
		}
		if r.Provider.Options == nil {
			r.Provider.Options = make(map[string]map[string]any)
		}
		r.Provider.Options[provider] = options
	}
}
