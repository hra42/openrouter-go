# Image Inputs Example

This example demonstrates how to send images to vision models through the OpenRouter API.

## Features Demonstrated

1. **Single Image with URL** - Send a publicly accessible image URL
2. **Multiple Images** - Send multiple images in a single request
3. **Detail Levels** - Control image analysis detail (low/high/auto)
4. **ContentBuilder** - Build complex messages with interleaved text and images
5. **Base64 Encoding** - Send local images using base64 encoding

## Running the Example

```bash
# Set your API key
export OPENROUTER_API_KEY="your-api-key"

# Run the example
go run examples/image-inputs/main.go
```

## Image Formats

OpenRouter supports the following image formats:
- PNG (image/png)
- JPEG (image/jpeg)
- WebP (image/webp)
- GIF (image/gif)

## Image Sources

### URL-based Images
Use publicly accessible image URLs. This is more efficient as it doesn't require local encoding:

```go
message := openrouter.CreateUserMessageWithImage(
    "What's in this image?",
    "https://example.com/image.jpg",
)
```

### Base64-encoded Images
For local files or private images:

```go
// Automatically encodes the image to base64
message, err := openrouter.CreateUserMessageWithBase64Image(
    "What's in this image?",
    "path/to/image.jpg",
)
```

## Detail Levels

Some models support detail level parameters:
- `"low"` - Faster and cheaper, suitable for general understanding
- `"high"` - More detailed analysis at higher cost
- `"auto"` - Let the model decide based on image size (default)

```go
message := openrouter.CreateUserMessageWithImageDetail(
    "Describe this image in detail.",
    "https://example.com/image.jpg",
    "high",
)
```

## Multiple Images

Send multiple images in a single request:

```go
message := openrouter.CreateUserMessageWithImages(
    "Compare these images",
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg",
    "https://example.com/image3.jpg",
)
```

## Content Builder

For complex messages with interleaved text and images:

```go
content := openrouter.NewContentBuilder().
    AddText("Here's the first image:").
    AddImage("https://example.com/image1.jpg").
    AddText("And here's the second:").
    AddImageWithDetail("https://example.com/image2.jpg", "high").
    AddText("What are the differences?")

message := content.BuildMessage("user")
```

## Model Support

Most modern vision models support image inputs, including:
- Google Gemini models (gemini-2.0-flash-thinking-exp, etc.)
- OpenAI GPT-4 Vision models
- Anthropic Claude 3 models
- And many others

Check the [OpenRouter documentation](https://openrouter.ai/docs) for the latest list of supported models.

## Important Notes

- The number of images you can send in a single request varies per provider and per model
- Image URLs must be publicly accessible
- Base64-encoded images increase request size, so use URLs when possible
- Some providers may have size limits on images
- Pricing may vary based on image size and detail level
