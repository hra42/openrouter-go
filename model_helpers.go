package openrouter

import "slices"

// hasParameter checks if a model supports a specific parameter.
func hasParameter(params []string, param string) bool {
	return slices.Contains(params, param)
}

// hasInputModality checks if a model supports a specific input modality.
func hasInputModality(modalities []string, modality string) bool {
	return slices.Contains(modalities, modality)
}

// hasOutputModality checks if a model supports a specific output modality.
func hasOutputModality(modalities []string, modality string) bool {
	return slices.Contains(modalities, modality)
}

// SupportsTools returns true if the model supports tool calling.
func (m *Model) SupportsTools() bool {
	return hasParameter(m.SupportedParameters, "tools")
}

// SupportsVision returns true if the model supports image inputs.
func (m *Model) SupportsVision() bool {
	return hasInputModality(m.Architecture.InputModalities, "image")
}

// SupportsAudio returns true if the model supports audio inputs.
func (m *Model) SupportsAudio() bool {
	return hasInputModality(m.Architecture.InputModalities, "audio")
}

// SupportsJSON returns true if the model supports structured JSON output via response_format.
func (m *Model) SupportsJSON() bool {
	return hasParameter(m.SupportedParameters, "response_format")
}

// SupportsStreaming returns true if the model supports streaming responses.
func (m *Model) SupportsStreaming() bool {
	return hasParameter(m.SupportedParameters, "stream")
}

// SupportsParameter returns true if the model supports the given parameter.
func (m *Model) SupportsParameter(param string) bool {
	return hasParameter(m.SupportedParameters, param)
}

// HasInputModality returns true if the model accepts the given input modality.
func (m *Model) HasInputModality(modality string) bool {
	return hasInputModality(m.Architecture.InputModalities, modality)
}

// HasOutputModality returns true if the model produces the given output modality.
func (m *Model) HasOutputModality(modality string) bool {
	return hasOutputModality(m.Architecture.OutputModalities, modality)
}

// FilterModels returns a subset of models that match the given predicate.
func FilterModels(models []Model, predicate func(*Model) bool) []Model {
	var result []Model
	for i := range models {
		if predicate(&models[i]) {
			result = append(result, models[i])
		}
	}
	return result
}
