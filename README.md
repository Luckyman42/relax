# Relax

[![CI](https://img.shields.io/github/actions/workflow/status/luckyman42/relax/ci.yml?branch=main\&label=ci)](https://github.com/luckyman42/relax/actions)
[![codecov](https://codecov.io/gh/luckyman42/relax/branch/main/graph/badge.svg)](https://codecov.io/gh/luckyman42/relax)
[![Go Report Card](https://goreportcard.com/badge/github.com/luckyman42/relax)](https://goreportcard.com/report/github.com/luckyman42/relax)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/luckyman42/relax)](https://pkg.go.dev/github.com/luckyman42/relax)
![Go Version](https://img.shields.io/github/go-mod/go-version/luckyman42/relax)
[![License](https://img.shields.io/github/license/luckyman42/relax)](LICENSE)

> Don’t panic — relax.

`relax` is a small, focused Go library for structured panic-based error propagation inside trusted internal call chains.

It reduces repetitive `if err != nil { return err }` forwarding while preserving:

* the original error
* a captured stack trace
* failure timestamp
* structured contextual metadata

Unlike generic panic wrappers, `relax` keeps failure propagation explicit and typed through the `Failer` type.

---

# Why?

Go’s explicit error handling is excellent at API boundaries.

But inside deep internal call chains, a large amount of code often becomes repetitive forwarding boilerplate:

```go
func service() error {
    data, err := load()
    if err != nil {
        return err
    }

    processed, err := transform(data)
    if err != nil {
        return err
    }

    err = save(processed)
    if err != nil {
        return err
    }

    return nil
}
```

`relax` provides an alternative for *internal-only* propagation:

```go
func service() {
    data := relax.FailCheck(load())
    processed := relax.FailCheck(transform(data))
    err := save(processed)
    if err != nil {
        relax.FailWith(err)
    }
}
```

Then recover once at a boundary:

```go
err := relax.Guard(service)
//or 
err := relax.Guard(func(){
  doSomething(42)
})
```

This keeps the happy-path readable while still preserving structured debugging information.

---

# Core Concepts

## `Failer`

`Failer` is the central error type used by the library.

```go
type Failer struct {
    Err       error
    Stack     []byte
    Timestamp time.Time
    Context   map[string]any
}
```

A `Failer`:

* implements `error`
* preserves the original error
* captures a stack trace at failure creation
* stores optional structured metadata
* interoperates with `errors.Is` and `errors.As`

---

## `FailWith`

`FailWith` converts an error into a `Failer` and panics with it.

```go
relax.FailWith(err)
```

Additional contextual metadata can be attached:

```go
relax.FailWith(err,
    "user_id", userID,
    "operation", "create_order",
)
```

If the error already is a `Failer`, context is merged instead of double-wrapping.

---

## `FailCheck`

`FailCheck` eliminates repetitive forwarding.

```go
value := relax.FailCheck(load())
```

Equivalent to:

```go
value, err := load()
if err != nil {
    relax.FailWith(err)
}
```

---

## `Guard*`

`Guard`, `GuardValue`, `GuardErr`, and `GuardResult` define explicit recovery boundaries.

They:

* recover only `Failer` panics
* convert them back into normal Go errors
* re-panic every other panic unchanged

This distinction is important.

Programmer bugs such as:

* nil dereferences
* out-of-bounds access
* invariant violations

still crash normally instead of being silently converted into errors.

---

# Installation

```bash
go get github.com/luckyman42/relax
```

```go
import "github.com/luckyman42/relax"
```

---

# Quick Example

```go
package main

import (
    "errors"
    "fmt"
    "log"

    "github.com/luckyman42/relax"
)

func fetchUser(id int) (string, error) {
    return "", errors.New("database unavailable")
}

func loadProfile(id int) string {
    return relax.FailCheck(fetchUser(id))
}

func main() {
    profile, err := relax.GuardValue(func() string {
        return loadProfile(42)
    })

    if err != nil {
        var failer relax.Failer
        if errors.As(err, &failer) {
            log.Printf("error: %v", failer.Err)
            log.Printf("context: %+v", failer.Context)
            log.Printf("timestamp: %v", failer.Timestamp)
        }
        return
    }

    fmt.Println(profile)
}
```

---

# API Overview

## Failure Propagation

### `FailWith`

```go
func FailWith(err error, keyVals ...any)
```

Panics with a `Failer`.

If `err` already is a `Failer`, it is re-thrown without double wrapping.

---

### `FailCheck`

```go
func FailCheck[T any](v T, err error) T
```

Returns `v` if `err == nil`, otherwise calls `FailWith`.

---

## Recovery Boundaries

### `Guard`

```go
func Guard(fn func()) error
```

Executes a function and converts `Failer` panics into errors.

---

### `GuardValue`

```go
func GuardValue[T any](fn func() T) (T, error)
```

Like `Guard`, but returns a value.

---

### `GuardErr`

```go
func GuardErr(fn func() error) error
```

Protects functions already returning `error`.

---

### `GuardResult`

```go
func GuardResult[T any](fn func() (T, error)) (T, error)
```

Protects `(T, error)` style functions.

---

## Utilities

### `ConvertToFailer`

```go
func ConvertToFailer(err error) Failer
```

Converts any error into a `Failer`.

Existing `Failer` values are preserved.

---

### `IsFailer`

```go
func IsFailer(err error) bool
```

Reports whether an error is or wraps a `Failer`.

---

# Design Philosophy

`relax` intentionally keeps the model small.

There is:

* no framework
* no dependency injection
* no custom runtime
* no logging abstraction
* no reflection-heavy machinery

The library focuses on one specific problem:

> reducing internal error-forwarding boilerplate while preserving structured debugging information.

The implementation follows several strict design rules:

* Only `Failer` panics are recovered.
* All non-`Failer` panics propagate normally.
* Errors are never silently swallowed.
* Stack traces are captured at failure creation.
* Existing `Failer` values are never double-wrapped.
* Standard `errors.Is` / `errors.As` interoperability is preserved.

This keeps the behavior explicit, predictable, and easy to audit.

---

# Recommended Usage

`relax` works best in:

* service-layer orchestration
* request pipelines
* CLI execution flows
* internal business logic
* goroutine entry points
* worker/job processing systems

A common pattern is:

1. Use `FailCheck` and `FailWith` internally.
2. Recover once at a clear boundary with `Guard*`.
3. Log or translate the returned error.

---

# When NOT to Use Relax

`relax` is intentionally opinionated.

Avoid it for:

* exported public APIs
* low-level libraries consumed by others
* hot performance-critical paths
* normal control flow
* cases where explicit error returns improve clarity

If a package is expected to behave like standard Go libraries, explicit `error` returns are usually the better choice.

---

# Error Context

Structured metadata can be attached to failures:

```go
relax.FailWith(err,
    "request_id", requestID,
    "user_id", userID,
    "operation", "payment_capture",
)
```

This context is preserved across propagation and accessible from the recovered `Failer`.

---

# Stack Traces

Each newly created `Failer` captures a stack trace using:

```go
runtime/debug.Stack()
```

This makes failures significantly easier to debug in concurrent systems and asynchronous execution flows.

---

# Panic Safety

`relax` does **not** attempt to hide programmer errors.

The following still panic normally:

* nil pointer dereferences
* index out of range
* type assertion failures
* explicit `panic("...")`
* runtime panics

Only panics carrying a `Failer` are recovered.

This separation is intentional and critical for operational safety.

---

# Testing

Run all tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. ./...
```

---

# Comparison

`relax` is not trying to replace Go's error model.

Instead, it introduces a constrained and explicit propagation mechanism for internal call chains.

Compared to traditional panic/recover usage:

* failures are typed
* metadata is structured
* stack traces are preserved
* recovery boundaries are explicit
* non-library panics are not swallowed

Compared to repetitive explicit forwarding:

* happy-path code becomes significantly cleaner
* deeply nested orchestration code becomes easier to read
* contextual metadata becomes easier to attach consistently

---

# Example Boundary Pattern

```go
func HandleRequest() error {
    return relax.GuardErr(func() error {
        user := relax.FailCheck(loadUser())

        if err := validate(user); err != nil{
            relax.FailWith(err)
        }
        // or just use like this:
        relax.FailWith(process(user))

        return nil
    })
}
```

The internal flow remains linear while the external API still returns standard Go errors.

---

# License

MIT — see `LICENSE`.
