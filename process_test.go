package wrappederror

import (
	"testing"
)

func TestCurrentProcess(t *testing.T) {
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
