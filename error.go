// Package wrappederror contains the Error type.
package wrappederror

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/colinc86/coding"
)

// Error types wrap an error.
type Error interface {
	error
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Caller returns the error's caller.
	Caller() Caller

	// The error's context.
	Context() interface{}

	// Walk calls the step function for each error in the error chain and
	// continues until either the last error is reached, or until the step
	// function returns false.
	Walk(step func(err error) bool)

	// Depth returns the number of nested errors in the receiver.
	Depth() uint

	// Trace returns a prettified string representation of the wrapped error.
	Trace() string

	// Unwrap unwraps the wrapped error.
	Unwrap() error

	// As finds the first error in err's chain that matches target, and if so,
	// sets target to that error value and returns true. Otherwise, it returns
	// false.
	As(target interface{}) bool

	// Is reports whether any error in err's chain matches target.
	//
	// The chain consists of err itself followed by the sequence of errors
	// obtained by repeatedly calling Unwrap.
	//
	// An error is considered to match a target if it is equal to that target or
	// if it implements a method Is(error) bool such that Is(target) returns true.
	Is(target error) bool
}

// wError types wrap an error.
type wError struct {
	context interface{}

	// The inner error that this wrapped error wraps.
	inner error

	// The caller that invoked the `New` function.
	caller *caller
}

// Initializers

// New creates and returns a new wrapped error.
func New(err error, context interface{}) Error {
	return &wError{
		context: context,
		inner:   err,
		caller:  currentCaller(2),
	}
}

// Exported methods

func (e wError) Context() interface{} {
	return e.context
}

func (e wError) Caller() Caller {
	return e.caller
}

func (e wError) Walk(step func(err error) bool) {
	var ce error = e
	for {
		if !step(ce) {
			break
		}

		ue := errors.Unwrap(ce)
		if ue == nil {
			break
		}

		ce = ue
	}
}

func (e wError) Depth() uint {
	var d uint
	e.Walk(func(err error) bool {
		if errors.Unwrap(err) == nil {
			return false
		}
		d++
		return true
	})
	return d
}

func (e wError) Trace() string {
	d := e.Depth()
	if d == 0 {
		return fmt.Sprintf("%s %s", e.caller, e.Error())
	}

	var msg string
	e.Walk(func(err error) bool {
		u := errors.Unwrap(err) != nil
		var p string
		if msg == "" {
			p = "┌"
		} else if u {
			p = "└"
		} else {
			p = "├"
		}

		var em string
		if we, ok := err.(Error); ok {
			em = fmt.Sprintf("%s %d: %s %s", p, d, we.Caller(), we.Error())
		} else {
			em = fmt.Sprintf("%s %d: %s", p, d, err.Error())
		}

		msg += em
		if u {
			msg += "\n"
		}

		return true
	})

	return msg
}

// Error interface methods

func (e wError) Error() string {
	if e.context == nil {
		return ""
	}
	return fmt.Sprintf("%+v", e.context)
}

func (e wError) Unwrap() error {
	return e.inner
}

func (e wError) As(target interface{}) bool {
	as := false
	e.Walk(func(err error) bool {
		if err == target {
			target = err
			as = true
			return false
		}
		return true
	})
	return as
}

func (e wError) Is(target error) bool {
	is := false
	e.Walk(func(err error) bool {
		if err == target {
			is = true
			return false
		}
		return true
	})
	return is
}

// Stringer interface methods

func (e wError) String() string {
	var s string
	e.Walk(func(err error) bool {
		s += err.Error()

		if errors.Unwrap(err) != nil {
			s += ": "
		}

		return true
	})
	return s
}

// TextMarshaler and TextUnmarshaler interface methods

// MarshalText marshals the wrapped error in to text, but not JSON or binary.
func (e wError) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText unmarshals in to a wrapped error. Since the wrapped error
// doesn't know what you want from it, all errors that the wrapped errors
// wrapped are now wrapped errors themselves. Say that 5 times fast.
func (e *wError) UnmarshalText(b []byte) error {
	c := strings.Split(strings.TrimSpace(string(b)), ":")
	l := len(c)

	if l == 0 {
		e.context = ""
		e.inner = nil
	} else if l == 1 {
		e.context = c[0]
		e.inner = nil
	} else if l > 1 {
		e.context = c[0]

		we := new(wError)
		_ = we.UnmarshalText([]byte(strings.Join(c[1:], ":")))
		e.inner = we
	}

	return nil
}

// BinaryMarshaler and BinaryUnmarshaler interface methods

// MarshalBinary marshals the wrapped error in to binary.
func (e wError) MarshalBinary() ([]byte, error) {
	en := coding.NewEncoder()
	e.Walk(func(err error) bool {
		if we, ok := err.(wError); ok {
			// Is this error a Error?
			en.EncodeBool(true)

			// Attempt to marshal the context in to JSON data
			if jsonData, jsonErr := json.Marshal(we.Context()); jsonErr == nil {
				en.EncodeData(jsonData)
			} else {
				en.EncodeData(nil)
			}

			// Attempt to marshal the caller in to binary data
			if callerData, callerErr := we.Caller().MarshalBinary(); callerErr == nil {
				en.EncodeData(callerData)
			} else {
				en.EncodeData(nil)
			}
		} else {
			// TODO: Figure out how to marshal/unmarshal any error type
			en.EncodeBool(false)
			en.EncodeString(err.Error())
		}
		return true
	})

	return en.Compress()
}

// UnmarshalBinary unmarshals the wrapped error from binary.
func (e *wError) UnmarshalBinary(d []byte) error {
	de := coding.NewDecoder(d)
	if err := de.Decompress(); err != nil {
		return err
	}

	if err := de.Validate(); err != nil {
		return err
	}

	var errs []error
	for {
		isWError, err := de.DecodeBool()
		if err != nil {
			if err == coding.ErrEOB {
				break
			}
			return err
		}

		if isWError {
			ctxData, ctxErr := de.DecodeData()
			if ctxErr != nil {
				return ctxErr
			}

			ctx := new(interface{})
			jsonErr := json.Unmarshal(ctxData, ctx)
			if jsonErr != nil {
				return jsonErr
			}

			callerData, callerErr := de.DecodeData()
			if callerErr != nil {
				return callerErr
			}

			caller := new(caller)
			calUnmErr := caller.UnmarshalBinary(callerData)
			if calUnmErr != nil {
				return calUnmErr
			}

			errs = append(errs, &wError{ctx, nil, caller})
		} else {
			errStr, err := de.DecodeString()
			if err != nil {
				return err
			}

			errs = append(errs, errors.New(errStr))
		}
	}

	topError := new(wError)
	currentError := topError
	for _, err := range errs {
		currentError.inner = err

		if we, ok := err.(*wError); ok {
			currentError = we
		}
	}

	if we, ok := topError.inner.(*wError); ok {
		e.context = we.context
		e.caller = we.caller
		e.inner = we.inner
	}
	return nil
}
