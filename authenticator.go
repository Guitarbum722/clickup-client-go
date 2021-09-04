package clickup

import "net/http"

type Authenticator interface {
	AuthenticateFor(req *http.Request)
}
