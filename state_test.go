package detective

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDependentState(t *testing.T) {
	type args struct {
		name         string
		dependencies []State
	}
	tests := []struct {
		name string
		args args
		want State
	}{
		{
			name: "returns ok state if all dependencies are successful",
			args: args{
				name: "sample",
				dependencies: []State{
					State{
						Name: "state1",
						Ok:   true,
					},
					State{
						Name: "state2",
						Ok:   true,
					},
				},
			},
			want: State{
				Name:   "sample",
				Status: "Ok",
				Ok:     true,
				Dependencies: []State{
					State{
						Name: "state1",
						Ok:   true,
					},
					State{
						Name: "state2",
						Ok:   true,
					},
				},
			},
		},
		{
			name: "returns error state if some dependencies are unsuccessful",
			args: args{
				name: "sample",
				dependencies: []State{
					State{
						Name: "state1",
						Ok:   true,
					},
					State{
						Name: "state2",
						Ok:   false,
					},
				},
			},
			want: State{
				Name:   "sample",
				Status: "Error: dependency failure",
				Ok:     false,
				Dependencies: []State{
					State{
						Name: "state1",
						Ok:   true,
					},
					State{
						Name: "state2",
						Ok:   false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := State{Name: tt.args.name}
			if got := s.withDependencies(tt.args.dependencies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DependentState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertStatesEqual(t *testing.T, s1, s2 State) {
	assert.Equal(t, s1.Name, s2.Name)
	assert.Equal(t, s1.Ok, s2.Ok)
	assert.Equal(t, s1.Status, s2.Status)
	require.Equal(t, len(s1.Dependencies), len(s2.Dependencies))
	for i := range s1.Dependencies {
		assertStatesEqual(t, s1.Dependencies[i], s2.Dependencies[i])
	}
}
