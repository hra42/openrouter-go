package openrouter

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestParseBroadcastPayload(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		spans   int
	}{
		{
			name:    "valid empty object",
			input:   `{}`,
			wantErr: false,
			spans:   0,
		},
		{
			name:    "empty resourceSpans",
			input:   `{"resourceSpans":[]}`,
			wantErr: false,
			spans:   0,
		},
		{
			name:    "invalid JSON",
			input:   `{not json}`,
			wantErr: true,
		},
		{
			name: "single span",
			input: `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
				"traceId":"abc123","spanId":"def456","name":"chat","kind":1,
				"startTimeUnixNano":"1700000000000000000","endTimeUnixNano":"1700000001000000000",
				"attributes":[]
			}]}]}]}`,
			wantErr: false,
			spans:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := ParseBroadcastPayload([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			traces := ExtractBroadcastTraces(payload)
			if len(traces) != tt.spans {
				t.Errorf("got %d traces, want %d", len(traces), tt.spans)
			}
		})
	}
}

func TestExtractBroadcastTraces_FullAttributes(t *testing.T) {
	payload := `{"resourceSpans":[{"resource":{"attributes":[
		{"key":"service.name","value":{"stringValue":"openrouter"}}
	]},"scopeSpans":[{"scope":{"name":"broadcast","version":"1.0"},"spans":[{
		"traceId":"aaaa","spanId":"bbbb","parentSpanId":"cccc","name":"chat completions",
		"kind":2,
		"startTimeUnixNano":"1700000000000000000",
		"endTimeUnixNano":"1700000002500000000",
		"attributes":[
			{"key":"gen_ai.request.model","value":{"stringValue":"openai/gpt-4o"}},
			{"key":"gen_ai.usage.prompt_tokens","value":{"intValue":"150"}},
			{"key":"gen_ai.usage.completion_tokens","value":{"intValue":"50"}},
			{"key":"gen_ai.usage.cost","value":{"doubleValue":0.0035}},
			{"key":"user.id","value":{"stringValue":"user-123"}},
			{"key":"session.id","value":{"stringValue":"sess-456"}},
			{"key":"trace.metadata.env","value":{"stringValue":"production"}},
			{"key":"trace.metadata.version","value":{"stringValue":"2.1"}},
			{"key":"custom.tag","value":{"stringValue":"hello"}}
		],
		"status":{"code":1}
	}]}]}]}` //nolint:lll

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(traces) != 1 {
		t.Fatalf("got %d traces, want 1", len(traces))
	}

	tr := traces[0]

	// Basic span fields
	assertEqual(t, "TraceID", tr.TraceID, "aaaa")
	assertEqual(t, "SpanID", tr.SpanID, "bbbb")
	assertEqual(t, "ParentSpanID", tr.ParentSpanID, "cccc")
	assertEqual(t, "SpanName", tr.SpanName, "chat completions")

	// Timing
	wantStart := time.Unix(0, 1700000000000000000)
	wantEnd := time.Unix(0, 1700000002500000000)
	if !tr.StartTime.Equal(wantStart) {
		t.Errorf("StartTime = %v, want %v", tr.StartTime, wantStart)
	}
	if !tr.EndTime.Equal(wantEnd) {
		t.Errorf("EndTime = %v, want %v", tr.EndTime, wantEnd)
	}
	if tr.Duration != 2500*time.Millisecond {
		t.Errorf("Duration = %v, want 2.5s", tr.Duration)
	}

	// Model & tokens
	assertEqual(t, "Model", tr.Model, "openai/gpt-4o")
	assertEqualInt(t, "PromptTokens", tr.PromptTokens, 150)
	assertEqualInt(t, "CompletionTokens", tr.CompletionTokens, 50)
	assertEqualInt(t, "TotalTokens", tr.TotalTokens, 200) // computed

	// Cost
	if tr.Cost != 0.0035 {
		t.Errorf("Cost = %v, want 0.0035", tr.Cost)
	}

	// User & session
	assertEqual(t, "UserID", tr.UserID, "user-123")
	assertEqual(t, "SessionID", tr.SessionID, "sess-456")

	// Metadata
	if tr.Metadata["env"] != "production" {
		t.Errorf("Metadata[env] = %q, want %q", tr.Metadata["env"], "production")
	}
	if tr.Metadata["version"] != "2.1" {
		t.Errorf("Metadata[version] = %q, want %q", tr.Metadata["version"], "2.1")
	}

	// Resource attributes
	if tr.ResourceAttributes["service.name"] != "openrouter" {
		t.Errorf("ResourceAttributes[service.name] = %q, want %q", tr.ResourceAttributes["service.name"], "openrouter")
	}

	// Raw attributes
	if tr.RawAttributes["custom.tag"] != "hello" {
		t.Errorf("RawAttributes[custom.tag] = %q, want %q", tr.RawAttributes["custom.tag"], "hello")
	}
}

func TestExtractBroadcastTraces_TotalTokensExplicit(t *testing.T) {
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t1","spanId":"s1","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.prompt_tokens","value":{"intValue":"10"}},
			{"key":"gen_ai.usage.completion_tokens","value":{"intValue":"20"}},
			{"key":"gen_ai.usage.total_tokens","value":{"intValue":"30"}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqualInt(t, "TotalTokens", traces[0].TotalTokens, 30)
}

func TestExtractBroadcastTraces_MultipleSpans(t *testing.T) {
	// Two resourceSpans, each with one scopeSpan containing one span
	payload := `{"resourceSpans":[
		{"resource":{"attributes":[{"key":"host","value":{"stringValue":"a"}}]},
		 "scopeSpans":[{"spans":[
			{"traceId":"t1","spanId":"s1","name":"span1","kind":1,"startTimeUnixNano":"0","endTimeUnixNano":"0","attributes":[]},
			{"traceId":"t1","spanId":"s2","name":"span2","kind":1,"startTimeUnixNano":"0","endTimeUnixNano":"0","attributes":[]}
		 ]}]},
		{"resource":{"attributes":[{"key":"host","value":{"stringValue":"b"}}]},
		 "scopeSpans":[{"spans":[
			{"traceId":"t2","spanId":"s3","name":"span3","kind":1,"startTimeUnixNano":"0","endTimeUnixNano":"0","attributes":[]}
		 ]}]}
	]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(traces) != 3 {
		t.Fatalf("got %d traces, want 3", len(traces))
	}

	// Check resource attributes are scoped correctly
	if traces[0].ResourceAttributes["host"] != "a" {
		t.Errorf("traces[0] host = %q, want %q", traces[0].ResourceAttributes["host"], "a")
	}
	if traces[2].ResourceAttributes["host"] != "b" {
		t.Errorf("traces[2] host = %q, want %q", traces[2].ResourceAttributes["host"], "b")
	}
}

func TestOTLPAnyValue_StringVal(t *testing.T) {
	tests := []struct {
		name string
		val  OTLPAnyValue
		want string
	}{
		{"string", OTLPAnyValue{StringValue: new("hello")}, "hello"},
		{"int", OTLPAnyValue{IntValue: flexIntPtr("42")}, "42"},
		{"double", OTLPAnyValue{DoubleValue: new(3.14)}, "3.14"},
		{"bool true", OTLPAnyValue{BoolValue: new(true)}, "true"},
		{"bool false", OTLPAnyValue{BoolValue: new(false)}, "false"},
		{"empty", OTLPAnyValue{}, ""},
		{"array", OTLPAnyValue{ArrayValue: &OTLPArrayValue{}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.val.StringVal()
			if got != tt.want {
				t.Errorf("StringVal() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBroadcastWebhookHandler_ValidPayload(t *testing.T) {
	var received []BroadcastTrace
	handler := BroadcastWebhookHandler(func(traces []BroadcastTrace) {
		received = traces
	})

	body := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t1","spanId":"s1","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[{"key":"gen_ai.request.model","value":{"stringValue":"test-model"}}]
	}]}]}]}`

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if len(received) != 1 {
		t.Fatalf("received %d traces, want 1", len(received))
	}
	if received[0].Model != "test-model" {
		t.Errorf("Model = %q, want %q", received[0].Model, "test-model")
	}
}

func TestBroadcastWebhookHandler_TestConnection(t *testing.T) {
	called := false
	handler := BroadcastWebhookHandler(func(_ []BroadcastTrace) {
		called = true
	})

	req := httptest.NewRequest(http.MethodPost, "/webhook", nil)
	req.Header.Set("X-Test-Connection", "true")
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	if called {
		t.Error("callback should not be called for test connection")
	}
}

func TestBroadcastWebhookHandler_InvalidJSON(t *testing.T) {
	handler := BroadcastWebhookHandler(func(_ []BroadcastTrace) {})

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader("{bad}"))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestBroadcastWebhookHandler_WrongMethod(t *testing.T) {
	handler := BroadcastWebhookHandler(func(_ []BroadcastTrace) {})

	req := httptest.NewRequest(http.MethodGet, "/webhook", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", w.Code, http.StatusMethodNotAllowed)
	}
}

func TestBroadcastWebhookHandlerWithError_CallbackError(t *testing.T) {
	handler := BroadcastWebhookHandlerWithError(func(_ []BroadcastTrace) error {
		return errors.New("processing failed")
	})

	body := `{"resourceSpans":[]}`
	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(body))
	w := httptest.NewRecorder()
	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestExtractBroadcastTraces_IntValueTokens(t *testing.T) {
	// OTLP encodes int64 as strings in intValue
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t","spanId":"s","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.prompt_tokens","value":{"intValue":"100"}},
			{"key":"gen_ai.usage.completion_tokens","value":{"intValue":"200"}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqualInt(t, "PromptTokens", traces[0].PromptTokens, 100)
	assertEqualInt(t, "CompletionTokens", traces[0].CompletionTokens, 200)
	assertEqualInt(t, "TotalTokens", traces[0].TotalTokens, 300)
}

func TestExtractBroadcastTraces_IntValueAsNumber(t *testing.T) {
	// Some implementations send intValue as a JSON number instead of a string
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t","spanId":"s","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.prompt_tokens","value":{"intValue":100}},
			{"key":"gen_ai.usage.completion_tokens","value":{"intValue":200}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqualInt(t, "PromptTokens", traces[0].PromptTokens, 100)
	assertEqualInt(t, "CompletionTokens", traces[0].CompletionTokens, 200)
	assertEqualInt(t, "TotalTokens", traces[0].TotalTokens, 300)
}

func TestParseBroadcastTraces_Convenience(t *testing.T) {
	// Error case
	_, err := ParseBroadcastTraces([]byte("{bad}"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	// Success case
	traces, err := ParseBroadcastTraces([]byte(`{"resourceSpans":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(traces) != 0 {
		t.Errorf("got %d traces, want 0", len(traces))
	}
}

// Helpers

func assertEqual(t *testing.T, field, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %q, want %q", field, got, want)
	}
}

func assertEqualInt(t *testing.T, field string, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %d, want %d", field, got, want)
	}
}

func flexIntPtr(s string) *FlexInt {
	f := FlexInt(s)
	return &f
}

