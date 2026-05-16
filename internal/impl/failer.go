package impl

import (
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
