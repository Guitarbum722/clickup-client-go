package clickup

import "net/http"

type clientDoerFunc func(request *http.Request) (*http.Response, error)

type mockHTTPClient struct {
	doFn clientDoerFunc
}

func (m *mockHTTPClient) Do(request *http.Request) (*http.Response, error) {
	return m.doFn(request)
}

func newMockClientDoer(fn clientDoerFunc) *mockHTTPClient {
	return &mockHTTPClient{fn}
}
