package relax

import (
	"errors"
	"testing"
)

func TestGuardValue_CatchAndRethrowPreservesFailCheckError(t *testing.T) {
	_, err := GuardValue(func() string {
		_, innerErr := GuardValue(func() string {
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

func TestGuardValue_Success(t *testing.T) {
	result, err := GuardValue(func() int {
		return 42
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGuardValue_Failer(t *testing.T) {
	result, err := GuardValue(func() int {
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

func TestGuardValue_FailerPointer(t *testing.T) {
	result, err := GuardValue(func() int {
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

func TestGuardValue_OtherPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Expected re-panic, but none occurred")
		}
		if r != "other panic" {
			t.Errorf("Expected 'other panic', got %v", r)
		}
	}()

	_, _ = GuardValue(func() int {
		panic("other panic")
	})
}

func TestGuardResult_Success(t *testing.T) {
	result, err := GuardResult(func() (int, error) {
		return 42, nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 42 {
		t.Errorf("Expected 42, got %d", result)
	}
}

func TestGuardResult_Failer(t *testing.T) {
	result, err := GuardResult(func() (int, error) {
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

func TestGuardErr_Success(t *testing.T) {
	err := GuardErr(func() error {
		return nil
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGuardErr_Failer(t *testing.T) {
	err := GuardErr(func() error {
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

func TestGuard_Success(t *testing.T) {
	err := Guard(func() {
		// no thrown error
	})
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGuard_Failer(t *testing.T) {
	err := Guard(func() {
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
