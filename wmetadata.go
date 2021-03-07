package wrappederror

import (
	"fmt"
	"time"
)

// A type containing error metadata.
type wMetadata struct {
	time          time.Time
	index         int
	similarErrors int
}

// Initializers

// newWMetadata creates new metadata with the specified components.
func newWMetadata(
	time time.Time,
	index int,
	similarErrors int,
) *wMetadata {
	return &wMetadata{
		time:          time,
		index:         index,
		similarErrors: similarErrors,
	}
}

// Methods

// currentMetadata gets the current metadata that should be added to an error.
// The function requires the error's inner error to find similar errors.
func currentMetadata(err error) *wMetadata {
	return &wMetadata{
		time:          time.Now(),
		index:         getAndIncrementNextErrorIndex(),
		similarErrors: getSimilarErrorCount(err),
	}
}

// Stringer interface methods

func (m wMetadata) String() string {
	return fmt.Sprintf(
		"(#%d) (≈%d) %s",
		m.index,
		m.similarErrors,
		m.time,
	)
}

// Process interface methods

func (m wMetadata) Time() time.Time {
	return m.time
}

func (m wMetadata) Index() int {
	return m.index
}

func (m wMetadata) Similar() int {
	return m.similarErrors
}
