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
	"strings"
	"errors"
	"fmt"
)

var (
	Verbose bool       // If true more verbose errors, as well as passed tests, are logged.
	MsgFmt  = "%s: %s" // Result output format. Can be changed to another 2-argument string.
)

/****************************/
/* General helper functions */
/****************************/

// Aliases to "fmt" functions.
func sprint(v ...interface{}) string                 { return fmt.Sprint(v...) }
func sprintf(format string, v ...interface{}) string { return fmt.Sprintf(format, v...) }
func errorf(format string, v ...interface{}) error   { return fmt.Errorf(format, v...) }

// Create a new error using the string representation of v.
func error_(v ...interface{}) error { return errors.New(sprint(v...)) }

/*********************************************/
/* Helper functions for the Test and friends */
/*********************************************/

// Create a named messages with name format MsgFmt.
func msg(name string, v ...interface{}) string {
	if name != "" {
		return sprintf(MsgFmt, name, sprint(v...))
	}
	return sprint(v...)
}
func msgname(name, typ string) string { return strings.Join([]string{name, typ}, " ") }
