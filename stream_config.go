package openrouter

import "time"

// StreamConfig configures stream reconnection and buffering behavior.
type StreamConfig struct {
	// MaxRetries is the maximum number of reconnection attempts. Default: 3.
	MaxRetries int
	// MaxBackoff is the maximum backoff duration between reconnections. Default: 10s.
	MaxBackoff time.Duration
	// ChannelBuffer is the buffer size for the events channel. Default: 10.
	ChannelBuffer int
	// InitialBackoff is the initial backoff duration for reconnection. Default: 1s.
	InitialBackoff time.Duration
}

// DefaultStreamConfig returns the default stream configuration.
func DefaultStreamConfig() *StreamConfig {
	return &StreamConfig{
		MaxRetries:     3,
		MaxBackoff:     10 * time.Second,
		ChannelBuffer:  10,
		InitialBackoff: 1 * time.Second,
	}
}

// resolveStreamConfig returns the effective stream config by checking
// request-level, client-level, and default configs in order.
func resolveStreamConfig(reqConfig, clientConfig *StreamConfig) *StreamConfig {
	if reqConfig != nil {
		return reqConfig
	}
	if clientConfig != nil {
		return clientConfig
	}
	return DefaultStreamConfig()
}
