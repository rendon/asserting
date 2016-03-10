package asserting

import "testing"

func TestAll(t *testing.T) {
	Run(TestCase{T: t})
	Run(&TestCase{T: t})
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
