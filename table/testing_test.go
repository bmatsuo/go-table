package table

/*  Filename:    testing_test.go
 *  Author:      Bryan Matsuo <bmatsuo@soe.ucsc.edu>
 *  Created:     Sat Dec 10 15:09:48 PST 2011
 *  Description: For testing testing.go
 */

import (
	"regexp"
	"strings"
	"testing"
)

type logItem struct{ v interface{} }

func (item logItem) AsError() error {
	switch item.v.(type) {
	case error:
		return item.v.(error)
	}
	return nil
}

func (item logItem) String() string { return sprint(item.v) }

var errFailed = error_("failed")

// Construct with new(fauxT).
type fauxT struct {
	failed bool
	log    []logItem
}

func (t *fauxT) Len() int                 { return len(t.log) }
func (t *fauxT) Index(i int) interface{}  { return t.log[i].v }
func (t *fauxT) IndexString(i int) string { return t.log[i].String() }

func (t *fauxT) doLog(fn func(int, string)) {
	for i, item := range t.log {
		fn(i, item.String())
	}
}

func (t *fauxT) logLike(patt string) (lnmatch bool) {
	r := regexp.MustCompile(patt)
	t.doLog(func(i int, ln string) { lnmatch = lnmatch || r.MatchString(ln) })
	return
}

func (t *fauxT) logLineLike(i int, patt string) bool {
	return regexp.MustCompile(patt).MatchString(t.IndexString(i))
}

func (t *fauxT) Fail()                                     { t.failed = true }
func (t *fauxT) FailNow()                                  { t.Fail(); panic(errFailed) }
func (t *fauxT) Failed() bool                              { return t.failed }
func (t *fauxT) Log(args ...interface{})                   { t.log = append(t.log, logItem{sprint(args...)}) }
func (t *fauxT) Error(args ...interface{})                 { t.Log(error_(sprint(args...))); t.Fail() }
func (t *fauxT) Fatal(args ...interface{})                 { t.Log(error_(sprint(args...))); t.FailNow() }
func (t *fauxT) Logf(format string, args ...interface{})   { t.Log(sprintf(format, args...)) }
func (t *fauxT) Errorf(format string, args ...interface{}) { t.Error(sprintf(format, args...)) }
func (t *fauxT) Fatalf(format string, args ...interface{}) { t.Fatal(sprintf(format, args...)) }

func catchfailed(e interface{}) {
	switch e.(type) {
	case nil:
		return
	case error:
		if e.(error) == errFailed {
			return
		}
	}
	panic(e)
}
func fauxTest(name string, test func(T)) (t *fauxT) {
	t = new(fauxT)
	defer func() { catchfailed(recover()) }()
	test(subT(name, t))
	return
}

// utility type for testing.
type metaTest struct {
	name     string
	test     func(T)
	failed   bool
	validate func(T, *fauxT)
}

func (test metaTest) Test(t T) {
	ft := fauxTest(test.name, test.test)
	if test.failed != ft.Failed() {
		success := "success"
		if ft.Failed() {
			success = "failure"
		}
		t.Errorf("unexpected %s: %v", success, ft.log)
	}
	test.validate(t, ft)
}

// A metaTest-like type that only tests error contents, assumes tests only
// fail when error patterns have been supplied, and assumes the test function
// do no logging.
type metaTestSimple struct {
	name string
	test func(T)
	errs []string
}

func (test metaTestSimple) Test(t T) {
	metaTest{"simple meta-test", test.test, len(test.errs) > 0, func(t T, ft *fauxT) {
		if len(test.errs) == 0 {
			return
		}
		for i, patt := range test.errs {
			if !ft.logLike(patt) {
				subT(sprintf("error %d", i), t).Errorf("missing error: %v", patt)
			}
		}
	}}.Test(t)
}

type tTest metaTest

func (test tTest) Test(t T) { metaTest(test).Test(t) }

func stringContains(t T, name, text, sub string) {
	if strings.Index(text, sub) < 0 {
		t.Errorf("%s missing %#v: %#v", name, sub, text)
	}
}

func stringMissing(t T, name, text, sub string) {
	if strings.Index(text, sub) >= 0 {
		t.Errorf("%s unexpected %#v: %#v", name, sub, text)
	}
}

func emptyLog(t T, log []logItem) {
	if len(log) > 0 {
		t.Error("non-empty log")
	}
}

func sizeLog(t T, log []logItem, size int) {
	if size == 0 {
		emptyLog(t, log)
	} else if len(log) != size {
		t.Errorf("unexpected log size %d != %d", len(log), size)
	}
}

var tTests = []tTest{
	{"testname", func(t T) {}, false, func(t T, ft *fauxT) { emptyLog(t, ft.log) }},
	{"testname", func(t T) { t.Fail() }, true, func(t T, ft *fauxT) { emptyLog(t, ft.log) }},
	{"testname", func(t T) { t.FailNow() }, true, func(t T, ft *fauxT) { emptyLog(t, ft.log) }},
	{"testname", func(t T) { t.Log("logmsg") }, false, func(t T, ft *fauxT) {
		sizeLog(t, ft.log, 1)
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "logmsg")
	}},
	{"testname", func(t T) { t.Error("errmsg") }, true, func(t T, ft *fauxT) {
		sizeLog(t, ft.log, 1)
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "errmsg")
	}},
	{"testname", func(t T) { t.Fatal("fatmsg") }, true, func(t T, ft *fauxT) {
		sizeLog(t, ft.log, 1)
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "fatmsg")
	}},
	{"testname", func(t T) { t.Fatal("fatmsg"); t.Error("errmsg") }, true, func(t T, ft *fauxT) {
		sizeLog(t, ft.log, 1)
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "fatmsg")
		stringMissing(t, "log", sprint(ft.log[0].v), "errmsg")
	}},
	{"testname", func(t T) { t.Error("errmsg"); t.Log("logmsg") }, true, func(t T, ft *fauxT) {
		sizeLog(t, ft.log, 2)
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "testname")
		stringContains(t, "log", sprint(ft.log[0].v), "errmsg")
		stringMissing(t, "log", sprint(ft.log[0].v), "logmsg")
		stringContains(t, "log", sprint(ft.log[1].v), "testname")
		stringContains(t, "log", sprint(ft.log[1].v), "logmsg")
		stringMissing(t, "log", sprint(ft.log[1].v), "errmsg")
	}},
}

func TestT(t *testing.T) {
	for i, test := range tTests {
		t.Log("test ", i)
		test.Test(t)
	}
}
