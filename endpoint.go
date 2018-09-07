package detective

import (
	"encoding/json"
	"errors"
	"net/http"
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
	currentReq := e.req
	res, err := e.client.Do(&currentReq)
	if err != nil {
		return ErrorState(e.name, err)
	}
	if res.StatusCode != http.StatusOK {
		return ErrorState(e.name, errors.New("service "+e.name+" returned http status: "+res.Status))
	}
	var state State
	if err := json.NewDecoder(res.Body).Decode(&state); err != nil {
		return ErrorState(e.name, err)
	}
	return state
}
