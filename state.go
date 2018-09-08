package detective

import (
	"errors"
)

type State struct {
	Name         string  `json:"name"`
	Ok           bool    `json:"active"`
	Status       string  `json:"status"`
	Dependencies []State `json:"dependencies,omitempty"`
}

func (s State) WithError(err error) State {
	if err == nil {
		return s.WithOk()
	}
	ns := s
	ns.Ok = false
	ns.Status = "Error: " + err.Error()
	return ns
}

func (s State) WithOk() State {
	ns := s
	ns.Ok = true
	ns.Status = "Ok"
	return ns
}

func (s State) WithDependencies(dependencies []State) State {
	finalState := s
	finalState.Dependencies = dependencies
	if !NoErrors(dependencies) {
		return finalState.WithError(errors.New("dependency failure"))
	}
	return finalState.WithOk()
}

func NoErrors(states []State) bool {
	for i := range states {
		if !states[i].Ok {
			return false
		}
	}
	return true
}
