package wrappederror

import (
	"fmt"
	"time"
)

// A type containing error metadata.
type wMetadata struct {
	ErrorTime     time.Time `json:"time"`
	ErrorIndex    int       `json:"index"`
	SimilarErrors int       `json:"similar"`
}

// Initializers

// newWMetadata creates new metadata with the specified components.
func newWMetadata(
	time time.Time,
	index int,
	similarErrors int,
) *wMetadata {
	return &wMetadata{
		ErrorTime:     time,
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
		ErrorIndex:    getAndIncrementNextErrorIndex(),
		SimilarErrors: getSimilarErrorCount(err),
	}
}

// Stringer interface methods

func (m wMetadata) String() string {
	return fmt.Sprintf(
		"(#%d) (≈%d) %s",
		m.ErrorIndex,
		m.SimilarErrors,
		m.ErrorTime,
	)
}

// Process interface methods

func (m wMetadata) Time() time.Time {
	return m.ErrorTime
}

func (m wMetadata) Index() int {
	return m.ErrorIndex
}

func (m wMetadata) Similar() int {
	return m.SimilarErrors
}
