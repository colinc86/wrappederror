# Package 🎁'd error

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror is an over-engineered `error` type for Go that utilizes the `error` interface's `Unwrap() error` method to chain as many errors together as you'd like.

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/wrappederror
```

## Using

Import the package:

```go
import we "github.com/colinc86/wrappederror"
```

### Wrapping Errors

Use the package's `Error` type to wrap errors and give them context.

```go
e := we.New(err, "oh no")
```

An error's context doesn't have to be a string.

```go
myObj := &MyObj{}
if data, err := json.Marshal(myObj); err != nil {
	return we.New(err, myObj)
}
```

### Examining Errors

There are a few ways to probe an `Error` for information...

#### Caller

Access the `Caller()` method to get information about where the error was created.

```go
fmt.Printf("Error at %s: %s.\n", e.Caller(), e.Error())
```

Output:

```
Error at main.function (main.go:19): oh no
```

#### Depth

Wrapped errors have _depth_. That is, the number of errors after itself in the error chain.

For eample, the following -

```go
e0 := we.New(nil, "error A")
e1 := we.New(e0, "error B")
e2 := we.New(e1, "error C")

fmt.Printf("e0 depth: %d\n", e0.Depth())
fmt.Printf("e1 depth: %d\n", e1.Depth())
fmt.Printf("e2 depth: %d\n", e2.Depth())
```

outputs

```
e0 depth: 0
e1 depth: 1
e2 depth: 2
```

#### Walk

Step through the error chain with `Walk`.

```go
e2.Walk(func (err error) {
	// Do something with the error.

	if errors.Unwrap(err) == nil {
		// This is the last error in the chain.
		// Do something else.
	}
})
```

#### Trace

Get an error trace by calling the `Trace` method. This method returns a prettified string representation of an error.

```go
fmt.Println(e2.Trace())
```

Output:

```
┌ 2: main.function (main.go:61) error C
├ 1: main.function (main.go:60) error B
└ 0: main.function (main.go:59) error A
```

#### String

Calling `String` returns an inline string representation of an error.

```go
fmt.Println(e2.String())
```

Output:

```
error C: error B: error A
```
