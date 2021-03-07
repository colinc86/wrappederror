package wrappederror

import (
	"errors"
	"runtime"
	"time"
)

const (
	jsonWErrorMinimalSize = "minimal"
	jsonWErrorFullSize    = "full"
)

type jsonError struct {
	Error string      `json:"error"`
	Inner interface{} `json:"wraps,omitempty"`
}

type jsonWErrorMinimal struct {
	Size     string      `json:"_size"`
	Context  interface{} `json:"context"`
	Depth    int         `json:"depth"`
	Time     time.Time   `json:"time"`
	Index    int         `json:"index"`
	Similar  int         `json:"simlar,omitempty"`
	File     string      `json:"file"`
	Function string      `json:"function"`
	Line     int         `json:"line"`
	Inner    interface{} `json:"wraps,omitempty"`
}

type jsonWErrorFull struct {
	Size     string         `json:"_size"`
	Caller   *jsonWCaller   `json:"caller"`
	Process  *jsonWProcess  `json:"process"`
	Metadata *jsonWMetadata `json:"metadata"`
	Context  interface{}    `json:"context"`
	Depth    int            `json:"depth"`
	Inner    interface{}    `json:"wraps"`
}

type jsonWCaller struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
	Stack    string `json:"stackTrace"`
	Source   string `json:"sourceFragment"`
}

type jsonWMetadata struct {
	Time    time.Time `json:"time"`
	Index   int       `json:"index"`
	Similar int       `json:"similar,omitempty"`
}

type jsonWProcess struct {
	Routines int               `json:"goroutines"`
	CPUs     int               `json:"cpus"`
	CGO      int               `json:"cgos"`
	MemStats *runtime.MemStats `json:"memory,omitempty"`
}

// Initializers

func newJSONWError(e wError) interface{} {
	if MarshalMinimalJSON() {
		return newJSONWErrorMinimal(e)
	}
	return newJSONWErrorFull(e)
}

func newJSONWErrorMinimal(e wError) *jsonWErrorMinimal {
	return &jsonWErrorMinimal{
		Size:     jsonWErrorMinimalSize,
		Context:  e.context,
		Depth:    int(e.Depth()),
		Time:     e.metadata.time,
		Index:    e.metadata.index,
		Similar:  e.metadata.similarErrors,
		File:     e.caller.fileName,
		Function: e.caller.functionName,
		Line:     e.caller.lineNumber,
		Inner:    newJSONErrorOrWError(e.inner),
	}
}

func newJSONWErrorFull(e wError) *jsonWErrorFull {
	return &jsonWErrorFull{
		Size:     jsonWErrorFullSize,
		Caller:   newJSONWCaller(e.caller),
		Process:  newJSONWProcess(e.process),
		Metadata: newJSONWMetadata(e.metadata),
		Context:  e.context,
		Depth:    int(e.Depth()),
		Inner:    newJSONErrorOrWError(e.inner),
	}
}

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

func newJSONWCaller(c *wCaller) *jsonWCaller {
	return &jsonWCaller{
		File:     c.fileName,
		Function: c.functionName,
		Line:     c.lineNumber,
		Stack:    string(c.stackTrace),
		Source:   string(c.source),
	}
}

func newJSONWProcess(p *wProcess) *jsonWProcess {
	return &jsonWProcess{
		Routines: p.numRoutines,
		CPUs:     p.numCPUs,
		CGO:      p.numCGO,
		MemStats: p.memStats,
	}
}

func newJSONWMetadata(m *wMetadata) *jsonWMetadata {
	return &jsonWMetadata{
		Time:    m.time,
		Index:   m.index,
		Similar: m.similarErrors,
	}
}
