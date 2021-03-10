package wrappederror

import (
	"testing"
)

func TestNewProcess(t *testing.T) {
	we := New(nil, "test")

	if we.Process.Routines < 1 {
		t.Errorf("Unreasonable routines: %d\n", we.Process.Routines)
	}

	if we.Process.CPUs < 1 {
		t.Errorf("Unreasonable cpus: %d\n", we.Process.CPUs)
	}

	if we.Process.CGO < 0 {
		t.Errorf("Unreasonable cgos: %d\n", we.Process.CGO)
	}

	if we.Process.Memory == nil {
		t.Error("Expected memory statistics.")
	}
}

func TestProcessIgnoreBreakpoints(t *testing.T) {
	packageState.config.SetIgnoreBreakpoints(true)
	we := New(nil, "test")
	we.Process.Break()
}

func TestProcessString(t *testing.T) {
	// Sanity check
	we := New(nil, "test")
	if len(we.Process.String()) == 0 {
		t.Errorf("Unexpected string length %d.\n", len(we.Process.String()))
	}
}
