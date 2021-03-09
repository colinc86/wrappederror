package wrappederror

import (
	"strings"
	"testing"
)

func TestNewCaller_Skip(t *testing.T) {
	c1 := newCaller(1, false, 2)
	c2 := newCaller(2, false, 2)

	if c1.File == c2.File {
		t.Errorf("Incorrect file names.")
	}

	if c1.Function == c2.Function {
		t.Errorf("Incorrect function names.")
	}

	if c1.Line == c2.Line {
		t.Errorf("Incorrect line numbers.")
	}
}

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

func TestCallerFile(t *testing.T) {
	c := newCaller(1, false, 2)
	if c.File != "caller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.File)
	}
}

func TestCallerFunction(t *testing.T) {
	c := newCaller(1, false, 2)
	if c.Function != "github.com/colinc86/wrappederror.TestCallerFunction" {
		t.Errorf("Incorrect function name: %s\n", c.Function)
	}
}

func TestCallerLine(t *testing.T) {
	c := newCaller(1, false, 2)
	if c.Line != 54 {
		t.Errorf("Incorrect line number: %d\n", c.Line)
	}
}

func TestCallerStack(t *testing.T) {
	c := newCaller(1, false, 2)
	if c.StackTrace == "" {
		t.Error("Expected a stack trace.")
	}
}

func TestCallerSource(t *testing.T) {
	c := newCaller(1, true, 2)
	if c.Fragment.Source == "" {
		t.Error("Expected a source trace.")
	}
}

func TestGetSourceFails(t *testing.T) {
	_, err := newSourceFragment("/something/that/does/not/exist", 0, 1)
	if err == nil {
		t.Error("Expected error.")
	}
}
