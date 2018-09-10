package detective

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Endpoint struct {
	name   string
	req    http.Request
	client Doer
}

func (e *Endpoint) getState() State {
	init := time.Now()
	currentReq := e.req
	res, err := e.client.Do(&currentReq)
	diff := time.Now().Sub(init)
	s := State{Name: e.name, Latency: diff}
	if err != nil {
		return s.WithError(err)
	}
	if res.StatusCode != http.StatusOK {
		return s.WithError(errors.New("service " + e.name + " returned http status: " + res.Status))
	}
	if res.Body == nil {
		return s.WithError(errors.New("service " + e.name + " returned no response body"))
	}
	defer res.Body.Close()
	var state State
	if err := json.NewDecoder(res.Body).Decode(&state); err != nil {
		return s.WithError(err)
	}
	return state.WithOk()
}
