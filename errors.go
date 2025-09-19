package openrouter

import (
	"fmt"
)

// RequestError represents an error returned by the OpenRouter API.
type RequestError struct {
	StatusCode int
	Message    string
	Type       string
	Code       string
}

// Error implements the error interface.
func (e *RequestError) Error() string {
	if e.Type != "" {
		return fmt.Sprintf("openrouter: %s (type: %s, status: %d)", e.Message, e.Type, e.StatusCode)
	}
	return fmt.Sprintf("openrouter: %s (status: %d)", e.Message, e.StatusCode)
}

// IsRateLimitError returns true if the error is a rate limit error.
func (e *RequestError) IsRateLimitError() bool {
	return e.StatusCode == 429
}

// IsAuthenticationError returns true if the error is an authentication error.
func (e *RequestError) IsAuthenticationError() bool {
	return e.StatusCode == 401
}

// IsPermissionError returns true if the error is a permission error.
func (e *RequestError) IsPermissionError() bool {
	return e.StatusCode == 403
}

// IsNotFoundError returns true if the error is a not found error.
func (e *RequestError) IsNotFoundError() bool {
	return e.StatusCode == 404
}

// IsServerError returns true if the error is a server error.
func (e *RequestError) IsServerError() bool {
	return e.StatusCode >= 500
}

// StreamError represents an error that occurs during streaming.
type StreamError struct {
	Err     error
	Message string
}

// Error implements the error interface.
func (e *StreamError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("stream error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("stream error: %s", e.Message)
}

// Unwrap returns the underlying error.
func (e *StreamError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation error for request parameters.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// ErrNoAPIKey is returned when no API key is provided.
var ErrNoAPIKey = &ValidationError{Field: "apiKey", Message: "API key is required"}

// ErrNoModel is returned when no model is specified.
var ErrNoModel = &ValidationError{Field: "model", Message: "model is required"}

// ErrNoMessages is returned when no messages are provided for chat completion.
var ErrNoMessages = &ValidationError{Field: "messages", Message: "at least one message is required"}

// ErrNoPrompt is returned when no prompt is provided for completion.
var ErrNoPrompt = &ValidationError{Field: "prompt", Message: "prompt is required"}

// IsRequestError checks if an error is a RequestError.
func IsRequestError(err error) bool {
	_, ok := err.(*RequestError)
	return ok
}

// IsStreamError checks if an error is a StreamError.
func IsStreamError(err error) bool {
	_, ok := err.(*StreamError)
	return ok
}

// IsValidationError checks if an error is a ValidationError.
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}