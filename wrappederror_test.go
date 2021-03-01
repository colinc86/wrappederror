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

func TestMarshalJSON(t *testing.T) {
	we := New("outer error", errors.New("inner error"))
	d, err := we.MarshalJSON()
	if err != nil {
		t.Errorf("Unable to marshal wrapped error: %s\n", err)
		return
	}

	if len(d) != 24 {
		t.Errorf("Expected byte length 24 but received %d\n", len(d))
		return
	}
}

func TestUnmarshalJSON_1(t *testing.T) {
	outerErrorText := "outer error"
	innerErrorText := "inner error"
	text := outerErrorText + ": " + innerErrorText
	d := []byte(text)

	we := new(WrappedError)
	if err := we.UnmarshalJSON(d); err != nil {
		t.Errorf("Unable to unmarshal wrapped error: %s\n", err)
		return
	}

	if we.message != outerErrorText {
		t.Errorf("Expected outer error \"%s\" but received \"%s\"\n", outerErrorText, we.message)
		return
	}

	if we.err.Error() != innerErrorText {
		t.Errorf("Expected inner error \"%s\" but received \"%s\"\n", innerErrorText, we.err.Error())
		return
	}
}

func TestUnmarshalJSON_2(t *testing.T) {
	var d []byte

	we := new(WrappedError)
	if err := we.UnmarshalJSON(d); err != nil {
		t.Errorf("Unable to unmarshal wrapped error: %s\n", err)
		return
	}

	if we.message != "" {
		t.Errorf("Expected outer error \"\" but received \"%s\"\n", we.message)
		return
	}

	if we.err != nil {
		t.Errorf("Expected inner error to be nil but received \"%s\"\n", we.err.Error())
		return
	}
}
