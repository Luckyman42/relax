# Relax

[![CI](https://img.shields.io/github/actions/workflow/status/luckyman42/relax/ci.yml?branch=main&label=ci)](https://github.com/luckyman42/relax/actions)
[![codecov](https://codecov.io/gh/luckyman42/relax/branch/main/graph/badge.svg)](https://codecov.io/gh/luckyman42/relax)
[![Go Report Card](https://goreportcard.com/badge/github.com/luckyman42/relax)](https://goreportcard.com/report/github.com/luckyman42/relax)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/luckyman42/relax)](https://pkg.go.dev/github.com/luckyman42/relax)
![Go Version](https://img.shields.io/github/go-mod/go-version/luckyman42/relax)
[![License](https://img.shields.io/github/license/luckyman42/relax)](LICENSE)
> Don't **panic**, just **relax**!

Typed, structured panic-based propagation for well-defined internal boundaries.

Relax provides small, focused helpers to reduce repetitive error-forwarding in 
internal call chains while preserving the original error, stack trace and
optional structured context.

## Quick summary

- Use `FailWith` / `FailCheck` to escalate errors inside trusted internal
  call chains without changing every function signature.
- Use `Guard*` at explicit boundary points (handlers, goroutine entry points)
  to recover `Failer` panics back into normal `error` values.
- `Failer` preserves the original `error`, a stack trace, timestamp and an
  optional `map[string]any` context.

This is not a replacement for idiomatic explicit error handling in public
APIs — it's a pragmatic tool to reduce boilerplate where forwarding is the
common case.

## Installation

```
go get github.com/luckyman42/relax
```

Import as:

```go
import "github.com/luckyman42/relax"
```

## Example

Minimal, idiomatic example using the helpers:

```go
result, err := relax.Guard(func() string {
    data := relax.FailCheck(fetch())          // throws if fetch() returns an error
    return relax.FailCheck(process(data))     // throws if process() returns an error
})

if err != nil {
    var f relax.Failer
    if errors.As(err, &f) {
        log.Printf("caught failer: %s, context=%v", f.Err, f.Context)
    } else {
        log.Printf("error: %v", err)
    }
    return
}

fmt.Println(result)
```

See runnable examples in [examples/simple/main.go](examples/simple/main.go)
and package examples in [example_test.go](example_test.go).

## API (short)

- `FailWith(err error, keyVals ...any)` — panic with a `Failer` wrapping `err`.
- `ConvertToFailer(err error) Failer` — always get a `Failer` value for an error.
- `IsFailer(err error) bool` — true when an error is or wraps a `Failer`.
- `FailCheck`, `FailCheck0`, `FailCheck2`, `FailCheck3` — helpers that throw
  when an `error` is non-nil and otherwise return the provided values.
- `Guard`, `Guard0`, `Guard2`, `Guard3` — recover `Failer` panics and return
  them as `error`; other panics are re-panicked.
- `GuardErr*` variants preserve an explicit `(T, error)` signature while
  still allowing `FailWith` inside the call chain.

Full signatures and docs are available in the source and GoDoc.

## Design & guarantees

- `Failer` is a public type that implements `error` and preserves the original
  error (see `Failer.Err`), a captured stack trace (`Failer.Stack`), a
  creation timestamp, and optional `Context`.
- The implementation lives under `internal/impl` so the public surface is a
  stable, thin wrapper. See `relax.go` for the exported API and
  `internal/impl` for the implementation.
- `Guard*` functions only recover panics that are `Failer` (or pointers to it).
  Any other panic is re-panicked immediately: this preserves expected
  runtime/panic semantics for programmer errors and nil-pointer dereferences.
- `ConvertToFailer` and `IsFailer` use `errors.As` to avoid double-wrapping and
  to interoperate with wrapped errors.

These choices make the pattern explicit and auditable: thrown failures are
visible, typed and carry debugging metadata.

## When to use

- Internal service layers where forwarding is the dominant behavior.
- Request/handler boundaries and goroutine entry points where you convert
  a thrown `Failer` back to an `error` and handle/log/record it.

## When NOT to use

- Public library APIs or exported functions that callers expect to check
  returned `error` values directly.
- For control-flow, resource cleanup, or when deterministic performance on a
  hot path is required — panic/recover has measurable cost.

## Safety and best practices

- Use `Guard*` only at clear boundaries (handlers, goroutine starts) — do not
  sprinkle `FailWith` across public-facing APIs.
- Keep `FailWith` usage local to packages or modules that agree on the
  convention.
- Avoid returning `Failer` from normal returns; if you must return the
  underlying error, use `fail.Err` or `errors.Unwrap`.
- Always use `defer` for cleanup; `panic` does not replace explicit resource
  management.

## Tests, benchmarks and quality

- Unit tests live in `relax_test.go` and exercise propagation, context
  merging and concurrency behavior.
- A micro-benchmark for `ConvertToFailer` is included.
- Run tests and benchmarks:

```bash
go test ./... -v
go test -bench=. ./...
```

Recommended CI checks:

- `go test` on supported Go versions
- `go vet` and `golangci-lint` (recommended linters)
- publish `pkg.go.dev` docs and a coverage badge

## Contributing

- Fork and submit a PR with a clear description and tests for behavior changes.
- Keep public API changes backward compatible when possible.
- Add examples and tests for any new behavior or edge cases.

See `CHANGELOG.md` for release notes and the `LICENSE` for terms.

## License

MIT — see `LICENSE`.
