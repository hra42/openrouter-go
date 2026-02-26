package openrouter

import (
	"testing"
)

// strPtr is declared in broadcast_test.go

func makeModelEndpoints() []ModelEndpoint {
	return []ModelEndpoint{
		{
			ProviderName:        "provider-a",
			Pricing:             ModelEndpointPricing{Prompt: "0.001", Completion: "0.002"},
			SupportedParameters: []string{"tools", "stream", "temperature"},
			Status:              1.0,
			Quantization:        new("fp16"),
		},
		{
			ProviderName:        "provider-b",
			Pricing:             ModelEndpointPricing{Prompt: "0.005", Completion: "0.010"},
			SupportedParameters: []string{"stream", "temperature"},
			Status:              0.8,
			Quantization:        new("int8"),
		},
		{
			ProviderName:        "provider-a",
			Pricing:             ModelEndpointPricing{Prompt: "0.0005", Completion: "0.001"},
			SupportedParameters: []string{"tools", "stream"},
			Status:              0.95,
			Quantization:        nil,
		},
	}
}

func makeZDREndpoints() []PublicEndpoint {
	return []PublicEndpoint{
		{
			ProviderName:        "provider-x",
			Pricing:             PublicEndpointPricing{Prompt: "0.002", Completion: "0.004"},
			SupportedParameters: []string{"tools", "stream"},
			Status:              1.0,
			Quantization:        new("fp16"),
		},
		{
			ProviderName:        "provider-y",
			Pricing:             PublicEndpointPricing{Prompt: "0.001", Completion: "0.002"},
			SupportedParameters: []string{"stream"},
			Status:              0.9,
			Quantization:        new("int4"),
		},
	}
}

func TestFilterModelEndpoints(t *testing.T) {
	endpoints := makeModelEndpoints()

	t.Run("filter by provider", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{ProviderName: "provider-a"})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by required params", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{RequiredParameters: []string{"tools"}})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by max prompt price", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{MaxPromptPrice: 0.002})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by max completion price", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{MaxCompletionPrice: 0.003})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by min status", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{MinStatus: 0.9})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("filter by quantization", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{Quantization: "fp16"})
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("combined filter", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{
			ProviderName:       "provider-a",
			RequiredParameters: []string{"tools"},
			MaxPromptPrice:     0.002,
		})
		if len(result) != 2 {
			t.Errorf("expected 2, got %d", len(result))
		}
	})

	t.Run("no matches", func(t *testing.T) {
		result := FilterModelEndpoints(endpoints, EndpointFilter{ProviderName: "nonexistent"})
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := FilterModelEndpoints(nil, EndpointFilter{})
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})
}

func TestFilterZDREndpoints(t *testing.T) {
	endpoints := makeZDREndpoints()

	t.Run("filter by provider", func(t *testing.T) {
		result := FilterZDREndpoints(endpoints, EndpointFilter{ProviderName: "provider-x"})
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by required params", func(t *testing.T) {
		result := FilterZDREndpoints(endpoints, EndpointFilter{RequiredParameters: []string{"tools"}})
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("filter by quantization", func(t *testing.T) {
		result := FilterZDREndpoints(endpoints, EndpointFilter{Quantization: "int4"})
		if len(result) != 1 {
			t.Errorf("expected 1, got %d", len(result))
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := FilterZDREndpoints(nil, EndpointFilter{})
		if len(result) != 0 {
			t.Errorf("expected 0, got %d", len(result))
		}
	})
}

func TestCheapestModelEndpoint(t *testing.T) {
	endpoints := makeModelEndpoints()

	t.Run("finds cheapest", func(t *testing.T) {
		result := CheapestModelEndpoint(endpoints)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Pricing.Prompt != "0.0005" {
			t.Errorf("expected prompt price 0.0005, got %s", result.Pricing.Prompt)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := CheapestModelEndpoint(nil)
		if result != nil {
			t.Error("expected nil for empty slice")
		}
	})
}

func TestCheapestZDREndpoint(t *testing.T) {
	endpoints := makeZDREndpoints()

	t.Run("finds cheapest", func(t *testing.T) {
		result := CheapestZDREndpoint(endpoints)
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if result.Pricing.Prompt != "0.001" {
			t.Errorf("expected prompt price 0.001, got %s", result.Pricing.Prompt)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		result := CheapestZDREndpoint(nil)
		if result != nil {
			t.Error("expected nil for empty slice")
		}
	})
}

func TestEndpointSupportsParameter(t *testing.T) {
	params := []string{"tools", "stream", "temperature"}

	if !EndpointSupportsParameter(params, "tools") {
		t.Error("expected true for 'tools'")
	}
	if EndpointSupportsParameter(params, "nonexistent") {
		t.Error("expected false for 'nonexistent'")
	}
	if EndpointSupportsParameter(nil, "tools") {
		t.Error("expected false for nil params")
	}
}
