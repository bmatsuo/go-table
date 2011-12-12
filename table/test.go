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

// The interface table elements must satisfy.
type Element interface {
	Test(T) // Execute the test described by the object.
}

// Act on uncaught panics from the type's Test method. Underlying type must be
// string, *regexp.Regexp, or func(T, interface{}). String and regexp values
// test against the panic's string value. String values search for a substring.
type PanicExpectation interface{}

func acceptablePanicExpectation(t T, exp PanicExpectation) (ok bool) {
	switch exp.(type) {
	case nil:
		t.Error("nil PanicExpectation")
		return
	case string, *regexp.Regexp, func(T, string):
		return true
	}
	t.Errorf("unacceptable PanicExpectation type %s", reflect.TypeOf(exp))
	return
}

func applyPanicExpectation(t T, exp PanicExpectation, panicv interface{}) {
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
	case func(T, interface{}):
		exp.(func(T, interface{}))(subT("callback function", t), panicv)
	}
}

type indexedError struct {
	index int
	Err   error
}

func (err indexedError) Error() string { return err.Err.Error() }

func applyPanicExpectations(t T, exps []PanicExpectation, panicv interface{}) {
	for i, exp := range exps {
		applyPanicExpectation(subT(sprintf("panic expectation %d", i), t), exp, panicv)
	}
}

type ElementPanics interface {
	Element                     // Is an Element.
	Panics() []PanicExpectation // When non-nil, certain panic values are expected.
}

func getElementPanicsExpectations(t T, test ElementPanics) (exps []PanicExpectation, ok bool) {
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

type ElementBefore interface {
	Element   // Is an Element.
	Before(T) // Callback executed before the Test method.
}

type ElementAfter interface {
	Element  // Is an Element.
	After(T) // Callback executed after the Test method.
}

type ElementBeforeAfter interface {
	Element   // Is an Element.
	Before(T) // Is an ElementBefore.
	After(T)  // Is an ElementAfter.
}

// Not an Element.
type ElementGenerator interface {
	Generate(T) []Element // A (possibly empty) slice of Elements.
}

// Cast an value as an Element, or create an error describing the failure.
func mustElement(t T, elem interface{}) (tests []Element, err error) {
	switch elem.(type) {
	case nil:
		err = error_("nil slice element")
		t.Error(err)
		return
	case ElementGenerator:
		tests = elem.(ElementGenerator).Generate(t)
		nils := make([]int, 0, len(tests)+1)
		for i := range tests {
			if tests[i] == nil {
				nils = append(nils, i)
			}
		}
		if len(nils) > 0 {
			t.Errorf("nil generated Elements %v", nils)
		}
	case Element:
		tests = []Element{elem.(Element)}
	default:
		err = errorf("element does not implement table.T %v", reflect.TypeOf(elem))
		t.Error(err)
		return
	}
	return
}

// Execute test's Test method. If test is an ElementBefore type execute
// test.Before() prior to test.Test(). If test is a ElementAfter type, execute
// test.After() after test.Test() returns. Handles runtimes panics resulting
// from any of these callback.
func elementTest(t T, test Element) {
	if test == nil {
		return
	}
	place := "before"
	defer func() {
		if e := recover(); e != nil {
			t.Errorf("panic %s test; %v", place, e)
		}
	}()
	switch test.(type) {
	case ElementBeforeAfter:
		test.(ElementBefore).Before(subT("before test", t))
		defer test.(ElementAfter).After(subT("after test", t))
	case ElementBefore:
		test.(ElementBefore).Before(subT("before test", t))
	case ElementAfter:
		defer test.(ElementAfter).After(subT("after test", t))
	}
	place = "during"
	defer func() { place = "after" }()
	defer func() {
		switch panicv := recover(); test.(type) {
		case ElementPanics:
			exps, _ := getElementPanicsExpectations(t, test.(ElementPanics))
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
