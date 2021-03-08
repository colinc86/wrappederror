package wrappederror

import (
	"fmt"
	"time"
)

// Metadata types contain metadata about an error.
type Metadata struct {

	// The time that the error was created.
	Time time.Time `json:"time"`

	// The duration since the process was launched and when the error was created.
	//
	// This value mimics the monotonic clock reading appended to the end of
	// strings returned by the `time.String() string` method. There will be a
	// slight difference in duration depending on the executable's package loading
	// order.
	Duration time.Duration `json:"duration"`

	// The index of the error.
	//
	// Error indexes begin at 1 and incriment for each error created during the
	// process's execution.
	//
	// To start at a different index, use the `SetNextErrorIndex` function.
	Index int `json:"index"`

	// The number of similar errors when this error was created.
	//
	// A similar error is an error that has the same inner error. A hashmap is
	// maintained of inner error `Error() string` value hashes.
	//
	// To turn this off, use the `SetTrackSimilarErrors` function. When tracking
	// is off, this method always returns 0.
	Similar int `json:"similar"`
}

// Initializers

// newMetadata creates metadata that should be added to an error. The function
// requires the error's inner error to find similar errors.
func newMetadata(err error) *Metadata {
	return &Metadata{
		Time:     time.Now(),
		Duration: packageState.getDurationSinceLaunch(),
		Index:    packageState.config.getAndIncrementNextErrorIndex(),
		Similar:  packageState.getSimilarErrorCount(err),
	}
}

// Stringer interface methods

func (m Metadata) String() string {
	return fmt.Sprintf(
		"(#%d) (≈%d) (+%f) %s",
		m.Index,
		m.Similar,
		m.Duration.Seconds(),
		m.Time,
	)
}
