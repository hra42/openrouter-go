# Text-to-Speech (Create Speech)

Synthesize audio from text using OpenRouter's `/audio/speech` endpoint. The response is raw audio bytes (`mp3` or `pcm`), which you can write straight to disk.

```go
resp, err := client.CreateSpeech(ctx,
    "Hello from the OpenRouter Go SDK.",
    "hexgrad/kokoro-82m",
    "af_bella",
)
if err != nil { return err }

if err := os.WriteFile("speech.pcm", resp.Audio, 0o644); err != nil {
    return err
}
fmt.Printf("%d bytes, %s (%s)\n", len(resp.Audio), resp.ContentType, resp.Format)
```

The default format is `pcm`. Request MP3 and tweak playback speed via options:

```go
resp, err := client.CreateSpeech(ctx, input, "hexgrad/kokoro-82m", "af_bella",
    openrouter.WithSpeechResponseFormat(openrouter.SpeechFormatMP3),
    openrouter.WithSpeechSpeed(1.25),
)
```

Provider-specific passthrough parameters can be supplied with `WithSpeechProviderOptions("openai", map[string]any{...})` — the map is spread into the upstream request body when the matching provider serves the request.

See [`examples/tts/main.go`](../../examples/tts/main.go) for a runnable example.
