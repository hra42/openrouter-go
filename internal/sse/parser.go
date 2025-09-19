// Package sse provides Server-Sent Events parsing functionality.
package sse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Event represents a Server-Sent Event.
type Event struct {
	ID      string
	Type    string
	Data    []byte
	Retry   *time.Duration
	Comment string
}

// Parser provides SSE event parsing from an io.Reader.
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new SSE parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

// Next reads and returns the next SSE event.
// Returns io.EOF when the stream ends.
func (p *Parser) Next() (*Event, error) {
	event := &Event{}
	var dataBuffer bytes.Buffer
	hasData := false

	for {
		line, err := p.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// If we have accumulated data, return it as an event
				if hasData {
					event.Data = dataBuffer.Bytes()
					return event, nil
				}
			}
			return nil, err
		}

		// Remove trailing newline and carriage return
		line = bytes.TrimSuffix(line, []byte("\n"))
		line = bytes.TrimSuffix(line, []byte("\r"))

		// Empty line signals end of event
		if len(line) == 0 {
			if hasData || event.ID != "" || event.Type != "" {
				event.Data = dataBuffer.Bytes()
				return event, nil
			}
			// Empty event, continue reading
			continue
		}

		// Skip comments (lines starting with :)
		if bytes.HasPrefix(line, []byte(":")) {
			event.Comment = string(bytes.TrimPrefix(line, []byte(":")))
			continue
		}

		// Parse field
		field, value, found := parseField(line)
		if !found {
			// Invalid line format, skip
			continue
		}

		switch field {
		case "id":
			event.ID = value
		case "event":
			event.Type = value
		case "data":
			if hasData {
				dataBuffer.WriteByte('\n')
			}
			dataBuffer.WriteString(value)
			hasData = true
		case "retry":
			if ms, err := strconv.ParseInt(value, 10, 64); err == nil {
				retry := time.Duration(ms) * time.Millisecond
				event.Retry = &retry
			}
		}
	}
}

// parseField parses a field line into field name and value.
func parseField(line []byte) (field, value string, found bool) {
	colonIndex := bytes.IndexByte(line, ':')
	if colonIndex == -1 {
		// No colon found, treat entire line as field name with empty value
		return string(line), "", true
	}

	field = string(line[:colonIndex])

	// Skip optional space after colon
	valueStart := colonIndex + 1
	if valueStart < len(line) && line[valueStart] == ' ' {
		valueStart++
	}

	if valueStart < len(line) {
		value = string(line[valueStart:])
	}

	return field, value, true
}

// Scanner provides a convenient interface for iterating over SSE events.
type Scanner struct {
	parser *Parser
	event  *Event
	err    error
}

// NewScanner creates a new SSE scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		parser: NewParser(r),
	}
}

// Scan advances to the next event.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	s.event, s.err = s.parser.Next()
	return s.err == nil
}

// Event returns the current event.
func (s *Scanner) Event() *Event {
	return s.event
}

// Err returns the error that caused Scan to return false.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// IsEndOfStream checks if the data indicates end of stream.
func IsEndOfStream(data []byte) bool {
	trimmed := bytes.TrimSpace(data)
	return bytes.Equal(trimmed, []byte("[DONE]"))
}

// ParseEventStream parses a complete SSE stream and returns all events.
func ParseEventStream(r io.Reader) ([]*Event, error) {
	var events []*Event
	scanner := NewScanner(r)

	for scanner.Scan() {
		event := scanner.Event()

		// Check for end of stream marker
		if IsEndOfStream(event.Data) {
			break
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return events, err
	}

	return events, nil
}

// Writer provides SSE event writing functionality.
type Writer struct {
	w io.Writer
}

// NewWriter creates a new SSE writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// WriteEvent writes an SSE event.
func (w *Writer) WriteEvent(event *Event) error {
	var buf bytes.Buffer

	// Write ID if present
	if event.ID != "" {
		fmt.Fprintf(&buf, "id: %s\n", event.ID)
	}

	// Write event type if present
	if event.Type != "" {
		fmt.Fprintf(&buf, "event: %s\n", event.Type)
	}

	// Write retry if present
	if event.Retry != nil {
		fmt.Fprintf(&buf, "retry: %d\n", int(event.Retry.Milliseconds()))
	}

	// Write data
	if len(event.Data) > 0 {
		lines := strings.Split(string(event.Data), "\n")
		for _, line := range lines {
			fmt.Fprintf(&buf, "data: %s\n", line)
		}
	}

	// Write comment if present
	if event.Comment != "" {
		fmt.Fprintf(&buf, ": %s\n", event.Comment)
	}

	// Write event separator
	buf.WriteByte('\n')

	_, err := w.w.Write(buf.Bytes())
	return err
}

// WriteComment writes a comment line.
func (w *Writer) WriteComment(comment string) error {
	_, err := fmt.Fprintf(w.w, ": %s\n", comment)
	return err
}

// WriteRetry writes a retry directive.
func (w *Writer) WriteRetry(duration time.Duration) error {
	_, err := fmt.Fprintf(w.w, "retry: %d\n\n", int(duration.Milliseconds()))
	return err
}