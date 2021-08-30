package clickup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// clickup doesn't document the structure of their error response and it also seems to return a 200 when they say it won't -_-
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("failed to parse response from clickup: %v %v", err, string(b))
	}
	return nil
	// return json.NewDecoder(res.Body).Decode(result)
}
