package wrappederror

import (
	"fmt"
)

// Caller types contain call information.
type Caller interface {
	fmt.Stringer

	// The file the caller was created in.
	File() string

	// The function the caller was created in.
	Function() string

	// The line the caller was created on.
	Line() int

	// Stack provides a stack trace of the goroutine the caller was created on.
	Stack() string

	// Source returns raw source code around the line that the caller was created
	// on. This function will return an empty string if the process is not
	// currently being debugged.
	Source() string
}
