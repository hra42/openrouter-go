package openrouter

import (
	"testing"
)

func TestModelSupportsTools(t *testing.T) {
	tests := []struct {
		name   string
		params []string
		want   bool
	}{
		{"supports tools", []string{"temperature", "tools", "stream"}, true},
		{"no tools", []string{"temperature", "stream"}, false},
		{"empty params", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{SupportedParameters: tt.params}
			if got := m.SupportsTools(); got != tt.want {
				t.Errorf("SupportsTools() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelSupportsVision(t *testing.T) {
	tests := []struct {
		name       string
		modalities []string
		want       bool
	}{
		{"supports image", []string{"text", "image"}, true},
		{"text only", []string{"text"}, false},
		{"empty", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{Architecture: ModelArchitecture{InputModalities: tt.modalities}}
			if got := m.SupportsVision(); got != tt.want {
				t.Errorf("SupportsVision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelSupportsAudio(t *testing.T) {
	tests := []struct {
		name       string
		modalities []string
		want       bool
	}{
		{"supports audio", []string{"text", "audio"}, true},
		{"text only", []string{"text"}, false},
		{"empty", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{Architecture: ModelArchitecture{InputModalities: tt.modalities}}
			if got := m.SupportsAudio(); got != tt.want {
				t.Errorf("SupportsAudio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelSupportsJSON(t *testing.T) {
	tests := []struct {
		name   string
		params []string
		want   bool
	}{
		{"supports response_format", []string{"temperature", "response_format"}, true},
		{"no response_format", []string{"temperature"}, false},
		{"empty", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{SupportedParameters: tt.params}
			if got := m.SupportsJSON(); got != tt.want {
				t.Errorf("SupportsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelSupportsStreaming(t *testing.T) {
	tests := []struct {
		name   string
		params []string
		want   bool
	}{
		{"supports stream", []string{"temperature", "stream"}, true},
		{"no stream", []string{"temperature"}, false},
		{"empty", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{SupportedParameters: tt.params}
			if got := m.SupportsStreaming(); got != tt.want {
				t.Errorf("SupportsStreaming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelSupportsParameter(t *testing.T) {
	m := &Model{SupportedParameters: []string{"temperature", "top_p", "tools"}}

	if !m.SupportsParameter("temperature") {
		t.Error("expected SupportsParameter('temperature') = true")
	}
	if m.SupportsParameter("nonexistent") {
		t.Error("expected SupportsParameter('nonexistent') = false")
	}
}

func TestModelHasInputModality(t *testing.T) {
	m := &Model{Architecture: ModelArchitecture{InputModalities: []string{"text", "image"}}}

	if !m.HasInputModality("text") {
		t.Error("expected HasInputModality('text') = true")
	}
	if !m.HasInputModality("image") {
		t.Error("expected HasInputModality('image') = true")
	}
	if m.HasInputModality("audio") {
		t.Error("expected HasInputModality('audio') = false")
	}
}

func TestModelHasOutputModality(t *testing.T) {
	m := &Model{Architecture: ModelArchitecture{OutputModalities: []string{"text", "image"}}}

	if !m.HasOutputModality("text") {
		t.Error("expected HasOutputModality('text') = true")
	}
	if m.HasOutputModality("audio") {
		t.Error("expected HasOutputModality('audio') = false")
	}
}

func TestFilterModels(t *testing.T) {
	models := []Model{
		{ID: "model-a", SupportedParameters: []string{"tools", "stream"}},
		{ID: "model-b", SupportedParameters: []string{"stream"}},
		{ID: "model-c", SupportedParameters: []string{"tools", "response_format"}},
	}

	t.Run("filter by tools support", func(t *testing.T) {
		result := FilterModels(models, func(m *Model) bool {
			return m.SupportsTools()
		})
		if len(result) != 2 {
			t.Errorf("expected 2 models, got %d", len(result))
		}
		if result[0].ID != "model-a" || result[1].ID != "model-c" {
			t.Error("unexpected models in result")
		}
	})

	t.Run("filter with no matches", func(t *testing.T) {
		result := FilterModels(models, func(m *Model) bool {
			return m.SupportsAudio()
		})
		if len(result) != 0 {
			t.Errorf("expected 0 models, got %d", len(result))
		}
	})

	t.Run("filter empty slice", func(t *testing.T) {
		result := FilterModels(nil, func(m *Model) bool { return true })
		if len(result) != 0 {
			t.Errorf("expected 0 models, got %d", len(result))
		}
	})

	t.Run("combined predicate", func(t *testing.T) {
		result := FilterModels(models, func(m *Model) bool {
			return m.SupportsTools() && m.SupportsStreaming()
		})
		if len(result) != 1 {
			t.Errorf("expected 1 model, got %d", len(result))
		}
		if len(result) > 0 && result[0].ID != "model-a" {
			t.Errorf("expected model-a, got %s", result[0].ID)
		}
	})
}
