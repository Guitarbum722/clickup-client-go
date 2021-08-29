// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

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
