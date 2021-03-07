package wrappederror

import (
	"errors"
	"testing"
)

func TestCurrentMetadata(t *testing.T) {
	SetNextErrorIndex(1)
	m1 := currentMetadata(nil)
	m2 := currentMetadata(nil)

	if m1.index != 1 {
		t.Errorf("Expected starting index 1 but received: %d\n", m1.index)
	}

	if m1.index+1 != m2.index {
		t.Errorf("Expected m2.index (%d) to be one greater than m1.index (%d).\n", m2.index, m1.index)
	}
}

func TestSimilarMetadata(t *testing.T) {
	e1 := errors.New("test error")
	e2 := errors.New("test error")
	e3 := errors.New("test error")
	e4 := errors.New("testerror")
	e5 := errors.New("testerror")

	_ = currentMetadata(e1)
	_ = currentMetadata(e2)
	m1 := currentMetadata(e3)
	_ = currentMetadata(e4)
	m2 := currentMetadata(e5)
	_ = currentMetadata(nil)
	m3 := currentMetadata(nil)

	if m1.similarErrors != 2 {
		t.Errorf("Expected 2 similar errors but received %d.\n", m1.similarErrors)
	}

	if m2.similarErrors != 1 {
		t.Errorf("Expected 1 similar error but received %d.\n", m2.similarErrors)
	}

	if m3.similarErrors != 0 {
		t.Errorf("Expected no similar errors but received %d.\n", m3.similarErrors)
	}
}
