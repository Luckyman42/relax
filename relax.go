package relax

import (
	impl "github.com/luckyman42/relax/internal/impl"
)

// Failer is the public, exported representation of a thrown failure.
//
// A `Failer` preserves the original `error` (Err), a captured stack trace
// (Stack), the time it was created (Timestamp), and an optional
// map[string]any Context for arbitrary key/value metadata. The library uses
// `Failer` values to implement structured panic-based propagation inside
// trusted internal call chains: callers may `panic` a `Failer` (via
// `FailWith` or `Failer.Fail`) and a `Guard` boundary will convert that panic
// back into a returned `error`.
//
// This symbol is a type alias to the internal implementation to keep the
// public surface stable while the implementation lives under
// `internal/impl`.
type Failer = impl.Failer

// FailWith panics with a `Failer` that wraps `err`.
//
// If `err` is already a `Failer` (value or pointer) it will be re-panicked
// directly; in that case any provided key/value pairs are merged into the
// existing `Failer.Context`. If `err` is nil, `FailWith` is a no-op.
//
// The `keyVals` are interpreted as alternating key, value pairs. Keys are
// stringified using `fmt.Sprint`; an odd number of `keyVals` is allowed and
// the final key will be assigned a `nil` value.
func FailWith(err error, keyVals ...any) { impl.FailWith(err, keyVals...) }

// ConvertToFailer converts any error into a `Failer` value.
//
// If `err` is already a `Failer` (or wraps one), the underlying `Failer` is
// returned unchanged. Otherwise a new `Failer` is created capturing the
// current stack trace and timestamp.
func ConvertToFailer(err error) Failer { return impl.ConvertToFailer(err) }

// IsFailer reports whether `err` is a `Failer` value or wraps a `Failer`.
func IsFailer(err error) bool { return impl.IsFailer(err) }

// Guard executes `fn` and recovers only `Failer` panics.
//
// When `fn` returns normally Guard returns the result and a nil error.
// If `fn` panics with a `Failer` the panic is recovered and returned as an
// `error`. Any other panic is re-panicked unchanged. Use `Guard` at clear
// boundary points (for example request handlers or goroutine entry points)
// where converting thrown `Failer`s back into errors is appropriate.
func Guard[T any](fn func() T) (T, error) { return impl.Guard(fn) }

// Guard0 is the `Guard` variant for functions that return no value.
func Guard0(fn func()) error { return impl.Guard0(fn) }

// Guard2 is the `Guard` variant for functions that return two values.
func Guard2[T1 any, T2 any](fn func() (T1, T2)) (T1, T2, error) { return impl.Guard2(fn) }

// Guard3 is the `Guard` variant for functions that return three values.
func Guard3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (T1, T2, T3, error) {
	return impl.Guard3(fn)
}

// GuardErr executes `fn` which returns (T, error). If `fn` returns a
// non-nil `error`, it is returned unchanged. If `fn` panics with a `Failer`
// the panic is recovered and returned as the `error` value.
//
// This form lets callers preserve the usual (value, error) signature while
// allowing internal call chains to use `FailWith` for propagation.
func GuardErr[T any](fn func() (T, error)) (T, error) { return impl.GuardErr(fn) }

// GuardErr0 is the `GuardErr` variant for functions that return only an error.
func GuardErr0(fn func() error) error { return impl.GuardErr0(fn) }

// GuardErr2 is the `GuardErr` variant for functions returning (T1, T2, error).
func GuardErr2[T1 any, T2 any](fn func() (T1, T2, error)) (T1, T2, error) { return impl.GuardErr2(fn) }

// GuardErr3 is the `GuardErr` variant for functions returning
// (T1, T2, T3, error).
func GuardErr3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3, error)) (T1, T2, T3, error) {
	return impl.GuardErr3(fn)
}

// FailCheck returns `v` if `err == nil`; otherwise it throws the error via
// `FailWith(err, keyVals...)`.
//
// This reduces error-forwarding boilerplate inside internal call chains where
// panic-based propagation is acceptable. Prefer explicit returns in public
// APIs.
func FailCheck[T any](v T, err error, keyVals ...any) T { return impl.FailCheck(v, err, keyVals...) }

// FailCheck0 throws if `err` is not nil for functions that only return an error.
func FailCheck0(err error, keyVals ...any) { impl.FailCheck0(err, keyVals...) }

// FailCheck2 throws if `err` is not nil and otherwise returns the two values.
func FailCheck2[T1 any, T2 any](v1 T1, v2 T2, err error, keyVals ...any) (T1, T2) {
	return impl.FailCheck2(v1, v2, err, keyVals...)
}

// FailCheck3 throws if `err` is not nil and otherwise returns three values.
func FailCheck3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error, keyVals ...any) (T1, T2, T3) {
	return impl.FailCheck3(v1, v2, v3, err, keyVals...)
}
