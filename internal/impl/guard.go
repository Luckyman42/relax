package impl

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
