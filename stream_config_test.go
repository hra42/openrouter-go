package openrouter

import (
	"testing"
	"time"
)

func TestDefaultStreamConfig(t *testing.T) {
	cfg := DefaultStreamConfig()

	if cfg.MaxRetries != 3 {
		t.Errorf("MaxRetries = %d, want 3", cfg.MaxRetries)
	}
	if cfg.MaxBackoff != 10*time.Second {
		t.Errorf("MaxBackoff = %v, want 10s", cfg.MaxBackoff)
	}
	if cfg.ChannelBuffer != 10 {
		t.Errorf("ChannelBuffer = %d, want 10", cfg.ChannelBuffer)
	}
	if cfg.InitialBackoff != 1*time.Second {
		t.Errorf("InitialBackoff = %v, want 1s", cfg.InitialBackoff)
	}
}

func TestResolveStreamConfig(t *testing.T) {
	reqConfig := &StreamConfig{MaxRetries: 5, MaxBackoff: 20 * time.Second, ChannelBuffer: 20, InitialBackoff: 2 * time.Second}
	clientConfig := &StreamConfig{MaxRetries: 7, MaxBackoff: 30 * time.Second, ChannelBuffer: 30, InitialBackoff: 3 * time.Second}

	t.Run("request config takes precedence", func(t *testing.T) {
		cfg := resolveStreamConfig(reqConfig, clientConfig)
		if cfg.MaxRetries != 5 {
			t.Errorf("expected request config, got MaxRetries=%d", cfg.MaxRetries)
		}
	})

	t.Run("client config used when no request config", func(t *testing.T) {
		cfg := resolveStreamConfig(nil, clientConfig)
		if cfg.MaxRetries != 7 {
			t.Errorf("expected client config, got MaxRetries=%d", cfg.MaxRetries)
		}
	})

	t.Run("default used when both nil", func(t *testing.T) {
		cfg := resolveStreamConfig(nil, nil)
		if cfg.MaxRetries != 3 {
			t.Errorf("expected default config, got MaxRetries=%d", cfg.MaxRetries)
		}
	})
}

func TestWithStreamConfig(t *testing.T) {
	cfg := &StreamConfig{MaxRetries: 10, MaxBackoff: 60 * time.Second, ChannelBuffer: 50, InitialBackoff: 5 * time.Second}
	client := NewClient(WithAPIKey("test"), WithStreamConfig(cfg))

	if client.streamConfig != cfg {
		t.Error("expected stream config to be set on client")
	}
}

func TestWithStreamConfigPerRequest(t *testing.T) {
	t.Run("chat completion options", func(t *testing.T) {
		req := &ChatCompletionRequest{}
		WithStreamMaxRetries(5)(req)
		WithStreamMaxBackoff(20 * time.Second)(req)
		WithStreamChannelBuffer(25)(req)

		if req.streamConfig == nil {
			t.Fatal("expected streamConfig to be initialized")
		}
		if req.streamConfig.MaxRetries != 5 {
			t.Errorf("MaxRetries = %d, want 5", req.streamConfig.MaxRetries)
		}
		if req.streamConfig.MaxBackoff != 20*time.Second {
			t.Errorf("MaxBackoff = %v, want 20s", req.streamConfig.MaxBackoff)
		}
		if req.streamConfig.ChannelBuffer != 25 {
			t.Errorf("ChannelBuffer = %d, want 25", req.streamConfig.ChannelBuffer)
		}
	})

	t.Run("completion options", func(t *testing.T) {
		req := &CompletionRequest{}
		WithCompletionStreamMaxRetries(5)(req)
		WithCompletionStreamMaxBackoff(20 * time.Second)(req)
		WithCompletionStreamChannelBuffer(25)(req)

		if req.streamConfig == nil {
			t.Fatal("expected streamConfig to be initialized")
		}
		if req.streamConfig.MaxRetries != 5 {
			t.Errorf("MaxRetries = %d, want 5", req.streamConfig.MaxRetries)
		}
	})
}

func TestClientDefaultStreamConfig(t *testing.T) {
	client := NewClient(WithAPIKey("test"))
	if client.streamConfig != nil {
		t.Error("expected nil streamConfig by default (uses DefaultStreamConfig at resolve time)")
	}
}
