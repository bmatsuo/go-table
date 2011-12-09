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
	"os"
)

// Non-nil values returned by the Test method will cause the table test that
// called Test to fail. A FatalError retured by Test stops the table test.
type T interface {
	Test() os.Error // Execute the test described by the object.
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

func tTest(t T) (err os.Error) {
	place := "before"
	defer func() {
		if e := recover(); e != nil {
			err = Fatalf("panic %s test; %v", place, e)
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
	defer func() {
		if e := recover(); e != nil {
			err = Fatalf("panic during test; %v", e)
		}
	}()
	err = t.Test()
	place = "after"
	return
}
