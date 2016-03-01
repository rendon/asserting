package asserting

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type TestCase struct {
	T *testing.T
}

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
