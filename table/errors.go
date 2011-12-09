// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

/*  Filename:    errors.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Fri Dec  9 08:52:30 PST 2011
 *  Description: 
 */

import (
	"testing"
	"fmt"
	"os"
)

var (
	Verbose bool       // If true more verbose errors, as well as passed tests, are logged.
	MsgFmt  = "%s: %s" // Result output format. Can be changed to another 2-argument string.
)

// An error which can be returned by a T's Test method if the test was skipped.
var ErrSkip = os.NewError("skip")

/****************************/
/* General helper functions */
/****************************/

// Aliases to "fmt" functions.
func sprint(v interface{}) string                     { return fmt.Sprint(v) }
func sprintf(format string, v ...interface{}) string  { return fmt.Sprintf(format, v...) }
func errorf(format string, v ...interface{}) os.Error { return fmt.Errorf(format, v...) }

/***********************************************/
/* API functions for Test method return values */
/***********************************************/

// Create a new os.Error using the string representation of v.
func Error(v interface{}) os.Error { return errorf("%v", v) }
// Create a new os.Error with a formatted string. Alias of fmt.Errorf.
func Errorf(format string, v ...interface{}) os.Error { return errorf(format, v...) }

// An error that causes a table test to issue a fatal "testing" error to package
// "testing".
type FatalError struct{ os.Error }

func (f FatalError) String() string { return f.Error.String() }

// Create a new FatalError using the string representation of v.
func Fatal(v interface{}) FatalError                    { return FatalError{errorf("%v", v)} }
// Like Fatal, but uses a formatted error message.
func Fatalf(format string, v ...interface{}) FatalError { return FatalError{errorf(format, v...)} }

/*********************************************/
/* Helper functions for the Test and friends */
/*********************************************/

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
		err = msg(sprintf("%s error", name), v)
	} else {
		err = msg(name, v)
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
func error(t *testing.T, name string, err os.Error) bool {
	return onerror(err, func(err os.Error) {
		if err == ErrSkip {
			if Verbose {
				t.Logf("%s skipped", name)
			}
			return
		}
		switch err.(type) {
		case FatalError:
			fatal(t, name, err)
		}
		t.Error(errorstr(name, err))
	})
}
func fatal(t *testing.T, name string, err os.Error) bool {
	return onerror(err, func(err os.Error) { t.Error(fatalstr(name, err)) })
}
