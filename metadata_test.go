package wrappederror

import (
	"errors"
	"fmt"
	"testing"
)

func TestCurrentMetadata(t *testing.T) {
	packageState.config.SetNextErrorIndex(1)
	m1 := newMetadata(nil)
	m2 := newMetadata(nil)

	if m1.Index != 1 {
		t.Errorf("Expected starting index 1 but received: %d\n", m1.Index)
	}

	if m1.Index+1 != m2.Index {
		t.Errorf("Expected m2.index (%d) to be one greater than m1.index (%d).\n", m2.Index, m1.Index)
	}

	fmt.Println(m2)
}

func TestSimilarMetadata(t *testing.T) {
	e1 := errors.New("test error")
	e2 := errors.New("test error")
	e3 := errors.New("test error")
	e4 := errors.New("testerror")
	e5 := errors.New("testerror")

	_ = newMetadata(e1)
	_ = newMetadata(e2)
	m1 := newMetadata(e3)
	_ = newMetadata(e4)
	m2 := newMetadata(e5)
	_ = newMetadata(nil)
	m3 := newMetadata(nil)

	if m1.Similar != 2 {
		t.Errorf("Expected 2 similar errors but received %d.\n", m1.Similar)
	}

	if m2.Similar != 1 {
		t.Errorf("Expected 1 similar error but received %d.\n", m2.Similar)
	}

	if m3.Similar != 0 {
		t.Errorf("Expected no similar errors but received %d.\n", m3.Similar)
	}
}
