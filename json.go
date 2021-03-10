package wrappederror

import (
	"errors"
	"time"
)

// The generic JSON error type.
type jsonError struct {
	Error string      `json:"error"`
	Inner interface{} `json:"wraps,omitempty"`
}

// The minimal JSON error type.
type jsonWErrorMinimal struct {
	Context  interface{}   `json:"context"`
	Depth    int           `json:"depth"`
	Time     time.Time     `json:"time"`
	Duration time.Duration `json:"duration"`
	Index    int           `json:"index"`
	Similar  int           `json:"simlar,omitempty"`
	File     string        `json:"file"`
	Function string        `json:"function"`
	Line     int           `json:"line"`
	Inner    interface{}   `json:"wraps,omitempty"`
}

// The full JSON error type.
type jsonWErrorFull struct {
	Caller   *Caller     `json:"caller"`
	Process  *Process    `json:"process"`
	Metadata *Metadata   `json:"metadata"`
	Context  interface{} `json:"context"`
	Depth    int         `json:"depth"`
	Inner    interface{} `json:"wraps"`
}

// Initializers

// newJSONErrorOrWError examines an error for its type and either initializes
// a new jsonError or jsonWError.
func newJSONErrorOrWError(err error) interface{} {
	if err == nil {
		return nil
	}

	if we, ok := err.(*Error); ok {
		return newJSONWError(*we)
	}

	return newJSONError(err)
}

// newJSONError creates a new jsonError.
func newJSONError(err error) *jsonError {
	return &jsonError{
		Error: err.Error(),
		Inner: newJSONErrorOrWError(errors.Unwrap(err)),
	}
}

// newJSONWError creates a new jsonWError.
func newJSONWError(e Error) interface{} {
	if packageState.config.MarshalMinimalJSON() {
		return newJSONWErrorMinimal(e)
	}
	return newJSONWErrorFull(e)
}

// newJSONWErrorMinimal creates a new minimal json error.
func newJSONWErrorMinimal(e Error) *jsonWErrorMinimal {
	return &jsonWErrorMinimal{
		Context:  e.context,
		Depth:    int(e.Depth()),
		Time:     e.Metadata.Time,
		Duration: e.Metadata.Duration,
		Index:    e.Metadata.Index,
		Similar:  e.Metadata.Similar,
		File:     e.Caller.File,
		Function: e.Caller.Function,
		Line:     e.Caller.Line,
		Inner:    newJSONErrorOrWError(e.inner),
	}
}

// newJSONWErrorFull creates a new full json error.
func newJSONWErrorFull(e Error) *jsonWErrorFull {
	return &jsonWErrorFull{
		Caller:   e.Caller,
		Process:  e.Process,
		Metadata: e.Metadata,
		Context:  e.context,
		Depth:    int(e.Depth()),
		Inner:    newJSONErrorOrWError(e.inner),
	}
}
