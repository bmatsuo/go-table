// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*  Filename:    table.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Thu Dec  8 10:10:58 PST 2011
 *  Description: Main source file in go-table
 */

/*
Package table provides a simple framework for executing table driven tests.
A table is a set (usually a slice) of tests. A table test element is usually a
simple struct that describes a singular test. In the table package, table test
elements implement their own test(s) as a method Test which returns an os.Error.
Non-nil errors returned by a table test's Test method cause errors to be logged
with the "testing" package.

For general information about table driven testing in Go, see

	http://code.google.com/p/go-wiki/wiki/TableDrivenTests
*/
package table

import (
	"testing"
	"reflect"
	"os"
)

// Iterate over a range of values, issuing a callback for each one. The callback
// fn is expected to take two arguments (index/key, value pair) and return an
// os.Error.
func doRange(v reflect.Value, fn interface{}) os.Error {
	switch k := v.Kind(); k {
	case reflect.Slice:
		for i, n := 0, v.Len(); i < n; i++ {
			if out := reflect.ValueOf(fn).Call([]reflect.Value{reflect.ValueOf(i), v.Index(i)})[0]; !out.IsNil() {
				return out.Interface().(os.Error)
			}
		}
	default:
		panic(errorf("unacceptable type for range %v", v.Type()))
	}
	return nil
}

// Test each value in a slice table.
func testSlice(t *testing.T, v reflect.Value) (err os.Error) {
	err = doRange(v, func(i int, elem interface{}) (err os.Error) {
		var e T
		prefix := sprintf("%v %d", reflect.TypeOf(elem), i)
		if e, err = mustT(t, prefix, elem); err != nil {
			if err == ErrSkip {
				err = nil
				if Verbose {
					t.Logf("%s skipped", prefix, i)
				}
			} else {
				t.Fatal(err)
			}
			return
		}
		if !error(t, prefix, tTest(e)) && Verbose {
			t.Logf("%s passed", prefix, i)
		}
		return
	})
	return
}

// Detect a value's reflect.Kind. Return the reflect.Value as well for good measure.
func kind(x interface{}) (reflect.Value, reflect.Kind) { v := reflect.ValueOf(x); return v, v.Kind() }

// A table driven test. The table must be a slice of values all implementing T.
// But, not all elements need be of the same type. And furthermore, the slice's
// element type does not need to satisfy T. For example, a slice v of type
// []interface{} can be a valid table if all its elements satisfy T.
//
// A feasible future enhancement would be to allow map tables. Possibly chan
// tables.
func Test(t *testing.T, table interface{}) {
	prefix := "table.Test" // name for internal errors.

	// A table must be a slice type.
	val, k := kind(table)
	switch k {
	case reflect.Invalid:
		fatal(t, prefix, errorf("table is invalid"))
	case reflect.Slice: // Allow chan/map?
		break
	default:
		fatal(t, prefix, errorf("table %s is not a slice", val.Type().String()))
	}

	// A table can't be empty.
	if val.Len() == 0 && k != reflect.Chan {
		fatal(t, prefix, errorf("empty table"))
	}

	// Execute table tests.
	switch k {
	case reflect.Slice:
		if testSlice(t, val) != nil {
			return
		}
	default:
		fatal(t, prefix, errorf("unexpected error"))
	}
}
