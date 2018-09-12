package detective

import (
	"time"
)

// The DetectorFunc type represents the function signature to check the health of a dependency
type DetectorFunc func() error

// The Dependency type represents a detectable unit. The function provided in the Detect method will be called to monitor the state of the dependency
type Dependency struct {
	name     string
	detector DetectorFunc
	state    State
}

func noopDetectorFunc() DetectorFunc {
	return func() error {
		return nil
	}
}

func newDependency(name string) *Dependency {
	return &Dependency{
		name:     name,
		detector: noopDetectorFunc(),
	}
}

// Detect registers a function that will be called to detect the health of a dependency. If the dependency is healthy, a nil value should be returned as the error.
func (d *Dependency) Detect(df DetectorFunc) {
	d.detector = df
}

func (d *Dependency) updateState() {
	d.state = d.getState()
}

func (d *Dependency) getState() State {
	init := time.Now()
	err := d.detector()
	diff := time.Now().Sub(init)
	s := State{Name: d.name, Latency: diff}
	if err != nil {
		return s.withError(err)
	}
	return s.withOk()
}
