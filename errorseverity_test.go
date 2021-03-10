package wrappederror

import (
	"errors"
	"testing"
)

var testErrorSeverities struct {
	es0 *ErrorSeverity
	es1 *ErrorSeverity
	es2 *ErrorSeverity
	es3 *ErrorSeverity
}

func setupErrorSeverityTests() {
	testErrorSeverities.es0, _ = NewErrorSeverity(
		"es0",
		"abc",
		ErrorSeverityLevelLow,
	)
	testErrorSeverities.es1, _ = NewErrorSeverity(
		"es1",
		"abc",
		ErrorSeverityLevelLow,
	)
	testErrorSeverities.es2, _ = NewErrorSeverity(
		"es2",
		"abcde",
		ErrorSeverityLevelHigh,
	)
	testErrorSeverities.es3, _ = NewErrorSeverity(
		"es3",
		"abab",
		ErrorSeverityLevelHigh,
	)
}

// Tests

func TestNewErrorSeverity(t *testing.T) {
	t.Run("New error severity 0", func(t *testing.T) {
		testNewErrorSeverity(t, "es0", "abc", ErrorSeverityLevelLow, false)
	})
	t.Run("New error severity 1", func(t *testing.T) {
		testNewErrorSeverity(t, "es1", "\\", ErrorSeverityLevelLow, true)
	})
	t.Run("New error severity 2", func(t *testing.T) {
		testNewErrorSeverity(t, "es2", "", ErrorSeverityLevelLow, true)
	})
}

func testNewErrorSeverity(
	t *testing.T,
	s, r string,
	l ErrorSeverityLevel,
	e bool,
) {
	es, err := NewErrorSeverity(s, r, l)
	if e {
		if es != nil {
			t.Errorf("Expected a nil error severity but received %s.\n", es)
		}
		if err == nil {
			t.Error("Expected error.")
		}
	} else {
		if es == nil {
			t.Error("Expected an error severity.\n")
		}
		if err != nil {
			t.Errorf("Expected a nil error but received %s.\n", err)
		}
	}
}

func TestErrorSeverityString(t *testing.T) {
	// Sanity check
	if len(testErrorSeverities.es0.String()) == 0 {
		t.Error("Unexepected string length.")
	}
}

func TestErrorSeverityMatch(t *testing.T) {
	if err := RegisterErrorSeverity(testErrorSeverities.es0); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e0 := errors.New("abcdef")
	t.Run("Error severity match 0", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es0, e0, 0.5)
	})

	e1 := errors.New("defg")
	t.Run("Error severity match 1", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es0, e1, 0.0)
	})

	if err := RegisterErrorSeverity(testErrorSeverities.es2); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e2 := errors.New("abcdefg")
	t.Run("Error severity match 2", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es2, e2, 5.0/7.0)
	})

	e3 := errors.New("abcd")
	t.Run("Error severity match 3", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es2, e3, 0.0)
	})

	e4 := errors.New("")
	t.Run("Error severity match 4", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es2, e4, 0.0)
	})

	if err := RegisterErrorSeverity(testErrorSeverities.es3); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e5 := errors.New("abababababab")
	t.Run("Error severity match 5", func(t *testing.T) {
		testErrorSeverityMatch(t, testErrorSeverities.es3, e5, 1.0)
	})

	UnregisterErrorSeverity(testErrorSeverities.es0)
	UnregisterErrorSeverity(testErrorSeverities.es2)
	UnregisterErrorSeverity(testErrorSeverities.es3)
}

func testErrorSeverityMatch(t *testing.T, es *ErrorSeverity, err error, ex float64) {
	if es.match(err) != ex {
		t.Errorf("Expected %f but received %f.\n", ex, es.match(err))
	}
}

func TestErrorSeverityEquals(t *testing.T) {
	t.Run("Error severity equals 0", func(t *testing.T) {
		testErrorSeverityEquals(
			t,
			testErrorSeverities.es0,
			testErrorSeverities.es1,
			true,
		)
	})
	t.Run("Error severity equals 1", func(t *testing.T) {
		testErrorSeverityEquals(
			t,
			testErrorSeverities.es0,
			testErrorSeverities.es2,
			false,
		)
	})
	t.Run("Error severity equals 2", func(t *testing.T) {
		testErrorSeverityEquals(
			t,
			testErrorSeverities.es1,
			testErrorSeverities.es2,
			false,
		)
	})
}

func testErrorSeverityEquals(
	t *testing.T,
	es1, es2 *ErrorSeverity,
	ex bool,
) {
	if es1.equals(es2) && !ex {
		t.Errorf("Expected %t.\n", ex)
	}
	if es2.equals(es1) && !es1.equals(es2) {
		t.Error("Unexpected non-commutativity.")
	}
}
