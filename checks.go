package relax

// FailCheck returns `v` if `err == nil`; otherwise it throws the error via
// `FailWith(err, keyVals...)`.
//
// This reduces error-forwarding boilerplate inside internal call chains where
// panic-based propagation is acceptable. Prefer explicit returns in public APIs.
func FailCheck[T any](v T, err error, keyVals ...any) T {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v
}

// FailCheck0 throws if `err` is not nil for functions that only return an error.
func FailCheck0(err error, keyVals ...any) {
	if err != nil {
		FailWith(err, keyVals...)
	}
}

// FailCheck2 throws if `err` is not nil and otherwise returns the two values.
func FailCheck2[T1 any, T2 any](v1 T1, v2 T2, err error, keyVals ...any) (T1, T2) {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v1, v2
}

// FailCheck3 throws if `err` is not nil and otherwise returns three values.
func FailCheck3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error, keyVals ...any) (T1, T2, T3) {
	if err != nil {
		FailWith(err, keyVals...)
	}
	return v1, v2, v3
}
