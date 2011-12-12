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
	"reflect"
	"regexp"
	"strings"
)

// Non-nil values returned by the Test method will cause the table test that
// called Test to fail. A FatalError retured by Test stops the table test.
type T interface {
	Test(Testing) // Execute the test described by the object.
}

// These types act on strings values of uncaught panics inside Test() the
// type's test method.  Acceptable value types are substrings,
// *regexp.Regexp, or func(Testing, interface{}) objects.
type PanicExpectation interface{}

func acceptablePanicExpectation(t Testing, exp PanicExpectation) (ok bool) {
	switch exp.(type) {
	case nil:
		t.Error("nil PanicExpectation")
		return
	case string, *regexp.Regexp, func(Testing, string):
		return true
	}
	t.Errorf("unacceptable PanicExpectation type %s", reflect.TypeOf(exp))
	return
}

func applyPanicExpectation(t Testing, exp PanicExpectation, panicv interface{}) {
	switch exp.(type) {
	case *regexp.Regexp:
		r := exp.(*regexp.Regexp)
		if p := sprint(panicv); !r.MatchString(p) {
			t.Errorf("unexpected panic (doesn't match %v): %s", r, p)
		}
	case string:
		if p := sprint(panicv); strings.Index(p, exp.(string)) < 0 {
			t.Errorf("unexpected panic (doesn't contain %#v): %s", exp, p)
		}
	case func(Testing, interface{}):
		exp.(func(Testing, interface{}))(subT("callback function", t), panicv)
	}
}

type indexedError struct {
	index int
	Err   error
}

func (err indexedError) Error() string { return err.Err.Error() }

func applyPanicExpectations(t Testing, exps []PanicExpectation, panicv interface{}) {
	for i, exp := range exps {
		applyPanicExpectation(subT(sprintf("panic expectation %d", i), t), exp, panicv)
	}
}

type TPanics interface {
	T                           // TPanics is a TType.
	Panics() []PanicExpectation // TPanics when non-nil, certain panics expected.
}

func getTPanicsExpectations(t Testing, test TPanics) (exps []PanicExpectation, ok bool) {
	if test == nil {
		t.Error("nil test")
		return
	}
	exps = test.Panics()
	for i, exp := range exps {
		ok = ok && acceptablePanicExpectation(subT(sprintf("table.PanicExpectation %d", i), t), exp)
	}
	return
}

type TBefore interface {
	T               // TBefore is a T type.
	Before(Testing) // Callback executed before the Test method.
}

type TAfter interface {
	T              // TAfter is a T type.
	After(Testing) // Callback executed after the Test method.
}

type TBeforeAfter interface {
	T               // TBeforeAfter is a T type.
	Before(Testing) // TBeforeAfter is a TBefore type.
	After(Testing)  // TBeforeAfter is a TAfter type.
}

// Cast an element as a T, or create an os.Error describing the failure.
func mustT(t Testing, elem interface{}) (test T, err error) {
	switch elem.(type) {
	case nil:
		err = error_("nil slice element")
		t.Error(err)
		return
	case T:
	default:
		err = errorf("element does not implement table.T %v", reflect.TypeOf(elem))
		t.Error(err)
		return
	}
	return elem.(T), nil
}

// Execute t's Test method. If t is a TBefore type execute t.Before() prior to
// t.Test(). If t is a TAfter type, execute t.After() after t.Test() returns.
// Handles runtimes panics resulting from any of these callback.
func tTest(t Testing, test T) {
	place := "before"
	defer func() {
		if e := recover(); e != nil {
			t.Errorf("panic %s test; %v", place, e)
		}
	}()
	switch test.(type) {
	case TBeforeAfter:
		test.(TBefore).Before(subT("before test", t))
		defer test.(TAfter).After(subT("after test", t))
	case TBefore:
		test.(TBefore).Before(subT("before test", t))
	case TAfter:
		defer test.(TAfter).After(subT("after test", t))
	}
	place = "during"
	defer func() { place = "after" }()
	defer func() {
		switch panicv := recover(); test.(type) {
		case TPanics:
			exps, _ := getTPanicsExpectations(t, test.(TPanics))
			switch hasexp := len(exps) > 0; {
			case hasexp && panicv != nil:
				applyPanicExpectations(t, exps, panicv)
			case panicv != nil:
				t.Errorf("unexpected panic: %v", panicv)
			case hasexp:
				t.Errorf("test did not panic as expected %v", exps)
			}
			return
		default:
			if panicv != nil {
				t.Errorf("panic: %v", panicv)
			}
		}
	}()
	test.Test(t)
}
