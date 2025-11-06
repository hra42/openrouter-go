# Image Input Tests

This document describes the image input testing capabilities integrated into the OpenRouter Go client test suite.

## Overview

The test suite now includes comprehensive tests for multimodal (image) inputs, covering both URL-based and base64-encoded local images.

## Test Image

**File:** `test-image.png` (1.1 MB)
**URL:** `https://hra42.com/test-image.png`

**Description:** A laboratory scene featuring five test tubes with vibrant colored liquids (red, orange, yellow, green, and blue) arranged in a test tube rack. The background shows blurred laboratory equipment, creating a professional scientific setting.

**Purpose:** This image is used for both URL-based and base64-encoded image tests, allowing direct comparison between the two approaches. The same image content is accessed via:
- Direct URL for testing URL-based image inputs
- Local file for testing base64 encoding functionality

## Available Image Tests

### 1. Image Input (URL)
**Command:** `go run . -test image -model google/gemini-2.0-flash-thinking-exp:free`

Tests sending a single image via public URL to a vision model.

### 2. Multiple Images
**Command:** `go run . -test multiimage -model google/gemini-2.0-flash-thinking-exp:free`

Tests sending multiple images in a single request for comparison or analysis.

### 3. Image with Detail Level
**Command:** `go run . -test imagedetail -model google/gemini-2.0-flash-thinking-exp:free`

Tests the detail parameter (low/high/auto) for controlling image analysis quality and cost.

### 4. Content Builder
**Command:** `go run . -test contentbuilder -model google/gemini-2.0-flash-thinking-exp:free`

Tests the flexible ContentBuilder API for constructing complex multimodal messages with interleaved text and images.

### 5. Base64 Local Image (NEW!)
**Command:** `go run . -test base64image -model google/gemini-2.0-flash-thinking-exp:free`

Tests encoding and sending the local `test-image.png` file as a base64 data URL. This test:
- Reads the local PNG file
- Automatically encodes it to base64 with proper data URL format
- Sends it to a vision model
- Verifies the model can analyze the image content

**Important:** This test must be run from the `cmd/openrouter-test` directory.

## Quick Verification

To verify the test image is properly encoded without making an API call:

```bash
cd cmd/openrouter-test
go run tools/verify_image.go
```

This will encode the test image and display the base64 data URL length and prefix.

To verify the hosted URL is accessible:

```bash
curl -I https://hra42.com/test-image.png
# Should return: HTTP/2 200
```

## Helper Functions Used

The tests utilize these helper functions from the library:

- `CreateUserMessageWithImage()` - Single image with text
- `CreateUserMessageWithImages()` - Multiple images with text
- `CreateUserMessageWithImageDetail()` - Image with detail level control
- `CreateUserMessageWithBase64Image()` - Automatic base64 encoding from file
- `EncodeImageToBase64()` - Manual base64 encoding with format detection
- `NewContentBuilder()` - Flexible message construction

## Expected Results

When running the base64 image test with a vision model, you should see output similar to:

```
🔄 Test: Base64-Encoded Local Image
✅ Success! (3.45s)
   Response: The image shows five test tubes containing colored liquids - red, orange, yellow, green, and blue - arranged in a laboratory test tube rack...
   Model: google/gemini-2.0-flash-thinking-exp:free
   Tokens: 1234 prompt, 56 completion, 1290 total
   Image: test-image.png (base64-encoded)
```

## Notes

- Vision tests require vision-capable models (Gemini, GPT-4 Vision, Claude 3, etc.)
- Base64 encoding increases payload size by ~33%
- For production use, prefer URLs when images are publicly accessible
- The test image (1.1 MB) encodes to approximately 1.5 MB base64
- Some models may have size limits on base64-encoded images
