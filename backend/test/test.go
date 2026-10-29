package test

import (
	"errors"
	"fmt"
	"log"
	"os"

	// "messenger/messenger/internal/app/api"
	"testing"
)

func Hello(s string) (string, error) {
	if s == "" {
		return "", errors.New("hello: name required")
	}
	// return fmt.Sprintf("Hello, %s", "s"), nil
	return fmt.Sprintf("Hello, %s", s), nil
}

func TestMain(m *testing.M) {
	log.Println("main start")
	code := m.Run()
	log.Println("main end")

	os.Exit(code)
}

func TestFoo(t *testing.T) {
	log.Println("foo statrt")
	tHelp(t)
	if false {
		t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, "", "err", "want")
	}
	log.Println("foo end")
}

func tHelp(t *testing.T) {
	log.Println("helper statrt")
	t.Helper()
	t.Cleanup(func() {
		log.Println("cleanup")
	})
	log.Println("helper end")
}
