package wrappederror

import (
	"fmt"
	"testing"
)

var testFormatter = newFormatter()

// Tests

func TestFormatterFormat(t *testing.T) {
	ef := ""
	ex := ""
	t.Run("Formatter format 0", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})

	ef = fmt.Sprintf("%s", ErrorFormatTokenChain)
	ex = testErrors.e1.Error()
	t.Run("Formatter format 1", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})

	ef = fmt.Sprintf(" %s", ErrorFormatTokenChain)
	ex = " " + testErrors.e1.Error()
	t.Run("Formatter format 2", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})

	ef = fmt.Sprintf("%s ", ErrorFormatTokenChain)
	ex = testErrors.e1.Error() + " "
	t.Run("Formatter format 3", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})

	ef = fmt.Sprintf("%s%s", ErrorFormatTokenChain, ErrorFormatTokenChain)
	ex = testErrors.e1.Error() + testErrors.e1.Error()
	t.Run("Formatter format 4", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})

	ef = fmt.Sprintf("%s %s", ErrorFormatTokenChain, ErrorFormatTokenChain)
	ex = testErrors.e1.Error() + " " + testErrors.e1.Error()
	t.Run("Formatter format 5", func(t *testing.T) {
		testFormatterFormat(t, testFormatter, *testErrors.e1, ef, ex)
	})
}

func testFormatterFormat(t *testing.T, f *formatter, e Error, ef, ex string) {
	if f.format(e, ef) != ex {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", ex, f.format(e, ef))
	}
}

func TestFormatterFindIndexes(t *testing.T) {
	ef := ""
	s := ""
	ex := []int{}
	t.Run("Formatter find indexes 0", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = ""
	s = "a"
	ex = []int{}
	t.Run("Formatter find indexes 1", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "a"
	s = ""
	ex = []int{}
	t.Run("Formatter find indexes 2", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "a"
	s = "a"
	ex = []int{0}
	t.Run("Formatter find indexes 3", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "aa"
	s = "a"
	ex = []int{0, 1}
	t.Run("Formatter find indexes 4", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "aa"
	s = "aa"
	ex = []int{0}
	t.Run("Formatter find indexes 5", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "aaaa"
	s = "aa"
	ex = []int{0, 1, 2}
	t.Run("Formatter find indexes 6", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})

	ef = "aa"
	s = "aaaa"
	ex = []int{}
	t.Run("Formatter find indexes 7", func(t *testing.T) {
		testFormatterFindIndexes(t, testFormatter, ef, s, ex)
	})
}

func testFormatterFindIndexes(
	t *testing.T,
	f *formatter,
	ef, s string,
	ex []int,
) {
	o := f.findIndexes(ef, s)
	if len(o) != len(ex) {
		t.Errorf("Expected length %d but received %d.\n", len(ex), len(o))
	}
	for i, v := range o {
		if v != ex[i] {
			t.Errorf("Expected value %d but received %d.\n", ex[i], v)
		}
	}
}

func TestFormatterReplaceTokens(t *testing.T) {
	ef := ""
	var idx []int
	exs := ""
	var ext []ErrorFormatToken
	t.Run("Formatter replace tokens 0", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s", ErrorFormatTokenContext)
	idx = []int{0}
	exs = "%+v"
	ext = []ErrorFormatToken{ErrorFormatTokenContext}
	t.Run("Formatter replace tokens 1", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf(" %s", ErrorFormatTokenInner)
	idx = []int{1}
	exs = " %+v"
	ext = []ErrorFormatToken{ErrorFormatTokenInner}
	t.Run("Formatter replace tokens 1", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s ", ErrorFormatTokenChain)
	idx = []int{0}
	exs = "%s "
	ext = []ErrorFormatToken{ErrorFormatTokenChain}
	t.Run("Formatter replace tokens 2", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s%s", ErrorFormatTokenFile, ErrorFormatTokenFunction)
	idx = []int{0, 8}
	exs = "%s%s"
	ext = []ErrorFormatToken{ErrorFormatTokenFile, ErrorFormatTokenFunction}
	t.Run("Formatter replace tokens 3", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s %s", ErrorFormatTokenLine, ErrorFormatTokenTrace)
	idx = []int{0, 9}
	exs = "%d %s"
	ext = []ErrorFormatToken{ErrorFormatTokenLine, ErrorFormatTokenTrace}
	t.Run("Formatter replace tokens 4", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf(
		"%s %s",
		ErrorFormatTokenStack,
		ErrorFormatTokenSourceLowerLine,
	)
	idx = []int{0, 10}
	exs = "%s " + fmt.Sprintf("%s", ErrorFormatTokenSourceLowerLine)
	ext = []ErrorFormatToken{ErrorFormatTokenStack}
	t.Run("Formatter replace tokens 5", func(t *testing.T) {
		testFormatterReplaceTokens(t, testFormatter, ef, idx, exs, ext)
	})
}

func testFormatterReplaceTokens(
	t *testing.T,
	f *formatter,
	ef string,
	idx []int,
	exs string,
	ext []ErrorFormatToken,
) {
	os, ot := f.replaceTokens(ef, idx)
	if os != exs {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", exs, os)
	}
	if len(ot) != len(ext) {
		t.Errorf("Expected length %d but received %d.\n", len(ext), len(ot))
	}
	for i, v := range ot {
		if v != ext[i] {
			t.Errorf("Expected token \"%s\" but received \"%s\".\n", ext[i], v)
		}
	}
}

func TestFormatterReplaceToken(t *testing.T) {
	ef := ""
	idx := 0
	exs := ""
	var ext ErrorFormatToken
	t.Run("Formatter replace token 0", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s", ErrorFormatTokenSourceLowerLine)
	idx = 0
	exs = "%d"
	ext = ErrorFormatTokenSourceLowerLine
	t.Run("Formatter replace token 1", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf(" %s", ErrorFormatTokenSourceUpperLine)
	idx = 1
	exs = " %d"
	ext = ErrorFormatTokenSourceUpperLine
	t.Run("Formatter replace token 2", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s ", ErrorFormatTokenSource)
	idx = 0
	exs = "%s "
	ext = ErrorFormatTokenSource
	t.Run("Formatter replace token 3", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})

	ef = fmt.Sprintf("%s %s", ErrorFormatTokenTime, ErrorFormatTokenDuration)
	idx = 9
	exs = fmt.Sprintf("%s ", ErrorFormatTokenTime) + "%f"
	ext = ErrorFormatTokenDuration
	t.Run("Formatter replace token 4", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})

	ef = "asdfasdfasdf"
	idx = 0
	exs = "asdfasdfasdf"
	ext = ""
	t.Run("Formatter replace token 5", func(t *testing.T) {
		testFormatterReplaceToken(t, testFormatter, ef, idx, exs, ext)
	})
}

func testFormatterReplaceToken(
	t *testing.T,
	f *formatter,
	ef string,
	idx int,
	exs string,
	ext ErrorFormatToken,
) {
	os, ot := f.replaceToken(ef, idx)
	if os != exs {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", exs, os)
	}
	if ot != ext {
		t.Errorf("Expected token \"%s\" but received \"%s\".\n", ext, ot)
	}
}
