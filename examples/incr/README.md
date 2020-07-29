increment test
=

This is an example of an increment test. This is the example used in the
Goconvey docs, so I use it as an example to show how patterns that people have
seen before may be expressed in tea.

The same tests are written using different testing tools. Build flags are used
to control which of the testing frameworks is being used. The reason we use
build flags here is because some of the examples are not expressible in some
of the other testing frameworks; those examples exist to show how something
might unexpectedly fail in one framework but work in tea.

Files beginning with `std_` are written using only the standard library
testing framework on its own. Files beginning with `convey_` are written using
Goconvey, a BDD framework. Files beginning with `tea_` are written using tea.

- `go test -tags std` runs the standard library tests
- `go test -tags convey` runs the Goconvey tests
- `go test` runs the tea tests. This is the default since we want to give
  examples of how to use tea; the tea tests should be runnable with simply `go test`
