package clickup

import (
	"net/http"
	"time"
)

type APITokenAuthenticator struct {
	APIToken string
}

func (d *APITokenAuthenticator) AuthenticateFor(req *http.Request) {
	req.Header.Add("Authorization", d.APIToken)
}

type ClientOpts struct {
	// APIToken      string
	Doer          ClientDoer
	Authenticator Authenticator
}

type Client struct {
	doer          ClientDoer
	authenticator Authenticator
	baseURL       string
}

// wrapper for internal authenticator for convenience <shrug>.
func (c *Client) AuthenticateFor(req *http.Request) {
	c.authenticator.AuthenticateFor(req)
}

const basePath = "https://api.clickup.com/api/v2"

func NewClient(opts *ClientOpts) *Client {
	if opts.Doer != nil {
		return &Client{
			doer:          opts.Doer,
			authenticator: opts.Authenticator,
			baseURL:       basePath,
		}
	}

	return &Client{
		doer: &http.Client{
			Timeout: time.Duration(time.Second * 20),
		},
		authenticator: opts.Authenticator,
		baseURL:       basePath,
	}
}
