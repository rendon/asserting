package yatf

import (
	"net/http"
	"testing"
)

type TestCase struct {
	T *testing.T
}

func (t TestCase) Assert(v bool) {
	if !v {
		t.T.Fatalf("Expected true, got false")
	}
}

func (t TestCase) AssertFalse(v bool) {
	if v {
		t.T.Fatalf("Expected false, got true")
	}
}

func (t TestCase) Assertf(ok bool, msg string) {
	if !ok {
		t.T.Fatalf("Assertion failed: %s", msg)
	}
}

func (t TestCase) AssertOK(code int) {
	if code != http.StatusOK {
		t.T.Fatalf("Expected 200, got %d", code)
	}
}

func (t TestCase) AssertCreated(code int) {
	if code != http.StatusCreated {
		t.T.Fatalf("Expected 201, got %d", code)
	}
}
