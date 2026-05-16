package relax_test

import (
	"errors"
	"fmt"
	"log"

	"github.com/luckyman42/relax"
)

// ExampleFailCheck demonstrates basic FailCheck usage for reducing boilerplate.
func ExampleFailCheck() {
	// Instead of:
	// v, err := someFunc()
	// if err != nil { return err }
	// Use:
	v := relax.FailCheck(someFunc())
	fmt.Println(v)
	// Output: success
}

// ExampleFailWith demonstrates throwing an error with context.
func ExampleFailWith() {
	defer func() {
		if r := recover(); r != nil {
			if throwable, ok := r.(relax.Failer); ok {
				log.Printf("Caught: %s\nStack: %s\nUser: %v", throwable.Err, throwable.Stack, throwable.Context["user"])
			}
		}
	}()
	relax.FailWith(errors.New("example error"), "user", "alice", "route", "/login")
	// This will panic and be caught by the defer
}

// ExampleGuardValue demonstrates safe recovery.
func ExampleGuardValue() {
	result, err := relax.GuardValue(func() string {
		// Simulate a call chain with FailCheck
		data := relax.FailCheck(parseData("input"))
		return relax.FailCheck(processData(data))
	})

	if err != nil {
		log.Printf("GuardValued error: %s", err)
	} else {
		fmt.Println(result)
	}
	// Output: processed:parsed:input
}

// Helper functions for examples
func someFunc() (string, error) {
	return "success", nil
}

func parseData(input string) (string, error) {
	if input == "" {
		return "", errors.New("empty input")
	}
	return "parsed:" + input, nil
}

func processData(data string) (string, error) {
	if data == "" {
		return "", errors.New("no data")
	}
	return "processed:" + data, nil
}

func doSomething() error {
	return nil
}
