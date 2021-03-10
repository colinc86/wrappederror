package wrappederror

import (
	"errors"
	"testing"
)

func TestNewJSONErrorOrWError(t *testing.T) {
	packageState.config.SetMarshalMinimalJSON(true)

	var err error
	err = errors.New("test")
	t.Run("New JSON error or wError 0", func(t *testing.T) {
		testNewJSONErrorOrWError(t, err)
	})

	err = New(nil, "test")
	t.Run("New JSON error or wError 1", func(t *testing.T) {
		testNewJSONErrorOrWError(t, err)
	})

	err = New(err, "test")
	t.Run("New JSON error or wError 2", func(t *testing.T) {
		testNewJSONErrorOrWError(t, err)
	})
}

func testNewJSONErrorOrWError(t *testing.T, err error) {
	e := newJSONErrorOrWError(err)
	if _, ok := err.(*Error); ok {
		if _, ok := e.(*jsonWErrorMinimal); !ok {
			t.Error("Unexpected JSON wError type.")
		}
	} else if _, ok := e.(*jsonError); !ok {
		t.Error("Unexpected JSON error type.")
	}
}
