package openrouter

import (
	"errors"
	"strings"
	"testing"
)

func TestRequestError(t *testing.T) {
	tests := []struct {
		name           string
		err            *RequestError
		expectedString string
		isRateLimit    bool
		isAuth         bool
		isPermission   bool
		isNotFound     bool
		isServer       bool
	}{
		{
			name: "rate limit error",
			err: &RequestError{
				StatusCode: 429,
				Message:    "Rate limit exceeded",
				Type:       "rate_limit_error",
			},
			expectedString: "openrouter: Rate limit exceeded (type: rate_limit_error, status: 429)",
			isRateLimit:    true,
		},
		{
			name: "authentication error",
			err: &RequestError{
				StatusCode: 401,
				Message:    "Invalid API key",
				Type:       "authentication_error",
			},
			expectedString: "openrouter: Invalid API key (type: authentication_error, status: 401)",
			isAuth:         true,
		},
		{
			name: "permission error",
			err: &RequestError{
				StatusCode: 403,
				Message:    "Access denied",
				Type:       "permission_error",
			},
			expectedString: "openrouter: Access denied (type: permission_error, status: 403)",
			isPermission:   true,
		},
		{
			name: "not found error",
			err: &RequestError{
				StatusCode: 404,
				Message:    "Model not found",
				Type:       "not_found_error",
			},
			expectedString: "openrouter: Model not found (type: not_found_error, status: 404)",
			isNotFound:     true,
		},
		{
			name: "server error",
			err: &RequestError{
				StatusCode: 500,
				Message:    "Internal server error",
				Type:       "server_error",
			},
			expectedString: "openrouter: Internal server error (type: server_error, status: 500)",
			isServer:       true,
		},
		{
			name: "error without type",
			err: &RequestError{
				StatusCode: 400,
				Message:    "Bad request",
			},
			expectedString: "openrouter: Bad request (status: 400)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Error() string
			if got := tt.err.Error(); got != tt.expectedString {
				t.Errorf("Error() = %q, want %q", got, tt.expectedString)
			}

			// Test type checking methods
			if got := tt.err.IsRateLimitError(); got != tt.isRateLimit {
				t.Errorf("IsRateLimitError() = %v, want %v", got, tt.isRateLimit)
			}
			if got := tt.err.IsAuthenticationError(); got != tt.isAuth {
				t.Errorf("IsAuthenticationError() = %v, want %v", got, tt.isAuth)
			}
			if got := tt.err.IsPermissionError(); got != tt.isPermission {
				t.Errorf("IsPermissionError() = %v, want %v", got, tt.isPermission)
			}
			if got := tt.err.IsNotFoundError(); got != tt.isNotFound {
				t.Errorf("IsNotFoundError() = %v, want %v", got, tt.isNotFound)
			}
			if got := tt.err.IsServerError(); got != tt.isServer {
				t.Errorf("IsServerError() = %v, want %v", got, tt.isServer)
			}
		})
	}
}

func TestStreamError(t *testing.T) {
	baseErr := errors.New("connection failed")

	tests := []struct {
		name           string
		err            *StreamError
		expectedString string
		hasUnwrap      bool
	}{
		{
			name: "with underlying error",
			err: &StreamError{
				Err:     baseErr,
				Message: "failed to read stream",
			},
			expectedString: "stream error: failed to read stream: connection failed",
			hasUnwrap:      true,
		},
		{
			name: "without underlying error",
			err: &StreamError{
				Message: "stream terminated unexpectedly",
			},
			expectedString: "stream error: stream terminated unexpectedly",
			hasUnwrap:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Error() string
			if got := tt.err.Error(); got != tt.expectedString {
				t.Errorf("Error() = %q, want %q", got, tt.expectedString)
			}

			// Test Unwrap()
			unwrapped := tt.err.Unwrap()
			if tt.hasUnwrap && unwrapped != baseErr {
				t.Errorf("Unwrap() returned wrong error")
			}
			if !tt.hasUnwrap && unwrapped != nil {
				t.Errorf("Unwrap() should return nil")
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name           string
		err            *ValidationError
		expectedString string
	}{
		{
			name: "with field",
			err: &ValidationError{
				Field:   "temperature",
				Message: "must be between 0 and 1",
			},
			expectedString: "validation error for field 'temperature': must be between 0 and 1",
		},
		{
			name: "without field",
			err: &ValidationError{
				Message: "invalid request format",
			},
			expectedString: "validation error: invalid request format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expectedString {
				t.Errorf("Error() = %q, want %q", got, tt.expectedString)
			}
		})
	}
}

func TestPredefinedErrors(t *testing.T) {
	// Test ErrNoAPIKey
	if !strings.Contains(ErrNoAPIKey.Error(), "API key is required") {
		t.Errorf("ErrNoAPIKey has unexpected message: %v", ErrNoAPIKey)
	}

	// Test ErrNoModel
	if !strings.Contains(ErrNoModel.Error(), "model is required") {
		t.Errorf("ErrNoModel has unexpected message: %v", ErrNoModel)
	}

	// Test ErrNoMessages
	if !strings.Contains(ErrNoMessages.Error(), "at least one message is required") {
		t.Errorf("ErrNoMessages has unexpected message: %v", ErrNoMessages)
	}

	// Test ErrNoPrompt
	if !strings.Contains(ErrNoPrompt.Error(), "prompt is required") {
		t.Errorf("ErrNoPrompt has unexpected message: %v", ErrNoPrompt)
	}
}

func TestErrorTypeChecking(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		isRequestError bool
		isStreamError  bool
		isValidation   bool
	}{
		{
			name:           "RequestError",
			err:            &RequestError{StatusCode: 400, Message: "bad request"},
			isRequestError: true,
			isStreamError:  false,
			isValidation:   false,
		},
		{
			name:           "StreamError",
			err:            &StreamError{Message: "stream failed"},
			isRequestError: false,
			isStreamError:  true,
			isValidation:   false,
		},
		{
			name:           "ValidationError",
			err:            &ValidationError{Field: "test", Message: "invalid"},
			isRequestError: false,
			isStreamError:  false,
			isValidation:   true,
		},
		{
			name:           "generic error",
			err:            errors.New("generic error"),
			isRequestError: false,
			isStreamError:  false,
			isValidation:   false,
		},
		{
			name:           "nil error",
			err:            nil,
			isRequestError: false,
			isStreamError:  false,
			isValidation:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRequestError(tt.err); got != tt.isRequestError {
				t.Errorf("IsRequestError() = %v, want %v", got, tt.isRequestError)
			}
			if got := IsStreamError(tt.err); got != tt.isStreamError {
				t.Errorf("IsStreamError() = %v, want %v", got, tt.isStreamError)
			}
			if got := IsValidationError(tt.err); got != tt.isValidation {
				t.Errorf("IsValidationError() = %v, want %v", got, tt.isValidation)
			}
		})
	}
}

func TestRequestErrorStatusChecks(t *testing.T) {
	tests := []struct {
		statusCode int
		checkFunc  func(*RequestError) bool
		expected   bool
	}{
		{429, (*RequestError).IsRateLimitError, true},
		{401, (*RequestError).IsAuthenticationError, true},
		{403, (*RequestError).IsPermissionError, true},
		{404, (*RequestError).IsNotFoundError, true},
		{500, (*RequestError).IsServerError, true},
		{502, (*RequestError).IsServerError, true},
		{503, (*RequestError).IsServerError, true},
		{400, (*RequestError).IsServerError, false},
		{200, (*RequestError).IsServerError, false},
	}

	for _, tt := range tests {
		err := &RequestError{StatusCode: tt.statusCode}
		if got := tt.checkFunc(err); got != tt.expected {
			t.Errorf("StatusCode %d: check returned %v, want %v", tt.statusCode, got, tt.expected)
		}
	}
}
