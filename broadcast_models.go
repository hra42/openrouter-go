package openrouter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// OTLPExportTraceRequest is the top-level OTLP JSON trace payload
// sent by OpenRouter's Broadcast webhook feature.
type OTLPExportTraceRequest struct {
	ResourceSpans []OTLPResourceSpan `json:"resourceSpans"`
}

// OTLPResourceSpan groups spans by their originating resource.
type OTLPResourceSpan struct {
	Resource   OTLPResource    `json:"resource"`
	ScopeSpans []OTLPScopeSpan `json:"scopeSpans"`
}

// OTLPResource describes the entity producing telemetry.
type OTLPResource struct {
	Attributes []OTLPAttribute `json:"attributes"`
}

// OTLPScopeSpan groups spans by instrumentation scope.
type OTLPScopeSpan struct {
	Scope *OTLPScope `json:"scope,omitempty"`
	Spans []OTLPSpan `json:"spans"`
}

// OTLPScope identifies the instrumentation library.
type OTLPScope struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

// OTLPSpan represents a single span in a trace.
type OTLPSpan struct {
	TraceID           string          `json:"traceId"`
	SpanID            string          `json:"spanId"`
	ParentSpanID      string          `json:"parentSpanId,omitempty"`
	Name              string          `json:"name"`
	Kind              int             `json:"kind"`
	StartTimeUnixNano string          `json:"startTimeUnixNano"`
	EndTimeUnixNano   string          `json:"endTimeUnixNano"`
	Attributes        []OTLPAttribute `json:"attributes"`
	Status            *OTLPStatus     `json:"status,omitempty"`
	Events            []OTLPEvent     `json:"events,omitempty"`
}

// OTLPAttribute is a key-value pair attached to a span or resource.
type OTLPAttribute struct {
	Key   string       `json:"key"`
	Value OTLPAnyValue `json:"value"`
}

// OTLPAnyValue represents a polymorphic OTLP value.
// The OTLP spec encodes int64 values as strings, but some implementations
// send them as JSON numbers. IntValue accepts both forms.
type OTLPAnyValue struct {
	StringValue *string         `json:"stringValue,omitempty"`
	IntValue    *FlexInt        `json:"intValue,omitempty"`
	DoubleValue *float64        `json:"doubleValue,omitempty"`
	BoolValue   *bool           `json:"boolValue,omitempty"`
	ArrayValue  *OTLPArrayValue `json:"arrayValue,omitempty"`
}

// FlexInt is a string-backed integer that can unmarshal from both
// a JSON string ("150") and a JSON number (150).
type FlexInt string

// UnmarshalJSON implements json.Unmarshaler for FlexInt.
func (f *FlexInt) UnmarshalJSON(data []byte) error {
	// Try string first (OTLP canonical form).
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexInt(s)
		return nil
	}
	// Fall back to number.
	var n int64
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexInt(strconv.FormatInt(n, 10))
		return nil
	}
	return fmt.Errorf("intValue: cannot unmarshal %s as string or number", string(data))
}

// String returns the string representation.
func (f FlexInt) String() string {
	return string(f)
}

// MarshalJSON implements json.Marshaler, encoding as a JSON string per the OTLP spec.
func (f FlexInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

// StringVal returns the value as a string regardless of its underlying type.
func (v OTLPAnyValue) StringVal() string {
	switch {
	case v.StringValue != nil:
		return *v.StringValue
	case v.IntValue != nil:
		return v.IntValue.String()
	case v.DoubleValue != nil:
		return fmt.Sprintf("%g", *v.DoubleValue)
	case v.BoolValue != nil:
		return fmt.Sprintf("%t", *v.BoolValue)
	default:
		return ""
	}
}

// OTLPArrayValue wraps a slice of OTLP values.
type OTLPArrayValue struct {
	Values []OTLPAnyValue `json:"values"`
}

// OTLPStatus describes the status of a span.
type OTLPStatus struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// OTLPEvent represents a timed event within a span.
type OTLPEvent struct {
	Name         string          `json:"name"`
	TimeUnixNano string          `json:"timeUnixNano,omitempty"`
	Attributes   []OTLPAttribute `json:"attributes,omitempty"`
}

// BroadcastTrace is a user-friendly representation of a single span
// extracted from an OTLP trace payload sent by OpenRouter Broadcast.
type BroadcastTrace struct {
	TraceID      string
	SpanID       string
	ParentSpanID string
	SpanName     string

	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	Model            string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64

	UserID    string
	SessionID string

	// Metadata contains values from trace.metadata.* attributes (prefix stripped).
	Metadata map[string]string
	// ResourceAttributes contains attributes from the OTLP resource.
	ResourceAttributes map[string]string
	// RawAttributes contains all other span attributes not mapped to named fields.
	RawAttributes map[string]string
}
