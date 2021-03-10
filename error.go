package wrappederror

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Error structure string constants.
var (
	errorChainDelimiter            string = ": "
	errorTraceFirstItemDecoration  string = "┌"
	errorTraceMiddleItemDecoration string = "├"
	errorTraceLastItemDecoration   string = "└"
)

// Error types wrap an error and provide context, caller and process
// information.
type Error struct {

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

	// The error's context.
	//
	// When an error is wrapped, it is given context. An error's context can be a
	// string description or any other type of information that is pertinent to
	// the error being wrapped.
	context interface{}

	// The inner error that this wrapped error wraps.
	inner error
}

// Initializers

// New creates and returns a new error with an inner error and context.
func New(err error, ctx interface{}) *Error {
	var caller *Caller
	if packageState.config.CaptureCaller() {
		caller = newCaller(
			2,
			packageState.config.CaptureSourceFragments(),
			packageState.config.SourceFragmentRadius(),
		)
	}

	var process *Process
	if packageState.config.CaptureProcess() {
		process = newProcess()
	}

	return &Error{
		context:  ctx,
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

// Chain returns the error chain as a slice with the receiver at index 0.
func (e Error) Chain() []error {
	var c []error
	var ce error = e
	for ce != nil {
		c = append(c, ce)
		ce = errors.Unwrap(ce)
	}
	return c
}

// Walk calls the step function for each error in the error chain, including the
// receiver, and continues until either the last error is reached, or the step
// function returns false.
func (e Error) Walk(step func(err error) bool) {
	for _, ecv := range e.Chain() {
		if !step(ecv) {
			break
		}
	}
}

// ErrorWithDepth returns the error in the chain with the given depth.
func (e Error) ErrorWithDepth(depth int) error {
	c := e.Chain()
	if depth < 0 || depth >= len(c) {
		return nil
	}
	return c[len(c)-depth-1]
}

// ErrorWithIndex returns the error in the chain with the given index. This is
// the inverse index of depth.
//
// For example, if the error chain is [e0, e1, e2] with depths [2, 1, 0], then
// the error at index 0 is e0, the error at index 1 is e1, and so forth.
func (e Error) ErrorWithIndex(index int) error {
	c := e.Chain()
	if index < 0 || index >= len(c) {
		return nil
	}
	return c[index]
}

// Depth returns the number of nested errors in the receiver. That is, the
// number of errors after, but not including, this error in the error chain.
//
// For example, if an error has no nested errors, then its depth is 0.
func (e Error) Depth() int {
	return len(e.Chain()) - 1
}

// Trace returns a prettified string representation of the error chain.
func (e Error) Trace() string {
	d := e.Depth()
	if d == 0 {
		return fmt.Sprintf("%s %s", e.Caller, e.Error())
	}

	var msg string
	c := e.Chain()

	for i, err := range c {
		end := i == len(c)-1
		var p string
		if i == 0 {
			p = errorTraceFirstItemDecoration
		} else if end {
			p = errorTraceLastItemDecoration
		} else {
			p = errorTraceMiddleItemDecoration
		}

		var em string
		if we, ok := err.(Error); ok {
			em = fmt.Sprintf("%s %d: %s %+v", p, d-i, we.Caller, we.context)
		} else if we, ok := err.(*Error); ok {
			em = fmt.Sprintf("%s %d: %s %+v", p, d-i, we.Caller, we.context)
		} else {
			em = fmt.Sprintf("%s %d: %s", p, d-i, err.Error())
		}

		msg += em
		if !end {
			msg += "\n"
		}
	}

	return msg
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

	if we, ok := target.(*Error); ok {
		e.Walk(func(err error) bool {
			if wei, ok := err.(Error); ok && wei.Error() == we.Error() {
				target = &wei
				as = true
				return false
			}
			return true
		})
	}

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
		if err.Error() == target.Error() {
			is = true
			return false
		}
		return true
	})
	return is
}

// Context returns the error's context.
func (e Error) Context() interface{} {
	return e.context
}

// Error interface methods

func (e Error) Error() string {
	var s string
	c := e.Chain()

	for i, err := range c {
		if we, ok := err.(Error); ok {
			s += fmt.Sprintf("%+v", we.context)
		} else if we, ok := err.(*Error); ok {
			s += fmt.Sprintf("%+v", we.context)
		} else {
			s += err.Error()
		}

		if i < len(c)-1 {
			s += errorChainDelimiter
		}
	}

	return s
}

// JSON Marshaler interface methods

// MarshalJSON marshals the error in to JSON data.
func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(newJSONWError(e))
}
