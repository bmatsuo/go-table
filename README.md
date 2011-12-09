
[install go]: http://golang.org/install.html "Install Go"
[the godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/go-table/table/ "the Godoc URL"
[table driven testing in go]: http://code.google.com/p/go-wiki/wiki/TableDrivenTests "table driven testing in Go"
[dry]: http://en.wikipedia.org/wiki/Don't_repeat_yourself "DRY"

About Go-Table
==============

Go-Table is a package to facilitate [table driven testing in go][]. Table
driven tests are not difficult to write once you get the have of them. But,
they're extremely quite effective once writtren.

Tests with Go-Table are implemented as methods of a type defined in the test
file. The generally reduce the actual Test&ast; function used by "testing" to
a single line.

Go-Tables encourages [DRY][] testing. The code contained in test files should
only be repetitive in structure (so that writing it is easy). Any
looping/logging constructs that are used in *every* test function should be
absolutely minimal. Table driven tests, combined with Go's relatively flexible
anonomous type inference can create a lightweight, powerful, and minimalistic
unit testing environment. Go-Tables aims to exploit these benefits while
providing an interface with robust logging and error handling.

Documentation
=============

Below, we see the contents of a file, `mylib_test.go`.

```go
package mylib
import (
    "github.com/bmatsuo/go-table/table"
    "testing"
    "fmt"
    "os"
)

type Point struct { x, y int64 }

type (p1 Point) Equals(p2 Point) bool { return p1.x == p2.x && p1.y == p2.y }
type (p1 Point) Plus(p2 Point) Point { return Point{p1.x+p2.x, p1.y+p2.y} }

type PointEqualsTest struct {
    p1, p2   Point
    areEqual bool
}

func (test PointEqualsTest) Test() os.Error {
    if equal := test.p1.Equals(test.p2); equal != test.areEqual {
        return table.Fatal(fmt.Errorf("%v.Equals(%v) returned %v", test.p1, test.p2, test.areEqual))
    }
    return nil
}

var PointEqualsTests = []PointEqualsTest{
    {Point{0,0}, Point{0,1}, false},
    {Point{0,1}, Point{0,1}, true},
    {Point{2,1}, Point{1,1}, false},
    {Point{2,2}, Point{2,2}, true},
}

func TestPointEquals(T *testing.T) { table.Test("Point.Equals", PointEqualsTests...) }

func (test PointEqualsTest)

type PointPlusTest struct { p1, p2, sum Point }

func (test PointPlusTest) Test() os.Error {
    if sum := test.p1.Plus(test.p2); !sum.Equals(test.sum) {
        return fmt.Errorf("%v + %v = %v != %v", test.p1, test.p2, sum, test.sum)
    }
    return nil
}

var PointPlusTests = []PointPlusTest{
    {Point{0,0}, Point{0,1}, Point{0,1}},
    {Point{1,0}, Point{0,1}, Point{1,1}},
    {Point{1,3}, Point{2,1}, Point{3,4}},
}

func TestPointPlus(T *testing.T) { table.Test("Point.Plus", PointPlusTests...) }
```

Above, unit tests have been implemented for the Equals and Plus methods of the
Point type. But this is fairly obvious given the naming convensions.

Prerequisites
-------------

[Install Go][].

Installation
-------------

Use goinstall to install go-table

    goinstall github.com/bmatsuo/go-table/table

General Documentation
---------------------

Use godoc to vew the documentation for go-table

    godoc github.com/bmatsuo/go-table/table

Or alternatively, use a godoc http server

    godoc -http=:6060

and visit [the Godoc URL][]


Author
======

Bryan Matsuo &lt;bmatsuo@soe.ucsc.edu&gt;

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.
Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.
