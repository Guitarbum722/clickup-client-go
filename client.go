package clickup

import (
	"net/http"
	"time"
)

type ClientOpts struct {
	APIToken string
	Doer     ClientDoer
}

type Client struct {
	doer    ClientDoer
	opts    *ClientOpts
	baseURL string
}

const basePath = "https://api.clickup.com/api/v2"

func NewClient(opts *ClientOpts) *Client {
	if opts.Doer != nil {
		return &Client{
			doer:    opts.Doer,
			opts:    opts,
			baseURL: basePath,
		}
	}

	return &Client{
		doer: &http.Client{
			Timeout: time.Duration(time.Second * 20),
		},
		opts:    opts,
		baseURL: basePath,
	}
}
