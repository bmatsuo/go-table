package examples

import "table" // "github.com/bmatsuo/go-table/table"

import (
	"fmt"
	"reflect"
	"testing"
)

type sprintfTest struct {
	flag string
	v    interface{}
	out  string
}

func (test sprintfTest) Test(t table.T) {
	if out := fmt.Sprintf(test.flag, test.v); out != test.out {
		t.Errorf("fmt.Sprintf(%q, %v) => %q != %q", test.flag, test.v, out, test.out)
	}
}

func TestSprintf(t *testing.T) {
	table.Test(t, []sprintfTest{
		{"%4s", "foo", " foo"},
		{"%-4s", "foo", "foo "},
		{"%5d", 1337, " 1337"},
		{"%-5d", 1337, "1337 "},
		{"%g", 2.125 + 1.5i, "(2.125+1.5i)"},
		{"%v", 2.125 + 1.5i, "(2.125+1.5i)"},
		{"%v", reflect.TypeOf(t), "*T"}, // error: `... examples.sprintfTest 6: fmt.Sprintf("%v", *testing.T) => "*testing.T" != "*T"`
		{"%v", []int{1, 2}, "[1 2]"},
		{"%v", []string{"1", "2"}, "[1 2]"},
		{"%v", struct{ a, b int }{1, 2}, "{1 2}"},
		{"%v", struct{ a, b string }{"1", "2"}, "{\"1\" \"2\"}"}, // error
		{"%#v", struct{ a, b int }{1, 2}, "struct { a int; b int }{a:1, b:2}"},
	})
}
