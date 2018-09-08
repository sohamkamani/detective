package detective

import (
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
			if got := s.WithDependencies(tt.args.dependencies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DependentState() = %v, want %v", got, tt.want)
			}
		})
	}
}
