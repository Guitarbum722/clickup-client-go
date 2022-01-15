// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"net/http"
)

// Authenticator adds authentication details to an http.Request.
// Typically the Authorization header will be added by the implementation.
type Authenticator interface {
	AuthenticateFor(*http.Request) error
}
