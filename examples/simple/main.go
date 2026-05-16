package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/luckyman42/relax"
)

func main() {
	result, err := relax.Guard(func() string {
		data := relax.FailCheck(fetch())
		return relax.FailCheck(process(data))
	})

	if err != nil {
		var f relax.Failer
		if errors.As(err, &f) {
			log.Printf("caught failer: %s, context=%v", f.Err, f.Context)
		} else {
			log.Printf("error: %v", err)
		}
		return
	}

	fmt.Println(result)
}

func fetch() (string, error) {
	return "input", nil
}

func process(s string) (string, error) {
	if s == "" {
		return "", errors.New("empty")
	}
	return "processed:" + s, nil
}
