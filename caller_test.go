package wrappederror

import (
	"testing"
)

func TestNewCaller(t *testing.T) {
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

func TestNewCallerFailure_1(t *testing.T) {
	c := newCaller(10, false, 2)
	if c.File != callerFileNameUnknown ||
		c.Function != callerFunctionNameUnknown ||
		c.Line != callerLineNumberUnknown ||
		c.Fragment != nil {
		t.Errorf("Unknown caller information %s.\n", c)
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
	if c.Line != 39 {
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
