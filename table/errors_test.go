package table

/*  Filename:    errors_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:52:30 PST 2011
 *  Description: For testing errors.go
 */

import (
	"testing"
	"reflect"
	"regexp"
	"os"
)

type errorStringTest struct {
	err     os.Error
	matches []string
}

func (test errorStringTest) Test(t Testing) {
	for i, r := range test.matches {
		if !regexp.MustCompile(r).MatchString(test.err.String()) {
			t.Errorf("pattern %d (%s) doesn't match error string: %v", i, r, test.err)
		}
	}
}

// This table test is a little wonky because I set the Verbose global option for half of the tests.
func TestErrorString(t *testing.T) {
	var errorStringTests = []errorStringTest{
		{os.NewError(errorstr("hello", Errorf("world"))), []string{`hello: world`}},
		{os.NewError(errorstrf("hello", "%q", "world")), []string{`hello: "world"`}},
		{os.NewError(fatalstr("hello", Errorf("world"))), []string{`hello: world`}},
		{os.NewError(fatalstrf("hello", "%q", "world")), []string{`hello: "world"`}},
		{Errorf("this is a %s error %v", "formatted", reflect.TypeOf("abc")), []string{"this is a formatted error string"}},
		{Fatalf("this is a %s error string", "formatted"), []string{"this is a formatted error string"}},
	}
	Verbose = true
	defer func() { Verbose = false }()
	errorStringTests = append(errorStringTests, []errorStringTest{
		{os.NewError(errorstr("hello", Errorf("world"))), []string{`hello error: world`}},
		{os.NewError(errorstrf("hello", "%q", "world")), []string{`hello error: "world"`}},
		{os.NewError(fatalstr("hello", Errorf("world"))), []string{`hello fatal error: world`}},
		{os.NewError(fatalstrf("hello", "%q", "world")), []string{`hello fatal error: "world"`}},
		{Errorf("this is a %s error %v", "formatted", reflect.TypeOf("abc")), []string{"this is a formatted error string"}},
		{Fatalf("this is a %s error string", "formatted"), []string{"this is a formatted error string"}},
	}...)

	for i, test := range errorStringTests {
		tTest(newTestingT(sprintf("errorStringTest %d", i), t), test)
	}
}
