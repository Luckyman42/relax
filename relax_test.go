package relax

import (
	"errors"
	"sync"
	"testing"
)

func TestFailWith(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", failer.Err.Error())
		}
		if len(failer.Stack) == 0 {
			t.Error("Expected stack trace, but it was empty")
		}
	}()

	FailWith(errors.New("test error"))
}
func TestFailWith_Pointer(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", failer.Err.Error())
		}
		if len(failer.Stack) == 1 {
			t.Error("Expected stack trace, but it was empty")
		}
	}()
	FailWith(&Failer{Err: errors.New("test error")})
}

func TestFailWith_Context(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "context error" {
			t.Errorf("Expected 'context error', got '%s'", failer.Err.Error())
		}
		if failer.Context["user"] != "alice" {
			t.Errorf("Expected context user='alice', got %v", failer.Context["user"])
		}
		if failer.Context["attempt"] != 3 {
			t.Errorf("Expected context attempt=3, got %v", failer.Context["attempt"])
		}
		if failer.Timestamp.IsZero() {
			t.Error("Expected Timestamp to be set")
		}
	}()

	FailWith(errors.New("context error"), "user", "alice", "attempt", 3)
}

func TestFailWith_ExistingFailerRepanicsPlainly(t *testing.T) {
	orig := Failer{Err: errors.New("existing failer"), Context: map[string]any{"a": 1}}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "existing failer" {
			t.Errorf("Expected 'existing failer', got '%s'", failer.Err.Error())
		}
		if failer.Context["a"] != 1 {
			t.Errorf("Expected context a=1, got %v", failer.Context["a"])
		}
	}()

	FailWith(orig)
}

func TestFailWith_ExistingFailerMergesContext(t *testing.T) {
	orig := Failer{Err: errors.New("existing failer"), Context: map[string]any{"a": 1}}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "existing failer" {
			t.Errorf("Expected 'existing failer', got '%s'", failer.Err.Error())
		}
		if failer.Context["a"] != 3 {
			t.Errorf("Expected context a=3, got %v", failer.Context["a"])
		}
		if failer.Context["b"] != "two" {
			t.Errorf("Expected context b='two', got %v", failer.Context["b"])
		}
	}()

	FailWith(orig, "a", 3, "b", "two")
}

func TestFailer_FailMethodMergesContext(t *testing.T) {
	orig := Failer{Err: errors.New("existing error"), Context: map[string]any{"a": 1}}
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Context["a"] != 3 {
			t.Errorf("Expected context a=3, got %v", failer.Context["a"])
		}
		if failer.Context["b"] != "two" {
			t.Errorf("Expected context b='two', got %v", failer.Context["b"])
		}
	}()

	orig.Fail("a", 3, "b", "two")
}

func TestGuard_CatchAndRethrowPreservesFailCheckError(t *testing.T) {
	_, err := Guard(func() string {
		_, innerErr := Guard(func() string {
			return FailCheck("ok", errors.New("must fail"))
		})
		if innerErr != nil {
			FailWith(innerErr)
		}
		return ""
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if errors.As(failer.Err, &Failer{}) {
		t.Error("Expected underlying Err not to be a nested Failer")
	}
}

func TestGuard_Success(t *testing.T) {
	result, err := Guard(func() int {
		return 42
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGuard_Failer(t *testing.T) {
	result, err := Guard(func() int {
		FailWith(errors.New("thrown error"))
		return 0
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGuard_FailerPointer(t *testing.T) {
	result, err := Guard(func() int {
		FailWith(&Failer{Err: errors.New("thrown error")})
		return 0
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGuard_OtherPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected re-panic, but none occurred")
		}
		if r != "other panic" {
			t.Errorf("Expected 'other panic', got %v", r)
		}
	}()

	_, _ = Guard(func() int {
		panic("other panic")
	})
}

func TestIsFailer(t *testing.T) {
	existing := Failer{Err: errors.New("test")}
	if !IsFailer(existing) {
		t.Error("Expected IsFailer to return true for Failer value")
	}
	if !IsFailer(&existing) {
		t.Error("Expected IsFailer to return true for Failer pointer")
	}
	if IsFailer(errors.New("plain error")) {
		t.Error("Expected IsFailer to return false for plain error")
	}
}

func TestGuardErr_Success(t *testing.T) {
	result, err := GuardErr(func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGuardErr_Failer(t *testing.T) {
	result, err := GuardErr(func() (int, error) {
		FailWith(errors.New("thrown error"))
		return 0, nil
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestGuardErr0_Success(t *testing.T) {
	err := GuardErr0(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGuardErr0_Failer(t *testing.T) {
	err := GuardErr0(func() error {
		FailWith(errors.New("thrown error"))
		return nil
	})
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", failer.Err.Error())
	}
}

func TestGuardErr2_Success(t *testing.T) {
	result1, result2, err := GuardErr2(func() (int, int, error) {
		return 42, 23, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result1 != 42 {
		t.Errorf("Expected 42, got %d", result1)
	}
	if result2 != 23 {
		t.Errorf("Expected 23, got %d", result2)
	}
}

func TestGuardErr2_Failer(t *testing.T) {
	_, _, err := GuardErr2(func() (int, string, error) {
		FailWith(errors.New("unwind2 error"))
		return 0, "", nil
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind2 error" {
		t.Errorf("Expected 'unwind2 error', got '%s'", failer.Err.Error())
	}
}

func TestGuardErr3_Success(t *testing.T) {
	result1, result2, result3, err := GuardErr3(func() (int, int, int, error) {
		return 42, 23, 18, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result1 != 42 {
		t.Errorf("Expected 42, got %d", result1)
	}
	if result2 != 23 {
		t.Errorf("Expected 23, got %d", result2)
	}
	if result3 != 18 {
		t.Errorf("Expected 18, got %d", result3)
	}
}

func TestGuardErr3_Failer(t *testing.T) {
	_, _, _, err := GuardErr3(func() (int, int, int, error) {
		FailWith(errors.New("unwind3 error"))
		return 0, 0, 0, nil
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind3 error" {
		t.Errorf("Expected 'unwind3 error', got '%s'", failer.Err.Error())
	}
}

func TestFailCheck_Success(t *testing.T) {
	result := FailCheck(42, nil)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestFailCheck0_Success(t *testing.T) {
	FailCheck0(nil)
}

func TestFailCheck0_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "must0 error" {
			t.Errorf("Expected 'must0 error', got '%s'", failer.Err.Error())
		}
	}()

	FailCheck0(errors.New("must0 error"))
}

func TestFailCheck2_Success(t *testing.T) {
	a, b := FailCheck2(1, "ok", nil)
	if a != 1 || b != "ok" {
		t.Errorf("Expected (1, ok), got (%d, %s)", a, b)
	}
}

func TestFailCheck2_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "must2 error" {
			t.Errorf("Expected 'must2 error', got '%s'", failer.Err.Error())
		}
	}()

	FailCheck2(0, "", errors.New("must2 error"))
}

func TestFailCheck3_Success(t *testing.T) {
	a, b, c := FailCheck3(1, "ok", true, nil)
	if a != 1 || b != "ok" || !c {
		t.Errorf("Expected (1, ok, true), got (%d, %s, %v)", a, b, c)
	}
}

func TestFailCheck3_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "must3 error" {
			t.Errorf("Expected 'must3 error', got '%s'", failer.Err.Error())
		}
	}()

	FailCheck3(0, "", false, errors.New("must3 error"))
}

func TestGuard0_Success(t *testing.T) {
	err := Guard0(func() {
		// no thrown error
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGuard0_Failer(t *testing.T) {
	err := Guard0(func() {
		FailWith(errors.New("unwind0 error"))
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind0 error" {
		t.Errorf("Expected 'unwind0 error', got '%s'", failer.Err.Error())
	}
}

func TestGuard2_Failer(t *testing.T) {
	_, _, err := Guard2(func() (int, string) {
		FailWith(errors.New("unwind2 error"))
		return 0, ""
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind2 error" {
		t.Errorf("Expected 'unwind2 error', got '%s'", failer.Err.Error())
	}
}

func TestGuard3_Failer(t *testing.T) {
	_, _, _, err := Guard3(func() (int, string, bool) {
		FailWith(errors.New("unwind3 error"))
		return 0, "", false
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var failer Failer
	if !errors.As(err, &failer) {
		t.Fatalf("Expected Failer, got %T", err)
	}
	if failer.Err.Error() != "unwind3 error" {
		t.Errorf("Expected 'unwind3 error', got '%s'", failer.Err.Error())
	}
}

func TestFailCheck_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		failer, ok := r.(Failer)
		if !ok {
			t.Fatalf("Expected Failer, got %T", r)
		}
		if failer.Err.Error() != "must error" {
			t.Errorf("Expected 'must error', got '%s'", failer.Err.Error())
		}
	}()

	FailCheck(0, errors.New("must error"))
}

func TestFailer_Error(t *testing.T) {
	err := errors.New("underlying")
	failer := Failer{Err: err}
	if failer.Error() != "underlying" {
		t.Errorf("Expected 'underlying', got '%s'", failer.Error())
	}
}

func TestFailer_Empty(t *testing.T) {
	failer := Failer{Err: nil}
	if failer.Error() != "" {
		t.Errorf("Expected empty error, got '%s'", failer.Error())
	}
}

func TestConvertToFailer_ReturnsExistingFailer(t *testing.T) {
	existing := Failer{Err: errors.New("existing")}
	parsed := ConvertToFailer(existing)
	if parsed.Err == nil || parsed.Err.Error() != "existing" {
		t.Fatalf("Expected existing failer error, got %v", parsed.Err)
	}
}

func TestConvertToFailer_ReturnsExistingFailerPointer(t *testing.T) {
	existing := Failer{Err: errors.New("existing")}
	parsed := ConvertToFailer(&existing)
	if parsed.Err == nil || parsed.Err.Error() != "existing" {
		t.Fatalf("Expected existing failer error, got %v", parsed.Err)
	}
}

func TestConvertToFailer_ReturnsEmptyFailer(t *testing.T) {
	parsed := ConvertToFailer(nil)
	if parsed.Err != nil {
		t.Fatalf("Expected nil Err, got %v", parsed.Err)
	}
}

func TestConvertToFailer_WrapsNonFailerError(t *testing.T) {
	err := errors.New("plain")
	parsed := ConvertToFailer(err)
	if parsed.Err == nil || parsed.Err.Error() != "plain" {
		t.Fatalf("Expected wrapped error, got %v", parsed.Err)
	}
	if parsed.Context != nil {
		t.Errorf("Expected nil context for wrapped error, got %v", parsed.Context)
	}
}

func TestFailer_Unwrap(t *testing.T) {
	err := errors.New("underlying")
	failer := Failer{Err: err}
	if !errors.Is(failer, err) {
		t.Error("Expected errors.Is to work with Unwrap")
	}
}

func TestConvertToFailer_Nil(t *testing.T) {
	var err error = nil
	f := ConvertToFailer(err)
	if f.Err != nil {
		t.Fatalf("Expected nil Err, got %v", f.Err)
	}
}

func TestRecoverInto_RepanicsNonFailer(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected re-panic, but none occurred")
		}
		if r != "panic-now" {
			t.Fatalf("Expected 'panic-now', got %v", r)
		}
	}()

	// Guard should re-panic non-Failer panics; exercise that behavior here.
	_, _ = Guard(func() int {
		panic("panic-now")
	})
}

func TestFailerConcurrency(t *testing.T) {
	const n = 8
	errs := make(chan error, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, err := Guard(func() string {
				FailWith(errors.New("concurrent"))
				return ""
			})
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)

	count := 0
	for e := range errs {
		if e == nil {
			t.Errorf("expected error but got nil")
		} else {
			count++
		}
	}
	if count != n {
		t.Fatalf("expected %d errors, got %d", n, count)
	}
}

func BenchmarkConvertToFailer(b *testing.B) {
	err := errors.New("bench")
	for i := 0; i < b.N; i++ {
		_ = ConvertToFailer(err)
	}
}
