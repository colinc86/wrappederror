package wrappederror

import (
	"fmt"
	"time"
)

// Metadata types contain metadata about an error.
type Metadata interface {
	fmt.Stringer

	// The time that the error was created.
	Time() time.Time

	// The index of the error.
	//
	// Error indexes begin at 1 and incriment for each error created during the
	// process's execution.
	//
	// To start at a different index, use the `SetNextErrorIndex` function.
	Index() int

	// The number of similar errors when this error was created.
	//
	// A similar error is an error that has the same inner error. A hashmap is
	// maintained of inner error `Error() string` value hashes.
	//
	// To turn this off, use the `SetTrackSimilarErrors` function. When tracking
	// is off, this method always returns 0.
	Similar() int
}
