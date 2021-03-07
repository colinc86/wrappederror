package wrappederror

import (
	"sync"
	"time"
)

// state types keep track of the package's current state.
type state struct {
	errorHashMap      *errorMap
	processLaunchTime *configValue
	configuration     *Configuration
}

// Initializers

// newState creates and returns a new state structure.
func newState() *state {
	return &state{
		errorHashMap:      newErrorMap(),
		processLaunchTime: newConfigValue(time.Now()),
		configuration:     newConfiguration(),
	}
}

// Non-exported methods.

// reset resets the state to its initial value.
func (s *state) reset() {
	s.errorHashMap.hashMap = new(sync.Map)
	s.processLaunchTime.set(time.Now())
	s.configuration = newConfiguration()
}

// getSimilarErrorCount gets and returns the number of errors in the error hash
// map equal to err.
func (s state) getSimilarErrorCount(err error) int {
	if !s.configuration.TrackSimilarErrors() || err == nil {
		return 0
	}

	st := s.errorHashMap.similarErrors(err)
	s.errorHashMap.addError(err)
	return st
}

// getDurationSinceLaunch gets the current duration since the process was
// launched.
func (s state) getDurationSinceLaunch() time.Duration {
	n := time.Now()
	lt := s.processLaunchTime.get().(time.Time)
	return n.Sub(lt)
}
