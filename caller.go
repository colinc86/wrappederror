package wrappederror

import (
	"bufio"
	"fmt"
	"os"
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

	// Source returns raw source code around the line that the caller was created
	// on. This function will return an empty string if the process is not
	// currently being debugged.
	SourceFragment string `json:"sourceFragment"`
}

// Initializers

// currentCaller gets the current caller with the given skip.
func newCaller(skip int, fragmentRadius int) *Caller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		so, _ := getSource(
			fp,
			ln,
			fragmentRadius,
		)

		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return &Caller{
				fin,
				f.Name(),
				ln,
				string(st),
				string(so),
			}
		}

		return &Caller{
			fin,
			callerFunctionNameUnknown,
			ln,
			string(st),
			string(so),
		}
	}

	// *Wah, wah, wah* sound effect.
	return &Caller{
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		string(st),
		"",
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

// Non-exported functions

// getSource gets lines of source from filePath around lineNumber.
func getSource(filePath string, lineNumber int, radius int) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	l := 0
	li := lineNumber - radius
	ui := lineNumber + radius
	var b []byte

	ali := 0
	aui := 0

	for s.Scan() {
		l++

		if l >= li && l <= ui {
			if ali == 0 {
				ali = l
			}
			aui = l
			lnb := append(s.Bytes(), []byte("\n")...)
			b = append(b, lnb...)
		} else if l > ui {
			break
		}
	}

	hb := []byte(fmt.Sprintf("[%d-%d] %s\n", ali, aui, filePath))

	return append(hb, b...), nil
}
