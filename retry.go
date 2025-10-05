package openrouter

import (
	"context"
	"math"
	"math/rand"
	"time"
)

const (
	// defaultJitterFactor is the default jitter factor for retry backoff (±25%).
	defaultJitterFactor = 0.25
	// maxReconnectBackoff is the maximum backoff duration for stream reconnection attempts.
	maxReconnectBackoff = 10 * time.Second
	// defaultMaxDelay is the default maximum delay for retry backoff.
	defaultMaxDelay = 30 * time.Second
	// defaultMultiplier is the default multiplier for exponential backoff.
	defaultMultiplier = 2.0
)

// RetryConfig configures retry behavior for API requests.
type RetryConfig struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	Multiplier     float64
	Jitter         bool
	RetryableError func(error) bool
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:   3,
		InitialDelay: 1 * time.Second,
		MaxDelay:     defaultMaxDelay,
		Multiplier:   defaultMultiplier,
		Jitter:       true,
		RetryableError: func(err error) bool {
			if reqErr, ok := err.(*RequestError); ok {
				// Retry on rate limit, server errors, and gateway timeouts
				return reqErr.IsRateLimitError() ||
					reqErr.IsServerError() ||
					reqErr.StatusCode == 502 ||
					reqErr.StatusCode == 503 ||
					reqErr.StatusCode == 504
			}
			return false
		},
	}
}

// calculateBackoff calculates the backoff duration for a retry attempt.
func (rc *RetryConfig) calculateBackoff(attempt int) time.Duration {
	if attempt <= 0 {
		return rc.InitialDelay
	}

	// Exponential backoff
	delay := float64(rc.InitialDelay) * math.Pow(rc.Multiplier, float64(attempt-1))

	// Cap at max delay
	if delay > float64(rc.MaxDelay) {
		delay = float64(rc.MaxDelay)
	}

	// Add jitter if enabled
	if rc.Jitter {
		// Add ±25% jitter
		jitter := delay * defaultJitterFactor
		delay = delay + (rand.Float64()*2-1)*jitter
	}

	return time.Duration(delay)
}

// RetryWithBackoff executes a function with exponential backoff retry logic.
func RetryWithBackoff(ctx context.Context, config *RetryConfig, fn func() error) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// Execute the function
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if we should retry
		if attempt >= config.MaxRetries {
			break
		}

		// Check if the error is retryable
		if config.RetryableError != nil && !config.RetryableError(err) {
			return err
		}

		// Calculate backoff duration
		backoff := config.calculateBackoff(attempt + 1)

		// Handle special case for rate limit errors
		if reqErr, ok := err.(*RequestError); ok && reqErr.IsRateLimitError() {
			// If we have a Retry-After header value, use it
			// For now, we'll use our calculated backoff
			// In a real implementation, we'd parse the Retry-After header
			backoff = config.calculateBackoff(attempt + 1)
		}

		// Wait with context cancellation support
		select {
		case <-time.After(backoff):
			// Continue to next retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if lastErr != nil {
		return lastErr
	}

	return nil
}

// RateLimiter provides rate limiting functionality.
type RateLimiter struct {
	requestsPerSecond float64
	burst             int
	tokens            chan struct{}
	ticker            *time.Ticker
	done              chan struct{}
}

// NewRateLimiter creates a new rate limiter.
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 10 // Default to 10 requests per second
	}
	if burst <= 0 {
		burst = 1
	}

	rl := &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		burst:             burst,
		tokens:            make(chan struct{}, burst),
		ticker:            time.NewTicker(time.Duration(float64(time.Second) / requestsPerSecond)),
		done:              make(chan struct{}),
	}

	// Fill the token bucket initially
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	// Start the token replenishment goroutine
	go rl.refill()

	return rl
}

// refill continuously adds tokens to the bucket.
func (rl *RateLimiter) refill() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case rl.tokens <- struct{}{}:
				// Token added
			default:
				// Bucket is full
			}
		case <-rl.done:
			rl.ticker.Stop()
			return
		}
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close stops the rate limiter.
func (rl *RateLimiter) Close() {
	close(rl.done)
}

// doRequest performs an HTTP request to the OpenRouter API with retry logic.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, v interface{}) error {
	config := &RetryConfig{
		MaxRetries:   c.maxRetries,
		InitialDelay: c.retryDelay,
		MaxDelay:     defaultMaxDelay,
		Multiplier:   defaultMultiplier,
		Jitter:       true,
		RetryableError: func(err error) bool {
			if reqErr, ok := err.(*RequestError); ok {
				// Don't retry client errors except rate limit
				if reqErr.StatusCode >= 400 && reqErr.StatusCode < 500 {
					return reqErr.IsRateLimitError()
				}
				// Retry server errors
				return true
			}
			// Retry network errors
			return true
		},
	}

	return RetryWithBackoff(ctx, config, func() error {
		return c.doRequestOnce(ctx, method, endpoint, body, v)
	})
}
