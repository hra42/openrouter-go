# Test Tools

This directory contains standalone utility tools for testing image functionality.

## verify_image.go

Quick verification tool to test base64 image encoding without making API calls.

**Usage:**
```bash
cd cmd/openrouter-test
go run tools/verify_image.go
```

**Output:**
```
✅ Image encoded successfully!
   Data URL length: 1490970 bytes
   Prefix: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAABYAA...

The test image is ready for use in e2e tests!
```

## test_local_image.go

Standalone function demonstrating local image encoding and API usage.

**Note:** This file contains a `TestLocalImage()` function that can be imported and called from other test files. It's not meant to be run directly as it doesn't have a main function.

## Purpose

These tools help verify the image functionality works correctly before running the full e2e test suite. They're useful for:

- Testing base64 encoding without API calls
- Debugging image loading issues
- Verifying file paths are correct
- Understanding the encoded data format
