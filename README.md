# Package 🎁'd error

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror is an `error` type for Go that utilizes the `errors` package's `Unwrap`, `Is` and `As` methods to chain as many errors together as you'd like.

It contains handy methods to examine the error chain, the stack and your source, and plays nicely with other `error` types.

## Features

- [x] Wrap/unwrap errors
- [x] Give errors context
- [x] Examine error chain
  - [x] Walk
  - [x] Depth
  - [x] Prettified trace
- [x] Examine caller
  - [x] File, function and line
  - [x] Stack trace
  - [x] Source code
- [x] Examine process
  - [x] Num routines, CPUs and cgo calls
  - [x] Memory statistics
  - [x] Programmatic breakpoints

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
err := errors.New("some error")
e := we.New(err, "oh no")

fmt.Println(e)
```

Output:

```
oh no: some error
```

An error's context doesn't have to be a string.

```go
myObj := &MyObj{}
if data, err := json.Marshal(myObj); err != nil {
  return we.New(err, myObj)
}
```

### Examining Errors

There are many ways to probe an `Error` for information...

#### Depth

Wrapped errors have _depth_. That is, the number of errors after itself in the error chain.

For eample, the following prints the depth of each error in the chain.

```go
e0 := we.New(nil, "error A")
e1 := we.New(e0, "error B")
e2 := we.New(e1, "error C")

fmt.Printf("e0 depth: %d\n", e0.Depth())
fmt.Printf("e1 depth: %d\n", e1.Depth())
fmt.Printf("e2 depth: %d\n", e2.Depth())
```

```
e0 depth: 0
e1 depth: 1
e2 depth: 2
```

#### Walk

Step through the error chain with the `Walk` method.

```go
e2.Walk(func (err error) bool {
  // Do something with the error.

  if err == ErrSomeParticularType {
    // Don't continue with the walk.
    return false
  }

  if errors.Unwrap(err) == nil {
    // This is the last error in the chain.
    // Do something else.
  }

  // Continue with the walk.
  return true
})
```

#### Trace

Get an error trace by calling the `Trace` method. This method returns a prettified string representation of an error with caller information.

```go
fmt.Println(e2.Trace())
```

```
┌ 2: main.function (main.go:61) error C
├ 1: main.function (main.go:60) error B
└ 0: main.function (main.go:59) error A
```

#### Error and Context

The error's `Error` method returns an inline string representation of the entire error chain.

```go
fmt.Println(e2.Error())
```

```
error C: error B: error A
```

To only examine the receivers context, use the `Context() interface{}` method.

```go
fmt.Printf("%+v", e2.Context())
```

```
error C
```

#### Caller

Errors contain call information accessible from the `Caller` method. 

```go
fmt.Println(e2.Caller())
```

```
main.function (main.go:19)
```

Along with basic file, function and line information, you can use the caller to provide a stack trace of the goroutine the error was created on.

```go
fmt.Println(e.Caller().Stack())
```

```
goroutine 18 [running]:
runtime/debug.Stack(0x0, 0x0, 0x0)
  /usr/local/Cellar/go/1.16/libexec/src/runtime/debug/stack.go:24 +0xa5
github.com/colinc86/wrappederror.currentCaller(0x1, 0x0)
  /Users/colin/Documents/Programming/Go/wrappederror/caller.go:65 +0x45
github.com/colinc86/wrappederror.TestStack(0xc000082600)
  /Users/colin/Documents/Programming/Go/wrappederror/caller_test.go:25 +0x3f
testing.tRunner(0xc000082600, 0x11acff0)
  /usr/local/Cellar/go/1.16/libexec/src/testing/testing.go:1194 +0x1a3
created by testing.(*T).Run
  /usr/local/Cellar/go/1.16/libexec/src/testing/testing.go:1239 +0x63c
```

When debugging, the caller type also collects source code information.

```go
fmt.Println(e2.Caller().Source())
```

```
...
e0 := we.New(nil, "error A")
e1 := we.New(e0, "error B")
e2 := we.New(e1, "error C")

fmt.Printf("e0 depth: %d\n", e0.Depth())
...
```

By default, when possible, the caller collects the immediate two lines above and below the caller. If you want more or less information you can set (and check) the radius with the `SetSourceFragmentRadius(radius int)` and `SourceFragmentRadius() int` functions.

```go
if we.SourceFragmentRadius() != 5 {
  we.SetSourceFragmentRadius(5)
}
```

#### Process

Use the `Process()` method to get information about the current process. Memory statistics are also available with the `e.Process().Memory()` method but don't print in the process's string value.

```go
fmt.Println(e.Process())
```

```
goroutines: 2, cpus: 16, cgos: 0
```

#### Debugging

It is also possible to trigger a breakpoint programatically when an error is received using the `Process` type.

```go
// doSomething returns a wrapped error
if e := doSomething(); e != nil {
  e.Process().Break()
  return we
}
```

By default, all calls to `Process.Break` are ignored. A call to `SetIgnoreBreakpoints(false)` must happen before `Process` types will attempt to break.

```go
if !we.IgnoreBreakpoints() {
  we.SetIgnoreBreakpoints(true) // Ignore all breakpoints
}

e := New(nil, err)
e.Break() // Does nothing

we.SetIgnoreBreakpoints(false) // Break on calls to Break

e.Break() // Time to debug!
```

## Contributing

Feel free to contribute either through reporting issues or submitting pull requests.

Thank you to @GregWWalters for ideas, tips and advice.
