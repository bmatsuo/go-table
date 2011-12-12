package table

/*  Filename:    test_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:58:18 PST 2011
 *  Description: For testing test.go
 */

import (
	"testing"
	"regexp"
	"os"
)

// Test the internal tTest function.
type tTestTest struct {
	fn    func(t Testing)
	esubs []string
}

func (test tTestTest) Test(t Testing) { test.fn(t) }

var tTestTests = []tTestTest{
	{func(t Testing) {}, nil},
	{func(t Testing) { t.Error(os.NewError("emsg")) }, []string{"emsg"}},
	{func(t Testing) { t.Fatal("fmsg") }, []string{"fmsg"}},
	{func(t Testing) { panic("pmsg") }, []string{"pmsg"}},
}

func TestTTest(t *testing.T) {
	for i, test := range tTestTests {
		prefix := sprintf("tTestTest %d", i)
		ft := fauxTest(prefix, func(t Testing) { tTest(t, test) })
		switch failed := ft.failed; {
		case !failed && test.esubs == nil:
			break
		case !failed:
			t.Errorf("%s: unexpected nil error (not %v)", prefix, test.esubs)
		case test.esubs == nil:
			t.Errorf("%s: unexpected error %v", prefix, ft.log)
		default:
			for j, sub := range test.esubs {
				if !ft.logLike(sub) {
					t.Errorf("%s: error missing pattern %d %#v; %#v", prefix, j, sub, ft.log)
				}
			}
		}
	}
}

type tTestExtraTest struct {
	before, after, verify func(Testing)
	exp                   interface{}
	errpatt               string
}

func (test tTestExtraTest) Test(t Testing) { test.verify(t) }

func testingCallNonNil(fn func(Testing), t Testing) {
	if fn != nil {
		fn(t)
	}
}
func (test tTestExtraTest) Before(t Testing) { testingCallNonNil(test.before, t) }
func (test tTestExtraTest) After(t Testing)  { testingCallNonNil(test.after, t) }
func (test tTestExtraTest) Panics() (exps []PanicExpectation) {
	switch test.exp.(type) {
	case nil:
	case string:
		if sub := test.exp.(string); sub != "" {
			return append(exps, sub)
		}
	default:
		return append(exps, test.exp)
	}
	return
}

var tTestExtraTestInt = new(int)

func tTestPanic(msg string) func(Testing)   { return func(t Testing) { panic(msg) } }
func tTestBeforeInt(plus int) func(Testing) { return func(t Testing) { (*tTestExtraTestInt) += plus } }
func tTestAfterInt(minus int) func(Testing) { return func(t Testing) { (*tTestExtraTestInt) -= minus } }
func tTestNoOp() func(Testing)              { return func(t Testing) {} }
func tTestVerifyInt(x int) func(Testing) {
	return func(t Testing) {
		if y := (*tTestExtraTestInt); y != x {
			t.Errorf("test integer value %d != %d", y, x)
		}
	}
}

var tTestExtraTests = []tTestExtraTest{
	{nil, nil, func(t Testing) {}, nil, ""}, // Sanity test.
	{nil, nil, tTestPanic("gophers"), "gophers", ""},
	{nil, nil, tTestPanic("gophers"), regexp.MustCompile("gophers?"), ""},
	{nil, nil, tTestPanic("gophers"), func(t Testing, panicv interface{}) {
		if p := sprint(panicv); p != "gophers" {
			t.Errorf("unexpected panic (missing \"gophers\"): %s", p)
		}
	}, ""},

	// Order is important for next group.
	{tTestBeforeInt(1), nil, tTestVerifyInt(1), nil, ""},              // Tests Before call.
	{nil, tTestAfterInt(1), tTestVerifyInt(1), nil, ""},               // Ensures the value of the integer is persistant.
	{tTestBeforeInt(3), tTestAfterInt(3), tTestVerifyInt(3), nil, ""}, // Tests After call.
	{tTestBeforeInt(1), tTestAfterInt(1), tTestVerifyInt(1), nil, ""}, // Tests both Before and After work togeter.

	{tTestPanic("gophers"), tTestNoOp(), tTestNoOp(), nil, "before.*gophers"},
	{tTestNoOp(), tTestPanic("gophers"), tTestNoOp(), nil, "after.*gophers"},
}

func TestTTestExtraTests(t *testing.T) {
	for i, test := range tTestExtraTests {
		prefix := sprintf("extra functionality test %d:", i)
		ft := fauxTest(prefix, func(t Testing) { tTest(t, test) })
		switch failed := ft.failed; {
		case !failed && test.errpatt == "":
			continue
		case !failed:
			t.Errorf("%s missing expected error: %v", prefix, test.errpatt)
		case test.errpatt == "":
			t.Errorf("%s unexpected error: %v", prefix, ft.log)
		case !ft.logLike(test.errpatt):
			t.Errorf("%s unexpected error (not %s): %v", prefix, test.errpatt, ft.log)
		}
	}
}
