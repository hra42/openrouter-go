package openrouter

import (
	"strings"
	"testing"
)

func TestGenerateCodeVerifier(t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if verifier == "" {
		t.Fatal("expected non-empty verifier")
	}

	// 32 bytes base64url-encoded without padding = 43 characters
	if len(verifier) != 43 {
		t.Errorf("expected verifier length 43, got %d", len(verifier))
	}

	// Verify uniqueness
	verifier2, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if verifier == verifier2 {
		t.Error("expected unique verifiers, got identical values")
	}
}

func TestCreateS256CodeChallenge(t *testing.T) {
	// Known test vector from RFC 7636 Appendix B:
	// verifier: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	// expected challenge: "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	expected := "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"

	challenge := CreateS256CodeChallenge(verifier)
	if challenge != expected {
		t.Errorf("expected challenge %q, got %q", expected, challenge)
	}
}

func TestBuildAuthURL(t *testing.T) {
	t.Run("all params", func(t *testing.T) {
		authURL, err := BuildAuthURL("https://openrouter.ai/auth", AuthURLParams{
			CallbackURL:         "https://myapp.com/callback",
			CodeChallenge:       "test-challenge",
			CodeChallengeMethod: CodeChallengeMethodS256,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.HasPrefix(authURL, "https://openrouter.ai/auth?") {
			t.Errorf("expected URL to start with 'https://openrouter.ai/auth?', got %q", authURL)
		}
		if !strings.Contains(authURL, "callback_url=") {
			t.Error("expected URL to contain callback_url parameter")
		}
		if !strings.Contains(authURL, "code_challenge=test-challenge") {
			t.Error("expected URL to contain code_challenge parameter")
		}
		if !strings.Contains(authURL, "code_challenge_method=S256") {
			t.Error("expected URL to contain code_challenge_method parameter")
		}
	})

	t.Run("minimal params", func(t *testing.T) {
		authURL, err := BuildAuthURL("https://openrouter.ai/auth", AuthURLParams{
			CallbackURL: "https://myapp.com/callback",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(authURL, "callback_url=") {
			t.Error("expected URL to contain callback_url parameter")
		}
		if strings.Contains(authURL, "code_challenge=") {
			t.Error("expected URL to not contain code_challenge parameter")
		}
		if strings.Contains(authURL, "code_challenge_method=") {
			t.Error("expected URL to not contain code_challenge_method parameter")
		}
	})

	t.Run("empty callback", func(t *testing.T) {
		_, err := BuildAuthURL("https://openrouter.ai/auth", AuthURLParams{})
		if err == nil {
			t.Fatal("expected error for empty callback_url")
		}
	})
}
