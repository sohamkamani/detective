package detective

import (
	"bytes"
	"encoding/json"
	dm "github.com/sohamkamani/detective/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetective(t *testing.T) {

	t.Run("get state", func(t *testing.T) {

		mockClient := &dm.MockClient{}
		mockClient.On("Do", mock.Anything).Return(dm.MockJSONResponse(`{"name":"sample","status":"Ok", "active":true}`, http.StatusOK), nil).Once()
		d := New("sample").WithHTTPClient(mockClient)
		d.Endpoint("http://sample")
		dep := d.Dependency("sampledep")
		depCalled := false
		dep.Detect(func() error {
			depCalled = true
			return nil
		})

		s := d.getState([]string{})
		expectedState := State{
			Name:   "sample",
			Ok:     true,
			Status: "Ok",
			Dependencies: []State{
				State{Name: "sampledep", Ok: true, Status: "Ok"},
				State{Name: "sample", Ok: true, Status: "Ok"},
			},
		}
		assertStatesEqual(t, expectedState, s)
		assert.True(t, depCalled)
	})

	t.Run("handler", func(t *testing.T) {
		mockClient := &dm.MockClient{}
		d := New("sample").WithHTTPClient(mockClient)
		d.Endpoint("http://sample")
		d.Dependency("sampledep")
		mockClient.On("Do", mock.Anything).Return(dm.MockJSONResponse(`{"name":"sample","status":"Ok", "active":true}`, http.StatusOK), nil).Once()
		rw := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
		req, err := http.NewRequest(http.MethodGet, "", nil)
		req.Header.Add(fromHeader, "abc|def")
		require.NoError(t, err)
		d.ServeHTTP(rw, req)
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		var gotState State
		json.NewDecoder(res.Body).Decode(&gotState)
		defer res.Body.Close()
		expectedState := State{
			Name:   "sample",
			Ok:     true,
			Status: "Ok",
			Dependencies: []State{
				State{Name: "sampledep", Ok: true, Status: "Ok"},
				State{Name: "sample", Ok: true, Status: "Ok"},
			},
		}
		assertStatesEqual(t, expectedState, gotState)

		clientRequest, ok := mockClient.Calls[0].Arguments[0].(*http.Request)
		require.True(t, ok)
		assert.Equal(t, "abc|def|sample", clientRequest.Header.Get(fromHeader))
	})

	t.Run("handler with name in fromHeader chain", func(t *testing.T) {
		mockClient := &dm.MockClient{}
		d := New("sample").WithHTTPClient(mockClient)
		d.Endpoint("http://sample")
		d.Dependency("sampledep")
		mockClient.On("Do", mock.Anything).Return(dm.MockJSONResponse(`{"name":"sample2","status":"Ok", "active":true}`, http.StatusOK), nil).Once()
		rw := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
		req, err := http.NewRequest(http.MethodGet, "", nil)
		req.Header.Add(fromHeader, "abc|sample|ghi|")
		require.NoError(t, err)
		d.ServeHTTP(rw, req)
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		var gotState State
		json.NewDecoder(res.Body).Decode(&gotState)
		defer res.Body.Close()
		expectedState := State{
			Name:   "sample",
			Ok:     true,
			Status: "Ok",
			Dependencies: []State{
				State{Name: "sampledep", Ok: true, Status: "Ok"},
			},
		}
		assertStatesEqual(t, expectedState, gotState)
	})
}
