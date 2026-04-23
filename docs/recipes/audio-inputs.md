# Audio inputs

Audio is sent as base64-encoded WAV or MP3 in a `ContentPart`. Helpers in `audio_utils.go` make loading from disk easy.

See [`examples/audio-inputs/main.go`](../../examples/audio-inputs/main.go) for the canonical flow, including the `AudioBuilder` helper and multi-format tests.

The model must support audio input.
