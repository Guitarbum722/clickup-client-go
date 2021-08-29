package clickup

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const baseURL = "https://api.clickup.com/api/v2"

type ClientOpts struct {
	APIToken   string
	HTTPClient *http.Client
}

type Client struct {
	hc   *http.Client
	opts *ClientOpts
}

func NewClient(opts *ClientOpts) *Client {
	if opts.HTTPClient != nil {
		return &Client{
			hc:   opts.HTTPClient,
			opts: opts,
		}
	}

	return &Client{
		hc: &http.Client{
			Timeout: time.Duration(time.Second * 20),
		},
		opts: opts,
	}
}

// TODO: put this in a submodule for information hiding?
func (c *Client) call(method, uri string, data *bytes.Buffer, result interface{}) error {
	var req *http.Request
	var err error

	endpoint := baseURL + uri

	switch method {
	case http.MethodGet:
		req, err = http.NewRequest(method, endpoint, nil)
		if err != nil {
			return err
		}

	case http.MethodPost:
		req, err = http.NewRequest(method, endpoint, data)
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			return err
		}
	case http.MethodDelete:
		req, err = http.NewRequest(method, endpoint, nil)
		if err != nil {
			return err
		}
	}
	req.Header.Set("Authorization", c.opts.APIToken)

	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(result)
}
