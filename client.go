package clickup

import (
	"net/http"
	"time"
)

type ClientOpts struct {
	APIToken   string
	HTTPClient *http.Client
}

type Client struct {
	hc      *http.Client
	opts    *ClientOpts
	baseURL string
}

func NewClient(opts *ClientOpts) *Client {
	if opts.HTTPClient != nil {
		return &Client{
			hc:      opts.HTTPClient,
			opts:    opts,
			baseURL: "https://api.clickup.com/api/v2",
		}
	}

	return &Client{
		hc: &http.Client{
			Timeout: time.Duration(time.Second * 20),
		},
		opts:    opts,
		baseURL: "https://api.clickup.com/api/v2",
	}
}
