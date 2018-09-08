package detective

import (
	"bytes"
	dm "github.com/sohamkamani/detective/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
	assert.Equal(t, expectedState, s)
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
		b, err := ioutil.ReadAll(res.Body)
		assert.JSONEq(t, `{"status":"Ok", "dependencies":[{"active":true, "status":"Ok", "name":"sampledep"}, {"name":"sample", "active":true, "status":"Ok"}], "name":"sample", "active":true}`, string(b))
	})
}
