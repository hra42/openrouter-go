# Test Tools

This directory contains standalone utility tools for testing multimodal functionality (images, audio, etc.).

## verify_image

Quick verification tool to test base64 image encoding without making API calls.

**Usage:**
```bash
go run cmd/openrouter-test/tools/verify_image/main.go
```

**Output:**
```
✅ Image encoded successfully!
   Data URL length: 1490970 bytes
   Prefix: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAABYAA...

The test image is ready for use in e2e tests!
```

## test_local_image

Standalone tool demonstrating local image encoding and API usage.

**Usage:**
```bash
export OPENROUTER_API_KEY="your-api-key"
go run cmd/openrouter-test/tools/test_local_image/main.go
```

## test_local_audio

Standalone tool to test audio file encoding without making API calls.

**Usage:**
```bash
go run cmd/openrouter-test/tools/test_local_audio/main.go
```

**Output:**
```
✅ Successfully encoded audio file!
   Path: cmd/openrouter-test/test-mp3.mp3
   Format: mp3
   Base64 length: 75852 characters
   ...
```

## Purpose

These tools help verify multimodal functionality works correctly before running the full e2e test suite. They're useful for:

- Testing base64 encoding without API calls
- Debugging file loading issues
- Verifying file paths are correct
- Understanding the encoded data format
