package wrappederror

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCurrentMetadata(t *testing.T) {
	packageState.configuration.SetNextErrorIndex(1)
	time.Sleep(time.Second * 10)
	m1 := currentMetadata(nil)
	m2 := currentMetadata(nil)

	if m1.ErrorIndex != 1 {
		t.Errorf("Expected starting index 1 but received: %d\n", m1.ErrorIndex)
	}

	if m1.ErrorIndex+1 != m2.ErrorIndex {
		t.Errorf("Expected m2.index (%d) to be one greater than m1.index (%d).\n", m2.ErrorIndex, m1.ErrorIndex)
	}

	fmt.Println(m2)
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

	if m1.SimilarErrors != 2 {
		t.Errorf("Expected 2 similar errors but received %d.\n", m1.SimilarErrors)
	}

	if m2.SimilarErrors != 1 {
		t.Errorf("Expected 1 similar error but received %d.\n", m2.SimilarErrors)
	}

	if m3.SimilarErrors != 0 {
		t.Errorf("Expected no similar errors but received %d.\n", m3.SimilarErrors)
	}
}
