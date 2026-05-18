/*
Package relax provides small helpers for structured, typed panic-based
propagation inside well-defined internal boundaries.

Use `FailWith` and `FailOnError*` to escalate errors with optional key/value
context without changing function signatures. Use `Check*` helpers at
boundary points to recover `Failer` panics back into normal `error` values.

See the package examples in `example_test.go` for runnable usage samples.
*/
package relax
