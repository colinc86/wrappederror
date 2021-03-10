package wrappederror

import (
	"errors"
	"testing"
)

func TestNewErrorSeverity(t *testing.T) {
	_, err := NewErrorSeverity("Test Severity", "abc", ErrorSeverityLevelLow)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
}

func TestErrorSeverityCreate_2(t *testing.T) {
	_, err := NewErrorSeverity("Test Severity", "\\", ErrorSeverityLevelLow)
	if err == nil {
		t.Errorf("Expected error.")
	}
}

func TestErrorSeverityRegister(t *testing.T) {
	s, _ := NewErrorSeverity("Test Severity", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	if err := RegisterErrorSeverity(s); err == nil {
		t.Error("Expected error.")
	}

	UnregisterErrorSeverity(s)
}

func TestErrorSeverityUnregister(t *testing.T) {
	s, _ := NewErrorSeverity("Test Severity", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	UnregisterErrorSeverity(s)
	UnregisterErrorSeverity(s)
}

func TestErrorSeverityMatch_1(t *testing.T) {
	s, _ := NewErrorSeverity("Test Severity", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e := New(errors.New("abcdefg"), "an error")
	if e.Metadata.Severity == nil {
		UnregisterErrorSeverity(s)
		t.Fatal("Severity should not be nil.")
	}

	if !s.equals(*e.Metadata.Severity) {
		t.Errorf("Expected severity %+v but received %+v\n", s, e.Metadata.Severity)
	}

	UnregisterErrorSeverity(s)
}

func TestErrorSeverityMatch_2(t *testing.T) {
	s, _ := NewErrorSeverity("Test Severity", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e := New(errors.New("defg"), "an error")
	if e.Metadata.Severity != nil {
		t.Errorf("Severity be nil but received %+v\n", e.Metadata.Severity)
	}

	UnregisterErrorSeverity(s)
}

func TestErrorSeverityMatch_3(t *testing.T) {
	s1, _ := NewErrorSeverity("Test Severity 1", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s1); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	s2, _ := NewErrorSeverity("Test Severity 2", "abcde", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s2); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e := New(errors.New("abcdefg"), "an error")
	if e.Metadata.Severity == nil {
		UnregisterErrorSeverity(s1)
		UnregisterErrorSeverity(s2)
		t.Fatal("Severity should not be nil.")
	}

	if !e.Metadata.Severity.equals(s2) {
		t.Errorf("Expected severity %+v but received %+v\n", s2, e.Metadata.Severity)
	}

	UnregisterErrorSeverity(s1)
	UnregisterErrorSeverity(s2)
}

func TestErrorSeverityMatch_4(t *testing.T) {
	s1, _ := NewErrorSeverity("Test Severity 1", "abc", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s1); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	s2, _ := NewErrorSeverity("Test Severity 2", "abcde", ErrorSeverityLevelLow)
	if err := RegisterErrorSeverity(s2); err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	e1 := New(errors.New("abcdefg"), "an error")
	if e1.Metadata.Severity == nil {
		UnregisterErrorSeverity(s1)
		UnregisterErrorSeverity(s2)
		t.Fatal("Severity should not be nil.")
	}
	if !e1.Metadata.Severity.equals(s2) {
		t.Errorf("Expected severity %+v but received %+v\n", s2, e1.Metadata.Severity)
	}

	e2 := New(e1, "abcd")
	if e2.Metadata.Severity == nil {
		UnregisterErrorSeverity(s1)
		UnregisterErrorSeverity(s2)
		t.Fatal("Severity should not be nil.")
	}
	if !e2.Metadata.Severity.equals(s2) {
		t.Errorf("Expected severity %+v but received %+v\n", s2, e2.Metadata.Severity)
	}

	UnregisterErrorSeverity(s1)
	UnregisterErrorSeverity(s2)
}
