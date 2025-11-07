# Audio Input Tests

This document describes the audio input tests for the OpenRouter Go client.

## Test Files

- `test-mp3.mp3` - A short MP3 audio file used for testing audio input functionality

## Available Tests

### 1. Audio Base64 Test (`-test audio`)

Tests basic audio input functionality by encoding a local MP3 file to base64 and sending it to the model.

```bash
go run cmd/openrouter-test/main.go -test audio -model google/gemini-2.5-flash
```

**What it tests:**
- Reading and encoding local audio files
- Creating messages with base64-encoded audio
- Sending audio to audio-capable models
- Handling models that don't support audio gracefully

### 2. Audio ContentBuilder Test (`-test audiobuilder`)

Tests audio input using the ContentBuilder pattern for more complex message construction.

```bash
go run cmd/openrouter-test/main.go -test audiobuilder -model google/gemini-2.5-flash
```

**What it tests:**
- Using ContentBuilder to construct messages with audio
- Mixing text and audio content parts
- The `AddBase64Audio()` method

### 3. Audio Formats Test (`-test audioformats`)

Tests support for different audio formats (currently WAV and MP3).

```bash
go run cmd/openrouter-test/main.go -test audioformats -model google/gemini-2.5-flash
```

**What it tests:**
- WAV format support (using synthetic test data)
- MP3 format support (using synthetic test data)
- Format validation and error handling

## Path Resolution

All audio tests automatically try multiple file paths to work regardless of where the test is run from:

1. `test-mp3.mp3` - When run from `cmd/openrouter-test/`
2. `cmd/openrouter-test/test-mp3.mp3` - When run from repo root
3. `../../test-mp3.mp3` - When run from nested directories

This follows the same pattern as the image and PDF tests.

## Model Compatibility

Not all models support audio input. The tests will gracefully skip if a model doesn't support audio (status codes 400, 403, or 404).

Models known to support audio input:
- `google/gemini-2.5-flash`
- `google/gemini-2.0-flash-thinking-exp`

To find more models that support audio, visit: https://openrouter.ai/models?fmt=cards&input_modalities=audio

## Supported Audio Formats

OpenRouter currently supports:
- **WAV** (`.wav`) - Waveform Audio File Format
- **MP3** (`.mp3`) - MPEG Audio Layer III

**Important:** Audio files must be base64-encoded. Direct URLs are not supported for audio content.

## Test Utility

A test utility is available to verify audio encoding without making API calls:

```bash
go run cmd/openrouter-test/tools/test_local_audio/main.go
```

This utility will:
- Find and encode the test MP3 file
- Display the base64 encoding details
- Create a test message with the audio
- Verify the message structure

## Running All Audio Tests

To run all audio-related tests in sequence:

```bash
# Run all tests (includes audio tests)
go run cmd/openrouter-test/main.go -test all -model google/gemini-2.5-flash

# Or run audio tests individually
go run cmd/openrouter-test/main.go -test audio -model google/gemini-2.5-flash
go run cmd/openrouter-test/main.go -test audiobuilder -model google/gemini-2.5-flash
go run cmd/openrouter-test/main.go -test audioformats -model google/gemini-2.5-flash
```

## Expected Behavior

When tests run successfully, you should see:

```
🔄 Test: Audio Input (Base64 from Local File)
   Using audio file: cmd/openrouter-test/test-mp3.mp3
✅ Success! (1.23s)
   Response: [Model's transcription or description of the audio]
   Model: google/gemini-2.5-flash
   Tokens: 150 prompt, 50 completion, 200 total
```

If a model doesn't support audio:

```
🔄 Test: Audio Input (Base64 from Local File)
⚠️  Skipped: Model openai/gpt-3.5-turbo not available or doesn't support audio
```

## Troubleshooting

**Error: "Failed to encode audio"**
- Ensure `test-mp3.mp3` exists in `cmd/openrouter-test/`
- Check file permissions
- Verify the file is a valid MP3

**Error: "unsupported audio format"**
- Only WAV and MP3 are supported
- Check file extension matches actual format

**Test always skipped**
- The model you're testing doesn't support audio input
- Try using `google/gemini-2.5-flash` which is known to support audio
