package wrappederror

import (
	"errors"
	"strings"
	"testing"
)

func TestNewWrappedError_1(t *testing.T) {
	outerErrorMessage := "outer error"
	we := New(outerErrorMessage, nil)

	if we.Error() != outerErrorMessage {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", outerErrorMessage, we.Error())
	}
}

func TestNewWrappedError_2(t *testing.T) {
	innerErrorMessage := "inner error"
	e := errors.New(innerErrorMessage)

	outerErrorMessage := "outer error"
	we := New(outerErrorMessage, e)

	composite := outerErrorMessage + ": " + innerErrorMessage
	if we.Error() != composite {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", composite, we.Error())
		return
	}
}

func TestNewWrappedError_3(t *testing.T) {
	innerErrorMessage := "inner error"
	e := errors.New(innerErrorMessage)

	middleErrorMessage := "middle error"
	wem := New(middleErrorMessage, e)

	outerErrorMessage := "outer error"
	weo := New(outerErrorMessage, wem)

	composite := outerErrorMessage + ": " + middleErrorMessage + ": " + innerErrorMessage
	if weo.Error() != composite {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", composite, weo.Error())
		return
	}
}

func TestDepth_0(t *testing.T) {
	we := New("single error", nil)
	if we.Depth() != 0 {
		t.Errorf("Expected depth 0 but received %d.\n", we.Depth())
	}
}

func TestDepth_1(t *testing.T) {
	we := New("error 1", errors.New("error 0"))
	if we.Depth() != 1 {
		t.Errorf("Expected depth 1 but received %d.\n", we.Depth())
	}
}

func TestDepth_2(t *testing.T) {
	e0 := errors.New("error 0")
	e1 := New("error 1", e0)
	e2 := New("error 2", e1)

	if e2.Depth() != 2 {
		t.Errorf("Expected depth 2 but received %d.\n", e2.Depth())
	}
}

func TestDepth_3(t *testing.T) {
	e0 := errors.New("error 0")
	e1 := New("error 1", e0)
	e2 := New("error 2", e1)
	e3 := New("error 3", e2)

	if e3.Depth() != 3 {
		t.Errorf("Expected depth 2 but received %d.\n", e3.Depth())
	}
}

func TestTrace(t *testing.T) {
	// Quick check for sanity
	e0 := errors.New("error 0")
	e1 := New("error 1", e0)
	e2 := New("error 2", e1)
	e3 := New("error 3", e2)

	tr := e3.Trace()
	nc := strings.Count(tr, "\n")
	if nc != 3 {
		t.Errorf("Expected 3 newlines but found %d.\n", nc)
	}
}

func TestCaller(t *testing.T) {
	we := New("test", nil)
	if we.File() != "werror_test.go" ||
		we.Function() != "github.com/colinc86/wrappederror.TestCaller" ||
		we.Line() != 101 {
		t.Errorf("Incorrect caller: %s\n", we.(*wError).caller)
	}
}

func TestUnwrap(t *testing.T) {
	e := errors.New("inner error")
	we := New("outer error", e)

	if we.Unwrap() != e {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", e, we.Unwrap())
	}
}

func TestString(t *testing.T) {
	e0 := errors.New("error A")
	e1 := New("error B", e0)
	if e1.Error() != e1.String() {
		t.Errorf("Expected equal strings %s != %s\n", e1.Error(), e1.String())
	}
}

func TestWrappedErrorMarshalText(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := New("error 2", e1)
	e3 := New("error 3", e2)
	e4 := New("error 4", e3)

	d, err := e4.MarshalText()
	if err != nil {
		t.Errorf("Error marshaling text: %s\n", err)
	}

	we := &wError{}
	if err = we.UnmarshalText(d); err != nil {
		t.Errorf("Error unmarshaling text: %s\n", err)
	}

	if string(d) != we.Error() {
		t.Error("Expected unmarshaled error.")
	}
}

func TestWrappedErrorMarshalBinary(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := New("error 2", e1)
	e3 := New("error 3", e2)
	e4 := New("error 4", e3)

	d, err := e4.MarshalBinary()
	if err != nil {
		t.Errorf("Error marshaling binary: %s\n", err)
	}

	we := &wError{}
	if err = we.UnmarshalBinary(d); err != nil {
		t.Errorf("Error unmarshaling binary: %s\n", err)
	}

	if e4.Error() != we.Error() {
		t.Error("Expected unmarshaled error.")
	}
}
