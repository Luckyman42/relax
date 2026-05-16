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

func TestFailCheck_Success(t *testing.T) {
	result := FailCheck(42, nil)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
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

func TestConvertToFailer_ReturnsNilFailerPointer(t *testing.T) {
	var p *Failer = nil
	var err error = p
	parsed := ConvertToFailer(err)
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

	// GuardValue should re-panic non-Failer panics; exercise that behavior here.
	_, _ = GuardValue(func() int {
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
			_, err := GuardValue(func() string {
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
