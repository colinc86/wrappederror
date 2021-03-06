// Package wrappederror contains the Error type.
package wrappederror

import (
	"encoding"
	"errors"
	"fmt"
	"strings"
)

// Error types wrap an error.
type Error interface {
	error
	encoding.TextMarshaler
	encoding.TextUnmarshaler

	// Caller returns the error's caller.
	//
	// Use this value to examine the error's call information such as file name,
	// function name and line number.
	Caller() Caller

	// The error's context.
	//
	// When an error is wrapped, it is given context. An error's context can be a
	// string description or any other type of information that is pertinent to
	// the error being wrapped.
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

	// As finds the first error in the error chain that matches target, and if so,
	// sets target to that error value and returns true. Otherwise, it returns
	// false.
	As(target interface{}) bool

	// Is reports whether any error in the error chain matches target.
	//
	// The chain consists of the receiver followed by the sequence of errors
	// obtained by repeatedly calling Unwrap.
	//
	// The receiver considered to match a target if it is equal to that target or
	// if it implements a method Is(error) bool such that Is(target) returns true.
	Is(target error) bool
}

// Error implementation

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

// Error interface methods

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
	var s string
	e.Walk(func(err error) bool {
		if we, ok := err.(wError); ok {
			s += fmt.Sprintf("%+v", we.Context())
		} else if we, ok := err.(*wError); ok {
			s += fmt.Sprintf("%+v", we.Context())
		} else {
			s += err.Error()
		}

		if errors.Unwrap(err) != nil {
			s += ": "
		}

		return true
	})
	return s
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
