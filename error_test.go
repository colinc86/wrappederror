package wrappederror

import (
	"errors"
	"strings"
	"testing"
)

func TestNewError_1(t *testing.T) {
	outerErrorMessage := "outer error"
	we := New(nil, outerErrorMessage)

	if we.Error() != outerErrorMessage {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", outerErrorMessage, we.Error())
	}
}

func TestNewError_2(t *testing.T) {
	innerErrorMessage := "inner error"
	e := errors.New(innerErrorMessage)

	outerErrorMessage := "outer error"
	we := New(e, outerErrorMessage)

	if we.Error() != outerErrorMessage {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", outerErrorMessage, we.Error())
		return
	}
}

func TestNewError_3(t *testing.T) {
	innerErrorMessage := "inner error"
	e := errors.New(innerErrorMessage)

	middleErrorMessage := "middle error"
	wem := New(e, middleErrorMessage)

	outerErrorMessage := "outer error"
	weo := New(wem, outerErrorMessage)

	composite := outerErrorMessage + ": " + middleErrorMessage + ": " + innerErrorMessage
	if weo.String() != composite {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", composite, weo.String())
		return
	}
}

func TestDepth_0(t *testing.T) {
	we := New(nil, "single error")
	if we.Depth() != 0 {
		t.Errorf("Expected depth 0 but received %d.\n", we.Depth())
	}
}

func TestDepth_1(t *testing.T) {
	we := New(errors.New("error 0"), "error 1")
	if we.Depth() != 1 {
		t.Errorf("Expected depth 1 but received %d.\n", we.Depth())
	}
}

func TestDepth_2(t *testing.T) {
	e0 := errors.New("error 0")
	e1 := New(e0, "error 1")
	e2 := New(e1, "error 2")

	if e2.Depth() != 2 {
		t.Errorf("Expected depth 2 but received %d.\n", e2.Depth())
	}
}

func TestDepth_3(t *testing.T) {
	e0 := errors.New("error 0")
	e1 := New(e0, "error 1")
	e2 := New(e1, "error 2")
	e3 := New(e2, "error 3")

	if e3.Depth() != 3 {
		t.Errorf("Expected depth 2 but received %d.\n", e3.Depth())
	}
}

func TestTrace(t *testing.T) {
	// Quick check for sanity
	e0 := errors.New("error 0")
	e1 := New(e0, "error 1")
	e2 := New(e1, "error 2")
	e3 := New(e2, "error 3")

	tr := e3.Trace()
	nc := strings.Count(tr, "\n")
	if nc != 3 {
		t.Errorf("Expected 3 newlines but found %d.\n", nc)
	}
}

func TestCaller(t *testing.T) {
	we := New(nil, "test")
	if we.Caller().File() != "werror_test.go" ||
		we.Caller().Function() != "github.com/colinc86/wrappederror.TestCaller" ||
		we.Caller().Line() != 98 {
		t.Errorf("Incorrect caller: %s\n", we.(*wError).caller)
	}
}

func TestUnwrap(t *testing.T) {
	e := errors.New("inner error")
	we := New(e, "outer error")

	if we.Unwrap() != e {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", e, we.Unwrap())
	}
}

func TestErrorMarshalText(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := New(e1, "error 2")
	e3 := New(e2, "error 3")
	e4 := New(e3, "error 4")

	d, err := e4.MarshalText()
	if err != nil {
		t.Errorf("Error marshaling text: %s\n", err)
	}

	we := &wError{}
	if err = we.UnmarshalText(d); err != nil {
		t.Errorf("Error unmarshaling text: %s\n", err)
	}

	if string(d) != we.String() {
		t.Error("Expected unmarshaled error.")
	}
}

func TestErrorMarshalBinary(t *testing.T) {
	e := New(nil, "test")

	d, err := e.MarshalBinary()
	if err != nil {
		t.Errorf("Error marshaling binary: %s\n", err)
	}

	we := &wError{}
	if err = we.UnmarshalBinary(d); err != nil {
		t.Errorf("Error unmarshaling binary: %s\n", err)
	}
}