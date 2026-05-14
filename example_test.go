package relax_test

import (
	"errors"
	"fmt"
	"log"

	"relax"
)

// ExampleMust demonstrates basic Must usage for reducing boilerplate.
func ExampleMust() {
	// Instead of:
	// v, err := someFunc()
	// if err != nil { return err }
	// Use:
	v := relax.Must(someFunc())
	fmt.Println(v)
	// Output: success
}

// ExampleThrow demonstrates throwing an error with context.
func ExampleThrow() {
	defer func() {
		if r := recover(); r != nil {
			if throwable, ok := r.(relax.Throwable); ok {
				log.Printf("Caught: %s\nStack: %s\nUser: %v", throwable.Err, throwable.Stack, throwable.Context["user"])
			}
		}
	}()
	relax.Throw(errors.New("example error"), "user", "alice", "route", "/login")
	// This will panic and be caught by the defer
}

// ExampleMust0 demonstrates Must0 usage for error-only functions.
func ExampleMust0() {
	relax.Must0(doSomething())
	fmt.Println("done")
	// Output: done
}

// ExampleHandle demonstrates safe recovery.
func ExampleHandle() {
	result, err := relax.Handle(func() string {
		// Simulate a call chain with Must
		data := relax.Must(parseData("input"))
		return relax.Must(processData(data))
	})

	if err != nil {
		log.Printf("Handled error: %s", err)
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