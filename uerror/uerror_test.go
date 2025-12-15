package uerror

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("test error")
	if err.Code != 0 {
		t.Errorf("expected code 0, got %d", err.Code)
	}
	if err.Message != "test error" {
		t.Errorf("expected message 'test error', got '%s'", err.Message)
	}
}

func TestNewWithCode(t *testing.T) {
	err := NewWithCode(404, "not found")
	if err.Code != 404 {
		t.Errorf("expected code 404, got %d", err.Code)
	}
	if err.Message != "not found" {
		t.Errorf("expected message 'not found', got '%s'", err.Message)
	}
}

func TestWrap(t *testing.T) {
	orig := errors.New("original error")
	err := Wrap(orig, "wrapped")
	if err.Err != orig {
		t.Error("expected underlying error to match original")
	}
	if err.Message != "wrapped" {
		t.Errorf("expected message 'wrapped', got '%s'", err.Message)
	}
}

func TestErrorString(t *testing.T) {
	err := NewWithCode(500, "internal error")
	expected := "[500] internal error"
	if err.Error() != expected {
		t.Errorf("expected string '%s', got '%s'", expected, err.Error())
	}

	orig := errors.New("db fail")
	wrapped := WrapWithCode(500, orig, "query failed")
	expectedWrapped := "[500] query failed: db fail"
	if wrapped.Error() != expectedWrapped {
		t.Errorf("expected string '%s', got '%s'", expectedWrapped, wrapped.Error())
	}
}

func TestFormat(t *testing.T) {
	orig := errors.New("root cause")
	err := WrapWithCode(500, orig, "something went wrong")

	s := fmt.Sprintf("%+v", err)
	// Base expectation part of the string
	expectedStart := "Code: 500\nMessage: something went wrong\nCause: root cause\n"

	if len(s) < len(expectedStart) || s[:len(expectedStart)] != expectedStart {
		t.Errorf("expected formatted string to start with:\n%q\ngot:\n%q", expectedStart, s)
	}

	// Check if stack trace is present by looking for the current file name
	if !contains(s, "uerror_test.go") {
		t.Error("expected stack trace to capture testing file 'uerror_test.go'")
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
