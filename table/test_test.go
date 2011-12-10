package table

/*  Filename:    test_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:58:18 PST 2011
 *  Description: For testing test.go
 */

import (
	"testing"
	"strings"
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
	{func() os.Error { panic("pmsg") }, true, []string{"pmsg"}},
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
