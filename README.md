# Package 🎁'd error

[![Go Tests](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/colinc86/wrappederror/actions/workflows/go-test.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/colinc86/wrappederror.svg)](https://pkg.go.dev/github.com/colinc86/wrappederror)

Package wrappederror is an over-engineered `error` type for Go that utilizes the `error` interface's `Unwrap() error` method to chain as many errors together as you'd like. It imports no non-standard library packages, and emits no errors, ensuring you always get _your_ errors. 🎤💧

## Installing

Navigate to your module and execute the following.

```bash
$ go get github.com/colinc86/wrappederror
```

## Examples

### Wrapping Errors

```go
package main

import we "github.com/colinc86/wrappederror"

func main() {
  e := errorChain()
  if e != nil {
    fmt.Printf("Got error: %s\n", e)
  }
}

func errorChain() error {
  return we.New("the error chain", functionA())
}

func functionA() error {
  return we.New("function A error", functionB())
}

func functionB() error {
  return errors.New("function B error")
}
```

Output

```
Got error: the error chain: function A error: function B error
```

### Encoding Wrapped Errors

```go
package main

import (
  "encoding/json"

  we "github.com/colinc86/wrappederror"
)

type ServerError struct {
  err we.WrappedError `json:"error"`
}

func main() {
  var err error
  // ... some process that returns an error "failed"...

  serverError := &ServerError{we.New("unexpected error", err)}
  b, _ := json.Marshal(serverError)
  fmt.Printf(string(b))
}
```

Output

```
{
  "error": "unexpected error: failed"
}
