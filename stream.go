package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/hra42/openrouter-go/internal/sse"
)

// eventStream handles Server-Sent Events (SSE) streaming.
type eventStream struct {
	ctx       context.Context
	cancel    context.CancelFunc
	response  *http.Response
	scanner   *sse.Scanner
	events    chan StreamEvent
	err       error
	errMu     sync.RWMutex
	closed    bool
	closeMu   sync.Mutex
	reconnect bool
	client    *Client
	endpoint  string
	body      interface{}
}

// createStream creates a new SSE stream for the given endpoint and request.
func (c *Client) createStream(ctx context.Context, endpoint string, body interface{}) (*eventStream, error) {
	url := c.baseURL + endpoint

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	if c.referer != "" {
		req.Header.Set("HTTP-Referer", c.referer)
	}

	if c.appName != "" {
		req.Header.Set("X-Title", c.appName)
	}

	// Add custom headers
	for key, value := range c.customHeaders {
		req.Header.Set(key, value)
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var errorResp ErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, &RequestError{
				StatusCode: resp.StatusCode,
				Message:    string(body),
			}
		}

		return nil, &RequestError{
			StatusCode: resp.StatusCode,
			Message:    errorResp.Error.Message,
			Type:       errorResp.Error.Type,
			Code:       errorResp.Error.Code,
		}
	}

	// Create stream context
	streamCtx, cancel := context.WithCancel(ctx)

	stream := &eventStream{
		ctx:       streamCtx,
		cancel:    cancel,
		response:  resp,
		scanner:   sse.NewScanner(resp.Body),
		events:    make(chan StreamEvent, 10),
		reconnect: true,
		client:    c,
		endpoint:  endpoint,
		body:      body,
	}

	// Start reading events
	go stream.readEvents()

	return stream, nil
}

// readEvents reads SSE events from the stream.
func (es *eventStream) readEvents() {
	defer close(es.events)
	defer es.response.Body.Close()

	retryCount := 0
	maxRetries := 3

	for {
		// Check if context is cancelled
		select {
		case <-es.ctx.Done():
			es.setError(es.ctx.Err())
			return
		default:
		}

		// Read next event
		if !es.scanner.Scan() {
			if err := es.scanner.Err(); err != nil {
				// Handle connection errors with reconnection
				if es.reconnect && retryCount < maxRetries {
					retryCount++
					if es.attemptReconnect(retryCount) {
						continue
					}
				}
				es.setError(&StreamError{
					Err:     err,
					Message: "stream ended unexpectedly",
				})
			}
			return
		}

		event := es.scanner.Event()

		// Check for end of stream
		if sse.IsEndOfStream(event.Data) {
			return
		}

		// Convert SSE event to StreamEvent
		streamEvent := StreamEvent{
			ID:    event.ID,
			Event: event.Type,
			Data:  string(event.Data),
			Retry: event.Retry,
		}

		// Send the event
		select {
		case es.events <- streamEvent:
			retryCount = 0 // Reset retry count on successful event
		case <-es.ctx.Done():
			return
		}
	}
}

// Events returns the channel of SSE events.
func (es *eventStream) Events() <-chan StreamEvent {
	return es.events
}

// Err returns any error that occurred during streaming.
func (es *eventStream) Err() error {
	es.errMu.RLock()
	defer es.errMu.RUnlock()
	return es.err
}

// setError sets the stream error.
func (es *eventStream) setError(err error) {
	es.errMu.Lock()
	defer es.errMu.Unlock()
	if es.err == nil {
		es.err = err
	}
}

// Close closes the stream.
func (es *eventStream) Close() error {
	es.closeMu.Lock()
	defer es.closeMu.Unlock()

	if es.closed {
		return nil
	}

	es.closed = true
	es.reconnect = false // Disable reconnection on explicit close
	es.cancel()

	if es.response != nil && es.response.Body != nil {
		return es.response.Body.Close()
	}

	return nil
}

// attemptReconnect attempts to reconnect to the stream.
func (es *eventStream) attemptReconnect(attempt int) bool {
	// Close current connection
	if es.response != nil && es.response.Body != nil {
		es.response.Body.Close()
	}

	// Calculate backoff
	backoff := time.Duration(attempt) * time.Second
	if backoff > 10*time.Second {
		backoff = 10 * time.Second
	}

	// Wait before reconnecting
	select {
	case <-time.After(backoff):
	case <-es.ctx.Done():
		return false
	}

	// Attempt to reconnect
	url := es.client.baseURL + es.endpoint

	jsonData, err := json.Marshal(es.body)
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(es.ctx, "POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return false
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+es.client.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	if es.client.referer != "" {
		req.Header.Set("HTTP-Referer", es.client.referer)
	}

	if es.client.appName != "" {
		req.Header.Set("X-Title", es.client.appName)
	}

	// Add custom headers
	for key, value := range es.client.customHeaders {
		req.Header.Set(key, value)
	}

	// Make the request
	resp, err := es.client.httpClient.Do(req)
	if err != nil {
		return false
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return false
	}

	// Update stream with new connection
	es.response = resp
	es.scanner = sse.NewScanner(resp.Body)

	return true
}

// parseSSEData parses the SSE data field into the given value.
func parseSSEData(data string, v interface{}) error {
	if data == "" {
		return nil
	}

	if err := json.Unmarshal([]byte(data), v); err != nil {
		return &StreamError{
			Err:     err,
			Message: "failed to parse SSE data",
		}
	}

	return nil
}

// Helper function to concatenate streaming chat responses.
func ConcatenateChatStreamResponses(responses []ChatCompletionResponse) string {
	var result strings.Builder

	for _, resp := range responses {
		for _, choice := range resp.Choices {
			if choice.Delta != nil && choice.Delta.Content != nil {
				if content, ok := choice.Delta.Content.(string); ok {
					result.WriteString(content)
				}
			}
		}
	}

	return result.String()
}

// Helper function to concatenate streaming completion responses.
func ConcatenateCompletionStreamResponses(responses []CompletionResponse) string {
	var result strings.Builder

	for _, resp := range responses {
		for _, choice := range resp.Choices {
			result.WriteString(choice.Text)
		}
	}

	return result.String()
}