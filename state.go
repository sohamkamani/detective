package detective

import (
	"errors"
	"fmt"
	"time"
)

// State describes the current status of an entity. This entity can be a Dependency, or a Detective instance. A State can contain other States as well.
type State struct {
	Name         string        `json:"name"`
	Ok           bool          `json:"active"`
	Status       string        `json:"status"`
	Latency      time.Duration `json:"latency"`
	Dependencies []State       `json:"dependencies,omitempty"`
}

func (s State) withError(err error) State {
	if err == nil {
		return s.withOk()
	}
	ns := s
	ns.Ok = false
	ns.Status = "Error: " + err.Error()
	return ns
}

func (s State) withOk() State {
	ns := s
	ns.Ok = true
	ns.Status = "Ok"
	return ns
}

func (s State) withDependencies(dependencies []State) State {
	finalState := s
	finalState.Dependencies = dependencies
	if !noErrors(dependencies) {
		return finalState.withError(errors.New("dependency failure"))
	}
	return finalState.withOk()
}

func noErrors(states []State) bool {
	for i := range states {
		if !states[i].Ok {
			return false
		}
	}
	return true
}
