package yatf

import (
	"reflect"
	"strings"
	"testing"
)

func runTests(i interface{}) {
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

func TestAll(t *testing.T) {
	runTests(TestCase{T: t})
	runTests(&TestCase{T: t})
}

func (t TestCase) TestAddition() {
	t.Assert(2 == 1*4/2)
	t.Assert(0 == -1+1)
}

func (t TestCase) TestDivision() {
	t.AssertFalse(4 == 10/3)
}

func (t TestCase) TestNonOKResponse() {
	t.AssertCreated(201)
}

func (t TestCase) TestOKResponse() {
	t.AssertOK(200)
}
