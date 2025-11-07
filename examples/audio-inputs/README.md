# Audio Inputs Example

This example demonstrates how to send audio files to OpenRouter models that support audio input.

## Features

- Base64-encoded audio from local files
- Content builder for complex messages with audio
- Manual audio encoding for fine-grained control
- Support for WAV and MP3 formats

## Usage

First, set your OpenRouter API key:

```bash
export OPENROUTER_API_KEY="your-api-key"
```

Then run the example:

```bash
go run examples/audio-inputs/main.go
```

**Note:** You'll need to update the `audioPath` variables in the code with actual paths to your audio files.

## Supported Audio Formats

OpenRouter currently supports the following audio formats:
- WAV (`.wav`)
- MP3 (`.mp3`)

## Important Notes

- Audio files must be **base64-encoded** - direct URLs are not supported for audio content
- Only models with audio processing capabilities will handle these requests
- You can search for models that support audio by filtering to audio input modality on the [OpenRouter Models page](https://openrouter.ai/models?fmt=cards&input_modalities=audio)

## Examples in this Demo

### 1. Base64-encoded Audio from Local File

Simple example showing how to transcribe an audio file:

```go
message, err := openrouter.CreateUserMessageWithBase64Audio(
    "Please transcribe this audio file.",
    "path/to/audio.wav",
)
```

### 2. Content Builder

More complex message construction with multiple content parts:

```go
content := openrouter.NewContentBuilder().
    AddText("I have an audio file to transcribe:")

content, err := content.AddBase64Audio("path/to/audio.wav")
content.AddText("Please transcribe the audio and provide a summary.")

message := content.BuildMessage("user")
```

### 3. Manual Encoding

Fine-grained control over the encoding process:

```go
base64Audio, format, err := openrouter.EncodeAudioToBase64("path/to/audio.mp3")

message := openrouter.CreateUserMessageWithAudio(
    "What is said in this audio?",
    base64Audio,
    format,
)
```

## API Reference

### Helper Functions

- `EncodeAudioToBase64(audioPath string) (string, string, error)` - Reads and encodes an audio file
- `EncodeAudioBytesToBase64(audioData []byte, format string) (string, error)` - Encodes audio bytes
- `CreateUserMessageWithAudio(text string, audioData string, format string) Message` - Creates a message with audio
- `CreateUserMessageWithBase64Audio(text string, audioPath string) (Message, error)` - Convenience function that combines encoding and message creation

### Content Builder Methods

- `AddAudio(audioData string, format string) *ContentBuilder` - Adds audio to the content builder
- `AddBase64Audio(audioPath string) (*ContentBuilder, error)` - Reads, encodes, and adds audio
