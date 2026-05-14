package relax

import (
	"errors"
	"testing"
)

func TestThrow(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", throwable.Err.Error())
		}
		if len(throwable.Stack) == 0 {
			t.Error("Expected stack trace, but it was empty")
		}
	}()

	Throw(errors.New("test error"))
}

func TestThrow_Context(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "context error" {
			t.Errorf("Expected 'context error', got '%s'", throwable.Err.Error())
		}
		if throwable.Context["user"] != "alice" {
			t.Errorf("Expected context user='alice', got %v", throwable.Context["user"])
		}
		if throwable.Context["attempt"] != 3 {
			t.Errorf("Expected context attempt=3, got %v", throwable.Context["attempt"])
		}
		if throwable.Timestamp.IsZero() {
			t.Error("Expected Timestamp to be set")
		}
	}()

	Throw(errors.New("context error"), "user", "alice", "attempt", 3)
}

func TestHandle_Success(t *testing.T) {
	result, err := Handle(func() int {
		return 42
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestHandle_Throwable(t *testing.T) {
	result, err := Handle(func() int {
		Throw(errors.New("thrown error"))
		return 0
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var throwable Throwable
	if !errors.As(err, &throwable) {
		t.Fatalf("Expected Throwable, got %T", err)
	}
	if throwable.Err.Error() != "thrown error" {
		t.Errorf("Expected 'thrown error', got '%s'", throwable.Err.Error())
	}
	if result != 0 {
		t.Errorf("Expected 0, got %d", result)
	}
}

func TestHandle_OtherPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected re-panic, but none occurred")
		}
		if r != "other panic" {
			t.Errorf("Expected 'other panic', got %v", r)
		}
	}()

	_, _ = Handle(func() int {
		panic("other panic")
	})
}

func TestMust_Success(t *testing.T) {
	result := Must(42, nil)
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestMust0_Success(t *testing.T) {
	Must0(nil)
}

func TestMust0_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "must0 error" {
			t.Errorf("Expected 'must0 error', got '%s'", throwable.Err.Error())
		}
		if throwable.Context["must"] != true {
			t.Errorf("Expected must=true, got %v", throwable.Context["must"])
		}
	}()

	Must0(errors.New("must0 error"))
}

func TestMust2_Success(t *testing.T) {
	a, b := Must2(1, "ok", nil)
	if a != 1 || b != "ok" {
		t.Errorf("Expected (1, ok), got (%d, %s)", a, b)
	}
}

func TestMust2_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "must2 error" {
			t.Errorf("Expected 'must2 error', got '%s'", throwable.Err.Error())
		}
		if throwable.Context["must"] != true {
			t.Errorf("Expected must=true, got %v", throwable.Context["must"])
		}
	}()

	Must2(0, "", errors.New("must2 error"))
}

func TestMust3_Success(t *testing.T) {
	a, b, c := Must3(1, "ok", true, nil)
	if a != 1 || b != "ok" || !c {
		t.Errorf("Expected (1, ok, true), got (%d, %s, %v)", a, b, c)
	}
}

func TestMust3_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "must3 error" {
			t.Errorf("Expected 'must3 error', got '%s'", throwable.Err.Error())
		}
		if throwable.Context["must"] != true {
			t.Errorf("Expected must=true, got %v", throwable.Context["must"])
		}
	}()

	Must3(0, "", false, errors.New("must3 error"))
}

func TestHandle0_Success(t *testing.T) {
	err := Handle0(func() {
		// no thrown error
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHandle0_Throwable(t *testing.T) {
	err := Handle0(func() {
		Throw(errors.New("unwind0 error"))
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var throwable Throwable
	if !errors.As(err, &throwable) {
		t.Fatalf("Expected Throwable, got %T", err)
	}
	if throwable.Err.Error() != "unwind0 error" {
		t.Errorf("Expected 'unwind0 error', got '%s'", throwable.Err.Error())
	}
}

func TestHandle2_Throwable(t *testing.T) {
	_, _, err := Handle2(func() (int, string) {
		Throw(errors.New("unwind2 error"))
		return 0, ""
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var throwable Throwable
	if !errors.As(err, &throwable) {
		t.Fatalf("Expected Throwable, got %T", err)
	}
	if throwable.Err.Error() != "unwind2 error" {
		t.Errorf("Expected 'unwind2 error', got '%s'", throwable.Err.Error())
	}
}

func TestHandle3_Throwable(t *testing.T) {
	_, _, _, err := Handle3(func() (int, string, bool) {
		Throw(errors.New("unwind3 error"))
		return 0, "", false
	})
	if err == nil {
		t.Fatal("Expected error, but got none")
	}
	var throwable Throwable
	if !errors.As(err, &throwable) {
		t.Fatalf("Expected Throwable, got %T", err)
	}
	if throwable.Err.Error() != "unwind3 error" {
		t.Errorf("Expected 'unwind3 error', got '%s'", throwable.Err.Error())
	}
}

func TestMust_Error(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected panic, but none occurred")
		}
		throwable, ok := r.(Throwable)
		if !ok {
			t.Fatalf("Expected Throwable, got %T", r)
		}
		if throwable.Err.Error() != "must error" {
			t.Errorf("Expected 'must error', got '%s'", throwable.Err.Error())
		}
		if throwable.Context["must"] != true {
			t.Errorf("Expected must=true, got %v", throwable.Context["must"])
		}
	}()

	Must(0, errors.New("must error"))
}

func TestThrowable_Error(t *testing.T) {
	err := errors.New("underlying")
	throwable := Throwable{Err: err}
	if throwable.Error() != "underlying" {
		t.Errorf("Expected 'underlying', got '%s'", throwable.Error())
	}
}

func TestThrowable_Unwrap(t *testing.T) {
	err := errors.New("underlying")
	throwable := Throwable{Err: err}
	if !errors.Is(throwable, err) {
		t.Error("Expected errors.Is to work with Unwrap")
	}
}