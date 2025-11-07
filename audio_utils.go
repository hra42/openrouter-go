package openrouter

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// EncodeAudioToBase64 reads an audio file and encodes it to base64.
// It automatically detects the audio format based on the file extension.
// Supported formats: wav, mp3
// Note: Audio must be base64-encoded; direct URLs are not supported.
func EncodeAudioToBase64(audioPath string) (string, string, error) {
	// Read the audio file
	audioData, err := os.ReadFile(audioPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read audio file: %w", err)
	}

	// Detect format from file extension
	ext := strings.ToLower(filepath.Ext(audioPath))
	var format string
	switch ext {
	case ".wav":
		format = "wav"
	case ".mp3":
		format = "mp3"
	default:
		return "", "", fmt.Errorf("unsupported audio format: %s (supported: wav, mp3)", ext)
	}

	// Encode to base64
	base64Audio := base64.StdEncoding.EncodeToString(audioData)

	return base64Audio, format, nil
}

// EncodeAudioBytesToBase64 encodes audio bytes to base64.
// The format should be one of: "wav", "mp3"
func EncodeAudioBytesToBase64(audioData []byte, format string) (string, error) {
	// Validate format
	switch format {
	case "wav", "mp3":
		// Valid format
	default:
		return "", fmt.Errorf("unsupported audio format: %s (supported: wav, mp3)", format)
	}

	base64Audio := base64.StdEncoding.EncodeToString(audioData)
	return base64Audio, nil
}

// CreateUserMessageWithAudio creates a user message with text and an audio input.
// The audioData should be base64-encoded audio data.
// The format should be one of: "wav", "mp3"
func CreateUserMessageWithAudio(text string, audioData string, format string) Message {
	return Message{
		Role: "user",
		Content: []ContentPart{
			{Type: "text", Text: text},
			{Type: "input_audio", InputAudio: &InputAudio{
				Data:   audioData,
				Format: format,
			}},
		},
	}
}

// CreateUserMessageWithBase64Audio creates a user message with a base64-encoded audio file.
// This is a convenience function that combines EncodeAudioToBase64 and CreateUserMessageWithAudio.
func CreateUserMessageWithBase64Audio(text string, audioPath string) (Message, error) {
	base64Audio, format, err := EncodeAudioToBase64(audioPath)
	if err != nil {
		return Message{}, err
	}
	return CreateUserMessageWithAudio(text, base64Audio, format), nil
}

// AddAudio adds an audio input to the content builder.
// The audioData should be base64-encoded audio data.
// The format should be one of: "wav", "mp3"
func (cb *ContentBuilder) AddAudio(audioData string, format string) *ContentBuilder {
	cb.parts = append(cb.parts, ContentPart{
		Type: "input_audio",
		InputAudio: &InputAudio{
			Data:   audioData,
			Format: format,
		},
	})
	return cb
}

// AddBase64Audio reads and encodes an audio file, then adds it to the content builder.
func (cb *ContentBuilder) AddBase64Audio(audioPath string) (*ContentBuilder, error) {
	base64Audio, format, err := EncodeAudioToBase64(audioPath)
	if err != nil {
		return cb, err
	}
	return cb.AddAudio(base64Audio, format), nil
}
