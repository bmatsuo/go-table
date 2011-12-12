// Copyright 2011, Bryan Matsuo. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package table

/*  Filename:    testing.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Dec 10 15:09:48 PST 2011
 *  Description: 
 */

import ()

type testingT struct {
	name string
	t    T
}

func subT(name string, t T) *testingT       { return &testingT{name, t} }
func (t *testingT) dup() (cp *testingT)           { cp = new(testingT); *cp = *t; return }
func (t *testingT) sub(name string) (s *testingT) { s = subT(name, t); return }

func (t *testingT) msg(v ...interface{}) (m string) {
	if m = sprint(v...); t.name != "" {
		m = msg(t.name, m)
	}
	return
}
func (t *testingT) errmsg(typ string, v ...interface{}) string {
	name, m := t.name, sprint(v...)
	switch {
	case Verbose && name != "":
		typ = sprintf("%s %s", t.name, typ)
		fallthrough
	case Verbose:
		name = msg(name, typ)
		fallthrough
	case name != "":
		m = msg(name, m)
	}
	return m
}
func (t *testingT) msgf(f string, v ...interface{}) string { return t.msg(sprintf(f, v...)) }

func (t *testingT) Fail()                                     { t.t.Fail() }
func (t *testingT) FailNow()                                  { t.t.FailNow() }
func (t *testingT) Failed() bool                              { return t.t.Failed() }
func (t *testingT) log(args ...interface{})                   { t.t.Log(sprint(args...)) }
func (t *testingT) error(args ...interface{})                 { t.t.Error(sprint(args...)) }
func (t *testingT) fatal(args ...interface{})                 { t.t.Fatal(sprint(args...)) }
func (t *testingT) Log(args ...interface{})                   { t.log(t.msg(args...)) }
func (t *testingT) Error(args ...interface{})                 { t.error(t.errmsg("error", args...)) }
func (t *testingT) Fatal(args ...interface{})                 { t.fatal(t.errmsg("fatal", args...)) }
func (t *testingT) Logf(format string, args ...interface{})   { t.log(t.msgf(format, args...)) }
func (t *testingT) Errorf(format string, args ...interface{}) { t.error(t.msgf(format, args...)) }
func (t *testingT) Fatalf(format string, args ...interface{}) { t.fatal(t.msgf(format, args...)) }

// Think *testing.T
type T interface {
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
