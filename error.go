package wrappederror

import "encoding"

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

	// Process returns the error's process.
	//
	// Use this value to examine the current process's information such as number
	// of goroutines when the error was created.
	Process() Process

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
