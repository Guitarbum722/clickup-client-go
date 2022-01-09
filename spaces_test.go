// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_SpacesForWorkspace(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		teamID          string
		includeArchived bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful spaces returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"spaces":[{"id":"14865529","name":"John's Happy Place"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				teamID:          "fakeTeamID",
				includeArchived: false,
			},
			wantErr: false,
		},
		{
			name: "TestUnauthorized-Invalid-API-Key",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"err": "Token invalid", "ECODE": "OAUTH_025"}`
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				teamID:          "fakeTeamID",
				includeArchived: false,
			},
			wantErr: true,
		},
		{
			name: "TestUnknownErrorResponseStructure200",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{{{{"first": "Ned", "last": "Schneebly"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				teamID:          "fakeTeamID",
				includeArchived: false,
			},
			wantErr: true,
		},
		{
			name: "TestUnknownErrorResponseStructure500",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{{{{"first": "Ned", "last": "Schneebly"}`
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				teamID:          "fakeTeamID",
				includeArchived: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.SpacesForWorkspace(context.Background(), tt.args.teamID, tt.args.includeArchived)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpacesForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_SpaceByID(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		spaceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful space returned by ID",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":"123","name":"John's Happy Place"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				spaceID: "fakeSpaceID",
			},
			wantErr: false,
		},
		{
			name: "TestUnauthorized-Invalid-API-Key",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"err": "Token invalid", "ECODE": "OAUTH_025"}`
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				spaceID: "fakeTeamID",
			},
			wantErr: true,
		},
		{
			name: "TestUnknownErrorResponseStructure-200",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{{{{"first": "Ned", "last": "Schneebly"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				spaceID: "fakeTeamID",
			},
			wantErr: true,
		},
		{
			name: "TestUnknownErrorResponseStructure-s500",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{{{{"first": "Ned", "last": "Schneebly"}`
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				spaceID: "fakeTeamID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.SpaceByID(context.Background(), tt.args.spaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SpaceByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestClient_DeleteSpace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx     context.Context
		spaceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success space id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:     context.Background(),
				spaceID: "test id",
			},
			wantErr: false,
		},
		{
			name: "Fail missing space ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:     context.Background(),
				spaceID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
			}
			if err := c.DeleteSpace(tt.args.ctx, tt.args.spaceID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteSpace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateSpaceForWorkspace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx   context.Context
		space UpdateSpaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success space id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id": "test id"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				space: UpdateSpaceRequest{
					ID: "test id",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing space ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				space: UpdateSpaceRequest{
					ID: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
			}
			_, err := c.UpdateSpaceForWorkspace(tt.args.ctx, tt.args.space)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateSpaceForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateSpaceForWorkspace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx   context.Context
		space CreateSpaceRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success workspace id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id": "test id"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				space: CreateSpaceRequest{
					WorkspaceID: "test id",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing workspace ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				space: CreateSpaceRequest{
					WorkspaceID: "",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
			}
			_, err := c.CreateSpaceForWorkspace(tt.args.ctx, tt.args.space)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateSpaceForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
