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
	"strings"
	"fmt"
	"os"
)

// An error that causes a table test to issue a fatal "testing" error to package
// "testing".
type FatalError struct{ os.Error }

func (f FatalError) String() string    { return f.Error.String() }

// Create a new FatalError.
func Fatal(err os.Error) FatalError { return FatalError{err} }

// Elements of table tests must individually satisfy type T. The table elements
// themselves describe tests, and their Test method should execute that test.
// Any non-nil os.Error retured by a Test method will cause the table test to fail.
// A FatalError retured by Test will cause the table test to issue a fatal error
// and stop the table test.
type T interface {
	Test() os.Error
}

var (
	Verbose bool       // If true more verbose errors, as well as passed tests, are logged.
	MsgFmt  = "%s: %s" // Result output format. Can be changed to another 2-argument string.
)

// Aliases to "fmt" functions.
func sprint(v interface{}) string                     { return fmt.Sprint(v) }
func sprintf(format string, v ...interface{}) string  { return fmt.Sprintf(format, v...) }
func errorf(format string, v ...interface{}) os.Error { return fmt.Errorf(format, v...) }

// Create a named messages with name format MsgFmt.
func msg(name string, v interface{}) string { return sprintf(MsgFmt, name, sprint(v)) }
func msgf(name, format string, v ...interface{}) string {
	return sprintf(MsgFmt, name, sprintf(format, v...))
}

// Execute a callback if the given error is non-nil.
func onerror(err os.Error, fn func(os.Error)) bool {
	if err != nil {
		fn(err)
		return true
	}
	return false
}

// Functions to generate error strings.
func errorstr(name string, v interface{}) (err string) {
	if Verbose {
		err = msg(sprintf("%s error", name), err)
	} else {
		err = msg(name, err)
	}
	return
}
func fatalstr(name string, v interface{}) string {
	if Verbose {
		name = sprintf("%s fatal", name)
	}
	return errorstr(name, v)
}
// Functions to generate error strings with formatted messages.
func errorstrf(name, format string, v ...interface{}) string {
	return errorstr(name, sprintf(format, v...))
}
func fatalstrf(name, format string, v ...interface{}) string {
	return fatalstr(name, sprintf(format, v...))
}

// Functions to log errors when they occur.
func error(T *testing.T, name string, err os.Error) bool {
	return onerror(err, func(err os.Error) { T.Error(errorstr(name, err)) })
}
func fatal(T *testing.T, name string, err os.Error) bool {
	return onerror(err, func(err os.Error) { T.Error(fatalstr(name, err)) })
}

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

// So far unused internal error telling the test to skip a table element.
var errSkip = os.NewError("skip")

// Cast an element as a T, or create an os.Error describing the failure.
func mustT(t *testing.T, name string, elem interface{}) (T, os.Error) {
	switch elem.(type) {
	case nil:
		return nil, errorf(strings.Replace(
			fatalstrf(name, "nil slice element"),
			"%", "%%", -1))
	case T:
	default:
		return nil, errorf(strings.Replace(
			fatalstrf(name, "element does not implement table.T %v", reflect.TypeOf(elem)),
			"%", "%%", -1))
	}
	return elem.(T), nil
}

// Test each value in a slice table.
func testSlice(t *testing.T, v reflect.Value) (err os.Error) {
	err = doRange(v, func(i int, elem interface{}) (err os.Error) {
		var e T
		prefix := sprintf("%v %d", reflect.TypeOf(e), i)
		if e, err = mustT(t, prefix, elem); err != nil {
			if err == errSkip {
				err = nil
				if Verbose {
					t.Log(sprintf("%s skipped", prefix, i))
				}
			} else {
				t.Fatal(err)
			}
			return
		}
		if !error(t, prefix, e.Test()) && Verbose {
			t.Log(sprintf("%s passed", prefix, i))
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
