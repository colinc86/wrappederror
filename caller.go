package wrappederror

import (
	"fmt"
	"path"
	"runtime"
)

// Caller types contain call information.
type Caller interface {
	fmt.Stringer

	File() string
	Function() string
	Line() int
}

// Values to use when we can't get components of the caller.
const (
	callerFileNameUnknown     string = "unknown file"
	callerFunctionNameUnknown string = "unknown function"
	callerLineNumberUnknown   int    = 0
)

// Used to encode callers in to binary data.
type callerGob struct {
	Fi string
	Fu string
	Li int
}

// A type containing call information.
type caller struct {
	fileName     string
	functionName string
	lineNumber   int
}

// newCaller creates a new caller with the specified components.
func newCaller(
	fileName string,
	functionName string,
	lineNumber int,
) *caller {
	return &caller{
		fileName:     fileName,
		functionName: functionName,
		lineNumber:   lineNumber,
	}
}

// currentCaller gets the current caller with the given depth.
func currentCaller(skip int) *caller {
	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return newCaller(fin, f.Name(), ln)
		}
		return newCaller(fin, callerFunctionNameUnknown, ln)
	}

	// *Wah, wah, wah* sound effect.
	return newCaller(
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
	)
}

// String interface methods

func (c caller) String() string {
	return fmt.Sprintf(
		"%s (%s:%d)",
		c.functionName,
		c.fileName,
		c.lineNumber,
	)
}

// Caller interface methods

func (c caller) File() string {
	return c.fileName
}

func (c caller) Function() string {
	return c.functionName
}

func (c caller) Line() int {
	return c.lineNumber
}
