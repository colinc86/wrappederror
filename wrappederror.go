// Package wrappederror contains the snarky WrappedError type.
package wrappederror

import (
	"fmt"
	"strings"
)

// WrappedError types wrap an error. Along with your wrapped error you can give
// the recipient of said error some context by including a message.
type WrappedError struct {

	// The error message that gives context to err for the programmers that need a
	// little bit of help.
	message string

	// The inner error that this wrapped error wraps.
	inner error

	// The caller that invoked the `New` function.
	caller caller
}

// Initializers

// New creates and returns a new wrapped error... like the name suggests.
func New(message string, err error) WrappedError {
	return WrappedError{
		message: message,
		inner:   err,
		caller:  currentCaller(2),
	}
}

// Exported methods

// Depth returns the number of nested errors in the receiver.
func (e WrappedError) Depth() uint {
	if e.inner == nil {
		return 0
	} else if we, ok := e.inner.(WrappedError); ok {
		return we.Depth() + 1
	}
	return 1
}

// Trace returns a prettified string representation of the wrapped error.
func (e WrappedError) Trace() string {
	msg := fmt.Sprintf("%s %s", e.caller, e.message)
	if e.Depth() == 0 {
		return msg
	}

	// If the current caller with depth 1 isn't the same as the caller with depth
	// 2, then we know we are the "top-most" error in the trace.
	p1 := ""
	c1 := currentCaller(1)
	c2 := currentCaller(2)
	if c1.functionName != c2.functionName {
		p1 = "┌ "
	}

	// Add to our message
	msg = fmt.Sprintf("%s%d: %s", p1, e.Depth(), msg)

	// Do some recursive stuff
	if e.inner != nil {
		if we, ok := e.inner.(WrappedError); ok {
			var p2 string
			if we.inner == nil {
				p2 = "└"
			} else {
				p2 = "├"
			}

			msg += "\n" + p2 + " " + we.Trace()
		} else {
			msg += "\n└ 0: " + e.inner.Error()
		}
	}

	return msg
}

// Error interface methods

func (e WrappedError) Error() string {
	if e.inner == nil {
		return e.message
	}
	return e.message + ": " + e.inner.Error()
}

func (e WrappedError) Unwrap() error {
	return e.inner
}

// JSON Marshaler and Unmarshaler interface methods

// MarshalJSON marshals the wrapped error in to JSON, but not text or binary.
func (e WrappedError) MarshalJSON() ([]byte, error) {
	return e.MarshalText()
}

// UnmarshalJSON unmarshals JSON in to a wrapped error.
func (e *WrappedError) UnmarshalJSON(b []byte) error {
	return e.UnmarshalText(b)
}

// TextMarshaler and TextUnmarshaler interface methods

// MarshalText marshals the wrapped error in to text, but not JSON or binary.
func (e WrappedError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// UnmarshalText unmarshals in to a wrapped error. Since the wrapped error
// doesn't know what you want from it, all errors that the wrapped errors
// wrapped are now wrapped errors themselves. Say that 5 times fast.
func (e *WrappedError) UnmarshalText(b []byte) error {
	c := strings.Split(strings.TrimSpace(string(b)), ":")
	l := len(c)

	if l == 0 {
		e.message = ""
		e.inner = nil
	} else if l == 1 {
		e.message = c[0]
		e.inner = nil
	} else if l > 1 {
		e.message = c[0]

		we := new(WrappedError)
		_ = we.UnmarshalText([]byte(strings.Join(c[1:], ":")))
		e.inner = we
	}

	return nil
}

// BinaryMarshaler and BinaryUnmarshaler interface methods

// MarshalBinary marshals the wrapped error in to binary.
func (e WrappedError) MarshalBinary() ([]byte, error) {
	en := newEncoder()
	en.encodeWrappedError(e)
	en.calculateCRC()

	if err := en.compress(); err != nil {
		return nil, err
	}

	return en.data, nil
}

// UnmarshalBinary unmarshals the wrapped error from binary.
func (e *WrappedError) UnmarshalBinary(d []byte) error {
	de := newDecoder(d)
	if err := de.decompress(); err != nil {
		return err
	}

	if !de.validate() {
		return ErrCRC
	}

	var errors []error
	for {
		decodedError, err := de.decodeError()
		if err != nil {
			return err
		}

		if decodedError == nil {
			break
		}

		errors = append(errors, decodedError)
	}

	topError := &WrappedError{}
	currentError := topError
	for _, err := range errors {
		currentError.inner = err

		if we, ok := err.(*WrappedError); ok {
			currentError = we
		}
	}

	if we, ok := topError.inner.(*WrappedError); ok {
		e.message = we.message
		e.caller = we.caller
		e.inner = we.inner
	}
	return nil
}
