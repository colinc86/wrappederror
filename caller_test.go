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
}

func TestCallerMarshalBinary(t *testing.T) {
	c1 := currentCaller(1)
	b, err := c1.MarshalBinary()
	if err != nil {
		t.Errorf("Unable to marshal caller in to binary: %s\n", err)
	}

	c2 := &caller{}
	if err = c2.UnmarshalBinary(b); err != nil {
		t.Errorf("Unable to unmarshal caller from binary: %s\n", err)
	}

	if c1.fileName != c2.fileName ||
		c1.functionName != c2.functionName ||
		c1.lineNumber != c2.lineNumber {
		t.Errorf("Expected \"%s\" and \"%s\" to be equal.\n", c1, c2)
	}
}
