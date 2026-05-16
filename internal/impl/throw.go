package impl

import "errors"

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
