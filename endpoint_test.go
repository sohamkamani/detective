package detective

import (
	"errors"
	dm "github.com/sohamkamani/detective/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		httpStatus    int
		httpError     error
		jsonResponse  string
		expectedState State
	}{
		{
			name:       "success",
			httpStatus: http.StatusOK,
			jsonResponse: `{
				"name":"sample",
				"active": true,
				"status":"Ok",
				"dependencies":[
					{
						"name":"dep1",
						"active": true,
						"status":"Ok"
					}
				]
			}`,
			expectedState: State{
				Name:   "sample",
				Ok:     true,
				Status: "Ok",
				Dependencies: []State{
					{
						Name:   "dep1",
						Ok:     true,
						Status: "Ok",
					},
				},
			},
		},
		{
			name:       "http error status",
			httpStatus: http.StatusInternalServerError,
			expectedState: State{
				Name:   "sample",
				Ok:     false,
				Status: "Error: service sample returned http status: 500 Internal Server Error",
			},
		},
		{
			name:       "empty body",
			httpStatus: http.StatusOK,
			expectedState: State{
				Name:   "sample",
				Ok:     false,
				Status: "Error: service sample returned no response body",
			},
		},
		{
			name:       "random response",
			httpStatus: http.StatusOK,
			jsonResponse: `{
				"some":"random",
				"response": 0
			}`,
			expectedState: State{Name: "", Ok: false, Status: ""},
		},
		{
			name:       "incorrect json",
			httpStatus: http.StatusOK,
			jsonResponse: `{
				"some":"random"`,
			expectedState: State{Name: "sample", Ok: false, Status: "Error: unexpected EOF"},
		},
		{
			name:          "http failure",
			httpError:     errors.New("failed"),
			expectedState: State{Name: "sample", Ok: false, Status: "Error: failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &dm.MockClient{}
			mockClient.On("Do", mock.Anything).Return(dm.MockJSONResponse(tt.jsonResponse, tt.httpStatus), tt.httpError)

			req, err := http.NewRequest(http.MethodGet, "http://mock.com/", nil)
			require.NoError(t, err)

			e := &endpoint{
				name:   "sample",
				client: mockClient,
				req:    *req,
			}

			s := e.getState("")
			assertStatesEqual(t, tt.expectedState, s)
		})
	}
}
