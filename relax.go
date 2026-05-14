package relax

import (
	"errors"
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
// If err is already a Throwable, it is re-panicked directly.
// If extra key/value pairs are provided, they are merged into the Throwable context.
// If err is nil, it does nothing.
func Throw(err error, keyVals ...any) {
	if err == nil {
		return
	}

	switch t := err.(type) {
	case Throwable:
		if len(keyVals) == 0 {
			panic(t)
		}
		if t.Context == nil {
			t.Context = make(map[string]any)
		}
		for i := 0; i < len(keyVals); i += 2 {
			key := fmt.Sprint(keyVals[i])
			var value any
			if i+1 < len(keyVals) {
				value = keyVals[i+1]
			}
			t.Context[key] = value
		}
		panic(t)
	case *Throwable:
		if t == nil {
			return
		}
		copy := *t
		if len(keyVals) == 0 {
			panic(copy)
		}
		if copy.Context == nil {
			copy.Context = make(map[string]any)
		}
		for i := 0; i < len(keyVals); i += 2 {
			key := fmt.Sprint(keyVals[i])
			var value any
			if i+1 < len(keyVals) {
				value = keyVals[i+1]
			}
			copy.Context[key] = value
		}
		panic(copy)
	default:
		panic(newThrowable(err, keyVals...))
	}
}

// ParseError converts any error into a Throwable.
// If err is already a Throwable, it is returned unchanged.
// Otherwise the error is wrapped into a new Throwable.
func ParseError(err error) Throwable {
	if err == nil {
		return Throwable{}
	}

	var throwable Throwable
	if errors.As(err, &throwable) {
		return throwable
	}

	return newThrowable(err)
}

// IsMust reports whether err is a Throwable propagated through Must.
func IsMust(err error) bool {
	if err == nil {
		return false
	}

	var throwable Throwable
	if !errors.As(err, &throwable) {
		return false
	}

	return throwable.Context != nil && throwable.Context["must"] == true
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
