package main

import (
	// "errors"
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

/*
In this example, `d` and `d2` are registered as dependencies of each other
Detective helps fix circular dependencies by appending information about the calling instance onto the "X_DETECTIVE_FROM_CHAIN"
http header.
*/
func main() {
	d := detective.New("your application")
	dep1 := d.Dependency("cache")
	dep1.Detect(func() error {
		time.Sleep(250 * time.Millisecond)
		return nil
	})

	d.Dependency("db").Detect(func() error {
		time.Sleep(250 * time.Millisecond)
		return nil //errors.New("failed")
	})

	// Here, the detective instance running on port 8080 (d2) is registered as a dependency
	d.Endpoint("http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8081", d); err != nil {
			panic(err)
		}
	}()

	d2 := detective.New("Another application")

	d2.Dependency("cache").Detect(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	// At the same time, the first instance (d) is registered as a dependency of d2 as well
	d2.Endpoint("http://localhost:8081")
	if err := http.ListenAndServe(":8080", d2); err != nil {
		panic(err)
	}

}
