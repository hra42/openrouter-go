package openrouter

// EndpointFilter defines criteria for filtering model endpoints.
type EndpointFilter struct {
	// ProviderName filters by provider name (exact match).
	ProviderName string
	// RequiredParameters filters to endpoints supporting all listed parameters.
	RequiredParameters []string
	// MaxPromptPrice filters to endpoints with prompt price at or below this value (per million tokens).
	// Zero means no filter.
	MaxPromptPrice float64
	// MaxCompletionPrice filters to endpoints with completion price at or below this value (per million tokens).
	// Zero means no filter.
	MaxCompletionPrice float64
	// MinStatus filters to endpoints with status at or above this value.
	// Zero means no filter.
	MinStatus float64
	// Quantization filters by quantization type (exact match).
	Quantization string
}

// EndpointSupportsParameter checks if a parameter list contains the given parameter.
func EndpointSupportsParameter(params []string, param string) bool {
	return hasParameter(params, param)
}

// FilterModelEndpoints returns endpoints matching all criteria in the filter.
func FilterModelEndpoints(endpoints []ModelEndpoint, filter EndpointFilter) []ModelEndpoint {
	var result []ModelEndpoint
	for i := range endpoints {
		if matchModelEndpoint(&endpoints[i], &filter) {
			result = append(result, endpoints[i])
		}
	}
	return result
}

// FilterZDREndpoints returns public endpoints matching all criteria in the filter.
func FilterZDREndpoints(endpoints []PublicEndpoint, filter EndpointFilter) []PublicEndpoint {
	var result []PublicEndpoint
	for i := range endpoints {
		if matchZDREndpoint(&endpoints[i], &filter) {
			result = append(result, endpoints[i])
		}
	}
	return result
}

// CheapestModelEndpoint returns the endpoint with the lowest prompt price.
// Returns nil if the slice is empty.
func CheapestModelEndpoint(endpoints []ModelEndpoint) *ModelEndpoint {
	if len(endpoints) == 0 {
		return nil
	}
	cheapest := &endpoints[0]
	cheapestPrice := ParsePricing(cheapest.Pricing.Prompt)
	for i := 1; i < len(endpoints); i++ {
		price := ParsePricing(endpoints[i].Pricing.Prompt)
		if price < cheapestPrice {
			cheapest = &endpoints[i]
			cheapestPrice = price
		}
	}
	return cheapest
}

// CheapestZDREndpoint returns the public endpoint with the lowest prompt price.
// Returns nil if the slice is empty.
func CheapestZDREndpoint(endpoints []PublicEndpoint) *PublicEndpoint {
	if len(endpoints) == 0 {
		return nil
	}
	cheapest := &endpoints[0]
	cheapestPrice := ParsePricing(cheapest.Pricing.Prompt)
	for i := 1; i < len(endpoints); i++ {
		price := ParsePricing(endpoints[i].Pricing.Prompt)
		if price < cheapestPrice {
			cheapest = &endpoints[i]
			cheapestPrice = price
		}
	}
	return cheapest
}

func matchModelEndpoint(ep *ModelEndpoint, f *EndpointFilter) bool {
	if f.ProviderName != "" && ep.ProviderName != f.ProviderName {
		return false
	}
	for _, param := range f.RequiredParameters {
		if !hasParameter(ep.SupportedParameters, param) {
			return false
		}
	}
	if f.MaxPromptPrice > 0 {
		price := ParsePricing(ep.Pricing.Prompt)
		if price > f.MaxPromptPrice {
			return false
		}
	}
	if f.MaxCompletionPrice > 0 {
		price := ParsePricing(ep.Pricing.Completion)
		if price > f.MaxCompletionPrice {
			return false
		}
	}
	if f.MinStatus > 0 && ep.Status < f.MinStatus {
		return false
	}
	if f.Quantization != "" {
		if ep.Quantization == nil || *ep.Quantization != f.Quantization {
			return false
		}
	}
	return true
}

func matchZDREndpoint(ep *PublicEndpoint, f *EndpointFilter) bool {
	if f.ProviderName != "" && ep.ProviderName != f.ProviderName {
		return false
	}
	for _, param := range f.RequiredParameters {
		if !hasParameter(ep.SupportedParameters, param) {
			return false
		}
	}
	if f.MaxPromptPrice > 0 {
		price := ParsePricing(ep.Pricing.Prompt)
		if price > f.MaxPromptPrice {
			return false
		}
	}
	if f.MaxCompletionPrice > 0 {
		price := ParsePricing(ep.Pricing.Completion)
		if price > f.MaxCompletionPrice {
			return false
		}
	}
	if f.MinStatus > 0 && ep.Status < f.MinStatus {
		return false
	}
	if f.Quantization != "" {
		if ep.Quantization == nil || *ep.Quantization != f.Quantization {
			return false
		}
	}
	return true
}
