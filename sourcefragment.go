package wrappederror

import (
	"bufio"
	"fmt"
	"os"
)

// SourceFragment types store information about a source fragment such as the
// file, line range and source code.
type SourceFragment struct {

	// The source fragment's file.
	File string `json:"file"`

	// The lower line index of the lines stored in Source.
	LowerLine int `json:"lowerLine"`

	// The upper line index of the lines stored in Source.
	UpperLine int `json:"upperLine"`

	// The fragment's source code extracted from the file with path File from
	// LowerLine through UpperLine.
	Source string `json:"source"`
}

// newSourceFragment creates and returns a new source fragment from the file
// with the given filePath at lineNumber with radius. Only valid lines in the
// file are scanned.
//
// For example, asking for the source fragment with radius 5 around the line -10
// will create a source fragment with lower and upper lines set to 0 and an
// empty Source value.
func newSourceFragment(
	filePath string,
	lineNumber int,
	radius int,
) (*SourceFragment, error) {
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

	return &SourceFragment{
		File:      filePath,
		LowerLine: ali,
		UpperLine: aui,
		Source:    string(b),
	}, nil
}

// Stringer interface methods

func (f SourceFragment) String() string {
	return fmt.Sprintf(
		"[%d - %d] %s\n%s",
		f.LowerLine,
		f.UpperLine,
		f.File,
		f.Source,
	)
}
