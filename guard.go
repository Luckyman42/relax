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

// GuardHandle executes fn inside a Guard boundary and forwards any recovered
// Failer as a standard error to onError.
//
// GuardHandle is useful for explicit failure-handling boundaries where
// panic-based propagation is used internally, but failures must be handled
// locally instead of escaping the current execution flow.
//
// Only panics carrying a Failer are recovered.
// Any other panic is re-panicked unchanged.
//
// onError must not be nil.
// Passing a nil handler causes GuardHandle to panic immediately.
//
// Example:
//
//	relax.GuardHandle(func() {
//	    process()
//	}, logger.Error)
//
// In this example, if process panics with a Failer, the panic is recovered and
// forwarded to logger.Error as a standard error. Any other panic from process
// is re-panicked unchanged. Could be safely used at the entry point of a goroutine,
// worker loop, or background job:
//
// go relax.GuardHandle(process, logger.Error)
//
// GuardHandle is especially useful at goroutine entry points, worker loops,
// background jobs, and asynchronous execution boundaries.
func GuardHandle(fn func(), onError func(error)) {
	if onError == nil {
		panic("relax: onError cannot be nil")
	}

	if err := Guard(fn); err != nil {
		onError(err)
	}
}

// GuardGo starts fn in a new goroutine protected by GuardHandle.
//
// If fn panics with a Failer, the panic is recovered and forwarded to
// onError as a standard error.
//
// Any non-Failer panic is re-panicked unchanged inside the goroutine.
//
// onError must not be nil.
// Passing a nil handler causes GuardGo to panic immediately before launching
// the goroutine.
//
// GuardGo provides a safe execution boundary for asynchronous tasks using
// panic-based propagation internally.
//
// Example:
//
//	relax.GuardGo(func() {
//	    syncUsers()
//	}, logger.Error)
//
// GuardGo is intended for fire-and-forget tasks, worker execution,
// background synchronization, and concurrent internal pipelines.
func GuardGo(fn func(), onError func(error)) {
	if onError == nil {
		panic("relax: onError cannot be nil")
	}

	go GuardHandle(fn, onError)
}
