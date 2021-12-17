// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"context"
	"net/http"
	"time"
)

type APITokenAuthenticator struct {
	APIToken string
}

func (d *APITokenAuthenticator) AuthenticateFor(ctx context.Context, req *http.Request) error {
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
func (c *Client) AuthenticateFor(ctx context.Context, req *http.Request) error {
	return c.authenticator.AuthenticateFor(ctx, req)
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
