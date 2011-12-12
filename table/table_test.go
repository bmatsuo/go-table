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
	//"errors"
	"reflect"
	"testing"
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
func onetestint()   { *testint1, *testint2 = 1, 1 }

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
		[]int{1, 2, 3}, func(i, x int) os.Error { *testint1 += i; *testint2 += x; return error_("fail") }, true,
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
		ft := fauxTest("doRange test", func(t Testing) {
			doRange(newTestingT("doRange slice", t), reflect.ValueOf(test.x), test.fn)
		})
		if ft.failed != test.isError {
			explain, qualify := "test is %san error", "not "
			if ft.failed {
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

type stringifyTest struct {
	i   int
	v   interface{}
	out string
}

func (test stringifyTest) Test(t Testing) {
	if str := stringifyIndex(test.i, test.v); str != test.out {
		t.Errorf("stringifyIndex(%d, %#v) => %#v != %#v", test.i, test.v, str, test.out)
	}
}

type testStringerTest struct{ in, out string }

func (s testStringerTest) String() string { return "simple string test" }
func (s testStringerTest) Test(t Testing) {}

var stringifyTests = []stringifyTest{
	{1, "abc", "abc"},
	{1, struct{ a, b int }{1, 2}, "struct { a int; b int } 1"},
	{0, stringifyTest{1, 1, "int 1"}, "table.stringifyTest 0"},
	{2, testStringerTest{"abc", "def"}, "simple string test 2"},
}

func TestStringify(t *testing.T) {
	for i, test := range stringifyTests {
		tTest(newTestingT(sprintf("stringify %d", i), t), test)
	}
}

type validateTableTest struct {
	table interface{}
	errs  []string
}

func (test validateTableTest) Test(t Testing) {
	meta := metaTestSimple{
		sprintf("validateTable(t, %#v)", test.table),
		func(t Testing) { validateTable(newTestingT("", t), test.table) },
		test.errs}
	meta.Test(t)
}

var validateTableTests = []validateTableTest{
	{[]int{1, 2, 3}, []string{}},
	{[]interface{}{1, "abc", 3 + 3i}, []string{}},
	{34, []string{"not a", "slice"}},
	{nil, []string{"invalid"}},
	{make([]int, 0), []string{"empty"}},
}

func TestValidateTable(t *testing.T) {
	for i, test := range validateTableTests {
		tTest(newTestingT(sprintf("vaidateTable %d", i), t), test)
	}
}
