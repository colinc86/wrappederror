package wrappederror

import (
	"strings"
	"testing"
)

func TestNewCaller_Skip(t *testing.T) {
	c1 := newCaller(1, 2)
	c2 := newCaller(2, 2)

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
	c := newCaller(1, 3)
	n := len(strings.Split(strings.TrimSpace(c.SourceFragment), "\n"))
	if n != 8 {
		t.Errorf("Expected 9 lines, but found %d.\n", n)
	}

	c = newCaller(1, 1)
	n = len(strings.Split(strings.TrimSpace(c.SourceFragment), "\n"))
	if n != 4 {
		t.Errorf("Expected 5 lines, but found %d.\n", n)
	}
}

func TestCallerFile(t *testing.T) {
	c := newCaller(1, 2)
	if c.File != "caller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.File)
	}
}

func TestCallerFunction(t *testing.T) {
	c := newCaller(1, 2)
	if c.Function != "github.com/colinc86/wrappederror.TestCallerFunction" {
		t.Errorf("Incorrect function name: %s\n", c.Function)
	}
}

func TestCallerLine(t *testing.T) {
	c := newCaller(1, 2)
	if c.Line != 54 {
		t.Errorf("Incorrect line number: %d\n", c.Line)
	}
}

func TestCallerStack(t *testing.T) {
	c := newCaller(1, 2)
	if c.StackTrace == "" {
		t.Error("Expected a stack trace.")
	}
}

func TestCallerSource(t *testing.T) {
	c := newCaller(1, 2)
	if c.SourceFragment == "" {
		t.Error("Expected a source trace.")
	}
}

func TestGetSourceFails(t *testing.T) {
	_, err := getSource("/something/that/does/not/exist", 0, 1)
	if err == nil {
		t.Error("Expected error.")
	}
}
