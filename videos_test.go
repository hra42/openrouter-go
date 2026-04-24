package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateVideo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/videos" {
			t.Errorf("expected path /videos, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("unexpected Authorization header: %q", got)
		}

		var reqBody VideoGenerationRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if reqBody.Model != "google/veo-3.1" {
			t.Errorf("unexpected model %q", reqBody.Model)
		}
		if reqBody.Prompt != "A cat on a beach" {
			t.Errorf("unexpected prompt %q", reqBody.Prompt)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(VideoGenerationResponse{
			ID:         "job-abc123",
			PollingURL: "https://openrouter.ai/api/v1/videos/job-abc123",
			Status:     VideoStatusPending,
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateVideo(context.Background(), "google/veo-3.1", "A cat on a beach")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "job-abc123" {
		t.Errorf("unexpected id %q", resp.ID)
	}
	if resp.Status != VideoStatusPending {
		t.Errorf("expected status pending, got %q", resp.Status)
	}
	if resp.PollingURL == "" {
		t.Errorf("expected polling url, got empty")
	}
}

func TestCreateVideoWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody VideoGenerationRequest
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}

		if reqBody.AspectRatio != VideoAspectRatio16x9 {
			t.Errorf("unexpected aspect_ratio %q", reqBody.AspectRatio)
		}
		if reqBody.Duration == nil || *reqBody.Duration != 8 {
			t.Errorf("unexpected duration %+v", reqBody.Duration)
		}
		if reqBody.Resolution != VideoResolution720p {
			t.Errorf("unexpected resolution %q", reqBody.Resolution)
		}
		if reqBody.Size != "1280x720" {
			t.Errorf("unexpected size %q", reqBody.Size)
		}
		if reqBody.Seed == nil || *reqBody.Seed != 42 {
			t.Errorf("unexpected seed %+v", reqBody.Seed)
		}
		if reqBody.GenerateAudio == nil || !*reqBody.GenerateAudio {
			t.Errorf("expected generate_audio true, got %+v", reqBody.GenerateAudio)
		}
		if reqBody.CallbackURL != "https://example.com/webhook" {
			t.Errorf("unexpected callback_url %q", reqBody.CallbackURL)
		}
		if len(reqBody.FrameImages) != 1 || reqBody.FrameImages[0].FrameType != VideoFrameTypeFirst {
			t.Errorf("frame_images not propagated: %+v", reqBody.FrameImages)
		}
		if len(reqBody.InputReferences) != 1 || reqBody.InputReferences[0].Type != "image_url" {
			t.Errorf("input_references not propagated: %+v", reqBody.InputReferences)
		}
		if reqBody.Provider == nil || reqBody.Provider.Options["google-vertex"]["foo"] != "bar" {
			t.Errorf("provider options not propagated: %+v", reqBody.Provider)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(VideoGenerationResponse{
			ID:         "job-xyz",
			PollingURL: "https://openrouter.ai/api/v1/videos/job-xyz",
			Status:     VideoStatusPending,
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	_, err := client.CreateVideo(context.Background(), "google/veo-3.1", "A serene landscape",
		WithVideoAspectRatio(VideoAspectRatio16x9),
		WithVideoDuration(8),
		WithVideoResolution(VideoResolution720p),
		WithVideoSize("1280x720"),
		WithVideoSeed(42),
		WithVideoGenerateAudio(true),
		WithVideoCallbackURL("https://example.com/webhook"),
		WithVideoFrameImages(VideoFrameImage{
			ImageURL:  VideoImageURL{URL: "https://example.com/first.png"},
			Type:      "image_url",
			FrameType: VideoFrameTypeFirst,
		}),
		WithVideoInputReferences(VideoContentPartImage{
			ImageURL: VideoImageURL{URL: "https://example.com/ref.png"},
			Type:     "image_url",
		}),
		WithVideoProviderOptions("google-vertex", map[string]any{"foo": "bar"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateVideoValidation(t *testing.T) {
	t.Run("no api key", func(t *testing.T) {
		client := NewClient()
		_, err := client.CreateVideo(context.Background(), "m", "p")
		if !errors.Is(err, ErrNoAPIKey) {
			t.Errorf("expected ErrNoAPIKey, got %v", err)
		}
	})

	t.Run("empty model", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.CreateVideo(context.Background(), "", "p")
		if !errors.Is(err, ErrNoModel) {
			t.Errorf("expected ErrNoModel, got %v", err)
		}
	})

	t.Run("empty prompt", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.CreateVideo(context.Background(), "m", "")
		var verr *ValidationError
		if !errors.As(err, &verr) || verr.Field != "prompt" {
			t.Errorf("expected ValidationError for prompt, got %v", err)
		}
	})
}

func TestCreateVideoServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "invalid aspect ratio",
				Type:    "invalid_request_error",
			},
		})
	}))
	defer server.Close()

	client := NewClient(
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithRetry(0, 0),
	)

	_, err := client.CreateVideo(context.Background(), "m", "p")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var reqErr *RequestError
	if !errors.As(err, &reqErr) {
		t.Fatalf("expected *RequestError, got %T: %v", err, err)
	}
	if reqErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", reqErr.StatusCode)
	}
	if reqErr.Message != "invalid aspect ratio" {
		t.Errorf("unexpected message %q", reqErr.Message)
	}
}

func TestGetVideo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/videos/job-abc" {
			t.Errorf("expected /videos/job-abc, got %s", r.URL.Path)
		}

		cost := 0.12
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(VideoGenerationResponse{
			ID:           "job-abc",
			PollingURL:   "https://openrouter.ai/api/v1/videos/job-abc",
			Status:       VideoStatusCompleted,
			GenerationID: "gen-1",
			UnsignedURLs: []string{"https://cdn.example.com/video.mp4"},
			Usage: &VideoGenerationUsage{
				Cost:   &cost,
				IsBYOK: false,
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.GetVideo(context.Background(), "job-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != VideoStatusCompleted {
		t.Errorf("expected completed, got %q", resp.Status)
	}
	if resp.GenerationID != "gen-1" {
		t.Errorf("unexpected generation_id %q", resp.GenerationID)
	}
	if len(resp.UnsignedURLs) != 1 {
		t.Errorf("expected 1 unsigned url, got %d", len(resp.UnsignedURLs))
	}
	if resp.Usage == nil || resp.Usage.Cost == nil || *resp.Usage.Cost != 0.12 {
		t.Errorf("unexpected usage %+v", resp.Usage)
	}
}

func TestGetVideoValidation(t *testing.T) {
	t.Run("no api key", func(t *testing.T) {
		client := NewClient()
		_, err := client.GetVideo(context.Background(), "job")
		if !errors.Is(err, ErrNoAPIKey) {
			t.Errorf("expected ErrNoAPIKey, got %v", err)
		}
	})

	t.Run("empty job id", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.GetVideo(context.Background(), "")
		var verr *ValidationError
		if !errors.As(err, &verr) || verr.Field != "jobID" {
			t.Errorf("expected ValidationError for jobID, got %v", err)
		}
	})
}

func TestGetVideoContent(t *testing.T) {
	wantBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/videos/job-abc/content" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if r.URL.RawQuery != "" {
			t.Errorf("expected no query for index 0, got %q", r.URL.RawQuery)
		}

		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(wantBytes)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.GetVideoContent(context.Background(), "job-abc", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(resp.Content, wantBytes) {
		t.Errorf("content mismatch: got %v want %v", resp.Content, wantBytes)
	}
	if resp.ContentType != "video/mp4" {
		t.Errorf("unexpected content type %q", resp.ContentType)
	}
}

func TestGetVideoContentWithIndex(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/videos/job-abc/content" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("index"); got != "2" {
			t.Errorf("expected index=2, got %q", got)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{0xAA})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.GetVideoContent(context.Background(), "job-abc", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Content) != 1 {
		t.Errorf("expected 1 byte, got %d", len(resp.Content))
	}
}

func TestGetVideoContentValidation(t *testing.T) {
	t.Run("no api key", func(t *testing.T) {
		client := NewClient()
		_, err := client.GetVideoContent(context.Background(), "job", 0)
		if !errors.Is(err, ErrNoAPIKey) {
			t.Errorf("expected ErrNoAPIKey, got %v", err)
		}
	})

	t.Run("empty job id", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.GetVideoContent(context.Background(), "", 0)
		var verr *ValidationError
		if !errors.As(err, &verr) || verr.Field != "jobID" {
			t.Errorf("expected ValidationError for jobID, got %v", err)
		}
	})
}

func TestListVideoModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/videos/models" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}

		supportsAudio := true
		supportsSeed := true

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(VideoModelsResponse{
			Data: []VideoModel{
				{
					ID:                           "google/veo-3.1",
					Name:                         "Veo 3.1",
					CanonicalSlug:                "google/veo-3.1",
					Created:                      1700000000,
					Description:                  "Google Veo 3.1",
					AllowedPassthroughParameters: []string{"enhance_prompt"},
					GenerateAudio:                &supportsAudio,
					Seed:                         &supportsSeed,
					SupportedAspectRatios:        []VideoAspectRatio{VideoAspectRatio16x9, VideoAspectRatio9x16},
					SupportedDurations:           []int{4, 8},
					SupportedFrameImages:         []VideoFrameType{VideoFrameTypeFirst, VideoFrameTypeLast},
					SupportedResolutions:         []VideoResolution{VideoResolution720p, VideoResolution1080p},
					SupportedSizes:               []string{"1280x720", "1920x1080"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.ListVideoModels(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 model, got %d", len(resp.Data))
	}
	model := resp.Data[0]
	if model.ID != "google/veo-3.1" {
		t.Errorf("unexpected id %q", model.ID)
	}
	if len(model.SupportedAspectRatios) != 2 {
		t.Errorf("expected 2 aspect ratios, got %d", len(model.SupportedAspectRatios))
	}
	if len(model.SupportedDurations) != 2 || model.SupportedDurations[1] != 8 {
		t.Errorf("unexpected durations %+v", model.SupportedDurations)
	}
}

func TestListVideoModelsNoAPIKey(t *testing.T) {
	client := NewClient()
	_, err := client.ListVideoModels(context.Background())
	if !errors.Is(err, ErrNoAPIKey) {
		t.Errorf("expected ErrNoAPIKey, got %v", err)
	}
}
