package wrappederror

import (
	"errors"
	"testing"
)

func TestNewErrorMap(t *testing.T) {
	m := newErrorMap()
	if m.hashMap == nil {
		t.Error("Found unexpected nil.")
	}
}

func TestErrorMapSimilarErrors(t *testing.T) {
	m := newErrorMap()

	e1 := errors.New("test")
	t.Run("Error map similar 0", func(t *testing.T) {
		testErrorMapSimilarErrors(t, m, e1, 0)
	})
	m.addError(e1)
	t.Run("Error map similar 1", func(t *testing.T) {
		testErrorMapSimilarErrors(t, m, e1, 1)
	})
	m.addError(e1)
	t.Run("Error map similar 2", func(t *testing.T) {
		testErrorMapSimilarErrors(t, m, e1, 2)
	})

	e2 := errors.New("testt")
	t.Run("Error map similar 3", func(t *testing.T) {
		testErrorMapSimilarErrors(t, m, e2, 0)
	})
}

func testErrorMapSimilarErrors(t *testing.T, m *errorMap, err error, i int) {
	if m.similarErrors(err) != i {
		t.Errorf("Expected %d but received %d.\n", i, m.similarErrors(err))
	}
}
