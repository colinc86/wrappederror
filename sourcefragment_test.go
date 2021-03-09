package wrappederror

import (
	"strings"
	"testing"
)

func TestNewCaller_Radius(t *testing.T) {
	c := newCaller(1, true, 3)
	n := len(strings.Split(strings.TrimSpace(c.Fragment.Source), "\n"))
	if n != 7 {
		t.Errorf("Expected 7 lines, but found %d.\n", n)
	}

	c = newCaller(1, true, 1)
	n = len(strings.Split(strings.TrimSpace(c.Fragment.Source), "\n"))
	if n != 2 {
		t.Errorf("Expected 2 lines, but found %d.\n", n)
	}
}

func TestGetSourceFails(t *testing.T) {
	_, err := newSourceFragment("/something/that/does/not/exist", 0, 1)
	if err == nil {
		t.Error("Expected error.")
	}
}
