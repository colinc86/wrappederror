package wrappederror

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error types wrap an error and provide context, caller and process
// information.
type Error struct {
	error
	json.Marshaler

	// The error's context.
	//
	// When an error is wrapped, it is given context. An error's context can be a
	// string description or any other type of information that is pertinent to
	// the error being wrapped.
	Context interface{}

	// The error's caller.
	//
	// Use this value to examine the error's call information such as file name,
	// function name, line number, stack trace and source fragment.
	//
	// Call information is captured by default. To not capture this information,
	// use the SetCaptureCaller configuration method.
	//
	// If call information is not captured, then this property is nil.
	Caller *Caller

	// Process information at the time the error was created.
	//
	// Use this value to examine the current process's information such as number
	// of goroutines, available logical CPUs, memory statistics, and so forth.
	//
	// Process information is captured by default. To not capture this
	// information, use the SetCaptureProcess configuration method.
	//
	// If process information is not captured, then this property is nil.
	Process *Process

	// Metadata pertaining to the error.
	//
	// Use this value to examine general properties about the error such as its
	// index, timestamp, and how many other similar errors have been created.
	//
	// Metadata is always captured, but some of its properties are configurable.
	Metadata *Metadata

	// The inner error that this wrapped error wraps.
	inner error
}

// Initializers

// New creates and returns a new error.
func New(err error, ctx interface{}) *Error {
	var caller *Caller
	if packageState.config.CaptureCaller() {
		caller = newCaller(2, packageState.config.SourceFragmentRadius())
	}

	var process *Process
	if packageState.config.CaptureProcess() {
		process = newProcess()
	}

	return &Error{
		Context:  ctx,
		Caller:   caller,
		Process:  process,
		Metadata: newMetadata(err),
		inner:    err,
	}
}

// Exported methods

// Format returns a formatted string representation of the error using the error
// format string, ef.
//
// You create an error format string by building a string with
// ErrorFormatToken types.
//
// Do not use formatting verbs supported by the fmt package in the error
// format string.
func (e Error) Format(ef string) string {
	return newFormatter().format(e, ef)
}

// Walk calls the step function for each error in the error chain, including the
// receiver, and continues until either the last error is reached, or the step
// function returns false.
func (e Error) Walk(step func(err error) bool) {
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

// Depth returns the number of nested errors in the receiver. That is, the
// number of errors after, but not including, this error in the error chain.
//
// For example, if an error has no nested errors, then its depth is 0.
func (e Error) Depth() uint {
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

// Trace returns a prettified string representation of the error chain.
func (e Error) Trace() string {
	d := e.Depth()
	if d == 0 {
		return fmt.Sprintf("%s %s", e.Caller, e.Error())
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
			em = fmt.Sprintf("%s %d: %s %s", p, d, we.Caller, we.Error())
		} else if we, ok := err.(*Error); ok {
			em = fmt.Sprintf("%s %d: %s %s", p, d, we.Caller, we.Error())
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

func (e Error) Error() string {
	var s string
	e.Walk(func(err error) bool {
		if we, ok := err.(Error); ok {
			s += fmt.Sprintf("%+v", we.Context)
		} else if we, ok := err.(*Error); ok {
			s += fmt.Sprintf("%+v", we.Context)
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

// Unwrap returns the wrapped error or nil if one doesn't exist.
func (e Error) Unwrap() error {
	return e.inner
}

// As finds the first error in the error chain that matches target, and if so,
// sets target to that error value and returns true. Otherwise, it returns
// false.
func (e Error) As(target interface{}) bool {
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

// Is reports whether any error in the error chain matches target.
//
// The chain consists of the receiver followed by the sequence of errors
// obtained by repeatedly calling Unwrap.
//
// The receiver considered to match a target if it is equal to that target or
// if it implements a method Is(error) bool such that Is(target) returns true.
func (e Error) Is(target error) bool {
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

// JSON Marshaler interface methods

// MarshalJSON marshals the error in to JSON data.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(newJSONWError(e))
}
