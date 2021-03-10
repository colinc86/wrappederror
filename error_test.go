package wrappederror

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	setupCallerTests()
	setupConfigurationTests()
	setupErrorTests()
	setupErrorSeverityTests()
	os.Exit(m.Run())
}

var testErrors struct {
	e0 *Error
	e1 *Error
	e2 *Error
}

func setupErrorTests() {
	testErrors.e0 = New(nil, "error 0")
	testErrors.e1 = New(testErrors.e0, "error 1")
	testErrors.e2 = New(testErrors.e1, "error 2")
}

// Tests

func TestNewError(t *testing.T) {
	e0m := testErrors.e0.context.(string)
	ex0m := e0m
	t.Run("Error message 0", func(t *testing.T) {
		testErrorMessage(t, testErrors.e0, ex0m)
	})

	e1m := testErrors.e1.context.(string)
	ex1m := e1m + errorChainDelimiter + e0m
	t.Run("Error message 1", func(t *testing.T) {
		testErrorMessage(t, testErrors.e1, ex1m)
	})

	e2m := testErrors.e2.context.(string)
	ex2m := e2m + errorChainDelimiter + e1m + errorChainDelimiter + e0m
	t.Run("Error message 2", func(t *testing.T) {
		testErrorMessage(t, testErrors.e2, ex2m)
	})
}

func testErrorMessage(t *testing.T, e *Error, s string) {
	if e.Error() != s {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", s, e.Error())
	}
}

func TestErrorDepth(t *testing.T) {
	t.Run("Error depth 0", func(t *testing.T) {
		testErrorDepth(t, testErrors.e0, 0)
	})
	t.Run("Error depth 1", func(t *testing.T) {
		testErrorDepth(t, testErrors.e1, 1)
	})
	t.Run("Error depth 2", func(t *testing.T) {
		testErrorDepth(t, testErrors.e2, 2)
	})
}

func testErrorDepth(t *testing.T, e *Error, i int) {
	if e.Depth() != i {
		t.Errorf("Expected depth %d but received %d.\n", i, e.Depth())
	}
}

func TestErrorFormat(t *testing.T) {
	es := "error"
	ef := string(ErrorFormatTokenChain)
	e := New(nil, es)

	t.Run("Error format 0", func(t *testing.T) { testErrorFormat(t, e, ef, es) })

	es += es
	ef += ef
	t.Run("Error format 1", func(t *testing.T) { testErrorFormat(t, e, ef, es) })

	es += es
	ef += ef
	t.Run("Error format 2", func(t *testing.T) { testErrorFormat(t, e, ef, es) })
}

func testErrorFormat(t *testing.T, e *Error, ef, s string) {
	if e.Format(ef) != s {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", s, e.Format(ef))
	}
}

func TestErrorChain(t *testing.T) {
	ex0 := []*Error{testErrors.e0}
	t.Run("Error chain 0", func(t *testing.T) {
		testErrorChain(t, testErrors.e0, 1, ex0)
	})

	ex1 := []*Error{testErrors.e1, testErrors.e0}
	t.Run("Error chain 1", func(t *testing.T) {
		testErrorChain(t, testErrors.e1, 2, ex1)
	})

	ex2 := []*Error{testErrors.e2, testErrors.e1, testErrors.e0}
	t.Run("Error chain 2", func(t *testing.T) {
		testErrorChain(t, testErrors.e2, 3, ex2)
	})
}

func testErrorChain(t *testing.T, e *Error, i int, s []*Error) {
	c := e.Chain()
	if len(c) != i {
		t.Fatalf("Expected %d errors but received %d.\n", i, len(c))
	}
	for i, v := range c {
		if v.Error() != s[i].Error() {
			t.Fatalf("Found unequal errors %s and %s at index %d.\n", v, s[i], i)
		}
	}
}

func TestErrorWalk(t *testing.T) {
	t.Run("Error chain 0", func(t *testing.T) { testErrorWalk(t, testErrors.e0) })
	t.Run("Error chain 1", func(t *testing.T) { testErrorWalk(t, testErrors.e1) })
	t.Run("Error chain 2", func(t *testing.T) { testErrorWalk(t, testErrors.e2) })
}

func testErrorWalk(t *testing.T, e *Error) {
	c := e.Chain()
	ci := 0
	e.Walk(func(err error) bool {
		if c[ci].Error() != err.Error() {
			t.Errorf("Expected error \"%s\" but received \"%s\".\n", c[ci], err)
			return false
		}
		ci++
		return true
	})

	ci = 0
	e.Walk(func(err error) bool {
		ci++
		return false
	})

	if ci != 1 {
		t.Errorf("Expected to walk 1 error, but took %d steps instead.\n", ci)
	}
}

func TestErrorWithDepthBadDepth_1(t *testing.T) {
	if testErrors.e2.ErrorWithDepth(-1) != nil {
		t.Error("Expected no error.")
	}
}

func TestErrorWithDepthBadDepth_2(t *testing.T) {
	if testErrors.e2.ErrorWithDepth(testErrors.e2.Depth()+1) != nil {
		t.Error("Expected no error.")
	}
}

func TestErrorErrorWithDepth(t *testing.T) {
	t.Run("Error with depth 0", func(t *testing.T) {
		testErrorErrorWithDepth(t, testErrors.e0, 0, testErrors.e0)
	})
	t.Run("Error with depth 1", func(t *testing.T) {
		testErrorErrorWithDepth(t, testErrors.e1, 0, testErrors.e0)
	})
	t.Run("Error with depth 2", func(t *testing.T) {
		testErrorErrorWithDepth(t, testErrors.e2, 0, testErrors.e0)
	})
}

func testErrorErrorWithDepth(t *testing.T, e *Error, i int, err error) {
	ewd := e.ErrorWithDepth(i)
	if ewd.Error() != err.Error() {
		t.Errorf("Expected error \"%s\" but received \"%s\".\n", err, ewd)
	}
}

func TestErrorWithIndexBadIndex_1(t *testing.T) {
	if testErrors.e2.ErrorWithIndex(-1) != nil {
		t.Error("Expected no error.")
	}
}

func TestErrorWithIndexBadIndex_2(t *testing.T) {
	if testErrors.e2.ErrorWithIndex(testErrors.e2.Depth()+1) != nil {
		t.Error("Expected no error.")
	}
}

func TestErrorErrorWithIndex(t *testing.T) {
	t.Run("Error with index 0", func(t *testing.T) {
		testErrorErrorWithIndex(t, testErrors.e0, 0, testErrors.e0)
	})
	t.Run("Error with index 1", func(t *testing.T) {
		testErrorErrorWithIndex(t, testErrors.e1, 0, testErrors.e1)
	})
	t.Run("Error with index 2", func(t *testing.T) {
		testErrorErrorWithIndex(t, testErrors.e2, 0, testErrors.e2)
	})
}

func testErrorErrorWithIndex(t *testing.T, e *Error, i int, err error) {
	ewd := e.ErrorWithIndex(i)
	if ewd.Error() != err.Error() {
		t.Errorf("Expected error \"%s\" but received \"%s\".\n", err, ewd)
	}
}

func TestErrorTrace(t *testing.T) {
	t.Run("Error trace 0", func(t *testing.T) {
		testErrorTrace(t, testErrors.e0)
	})
	t.Run("Error trace 1", func(t *testing.T) {
		testErrorTrace(t, testErrors.e1)
	})
	t.Run("Error trace 2", func(t *testing.T) {
		testErrorTrace(t, testErrors.e2)
	})

	e0 := errors.New("error 0")
	e1 := New(e0, "error 1")
	t.Run("Error trace 3", func(t *testing.T) {
		testErrorTrace(t, e1)
	})
}

func testErrorTrace(t *testing.T, e *Error) {
	d := e.Depth()
	lines := strings.Split(e.Trace(), "\n")

	if len(lines) != d+1 {
		t.Errorf("Expected %d lines but received %d.\n", d, len(lines))
	}

	for i, l := range lines {
		if len(l) < 3 {
			t.Errorf("Malformed trace line with length %d.\n", len(l))
			break
		}

		dec := ""
		if i == 0 {
			if d > 0 {
				dec = errorTraceFirstItemDecoration
			}
		} else if i == e.Depth() {
			dec = errorTraceLastItemDecoration
		} else {
			dec = errorTraceMiddleItemDecoration
		}

		p := ""
		if len(dec) > 0 {
			p = fmt.Sprintf("%s %d:", dec, d)
		}

		if !strings.HasPrefix(l, p) {
			t.Errorf("Epected prefix \"%s\" in line \"%s\".\n", p, l)
		}

		d--
	}
}

func TestErrorUnwrap(t *testing.T) {
	t.Run("Error unwrap 0", func(t *testing.T) {
		testErrorUnwrap(t, testErrors.e0, nil)
	})
	t.Run("Error unwrap 1", func(t *testing.T) {
		testErrorUnwrap(t, testErrors.e1, testErrors.e0)
	})
	t.Run("Error unwrap 2", func(t *testing.T) {
		testErrorUnwrap(t, testErrors.e2, testErrors.e1)
	})
}

func testErrorUnwrap(t *testing.T, e *Error, err error) {
	if err == nil {
		if e.Unwrap() != nil {
			t.Errorf("Expected nil but received \"%s\".\n", e.Unwrap())
		}
		return
	}

	if e.Unwrap().Error() != err.Error() {
		t.Errorf("Expected \"%s\" but received \"%s\".\n", err, e.Unwrap())
	}
}

func TestErrorAs(t *testing.T) {
	e := errors.New("error 0")
	t.Run("Error as 0", func(t *testing.T) {
		testErrorAs(t, testErrors.e0, e, false)
	})
	t.Run("Error as 1", func(t *testing.T) {
		testErrorAs(t, testErrors.e1, e, false)
	})
	t.Run("Error as 2", func(t *testing.T) {
		testErrorAs(t, testErrors.e2, e, false)
	})
	t.Run("Error as 3", func(t *testing.T) {
		testErrorAs(t, testErrors.e0, testErrors.e0, true)
	})
	t.Run("Error as 4", func(t *testing.T) {
		testErrorAs(t, testErrors.e1, testErrors.e1, true)
	})
	t.Run("Error as 5", func(t *testing.T) {
		testErrorAs(t, testErrors.e2, testErrors.e2, true)
	})
	ee1 := New(e, "error 1")
	ee2 := New(ee1, "error 2")
	t.Run("Error as 6", func(t *testing.T) {
		testErrorAs(t, ee2, ee1, false)
	})
}

func testErrorAs(t *testing.T, e *Error, as interface{}, ex bool) {
	if e.As(as) != ex {
		t.Errorf("Expected %t but received %t.\n", ex, e.As(as))
	}

	if _, ok := as.(*Error); !ok && ex {
		t.Error("Expected error.")
	}
}

func TestErrorIs(t *testing.T) {
	e := errors.New("error")
	t.Run("Error is 0", func(t *testing.T) {
		testErrorIs(t, testErrors.e0, e, false)
	})
	t.Run("Error is 1", func(t *testing.T) {
		testErrorIs(t, testErrors.e1, e, false)
	})
	t.Run("Error is 2", func(t *testing.T) {
		testErrorIs(t, testErrors.e2, e, false)
	})
	t.Run("Error is 3", func(t *testing.T) {
		testErrorIs(t, testErrors.e0, testErrors.e0, true)
	})
	t.Run("Error is 4", func(t *testing.T) {
		testErrorIs(t, testErrors.e1, testErrors.e1, true)
	})
	t.Run("Error is 5", func(t *testing.T) {
		testErrorIs(t, testErrors.e2, testErrors.e2, true)
	})
}

func testErrorIs(t *testing.T, e *Error, is error, ex bool) {
	if e.Is(is) != ex {
		t.Errorf("Expected %t but received %t.\n", ex, e.Is(is))
	}
}

func TestErrorContext(t *testing.T) {
	t.Run("Error context 0", func(t *testing.T) {
		testErrorContext(t, testErrors.e0, "error 0")
	})
	t.Run("Error context 1", func(t *testing.T) {
		testErrorContext(t, testErrors.e1, "error 1")
	})
	t.Run("Error context 2", func(t *testing.T) {
		testErrorContext(t, testErrors.e2, "error 2")
	})
}

func testErrorContext(t *testing.T, e *Error, c interface{}) {
	if e.Context() != c {
		t.Errorf("Expected %+v but received %v.\n", c, e.Context())
	}
}

func TestErrorMarshalJSONMinimal(t *testing.T) {
	packageState.config.SetMarshalMinimalJSON(true)
	_, err := json.Marshal(testErrors.e2)
	if err != nil {
		t.Fatalf("Error marshaling json: %s\n", err)
	}
}

func TestErrorMarshalJSONFull(t *testing.T) {
	packageState.config.SetMarshalMinimalJSON(false)
	_, err := json.Marshal(testErrors.e2)
	if err != nil {
		t.Fatalf("Error marshaling json: %s\n", err)
	}
}

// Benchmarks

func BenchmarkNewError_Defaults(b *testing.B) {
	packageState.reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nil, "")
	}
}

func BenchmarkNewError_NoCaller(b *testing.B) {
	packageState.reset()
	packageState.config.SetCaptureCaller(false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nil, "")
	}
}

func BenchmarkNewError_NoProcess(b *testing.B) {
	packageState.reset()
	packageState.config.SetCaptureProcess(false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nil, "")
	}
}

func BenchmarkNewError_NoCallerNoProcess(b *testing.B) {
	packageState.reset()
	packageState.config.SetCaptureCaller(false)
	packageState.config.SetCaptureProcess(false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nil, "")
	}
}

func BenchmarkNewError_NoFeatures(b *testing.B) {
	packageState.config.Set(false, false, false, true, false, true, 0, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = New(nil, "")
	}
}
