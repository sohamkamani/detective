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
	// Our mock database
}

func (d db) ping() error {
	return nil
}

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

	if err := http.ListenAndServe(":8081", d); err != nil {
		panic(err)
	}

}
