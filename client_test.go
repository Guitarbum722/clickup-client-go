// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

type mockAuthenticator struct{}

func (m *mockAuthenticator) AuthenticateFor(req *http.Request) error {
	return errors.New("mock error")
}

func TestClient_call(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx    context.Context
		method string
		uri    string
		data   *bytes.Buffer
		result interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success call",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"lists":[{"id":"12345678","name":"Hippy Dippy List"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				uri:    "/list",
				data:   nil,
				result: &ListsResponse{},
			},
			wantErr: false,
		},
		{
			name: "Error for unsupported method",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodPatch,
				uri:    "/list",
				data:   nil,
				result: &ListsResponse{},
			},
			wantErr: true,
		},
		{
			name: "Error for authenticator failure",
			fields: fields{
				doer:          nil,
				authenticator: &mockAuthenticator{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				uri:    "/list",
				data:   nil,
				result: &ListsResponse{},
			},
			wantErr: true,
		},
		{
			name: "Error for doer Do()",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("fake error")
				}),
				authenticator: &mockAuthenticator{},
			},
			args: args{
				ctx:    context.Background(),
				method: http.MethodGet,
				uri:    "/list",
				data:   nil,
				result: &ListsResponse{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
				baseURL:       tt.fields.baseURL,
			}
			if err := c.call(tt.args.ctx, tt.args.method, tt.args.uri, tt.args.data, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("Client.call() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		opts *ClientOpts
	}
	tests := []struct {
		name string
		args args
		want *Client
	}{
		{
			name: "Success new client with injected opts",
			args: args{
				opts: &ClientOpts{
					Doer: &http.Client{
						Timeout: time.Duration(time.Second * 2),
					},
					Authenticator: &mockAuthenticator{},
				},
			},
			want: &Client{
				doer: &http.Client{
					Timeout: time.Duration(time.Second * 2),
				},
				authenticator: &mockAuthenticator{},
				baseURL:       "https://api.clickup.com/api/v2",
			},
		},
		{
			name: "Success new client with defaults",
			args: args{
				opts: &ClientOpts{
					Doer:          nil,
					Authenticator: nil,
				},
			},
			want: &Client{
				doer: &http.Client{
					Timeout: time.Duration(time.Second * 20),
				},
				authenticator: &APITokenAuthenticator{},
				baseURL:       "https://api.clickup.com/api/v2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClient(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
