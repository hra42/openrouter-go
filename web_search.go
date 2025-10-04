package openrouter

// WebSearchEngine represents the search engine options for the web plugin.
type WebSearchEngine string

const (
	// WebSearchEngineNative uses the model provider's built-in web search
	WebSearchEngineNative WebSearchEngine = "native"
	// WebSearchEngineExa uses Exa's search API for web results
	WebSearchEngineExa WebSearchEngine = "exa"
	// WebSearchEngineAuto automatically selects the best available engine
	WebSearchEngineAuto WebSearchEngine = ""
)

// WebSearchContextSize represents the amount of search context to retrieve.
type WebSearchContextSize string

const (
	// WebSearchContextLow provides minimal search context
	WebSearchContextLow WebSearchContextSize = "low"
	// WebSearchContextMedium provides moderate search context
	WebSearchContextMedium WebSearchContextSize = "medium"
	// WebSearchContextHigh provides extensive search context
	WebSearchContextHigh WebSearchContextSize = "high"
)

// NewWebPlugin creates a new web search plugin configuration.
// By default, uses automatic engine selection and 5 max results.
func NewWebPlugin() Plugin {
	return Plugin{
		ID:         "web",
		MaxResults: 5,
	}
}

// NewWebPluginWithOptions creates a web search plugin with custom options.
func NewWebPluginWithOptions(engine WebSearchEngine, maxResults int, searchPrompt string) Plugin {
	return Plugin{
		ID:           "web",
		Engine:       string(engine),
		MaxResults:   maxResults,
		SearchPrompt: searchPrompt,
	}
}

// DefaultWebSearchPrompt returns the default search prompt template.
// The date should be substituted with the current date when used.
func DefaultWebSearchPrompt(date string) string {
	return "A web search was conducted on `" + date + "`. Incorporate the following web search results into your response.\n\n" +
		"IMPORTANT: Cite them using markdown links named using the domain of the source.\n" +
		"Example: [nytimes.com](https://nytimes.com/some-page)."
}

// WithOnlineModel appends ":online" to a model name to enable web search.
// This is a shortcut for using the web plugin.
func WithOnlineModel(model string) string {
	return model + ":online"
}

// ParseAnnotations extracts URL citations from message annotations.
func ParseAnnotations(annotations []Annotation) []URLCitation {
	var citations []URLCitation
	for _, annotation := range annotations {
		if annotation.Type == "url_citation" && annotation.URLCitation != nil {
			citations = append(citations, *annotation.URLCitation)
		}
	}
	return citations
}
