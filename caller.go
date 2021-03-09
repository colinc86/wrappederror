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

// Caller types contain call information.
type Caller struct {

	// The caller's file.
	File string `json:"file"`

	// The caller's function name.
	Function string `json:"function"`

	// The caller's line number.
	Line int `json:"line"`

	// A stack trace of the goroutine that created the caller.
	StackTrace string `json:"stackTrace"`

	// Fragment returns raw source code around the line that the caller was
	// created on. This function will return an empty string if the process is not
	// currently being debugged.
	Fragment *SourceFragment `json:"sourceFragment"`
}

// Initializers

// currentCaller gets the current caller with the given skip.
func newCaller(skip int, captureFragment bool, fragmentRadius int) *Caller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		var sf *SourceFragment
		if captureFragment {
			sf, _ = newSourceFragment(fp, ln, fragmentRadius)
		}

		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return &Caller{
				fin,
				f.Name(),
				ln,
				string(st),
				sf,
			}
		}

		return &Caller{
			fin,
			callerFunctionNameUnknown,
			ln,
			string(st),
			sf,
		}
	}

	// *Wah, wah, wah* sound effect.
	return &Caller{
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		string(st),
		nil,
	}
}

// Stringer interface methods

func (c Caller) String() string {
	return fmt.Sprintf(
		"%s (%s:%d)",
		c.Function,
		c.File,
		c.Line,
	)
}
