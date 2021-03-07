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

// A type containing call information.
type wCaller struct {
	fileName     string
	functionName string
	lineNumber   int
	stackTrace   []byte
	source       []byte
}

// Initializers

// newWCaller creates a new caller with the specified components.
func newWCaller(
	fileName string,
	functionName string,
	lineNumber int,
	stackTrace []byte,
	source []byte,
) *wCaller {
	return &wCaller{
		fileName:     fileName,
		functionName: functionName,
		lineNumber:   lineNumber,
		stackTrace:   stackTrace,
		source:       source,
	}
}

// Methods

// currentCaller gets the current caller with the given depth.
func currentCaller(skip int) *wCaller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		so, _ := getSource(fp, ln, SourceFragmentRadius())

		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return newWCaller(fin, f.Name(), ln, st, so)
		}
		return newWCaller(fin, callerFunctionNameUnknown, ln, st, so)
	}

	// *Wah, wah, wah* sound effect.
	return newWCaller(
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		st,
		nil,
	)
}

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

	for s.Scan() {
		l++

		if l >= li && l <= ui {
			lnb := append(s.Bytes(), []byte("\n")...)
			b = append(b, lnb...)
		} else if l > ui {
			break
		}
	}

	if len(b) > 0 {
		eb := []byte("...\n")
		b = append(eb, b...)
		b = append(b, eb...)
	}

	return b, nil
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

func (c wCaller) Source() string {
	return string(c.source)
}
