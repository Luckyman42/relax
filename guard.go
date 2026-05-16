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

// Guard executes `fn` and recovers only `Failer` panics.
//
// When `fn` returns normally Guard returns the result and a nil error.
// If `fn` panics with a `Failer` the panic is recovered and returned as an
// `error`. Any other panic is re-panicked unchanged. Use `Guard` at clear
// boundary points (for example request handlers or goroutine entry points)
// where converting thrown `Failer`s back into errors is appropriate.
func Guard[T any](fn func() T) (result T, err error) {
	defer recoverInto(&err)

	result = fn()
	return
}

// Guard0 is the `Guard` variant for functions that return no value.
func Guard0(fn func()) (err error) {
	defer recoverInto(&err)

	fn()
	return
}

// Guard2 is the `Guard` variant for functions that return two values.
func Guard2[T1 any, T2 any](fn func() (T1, T2)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2 = fn()
	return
}

// Guard3 is the `Guard` variant for functions that return three values.
func Guard3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3 = fn()
	return
}

// GuardErr executes `fn` which returns (T, error). If `fn` returns a
// non-nil `error`, it is returned unchanged. If `fn` panics with a `Failer`
// the panic is recovered and returned as the `error` value.
//
// This form lets callers preserve the usual (value, error) signature while
// allowing internal call chains to use `FailWith` for propagation.
func GuardErr[T any](fn func() (T, error)) (result T, err error) {
	defer recoverInto(&err)

	result, err = fn()
	return
}

// GuardErr0 is the `GuardErr` variant for functions that return only an error.
func GuardErr0(fn func() error) (err error) {
	defer recoverInto(&err)

	err = fn()
	return
}

// GuardErr2 is the `GuardErr` variant for functions returning (T1, T2, error).
func GuardErr2[T1 any, T2 any](fn func() (T1, T2, error)) (result1 T1, result2 T2, err error) {
	defer recoverInto(&err)

	result1, result2, err = fn()
	return
}

// GuardErr3 is the `GuardErr` variant for functions returning
// (T1, T2, T3, error).
func GuardErr3[T1 any, T2 any, T3 any](fn func() (T1, T2, T3, error)) (result1 T1, result2 T2, result3 T3, err error) {
	defer recoverInto(&err)

	result1, result2, result3, err = fn()
	return
}
