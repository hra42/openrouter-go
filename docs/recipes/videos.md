# Video Generation

Generate videos from a text prompt via OpenRouter's async `/videos` endpoint. Jobs are submitted, polled until terminal, and the resulting bytes fetched separately.

```go
job, err := client.CreateVideo(ctx, "google/veo-3.1-lite",
    "A serene mountain landscape at sunset, cinematic",
    openrouter.WithVideoAspectRatio(openrouter.VideoAspectRatio16x9),
    openrouter.WithVideoResolution(openrouter.VideoResolution720p),
    openrouter.WithVideoDuration(4),
)
if err != nil { return err }
fmt.Printf("job=%s status=%s\n", job.ID, job.Status)
```

Poll until the job reaches a terminal state (`completed`, `failed`, `cancelled`, `expired`):

```go
for {
    resp, err := client.GetVideo(ctx, job.ID)
    if err != nil { return err }
    switch resp.Status {
    case openrouter.VideoStatusCompleted:
        goto done
    case openrouter.VideoStatusFailed,
         openrouter.VideoStatusCancelled,
         openrouter.VideoStatusExpired:
        return fmt.Errorf("job %s: %s", resp.Status, resp.Error)
    }
    time.Sleep(5 * time.Second)
}
done:
```

Then download the bytes and write to disk:

```go
content, err := client.GetVideoContent(ctx, job.ID, 0) // index=0 for default output
if err != nil { return err }
_ = os.WriteFile("video.mp4", content.Content, 0o644)
```

## Options

| Option | Purpose |
|---|---|
| `WithVideoAspectRatio(ratio)` | `16:9`, `9:16`, `1:1`, `4:3`, `3:4`, `21:9`, `9:21` |
| `WithVideoResolution(res)` | `480p`, `720p`, `1080p`, `1K`, `2K`, `4K` |
| `WithVideoSize("1280x720")` | Exact pixel dimensions (alternative to resolution+ratio) |
| `WithVideoDuration(seconds)` | Target duration; supported values depend on the model |
| `WithVideoGenerateAudio(bool)` | Enable audio (model-dependent) |
| `WithVideoSeed(int)` | Deterministic sampling seed |
| `WithVideoFrameImages(...)` | First/last frame references |
| `WithVideoInputReferences(...)` | Reference images to guide generation |
| `WithVideoCallbackURL(url)` | HTTPS webhook to notify on completion |
| `WithVideoProviderOptions(provider, opts)` | Provider-specific passthrough |

Discover model capabilities with `client.ListVideoModels(ctx)` — the returned `VideoModel` exposes `SupportedAspectRatios`, `SupportedResolutions`, `SupportedDurations`, and the allowlisted passthrough parameters.

See [`examples/videos/main.go`](../../examples/videos/main.go) for a runnable submit-poll-download flow.
