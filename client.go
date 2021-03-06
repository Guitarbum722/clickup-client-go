// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type APITokenAuthenticator struct {
	APIToken string
}

// AuthenticateFor adds an API token to the Authorization header of req.
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

// NewClient initializes and returns a pointer to a Client.
// If opts.Doer is not provided, an http.Client with a 20 seconds timeout is used.
// If opts.Authenticator is not provided, an APITokenAuthenticator is used.
func NewClient(opts *ClientOpts) *Client {
	auth := opts.Authenticator

	if auth == nil {
		auth = &APITokenAuthenticator{}
	}

	if opts.Doer != nil {
		return &Client{
			doer:          opts.Doer,
			authenticator: auth,
			baseURL:       basePath,
		}
	}

	return &Client{
		doer: &http.Client{
			Timeout: time.Duration(time.Second * 20),
		},
		authenticator: auth,
		baseURL:       basePath,
	}
}

func (c *Client) call(ctx context.Context, method, uri string, data *bytes.Buffer, result interface{}) error {
	var req *http.Request
	var err error

	endpoint := fmt.Sprintf("%s%s", c.baseURL, uri)

	switch method {
	case http.MethodGet:
		req, err = http.NewRequestWithContext(ctx, method, endpoint, nil)
		if err != nil {
			return err
		}

	case http.MethodPost:
		req, err = http.NewRequestWithContext(ctx, method, endpoint, data)
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			return err
		}
	case http.MethodPut:
		req, err = http.NewRequestWithContext(ctx, method, endpoint, data)
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			return err
		}
	case http.MethodDelete:
		req, err = http.NewRequestWithContext(ctx, method, endpoint, nil)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported http method")
	}

	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
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
