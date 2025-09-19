package sse

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Event
	}{
		{
			name: "simple event",
			input: `data: hello world

`,
			expected: []Event{
				{Data: []byte("hello world")},
			},
		},
		{
			name: "event with type",
			input: `event: message
data: hello

`,
			expected: []Event{
				{Type: "message", Data: []byte("hello")},
			},
		},
		{
			name: "event with id",
			input: `id: 123
data: test

`,
			expected: []Event{
				{ID: "123", Data: []byte("test")},
			},
		},
		{
			name: "multiline data",
			input: `data: line 1
data: line 2
data: line 3

`,
			expected: []Event{
				{Data: []byte("line 1\nline 2\nline 3")},
			},
		},
		{
			name: "event with retry",
			input: `retry: 5000
data: reconnect test

`,
			expected: []Event{
				{Data: []byte("reconnect test"), Retry: durationPtr(5 * time.Second)},
			},
		},
		{
			name: "multiple events",
			input: `data: first

data: second

event: custom
data: third

`,
			expected: []Event{
				{Data: []byte("first")},
				{Data: []byte("second")},
				{Type: "custom", Data: []byte("third")},
			},
		},
		{
			name: "event with comment",
			input: `: this is a comment
data: actual data

`,
			expected: []Event{
				{Data: []byte("actual data"), Comment: " this is a comment"},
			},
		},
		{
			name: "end of stream",
			input: `data: [DONE]

`,
			expected: []Event{
				{Data: []byte("[DONE]")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewParser(strings.NewReader(tt.input))
			var events []Event

			for i := 0; i < len(tt.expected); i++ {
				event, err := parser.Next()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				events = append(events, *event)
			}

			if len(events) != len(tt.expected) {
				t.Fatalf("expected %d events, got %d", len(tt.expected), len(events))
			}

			for i, event := range events {
				expected := tt.expected[i]
				if event.ID != expected.ID {
					t.Errorf("event %d: expected ID %q, got %q", i, expected.ID, event.ID)
				}
				if event.Type != expected.Type {
					t.Errorf("event %d: expected Type %q, got %q", i, expected.Type, event.Type)
				}
				if !bytes.Equal(event.Data, expected.Data) {
					t.Errorf("event %d: expected Data %q, got %q", i, expected.Data, event.Data)
				}
				if (event.Retry == nil) != (expected.Retry == nil) {
					t.Errorf("event %d: retry nil mismatch", i)
				}
				if event.Retry != nil && expected.Retry != nil && *event.Retry != *expected.Retry {
					t.Errorf("event %d: expected Retry %v, got %v", i, *expected.Retry, *event.Retry)
				}
			}
		})
	}
}

func TestScanner(t *testing.T) {
	input := `data: first

event: update
data: second

id: 42
data: third

`

	scanner := NewScanner(strings.NewReader(input))

	// First event
	if !scanner.Scan() {
		t.Fatal("expected first event")
	}
	event := scanner.Event()
	if string(event.Data) != "first" {
		t.Errorf("expected first event data to be 'first', got %q", event.Data)
	}

	// Second event
	if !scanner.Scan() {
		t.Fatal("expected second event")
	}
	event = scanner.Event()
	if event.Type != "update" || string(event.Data) != "second" {
		t.Errorf("unexpected second event: type=%q, data=%q", event.Type, event.Data)
	}

	// Third event
	if !scanner.Scan() {
		t.Fatal("expected third event")
	}
	event = scanner.Event()
	if event.ID != "42" || string(event.Data) != "third" {
		t.Errorf("unexpected third event: id=%q, data=%q", event.ID, event.Data)
	}

	// No more events
	if scanner.Scan() {
		t.Error("expected no more events")
	}

	if err := scanner.Err(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestIsEndOfStream(t *testing.T) {
	tests := []struct {
		input    []byte
		expected bool
	}{
		{[]byte("[DONE]"), true},
		{[]byte(" [DONE] "), true},
		{[]byte("\n[DONE]\n"), true},
		{[]byte("not done"), false},
		{[]byte(""), false},
		{[]byte("DONE"), false},
	}

	for _, tt := range tests {
		result := IsEndOfStream(tt.input)
		if result != tt.expected {
			t.Errorf("IsEndOfStream(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := NewWriter(&buf)

	// Write a simple event
	err := writer.WriteEvent(&Event{
		ID:   "123",
		Type: "message",
		Data: []byte("hello world"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "id: 123\nevent: message\ndata: hello world\n\n"
	if buf.String() != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, buf.String())
	}

	// Write multiline data
	buf.Reset()
	err = writer.WriteEvent(&Event{
		Data: []byte("line 1\nline 2\nline 3"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected = "data: line 1\ndata: line 2\ndata: line 3\n\n"
	if buf.String() != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, buf.String())
	}

	// Write retry
	buf.Reset()
	err = writer.WriteRetry(3 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected = "retry: 3000\n\n"
	if buf.String() != expected {
		t.Errorf("expected:\n%q\ngot:\n%q", expected, buf.String())
	}
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}