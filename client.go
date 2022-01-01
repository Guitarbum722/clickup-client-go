// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type APITokenAuthenticator struct {
	APIToken string
}

func (d *APITokenAuthenticator) AuthenticateFor(req *http.Request) error {
	req.Header.Add("Authorization", d.APIToken)
	return nil
}

type ClientOpts struct {
	Doer          ClientDoer
	Authenticator Authenticator
}

type Client struct {
	doer          ClientDoer
	authenticator Authenticator
	baseURL       string
}

// wrapper for internal authenticator for convenience <shrug>.
func (c *Client) AuthenticateFor(req *http.Request) error {
	return c.authenticator.AuthenticateFor(req)
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

func (c *Client) call(method, uri string, data *bytes.Buffer, result interface{}) error {
	var req *http.Request
	var err error

	endpoint := fmt.Sprintf("%s%s", c.baseURL, uri)

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
	case http.MethodPut:
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
	default:
		return errors.New("unsupported http method")
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
