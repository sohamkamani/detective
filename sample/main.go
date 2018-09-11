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
	d := detective.New("sample")
	d1 := d.Dependency("cache")
	d1.Detect(func() error {
		time.Sleep(2500 * time.Millisecond)
		return nil
	})
	d.Endpoint("http://localhost:8080")
	go func() {
		if err := http.ListenAndServe(":8081", d); err != nil {
			panic(err)
		}
	}()

	d2 := detective.New("sample2")
	g1 := d2.Dependency("db")
	g1.Detect(func() error {
		time.Sleep(500 * time.Millisecond)
		return errors.New("dkcndkcn")
	})

	if err := http.ListenAndServe(":8080", d2); err != nil {
		panic(err)
	}

}
