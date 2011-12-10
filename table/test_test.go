package table

/*  Filename:    test_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:58:18 PST 2011
 *  Description: For testing test.go
 */

import (
	"testing"
	"strings"
	"regexp"
	"os"
)

// Test the internal tTest function.
type tTestTest struct {
	fn    func() os.Error
	fatal bool
	esubs []string
}

func (test tTestTest) Test() os.Error { return test.fn() }

var tTestTests = []tTestTest{
	{func() os.Error { return nil }, false, nil},
	{func() os.Error { return os.NewError("emsg") }, false, []string{"emsg"}},
	{func() os.Error { return Fatal("fmsg") }, true, []string{"fmsg"}},
	{func() os.Error { panic("pmsg") }, false, []string{"pmsg"}},
}

func TestTTest(t *testing.T) {
	for i, test := range tTestTests {
		switch err := tTest(test); {
		case err == nil && (test.esubs != nil || test.fatal):
			if test.fatal {
				t.Errorf("tTestTest %d: unexpected nil error (not fatal %v)", i, test.esubs)
			} else {
				t.Errorf("tTestTest %d: unexpected nil error (not %v)", i, test.esubs)
			}
		case err == nil:
			break
		case test.esubs == nil && !test.fatal:
			t.Errorf("tTestTest %d: unexpected error %v", i, err)
		default:
			var fatal bool
			switch err.(type) {
			case FatalError:
				fatal = true
			default:
			}
			if test.fatal != fatal {
				if fatal {
					t.Errorf("tTestTest %d: unexpected fatal error %v", err)
				} else {
					t.Errorf("tTestTest %d: unexpected non-fatal error %v", err)
				}
			}
			for j, sub := range test.esubs {
				estr := err.String()
				if strings.Index(estr, sub) < 0 {
					t.Errorf("tTestTest %d: error missing substring %d; %s", i, j, estr)
				}
			}
		}
	}
}

type tTestExtraTest struct {
	before, after func()
	verify        func() os.Error
	exp           interface{}
	errpatt       string
}

func (test tTestExtraTest) Test() (err os.Error) {
	return test.verify()
	return
}

func fnCallNonNil(fn func()) {
	if fn != nil {
		fn()
	}
}
func (test tTestExtraTest) Before() { fnCallNonNil(test.before) }
func (test tTestExtraTest) After()  { fnCallNonNil(test.after) }
func (test tTestExtraTest) Panics() (exps []PanicExpectation) {
	switch test.exp.(type) {
	case string:
		if sub := test.exp.(string); sub != "" {
			return append(exps, sub)
		}
	default:
		if p := test.exp; p == nil {
			return append(exps, p)
		}
	}
	return
}

var tTestExtraTestInt = new(int)

func tTestPanicTest(msg string) func() os.Error { return func() os.Error { panic(msg) } }
func tTestPanic(msg string) func()              { return func() { panic(msg) } }
func tTestBeforeInt(plus int) func()            { return func() { (*tTestExtraTestInt) += plus } }
func tTestAfterInt(minus int) func()            { return func() { (*tTestExtraTestInt) -= minus } }
func tTestNoOpTest() func() os.Error            { return func() os.Error { return nil } }
func tTestNoOp() func()                         { return func() {} }
func tTestVerifyInt(x int) func() os.Error {
	return func() os.Error {
		if y := (*tTestExtraTestInt); y != x {
			return Errorf("test integer value %d != %d", y, x)
		}
		return nil
	}
}

var tTestExtraTests = []tTestExtraTest{
	{nil, nil, func() os.Error { return nil }, "", ""}, // Sanity test.
	{nil, nil, tTestPanicTest("gophers"), "gophers", ""},
	{nil, nil, tTestPanicTest("gophers"), regexp.MustCompile("gophers?"), ""},
	{nil, nil, tTestPanicTest("gophers"), func(pstr string) os.Error {
		if pstr != "gophers" {
			return Errorf("unexpected panic: %s", pstr)
		}
		return nil
	}, ""},

	// Order is important for next group.
	{tTestBeforeInt(1), nil, tTestVerifyInt(1), "", ""},              // Tests Before call.
	{nil, tTestAfterInt(1), tTestVerifyInt(1), "", ""},               // Ensures the value of the integer is persistant.
	{tTestBeforeInt(3), tTestAfterInt(3), tTestVerifyInt(3), "", ""}, // Tests After call.
	{tTestBeforeInt(1), tTestAfterInt(1), tTestVerifyInt(1), "", ""}, // Tests both Before and After work togeter.

	{tTestPanic("gophers"), tTestNoOp(), tTestNoOpTest(), "", "before.*gophers"},
	{tTestNoOp(), tTestPanic("gophers"), tTestNoOpTest(), "", "after.*gophers"},
}

func TestTTestExtraTests(t *testing.T) {
	for i, test := range tTestExtraTests {
		prefix := sprintf("extra functionality test %d:", i)
		var err os.Error
		switch err = tTest(test); {
		case err == nil && test.errpatt == "":
			return
		case err == nil:
			t.Errorf("%s missing expected error: %s", prefix, test.errpatt)
		case !regexp.MustCompile(test.errpatt).MatchString(err.String()):
			t.Errorf("%s unexpected error (not %s): %v", prefix, test.errpatt, err)
		}
		t.Errorf("%s unexpected error: %v", prefix, err)
	}
}
