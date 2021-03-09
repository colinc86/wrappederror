package wrappederror

import "fmt"

// ErrorFormatToken types specify a token that can be used to format an error
// string.
type ErrorFormatToken string

// A group of error format tokens.
const (
	// errorFormatTokenNone is used when no other token match is found.
	errorFormatTokenNone ErrorFormatToken = ""

	// ErrorFormatTokenContext prints the error's context.
	ErrorFormatTokenContext ErrorFormatToken = "${{CTX}}"

	// ErrorFormatTokenInner prints the error's context.
	ErrorFormatTokenInner ErrorFormatToken = "${{INN}}"

	// ErrorFormatTokenChain prints the error's context.
	ErrorFormatTokenChain ErrorFormatToken = "${{CHN}}"

	// ErrorFormatTokenFile prints the error's context.
	ErrorFormatTokenFile ErrorFormatToken = "${{FIL}}"

	// ErrorFormatTokenFunction prints the error's context.
	ErrorFormatTokenFunction ErrorFormatToken = "${{FUN}}"

	// ErrorFormatTokenLine prints the error's context.
	ErrorFormatTokenLine ErrorFormatToken = "${{LIN}}"

	// ErrorFormatTokenStack prints the error's context.
	ErrorFormatTokenStack ErrorFormatToken = "${{STK}}"

	// ErrorFormatTokenSource prints the error's context.
	ErrorFormatTokenSource ErrorFormatToken = "${{SRC}}"

	// ErrorFormatTokenTime prints the error's context.
	ErrorFormatTokenTime ErrorFormatToken = "${{TIM}}"

	// ErrorFormatTokenDuration prints the error's context.
	ErrorFormatTokenDuration ErrorFormatToken = "${{DUR}}"

	// ErrorFormatTokenIndex prints the error's context.
	ErrorFormatTokenIndex ErrorFormatToken = "${{IDX}}"

	// ErrorFormatTokenSimilar prints the error's context.
	ErrorFormatTokenSimilar ErrorFormatToken = "${{SIM}}"

	// ErrorFormatTokenRoutines prints the error's context.
	ErrorFormatTokenRoutines ErrorFormatToken = "${{RTS}}"

	// ErrorFormatTokenCPUs prints the error's context.
	ErrorFormatTokenCPUs ErrorFormatToken = "${{CPU}}"

	// ErrorFormatTokenCGO prints the error's context.
	ErrorFormatTokenCGO ErrorFormatToken = "${{CGO}}"

	// ErrorFormatTokenMemory prints the error's context.
	ErrorFormatTokenMemory ErrorFormatToken = "${{MEM}}"

	// ErrorFormatTokenSeverityTitle prints the error's severity title.
	ErrorFormatTokenSeverityTitle ErrorFormatToken = "${{SET}}"

	// ErrorFormatTokenSeverityLevel prints the error's severity level.
	ErrorFormatTokenSeverityLevel ErrorFormatToken = "${{SEL}}"
)

const (
	// The leading substring to match against.
	tokenLeadingSubstring = "${{"

	// The length of all error format tokens.
	tokenLength = 8
)

// formatter types format an error according to an error format string.
type formatter struct{}

// Initializers

// newFormatter creates and returns a new formatter.
func newFormatter() *formatter {
	return &formatter{}
}

// Methods

// format returns a formatted version of the error according to the given error
// format string.
func (f formatter) format(e Error, ef string) string {
	indexes := f.findIndexes(ef, tokenLeadingSubstring)
	format, tokens := f.replaceTokens(ef, indexes)

	var values []interface{}
	for _, token := range tokens {
		values = append(values, f.value(e, token))
	}

	return fmt.Sprintf(format, values...)
}

// findIndexes finds the indexes of the substring, s, in the error format
// string, ef.
func (f formatter) findIndexes(ef string, s string) []int {
	if len(ef) < len(s) {
		return nil
	}

	var idx []int

	for i := 0; i < len(ef)-len(s); i++ {
		if ef[i:i+len(s)] == s {
			idx = append(idx, i)
		}
	}

	return idx
}

// replaceTokens replaces the tokens at the given indexes in the error format
// string.
//
// It returns, in order, the error format tokens that were replaced.
func (f formatter) replaceTokens(
	ef string,
	idx []int,
) (string, []ErrorFormatToken) {
	if len(idx) == 0 {
		return ef, nil
	}

	efc := ef
	var tokens []ErrorFormatToken

	for j := len(idx) - 1; j >= 0; j-- {
		i := idx[j]

		if i+tokenLength > len(efc) {
			continue
		}

		var t ErrorFormatToken
		efc, t = f.replaceToken(efc, i)
		if t != errorFormatTokenNone {
			tokens = append([]ErrorFormatToken{t}, tokens...)
		}
	}

	return efc, tokens
}

// replaceToken replaces the token at the given index and returns the new error
// format string and the token that was replaced.
func (f formatter) replaceToken(ef string, idx int) (string, ErrorFormatToken) {
	if idx+8 > len(ef) {
		return ef, errorFormatTokenNone
	}

	if ef[idx:idx+3] != tokenLeadingSubstring {
		return ef, errorFormatTokenNone
	}

	formatToken, verb := f.newFormat(ef[idx : idx+tokenLength])
	efc := ef[:idx] + verb + ef[idx+tokenLength:]
	return efc, formatToken
}

// newFormat returns the token as an error format token, and its format verb.
func (f formatter) newFormat(t string) (ErrorFormatToken, string) {
	switch ErrorFormatToken(t) {
	case ErrorFormatTokenContext:
		return ErrorFormatTokenContext, "%+v"
	case ErrorFormatTokenInner:
		return ErrorFormatTokenInner, "%+v"
	case ErrorFormatTokenChain:
		return ErrorFormatTokenChain, "%s"
	case ErrorFormatTokenFile:
		return ErrorFormatTokenFile, "%s"
	case ErrorFormatTokenFunction:
		return ErrorFormatTokenFunction, "%s"
	case ErrorFormatTokenLine:
		return ErrorFormatTokenLine, "%d"
	case ErrorFormatTokenStack:
		return ErrorFormatTokenStack, "%s"
	case ErrorFormatTokenSource:
		return ErrorFormatTokenSource, "%s"
	case ErrorFormatTokenTime:
		return ErrorFormatTokenTime, "%s"
	case ErrorFormatTokenDuration:
		return ErrorFormatTokenDuration, "%f"
	case ErrorFormatTokenIndex:
		return ErrorFormatTokenIndex, "%d"
	case ErrorFormatTokenSimilar:
		return ErrorFormatTokenSimilar, "%d"
	case ErrorFormatTokenRoutines:
		return ErrorFormatTokenRoutines, "%d"
	case ErrorFormatTokenCPUs:
		return ErrorFormatTokenCPUs, "%d"
	case ErrorFormatTokenCGO:
		return ErrorFormatTokenCGO, "%d"
	case ErrorFormatTokenMemory:
		return ErrorFormatTokenMemory, "%+v"
	case ErrorFormatTokenSeverityTitle:
		return ErrorFormatTokenSeverityTitle, "%s"
	case ErrorFormatTokenSeverityLevel:
		return ErrorFormatTokenSeverityLevel, "%s"
	default:
		return errorFormatTokenNone, ""
	}
}

// value gets the value of the error for the given error format token.
func (f formatter) value(e Error, t ErrorFormatToken) interface{} {
	switch ErrorFormatToken(t) {
	case ErrorFormatTokenContext:
		return e.Context
	case ErrorFormatTokenInner:
		return e.inner.Error()
	case ErrorFormatTokenChain:
		return e.Error()
	case ErrorFormatTokenFile:
		if e.Caller == nil {
			return callerFileNameUnknown
		}
		return e.Caller.File
	case ErrorFormatTokenFunction:
		if e.Caller == nil {
			return callerFunctionNameUnknown
		}
		return e.Caller.Function
	case ErrorFormatTokenLine:
		if e.Caller == nil {
			return callerLineNumberUnknown
		}
		return e.Caller.Line
	case ErrorFormatTokenStack:
		if e.Caller == nil {
			return ""
		}
		return e.Caller.StackTrace
	case ErrorFormatTokenSource:
		if e.Caller == nil {
			return ""
		}
		if e.Caller.Fragment == nil {
			return ""
		}
		return e.Caller.Fragment.Source
	case ErrorFormatTokenTime:
		return e.Metadata.Time
	case ErrorFormatTokenDuration:
		return e.Metadata.Duration.Seconds()
	case ErrorFormatTokenIndex:
		return e.Metadata.Index
	case ErrorFormatTokenSimilar:
		return e.Metadata.Similar
	case ErrorFormatTokenRoutines:
		if e.Process == nil {
			return processRoutinesNumberUnknown
		}
		return e.Process.Routines
	case ErrorFormatTokenCPUs:
		if e.Process == nil {
			return processCPUsNumberUnknown
		}
		return e.Process.CPUs
	case ErrorFormatTokenCGO:
		if e.Process == nil {
			return processCGONumberUnknown
		}
		return e.Process.CGO
	case ErrorFormatTokenMemory:
		if e.Process == nil {
			return nil
		}
		return e.Process.Memory
	case ErrorFormatTokenSeverityTitle:
		if e.Metadata.Severity == nil {
			return ""
		}
		return e.Metadata.Severity.Title
	case ErrorFormatTokenSeverityLevel:
		if e.Metadata.Severity == nil {
			return ""
		}
		return e.Metadata.Severity.Level
	default:
		return nil
	}
}
