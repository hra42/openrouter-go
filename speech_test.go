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

func TestCreateSpeech(t *testing.T) {
	wantAudio := []byte{0x00, 0x01, 0x02, 0x03, 0x04}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/audio/speech" {
			t.Errorf("expected path /audio/speech, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("unexpected Authorization header: %q", got)
		}

		var reqBody SpeechRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if reqBody.Input != "Hello world" {
			t.Errorf("expected input 'Hello world', got %q", reqBody.Input)
		}
		if reqBody.Model != "openai/gpt-4o-mini-tts" {
			t.Errorf("unexpected model %q", reqBody.Model)
		}
		if reqBody.Voice != "alloy" {
			t.Errorf("unexpected voice %q", reqBody.Voice)
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(wantAudio)
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateSpeech(context.Background(), "Hello world", "openai/gpt-4o-mini-tts", "alloy")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(resp.Audio, wantAudio) {
		t.Errorf("audio bytes mismatch: got %v want %v", resp.Audio, wantAudio)
	}
	if resp.ContentType != "application/octet-stream" {
		t.Errorf("unexpected content type %q", resp.ContentType)
	}
	if resp.Format != SpeechFormatPCM {
		t.Errorf("expected default format %q, got %q", SpeechFormatPCM, resp.Format)
	}
}

func TestCreateSpeechWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody SpeechRequest
		body, _ := io.ReadAll(r.Body)
		if err := json.Unmarshal(body, &reqBody); err != nil {
			t.Fatalf("failed to decode body: %v", err)
		}
		if reqBody.ResponseFormat != SpeechFormatMP3 {
			t.Errorf("expected mp3 format, got %q", reqBody.ResponseFormat)
		}
		if reqBody.Speed == nil || *reqBody.Speed != 1.25 {
			t.Errorf("unexpected speed: %+v", reqBody.Speed)
		}
		if reqBody.Provider == nil || reqBody.Provider.Options["openai"]["instructions"] != "cheerful" {
			t.Errorf("provider options not propagated: %+v", reqBody.Provider)
		}

		w.Header().Set("Content-Type", "audio/mpeg")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte{0xFF, 0xFB})
	}))
	defer server.Close()

	client := NewClient(WithAPIKey("test-key"), WithBaseURL(server.URL))

	resp, err := client.CreateSpeech(context.Background(), "hi", "openai/gpt-4o-mini-tts", "alloy",
		WithSpeechResponseFormat(SpeechFormatMP3),
		WithSpeechSpeed(1.25),
		WithSpeechProviderOptions("openai", map[string]any{"instructions": "cheerful"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Format != SpeechFormatMP3 {
		t.Errorf("expected format mp3, got %q", resp.Format)
	}
	if resp.ContentType != "audio/mpeg" {
		t.Errorf("unexpected content type %q", resp.ContentType)
	}
	if len(resp.Audio) != 2 {
		t.Errorf("expected 2 bytes, got %d", len(resp.Audio))
	}
}

func TestCreateSpeechValidation(t *testing.T) {
	t.Run("no api key", func(t *testing.T) {
		client := NewClient()
		_, err := client.CreateSpeech(context.Background(), "hi", "m", "v")
		if !errors.Is(err, ErrNoAPIKey) {
			t.Errorf("expected ErrNoAPIKey, got %v", err)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.CreateSpeech(context.Background(), "", "m", "v")
		var verr *ValidationError
		if !errors.As(err, &verr) || verr.Field != "input" {
			t.Errorf("expected ValidationError for input, got %v", err)
		}
	})

	t.Run("empty model", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.CreateSpeech(context.Background(), "hi", "", "v")
		if !errors.Is(err, ErrNoModel) {
			t.Errorf("expected ErrNoModel, got %v", err)
		}
	})

	t.Run("empty voice", func(t *testing.T) {
		client := NewClient(WithAPIKey("k"))
		_, err := client.CreateSpeech(context.Background(), "hi", "m", "")
		var verr *ValidationError
		if !errors.As(err, &verr) || verr.Field != "voice" {
			t.Errorf("expected ValidationError for voice, got %v", err)
		}
	})
}

func TestCreateSpeechServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(ErrorResponse{
			Error: APIError{
				Message: "invalid voice",
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

	_, err := client.CreateSpeech(context.Background(), "hi", "m", "v")
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
	if reqErr.Message != "invalid voice" {
		t.Errorf("unexpected message %q", reqErr.Message)
	}
}
