package wrappederror

import "testing"

func TestGlobalsConfig(t *testing.T) {
	// Sanity check
	if Config() == nil {
		t.Error("Unexpected nil config.")
	}
}
