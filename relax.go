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

	var failer Failer
	if errors.As(err, &failer) {
		return failer
	}

	return newFailer(err)
}

// Guard executes the given function and recovers only Failer panics,
// converting them back to errors. Other panics are re-panicked.
// This provides a safe recovery boundary.
func Guard[T any](fn func() T) (result T, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		failer, ok := r.(Failer)
		if !ok {
			panic(r)
		}

		err = failer
	}()

	result = fn()
	return
}

// Guard0 recovers from Failer panics in a function that returns no values.
func Guard0(fn func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		failer, ok := r.(Failer)
		if !ok {
			panic(r)
		}

		err = failer
	}()

	fn()
	return
}

// Guard2 recovers from Failer panics in a function that returns two values.
func Guard2[T1 any, T2 any](fn func() (T1, T2)) (result1 T1, result2 T2, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		failer, ok := r.(Failer)
		if !ok {
			panic(r)
		}

		err = failer
	}()

	result1, result2 = fn()
	return
}

// Guard3 recovers from Failer panics in a function that returns three values.
func Guard3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (result1 T1, result2 T2, result3 T3, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		failer, ok := r.(Failer)
		if !ok {
			panic(r)
		}

		err = failer
	}()

	result1, result2, result3 = fn()
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
