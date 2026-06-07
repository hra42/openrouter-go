package openrouter

import (
	"encoding/json"
	"testing"
)

// TestUsageCostDecoding verifies that the usage accounting fields decode
// correctly from a real OpenRouter response. OpenRouter always returns these
// fields; no request flag is required. The fixture matches the live wire format
// observed from /chat/completions (including BYOK and upstream cost details).
func TestUsageCostDecoding(t *testing.T) {
	body := `{
		"id": "gen-123",
		"model": "openai/gpt-4o-mini",
		"choices": [{"index": 0, "message": {"role": "assistant", "content": "hi"}, "finish_reason": "stop"}],
		"usage": {
			"prompt_tokens": 13,
			"completion_tokens": 2,
			"total_tokens": 15,
			"cost": 0,
			"is_byok": true,
			"prompt_tokens_details": {"cached_tokens": 1, "cache_write_tokens": 2, "audio_tokens": 3, "video_tokens": 4},
			"cost_details": {
				"upstream_inference_cost": 0.00000315,
				"upstream_inference_prompt_cost": 0.00000195,
				"upstream_inference_completions_cost": 0.0000012
			},
			"completion_tokens_details": {"reasoning_tokens": 5, "image_tokens": 6, "audio_tokens": 7}
		}
	}`

	var resp ChatCompletionResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	u := resp.Usage
	if u.PromptTokens != 13 || u.CompletionTokens != 2 || u.TotalTokens != 15 {
		t.Errorf("token counts wrong: %+v", u)
	}
	if u.Cost == nil || *u.Cost != 0 {
		t.Errorf("Cost = %v, want pointer to 0", u.Cost)
	}
	if !u.IsBYOK {
		t.Error("IsBYOK = false, want true")
	}

	cd := u.CostDetails
	if cd == nil {
		t.Fatal("CostDetails is nil")
	}
	if cd.UpstreamInferenceCost == nil || *cd.UpstreamInferenceCost != 0.00000315 {
		t.Errorf("UpstreamInferenceCost wrong: %v", cd.UpstreamInferenceCost)
	}
	if cd.UpstreamInferencePromptCost == nil || *cd.UpstreamInferencePromptCost != 0.00000195 {
		t.Errorf("UpstreamInferencePromptCost wrong: %v", cd.UpstreamInferencePromptCost)
	}
	if cd.UpstreamInferenceCompletionsCost == nil || *cd.UpstreamInferenceCompletionsCost != 0.0000012 {
		t.Errorf("UpstreamInferenceCompletionsCost wrong: %v", cd.UpstreamInferenceCompletionsCost)
	}

	pd := u.PromptTokensDetails
	if pd == nil || pd.CachedTokens != 1 || pd.CacheWriteTokens != 2 || pd.AudioTokens != 3 || pd.VideoTokens != 4 {
		t.Errorf("PromptTokensDetails wrong: %+v", pd)
	}

	ct := u.CompletionTokensDetails
	if ct == nil || ct.ReasoningTokens != 5 || ct.ImageTokens != 6 || ct.AudioTokens != 7 {
		t.Errorf("CompletionTokensDetails wrong: %+v", ct)
	}
}

// TestUsageNonBYOK verifies a non-BYOK response decodes with IsBYOK false and no
// upstream cost details.
func TestUsageNonBYOK(t *testing.T) {
	var resp ChatCompletionResponse
	body := `{"id":"x","usage":{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3,"cost":0.001}}`
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.Usage.IsBYOK {
		t.Error("IsBYOK = true, want false")
	}
	if resp.Usage.CostDetails != nil {
		t.Errorf("CostDetails should be nil when absent: %+v", resp.Usage.CostDetails)
	}
	if resp.Usage.Cost == nil || *resp.Usage.Cost != 0.001 {
		t.Errorf("Cost = %v, want pointer to 0.001", resp.Usage.Cost)
	}
}

// TestUsageCostOmittedWhenAbsent verifies the optional fields are omitted from
// JSON when not set, so the SDK does not emit noise.
func TestUsageCostOmittedWhenAbsent(t *testing.T) {
	u := Usage{PromptTokens: 1, CompletionTokens: 2, TotalTokens: 3}
	out, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := string(out); got != `{"prompt_tokens":1,"completion_tokens":2,"total_tokens":3}` {
		t.Errorf("optional fields should be omitted, got %s", got)
	}
}
