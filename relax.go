package relax

import (
	"fmt"
	"runtime/debug"
	"time"
)

// Throwable represents a structured error that can be thrown and caught.
// It implements the error interface and preserves stack traces and metadata.
type Throwable struct {
	Err       error
	Stack     []byte
	Timestamp time.Time
	Context   map[string]any
}

// Error implements the error interface.
func (t Throwable) Error() string {
	if t.Err != nil {
		return t.Err.Error()
	}
	return ""
}

// Unwrap returns the underlying error for compatibility with errors.As and errors.Is.
func (t Throwable) Unwrap() error {
	return t.Err
}

func newThrowable(err error, keyVals ...any) Throwable {
	throwable := Throwable{
		Err:       err,
		Stack:     debug.Stack(),
		Timestamp: time.Now(),
	}

	if len(keyVals) == 0 {
		return throwable
	}

	context := make(map[string]any, (len(keyVals)+1)/2)
	for i := 0; i < len(keyVals); i += 2 {
		key := fmt.Sprint(keyVals[i])
		var value any
		if i+1 < len(keyVals) {
			value = keyVals[i+1]
		}
		context[key] = value
	}
	throwable.Context = context
	return throwable
}

// Throw panics with a Throwable wrapping the given error.
// Any extra key/value pairs are stored in Throwable.Context.
// If err is nil, it does nothing.
func Throw(err error, keyVals ...any) {
	if err == nil {
		return
	}
	panic(newThrowable(err, keyVals...))
}

// Handle executes the given function and recovers only Throwable panics,
// converting them back to errors. Other panics are re-panicked.
// This provides a safe recovery boundary.
func Handle[T any](fn func() T) (result T, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		throwable, ok := r.(Throwable)
		if !ok {
			panic(r)
		}

		err = throwable
	}()

	result = fn()
	return
}

// Handle0 recovers from Throwable panics in a function that returns no values.
func Handle0(fn func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		throwable, ok := r.(Throwable)
		if !ok {
			panic(r)
		}

		err = throwable
	}()

	fn()
	return
}

// Handle2 recovers from Throwable panics in a function that returns two values.
func Handle2[T1 any, T2 any](fn func() (T1, T2)) (result1 T1, result2 T2, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		throwable, ok := r.(Throwable)
		if !ok {
			panic(r)
		}

		err = throwable
	}()

	result1, result2 = fn()
	return
}

// Handle3 recovers from Throwable panics in a function that returns three values.
func Handle3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (result1 T1, result2 T2, result3 T3, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		throwable, ok := r.(Throwable)
		if !ok {
			panic(r)
		}

		err = throwable
	}()

	result1, result2, result3 = fn()
	return
}

// Must checks if err is not nil and throws it if so.
// Otherwise, returns v. This reduces boilerplate in error propagation.
func Must[T any](v T, err error) T {
	if err != nil {
		Throw(err, "must", true)
	}
	return v
}

// Must0 throws if err is not nil and otherwise returns nothing.
func Must0(err error) {
	if err != nil {
		Throw(err, "must", true)
	}
}

// Must2 throws if err is not nil and otherwise returns two values.
func Must2[T1 any, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		Throw(err, "must", true)
	}
	return v1, v2
}

// Must3 throws if err is not nil and otherwise returns three values.
func Must3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		Throw(err, "must", true)
	}
	return v1, v2, v3
}
