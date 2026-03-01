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
	// Uses actual OTLP keys from real OpenRouter payloads
	payload := `{"resourceSpans":[{"resource":{"attributes":[
		{"key":"service.name","value":{"stringValue":"openrouter"}}
	]},"scopeSpans":[{"scope":{"name":"broadcast","version":"1.0"},"spans":[{
		"traceId":"aaaa","spanId":"bbbb","parentSpanId":"cccc","name":"chat completions",
		"kind":2,
		"startTimeUnixNano":"1700000000000000000",
		"endTimeUnixNano":"1700000002500000000",
		"attributes":[
			{"key":"gen_ai.response.model","value":{"stringValue":"openai/gpt-4o"}},
			{"key":"gen_ai.request.model","value":{"stringValue":"openai/gpt-4o"}},
			{"key":"gen_ai.usage.input_tokens","value":{"intValue":"150"}},
			{"key":"gen_ai.usage.output_tokens","value":{"intValue":"50"}},
			{"key":"gen_ai.usage.total_cost","value":{"doubleValue":0.0035}},
			{"key":"gen_ai.usage.input_cost","value":{"doubleValue":0.002}},
			{"key":"gen_ai.usage.output_cost","value":{"doubleValue":0.0015}},
			{"key":"gen_ai.usage.input_tokens.cached","value":{"intValue":"10"}},
			{"key":"gen_ai.usage.output_tokens.reasoning","value":{"intValue":"5"}},
			{"key":"gen_ai.operation.name","value":{"stringValue":"chat"}},
			{"key":"gen_ai.system","value":{"stringValue":"openai"}},
			{"key":"gen_ai.provider.name","value":{"stringValue":"OpenAI"}},
			{"key":"gen_ai.response.finish_reason","value":{"stringValue":"stop"}},
			{"key":"gen_ai.prompt","value":{"stringValue":"Hello world"}},
			{"key":"gen_ai.completion","value":{"stringValue":"Hi there"}},
			{"key":"openrouter.provider_slug","value":{"stringValue":"openai"}},
			{"key":"openrouter.api_key_name","value":{"stringValue":"my-key"}},
			{"key":"openrouter.entity_id","value":{"stringValue":"entity-1"}},
			{"key":"openrouter.finish_reason","value":{"stringValue":"end_turn"}},
			{"key":"openrouter.input_unit_price","value":{"doubleValue":0.000005}},
			{"key":"openrouter.output_unit_price","value":{"doubleValue":0.000015}},
			{"key":"openrouter.source","value":{"stringValue":"api"}},
			{"key":"span.type","value":{"stringValue":"llm"}},
			{"key":"span.level","value":{"stringValue":"DEFAULT"}},
			{"key":"trace.name","value":{"stringValue":"my-trace"}},
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

	// Model fields (new)
	assertEqual(t, "ResponseModel", tr.ResponseModel, "openai/gpt-4o")
	assertEqual(t, "RequestModel", tr.RequestModel, "openai/gpt-4o")
	// Backward compat: Model set from gen_ai.response.model
	assertEqual(t, "Model", tr.Model, "openai/gpt-4o")

	// Token fields (new canonical names)
	assertEqualInt(t, "InputTokens", tr.InputTokens, 150)
	assertEqualInt(t, "OutputTokens", tr.OutputTokens, 50)
	// Backward compat aliases
	assertEqualInt(t, "PromptTokens", tr.PromptTokens, 150)
	assertEqualInt(t, "CompletionTokens", tr.CompletionTokens, 50)
	assertEqualInt(t, "TotalTokens", tr.TotalTokens, 200) // computed

	// Cost fields
	assertEqualFloat(t, "TotalCost", tr.TotalCost, 0.0035)
	assertEqualFloat(t, "Cost", tr.Cost, 0.0035) // backward compat
	assertEqualFloat(t, "InputCost", tr.InputCost, 0.002)
	assertEqualFloat(t, "OutputCost", tr.OutputCost, 0.0015)

	// Token detail fields
	assertEqualInt(t, "CachedTokens", tr.CachedTokens, 10)
	assertEqualInt(t, "ReasoningTokens", tr.ReasoningTokens, 5)

	// GenAI semantic convention fields
	assertEqual(t, "OperationName", tr.OperationName, "chat")
	assertEqual(t, "System", tr.System, "openai")
	assertEqual(t, "ProviderName", tr.ProviderName, "OpenAI")
	assertEqual(t, "FinishReason", tr.FinishReason, "stop")

	// Content fields
	assertEqual(t, "Prompt", tr.Prompt, "Hello world")
	assertEqual(t, "Completion", tr.Completion, "Hi there")

	// OpenRouter-specific fields
	assertEqual(t, "ProviderSlug", tr.ProviderSlug, "openai")
	assertEqual(t, "APIKeyName", tr.APIKeyName, "my-key")
	assertEqual(t, "EntityID", tr.EntityID, "entity-1")
	assertEqual(t, "OpenRouterFinishReason", tr.OpenRouterFinishReason, "end_turn")
	assertEqualFloat(t, "InputUnitPrice", tr.InputUnitPrice, 0.000005)
	assertEqualFloat(t, "OutputUnitPrice", tr.OutputUnitPrice, 0.000015)
	assertEqual(t, "Source", tr.Source, "api")

	// Span-level fields
	assertEqual(t, "SpanType", tr.SpanType, "llm")
	assertEqual(t, "SpanLevel", tr.SpanLevel, "DEFAULT")

	// Trace-level fields
	assertEqual(t, "TraceName", tr.TraceName, "my-trace")

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

func TestExtractBroadcastTraces_BackwardCompatOldKeys(t *testing.T) {
	// Verify that old attribute keys still work for backward compatibility
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t1","spanId":"s1","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.request.model","value":{"stringValue":"old/model"}},
			{"key":"gen_ai.usage.prompt_tokens","value":{"intValue":"100"}},
			{"key":"gen_ai.usage.completion_tokens","value":{"intValue":"200"}},
			{"key":"gen_ai.usage.total_tokens","value":{"intValue":"300"}},
			{"key":"gen_ai.usage.cost","value":{"doubleValue":0.005}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr := traces[0]

	// Old fields still populated
	assertEqual(t, "Model", tr.Model, "old/model")
	assertEqualInt(t, "PromptTokens", tr.PromptTokens, 100)
	assertEqualInt(t, "CompletionTokens", tr.CompletionTokens, 200)
	assertEqualInt(t, "TotalTokens", tr.TotalTokens, 300)
	assertEqualFloat(t, "Cost", tr.Cost, 0.005)

	// New fields also populated from old keys
	assertEqual(t, "RequestModel", tr.RequestModel, "old/model")
	assertEqualInt(t, "InputTokens", tr.InputTokens, 100)
	assertEqualInt(t, "OutputTokens", tr.OutputTokens, 200)
	assertEqualFloat(t, "TotalCost", tr.TotalCost, 0.005)
}

func TestExtractBroadcastTraces_ResponseModelOverridesRequestModel(t *testing.T) {
	// When both gen_ai.response.model and gen_ai.request.model are present,
	// Model should be set from gen_ai.response.model
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t1","spanId":"s1","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.request.model","value":{"stringValue":"requested/model"}},
			{"key":"gen_ai.response.model","value":{"stringValue":"actual/model"}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr := traces[0]

	assertEqual(t, "Model", tr.Model, "actual/model")
	assertEqual(t, "ResponseModel", tr.ResponseModel, "actual/model")
	assertEqual(t, "RequestModel", tr.RequestModel, "requested/model")
}

func TestExtractBroadcastTraces_TotalTokensExplicit(t *testing.T) {
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t1","spanId":"s1","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.input_tokens","value":{"intValue":"10"}},
			{"key":"gen_ai.usage.output_tokens","value":{"intValue":"20"}},
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
		"attributes":[{"key":"gen_ai.response.model","value":{"stringValue":"test-model"}}]
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
	if received[0].ResponseModel != "test-model" {
		t.Errorf("ResponseModel = %q, want %q", received[0].ResponseModel, "test-model")
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
	// OTLP encodes int64 as strings in intValue — use new canonical keys
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t","spanId":"s","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.input_tokens","value":{"intValue":"100"}},
			{"key":"gen_ai.usage.output_tokens","value":{"intValue":"200"}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqualInt(t, "InputTokens", traces[0].InputTokens, 100)
	assertEqualInt(t, "OutputTokens", traces[0].OutputTokens, 200)
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
			{"key":"gen_ai.usage.input_tokens","value":{"intValue":100}},
			{"key":"gen_ai.usage.output_tokens","value":{"intValue":200}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertEqualInt(t, "InputTokens", traces[0].InputTokens, 100)
	assertEqualInt(t, "OutputTokens", traces[0].OutputTokens, 200)
	assertEqualInt(t, "PromptTokens", traces[0].PromptTokens, 100)
	assertEqualInt(t, "CompletionTokens", traces[0].CompletionTokens, 200)
	assertEqualInt(t, "TotalTokens", traces[0].TotalTokens, 300)
}

func TestExtractBroadcastTraces_NewFields(t *testing.T) {
	payload := `{"resourceSpans":[{"resource":{"attributes":[]},"scopeSpans":[{"spans":[{
		"traceId":"t","spanId":"s","name":"test","kind":1,
		"startTimeUnixNano":"0","endTimeUnixNano":"0",
		"attributes":[
			{"key":"gen_ai.usage.input_tokens.audio","value":{"intValue":"25"}},
			{"key":"gen_ai.usage.input_tokens.video","value":{"intValue":"30"}},
			{"key":"gen_ai.usage.output_tokens.image","value":{"intValue":"15"}},
			{"key":"span.input","value":{"stringValue":"input data"}},
			{"key":"span.output","value":{"stringValue":"output data"}},
			{"key":"trace.input","value":{"stringValue":"trace in"}},
			{"key":"trace.output","value":{"stringValue":"trace out"}}
		]
	}]}]}]}`

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tr := traces[0]

	assertEqualInt(t, "AudioInputTokens", tr.AudioInputTokens, 25)
	assertEqualInt(t, "VideoInputTokens", tr.VideoInputTokens, 30)
	assertEqualInt(t, "ImageOutputTokens", tr.ImageOutputTokens, 15)
	assertEqual(t, "SpanInput", tr.SpanInput, "input data")
	assertEqual(t, "SpanOutput", tr.SpanOutput, "output data")
	assertEqual(t, "TraceInput", tr.TraceInput, "trace in")
	assertEqual(t, "TraceOutput", tr.TraceOutput, "trace out")
}

func TestExtractBroadcastTraces_RealOpenRouterPayload(t *testing.T) {
	// Real trace attributes from the OpenRouter test-connection button.
	// All values arrive as stringValue — this validates string-based parsing
	// for numeric fields like tokens and costs.
	payload := `{"resourceSpans":[{"resource":{"attributes":[
		{"key":"service.name","value":{"stringValue":"openrouter"}}
	]},"scopeSpans":[{"spans":[{
		"traceId":"real-trace-id","spanId":"real-span-id","name":"generation",
		"kind":3,
		"startTimeUnixNano":"1772321746000000000",
		"endTimeUnixNano":"1772321748000000000",
		"attributes":[
			{"key":"environment","value":{"stringValue":"test"}},
			{"key":"gen_ai.completion","value":{"stringValue":"{\"id\":\"chatcmpl-test123\"}"}},
			{"key":"gen_ai.operation.name","value":{"stringValue":"chat"}},
			{"key":"gen_ai.prompt","value":{"stringValue":"{\"messages\":[{\"role\":\"user\",\"content\":\"What is the capital of France?\"}]}"}},
			{"key":"gen_ai.provider.name","value":{"stringValue":"openai"}},
			{"key":"gen_ai.request.frequency_penalty","value":{"stringValue":"0"}},
			{"key":"gen_ai.request.max_tokens","value":{"stringValue":"150"}},
			{"key":"gen_ai.request.presence_penalty","value":{"stringValue":"0"}},
			{"key":"gen_ai.request.temperature","value":{"stringValue":"0.7"}},
			{"key":"gen_ai.request.top_p","value":{"stringValue":"1"}},
			{"key":"gen_ai.response.finish_reason","value":{"stringValue":"stop"}},
			{"key":"gen_ai.response.finish_reasons","value":{"stringValue":"[\"stop\"]"}},
			{"key":"gen_ai.response.model","value":{"stringValue":"openai/gpt-4-turbo"}},
			{"key":"gen_ai.system","value":{"stringValue":"openai"}},
			{"key":"gen_ai.usage.input_cost","value":{"stringValue":"0.005"}},
			{"key":"gen_ai.usage.input_tokens","value":{"stringValue":"50"}},
			{"key":"gen_ai.usage.input_tokens.cached","value":{"stringValue":"20"}},
			{"key":"gen_ai.usage.output_cost","value":{"stringValue":"0.015"}},
			{"key":"gen_ai.usage.output_tokens","value":{"stringValue":"100"}},
			{"key":"gen_ai.usage.output_tokens.reasoning","value":{"stringValue":"10"}},
			{"key":"gen_ai.usage.total_cost","value":{"stringValue":"0.02"}},
			{"key":"openrouter.api_key_name","value":{"stringValue":"Test API Key"}},
			{"key":"openrouter.entity_id","value":{"stringValue":"user_2gNFMdypK9X7tRGUKgA4B9dyXQt"}},
			{"key":"openrouter.finish_reason","value":{"stringValue":"stop"}},
			{"key":"openrouter.input_unit_price","value":{"stringValue":"0.0001"}},
			{"key":"openrouter.output_unit_price","value":{"stringValue":"0.00015"}},
			{"key":"openrouter.provider_name","value":{"stringValue":"OpenAI"}},
			{"key":"openrouter.provider_slug","value":{"stringValue":"openai"}},
			{"key":"openrouter.source","value":{"stringValue":"openrouter"}},
			{"key":"openrouter.user_id","value":{"stringValue":"user_2gNFMdypK9X7tRGUKgA4B9dyXQt"}},
			{"key":"source","value":{"stringValue":"observability-test-button"}},
			{"key":"span.input","value":{"stringValue":"{\"messages\":[{\"role\":\"user\",\"content\":\"What is the capital of France?\"}]}"}},
			{"key":"span.level","value":{"stringValue":"DEFAULT"}},
			{"key":"span.metadata.test","value":{"stringValue":"true"}},
			{"key":"span.output","value":{"stringValue":"{\"id\":\"chatcmpl-test123\"}"}},
			{"key":"span.type","value":{"stringValue":"generation"}},
			{"key":"testId","value":{"stringValue":"0dc588ef-d837-4089-9420-b697a94922e6"}},
			{"key":"trace.input","value":{"stringValue":"{\"messages\":[]}"}},
			{"key":"trace.name","value":{"stringValue":"Test Trace - OpenRouter Observability"}},
			{"key":"trace.output","value":{"stringValue":"{\"id\":\"chatcmpl-test123\"}"}},
			{"key":"trace.tags","value":{"stringValue":"[\"test\",\"observability\",\"gpt-4-turbo\"]"}}
		]
	}]}]}]}` //nolint:lll

	traces, err := ParseBroadcastTraces([]byte(payload))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(traces) != 1 {
		t.Fatalf("got %d traces, want 1", len(traces))
	}

	tr := traces[0]

	// Span identity
	assertEqual(t, "TraceID", tr.TraceID, "real-trace-id")
	assertEqual(t, "SpanID", tr.SpanID, "real-span-id")
	assertEqual(t, "SpanName", tr.SpanName, "generation")

	// Model
	assertEqual(t, "ResponseModel", tr.ResponseModel, "openai/gpt-4-turbo")
	assertEqual(t, "Model", tr.Model, "openai/gpt-4-turbo") // backward compat

	// Tokens (string-encoded values)
	assertEqualInt(t, "InputTokens", tr.InputTokens, 50)
	assertEqualInt(t, "OutputTokens", tr.OutputTokens, 100)
	assertEqualInt(t, "PromptTokens", tr.PromptTokens, 50)
	assertEqualInt(t, "CompletionTokens", tr.CompletionTokens, 100)
	assertEqualInt(t, "TotalTokens", tr.TotalTokens, 150) // computed

	// Token details
	assertEqualInt(t, "CachedTokens", tr.CachedTokens, 20)
	assertEqualInt(t, "ReasoningTokens", tr.ReasoningTokens, 10)

	// Costs (string-encoded floats)
	assertEqualFloat(t, "TotalCost", tr.TotalCost, 0.02)
	assertEqualFloat(t, "Cost", tr.Cost, 0.02) // backward compat
	assertEqualFloat(t, "InputCost", tr.InputCost, 0.005)
	assertEqualFloat(t, "OutputCost", tr.OutputCost, 0.015)

	// GenAI semantic convention fields
	assertEqual(t, "OperationName", tr.OperationName, "chat")
	assertEqual(t, "System", tr.System, "openai")
	assertEqual(t, "ProviderName", tr.ProviderName, "openai")
	assertEqual(t, "FinishReason", tr.FinishReason, "stop")
	assertEqual(t, "FinishReasons", tr.FinishReasons, `["stop"]`)

	// Content
	if !strings.Contains(tr.Prompt, "What is the capital of France?") {
		t.Errorf("Prompt does not contain expected content: %s", tr.Prompt)
	}
	if !strings.Contains(tr.Completion, "chatcmpl-test123") {
		t.Errorf("Completion does not contain expected content: %s", tr.Completion)
	}

	// OpenRouter-specific fields
	assertEqual(t, "ProviderSlug", tr.ProviderSlug, "openai")
	assertEqual(t, "OpenRouterProviderName", tr.OpenRouterProviderName, "OpenAI")
	assertEqual(t, "APIKeyName", tr.APIKeyName, "Test API Key")
	assertEqual(t, "EntityID", tr.EntityID, "user_2gNFMdypK9X7tRGUKgA4B9dyXQt")
	assertEqual(t, "OpenRouterUserID", tr.OpenRouterUserID, "user_2gNFMdypK9X7tRGUKgA4B9dyXQt")
	assertEqual(t, "OpenRouterFinishReason", tr.OpenRouterFinishReason, "stop")
	assertEqualFloat(t, "InputUnitPrice", tr.InputUnitPrice, 0.0001)
	assertEqualFloat(t, "OutputUnitPrice", tr.OutputUnitPrice, 0.00015)
	assertEqual(t, "Source", tr.Source, "openrouter")

	// Span-level fields
	assertEqual(t, "SpanType", tr.SpanType, "generation")
	assertEqual(t, "SpanLevel", tr.SpanLevel, "DEFAULT")
	if !strings.Contains(tr.SpanInput, "What is the capital of France?") {
		t.Errorf("SpanInput does not contain expected content: %s", tr.SpanInput)
	}
	if !strings.Contains(tr.SpanOutput, "chatcmpl-test123") {
		t.Errorf("SpanOutput does not contain expected content: %s", tr.SpanOutput)
	}

	// Span metadata
	if tr.SpanMetadata["test"] != "true" {
		t.Errorf("SpanMetadata[test] = %q, want %q", tr.SpanMetadata["test"], "true")
	}

	// Trace-level fields
	assertEqual(t, "TraceName", tr.TraceName, "Test Trace - OpenRouter Observability")
	assertEqual(t, "TraceTags", tr.TraceTags, `["test","observability","gpt-4-turbo"]`)
	if tr.TraceInput == "" {
		t.Error("TraceInput should not be empty")
	}
	if tr.TraceOutput == "" {
		t.Error("TraceOutput should not be empty")
	}

	// Unmapped attributes go to RawAttributes
	if tr.RawAttributes["source"] != "observability-test-button" {
		t.Errorf("RawAttributes[source] = %q, want %q", tr.RawAttributes["source"], "observability-test-button")
	}
	if tr.RawAttributes["testId"] != "0dc588ef-d837-4089-9420-b697a94922e6" {
		t.Errorf("RawAttributes[testId] = %q, want %q", tr.RawAttributes["testId"], "0dc588ef-d837-4089-9420-b697a94922e6")
	}
	if tr.RawAttributes["environment"] != "test" {
		t.Errorf("RawAttributes[environment] = %q, want %q", tr.RawAttributes["environment"], "test")
	}
	// gen_ai.request.* params should end up in RawAttributes
	if tr.RawAttributes["gen_ai.request.temperature"] != "0.7" {
		t.Errorf("RawAttributes[gen_ai.request.temperature] = %q, want %q", tr.RawAttributes["gen_ai.request.temperature"], "0.7")
	}
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

func assertEqualFloat(t *testing.T, field string, got, want float64) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}

func flexIntPtr(s string) *FlexInt {
	f := FlexInt(s)
	return &f
}
