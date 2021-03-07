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

	// The duration since the process was launched and when the error was created.
	//
	// This value mimics the monotonic clock reading appended to the end of
	// strings returned by the `time.String() string` method. There will be a
	// slight difference in duration depending on the executable's package loading
	// order.
	Duration() time.Duration

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

// Implementation

// A type containing error metadata.
type wMetadata struct {
	ErrorTime     time.Time     `json:"time"`
	ErrorDuration time.Duration `json:"duration"`
	ErrorIndex    int           `json:"index"`
	SimilarErrors int           `json:"similar"`
}

// Initializers

// newWMetadata creates new metadata with the specified components.
func newWMetadata(
	time time.Time,
	duration time.Duration,
	index int,
	similarErrors int,
) *wMetadata {
	return &wMetadata{
		ErrorTime:     time,
		ErrorDuration: duration,
		ErrorIndex:    index,
		SimilarErrors: similarErrors,
	}
}

// Methods

// currentMetadata gets the current metadata that should be added to an error.
// The function requires the error's inner error to find similar errors.
func currentMetadata(err error) *wMetadata {
	return &wMetadata{
		ErrorTime:     time.Now(),
		ErrorDuration: packageState.getDurationSinceLaunch(),
		ErrorIndex:    packageState.configuration.getAndIncrementNextErrorIndex(),
		SimilarErrors: packageState.getSimilarErrorCount(err),
	}
}

// Stringer interface methods

func (m wMetadata) String() string {
	return fmt.Sprintf(
		"(#%d) (≈%d) (+%f) %s",
		m.ErrorIndex,
		m.SimilarErrors,
		m.ErrorDuration.Seconds(),
		m.ErrorTime,
	)
}

// Process interface methods

func (m wMetadata) Time() time.Time {
	return m.ErrorTime
}

func (m wMetadata) Duration() time.Duration {
	return m.ErrorDuration
}

func (m wMetadata) Index() int {
	return m.ErrorIndex
}

func (m wMetadata) Similar() int {
	return m.SimilarErrors
}
