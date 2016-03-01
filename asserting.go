// Package asserting provides a testing framework built on top of the "testing"
// package. The goal is to make testing a little bit easier, reducing the amount
// of boilerplate code.
package asserting

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// TestCase describes a test case with various assertion methods.
type TestCase struct {
	T *testing.T
}

// Run expects a type that extends TestCase and calls all methods with prefix
// "Test".
func Run(i interface{}) {
	value := reflect.ValueOf(i)
	testType := reflect.TypeOf(i)
	for i := 0; i < testType.NumMethod(); i++ {
		method := testType.Method(i)
		if strings.HasPrefix(method.Name, "Test") {
			finalMethod := value.MethodByName(method.Name)
			finalMethod.Call(nil)
		}
	}
}

// Assert tests v's truthiness.
func (t TestCase) Assert(v bool) {
	if !v {
		t.T.Fatalf("Expected true, got false")
	}
}

// AssertFalse tests v's falseness.
func (t TestCase) AssertFalse(v bool) {
	if v {
		t.T.Fatalf("Expected false, got true")
	}
}

// Assertf tests v's truthiness, shows msg in case of failure.
func (t TestCase) Assertf(ok bool, msg string) {
	if !ok {
		t.T.Fatalf("Assertion failed: %s", msg)
	}
}

// AssertOK tests for HTTP OK code.
func (t TestCase) AssertOK(code int) {
	if code != http.StatusOK {
		t.T.Fatalf("Expected 200, got %d", code)
	}
}

// AssertCreated tests for HTTP CREATED code.
func (t TestCase) AssertCreated(code int) {
	if code != http.StatusCreated {
		t.T.Fatalf("Expected 201, got %d", code)
	}
}
