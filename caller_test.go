package wrappederror

import (
	"testing"
)

func TestCurrentCaller(t *testing.T) {
	c := currentCaller(1)

	if c.functionName != "github.com/colinc86/wrappederror.TestCurrentCaller" {
		t.Errorf("Incorrect function name: %s\n", c.functionName)
	}

	if c.fileName != "caller_test.go" {
		t.Errorf("Incorrect file name: %s\n", c.fileName)
	}

	if c.lineNumber != 8 {
		t.Errorf("Incorrect line number: %d\n", c.lineNumber)
	}
}
