// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

/*  Filename:    testing.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Dec 10 15:09:48 PST 2011
 *  Description: 
 */

import (
//"testing"
)

type testingT struct {
	name string
	t    Testing
}

func newTestingT(name string, t Testing) *testingT { return &testingT{name, t} }
func (t *testingT) dup() (cp *testingT)            { cp = new(testingT); *cp = *t; return cp }
func (t *testingT) sub(name string) (s *testingT)  { s = t.dup(); s.name = t.msg(name); return s }

func (t *testingT) msg(v ...interface{}) string            { return msg(t.name, sprint(v...)) }
func (t *testingT) msgf(f string, v ...interface{}) string { return msg(t.name, sprintf(f, v...)) }
func (t *testingT) errmsg(typ string, v ...interface{}) string {
	if Verbose {
		return msg(sprintf("%s %s", t.name, typ), sprint(v...))
	}
	return msg(t.name, sprint(v...))
}

func (t *testingT) Fail()                                     { t.t.Fail() }
func (t *testingT) FailNow()                                  { t.t.FailNow() }
func (t *testingT) Failed() bool                              { return t.t.Failed() }
func (t *testingT) Log(args ...interface{})                   { t.t.Log(t.msg(args...)) }
func (t *testingT) Error(args ...interface{})                 { t.t.Error(t.errmsg("error", args...)) }
func (t *testingT) Fatal(args ...interface{})                 { t.t.Fatal(t.errmsg("fatal", args...)) }
func (t *testingT) Logf(format string, args ...interface{})   { t.Log(t.msgf(format, args...)) }
func (t *testingT) Errorf(format string, args ...interface{}) { t.Error(t.msgf(format, args...)) }
func (t *testingT) Fatalf(format string, args ...interface{}) { t.Fatal(t.msgf(format, args...)) }

type Testing interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}
