// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

/*  Filename:    test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:58:18 PST 2011
 *  Description: 
 */

import (
	"testing"
	"reflect"
	"strings"
	"regexp"
	"os"
)

// Non-nil values returned by the Test method will cause the table test that
// called Test to fail. A FatalError retured by Test stops the table test.
type T interface {
	Test() os.Error // Execute the test described by the object.
}

// These types act on strings values of uncaught panics inside Test() the
// type's test method.  Acceptable value types are substrings,
// *regexp.Regexp, or func(string) os.Error objects.
type PanicExpectation interface{}

func acceptablePanicExpectation(exp PanicExpectation) os.Error {
	switch exp.(type) {
	case nil:
		return Error("nil PanicExpectation")
	case string, *regexp.Regexp, func(string) os.Error:
		return nil
	}
	return Errorf("unacceptable PanicExpectation type %s", reflect.TypeOf(exp))
}

func applyPanicExpectation(exp PanicExpectation, pstr string) (err os.Error) {
	switch exp.(type) {
	case *regexp.Regexp:
		r := exp.(*regexp.Regexp)
		if !r.MatchString(pstr) {
			err = Fatalf("unexpected panic (doesn't match %v): %s", r, pstr)
		}
	case string:
		if strings.Index(pstr, exp.(string)) < 0 {
			err = Fatalf("unexpected panic (doesn't contain %#v): %s", exp, pstr)
		}
	case func(string) os.Error:
		if err := exp.(func(string) os.Error)(pstr); err != nil {
			err = Fatalf("panic callback error: %s", err)
		}
	}
	return
}

type indexedError struct {
	index int
	Err   os.Error
}

func (err indexedError) String() string { return err.Err.String() }

func applyPanicExpectations(exps []PanicExpectation, pstr string) (err os.Error) {
	var errs []indexedError
	for i, exp := range exps {
		if e := applyPanicExpectation(exp, pstr); e != nil {
			errs = append(errs, indexedError{i, e})
		}
	}
	if len(errs) > 0 {
		if e := errs[0]; len(errs) == 1 {
			err = Fatalf("PanicExpectation %d failure: %v", e.index, e)
			return
		}
		err = Fatal("PanicExpectation failure")
		for _, e := range errs {
			err = Fatalf("%s\n\texpectation %d: %v", err, e.index, e)
		}
	}
	return
}

type TPanics interface {
	T // TPanics is a TType.
	Panics() []PanicExpectation
}

func getTPanicsExpectations(test TPanics) (exps []PanicExpectation, err os.Error) {
	if test == nil {
		return nil, Error("nil test")
	}
	exps = test.Panics()
	for i, exp := range exps {
		if err = acceptablePanicExpectation(exp); err == nil {
			err = Errorf("PanicExpectation %d type error: %v", i, err)
			return
		}
	}
	return
}

type TBefore interface {
	T        // TBefore is a T type.
	Before() // Callback executed before the Test method.
}

type TAfter interface {
	T       // TAfter is a T type.
	After() // Callback executed after the Test method.
}

type TBeforeAfter interface {
	T        // TBeforeAfter is a T type.
	Before() // TBeforeAfter is a TBefore type.
	After()  // TBeforeAfter is a TAfter type.
}

// Cast an element as a T, or create an os.Error describing the failure.
func mustT(t *testing.T, name string, elem interface{}) (T, os.Error) {
	switch elem.(type) {
	case nil:
		return nil, os.NewError(fatalstrf(name, "nil slice element"))
	case T:
	default:
		errf := "element does not implement table.T %v"
		return nil, os.NewError(fatalstrf(name, errf, reflect.TypeOf(elem)))
	}
	return elem.(T), nil
}

// Execute t's Test method. If t is a TBefore type execute t.Before() prior to
// t.Test(). If t is a TAfter type, execute t.After() after t.Test() returns.
// Handles runtimes panics resulting from any of these callback.
func tTest(t T) (err os.Error) {
	place := "before"
	defer func() {
		if e := recover(); e != nil {
			err = Errorf("panic %s test; %v", place, e)
		}
	}()
	switch t.(type) {
	case TBeforeAfter:
		t.(TBefore).Before()
		defer t.(TAfter).After()
	case TBefore:
		t.(TBefore).Before()
	case TAfter:
		defer t.(TAfter).After()
	}
	place = "during"
	defer func() { place = "after" }()
	defer func() {
		panicv := recover()
		switch t.(type) {
		case TPanics:
			exps, err := getTPanicsExpectations(t.(TPanics))
			if err != nil {
				err = Errorf("error retrieving PanicExpectations: %v", err)
			} else if hasexp := len(exps) > 0; panicv != nil {
				if hasexp {
					err = applyPanicExpectations(exps, sprint(panicv))
				} else {
					err = Errorf("unexpected panic: %v", panicv)
				}
			} else {
				if hasexp {
					err = Errorf("test did not panic as expected %v", exps)
				}
			}
			return
		}
		if panicv != nil {
			err = Errorf("panic: %v", panicv)
		}
	}()
	err = t.Test()
	return
}
