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

func TestClient_TeamsForWorkspace(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestSuccess Retrieve Teams for authenticated client",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"teams":[{"id":"1234567","name":"Schrute Farms"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			wantErr: false,
		},
		{
			name: "TestFailure unauthenticated client",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"err": "Token invalid", "ECODE": "OAUTH_025"}`
					return &http.Response{
						StatusCode: http.StatusUnauthorized,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
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
			_, err := c.TeamsForWorkspace(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Teams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GroupsForWorkspace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx              context.Context
		workspaceID      string
		optionalGroupIDs []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccess Retrieve Groups for Workspace",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"groups":[{"id":"1234567","name":"Schrute Farms"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				ctx:              context.Background(),
				workspaceID:      "test workspace id",
				optionalGroupIDs: nil,
			},
			wantErr: false,
		},
		{
			name: "TestSuccess Retrieve Groups for Workspace with group IDs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"groups":[{"id":"1234567","name":"Schrute Farms"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				ctx:              context.Background(),
				workspaceID:      "test workspace id",
				optionalGroupIDs: []string{"1234567"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
			}
			_, err := c.GroupsForWorkspace(tt.args.ctx, tt.args.workspaceID, tt.args.optionalGroupIDs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GroupsForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateGroup(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx   context.Context
		group CreateGroupRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateGroupResponse
		wantErr bool
	}{
		{
			name: "TestSuccess create new group",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"groups":[{"id":"1234567","name":"Schrute Farms"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				ctx: context.Background(),
				group: CreateGroupRequest{
					WorkspaceID: "123",
					Name:        "test group",
					MemberIDs:   nil,
				},
			},
			wantErr: false,
		},
		{
			name: "Fail create new group missing workspace id",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer:          nil,
			},
			args: args{
				ctx: context.Background(),
				group: CreateGroupRequest{
					WorkspaceID: "",
					Name:        "test group",
					MemberIDs:   nil,
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
			_, err := c.CreateGroup(tt.args.ctx, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateGroup(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx   context.Context
		group UpdateGroupRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UpdateGroupResponse
		wantErr bool
	}{
		{
			name: "TestSuccess update existing group",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"groups":[{"id":"1234567","name":"Schrute Farms"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				ctx: context.Background(),
				group: UpdateGroupRequest{
					ID:   "123",
					Name: "test group",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail update group missing group id",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer:          nil,
			},
			args: args{
				ctx: context.Background(),
				group: UpdateGroupRequest{
					ID:   "",
					Name: "test group",
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
			_, err := c.UpdateGroup(tt.args.ctx, tt.args.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteGroup(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx     context.Context
		groupID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success remove existing group",
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
				groupID: "test id",
			},
			wantErr: false,
		},
		{
			name: "Fail missing group ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:     context.Background(),
				groupID: "",
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
			if err := c.DeleteGroup(tt.args.ctx, tt.args.groupID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
