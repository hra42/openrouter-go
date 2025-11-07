package openrouter

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncodeAudioToBase64(t *testing.T) {
	tests := []struct {
		name        string
		extension   string
		content     []byte
		wantFormat  string
		wantErr     bool
		errContains string
	}{
		{
			name:       "Valid WAV file",
			extension:  ".wav",
			content:    []byte("fake wav audio data"),
			wantFormat: "wav",
			wantErr:    false,
		},
		{
			name:       "Valid MP3 file",
			extension:  ".mp3",
			content:    []byte("fake mp3 audio data"),
			wantFormat: "mp3",
			wantErr:    false,
		},
		{
			name:        "Unsupported format",
			extension:   ".ogg",
			content:     []byte("fake ogg audio data"),
			wantErr:     true,
			errContains: "unsupported audio format",
		},
		{
			name:       "Uppercase extension",
			extension:  ".WAV",
			content:    []byte("fake wav audio data"),
			wantFormat: "wav",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test"+tt.extension)
			if err := os.WriteFile(tmpFile, tt.content, 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}

			base64Audio, format, err := EncodeAudioToBase64(tmpFile)

			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeAudioToBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("EncodeAudioToBase64() error = %v, should contain %v", err, tt.errContains)
					}
				}
				return
			}

			if format != tt.wantFormat {
				t.Errorf("EncodeAudioToBase64() format = %v, want %v", format, tt.wantFormat)
			}

			// Verify base64 encoding
			decoded, err := base64.StdEncoding.DecodeString(base64Audio)
			if err != nil {
				t.Errorf("Failed to decode base64: %v", err)
			}

			if string(decoded) != string(tt.content) {
				t.Errorf("Decoded content doesn't match original")
			}
		})
	}
}

func TestEncodeAudioToBase64_NonExistentFile(t *testing.T) {
	_, _, err := EncodeAudioToBase64("/nonexistent/audio.wav")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to read audio file") {
		t.Errorf("Error should mention failed to read audio file, got: %v", err)
	}
}

func TestEncodeAudioBytesToBase64(t *testing.T) {
	tests := []struct {
		name        string
		audioData   []byte
		format      string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Valid WAV data",
			audioData: []byte("fake wav audio data"),
			format:    "wav",
			wantErr:   false,
		},
		{
			name:      "Valid MP3 data",
			audioData: []byte("fake mp3 audio data"),
			format:    "mp3",
			wantErr:   false,
		},
		{
			name:        "Invalid format",
			audioData:   []byte("fake audio data"),
			format:      "ogg",
			wantErr:     true,
			errContains: "unsupported audio format",
		},
		{
			name:      "Empty data",
			audioData: []byte{},
			format:    "wav",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base64Audio, err := EncodeAudioBytesToBase64(tt.audioData, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeAudioBytesToBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("EncodeAudioBytesToBase64() error = %v, should contain %v", err, tt.errContains)
					}
				}
				return
			}

			// Verify base64 encoding
			decoded, err := base64.StdEncoding.DecodeString(base64Audio)
			if err != nil {
				t.Errorf("Failed to decode base64: %v", err)
			}

			if string(decoded) != string(tt.audioData) {
				t.Errorf("Decoded content doesn't match original")
			}
		})
	}
}

func TestCreateUserMessageWithAudio(t *testing.T) {
	text := "Please transcribe this audio."
	audioData := "ZmFrZSBhdWRpbyBkYXRh" // base64 encoded "fake audio data"
	format := "wav"

	msg := CreateUserMessageWithAudio(text, audioData, format)

	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatalf("Expected Content to be []ContentPart, got %T", msg.Content)
	}

	if len(parts) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(parts))
	}

	// Check text part
	if parts[0].Type != "text" {
		t.Errorf("Expected first part type 'text', got '%s'", parts[0].Type)
	}
	if parts[0].Text != text {
		t.Errorf("Expected text '%s', got '%s'", text, parts[0].Text)
	}

	// Check audio part
	if parts[1].Type != "input_audio" {
		t.Errorf("Expected second part type 'input_audio', got '%s'", parts[1].Type)
	}
	if parts[1].InputAudio == nil {
		t.Fatal("Expected InputAudio to be non-nil")
	}
	if parts[1].InputAudio.Data != audioData {
		t.Errorf("Expected audio data '%s', got '%s'", audioData, parts[1].InputAudio.Data)
	}
	if parts[1].InputAudio.Format != format {
		t.Errorf("Expected audio format '%s', got '%s'", format, parts[1].InputAudio.Format)
	}
}

func TestCreateUserMessageWithBase64Audio(t *testing.T) {
	// Create a temporary audio file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.wav")
	testContent := []byte("fake wav audio data")
	if err := os.WriteFile(tmpFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	text := "Please transcribe this audio."
	msg, err := CreateUserMessageWithBase64Audio(text, tmpFile)
	if err != nil {
		t.Fatalf("CreateUserMessageWithBase64Audio() error = %v", err)
	}

	if msg.Role != "user" {
		t.Errorf("Expected role 'user', got '%s'", msg.Role)
	}

	parts, ok := msg.Content.([]ContentPart)
	if !ok {
		t.Fatalf("Expected Content to be []ContentPart, got %T", msg.Content)
	}

	if len(parts) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(parts))
	}

	// Check audio part
	if parts[1].Type != "input_audio" {
		t.Errorf("Expected second part type 'input_audio', got '%s'", parts[1].Type)
	}
	if parts[1].InputAudio == nil {
		t.Fatal("Expected InputAudio to be non-nil")
	}
	if parts[1].InputAudio.Format != "wav" {
		t.Errorf("Expected audio format 'wav', got '%s'", parts[1].InputAudio.Format)
	}

	// Verify the encoded data
	decoded, err := base64.StdEncoding.DecodeString(parts[1].InputAudio.Data)
	if err != nil {
		t.Errorf("Failed to decode base64: %v", err)
	}
	if string(decoded) != string(testContent) {
		t.Errorf("Decoded content doesn't match original")
	}
}

func TestContentBuilder_AddAudio(t *testing.T) {
	cb := NewContentBuilder()
	audioData := "ZmFrZSBhdWRpbyBkYXRh" // base64 encoded "fake audio data"
	format := "wav"

	cb.AddText("Please transcribe this audio.").AddAudio(audioData, format)

	parts := cb.Build()
	if len(parts) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(parts))
	}

	// Check audio part
	if parts[1].Type != "input_audio" {
		t.Errorf("Expected second part type 'input_audio', got '%s'", parts[1].Type)
	}
	if parts[1].InputAudio == nil {
		t.Fatal("Expected InputAudio to be non-nil")
	}
	if parts[1].InputAudio.Data != audioData {
		t.Errorf("Expected audio data '%s', got '%s'", audioData, parts[1].InputAudio.Data)
	}
	if parts[1].InputAudio.Format != format {
		t.Errorf("Expected audio format '%s', got '%s'", format, parts[1].InputAudio.Format)
	}
}

func TestContentBuilder_AddBase64Audio(t *testing.T) {
	// Create a temporary audio file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.mp3")
	testContent := []byte("fake mp3 audio data")
	if err := os.WriteFile(tmpFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cb := NewContentBuilder()
	cb.AddText("Please transcribe this audio.")

	cb, err := cb.AddBase64Audio(tmpFile)
	if err != nil {
		t.Fatalf("AddBase64Audio() error = %v", err)
	}

	parts := cb.Build()
	if len(parts) != 2 {
		t.Fatalf("Expected 2 content parts, got %d", len(parts))
	}

	// Check audio part
	if parts[1].Type != "input_audio" {
		t.Errorf("Expected second part type 'input_audio', got '%s'", parts[1].Type)
	}
	if parts[1].InputAudio == nil {
		t.Fatal("Expected InputAudio to be non-nil")
	}
	if parts[1].InputAudio.Format != "mp3" {
		t.Errorf("Expected audio format 'mp3', got '%s'", parts[1].InputAudio.Format)
	}

	// Verify the encoded data
	decoded, err := base64.StdEncoding.DecodeString(parts[1].InputAudio.Data)
	if err != nil {
		t.Errorf("Failed to decode base64: %v", err)
	}
	if string(decoded) != string(testContent) {
		t.Errorf("Decoded content doesn't match original")
	}
}

func TestContentBuilder_AddBase64Audio_NonExistentFile(t *testing.T) {
	cb := NewContentBuilder()
	_, err := cb.AddBase64Audio("/nonexistent/audio.wav")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
