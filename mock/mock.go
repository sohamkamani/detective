package mock

import (
	"bytes"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	result := m.Called(req)
	res := result.Get(0)
	err := result.Error(1)
	if res == nil {
		return nil, err
	}
	return res.(*http.Response), err
}

func MockJSONResponse(body string, status int) *http.Response {
	r := httptest.ResponseRecorder{}
	if body != "" {
		r.Body = bytes.NewBuffer([]byte(body))
	}
	r.WriteHeader(status)
	return r.Result()
}
