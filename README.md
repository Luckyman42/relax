# Relax

> Don't **Panic**, just **Relax**!

Relax is a small Go toolkit for structured, typed panic-based propagation inside trusted internal paths.
It is designed to reduce boilerplate where many layers only forward errors without handling them.

## What Relax is

Relax is not a replacement for Go's error handling.
It is a companion for internal call chains where explicit forwarding becomes noisy, but the error should still be handled at a boundary.

The library provides:
- a typed propagation wrapper: `Failer`
- `FailWith(...)` for intentional propagation
- `Guard(...)` boundaries that recover only `Failer`
- helper forms like `FailCheck`, `FailCheck0`, `FailCheck2`, and `FailCheck3`

That means:
- `FailCheck` escalates failure through structured propagation, not "program invalid" semantics.
- only errors intentionally thrown through Relax are recovered
- runtime panics and programmer errors still behave like normal panics
- metadata can be attached without changing function signatures

## Why this exists

In deep service stacks, middleware, or parser-like flows, many intermediate functions do nothing but bubble failures upward.
In those cases, the usual `if err != nil { return err }` pattern adds a lot of repetition.

```
A -> B -> C -> D -> E
```

If `E()` returns an error, and none of the intermediate layers (`B/C/D`) can meaningfully handle it, then traditional Go error forwarding becomes repetitive boilerplate:

```go
v, err := Next()
if err != nil {
    return ..., err
}
```

Repeated over many layers.

In many cases, these layers:
- do not recover,
- do not retry,
- do not add context,
- do not compensate,
- do not transform the error.

They simply propagate it.

This creates what can reasonably be described as:
- forwarding noise,
- ceremony,
- cognitive overhead.

Relax lets you keep the propagation path cleaner while still preserving explicit recovery at the edges.

One important benefit is visibility: if a thrown error is not recovered, it will surface as a `Failer` panic instead of silently disappearing as an ignored error return.
That makes it easier to catch forgotten handling paths during runtime, rather than letting a dropped `error` value hide a bug.

In other words, an unhandled `Failer` is noisier than an ignored error return, which helps force the developer to treat failures explicitly.

Traditional explicit error handling allows developers to accidentally ignore errors:

```go
v, _ := Something()
```

or:

```go
if err != nil {
    return nil
}
```

Potentially losing:
- logging,
- telemetry,
- context,
- stack information.

With structured panic propagation:
- unhandled failures become visible,
- failures are noisy,
- silent corruption becomes less likely.

In distributed/backend systems, this can be a very reasonable tradeoff.

## When to use Relax

Use Relax for:
- request/handler boundaries
- goroutine entry points
- internal service or pipeline layers
- parser/validation flows where most layers do not handle the error

## When not to use Relax

Avoid Relax for:
- public library APIs
- code where explicit error returns are the expected contract
- control flow that should not be modeled as panic propagation
- hiding bugs or resource management issues

## Installation

```bash
go get github.com/luckyman42/relax
```

## Basic usage

```go
import "relax"

func Service() string {
    data := relax.FailCheck(fetchData())
    return relax.FailCheck(processData(data))
}

func Main() {
    result, err := relax.Guard(Service)
    if err != nil {
        log.Printf("request failed: %s", err)
        return
    }
    fmt.Println(result)
}
```

## FailWithing with metadata

`FailWith` accepts optional key/value pairs that are stored in `Failer.Context`.
This is useful when you want to attach extra information without changing many function signatures.

```go
func validateInput(input string) {
    if input == "" {
        relax.FailWith(errors.New("input required"), "field", "username", "retry", 1)
    }
}
```

## Inspecting failures

`FailCheck` is a helper that escalates errors through structured propagation.
When `FailCheck` throws, it panics with a `Failer` just like `FailWith` would.
At a recovery boundary, you can inspect the error normally:

```go
func Guard() {
    result, err := relax.Guard(func() string {
        return relax.FailCheck(fetchData())
    })
    if err != nil {
        var failer relax.Failer
        if errors.As(err, &failer) {
            log.Printf("failure with metadata: %s", failer.Err)
            log.Printf("context: %v", failer.Context)
        }
        return
    }
    fmt.Println(result)
}
```

## API

- `FailWith(err error, keyVals ...any)`: Panics with a `Failer` wrapping the error and optional context.
- `Guard[T any](fn func() T) (T, error)`: Executes `fn` and recovers only `Failer` panics.
- `Guard0(fn func()) error`: Executes a function with no return values and recovers `Failer` panics.
- `Guard2[T1 any, T2 any](fn func() (T1, T2)) (T1, T2, error)`: Recovers `Failer` panics from a two-value function.
- `Guard3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (T1, T2, T3, error)`: Recovers `Failer` panics from a three-value function.
- `ConvertToFailer(err error) Failer`: Converts any error into a `Failer`, preserving existing `Failer` values.
- `Failer.Fail(keyVals ...any)`: Panics with the `Failer`, optionally merging extra context.
- `FailCheck[T any](v T, err error, keyVals ...any) T`: Fails if `err` is not nil, otherwise returns `v`.
  It escalates failure through structured propagation, not "program invalid" semantics.
- `FailCheck0(err error, keyVals ...any)`: Fails if `err` is not nil for functions that return only error.
- `FailCheck2[T1 any, T2 any](v1 T1, v2 T2, err error, keyVals ...any) (T1, T2)`: Fails if `err` is not nil and returns two values.
- `FailCheck3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error, keyVals ...any) (T1, T2, T3)`: Fails if `err` is not nil and returns three values.
- `Failer`: Implements `error`, preserves stack traces, timestamp, and optional context.

## Safety

- Only `Failer` panics are caught; runtime panics are re-panicked.
- Use `Guard` at well-defined boundaries, not inside every helper.
- Keep cleanup explicit with `defer`.

## Trade-offs

- `panic`/`recover` is more expensive than direct error returns.
- This is a pattern for internal propagation, not a general replacement for errors.
- The main value is cleaner propagation paths, not performance.
- It can slow down the hot path, but it may improve maintainability by reducing repetitive forwarding code.

## Examples

See `example_test.go` for concrete patterns.

## License

See LICENSE file.