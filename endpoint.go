package detective

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Doer represents the standard HTTP client interface
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type endpoint struct {
	name   string
	req    http.Request
	client Doer
}

func (e *endpoint) getState() State {
	init := time.Now()
	currentReq := e.req
	res, err := e.client.Do(&currentReq)
	diff := time.Now().Sub(init)
	s := State{Name: e.name, Latency: diff}
	if err != nil {
		return s.withError(err)
	}
	if res.StatusCode != http.StatusOK {
		return s.withError(errors.New("service " + e.name + " returned http status: " + res.Status))
	}
	if res.Body == nil {
		return s.withError(errors.New("service " + e.name + " returned no response body"))
	}
	defer res.Body.Close()
	var state State
	if err := json.NewDecoder(res.Body).Decode(&state); err != nil {
		return s.withError(err)
	}
	return state.withOk()
}
