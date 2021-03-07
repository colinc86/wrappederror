package wrappederror

import (
	"bufio"
	"fmt"
	"os"
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

	// Stack provides a stack trace of the goroutine the caller was created on.
	Stack() string

	// Source returns raw source code around the line that the caller was created
	// on. This function will return an empty string if the process is not
	// currently being debugged.
	Source() string
}

// Implementation

// Values to use when we can't get components of the caller.
const (
	callerFileNameUnknown     string = "unknown file"
	callerFunctionNameUnknown string = "unknown function"
	callerLineNumberUnknown   int    = 0
)

// A type containing call information.
type wCaller struct {
	FileName       string `json:"file"`
	FunctionName   string `json:"function"`
	LineNumber     int    `json:"line"`
	StackTrace     string `json:"stackTrace"`
	SourceFragment string `json:"sourceFragment"`
}

// Initializers

// newWCaller creates a new caller with the specified components.
func newWCaller(
	fileName string,
	functionName string,
	lineNumber int,
	stackTrace string,
	source string,
) *wCaller {
	return &wCaller{
		FileName:       fileName,
		FunctionName:   functionName,
		LineNumber:     lineNumber,
		StackTrace:     stackTrace,
		SourceFragment: source,
	}
}

// Methods

// currentCaller gets the current caller with the given depth.
func currentCaller(skip int) *wCaller {
	st := debug.Stack()

	if pc, fp, ln, ok := runtime.Caller(skip); ok {
		so, _ := getSource(
			fp,
			ln,
			packageState.configuration.SourceFragmentRadius(),
		)

		_, fin := path.Split(fp)
		if f := runtime.FuncForPC(pc); f != nil {
			return newWCaller(fin, f.Name(), ln, string(st), string(so))
		}
		return newWCaller(
			fin,
			callerFunctionNameUnknown,
			ln,
			string(st),
			string(so),
		)
	}

	// *Wah, wah, wah* sound effect.
	return newWCaller(
		callerFileNameUnknown,
		callerFunctionNameUnknown,
		callerLineNumberUnknown,
		string(st),
		"",
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

// Stringer interface methods

func (c wCaller) String() string {
	return fmt.Sprintf(
		"%s (%s:%d)",
		c.FunctionName,
		c.FileName,
		c.LineNumber,
	)
}

// Caller interface methods

func (c wCaller) File() string {
	return c.FileName
}

func (c wCaller) Function() string {
	return c.FunctionName
}

func (c wCaller) Line() int {
	return c.LineNumber
}

func (c wCaller) Stack() string {
	return string(c.StackTrace)
}

func (c wCaller) Source() string {
	return string(c.SourceFragment)
}
