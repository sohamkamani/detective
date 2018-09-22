package detective

import (
	"encoding/json"
	"net/http"
	"sync"
)

// A Detective instance manages registered dependencies and endpoints.
// Dependencies can be registered with an instance.
// Each instance has a state which represents the health of its components.
type Detective struct {
	name         string
	client       Doer
	dependencies []*Dependency
	endpoints    []*endpoint
}

// New creates a new Detective instance. To avoid confusion, the name provided should preferably be unique among dependent detective instances.
func New(name string) *Detective {
	return &Detective{
		name:   name,
		client: &http.Client{},
	}
}

// WithHTTPClient sets the HTTP Client to be used while hitting the endpoint of another detective HTTP ping handler.
func (d *Detective) WithHTTPClient(c Doer) *Detective {
	d.client = c
	return d
}

// Dependency adds a new dependency to the Detective instance. The name provided should preferably be unique among dependencies registered within the same detective instance.
func (d *Detective) Dependency(name string) *Dependency {
	dependency := newDependency(name)
	d.dependencies = append(d.dependencies, dependency)
	return dependency
}

// Endpoint adds an HTTP endpoint as a dependency to the Detective instance, thereby allowing you to compose detective instances. This method creates a GET request to the provided url. If you want to customize the request (like using a different HTTP method, or adding headers), consider using the EndpointReq method instead.
func (d *Detective) Endpoint(url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	d.EndpointReq(req)
	return nil
}

// EndpointReq is similar to Endpoint, but takes an HTTP request object instead of a URL. Use this method if you want to customize the request to the ping handler of another detective instance.
func (d *Detective) EndpointReq(req *http.Request) {
	e := &endpoint{
		name:   d.name,
		client: d.client,
		req:    *req,
	}
	d.endpoints = append(d.endpoints, e)
}

func (d *Detective) getState() State {
	depLength := len(d.dependencies)
	totalDependencyLength := depLength + len(d.endpoints)
	subStates := make([]State, totalDependencyLength)
	var wg sync.WaitGroup
	wg.Add(totalDependencyLength)
	for iDep, dep := range d.dependencies {
		go func(dep *Dependency, i int) {
			s := dep.getState()
			subStates[i] = s
			wg.Done()
		}(dep, iDep)
	}
	for iEp, e := range d.endpoints {
		go func(e *endpoint, i int) {
			s := e.getState()
			subStates[depLength+i] = s
			wg.Done()
		}(e, iEp)
	}
	wg.Wait()
	s := State{Name: d.name}
	return s.withDependencies(subStates)
}

// ServeHTTP is the HTTP handler function for getting the state of the Detective instance
func (d *Detective) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := d.getState()
	sBody, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(sBody)
	return
}
