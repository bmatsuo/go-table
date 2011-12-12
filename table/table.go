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
	"reflect"
	"testing"
	"os"
)

func validValue(t *testingT, v reflect.Value, zero reflect.Value) reflect.Value {
	switch {
	case v.Interface() == nil:
		return zero
	case !v.IsValid():
		t.Error("invalid value")
	}
	return v
}

// Iterate over a range of values, issuing a callback for each one. The callback
// fn is expected to take two arguments (index/key, value pair) and return an
// os.Error.
func doRange(t *testingT, v reflect.Value, fn interface{}) {
	t = t.sub("internal doRange")
	fnval := reflect.ValueOf(fn)
	fntyp := fnval.Type()
	if numin := fntyp.NumIn(); numin != 2 {
		t.Fatalf("function of %d arguments %v", numin, fn)
	} else if numout := fntyp.NumOut(); numout != 1 {
		t.Fatalf("function of %d return values %v", numout, fn)
	}
	zero := reflect.Zero(fnval.Type().In(1))
	var out reflect.Value
	switch k := v.Kind(); k {
	case reflect.Slice:
		for i, n := 0, v.Len(); i < n; i++ {
			ival, vval := reflect.ValueOf(i), v.Index(i)
			arg := validValue(t.sub(sprintf("index %d", i)), vval, zero)
			if !arg.IsValid() {
				continue
			}
			out = fnval.Call([]reflect.Value{ival, arg})[0]
			if !out.IsNil() {
				t.Error(out.Interface())
				break
			}

		}
	case reflect.Map:
		for _, kval := range v.MapKeys() {
			vval := v.MapIndex(kval)
			arg := validValue(t.sub(sprintf("index %v", kval.Interface())), vval, zero)
			if !arg.IsValid() {
				continue
			}
			out = fnval.Call([]reflect.Value{kval, arg})[0]
			if !out.IsNil() {
				t.Error(out.Interface())
				break
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
			arg := validValue(t.sub(sprintf("received value %d", i)), vval, zero)
			if !arg.IsValid() {
				continue
			}
			out = fnval.Call([]reflect.Value{ival, arg})[0]
			if !out.IsNil() {
				t.Error(out.Interface())
				break
			}
		}
	default:
		t.Fatalf("unacceptable type for range %v", v.Type())
	}
}

func testMap(t *testingT, v reflect.Value) {
	doRange(t.sub("map"), v, func(k, v interface{}) os.Error {
		sub := t.sub(sprint(k))
		if test, err := mustT(sub, v); err == nil {
			tTest(sub, test)
		}
		return nil
	})
}

type stringer interface {
	String() string
}

func stringifyIndex(i int, v interface{}) string {
	switch v.(type) {
	case string:
		return v.(string)
	case stringer:
		return sprintf("%v %d", v, i)
	default:
	}
	return sprintf("%v %d", reflect.TypeOf(v), i)
}

// Test each value in a slice table.
func testSlice(t *testingT, v reflect.Value) {
	doRange(t.sub("slice"), v, func(i int, elem interface{}) os.Error {
		sub := t.sub(stringifyIndex(i, elem))
		if test, err := mustT(sub, elem); err == nil {
			tTest(sub, test)
		}
		return nil
	})
}

// Detect a value's reflect.Kind. Return the reflect.Value as well for good measure.
func kind(x interface{}) (reflect.Value, reflect.Kind) { v := reflect.ValueOf(x); return v, v.Kind() }

func validateTable(t *testingT, table interface{}) (tab reflect.Value, k reflect.Kind) {
	tab, k = kind(table)
	switch k {
	case reflect.Invalid:
		t.Fatal("table is invalid")
	case reflect.Slice, reflect.Map: // Allow chan?
		break
	default:
		t.Fatalf("table %v is not a slice", tab.Type())
	}

	// A table can't be empty.
	if tab.Len() == 0 && k != reflect.Chan {
		t.Fatal("empty table")
	}
	return
}

func testHelper(t *testingT, table interface{}) {
	tinternal := subT("internal table.Test", t)
	val, k := validateTable(tinternal.sub("table validation"), table)
	switch k {
	case reflect.Slice:
		testSlice(t, val)
	case reflect.Map:
		testMap(t, val)
	default:
		tinternal.Fatalf("unexpected table kind %v", k)
	}
}

// A table driven test. The table must be a slice of values all implementing T.
// But, not all elements need be of the same type. And furthermore, the slice's
// element type does not need to satisfy T. For example, a slice v of type
// []interface{} can be a valid table if all its elements satisfy T.
//
// A feasible future enhancement would be to allow map tables. Possibly chan
// tables.
func Test(t *testing.T, table interface{}) { testHelper(subT("", t), table) }
