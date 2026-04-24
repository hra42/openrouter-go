package openrouter

// VideoAspectRatio represents a supported aspect ratio for video generation.
type VideoAspectRatio string

// Supported video aspect ratios.
const (
	VideoAspectRatio16x9 VideoAspectRatio = "16:9"
	VideoAspectRatio9x16 VideoAspectRatio = "9:16"
	VideoAspectRatio1x1  VideoAspectRatio = "1:1"
	VideoAspectRatio4x3  VideoAspectRatio = "4:3"
	VideoAspectRatio3x4  VideoAspectRatio = "3:4"
	VideoAspectRatio21x9 VideoAspectRatio = "21:9"
	VideoAspectRatio9x21 VideoAspectRatio = "9:21"
)

// VideoResolution represents a supported output resolution for video generation.
type VideoResolution string

// Supported video resolutions.
const (
	VideoResolution480p  VideoResolution = "480p"
	VideoResolution720p  VideoResolution = "720p"
	VideoResolution1080p VideoResolution = "1080p"
	VideoResolution1K    VideoResolution = "1K"
	VideoResolution2K    VideoResolution = "2K"
	VideoResolution4K    VideoResolution = "4K"
)

// VideoFrameType identifies whether a supplied image is the first or last frame of the generated video.
type VideoFrameType string

// Supported frame image types.
const (
	VideoFrameTypeFirst VideoFrameType = "first_frame"
	VideoFrameTypeLast  VideoFrameType = "last_frame"
)

// VideoStatus represents the status of a video generation job.
type VideoStatus string

// Video generation job statuses.
const (
	VideoStatusPending    VideoStatus = "pending"
	VideoStatusInProgress VideoStatus = "in_progress"
	VideoStatusCompleted  VideoStatus = "completed"
	VideoStatusFailed     VideoStatus = "failed"
	VideoStatusCancelled  VideoStatus = "cancelled"
	VideoStatusExpired    VideoStatus = "expired"
)

// VideoImageURL wraps a URL pointing to an image used for video generation inputs.
type VideoImageURL struct {
	URL string `json:"url"`
}

// VideoContentPartImage is an image reference used to guide video generation.
type VideoContentPartImage struct {
	ImageURL VideoImageURL `json:"image_url"`
	// Type is always "image_url".
	Type string `json:"type"`
}

// VideoFrameImage is an image supplied as either the first or last frame of the generated video.
type VideoFrameImage struct {
	ImageURL VideoImageURL `json:"image_url"`
	// Type is always "image_url".
	Type string `json:"type"`
	// FrameType indicates whether this is the first or last frame.
	FrameType VideoFrameType `json:"frame_type"`
}

// VideoProvider holds provider-specific passthrough configuration for video generation requests.
type VideoProvider struct {
	// Options is a map keyed by provider slug. The map for the matched provider
	// is spread into the upstream request body.
	Options map[string]map[string]any `json:"options,omitempty"`
}

// VideoGenerationRequest represents a request to the POST /videos endpoint.
type VideoGenerationRequest struct {
	// Model is the video generation model identifier (required).
	Model string `json:"model"`
	// Prompt is the text prompt describing the video to generate (required).
	Prompt string `json:"prompt"`
	// AspectRatio selects the output aspect ratio.
	AspectRatio VideoAspectRatio `json:"aspect_ratio,omitempty"`
	// CallbackURL is an HTTPS URL that will receive a webhook notification when the job completes.
	CallbackURL string `json:"callback_url,omitempty"`
	// Duration is the requested video duration in seconds.
	Duration *int `json:"duration,omitempty"`
	// FrameImages provides first and/or last frame references for the generated video.
	FrameImages []VideoFrameImage `json:"frame_images,omitempty"`
	// GenerateAudio controls whether the provider also generates audio for the video.
	GenerateAudio *bool `json:"generate_audio,omitempty"`
	// InputReferences is a list of reference images used to guide generation.
	InputReferences []VideoContentPartImage `json:"input_references,omitempty"`
	// Provider contains provider-specific passthrough options.
	Provider *VideoProvider `json:"provider,omitempty"`
	// Resolution selects the output resolution.
	Resolution VideoResolution `json:"resolution,omitempty"`
	// Seed is an optional deterministic sampling seed.
	Seed *int `json:"seed,omitempty"`
	// Size specifies exact pixel dimensions as "WIDTHxHEIGHT" (e.g. "1280x720").
	// Interchangeable with Resolution + AspectRatio.
	Size string `json:"size,omitempty"`
}

// VideoGenerationUsage reports cost and BYOK information for a completed video generation job.
type VideoGenerationUsage struct {
	// Cost is the generation cost in USD.
	Cost *float64 `json:"cost,omitempty"`
	// IsBYOK indicates whether the request was made using a Bring Your Own Key configuration.
	IsBYOK bool `json:"is_byok,omitempty"`
}

// VideoGenerationResponse is returned by POST /videos and GET /videos/{jobId}.
type VideoGenerationResponse struct {
	// ID is the unique identifier for the video generation job.
	ID string `json:"id"`
	// PollingURL is the URL to poll for job status.
	PollingURL string `json:"polling_url"`
	// Status is the current job status.
	Status VideoStatus `json:"status"`
	// Error contains an error message if the job failed.
	Error string `json:"error,omitempty"`
	// GenerationID is the generation ID assigned to this job once processed.
	GenerationID string `json:"generation_id,omitempty"`
	// UnsignedURLs contains unsigned content URLs, when available.
	UnsignedURLs []string `json:"unsigned_urls,omitempty"`
	// Usage reports cost information once the job is complete.
	Usage *VideoGenerationUsage `json:"usage,omitempty"`
}

// VideoContentResponse is returned by GET /videos/{jobId}/content.
type VideoContentResponse struct {
	// Content holds the raw video bytes.
	Content []byte
	// ContentType is the Content-Type header returned by the server (typically "application/octet-stream").
	ContentType string
}

// VideoModel describes a single video generation model returned by GET /videos/models.
type VideoModel struct {
	ID                           string             `json:"id"`
	Name                         string             `json:"name"`
	CanonicalSlug                string             `json:"canonical_slug"`
	Created                      int64              `json:"created"`
	Description                  string             `json:"description,omitempty"`
	HuggingFaceID                *string            `json:"hugging_face_id,omitempty"`
	AllowedPassthroughParameters []string           `json:"allowed_passthrough_parameters"`
	GenerateAudio                *bool              `json:"generate_audio"`
	Seed                         *bool              `json:"seed"`
	PricingSKUs                  map[string]string  `json:"pricing_skus,omitempty"`
	SupportedAspectRatios        []VideoAspectRatio `json:"supported_aspect_ratios"`
	SupportedDurations           []int              `json:"supported_durations"`
	SupportedFrameImages         []VideoFrameType   `json:"supported_frame_images"`
	SupportedResolutions         []VideoResolution  `json:"supported_resolutions"`
	SupportedSizes               []string           `json:"supported_sizes"`
}

// VideoModelsResponse is returned by GET /videos/models.
type VideoModelsResponse struct {
	Data []VideoModel `json:"data"`
}
