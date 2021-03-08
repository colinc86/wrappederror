package wrappederror

import (
	"errors"
	"sync"
)

// ErrSeverityAlreadyRegistered indicates that the error severity has already
// been registered.
var ErrSeverityAlreadyRegistered = errors.New("severity already registered")

// severityTable keeps track of error severities.
type severityTable struct {
	severities      []ErrorSeverity
	severitiesMutex *sync.RWMutex
}

// Initializers

// newSeverityTable creates and returns a new severity table.
func newSeverityTable() *severityTable {
	return &severityTable{
		severitiesMutex: new(sync.RWMutex),
	}
}

// Methods

// register registers a new error severity. If the severity already exists, then
// it returns an ErrSeverityAlreadyRegistered error.
func (t *severityTable) register(severity ErrorSeverity) error {
	t.severitiesMutex.Lock()
	defer t.severitiesMutex.Unlock()

	for _, s := range t.severities {
		if s.equals(severity) {
			return ErrSeverityAlreadyRegistered
		}
	}

	t.severities = append(t.severities, severity)
	return nil
}

// unregister unregisters the severity.
func (t *severityTable) unregister(severity ErrorSeverity) {
	t.severitiesMutex.Lock()
	defer t.severitiesMutex.Unlock()

	for i, s := range t.severities {
		if s.equals(severity) {
			t.severities = append(t.severities[:i], t.severities[i+1:]...)
			return
		}
	}
}

// bestMatch returns the best match or errorSeverityUnknown if none was found
// with a match greater than 0.0.
//
// It walks the entire error chain beginning with err and finds the best match
// error severity.
func (t *severityTable) bestMatch(err error) ErrorSeverity {
	t.severitiesMutex.RLock()
	defer t.severitiesMutex.RUnlock()

	bestMatch := 0.0
	bestMatchErrorSeverity := errorSeverityUnknown

	e := err
	for e != nil {
		for _, s := range t.severities {
			m := s.match(e)
			if m > bestMatch {
				bestMatch = m
				bestMatchErrorSeverity = s
			}
		}

		e = errors.Unwrap(e)
	}

	return bestMatchErrorSeverity
}
