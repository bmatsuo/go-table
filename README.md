[install go]: http://golang.org/install.html "Install Go"
[gopkgdoc]: http://gopkgdoc.appspot.com/pkg/github.com/bmatsuo/go-table/table "GoPkgDoc"
[the godoc url]: http://localhost:6060/pkg/github.com/bmatsuo/go-table/table/ "the Godoc URL"
[table driven testing in go]: http://code.google.com/p/go-wiki/wiki/TableDrivenTests "table driven testing in Go"
[dry]: http://en.wikipedia.org/wiki/Don't_repeat_yourself "DRY"

Go-table
========

Don't write boiler-plate code. Write test code.

What is Go-table?
-----------------

Go-table provides package "table" to facilitate [table driven testing in Go][].
Table driven tests are not difficult to write once you get the hang of it.
Once that happens, they are a very effective tool for many unit-testing
applications.

Table tests in Go are awesome witout Go-table. Why use it?
----------------------------------------------------------

Table driven testing, combined with Go's flexible anonomous type inference is
already quite clean and conise even without Go-table. Go-table is merely trying
to squeeze a little more awesomeness out of this Go-friendly testing paradigm.

Go-table encourages [DRY][] coding. The code contained in your test files should
only be repetitive in structure (so writing them is easy). Any looping or
logging constructs that are used in *every* test function should be absolutely
minimal, if not handled for you automatically.

When you're writing tests, the *only* code you should be writing is test code.
Go-table tries to help get that done.

Features
--------

- Obviously, easy automation of table tests is a key feature.
- Convenient logging features to automatically name/identify (failed) tests.
- Automatic handling of runtime panics uncaught by test code.
- Custom callbacks that can run before or after individual tests.

Documentation
=============

API Documentation
-----------------

Documentation at [GoPkgDoc][] should be fairly recent.

To view documentation for your local Go-table installation, start a Godoc web server

    godoc -http=:6060

An example
----------

The following re-implements the example table test from [the wiki article
metioned earlier][table driven testing in go].

**Note:** It's important to remember this example is supposed to be from
package fmt, and has access to type flagPrinter, and funcs Sprintf, and Errorf.
Creating the test in the first place would make a cyclic dependency. But it is
good at illustrating the differences between handrolled table tests and the
table package.

```go
type flagtest struct { in, out string }

func (test flagtest) Test() (err os.Error) {
    if s := Sprintf(test.in, new(flagPrinter)); s != test.out {
        err = Errorf("Sprintf(%q, new(flagPrinter)) => %q, want %q", test.in, s, test.out)
    }
    return
}

func TestFlagParser(t *testing.T) {
    table.Test(t, []flagtest{
        {"%a", "[%a]"},
        {"%-a", "[%-a]"},
        {"%+a", "[%+a]"},
        {"%#a", "[%#a]"},
        {"% a", "[% a]"},
        {"%0a", "[%0a]"},
        {"%1.2a", "[%1.2a]"},
        {"%-1.2a", "[%-1.2a]"},
        {"%+1.2a", "[%+1.2a]"},
        {"%-+1.2a", "[%+-1.2a]"},
        {"%-+1.2abc", "[%+-1.2a]bc"},
        {"%-1.2abc", "[%-1.2a]bc"},
    })
}
```

[original example][table driven testing in go]

The example isn't vastly shorter, or less complex than the original, but the
boiler-plate looping is gone and the modularity eliminates the need for any
extra helper functions.

A close inspection also reveals that the test index `i` is no longer included in
the error message, because the table package prepends test index information for
you automatically.

These improvements are small, but eliminating this repitive code from virtually
every unit test written is definitely *a win* in the big picture.

Installation
============

Prerequisites
-------------

[Go][install go].

Installing
----------

    goinstall github.com/bmatsuo/go-table/table

or

    git clone https://github.com/bmatsuo/go-table
    pushd go-table/table
    gomake install
    popd

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
