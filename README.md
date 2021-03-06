# Package 🎁'd error

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror is an `error` type for Go that utilizes the `errors` package's `Unwrap`, `Is` and `As` methods to chain as many errors together as you'd like.

It contains handy methods to examine the error chain and plays nicely with other `error` types.

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

There are a few ways to probe an `Error` for information...

#### Caller

Access the `Caller()` method to get information about where the error was created.

```go
fmt.Printf("%s: %s", e.Caller(), e.Error())
```

Output:

```
main.function (main.go:19): oh no: some error
```

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

Output:

```
e0 depth: 0
e1 depth: 1
e2 depth: 2
```

#### Walk

Step through the error chain with `Walk`.

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

#### Error

The error's `Error` method returns an inline string representation of the entire error chain.

To only examine the receivers context, use the `Context() interface{}` method.

```go
fmt.Println(e2.Error())
```

Output:

```
error C: error B: error A
```

## Contributing

Feel free to contribute either through reporting issues or submitting pull requests.

Thank you to @GregWWalters for ideas, tips and advice.
