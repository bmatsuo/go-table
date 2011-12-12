package table

/*  Filename:    errors_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:52:30 PST 2011
 *  Description: For testing errors.go
 */

import (
	"reflect"
	"regexp"
	"testing"
	"os"
)

// Functions to generate error strings.
func esprint(name string, v ...interface{}) (err string) {
	if Verbose {
		name = msgname(name, "error")
	}
	err = msg(name, v...)
	return
}
func fsprint(name string, v ...interface{}) string {
	if Verbose {
		name = msgname(name, "fatal")
	}
	return esprint(name, v...)
}
// Functions to generate error strings with formatted messages.
func esprintf(name, format string, v ...interface{}) string {
	return esprint(name, sprintf(format, v...))
}
func fsprintf(name, format string, v ...interface{}) string {
	return fsprint(name, sprintf(format, v...))
}

type errorStringTest struct {
	err     os.Error
	matches []string
}

func (test errorStringTest) Test(t T) {
	for i, r := range test.matches {
		if !regexp.MustCompile(r).MatchString(test.err.String()) {
			t.Errorf("pattern %d (%s) doesn't match error string: %v", i, r, test.err)
		}
	}
}

// This table test is a little wonky because I set the Verbose global option for half of the tests.
func TestErrorString(t *testing.T) {
	var errorStringTests = []errorStringTest{
		{error_(esprint("hello", errorf("world"))), []string{`hello: world`}},
		{error_(esprintf("hello", "%q", "world")), []string{`hello: "world"`}},
		{error_(fsprint("hello", errorf("world"))), []string{`hello: world`}},
		{error_(fsprintf("hello", "%q", "world")), []string{`hello: "world"`}},
		{errorf("this is a %s error %v", "formatted", reflect.TypeOf("abc")), []string{"this is a formatted error string"}},
		//{fatalf("this is a %s error string", "formatted"), []string{"this is a formatted error string"}},
	}
	Verbose = true
	defer func() { Verbose = false }()
	errorStringTests = append(errorStringTests, []errorStringTest{
		{error_(esprint("hello", errorf("world"))), []string{`hello error: world`}},
		{error_(esprintf("hello", "%q", "world")), []string{`hello error: "world"`}},
		{error_(fsprint("hello", errorf("world"))), []string{`hello fatal error: world`}},
		{error_(fsprintf("hello", "%q", "world")), []string{`hello fatal error: "world"`}},
		{errorf("this is a %s error %v", "formatted", reflect.TypeOf("abc")), []string{"this is a formatted error string"}},
		//{fatalf("this is a %s error string", "formatted"), []string{"this is a formatted error string"}},
	}...)

	for i, test := range errorStringTests {
		elementTest(subT(sprintf("errorStringTest %d", i), t), test)
	}
}
