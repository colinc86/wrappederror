package wrappederror

import (
	"testing"
)

var testCallers struct {
	c0 *Caller
	c1 *Caller
	c2 *Caller
}

func setupCallerTests() {
	testCallers.c0 = newCaller(1, false, 2)
	testCallers.c1 = newCaller(2, false, 2)
	testCallers.c2 = newCaller(1, true, 2)
}

// Tests

func TestNewCaller(t *testing.T) {
	if testCallers.c0.File == testCallers.c1.File {
		t.Errorf("Incorrect file names.")
	}

	if testCallers.c0.Function == testCallers.c1.Function {
		t.Errorf("Incorrect function names.")
	}

	if testCallers.c0.Line == testCallers.c1.Line {
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
	if testCallers.c0.File != "caller_test.go" {
		t.Errorf("Incorrect file name: %s\n", testCallers.c0.File)
	}
}

func TestCallerFunction(t *testing.T) {
	if testCallers.c0.Function != "github.com/colinc86/wrappederror.setupCallerTests" {
		t.Errorf("Incorrect function name: %s\n", testCallers.c0.Function)
	}
}

func TestCallerLine(t *testing.T) {
	if testCallers.c0.Line != 14 {
		t.Errorf("Incorrect line number: %d\n", testCallers.c0.Line)
	}
}

func TestCallerStack(t *testing.T) {
	if testCallers.c0.StackTrace == "" {
		t.Error("Expected a stack trace.")
	}
}

func TestCallerSource(t *testing.T) {
	if testCallers.c2.Fragment.Source == "" {
		t.Error("Expected a source trace.")
	}
}
