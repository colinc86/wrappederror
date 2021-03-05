// Package wrappederror contains the WrappedError type.
package wrappederror

import (
	"encoding"
	"fmt"
)

// WrappedError types wrap an error.
type WrappedError interface {
	error
	fmt.Stringer
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	// Depth returns the number of nested errors in the receiver.
	Depth() uint

	// Trace returns a prettified string representation of the wrapped error.
	Trace() string

	// The file that the error was created in.
	File() string

	// The function that the error was created in.
	Function() string

	// The line that the error was created on.
	Line() int

	// Unwrap unwraps the wrapped error.
	Unwrap() error
}
