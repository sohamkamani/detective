package detective

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Detective struct {
	name         string
	client       Doer
	dependencies []*Dependency
	endpoints    []*Endpoint
}

func New(name string) *Detective {
	return &Detective{
		name:   name,
		client: &http.Client{},
	}
}

func (d *Detective) Dependency(name string) *Dependency {
	dependency := NewDependency(name)
	d.dependencies = append(d.dependencies, dependency)
	return dependency
}

func (d *Detective) Endpoint(url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	d.EndpointReq(req)
	return nil
}

func (d *Detective) EndpointReq(req *http.Request) {
	e := &Endpoint{
		client: d.client,
		req:    *req,
	}
	d.endpoints = append(d.endpoints, e)
}

func (d *Detective) getState() State {
	totalDependencyLength := len(d.dependencies) + len(d.endpoints)
	subStates := make([]State, 0, totalDependencyLength)
	var wg sync.WaitGroup
	wg.Add(totalDependencyLength)
	for _, dep := range d.dependencies {
		go func() {
			s := dep.getState()
			subStates = append(subStates, s)
			wg.Done()
		}()
	}
	for _, e := range d.endpoints {
		go func() {
			s := e.getState()
			subStates = append(subStates, s)
			wg.Done()
		}()
	}
	wg.Wait()
	return DependentState(d.name, subStates)
}

func (d *Detective) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := d.getState()
		sBody, err := json.Marshal(s)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(sBody)
		return
	}
}
