// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

/*  Filename:    table_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Thu Dec  8 10:10:58 PST 2011
 *  Description: Main test file for go-table
 */

import (
	"testing"
	"reflect"
	"os"
)

func TestTable(T *testing.T) {

}

type doRangeTest struct {
	before, after func()
	x, fn         interface{}
	isError       bool
	verify        func() os.Error
}

var testint1 = new(int)
var testint2 = new(int)

func cleartestint() { *testint1, *testint2 = 0, 0 }
func onetestint() { *testint1, *testint2 = 1, 1 }

var doRangeTests = []doRangeTest{
	// Make sure range callbacks are being called correctly.
	{cleartestint, cleartestint,
		[]int{1, 2, 3},
		func(i, x int) os.Error { *testint1 += i; *testint2 += x; return nil }, false,
		func() (err os.Error) {
			if !(*testint1 == 3 && *testint2 == 6) {
				err = errorf("testint1: %d, testint2: %d", *testint1, *testint2)
			}
			return
		}},
	// Meta test to ensure test callbacks are being called correctly.
	{onetestint, cleartestint,
		[]int{1, 2, 3, 4},
		func(i, x int) os.Error { *testint1 *= i; *testint2 *= x; return nil }, false,
		func() (err os.Error) {
			if !(*testint1 == 0 && *testint2 == 24) {
				err = errorf("testint1: %d, testint2: %d", *testint1, *testint2)
			}
			return
		}},
	// Ensure that errors from the range callback stop the iteration.
	{onetestint, cleartestint,
		[]int{1, 2, 3}, func(i, x int) os.Error{ *testint1 += i; *testint2 += x; return os.NewError("fail") }, true,
		func() (err os.Error) {
			if !(*testint1 == 1 && *testint2 == 2) {
				err = errorf("testint1: %d, testint2: %d", *testint1, *testint2)
			}
			return
		}},
}

func TestDoRange(t *testing.T) {
	for i, test := range doRangeTests {
		name := sprintf("doRange %d", i)
		if test.before != nil {
			test.before()
		}
		if err := doRange(reflect.ValueOf(test.x), test.fn); (err != nil) != test.isError {
			explain, qualify := "test is %san error", "not "
			if err != nil {
				qualify = ""
			}
			t.Error(errorstrf(name, explain, qualify))
		}
		if err := test.verify(); err != nil {
			t.Error(errorstrf(name, "verification failure: %v", err))
		}
		if test.after != nil {
			test.after()
		}
	}
}
