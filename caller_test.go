package wrappederror

import (
	"testing"
)

func TestCurrentCaller(t *testing.T) {
	c := currentCaller(1)

	if c.FunctionName != "github.com/colinc86/wrappederror.TestCurrentCaller" {
		t.Errorf("Incorrect function name: %s\n", c.FunctionName)
	}

	if c.FileName != "wcaller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.FileName)
	}

	if c.LineNumber != 8 {
		t.Errorf("Incorrect line number: %d\n", c.LineNumber)
	}
}

func TestCurrentCaller_Skip(t *testing.T) {
	c1 := currentCaller(1)
	c2 := currentCaller(2)

	if c1.FileName == c2.FileName {
		t.Errorf("Incorrect file names.")
	}

	if c1.FunctionName == c2.FunctionName {
		t.Errorf("Incorrect function names.")
	}

	if c1.LineNumber == c2.LineNumber {
		t.Errorf("Incorrect line numbers.")
	}
}

func TestCallerStack(t *testing.T) {
	c := currentCaller(1)
	if c.Stack() == "" {
		t.Error("Expected a stack trace.")
	}
}

func TestCallerSource(t *testing.T) {
	c := currentCaller(1)
	if c.Source() == "" {
		t.Error("Expected a source trace.")
	}
}
