# Package 🎁'd error

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror is an over-engineered `error` type for Go that utilizes the `error` interface's `Unwrap() error` method to chain as many errors together as you'd like.

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/wrappederror
```

## Example

```go
package main

import (
	"errors"
	"fmt"

	we "github.com/colinc86/wrappederror"
)

func main() {
	e := functionA()
	if e != nil {
		fmt.Printf("Got error: %s\n\n", e)
		fmt.Printf("Got trace:\n%s\n", e.Trace())
	}
}

func functionA() we.WrappedError {
	return we.New("error A", functionB())
}

func functionB() we.WrappedError {
	return we.New("error B", functionC())
}

func functionC() error {
	return errors.New("error C")
}
```

Output

```
Got error: error A: error B: error C

Got trace:
┌ 2: main.functionA (main.go:19) error A
├ 1: main.functionB (main.go:23) error B
└ 0: error C
```
