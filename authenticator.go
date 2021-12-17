// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"net/http"
)

type Authenticator interface {
	AuthenticateFor(*http.Request) error
}
