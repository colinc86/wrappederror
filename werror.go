package wrappederror

import (
	"fmt"
	"strings"
)

// wError types wrap an error.
type wError struct {
	message string

	// The inner error that this wrapped error wraps.
	inner error

	// The caller that invoked the `New` function.
	caller *caller
}

// Initializers

// New creates and returns a new wrapped error.
func New(message string, err error) WrappedError {
	return &wError{
		message: message,
		inner:   err,
		caller:  currentCaller(2),
	}
}

// Exported methods

func (e wError) Depth() uint {
	if e.inner == nil {
		return 0
	} else if we, ok := e.inner.(WrappedError); ok {
		return we.Depth() + 1
	}
	return 1
}

func (e wError) Trace() string {
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
			if we.Unwrap() == nil {
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

func (e wError) File() string {
	return e.caller.fileName
}

func (e wError) Function() string {
	return e.caller.functionName
}

func (e wError) Line() int {
	return e.caller.lineNumber
}

// Error interface methods

func (e wError) Error() string {
	if e.inner == nil {
		return e.message
	}
	return e.message + ": " + e.inner.Error()
}

func (e wError) Unwrap() error {
	return e.inner
}

// String interface methods

func (e wError) String() string {
	return e.Error()
}

// TextMarshaler and TextUnmarshaler interface methods

// MarshalText marshals the wrapped error in to text, but not JSON or binary.
func (e wError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// UnmarshalText unmarshals in to a wrapped error. Since the wrapped error
// doesn't know what you want from it, all errors that the wrapped errors
// wrapped are now wrapped errors themselves. Say that 5 times fast.
func (e *wError) UnmarshalText(b []byte) error {
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

		we := new(wError)
		_ = we.UnmarshalText([]byte(strings.Join(c[1:], ":")))
		e.inner = we
	}

	return nil
}

// BinaryMarshaler and BinaryUnmarshaler interface methods

// MarshalBinary marshals the wrapped error in to binary.
func (e wError) MarshalBinary() ([]byte, error) {
	en := newEncoder()
	en.encodeWrappedError(e)
	en.calculateCRC()

	if err := en.compress(); err != nil {
		return nil, err
	}

	return en.data, nil
}

// UnmarshalBinary unmarshals the wrapped error from binary.
func (e *wError) UnmarshalBinary(d []byte) error {
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

	topError := new(wError)
	currentError := topError
	for _, err := range errors {
		currentError.inner = err

		if we, ok := err.(*wError); ok {
			currentError = we
		}
	}

	if we, ok := topError.inner.(*wError); ok {
		e.message = we.message
		e.caller = we.caller
		e.inner = we.inner
	}
	return nil
}
