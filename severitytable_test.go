package wrappederror

import (
	"errors"
	"testing"
)

func TestSeverityTableRegister(t *testing.T) {
	st := newSeverityTable()
	if len(st.severities) != 0 {
		t.Errorf("Unexpected length %d.\n", len(st.severities))
	}

	t.Run("Severity table register 0", func(t *testing.T) {
		testSeverityTableRegister(t, st, testErrorSeverities.es0, nil)
	})

	if len(st.severities) == 0 {
		t.Errorf("Unexpected length %d.\n", len(st.severities))
	}

	t.Run("Severity table register 1", func(t *testing.T) {
		testSeverityTableRegister(
			t,
			st,
			testErrorSeverities.es0,
			ErrSeverityAlreadyRegistered,
		)
	})
}

func testSeverityTableRegister(
	t *testing.T,
	st *severityTable,
	s *ErrorSeverity,
	ex error,
) {
	o := st.register(s)
	if o != ex {
		t.Errorf("Expected error %+v but received %+v.\n", ex, o)
	}
}

func TestSeverityTableUnegister(t *testing.T) {
	st := newSeverityTable()
	_ = st.register(testErrorSeverities.es1)

	if len(st.severities) != 1 {
		t.Errorf("Unexpected length %d.\n", len(st.severities))
	}

	st.unregister(testErrorSeverities.es1)

	if len(st.severities) != 0 {
		t.Errorf("Unexpected length %d.\n", len(st.severities))
	}
}

func TestSeverityTableBestMatch(t *testing.T) {
	st := newSeverityTable()
	_ = st.register(testErrorSeverities.es1)
	_ = st.register(testErrorSeverities.es2)

	es := st.bestMatch(errors.New("abc"))
	if !es.equals(testErrorSeverities.es1) {
		t.Errorf("Unexpected severity %s.\n", es)
	}

	es = st.bestMatch(errors.New("abcde"))
	if !es.equals(testErrorSeverities.es2) {
		t.Errorf("Unexpected severity %s.\n", es)
	}
}
