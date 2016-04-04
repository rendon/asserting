// Package asserting provides a testing framework built on top of the "testing"
// package. The goal is to make testing a little bit easier, reducing the amount
// of boilerplate code.
package asserting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// TestCase describes a test case with various assertion methods.
type TestCase struct {
	ResponseBody []byte
	T            *testing.T
	err          error
	server       *httptest.Server
	response     *http.Response
}

// NewWebTestCase returns an initialized TestCase for Web API testing.
func NewWebTestCase(t *testing.T, handlers http.Handler) *TestCase {
	return &TestCase{
		T:      t,
		server: httptest.NewServer(handlers),
	}
}

// NewTestCase returns an initialized TestCase.
func NewTestCase(t *testing.T) *TestCase {
	return &TestCase{T: t}
}

// Run expects a type that extends TestCase and calls all methods with prefix
// "Test".
func Run(i interface{}) {
	value := reflect.ValueOf(i)
	testType := reflect.TypeOf(i)
	var be bool
	var ba bool
	for i := 0; i < testType.NumMethod(); i++ {
		method := testType.Method(i)
		if strings.HasPrefix(method.Name, "BeforeEach") {
			be = true
		} else if strings.HasPrefix(method.Name, "BeforeAll") {
			ba = true
		}
	}
	if ba {
		value.MethodByName("BeforeAll").Call(nil)
	}
	fmt.Printf("===> Running tests...\n")
	for i := 0; i < testType.NumMethod(); i++ {
		method := testType.Method(i)
		if strings.HasPrefix(method.Name, "Test") {
			if be {
				value.MethodByName("BeforeEach").Call(nil)
			}
			finalMethod := value.MethodByName(method.Name)
			finalMethod.Call(nil)
		}
	}
}

// Assert tests v's truthiness.
func (t TestCase) Assert(v bool) {
	if !v {
		t.T.Fatalf("Expected true, got false [%s]", CallerInfo())
	}
}

// AssertError tests that err is non-nil.
func (t TestCase) AssertError(err error) {
	if err == nil {
		t.T.Fatalf("Expected error, got nil [%s]", CallerInfo())
	}
}

// AssertNil tests that i is nil.
func (t TestCase) AssertNil(i interface{}) {
	if i != nil {
		t.T.Fatalf("Expected %v to be nil [%s]", i, CallerInfo())
	}
}

// AssertNotNil tests that i is not nil.
func (t TestCase) AssertNotNil(i interface{}) {
	if i == nil {
		t.T.Fatalf("Expected %v NOT to be nil [%s]", i, CallerInfo())
	}
}

// AssertFalse tests v's falseness.
func (t TestCase) AssertFalse(v bool) {
	if v {
		t.T.Fatalf("Expected false, got true [%s]", CallerInfo())
	}
}

// Assertf tests v's truthiness, shows msg in case of failure.
func (t TestCase) Assertf(ok bool, msg string) {
	if !ok {
		t.T.Fatalf("Assertion failed: %s [%s]", msg, CallerInfo())
	}
}

// AssertOK tests for HTTP OK code.
func (t TestCase) AssertOK() {
	if t.err != nil {
		t.T.Fatalf("Request error is not nil: %s [%s]", t.err, CallerInfo())
	}
	if t.response == nil {
		t.T.Fatalf("Response is nil [%s]", CallerInfo())
	}
	if t.response.StatusCode != http.StatusOK {
		info := CallerInfo()
		t.T.Fatalf("Expected 200, got %d [%s]", t.response.StatusCode, info)
	}
}

// AssertStatus tests for some specific HTTP response code against the previous
// HTTP request.
func (t TestCase) AssertStatus(code int) {
	if t.response.StatusCode != code {
		i := CallerInfo()
		t.T.Fatalf("Expected %v, got %d [%s]", code, t.response.StatusCode, i)
	}
}

// Get issues an HTTP GET request and keeps the response for later assertions.
func (t *TestCase) Get(url string) {
	if t.server == nil {
		t.T.Fatalf("Uninitialized test server [%s]", CallerInfo())
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
		t.T.Fatalf("Uninitialized test server [%s]", CallerInfo())
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

// Put issues an HTTP PUT request and keeps the response for later assertions
func (t *TestCase) Put(url string, contentType string, body []byte) {
	if t.server == nil {
		t.T.Fatalf("Uninitialized test server [%s]", CallerInfo())
	}
	url = t.server.URL + url
	req, err := http.NewRequest("PUT", url, bytes.NewReader(body))
	if err != nil {
		t.T.Fatalf("Failed to create new request: %s", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
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
		t.T.Fatalf("Response is nil [%s]", CallerInfo())
	}
	if err := json.Unmarshal(t.ResponseBody, i); err != nil {
		t.T.Fatalf("Failed to unmarshal response body data: %s", err)
	}
}

// Marshal converts interface i to JSON, an error makes the test fail.
func (t *TestCase) Marshal(i interface{}) []byte {
	body, err := json.Marshal(i)
	if err != nil {
		t.T.Fatalf("Failed to marshal data: %s [%s]", err, CallerInfo())
	}
	return body
}

// AssertEqualInt tests equality for int values.
func (t *TestCase) AssertEqualInt(expected, actual int) {
	if expected != actual {
		info := CallerInfo()
		t.T.Fatalf("Expected %d to equal %d [%s]", actual, expected, info)
	}
}

// AssertEqualInt64 tests equality for int64 values.
func (t *TestCase) AssertEqualInt64(expected, actual int64) {
	if expected != actual {
		info := CallerInfo()
		t.T.Fatalf("Expected %d to equal %d [%s]", actual, expected, info)
	}
}

// AssertEqualStr tests if expected is equal to actual.
func (t *TestCase) AssertEqualStr(expected, actual string) {
	if expected != actual {
		t.T.Fatalf("Expected %q, got %q [%s]", expected, actual, CallerInfo())
	}
}

// Stolen from stretchr/testify with a few changes.
// CallerInfo is necessary  because the assert functions use  the testing object
// internally, causing  it to print the  file:line of the assert  method, rather
// than where the problem actually occured in calling code.

// CallerInfo returns an array of strings containing the file and line number
// of each stack frame leading from the current test to the assert call that
// failed.
func CallerInfo() string {

	pc := uintptr(0)
	file := ""
	line := 0
	ok := false
	name := ""

	callers := []string{}
	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			return ""
		}

		// This is a huge edge case, but it will panic if this is the case,
		// see #180
		if file == "<autogenerated>" {
			break
		}

		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if (dir != "assert" && dir != "mock" && dir != "require") ||
			file == "mock_test.go" {
			callers = append(callers, fmt.Sprintf("%s:%d", file, line))
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()
		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
		if isTest(name, "Test") ||
			isTest(name, "Benchmark") ||
			isTest(name, "Example") {
			break
		}
	}
	if len(callers) > 0 {
		return callers[len(callers)-1]
	}
	return ""
}

// Stolen from the `go  test` tool isTest tells whether name  looks like a test.
// (or  benchmark, according  to prefix)  It  is a  Test  (say) if  there is  a.
// character  after  Test  that  is  not a  lower-case  letter  We  don't  want.
// TesticularCancer                                                            .
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}
