package impl

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
