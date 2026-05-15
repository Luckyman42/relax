package relax

import (
	"errors"
	"fmt"
	"runtime/debug"
	"time"
)

// Failer represents a structured error that can be thrown and caught.
// It implements the error interface and preserves stack traces and metadata.
type Failer struct {
	Err       error
	Stack     []byte
	Timestamp time.Time
	Context   map[string]any
}

// Error returns the underlying error message for this Failer.
func (f Failer) Error() string {
	if f.Err != nil {
		return f.Err.Error()
	}
	return ""
}

// Unwrap returns the underlying error for compatibility with errors.As and errors.Is.
func (f Failer) Unwrap() error {
	return f.Err
}

// Fail panics with this Failer.
// If extra key/value pairs are provided, they are merged into the Failer context.
// This avoids wrapping a Failer inside another Failer when rethrowing.
func (f Failer) Fail(keyVals ...any) {
	if len(keyVals) > 0 {
		if f.Context == nil {
			f.Context = make(map[string]any)
		}
		for i := 0; i < len(keyVals); i += 2 {
			key := fmt.Sprint(keyVals[i])
			var value any
			if i+1 < len(keyVals) {
				value = keyVals[i+1]
			}
			f.Context[key] = value
		}
	}
	panic(f)
}

func newFailer(err error, keyVals ...any) Failer {
	failer := Failer{
		Err:       err,
		Stack:     debug.Stack(),
		Timestamp: time.Now(),
	}

	if len(keyVals) == 0 {
		return failer
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
	failer.Context = context
	return failer
}

func recoverFailer(r any) (Failer, bool) {
	switch f := r.(type) {
	case Failer:
		return f, true
	case *Failer:
		if f == nil {
			return Failer{}, false
		}
		return *f, true
	default:
		return Failer{}, false
	}
}

func recoverInto(err *error) {
	r := recover()
	if r == nil {
		return
	}
	failer, ok := recoverFailer(r)
	if !ok {
		panic(r)
	}
	*err = failer
}

// FailWith panics with a Failer wrapping the given error.
// If err is already a Failer, it is re-panicked directly.
// If extra key/value pairs are provided, they are merged into the Failer context.
// If err is nil, it does nothing.
func FailWith(err error, keyVals ...any) {
	if err == nil {
		return
	}

	switch f := err.(type) {
	case Failer:
		f.Fail(keyVals...)
	case *Failer:
		if f == nil {
			return
		}
		(*f).Fail(keyVals...)
	default:
		panic(newFailer(err, keyVals...))
	}
}

// ConvertToFailer converts any error into a Failer.
// If err is already a Failer, it is returned unchanged.
// Otherwise the error is wrapped into a new Failer.
func ConvertToFailer(err error) Failer {
	if err == nil {
		return Failer{}
	}

	switch f := err.(type) {
	case Failer:
		return f
	case *Failer:
		if f == nil {
			return Failer{}
		}
		return *f
	}

	var pointerFailer *Failer
	if errors.As(err, &pointerFailer) {
		if pointerFailer == nil {
			return Failer{}
		}
		return *pointerFailer
	}

	var failer Failer
	if errors.As(err, &failer) {
		return failer
	}

	return newFailer(err)
}

// IsFailer reports whether err is a Failer or wraps a Failer.
func IsFailer(err error) bool {
	if err == nil {
		return false
	}

	var pointerFailer *Failer
	if errors.As(err, &pointerFailer) {
		return pointerFailer != nil
	}

	var failer Failer
	return errors.As(err, &failer)
}

// Guard executes the given function and recovers only Failer panics,
// converting them back to errors. Other panics are re-panicked.
// This provides a safe recovery boundary.
func Guard[T any](fn func() T) (result T, err error) {
	defer recoverInto(&err)

	result = fn()
	return
}

// Guard0 recovers from Failer panics in a function that returns no values.
func Guard0(fn func()) (err error) {
	defer recoverInto(&err)

	fn()
	return
}

// Guard2 recovers from Failer panics in a function that returns two values.
func Guard2[T1 any, T2 any](fn func() (T1, T2)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2 = fn()
	return
}

// Guard3 recovers from Failer panics in a function that returns three values.
func Guard3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3 = fn()
	return
}

// GuardErr recovers from Failer panics in a function that returns (T, error).
// It preserves the normal returned error or returns a Failer if one was thrown.
func GuardErr[T any](fn func() (T, error)) (result T, err error) {
	defer recoverInto(&err)

	result, err = fn()
	return
}

// GuardErr0 recovers from Failer panics in a function that returns an error.
func GuardErr0(fn func() error) (err error) {
	defer recoverInto(&err)

	err = fn()
	return
}

// GuardErr2 recovers from Failer panics in a function that returns (T1, T2, error).
func GuardErr2[T1 any, T2 any](fn func() (T1, T2, error)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2, err = fn()
	return
}

// GuardErr3 recovers from Failer panics in a function that returns (T1, T2, T3, error).
func GuardErr3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3, error)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3, err = fn()
	return
}

// FailCheck checks if err is not nil and fails with it if so.
// Otherwise, returns v. This reduces boilerplate in error propagation.
func FailCheck[T any](v T, err error, keyVals ...any) T {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v
}

// FailCheck0 fails if err is not nil and otherwise returns nothing.
func FailCheck0(err error, keyVals ...any) {
	if err != nil {
		FailWith(err, keyVals...)
	}
}

// FailCheck2 fails if err is not nil and otherwise returns two values.
func FailCheck2[T1 any, T2 any](v1 T1, v2 T2, err error, keyVals ...any) (T1, T2) {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v1, v2
}

// FailCheck3 fails if err is not nil and otherwise returns three values.
func FailCheck3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error, keyVals ...any) (T1, T2, T3) {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v1, v2, v3
}
