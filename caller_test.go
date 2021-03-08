package wrappederror

import (
	"testing"
)

func TestCurrentCaller(t *testing.T) {
	c := newCaller(1, 2)

	if c.Function != "github.com/colinc86/wrappederror.TestCurrentCaller" {
		t.Errorf("Incorrect function name: %s\n", c.Function)
	}

	if c.File != "caller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.File)
	}

	if c.Line != 8 {
		t.Errorf("Incorrect line number: %d\n", c.Line)
	}
}

func TestCurrentCaller_Skip(t *testing.T) {
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
