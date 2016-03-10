// Package asserting provides a testing framework built on top of the "testing"
// package. The goal is to make testing a little bit easier, reducing the amount
// of boilerplate code.
package asserting

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

// TestCase describes a test case with various assertion methods.
type TestCase struct {
	ResponseBody []byte
	T            *testing.T
	err          error
	server       *httptest.Server
	response     *http.Response
}

// NewTestCase returns an initialized TestCase.
func NewTestCase(t *testing.T, handlers *mux.Router) *TestCase {
	return &TestCase{
		T:      t,
		server: httptest.NewServer(handlers),
	}
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

// AssertNoError tests if err is not nil, i.e., no error has occurred.
func (t TestCase) AssertNoError(err error) {
	if err != nil {
		t.T.Fatalf("Unexpected error: %s", err)
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
func (t TestCase) AssertOK() {
	if t.err != nil {
		t.T.Fatalf("Request error is not nil: %s", t.err)
	}
	if t.response == nil {
		t.T.Fatalf("Response is nil")
	}
	if t.response.StatusCode != http.StatusOK {
		t.T.Fatalf("Expected 200, got %d", t.response.StatusCode)
	}
}

// AssertCreated tests for HTTP CREATED code.
func (t TestCase) AssertCreated(code int) {
	if code != http.StatusCreated {
		t.T.Fatalf("Expected 201, got %d", code)
	}
}

// AssertStatus tests for some specific HTTP response code against the previous
// HTTP request.
func (t TestCase) AssertStatus(code int) {
	if t.response.StatusCode != code {
		t.T.Fatalf("Expected %v, got %d", code, t.response.StatusCode)
	}
}

// Get issues an HTTP GET request and keeps the response for later assertions.
func (t *TestCase) Get(url string) {
	if t.server == nil {
		t.T.Fatalf("Uninitialized test server")
	}
	url = t.server.URL + url
	resp, err := http.Get(url)
	t.response = resp
	t.err = err
	if err == nil {
		defer t.response.Body.Close()
		t.ResponseBody, t.err = ioutil.ReadAll(t.response.Body)
	}
}

// Post issues an HTTP POST request and keeps the response for later assertions.
func (t *TestCase) Post(url string, contentType string, body []byte) {
	if t.server == nil {
		t.T.Fatalf("Uninitialized test server")
	}
	url = t.server.URL + url
	resp, err := http.Post(url, contentType, bytes.NewReader(body))
	t.response = resp
	t.err = err
	if err == nil {
		defer t.response.Body.Close()
		t.ResponseBody, t.err = ioutil.ReadAll(t.response.Body)
	}
}

// Unmarshal unmarshals response  body into and store it into  i, the test fails
// if some error occurs.
func (t *TestCase) Unmarshal(i interface{}) {
	if t.response == nil {
		t.T.Fatalf("Response is nil")
	}
	if err := json.Unmarshal(t.ResponseBody, i); err != nil {
		t.T.Fatalf("Failed to unmarshal response body data: %s", err)
	}
}

// AssertEqualInt tests if expected is equal to actual.
func (t *TestCase) AssertEqualInt(expected, actual int) {
	if expected != actual {
		t.T.Fatalf("Expected %d, got %d", expected, actual)
	}
}

// AssertEqualStr tests if expected is equal to actual.
func (t *TestCase) AssertEqualStr(expected, actual string) {
	if expected != actual {
		t.T.Fatalf("Expected %s, got %s", expected, actual)
	}
}
