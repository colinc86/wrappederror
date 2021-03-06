package wrappederror

import (
	"fmt"
	"path"
	"runtime"
	"runtime/debug"
)

// Caller types contain call information.
type Caller interface {
	fmt.Stringer

	// The file the caller was created in.
	File() string

	// The function the caller was created in.
	Function() string

	// The line the caller was created on.
	Line() int

	// The stack trace of the calling goroutine.
	Stack() string
}

// Caller implementation

// Values to use when we can't get components of the caller.
const (
	callerFileNameUnknown     string = "unknown file"
	callerFunctionNameUnknown string = "unknown function"
	callerLineNumberUnknown   int    = 0
)

// A type containing call information.
type caller struct {
	fileName     string
	functionName string
	lineNumber   int
	stackTrace   []byte
}

// Initializers

// newCaller creates a new caller with the specified components.
func newCaller(
	fileName string,
	functionName string,
	lineNumber int,
	stackTrace []byte,
) *caller {
	return &caller{
		fileName:     fileName,
		functionName: functionName,
		lineNumber:   lineNumber,
		stackTrace:   stackTrace,
	}
}

// Methods

// currentCaller gets the current caller with the given depth.
func currentCaller(skip int) *caller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return newCaller(fin, f.Name(), ln, st)
		}
		return newCaller(fin, callerFunctionNameUnknown, ln, st)
	}

	// *Wah, wah, wah* sound effect.
	return newCaller(
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		st,
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

func (c caller) Stack() string {
	return string(c.stackTrace)
}
