package openrouter

import (
	"sync"
	"testing"
	"time"
)

func TestCircuitBreakerStateTransitions(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    3,
		ResetTimeout:        100 * time.Millisecond,
		HalfOpenMaxAttempts: 1,
	})

	// Initially closed
	if cb.State() != CircuitClosed {
		t.Errorf("expected CircuitClosed, got %v", cb.State())
	}

	// Allows requests when closed
	if !cb.AllowRequest() {
		t.Error("expected AllowRequest() = true when closed")
	}

	// Stays closed after failures below threshold
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitClosed {
		t.Errorf("expected CircuitClosed after 2 failures, got %v", cb.State())
	}

	// Opens after reaching threshold
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Errorf("expected CircuitOpen after 3 failures, got %v", cb.State())
	}

	// Blocks requests when open
	if cb.AllowRequest() {
		t.Error("expected AllowRequest() = false when open")
	}

	// Transitions to half-open after reset timeout
	time.Sleep(150 * time.Millisecond)
	if !cb.AllowRequest() {
		t.Error("expected AllowRequest() = true after reset timeout (half-open)")
	}
	if cb.State() != CircuitHalfOpen {
		t.Errorf("expected CircuitHalfOpen, got %v", cb.State())
	}

	// Blocks additional requests in half-open (max 1 attempt)
	if cb.AllowRequest() {
		t.Error("expected AllowRequest() = false when half-open max attempts reached")
	}
}

func TestCircuitBreakerSuccessInHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    2,
		ResetTimeout:        50 * time.Millisecond,
		HalfOpenMaxAttempts: 1,
	})

	// Open the circuit
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatal("expected open")
	}

	// Wait for half-open
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest() // triggers half-open

	// Record success closes the circuit
	cb.RecordSuccess()
	if cb.State() != CircuitClosed {
		t.Errorf("expected CircuitClosed after success in half-open, got %v", cb.State())
	}

	// Should allow requests again
	if !cb.AllowRequest() {
		t.Error("expected AllowRequest() = true after closing")
	}
}

func TestCircuitBreakerFailureInHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    2,
		ResetTimeout:        50 * time.Millisecond,
		HalfOpenMaxAttempts: 1,
	})

	// Open the circuit
	cb.RecordFailure()
	cb.RecordFailure()

	// Wait for half-open
	time.Sleep(60 * time.Millisecond)
	cb.AllowRequest()

	// Failure in half-open reopens the circuit
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Errorf("expected CircuitOpen after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    2,
		ResetTimeout:        time.Minute,
		HalfOpenMaxAttempts: 1,
	})

	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitOpen {
		t.Fatal("expected open")
	}

	cb.Reset()
	if cb.State() != CircuitClosed {
		t.Errorf("expected CircuitClosed after reset, got %v", cb.State())
	}
	if !cb.AllowRequest() {
		t.Error("expected AllowRequest() = true after reset")
	}
}

func TestCircuitBreakerDefaultConfig(t *testing.T) {
	cb := NewCircuitBreaker(nil)
	if cb.config.FailureThreshold != 5 {
		t.Errorf("expected default FailureThreshold=5, got %d", cb.config.FailureThreshold)
	}
	if cb.config.ResetTimeout != 30*time.Second {
		t.Errorf("expected default ResetTimeout=30s, got %v", cb.config.ResetTimeout)
	}
	if cb.config.HalfOpenMaxAttempts != 1 {
		t.Errorf("expected default HalfOpenMaxAttempts=1, got %d", cb.config.HalfOpenMaxAttempts)
	}
}

func TestCircuitBreakerThreadSafety(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    100,
		ResetTimeout:        time.Second,
		HalfOpenMaxAttempts: 10,
	})

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(3)
		go func() {
			defer wg.Done()
			cb.AllowRequest()
		}()
		go func() {
			defer wg.Done()
			cb.RecordFailure()
		}()
		go func() {
			defer wg.Done()
			cb.State()
		}()
	}
	wg.Wait()
}

func TestCircuitBreakerSuccessResetsClosed(t *testing.T) {
	cb := NewCircuitBreaker(&CircuitBreakerConfig{
		FailureThreshold:    5,
		ResetTimeout:        time.Second,
		HalfOpenMaxAttempts: 1,
	})

	// Record some failures (but not enough to open)
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	// Success resets failure count
	cb.RecordSuccess()

	// Now need full threshold again to open
	cb.RecordFailure()
	cb.RecordFailure()
	if cb.State() != CircuitClosed {
		t.Errorf("expected CircuitClosed (failures reset), got %v", cb.State())
	}
}

func TestCircuitStateString(t *testing.T) {
	tests := []struct {
		state CircuitState
		want  string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.state.String(); got != tt.want {
			t.Errorf("CircuitState(%d).String() = %q, want %q", tt.state, got, tt.want)
		}
	}
}

func TestCircuitBreakerError(t *testing.T) {
	err := &CircuitBreakerError{
		State:   CircuitOpen,
		Message: "circuit is open",
	}

	if err.Error() != "circuit breaker error (state: open): circuit is open" {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	// Test IsCircuitBreakerError
	cbErr, ok := IsCircuitBreakerError(err)
	if !ok {
		t.Error("expected IsCircuitBreakerError to return true")
	}
	if cbErr.State != CircuitOpen {
		t.Errorf("expected CircuitOpen, got %v", cbErr.State)
	}

	// Test with non-circuit-breaker error
	_, ok = IsCircuitBreakerError(&RequestError{})
	if ok {
		t.Error("expected IsCircuitBreakerError to return false for RequestError")
	}
}

func TestWithCircuitBreaker(t *testing.T) {
	cfg := &CircuitBreakerConfig{FailureThreshold: 10}
	client := NewClient(WithAPIKey("test"), WithCircuitBreaker(cfg))

	if client.circuitBreaker == nil {
		t.Fatal("expected circuit breaker to be set on client")
	}
	if client.circuitBreaker.config.FailureThreshold != 10 {
		t.Errorf("expected FailureThreshold=10, got %d", client.circuitBreaker.config.FailureThreshold)
	}
}

func TestClientWithoutCircuitBreaker(t *testing.T) {
	client := NewClient(WithAPIKey("test"))
	if client.circuitBreaker != nil {
		t.Error("expected nil circuit breaker by default")
	}
}
