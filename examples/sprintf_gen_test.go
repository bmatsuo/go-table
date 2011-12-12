// See spirntf_test.go for the simpler version of these tests.

package examples

import "table" // "github.com/bmatsuo/go-table/table"

import (
	"reflect"
	"testing"
)

/*****************/
/* Begin testing */
/*****************/

// A table.Element that contains 'complex', constructed type(s).
type sprinterOutputTest struct {
	*sprinter
	out, expected string
}

// This kind of test usually has a constructor
func newSprinterOutputTest(flag string, v interface{}, expected string) *sprinterOutputTest {
	return &sprinterOutputTest{sprinter:newSprinter(flag, v), expected:expected}
}
// Setup extra fields with before test.Test() is called.
func (test *sprinterOutputTest) Before(t table.T) {
	test.out = test.sprinter.Out()
}
// Compare output against expected.
func (test *sprinterOutputTest) Test(t table.T) {
	if test.out != test.expected {
		t.Errorf("%s => %q != %q", test.sprinter.String(), test.out, test.expected)
	}
	// Usually a test this simple doesn't use such a complex setup.
	// Use constructors & callbacks ONLY IF THE MODULARITY GREATLY SIMPLIFIES THINGS.
	// See sprintf_test.go for the simple version of this test.
}

// Like a sprinterOutputTest. But, its a flat struct that is convient for making
// tables with.
type sprinterTest struct {
	flag string
	v    interface{}
	out  string
}

// Create a new complex table.Element to test.
func (test sprinterTest) Generate(t table.T) []table.Element {
	return []table.Element{
		newSprinterOutputTest(test.flag, test.v, test.out),
		// Possible to modularize and generate multiple sub-tests here.
		// ...
	}
}

func TestSprinter(t *testing.T) {
	table.Test(t, []sprinterTest{
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
