package wrappederror

import (
	"testing"
)

func TestCurrentCaller(t *testing.T) {
	c := currentCaller(1)

	if c.functionName != "github.com/colinc86/wrappederror.TestCurrentCaller" {
		t.Errorf("Incorrect function name: %s\n", c.functionName)
	}

	if c.fileName != "wcaller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.fileName)
	}

	if c.lineNumber != 8 {
		t.Errorf("Incorrect line number: %d\n", c.lineNumber)
	}
}

func TestCurrentCaller_Skip(t *testing.T) {
	c1 := currentCaller(1)
	c2 := currentCaller(2)

	if c1.fileName == c2.fileName {
		t.Errorf("Incorrect file names.")
	}

	if c1.functionName == c2.functionName {
		t.Errorf("Incorrect function names.")
	}

	if c1.lineNumber == c2.lineNumber {
		t.Errorf("Incorrect line numbers.")
	}
}

func TestCallerStack(t *testing.T) {
	c := currentCaller(1)
	if c.Stack() == "" {
		t.Error("Expected a stack trace.")
	}
}
