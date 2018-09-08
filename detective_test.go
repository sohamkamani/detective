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

	s := d.getState()
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

	t.Run("handler", func(t *testing.T) {
		handler := d.Handler()
		mockClient.On("Do", mock.Anything).Return(dm.MockJSONResponse(`{"name":"sample","status":"Ok", "active":true}`, http.StatusOK), nil).Once()
		rw := &httptest.ResponseRecorder{Body: bytes.NewBuffer([]byte{})}
		req, err := http.NewRequest(http.MethodGet, "", nil)
		require.NoError(t, err)
		handler(rw, req)
		res := rw.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		var gotState State
		json.NewDecoder(res.Body).Decode(&gotState)
		defer res.Body.Close()
		assertStatesEqual(t, expectedState, gotState)
	})
}
