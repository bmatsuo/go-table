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
func (t *testingT) dup() (cp *testingT)            { cp = new(testingT); *cp = *t; return }
func (t *testingT) sub(name string) (s *testingT)  { s = newTestingT(name, t); return }

func (t *testingT) msg(v ...interface{}) string            {
	m := sprint(v...)
	if t.name != "" {
		m = msg(t.name, m)
	}
	return m
}
func (t *testingT) errmsg(typ string, v ...interface{}) string {
	name, m := t.name, sprint(v...)
	if Verbose {
		prefix := typ
		if name != "" {
			prefix = sprintf("%s %s", t.name, typ)
		}
		name = msg(name, prefix)
	}
	if name != "" {
		m = msg(name, m)
	}
	return m
}
func (t *testingT) msgf(f string, v ...interface{}) string { return msg(t.name, sprintf(f, v...)) }


func (t *testingT) Fail()                                     { t.t.Fail() }
func (t *testingT) FailNow()                                  { t.t.FailNow() }
func (t *testingT) Failed() bool                              { return t.t.Failed() }
func (t *testingT) log(args ...interface{})                   { t.t.Log(sprint(args...)) }
func (t *testingT) error(args ...interface{})                 { t.t.Error(sprint(args...)) }
func (t *testingT) fatal(args ...interface{})                 { t.t.Fatal(sprint(args...)) }
func (t *testingT) Log(args ...interface{})                   { t.log(t.msg(args...)) }
func (t *testingT) Error(args ...interface{})                 { t.error(t.errmsg("error", args...)) }
func (t *testingT) Fatal(args ...interface{})                 { t.fatal(t.errmsg("fatal", args...)) }
func (t *testingT) Logf(format string, args ...interface{})   { t.Log(sprintf(format, args...)) }
func (t *testingT) Errorf(format string, args ...interface{}) { t.Error(sprintf(format, args...)) }
func (t *testingT) Fatalf(format string, args ...interface{}) { t.Fatal(sprintf(format, args...)) }

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
