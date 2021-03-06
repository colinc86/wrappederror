package wrappederror

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
)

// Values to use when we can't get components of the caller.
const (
	callerFileNameUnknown     string = "unknown file"
	callerFunctionNameUnknown string = "unknown function"
	callerLineNumberUnknown   int    = 0
)

// A type containing call information.
type wCaller struct {
	fileName     string
	functionName string
	lineNumber   int
	stackTrace   []byte
}

// Initializers

// newWCaller creates a new caller with the specified components.
func newWCaller(
	fileName string,
	functionName string,
	lineNumber int,
	stackTrace []byte,
) *wCaller {
	return &wCaller{
		fileName:     fileName,
		functionName: functionName,
		lineNumber:   lineNumber,
		stackTrace:   stackTrace,
	}
}

// Methods

// currentCaller gets the current caller with the given depth.
func currentCaller(skip int) *wCaller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return newWCaller(fin, f.Name(), ln, st)
		}
		return newWCaller(fin, callerFunctionNameUnknown, ln, st)
	}

	// *Wah, wah, wah* sound effect.
	return newWCaller(
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		st,
	)
}

// String interface methods

func (c wCaller) String() string {
	return fmt.Sprintf(
		"%s (%s:%d)",
		c.functionName,
		c.fileName,
		c.lineNumber,
	)
}

// Caller interface methods

func (c wCaller) File() string {
	return c.fileName
}

func (c wCaller) Function() string {
	return c.functionName
}

func (c wCaller) Line() int {
	return c.lineNumber
}

func (c wCaller) Stack() string {
	return string(c.stackTrace)
}
