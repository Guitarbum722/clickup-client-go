package clickup

import "net/http"

type ClientDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
