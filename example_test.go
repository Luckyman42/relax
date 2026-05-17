package relax_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/luckyman42/relax"
)

func fetchUser(id int) (string, error) {
	return "", errors.New("database unavailable")
}

func ExampleGuardValue() {
	profile, err := relax.GuardValue(func() string {
		return relax.FailCheck(fetchUser(42))
	})

	if err != nil {
		var failer relax.Failer
		if errors.As(err, &failer) {
			fmt.Println(failer.Err)
		}
		return
	}

	fmt.Println(profile)

	// Output:
	// database unavailable
}

func ExampleGuard() {
	err := relax.Guard(func() {
		relax.FailWith(errors.New("something failed"))
	})

	fmt.Println(err)

	// Output:
	// something failed
}

func ExampleConvertToFailer() {
	err := errors.New("boom")

	failer := relax.ConvertToFailer(err)

	fmt.Println(failer.Err)
	fmt.Println(failer.Timestamp.IsZero())
	fmt.Println(len(failer.Stack) > 0)

	// Output:
	// boom
	// false
	// true
}

func ExampleIsFailer() {
	failer := relax.ConvertToFailer(errors.New("failure"))

	fmt.Println(relax.IsFailer(failer))
	fmt.Println(relax.IsFailer(errors.New("normal error")))

	// Output:
	// true
	// false
}

func ExampleFailer_Fail() {
	err := relax.Guard(func() {
		failer := relax.ConvertToFailer(errors.New("repository failed"))

		failer.Fail(
			"repository", "users",
			"operation", "find",
		)
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
		fmt.Println(failer.Context["repository"])
	}

	// Output:
	// repository failed
	// users
}

func ExampleGuardResult() {
	value, err := relax.GuardResult(func() (int, error) {
		if true {
			relax.FailWith(errors.New("calculation failed"))
		}

		return 42, nil
	})

	fmt.Println(value)
	fmt.Println(err)

	// Output:
	// 0
	// calculation failed
}

func Example_realisticServiceFlow() {
	loadUser := func(id int) (string, error) {
		return "", errors.New("user not found")
	}

	err := relax.GuardErr(func() error {
		user := relax.FailCheck(loadUser(99))

		log.Println(user)
		return nil
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
	}

	// Output:
	// user not found
}

func Example_failerWithContext() {
	loadUser := func(id int) (string, error) {
		return "", errors.New("user not found")
	}

	err := relax.GuardErr(func() error {
		user, err := loadUser(99)

		if err != nil {
			relax.FailWith(err, "userid", 99)
		}

		log.Println(user)
		return nil
	})

	var failer relax.Failer
	if errors.As(err, &failer) {
		fmt.Println(failer.Err)
		fmt.Println(failer.Context["userid"])
	}

	// Output:
	// user not found
	// 99
}

func ExampleGuardHandle() {
	relax.GuardHandle(func() {
		relax.FailWith(errors.New("worker failed"))
	}, func(err error) {
		fmt.Println(err)
	})

	// Output:
	// worker failed
}

func ExampleGuardGo() {
	done := make(chan struct{})

	relax.GuardGo(func() {
		relax.FailWith(errors.New("async failure"))
	}, func(err error) {
		fmt.Println(err)
		close(done)
	})

	<-done

	// Output:
	// async failure
}
