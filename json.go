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
	Caller   *wCaller    `json:"caller"`
	Process  *wProcess   `json:"process"`
	Metadata *wMetadata  `json:"metadata"`
	Context  interface{} `json:"context"`
	Depth    int         `json:"depth"`
	Inner    interface{} `json:"wraps"`
}

// Initializers

func newJSONErrorOrWError(err error) interface{} {
	if err == nil {
		return nil
	}

	if we, ok := err.(wError); ok {
		return newJSONWError(we)
	} else if we, ok := err.(*wError); ok {
		return newJSONWError(*we)
	}

	return newJSONError(err)
}

func newJSONError(err error) *jsonError {
	return &jsonError{
		Error: err.Error(),
		Inner: newJSONErrorOrWError(errors.Unwrap(err)),
	}
}

func newJSONWError(e wError) interface{} {
	if packageState.configuration.MarshalMinimalJSON() {
		return newJSONWErrorMinimal(e)
	}
	return newJSONWErrorFull(e)
}

func newJSONWErrorMinimal(e wError) *jsonWErrorMinimal {
	return &jsonWErrorMinimal{
		Context:  e.context,
		Depth:    int(e.Depth()),
		Time:     e.metadata.ErrorTime,
		Duration: e.metadata.ErrorDuration,
		Index:    e.metadata.ErrorIndex,
		Similar:  e.metadata.SimilarErrors,
		File:     e.caller.FileName,
		Function: e.caller.FunctionName,
		Line:     e.caller.LineNumber,
		Inner:    newJSONErrorOrWError(e.inner),
	}
}

func newJSONWErrorFull(e wError) *jsonWErrorFull {
	return &jsonWErrorFull{
		Caller:   e.caller,
		Process:  e.process,
		Metadata: e.metadata,
		Context:  e.context,
		Depth:    int(e.Depth()),
		Inner:    newJSONErrorOrWError(e.inner),
	}
}
