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

// CheckFailer executes fn and returns any error produced during execution.
//
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// CheckFailer is intended for functions that do not return a value but may fail,
// and where failures should be handled as errors instead of panics.
func CheckFailer(fn func()) (err error) {
	defer recoverInto(&err)

	fn()
	return
}

// CheckValue executes fn and converts its execution into a checked call.
//
// If fn completes successfully, its return value is returned and err is nil.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// CheckValue is intended for boundary layers (such as HTTP handlers or goroutine
// entry points) where panics of type Failer should be translated into errors.
func CheckValue[T any](fn func() T) (result T, err error) {
	defer recoverInto(&err)

	result = fn()
	return
}

// CheckValue2 executes fn and converts its execution into a checked call for two returned values.
func CheckValue2[T1, T2 any](fn func() (T1, T2)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2 = fn()
	return
}

// CheckValue3 executes fn and converts its execution into a checked call for three returned values.
func CheckValue3[T1, T2, T3 any](fn func() (T1, T2, T3)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3 = fn()
	return
}

// CheckError executes fn which returns only an error.
//
// If fn returns a non-nil error, it is returned unchanged.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// CheckError is used for command-style functions where no value is returned,
// but execution may still fail.
func CheckError(fn func() error) (err error) {
	defer recoverInto(&err)

	err = fn()
	return
}

// CheckResult executes fn which returns a value and an error.
//
// If fn returns a non-nil error, it is returned unchanged.
// If fn panics with a Failer, the panic is recovered and converted into an error.
// Any other panic is re-panicked unchanged.
//
// CheckResult is used when the underlying function already follows Go's
// (T, error) convention but still needs panic-to-error boundary protection.
func CheckResult[T any](fn func() (T, error)) (result T, err error) {
	defer recoverInto(&err)

	result, err = fn()
	return
}

// CheckResult2 executes fn which returns two values and an error.
func CheckResult2[T1, T2 any](fn func() (T1, T2, error)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2, err = fn()
	return
}

// CheckResult3 executes fn which returns three values and an error.
func CheckResult3[T1, T2, T3 any](fn func() (T1, T2, T3, error)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3, err = fn()
	return
}

// HandleFailer executes fn inside a CheckFailer boundary and forwards any recovered
// Failer as a standard error to onError.
//
// HandleFailer is useful for explicit failure-handling boundaries where
// panic-based propagation is used internally, but failures must be handled
// locally instead of escaping the current execution flow.
//
// Only panics carrying a Failer are recovered.
// Any other panic is re-panicked unchanged.
//
// onError must not be nil.
// Passing a nil handler causes HandleFailer to panic immediately.
//
// Example:
//
//	relax.HandleFailer(func() {
//	    process()
//	}, logger.Error)
//
// In this example, if process panics with a Failer, the panic is recovered and
// forwarded to logger.Error as a standard error. Any other panic from process
// is re-panicked unchanged. Could be safely used at the entry point of a goroutine,
// worker loop, or background job:
//
// go relax.HandleFailer(process, logger.Error)
//
// HandleFailer is especially useful at goroutine entry points, worker loops,
// background jobs, and asynchronous execution boundaries.
func HandleFailer(fn func(), onError func(error)) {
	if onError == nil {
		panic("relax: onError cannot be nil")
	}

	if err := CheckFailer(fn); err != nil {
		onError(err)
	}
}
