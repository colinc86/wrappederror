// Package wrappederror contains the snarky WrappedError type.
package wrappederror

import (
	"strings"
)

// WrappedError types wrap an error. Along with your wrapped error you can give
// the recipient of said error some context by including a message.
type WrappedError struct {

	// The error message that gives context to err for the programmers that need a
	// little bit of help.
	message string

	// The inner error that this wrapped error wraps.
	err error
}

// Initializers

// New creates and returns a new wrapped error... like the name suggests.
func New(message string, err error) WrappedError {
	return WrappedError{
		message: message,
		err:     err,
	}
}

// Error interface methods

func (e WrappedError) Error() string {
	if e.err == nil {
		return e.message
	}
	return e.message + ": " + e.err.Error()
}

func (e WrappedError) Unwrap() error {
	return e.err
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
	return e.MarshalBinary()
}

// UnmarshalText unmarshals text in to a wrapped error.
func (e *WrappedError) UnmarshalText(b []byte) error {
	return e.UnmarshalBinary(b)
}

// BinaryMarshaler and BinaryUnmarshaler interface methods

// MarshalBinary marshals the wrapped error.
func (e WrappedError) MarshalBinary() ([]byte, error) {
	return []byte(e.Error()), nil
}

// UnmarshalBinary unmarshals in to a wrapped error. Since the wrapped error
// doesn't know what you want from it, all errors that the wrapped errors
// wrapped are now wrapped errors themselves. Say that 5 times fast.
func (e *WrappedError) UnmarshalBinary(d []byte) error {
	c := strings.Split(strings.TrimSpace(string(d)), ":")
	l := len(c)

	if l == 0 {
		e.message = ""
		e.err = nil
	} else if l == 1 {
		e.message = c[0]
		e.err = nil
	} else if l > 1 {
		e.message = c[0]

		we := new(WrappedError)
		_ = we.UnmarshalJSON([]byte(strings.Join(c[1:], ":")))
		e.err = we
	}

	return nil
}
