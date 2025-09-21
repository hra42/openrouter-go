package openrouter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://openrouter.ai/api/v1"
	defaultTimeout = 30 * time.Second
)

// Client is the main client for interacting with the OpenRouter API.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Optional configurations
	defaultModel      string
	referer           string
	appName           string
	maxRetries        int
	retryDelay        time.Duration
	customHeaders     map[string]string
}

// NewClient creates a new OpenRouter API client.
// The API key can be provided either as the first argument (for backwards compatibility)
// or via WithAPIKey option. If both are provided, the option takes precedence.
func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		maxRetries:    3,
		retryDelay:    time.Second,
		customHeaders: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request to the OpenRouter API.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body interface{}, v interface{}) error {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = strings.NewReader(string(jsonData))
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

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

	// Add metadata headers if present
	if reqStruct, ok := body.(interface{ GetMetadata() map[string]interface{} }); ok {
		if metadata := reqStruct.GetMetadata(); metadata != nil {
			for key, value := range metadata {
				headerKey := "X-" + key
				if strValue, ok := value.(string); ok {
					req.Header.Set(headerKey, strValue)
				}
			}
		}
	}

	// Perform request with retries
	var resp *http.Response
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-time.After(c.retryDelay * time.Duration(attempt)):
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		resp, err = c.httpClient.Do(req)
		if err != nil {
			if attempt == c.maxRetries {
				return fmt.Errorf("request failed after %d retries: %w", c.maxRetries, err)
			}
			continue
		}

		// Don't retry on client errors (4xx)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			break
		}

		// Retry on server errors (5xx) or rate limit
		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			if attempt < c.maxRetries {
				resp.Body.Close()
				continue
			}
		}

		break
	}

	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		var errorResp ErrorResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return &RequestError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		return &RequestError{
			StatusCode: resp.StatusCode,
			Message:    errorResp.Error.Message,
			Type:       errorResp.Error.Type,
			Code:       errorResp.Error.Code,
		}
	}

	// Unmarshal response
	if v != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, v); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// GetMetadata helper methods for request types
func (r *ChatCompletionRequest) GetMetadata() map[string]interface{} {
	return r.Metadata
}

func (r *CompletionRequest) GetMetadata() map[string]interface{} {
	return r.Metadata
}