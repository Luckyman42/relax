package relax

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

// Guard executes fn and returns any error produced during execution.
//
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// Guard is intended for functions that do not return a value but may fail,
// and where failures should be handled as errors instead of panics.
func Guard(fn func()) (err error) {
	defer recoverInto(&err)

	fn()
	return
}

// GuardValue executes fn and converts its execution into a guarded call.
//
// If fn completes successfully, its return value is returned and err is nil.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// GuardValue is intended for boundary layers (such as HTTP handlers or goroutine
// entry points) where panics of type Failer should be translated into errors.
func GuardValue[T any](fn func() T) (result T, err error) {
	defer recoverInto(&err)

	result = fn()
	return
}

// GuardErr executes fn which returns only an error.
//
// If fn returns a non-nil error, it is returned unchanged.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// GuardErr is used for command-style functions where no value is returned,
// but execution may still fail.
func GuardErr(fn func() error) (err error) {
	defer recoverInto(&err)

	err = fn()
	return
}

// GuardResult executes fn which returns a value and an error.
//
// If fn returns a non-nil error, it is returned unchanged.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// GuardResult is used when the underlying function already follows Go's
// (T, error) convention but still needs panic-to-error boundary protection.
func GuardResult[T any](fn func() (T, error)) (result T, err error) {
	defer recoverInto(&err)

	result, err = fn()
	return
}
