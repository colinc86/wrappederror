package wrappederror

import (
	"time"
)

// state types keep track of the package's current state.
type state struct {
	errorMap          *errorMap
	serverityTable    *severityTable
	processLaunchTime *safeValue
	config            *Configuration
}

// Initializers

// newState creates and returns a new state structure.
func newState() *state {
	s := new(state)
	s.reset()
	return s
}

// Methods

// reset resets the state to its initial value.
func (s *state) reset() {
	s.errorMap = newErrorMap()
	s.serverityTable = newSeverityTable()
	s.processLaunchTime = newSafeValue(time.Now())
	s.config = newConfiguration()
}

// getSimilarErrorCount gets and returns the number of errors in the error hash
// map equal to err.
func (s state) getSimilarErrorCount(err error) int {
	if !s.config.TrackSimilarErrors() || err == nil {
		return 0
	}

	st := s.errorMap.similarErrors(err)
	s.errorMap.addError(err)
	return st
}

// registerSeverity registers the severity with the state's severity table.
func (s state) registerSeverity(severity ErrorSeverity) error {
	return s.serverityTable.register(severity)
}

// unregisterSeverity unregisters the severity from the state's severity table.
func (s state) unregisterSeverity(severity ErrorSeverity) {
	s.serverityTable.unregister(severity)
}

// getBestMatchSeverity gets the best match severity for the given error.
func (s state) getBestMatchSeverity(err error) *ErrorSeverity {
	if s := s.serverityTable.bestMatch(err); s != errorSeverityUnknown {
		return &s
	}
	return nil
}

// getDurationSinceLaunch gets the current duration since the process was
// launched.
func (s state) getDurationSinceLaunch() time.Duration {
	n := time.Now()
	lt := s.processLaunchTime.get().(time.Time)
	return n.Sub(lt)
}
