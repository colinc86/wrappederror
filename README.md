# Package 🎁'derror

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror implements an `error` type for Go that utilizes the `errors` package's `Unwrap`, `Is` and `As` functions to chain as many errors together as you'd like.

It contains handy methods to examine the error chain, the stack and your source, and plays nicely with other `error` types.

## Features

- 🎁 [Wrap/unwrap errors](#wrapping-errors)
- 📎 [Give errors context](#wrapping-errors)
- 🎛 [Configurable](#configuring-errors)
- 🔍 [Examine errors](#examining-errors)
  - 🗂 [Metadata](#metadata)
  - 📏 [Depth](#depth)
  - 👣 [Walk](#walk)
  - ⛓ [Trace](#trace)
  - 🖇 [Context](#error-and-context)
- 📞 [Caller](#caller)
  - 📄 [File, function and line](#file-function-and-line)
  - 🧬 [Stack trace](#stack-trace)
  - 🧩 [Source fragment](#source-fragment)
- 🔬 [Process](#process)
  - 💻 [Num routines, CPUs and cgo calls](#goroutines-cpus-and-cgo)
  - 📊 [Memory statistics](#memory-statistics)
  - 📌 [Programmatic breakpoints](#debugging)

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/wrappederror
```

Import the package:

```go
import we "github.com/colinc86/wrappederror"
```

## Using

### Wrapping Errors

Use the package's `New(err error, ctx interface{}) Error` function to wrap errors and give them context.

```go
// Get an error
err := errors.New("some error")

// Wrap the error
e := we.New(err, "oh no")

// Print the wrapped error
fmt.Println(e)
```

```
oh no: some error
```

An error's context doesn't have to be a string.

```go
myObj := &MyObj{}
if data, err := json.Marshal(myObj); err != nil {
  // If we failed to marshal myObj, then attach it to the error for context
  return we.New(err, myObj)
}
```

### Examining Errors

There are many ways to probe an error for information...

#### Metadata

Errors come attached with metadata. `Metadata` types contain information about the error that can be useful when debugging such as

- the error's index during the process's execution,
- the number of similar non-nil errors that have been wrapped,
- and the time that the error was created.

```go
// Print the error's metadata
fmt.Println(e.Metadata())
```

```
(#1) (≈0) 2021-03-06 23:26:58.760018 -0600 CST m=+0.000599020
```

The package keeps track of the number of similar errors by keeping a hash map of the errors that have been wrapped. It creates a 128-bit hash of an error's `Error() string` method and keeps a count of the number of identical hashes. You can turn this behavior on/off by using the `SetTrackSimilarErrors(track bool)` function.

#### Depth

Errors have _depth_. That is, the number of errors after itself in the error chain.

For eample, the following prints the depth of each error in the chain.

```go
// Create some errors
e0 := we.New(nil, "error A")
e1 := we.New(e0, "error B")
e2 := we.New(e1, "error C")

// Print their depths
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
  // Do something with the error

  if err == ErrSomeParticularType {
    // Don't continue with the walk
    return false
  }

  if errors.Unwrap(err) == nil {
    // This is the last error in the chain...
    // Do something else
  }

  // Continue with the walk
  return true
})
```

#### Trace

Get an error trace by calling the `Trace() string` method. This method returns a prettified string representation of an error with caller information.

```go
// Print an error trace
fmt.Println(e2.Trace())
```

```
┌ 2: main.function (main.go:61) error C
├ 1: main.function (main.go:60) error B
└ 0: main.function (main.go:59) error A
```

#### Error and Context

The error's `Error() string` method returns an inline string representation of the entire error chain.

```go
// Print the entire error chain
fmt.Println(e2.Error())
```

```
error C: error B: error A
```

To only examine the receiver's context, use the `Context() interface{}` method.

```go
// Only print the error's context
fmt.Printf("%+v", e2.Context())
```

```
error C
```

### Caller

By default, errors contain call information accessible from the `Caller() interface{}` method. See the [Configuring Errors](#configuring-errors) section for more information.

#### File, Function and Line

```go
// Print call information
fmt.Println(e2.Caller())
```

```
main.function (main.go:19)
```

#### Stack Trace

Along with basic file, function and line information, you can use the caller to provide a stack trace of the goroutine the error was created on.

```go
// Print a stack trace
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

#### Source Fragment

When debugging, the caller type also collects source code information.

```go
// Print the source code around the line that the error was created on
fmt.Println(e2.Caller().Source())
```

```
[47-51] /Users/colin/Documents/Programming/Go/wrappederror/wcaller_test.go

func TestCallerSource(t *testing.T) {
	c := currentCaller(1)
	if c.Source() == "" {
		t.Error("Expected a source trace.")

```

By default, when possible, the caller collects the immediate two lines above and below the caller. If you want more or less information you can set (and check) the radius with the `SetSourceFragmentRadius(radius int)` and `SourceFragmentRadius() int` functions.

```go
// If the radius hasn't been set to 5...
if we.SourceFragmentRadius() != 5 {
  // Set the radius to 5
  we.SetSourceFragmentRadius(5)
}
```

### Process

Use the `Process() Process` method to get information about the current process. See the [Configuring Errors](#configuring-errors) section for more information.

#### Goroutines, CPUs and CGO

Process types contain some general process information like the number of current goroutines, the number of available CPUs, and the number of cgo functions executed.

```go
// Print the process information when the error was created
fmt.Println(e.Process())
```

```
goroutines: 2, cpus: 16, cgos: 0
```

#### Memory Statistics

Memory statistics are also available with the `e.Process().Memory() *runtime.MemStats` method.

```go
// Print the allocated memory at the time of the error
fmt.Printf("Allocated memory at %s: %d bytes\n", e.Time(), e.Process().Memory().Alloc)
```

#### Debugging

It is also possible to trigger a breakpoint programatically when an error is received using the `Process` type.

```go
// doSomething returns a wrapped error
if e := doSomething(); e != nil {
  // Initiate a breakpoint
  e.Process().Break()

  // Continue
  return we
}
```

By default, all calls to `Process().Break()` are ignored. A call to `SetIgnoreBreakpoints(false)` must happen before `Process` types will attempt to break.

```go
// Ignore all breakpoints if we aren't debugging
we.SetIgnoreBreakpoints(os.Getenv("DEBUG") != "true")

e := New(nil, err)

// Only attempts to break if the env var DEBUG is "true"
e.Break()
```

### Configuring Errors

Some of the behaviors of new errors can be configured using the follwing table of functions. Only the getters are listed, but setters exist for each.

| Function                     | Initial Value | Description |
|:-----------------------------|:--------------|:------------|
| `CaptureCaller() bool`       | `true`        | Determines whether or not new errors will capture their call information. If you don't need to capture call information, you can set this to `false`. Be advised, future calls to `Caller()` on new errors will return `nil`. |
| `CaptureProcess() bool`      | `true`        | Determines whether or not new errors will capture process information. If you don't need to capture process information, you can set this to `false`. Same as `CaptureCaller`, future calls to `Process()` on new errors will return `nil`. |
| `SourceFragmentRadius() int` | `2`           | The line radius of source fragments collected during debugging. For example, if the error is created on line 15 in a file, then (using the default radius of 2) source would be collected from lines 13 through 17. |
| `IgnoreBreakpoints() bool`   | `true`        | Determines whether or not breakpoints should be ignored when calling `Process().Break()`. |
| `NextErrorIndex() int`       | `1`           | The next index that will be used when creating an error in the error's metadata. |
| `TrackSimilarErrors() bool`  | `true`        | Whether or not errors that are wrapped should be tracked for similarity. |

## Contributing

Feel free to contribute either through reporting issues or submitting pull requests.

Thank you to @GregWWalters for ideas, tips and advice.
