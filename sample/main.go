package main

import (
	"errors"
	"github.com/sohamkamani/detective"
	"net/http"
	"time"
)

type cache struct {
}

func (c cache) ping() error {
	return nil
}

type db struct {
}

func (d db) ping() error {
	return nil
}

func main() {
	d := detective.New("your application")
	d1 := d.Dependency("cache")
	d1.Detect(func() error {
		time.Sleep(250 * time.Millisecond)
		return nil
	})

	d.Dependency("db").Detect(func() error {
		time.Sleep(250 * time.Millisecond)
		return errors.New("failed")
	})
	d.Endpoint("http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8081", d); err != nil {
			panic(err)
		}
	}()

	// Initialize a new detective instance
	d2 := detective.New("Another application")

	// Create a dependency, and register its detector function
	d2.Dependency("cache").Detect(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	if err := http.ListenAndServe(":8080", d2); err != nil {
		panic(err)
	}

}
