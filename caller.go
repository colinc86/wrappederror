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

	// The stack trace of the calling goroutine.
	Stack() string
}
