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
	"strings"
	"regexp"
	"os"
)

// Non-nil values returned by the Test method will cause the table test that
// called Test to fail. A FatalError retured by Test stops the table test.
type T interface {
	Test(Testing) // Execute the test described by the object.
}

// These types act on strings values of uncaught panics inside Test() the
// type's test method.  Acceptable value types are substrings,
// *regexp.Regexp, or func(string) os.Error objects.
type PanicExpectation interface{}

func acceptablePanicExpectation(exp PanicExpectation) os.Error {
	switch exp.(type) {
	case nil:
		return error_("nil PanicExpectation")
	case string, *regexp.Regexp, func(string) os.Error:
		return nil
	}
	return errorf("unacceptable PanicExpectation type %s", reflect.TypeOf(exp))
}

func applyPanicExpectation(exp PanicExpectation, pstr string) (err os.Error) {
	switch exp.(type) {
	case *regexp.Regexp:
		r := exp.(*regexp.Regexp)
		if !r.MatchString(pstr) {
			err = fatalf("unexpected panic (doesn't match %v): %s", r, pstr)
		}
	case string:
		if strings.Index(pstr, exp.(string)) < 0 {
			err = fatalf("unexpected panic (doesn't contain %#v): %s", exp, pstr)
		}
	case func(string) os.Error:
		if err := exp.(func(string) os.Error)(pstr); err != nil {
			err = fatalf("panic callback error: %s", err)
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
			err = fatalf("PanicExpectation %d failure: %v", e.index, e)
			return
		}
		err = fatal("PanicExpectation failure")
		for _, e := range errs {
			err = fatalf("%s\n\texpectation %d: %v", err, e.index, e)
		}
	}
	return
}

type TPanics interface {
	T                           // TPanics is a TType.
	Panics() []PanicExpectation // TPanics when non-nil, certain panics expected.
}

func getTPanicsExpectations(test TPanics) (exps []PanicExpectation, err os.Error) {
	if test == nil {
		return nil, error_("nil test")
	}
	exps = test.Panics()
	for i, exp := range exps {
		if err = acceptablePanicExpectation(exp); err == nil {
			err = errorf("PanicExpectation %d type error: %v", i, err)
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
func mustT(t Testing, elem interface{}) (test T, err os.Error) {
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
		test.(TBefore).Before()
		defer test.(TAfter).After()
	case TBefore:
		test.(TBefore).Before()
	case TAfter:
		defer test.(TAfter).After()
	}
	place = "during"
	defer func() { place = "after" }()
	defer func() {
		panicv := recover()
		switch test.(type) {
		case TPanics:
			if exps, err := getTPanicsExpectations(test.(TPanics)); err != nil {
				t.Errorf("error retrieving PanicExpectations: %v", err)
			} else if hasexp := len(exps) > 0; panicv != nil {
				if hasexp {
					if err = applyPanicExpectations(exps, sprint(panicv)); err != nil {
						t.Error(err)
					}
				} else {
					t.Errorf("unexpected panic: %v", panicv)
				}
			} else {
				if hasexp {
					t.Errorf("test did not panic as expected %v", exps)
				}
			}
			return
		}
		if panicv != nil {
			t.Errorf("panic: %v", panicv)
		}
	}()
	test.Test(t)
}
