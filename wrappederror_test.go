package wrappederror

import (
	"errors"
	"testing"
)

func TestNewWrappedError_1(t *testing.T) {
	outerErrorMessage := "outer error"
	we := New(outerErrorMessage, nil)

	if we.Error() != outerErrorMessage {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", outerErrorMessage, we.Error())
	}
}

func TestUnwrap(t *testing.T) {
	e := errors.New("inner error")
	we := New("outer error", e)

	if we.Unwrap() != e {
		t.Errorf("Expected \"%s\" but received \"%s\"\n", e, we.Unwrap())
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

func TestWrappedErrorMarshalBinary(t *testing.T) {
	e1 := errors.New("error 1")
	e2 := New("error 2", e1)
	e3 := New("error 3", e2)
	e4 := New("error 4", e3)

	d, err := e4.MarshalBinary()
	if err != nil {
		t.Errorf("Error marshaling binary: %s\n", err)
	}

	we := &WrappedError{}
	if err = we.UnmarshalBinary(d); err != nil {
		t.Errorf("Error unmarshaling binary: %s\n", err)
	}

	if e4.Error() != we.Error() {
		t.Error("Expected unmarshaled error.")
	}
}
