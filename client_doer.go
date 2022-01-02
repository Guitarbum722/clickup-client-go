// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import "net/http"

type ClientDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
