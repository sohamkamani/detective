package detective

import (
	"errors"
)

type State struct {
	Name         string  `json:"name"`
	Ok           bool    `json:"active"`
	Status       string  `json:"status"`
	Dependencies []State `json:"dependencies"`
}

func ErrorState(name string, err error) State {
	if err == nil {
		return NormalState(name)
	}
	return State{
		Name:   name,
		Ok:     false,
		Status: "Error: " + err.Error(),
	}
}

func NormalState(name string) State {
	return State{
		Name:   name,
		Ok:     true,
		Status: "Ok",
	}
}

func DependentState(name string, dependencies []State) State {
	var finalState State
	if !NoErrors(dependencies) {
		finalState = ErrorState(name, errors.New("dependency failure"))
	} else {
		finalState = NormalState(name)
	}
	finalState.Name = name
	finalState.Dependencies = dependencies
	return finalState
}

func NoErrors(states []State) bool {
	for i := range states {
		if !states[i].Ok {
			return false
		}
	}
	return true
}
