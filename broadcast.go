package openrouter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ParseBroadcastPayload unmarshals raw JSON bytes into an OTLPExportTraceRequest.
func ParseBroadcastPayload(data []byte) (*OTLPExportTraceRequest, error) {
	var payload OTLPExportTraceRequest
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse broadcast payload: %w", err)
	}
	return &payload, nil
}

// ExtractBroadcastTraces converts an OTLP payload into user-friendly BroadcastTrace values.
// Missing attributes produce zero values; extraction is best-effort.
func ExtractBroadcastTraces(payload *OTLPExportTraceRequest) []BroadcastTrace {
	var traces []BroadcastTrace

	for _, rs := range payload.ResourceSpans {
		resAttrs := extractAttributeMap(rs.Resource.Attributes)

		for _, ss := range rs.ScopeSpans {
			for _, span := range ss.Spans {
				t := BroadcastTrace{
					TraceID:            span.TraceID,
					SpanID:             span.SpanID,
					ParentSpanID:       span.ParentSpanID,
					SpanName:           span.Name,
					Metadata:           make(map[string]string),
					ResourceAttributes: resAttrs,
					RawAttributes:      make(map[string]string),
				}

				// Parse timestamps
				if nanos, err := strconv.ParseInt(span.StartTimeUnixNano, 10, 64); err == nil {
					t.StartTime = time.Unix(0, nanos)
				}
				if nanos, err := strconv.ParseInt(span.EndTimeUnixNano, 10, 64); err == nil {
					t.EndTime = time.Unix(0, nanos)
				}
				if !t.StartTime.IsZero() && !t.EndTime.IsZero() {
					t.Duration = t.EndTime.Sub(t.StartTime)
				}

				// Map known attributes
				for _, attr := range span.Attributes {
					val := attr.Value.StringVal()
					switch attr.Key {
					case "gen_ai.request.model":
						t.Model = val
					case "gen_ai.usage.prompt_tokens":
						t.PromptTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.completion_tokens":
						t.CompletionTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.total_tokens":
						t.TotalTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.cost":
						t.Cost, _ = strconv.ParseFloat(val, 64)
					case "user.id":
						t.UserID = val
					case "session.id":
						t.SessionID = val
					default:
						if strings.HasPrefix(attr.Key, "trace.metadata.") {
							t.Metadata[strings.TrimPrefix(attr.Key, "trace.metadata.")] = val
						} else {
							t.RawAttributes[attr.Key] = val
						}
					}
				}

				// Compute TotalTokens if absent
				if t.TotalTokens == 0 && (t.PromptTokens > 0 || t.CompletionTokens > 0) {
					t.TotalTokens = t.PromptTokens + t.CompletionTokens
				}

				traces = append(traces, t)
			}
		}
	}

	return traces
}

// ParseBroadcastTraces is a convenience function that parses raw JSON and
// extracts traces in one call.
func ParseBroadcastTraces(data []byte) ([]BroadcastTrace, error) {
	payload, err := ParseBroadcastPayload(data)
	if err != nil {
		return nil, err
	}
	return ExtractBroadcastTraces(payload), nil
}

// BroadcastWebhookHandler returns an http.HandlerFunc that receives OTLP trace
// payloads from OpenRouter Broadcast webhooks. It handles test-connection pings
// automatically and invokes the callback with parsed traces.
func BroadcastWebhookHandler(callback func([]BroadcastTrace)) http.HandlerFunc {
	return BroadcastWebhookHandlerWithError(func(traces []BroadcastTrace) error {
		callback(traces)
		return nil
	})
}

// BroadcastWebhookHandlerWithError is like BroadcastWebhookHandler but the
// callback may return an error, which results in an HTTP 500 response.
func BroadcastWebhookHandlerWithError(callback func([]BroadcastTrace) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		// Handle test-connection ping
		if r.Header.Get("X-Test-Connection") == "true" {
			w.WriteHeader(http.StatusOK)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer func() { _ = r.Body.Close() }()

		traces, err := ParseBroadcastTraces(body)
		if err != nil {
			http.Error(w, "Invalid payload: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := callback(traces); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// extractAttributeMap converts OTLP attributes into a simple string map.
func extractAttributeMap(attrs []OTLPAttribute) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Value.StringVal()
	}
	return m
}
