package wrappederror

import (
	"fmt"
	"path"
	"runtime"
)

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

// BinaryMarshaler and BinaryUnmarshaler interface methods

func (c caller) MarshalBinary() ([]byte, error) {
	e := newEncoder()
	e.encodeCaller(&c)
	e.calculateCRC()
	return e.data, nil
}

func (c *caller) UnmarshalBinary(b []byte) error {
	d := newDecoder(b)
	if !d.validate() {
		return ErrCRC
	}

	ca, err := d.decodeCaller()
	if err != nil {
		return err
	}

	c.fileName = ca.fileName
	c.functionName = ca.functionName
	c.lineNumber = ca.lineNumber
	return nil
}
