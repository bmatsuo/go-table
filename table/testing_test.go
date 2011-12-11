package table

/*  Filename:    testing_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Dec 10 15:09:48 PST 2011
 *  Description: For testing testing.go
 */

import (
	"testing"
	"strings"
	"os"
)

type logItem struct{ v interface{} }

func (item logItem) AsError() os.Error {
	switch item.v.(type) {
	case os.Error:
		return item.v.(os.Error)
	}
	return nil
}

var errFailed = Error("failed")

// Construct with new(fauxT).
type fauxT struct {
	failed bool
	log    []logItem
}

func (t *fauxT) Fail()                                     { t.failed = true }
func (t *fauxT) FailNow()                                  { t.Fail(); panic(errFailed) }
func (t *fauxT) Failed() bool                              { return t.failed }
func (t *fauxT) Log(args ...interface{})                   { t.log = append(t.log, logItem{sprint(args...)}) }
func (t *fauxT) Error(args ...interface{})                 { t.Log(Error(sprint(args...))); t.Fail() }
func (t *fauxT) Fatal(args ...interface{})                 { t.Log(Error(sprint(args...))); t.FailNow() }
func (t *fauxT) Logf(format string, args ...interface{})   { t.Log(sprintf(format, args...)) }
func (t *fauxT) Errorf(format string, args ...interface{}) { t.Error(sprintf(format, args...)) }
func (t *fauxT) Fatalf(format string, args ...interface{}) { t.Fatal(sprintf(format, args...)) }

func catchfailed(e interface{}) {
	switch e.(type) {
	case nil:
		return
	case os.Error:
		if e.(os.Error) == errFailed {
			return
		}
	}
	panic(e)
}
func fauxTest(name string, test func(Testing)) (t *fauxT) {
	t = new(fauxT)
	defer func() { catchfailed(recover()) }()
	test(newTestingT(name, t))
	return
}

// utility type for testing.
type metaTest struct {
	name     string
	test     func(Testing)
	failed   bool
	validate func(*fauxT) os.Error
}

func (test metaTest) Test() os.Error {
	var err os.Error
	t := fauxTest(test.name, test.test)
	if test.failed != t.Failed() {
		success := "success"
		if t.Failed() {
			success = "failure"
		}
		err = Errorf("unexpected %s", success)
	}
	if e := test.validate(t); e != nil {
		if err != nil {
			err = Errorf("\n\t%v\n\t%v", err, e)
		} else {
			err = e
		}
	}
	return err
}

type testingTTest metaTest

func (test testingTTest) Test() os.Error { return metaTest(test).Test() }

func stringContains(name, text, sub string) os.Error {
	if strings.Index(text, sub) < 0 {
		return Errorf("%s missing %#v: %#v", name, sub, text)
	}
	return nil
}

func stringMissing(name, text, sub string) os.Error {
	if strings.Index(text, sub) >= 0 {
		return Errorf("%s unexpected %#v: %#v", name, sub, text)
	}
	return nil
}

func emptyLog(log []logItem) os.Error {
	if len(log) > 0 {
		return Error("non-empty log")
	}
	return nil
}

func sizeLog(log []logItem, size int) os.Error {
	if size == 0 {
		return emptyLog(log)
	}
	if len(log) != size {
		return Errorf("unexpected log size %d != %d", len(log), size)
	}
	return nil
}

var testingTTests = []testingTTest{
	{"testname", func(t Testing) {}, false, func(t *fauxT) os.Error {
		if err := emptyLog(t.log); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Fail() }, true, func(t *fauxT) os.Error {
		if err := emptyLog(t.log); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.FailNow() }, true, func(t *fauxT) os.Error {
		if err := emptyLog(t.log); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Log("logmsg") }, false, func(t *fauxT) os.Error {
		if err := sizeLog(t.log, 1); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "logmsg"); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Error("errmsg") }, true, func(t *fauxT) os.Error {
		if err := sizeLog(t.log, 1); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "errmsg"); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Fatal("fatmsg") }, true, func(t *fauxT) os.Error {
		if err := sizeLog(t.log, 1); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "fatmsg"); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Fatal("fatmsg"); t.Error("errmsg") }, true, func(t *fauxT) os.Error {
		if err := sizeLog(t.log, 1); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "fatmsg"); err != nil {
			return err
		}
		if err := stringMissing("log", sprint(t.log[0].v), "errmsg"); err != nil {
			return err
		}
		return nil
	}},
	{"testname", func(t Testing) { t.Error("errmsg"); t.Log("logmsg") }, true, func(t *fauxT) os.Error {
		if err := sizeLog(t.log, 2); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[0].v), "errmsg"); err != nil {
			return err
		}
		if err := stringMissing("log", sprint(t.log[0].v), "logmsg"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[1].v), "testname"); err != nil {
			return err
		}
		if err := stringContains("log", sprint(t.log[1].v), "logmsg"); err != nil {
			return err
		}
		if err := stringMissing("log", sprint(t.log[1].v), "errmsg"); err != nil {
			return err
		}
		return nil
	}},
}

func TestTestingT(t *testing.T) {
	for i, test := range testingTTests {
		error(t, sprintf("testingT %d", i), test.Test())
	}
}
