// Package wrappederror contains the Error type.
package wrappederror

import (
	"encoding/json"
	"errors"
	"fmt"
)

// wError types wrap an error.
type wError struct {
	context interface{}

	// The inner error that this wrapped error wraps.
	inner error

	// The caller that invoked the `New` function.
	caller *wCaller

	// Information about the current process when the error was created.
	process *wProcess

	// Metadata pertaining to the error.
	metadata *wMetadata
}

// Initializers

// New creates and returns a new wrapped error.
func New(err error, ctx interface{}) Error {
	var caller *wCaller
	if CaptureCaller() {
		caller = currentCaller(2)
	}

	var process *wProcess
	if CaptureProcess() {
		process = currentProcess()
	}

	return &wError{
		context:  ctx,
		inner:    err,
		caller:   caller,
		process:  process,
		metadata: currentMetadata(err),
	}
}

// (Wrapped) Error interface methods

func (e wError) Context() interface{} {
	return e.context
}

func (e wError) Caller() Caller {
	return e.caller
}

func (e wError) Process() Process {
	return e.process
}

func (e wError) Metadata() Metadata {
	return e.metadata
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

// TextMarshaler interface methods

func (e wError) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// JSON Marshaler interface methods

func (e wError) MarshalJSON() ([]byte, error) {
	return json.Marshal(newJSONWError(e))
}
