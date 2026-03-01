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
					SpanMetadata:       make(map[string]string),
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
					// Model fields
					case "gen_ai.response.model":
						t.ResponseModel = val
						t.Model = val // backward compat
					case "gen_ai.request.model":
						t.RequestModel = val
						if t.Model == "" {
							t.Model = val // backward compat fallback
						}

					// Token usage (new canonical keys)
					case "gen_ai.usage.input_tokens":
						t.InputTokens, _ = strconv.Atoi(val)
						t.PromptTokens = t.InputTokens // backward compat
					case "gen_ai.usage.output_tokens":
						t.OutputTokens, _ = strconv.Atoi(val)
						t.CompletionTokens = t.OutputTokens // backward compat

					// Token usage (old keys, backward compat)
					case "gen_ai.usage.prompt_tokens":
						n, _ := strconv.Atoi(val)
						t.PromptTokens = n
						if t.InputTokens == 0 {
							t.InputTokens = n
						}
					case "gen_ai.usage.completion_tokens":
						n, _ := strconv.Atoi(val)
						t.CompletionTokens = n
						if t.OutputTokens == 0 {
							t.OutputTokens = n
						}

					case "gen_ai.usage.total_tokens":
						t.TotalTokens, _ = strconv.Atoi(val)

					// Cost fields (new canonical key)
					case "gen_ai.usage.total_cost":
						t.TotalCost, _ = strconv.ParseFloat(val, 64)
						t.Cost = t.TotalCost // backward compat
					// Cost (old key, backward compat)
					case "gen_ai.usage.cost":
						f, _ := strconv.ParseFloat(val, 64)
						t.Cost = f
						if t.TotalCost == 0 {
							t.TotalCost = f
						}
					case "gen_ai.usage.input_cost":
						t.InputCost, _ = strconv.ParseFloat(val, 64)
					case "gen_ai.usage.output_cost":
						t.OutputCost, _ = strconv.ParseFloat(val, 64)

					// Token detail fields
					case "gen_ai.usage.input_tokens.cached":
						t.CachedTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.input_tokens.audio":
						t.AudioInputTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.input_tokens.video":
						t.VideoInputTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.output_tokens.image":
						t.ImageOutputTokens, _ = strconv.Atoi(val)
					case "gen_ai.usage.output_tokens.reasoning":
						t.ReasoningTokens, _ = strconv.Atoi(val)

					// GenAI semantic convention fields
					case "gen_ai.operation.name":
						t.OperationName = val
					case "gen_ai.system":
						t.System = val
					case "gen_ai.provider.name":
						t.ProviderName = val
					case "gen_ai.response.finish_reason":
						t.FinishReason = val
					case "gen_ai.response.finish_reasons":
						t.FinishReasons = val

					// OpenRouter-specific fields
					case "openrouter.provider_slug":
						t.ProviderSlug = val
					case "openrouter.provider_name":
						t.OpenRouterProviderName = val
					case "openrouter.api_key_name":
						t.APIKeyName = val
					case "openrouter.entity_id":
						t.EntityID = val
					case "openrouter.user_id":
						t.OpenRouterUserID = val
					case "openrouter.finish_reason":
						t.OpenRouterFinishReason = val
					case "openrouter.input_unit_price":
						t.InputUnitPrice, _ = strconv.ParseFloat(val, 64)
					case "openrouter.output_unit_price":
						t.OutputUnitPrice, _ = strconv.ParseFloat(val, 64)
					case "openrouter.source":
						t.Source = val

					// Content fields
					case "gen_ai.prompt":
						t.Prompt = val
					case "gen_ai.completion":
						t.Completion = val

					// Span-level fields
					case "span.type":
						t.SpanType = val
					case "span.level":
						t.SpanLevel = val
					case "span.input":
						t.SpanInput = val
					case "span.output":
						t.SpanOutput = val

					// Trace-level fields
					case "trace.name":
						t.TraceName = val
					case "trace.input":
						t.TraceInput = val
					case "trace.output":
						t.TraceOutput = val
					case "trace.tags":
						t.TraceTags = val

					// Identity fields
					case "user.id":
						t.UserID = val
					case "session.id":
						t.SessionID = val

					default:
						if after, ok := strings.CutPrefix(attr.Key, "trace.metadata."); ok {
							t.Metadata[after] = val
						} else if after, ok := strings.CutPrefix(attr.Key, "span.metadata."); ok {
							t.SpanMetadata[after] = val
						} else {
							t.RawAttributes[attr.Key] = val
						}
					}
				}

				// Compute TotalTokens if absent
				if t.TotalTokens == 0 && (t.InputTokens > 0 || t.OutputTokens > 0) {
					t.TotalTokens = t.InputTokens + t.OutputTokens
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
