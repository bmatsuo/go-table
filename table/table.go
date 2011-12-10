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

func value(src string, index, v reflect.Value, zero reflect.Value) reflect.Value {
	switch {
	case v.Interface() == nil:
		return zero
	case !v.IsValid():
		panic(errorf("invalid value in %s index %v", src, index.Interface()))
	}
	return v
}

// Iterate over a range of values, issuing a callback for each one. The callback
// fn is expected to take two arguments (index/key, value pair) and return an
// os.Error.
func doRange(v reflect.Value, fn interface{}) os.Error {
	fnval := reflect.ValueOf(fn)
	fntyp := fnval.Type()
	if numin := fntyp.NumIn(); numin != 2 {
		panic(errorf("doRange function of %d arguments %v", numin, fn))
	}
	if numout := fntyp.NumOut(); numout != 1 {
		panic(errorf("doRange function of %d return values %v", numout, fn))
	}
	zero := reflect.Zero(fnval.Type().In(1))
	var out reflect.Value
	switch k := v.Kind(); k {
	case reflect.Slice:
		for i, n := 0, v.Len(); i < n; i++ {
			ival, vval := reflect.ValueOf(i), v.Index(i)
			arg := value("slice", ival, vval, zero)
			out = fnval.Call([]reflect.Value{ival, arg})[0]
			if !out.IsNil() {
				return out.Interface().(os.Error)
			}

		}
	case reflect.Map:
		for _, kval := range v.MapKeys() {
			vval := v.MapIndex(kval)
			arg := value("map", kval, vval, zero)
			out = fnval.Call([]reflect.Value{kval, arg})[0]
			if !out.IsNil() {
				return out.Interface().(os.Error)
			}
		}
	case reflect.Chan:
		var vval reflect.Value
		var ok bool
		for i := 0; true; i++ {
			ival := reflect.ValueOf(i)
			if vval, ok = v.Recv(); !ok {
				break
			}
			arg := value("chan", ival, vval, zero)
			out = fnval.Call([]reflect.Value{ival, arg})[0]
			if !out.IsNil() {
				return out.Interface().(os.Error)
			}
		}
	default:
		panic(errorf("unacceptable type for range %v", v.Type()))
	}
	return nil
}

type stringer interface {
	String() string
}

func testCastT(t *testing.T, prefix string, v interface{}) (test T, err os.Error) {
	if test, err = mustT(t, prefix, v); err != nil {
		if err == ErrSkip {
			if Verbose {
				t.Logf("%s skipped", prefix)
			}
			err = nil
		} else {
			t.Fatal(err)
		}
		return
	}
	return
}

func testExecute(t *testing.T, prefix string, test T) {
	if !error(t, prefix, tTest(test)) && Verbose {
		t.Logf("%s passed", prefix)
	}
}

func testMap(t *testing.T, v reflect.Value) os.Error {
	i := new(int)
	return doRange(v, func(k, v interface{}) (err os.Error) {
		var prefix string
		switch k.(type) {
		case string:
			prefix = k.(string)
		case stringer:
			prefix = k.(stringer).String()
		default:
			prefix = sprintf("%v %d", reflect.TypeOf(k).String(), *i)
		}
		var test T
		test, err = testCastT(t, prefix, v)
		testExecute(t, prefix, test)
		(*i)++
		return
	})
}

// Test each value in a slice table.
func testSlice(t *testing.T, v reflect.Value) os.Error {
	return doRange(v, func(i int, elem interface{}) (err os.Error) {
		prefix := sprintf("%v %d", reflect.TypeOf(elem), i)
		var test T
		test, err = testCastT(t, prefix, elem)
		testExecute(t, prefix, test)
		return
	})
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
	case reflect.Slice, reflect.Map: // Allow chan?
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
	case reflect.Map:
		if testMap(t, val) != nil {
			return
		}
	default:
		fatal(t, prefix, errorf("unexpected error"))
	}
}
