# Package 🎁'derror

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) ![Go Coverage](https://img.shields.io/badge/Go%20Coverage-68%25-lightgreen.svg?style=flat) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror implements an `error` type in Go for wrapping errors.

It contains handy methods to examine the error chain, stack and your source, and it plays nicely with other `error` types.

- 🎁 [Wrapping Errors](#-wrapping-errors)
- 🔍 [Examining Errors](#-examining-errors)
  - 📏 [Depth](#-depth)
  - 🔗 [Chain](#-chain)
  - 👣 [Walk](#-walk)
  - 🗺 [Trace](#-trace)
  - 🖇 [Error and Context](#-error-and-context)
  - 🗂 [Metadata](#-metadata)
  - 📇 [Caller](#-caller)
    - 📄 [File, Function and Line](#-file-function-and-line)
    - 🧬 [Stack Trace](#-stack-trace)
    - 🧩 [Source Fragments](#-source-fragments)
  - 🔬 [Process](#-process)
    - 💻 [Num Routines, CPUs and CGO](#-goroutines-cpus-and-cgo)
    - 📊 [Memory Statistics](#-memory-statistics)
    - 📌 [Programmatic Breakpoints](#-debugging)
- 🚨 [Severity Detection](#-severity-detection)
- 🧱 [Marshaling Errors](#-marshaling-errors)
- 🗒 [Formatting Errors](#-formatting-errors)
- 🎛 [Configuring Errors](#-configuring-errors)
- 🧵 [Thread Safety](#-thread-safety)

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/wrappederror
```

Import the package:

```go
import we "github.com/colinc86/wrappederror"
```

## 🎁 Wrapping Errors

Use the package's `New` function to wrap errors and give them context.

```go
// Get an error
err := errors.New("some error")

// Wrap the error with some context
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

## 🔍 Examining Errors

There are many ways to examine an error...

### 📏 Depth

Errors have depth. That is, the number of errors, not including itself, in the error chain.

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

### 🔗 Chain

Access the error chain as a flattened slice instead of wrapped errors using the `Chain` method.

```go
// Store the slice [e2, e1, e0] in c
c := e2.Chain()
```

Optionally, directly access an error with a given depth or index.

```go
// Get the last error in the chain
errA := e2.ErrorWithDepth(0)

// Get the first error in the chain
errB := e2.ErrorWithIndex(0)

// Gets nil
errC := e2.ErrorWithDepth(3)
errD := e2.ErrorWithIndex(-1)
```

### 👣 Walk

Step through the error chain with the `Walk` method. `Walk` calls the step function for every error in the chain until either the last error unwraps to `nil`, or the step function returns `false`.

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

### 🗺 Trace

Get an error trace by calling the `Trace` method. This method returns a prettified string representation of the error chain with caller information. Errors in the chain not defined by this package log their depth and result of calling `Error`.

```go
// Print an error trace
fmt.Println(e2.Trace())
```

```
┌ 2: main.function (main.go:61) error C
├ 1: main.function (main.go:60) error B
└ 0: main.function (main.go:59) error A
```

### 🖇 Error and Context

The error's `Error` method returns an inline string representation of the entire error chain with each component separated by the characters `: ` (colon, space).

```go
// Print the entire error chain
fmt.Println(e2.Error())
```

```
error C: error B: error A
```

To only examine the receiver's context, use the `Context` method.

```go
// Only print the error's context
fmt.Printf("%+v", e2.Context())
```

```
error C
```

### 🗂 Metadata

Errors come attached with metadata. `Metadata` types contain information about the error that can be useful when debugging such as

- the severity of the error, if enabled, (see [Severity Detection](#-severity-detection)),
- the error's index during the process's execution created by this package,
- the number of similar non-nil errors that have been wrapped,
- the duration since the process was launched and when the error was created,
- and the time that the error was created.

```go
// Print the error's metadata
fmt.Println(e.Metadata)
```

```
[moderate] Network Timeout (#1) (≈0) (+10.000280) 2021-03-07 13:29:07.179446 -0600 CST m=+10.000589560
```

The package keeps track of the number of similar errors by keeping a hash map of the errors that have been wrapped. It creates a 128-bit hash of an error's `Error` method and keeps a count of the number of identical hashes. You can turn this behavior on/off by using the `SetTrackSimilarErrors` configuration method.

### 📇 Caller

Errors capture call information accessible through the `Caller` property. Examine information such as code metadata, a stack trace and source fragment.

#### 📄 File, Function and Line

```go
// Print call information
fmt.Println(e2.Caller)
```

```
main.function (main.go:19)
```

#### 🧬 Stack Trace

Along with basic file, function and line information, you can use the caller to provide a stack trace of the goroutine the error was created on.

```go
// Print a stack trace
fmt.Println(e.Caller.StackTrace)
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

#### 🧩 Source Fragments

When possible, and permitted, the caller type also captures source code information.

```go
// Print the source code around the line that the error was created on
fmt.Println(e2.Caller.Fragment)
```

```
[47-51] /Users/colin/Documents/Programming/Go/wrappederror/wcaller_test.go

func TestCallerSource(t *testing.T) {
	c := currentCaller(1)
	if c.Source() == "" {
		t.Error("Expected a source trace.")

```

The caller collects the immediate two lines above to two lines below the calling line. If you want more or less information you can set (and check) the radius with the `SetSourceFragmentRadius` and `SourceFragmentRadius` functions. You can also turn the feature off altogether with `SetCollectSourceFragments`.

```go
// If the radius hasn't been set to 5...
if we.Config().SourceFragmentRadius() != 5 {
  // Set the radius to 5
  we.Config().SetSourceFragmentRadius(5)
}
```

### 🔬 Process

Use the error's `Process` property to get information about the current process at the time the error was created.

#### 💻 Goroutines, CPUs and CGO

Process types contain some general process information like the number of current goroutines, the number of available CPUs, and the number of cgo functions executed.

```go
// Print the process information when the error was created
fmt.Println(e.Process)
```

```
goroutines: 2, cpus: 16, cgos: 0
```

#### 📊 Memory Statistics

Memory statistics are available with the `e.Process.Memory` property.

```go
// Print the allocated memory at the time of the error
fmt.Printf("Allocated memory at %s: %d bytes\n", e.Metadata.Time, e.Process.Memory.Alloc)
```

#### 📌 Debugging

It is also possible to trigger a breakpoint programatically when an error is received using the `Process` type.

```go
// doSomething returns a wrapped error
if e := doSomething(); e != nil {
  // Initiate a breakpoint
  e.Process.Break()

  // Continue
  return we
}
```

All calls to `Process.Break()` are ignored by default. A call to the configuration's `SetIgnoreBreakpoints` with a value of `false` must happen before `Process` types will attempt to break.

```go
// Ignore all breakpoints if we aren't debugging
we.Config().SetIgnoreBreakpoints(os.Getenv("DEBUG") != "true")

e := New(nil, err)

// Only attempts to break if the env var DEBUG is "true"
e.Break()
```

## 🚨 Severity Detection

The package can detect the severity of newly wrapped errors using a table of registered `ErrorSeverity` types. The package matches the severity's regular expression against the output of each error's `Error` method in the error chain. A score in the interval [0.0, 1.0] is calculated by calculating the ratio of the number of matched characters in the string to the total number of characters in the string.

To register a new error severity, first create a new instance of the structure such that no error is returned.

```go
// Create two error severities
s1, err := we.NewErrorSeverity("Network Timeout", "i/o timeout", we.ErrorSeverityLevelModerate)
if err != nil {
  fmt.Printf("Invalid regex: %s\n", err)
}

s2, err := we.NewErrorSeverity("🚨", "fail", we.ErrorSeverityLevelHigh)
if err != nil {
  fmt.Printf("Invalid regex: %s\n", err)
}
```

and then register the error severity with the package.

```go
if err := we.RegisterErrorSeverity(s1); err != nil {
  fmt.Printf("Unable to register error severity: %s\n", err)
}

if err := we.RegisterErrorSeverity(s2); err != nil {
  fmt.Printf("Unable to register error severity: %s\n", err)
}
```

Now, when new errors are created, they will be matched against the registered error severities and the error's `Metadata.Severity` may contain a non-nil value.

```go
// Got a network timeout
e1 := errors.New("dial tcp 0.0.0.0:3000: i/o timeout")
e2 := New(e1, "get request failed")
if e2.Metadata.Severity != nil {
  fmt.Println(e2.Metadata.Severity)
}

// Got an error saving a file
e3 := errors.New("save failed because file does not exist")
e4 := New(e3, "unable to save file")
e5 := New(e4, "file error")
if e5.Metadata.Severity != nil {
  fmt.Println(e5.Metadata.Severity)
}
```

```
[moderate] Network Timeout
[high] 🚨
```

To unregister error severities, call the `UnregisterErrorSeverity` function with the severity you want to unregister.

```go
we.UnregisterErrorSeverity(s1)
we.UnregisterErrorSeverity(s2)
```

The available `ErrorSeverityLevel` constants are

| Level                        |
|:-----------------------------|
| `ErrorSeverityLevelNone`     |
| `ErrorSeverityLevelLow`      |
| `ErrorSeverityLevelModerate` |
| `ErrorSeverityLevelHigh`     |
| `ErrorSeverityLevelSevere`   |

## 🧱 Marshaling Errors

The package supports marshaling errors into JSON, but because the error type defined in this package wraps errors of type `error`, a bijective `UnmarshalJSON` method isn't possible. Intead of attempting to guess at wrapped types, the package just doesn't try.

The types `Caller`, `Process`, `Metadata` and `ErrorSeverity` _do_ implement both JSON marshaling and unmarshaling.

The error chain can get long, and if errors are collecting caller and process information, then JSON objects for a "single" top-level error may be disproportionately large compared to the rest of the JSON object they're embedded in. The package provides a method for determining how errors are marshaled in to JSON data.

```go
// Marshal full JSON objects
we.Config().SetMarshalMinimalJSON(false)

// Marshal a slimmed-down version of errors
we.Config().SetMarshalMinimalJSON(true)
```

The package marshals its error type in to one of two versions of JSON (defined by the `MarshalMinimalJSON` configuration value):

```jsonc
// The "full" version of an error
{
  "context": "the error's context",
  "depth": 0,
  "wraps": { /* another error or null */ },
  "caller": {
    "file": "/path/to/file",
    "function": "function",
    "line": 0,
    "stackTrace": "trace",
    "sourceFragment:" "source"
  },
  "process": {
    "goroutines": 1,
    "cpus": 8,
    "cgos": 0,
    "memory": { /* runtime.MemStats */ },
  },
  "metadata": {
    "time": "the time",
    "duration": 0.0,
    "index": 0,
    "similar": 0
  }
}
```

```jsonc
// The "minimal" version of an error
{
  "context": "the error's context",
  "depth": 0,
  "wraps": { /* another error or null */ },
  "time": "the time",
  "duration": 0.0,
  "index": 0,
  "similar": 0,
  "file": "/path/to/file",
  "function": "function",
  "line": 0
}
```

All other errors are marshaled in to a generic JSON object:

```jsonc
// The "generic" version of an error
{
  "error": "the output of Error() string",
  "wraps": { /* another error or null */ }
}
```

## 🗒 Formatting Errors

Errors have a `Format` method that returns a string with a custom format. It takes an error format string, `ef`, that is built using error format tokens.

For example, you can achieve the same output as the caller's description by using the following format,

```go
ef := fmt.Sprintf(
  "%s (%s:%s)",
  ErrorFormatTokenFunction,
  ErrorFormatTokenFile,
  ErrorFormatTokenLine,
)

// The following statments have the same output
fmt.Println(e.Format(ef))
fmt.Println(e.Caller)
```

or you can create more complex/custom error formats.

```go
ef := fmt.Sprintf(
  "Error #%s at %s (%s:%s): %s",
  ErrorFormatTokenIndex,
  ErrorFormatTokenTime,
  ErrorFormatTokenFile,
  ErrorFormatTokenLine,
  ErrorFormatTokenChain,
)

fmt.Println(e.Format(ef))
```

```
Error #2 at 2021-03-07 16:39:56.393366 -0600 CST m=+0.001129674 (formatter_test.go:10): error 2: error 1
```

The available tokens are as follows.

| Token                           | Description |
|:--------------------------------|:------------|
| `ErrorFormatTokenContext`       | The error's context. |
| `ErrorFormatTokenInner`         | The output of the inner error's `Error` method. |
| `ErrorFormatTokenChain`         | The error chain as returned by the error's `Error` method. |
| `ErrorFormatTokenFile`          | The file name from the error's caller. |
| `ErrorFormatTokenFunction`      | The function from the error's caller. |
| `ErrorFormatTokenLine`          | The line number from the error's caller. |
| `ErrorFormatTokenStack`         | The stack trace from the error's caller. |
| `ErrorFormatTokenSource`        | The source fragment from the error's caller. |
| `ErrorFormatTokenTime`          | The time from the error's metadata. |
| `ErrorFormatTokenDuration`      | The duration (in seconds) from the error's metadata. |
| `ErrorFormatTokenIndex`         | The error's index. |
| `ErrorFormatTokenSimilar`       | The number of similar errors. |
| `ErrorFormatTokenRoutines`      | The number of goroutines when the error was created. |
| `ErrorFormatTokenCPUs`          | The number of available CPUs when the error was created. |
| `ErrorFormatTokenCGO`           | The number of cgo calls when the error was created. |
| `ErrorFormatTokenMemory`        | The process memory statistics when the error was created. |
| `ErrorFormatTokenSeverityTitle` | The detected error severity title. |
| `ErrorFormatTokenSeverityLevel` | The detected error severity level. |

## 🎛 Configuring Errors

The package's configuration is accessible through the global `Config` function.

Use the `Set` method to configure everything at once, or any of the corresponding setters to the getters listed in the table below.

```go
// Configure the package to
// - Capture call information
// - Ignore process information
// - Collect source fragments
// - Get 9 lines of source
// - Ignore breakpoints
// - Start indexing errors at 1
// - Track similar errors
// - Marshal full errors in to JSON
we.Config().Set(true, false, true, 4, true, 1, true, false)
```

To return to the initial state upon launch, use the `ResetState` function. Resetting state resets all configuration variables, the process launch time, and the hash map for keeping track for similar wrapped errors.

```go
// Reset the package's state and configuration.
we.ResetState()
```

| Function                     | Initial Value | Description |
|:-----------------------------|:--------------|:------------|
| `CaptureCaller() bool`       | `true`        | Determines whether or not new errors will capture their call information. If you don't need to capture call information, you can set this to `false`. Be advised, future calls to `Caller` on new errors will return `nil`. |
| `CaptureProcess() bool`      | `true`        | Determines whether or not new errors will capture process information. If you don't need to capture process information, you can set this to `false`. Same as `CaptureCaller`, future calls to `Process` on new errors will return `nil`. |
| `CaptureSourceFragments`     | `true`        | Determines whether or not new errors will capture source code around the line that the error was created on. |
| `SourceFragmentRadius() int` | `2`           | The line radius of source fragments collected during debugging. For example, if the error is created on line 15 in a file, then (using the default radius of 2) source would be collected from lines 13 through 17. |
| `IgnoreBreakpoints() bool`   | `true`        | Determines whether or not breakpoints should be ignored when calling `Process.Break`. |
| `NextErrorIndex() int`       | `1`           | The next index that will be used when creating an error in the error's metadata. |
| `TrackSimilarErrors() bool`  | `true`        | Whether or not errors that are wrapped should be tracked for similarity. |
| `MarshalMinimalJSON() bool`  | `true`        | Determines how errors are marshaled in to JSON. When this value is true, a smaller JSON object is created without size-inflating data like stack traces and source fragments. |

## 🧵 Thread Safety

The package was built with thread-safety in mind. You can modify configuration settings and create errors from any goroutine without worrying about locks.

## 👥 Contributing

Feel free to contribute either through reporting issues or submitting pull requests.

Thank you to @GregWWalters for ideas, tips and advice.
