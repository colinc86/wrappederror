package wrappederror

import (
	"errors"
	"fmt"
	"regexp"
)

// ErrRegexRequired indicates that a regular expression is required.
var ErrRegexRequired = errors.New("regex required")

// ErrorSeverityLevel types define an error severity level.
type ErrorSeverityLevel string

// A group of error severity levels.
const (
	ErrorSeverityLevelNone     ErrorSeverityLevel = "none"
	ErrorSeverityLevelLow      ErrorSeverityLevel = "low"
	ErrorSeverityLevelModerate ErrorSeverityLevel = "moderate"
	ErrorSeverityLevelHigh     ErrorSeverityLevel = "high"
	ErrorSeverityLevelSevere   ErrorSeverityLevel = "severe"
)

// ErrorSeverity types define an error severity with a title, level, and a
// regular expression that is used to find matches in an error's `Error` method
// output.
type ErrorSeverity struct {

	// The severity's title.
	Title string `json:"title"`

	// The regular expression used to match against `Error` method strings.
	Regex *regexp.Regexp `json:"regex"`

	// The severity level.
	Level ErrorSeverityLevel `json:"level"`
}

// Initializers

// NewErrorSeverity creates and returns a new error severity with the given
// title, level and regular expression. If the regular expression is invalid,
// this function will return an error and an unknown error severity.
func NewErrorSeverity(
	title string,
	regex string,
	level ErrorSeverityLevel,
) (*ErrorSeverity, error) {
	if len(regex) == 0 {
		return nil, ErrRegexRequired
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	return &ErrorSeverity{
		Title: title,
		Regex: r,
		Level: level,
	}, nil
}

// Stringer interface methods

func (s ErrorSeverity) String() string {
	return fmt.Sprintf("[%s] %s", s.Level, s.Title)
}

// Non-exported methods

// match matches the error against the error severity.
func (s ErrorSeverity) match(err error) float64 {
	es := err.Error()
	if len(es) == 0 {
		return 0.0
	}

	matches := s.Regex.FindAllStringIndex(es, -1)
	cl := 0

	for _, m := range matches {
		cl += m[1] - m[0]
	}

	return float64(cl) / float64(len(es))
}

// equals returns whether or not the receiver is equal to severity. Two error
// severities are considered equal if their regular expressesions are equal.
func (s ErrorSeverity) equals(severity *ErrorSeverity) bool {
	return s.Regex.String() == severity.Regex.String()
}
