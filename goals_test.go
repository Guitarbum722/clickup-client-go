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

func TestClient_CreateGoal(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx  context.Context
		goal CreateGoalRequest
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
					body := `{"name": "test goal"}`
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
				goal: CreateGoalRequest{
					WorkspaceID: "Test WorkspaceID",
					Name:        "test goal",
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
				goal: CreateGoalRequest{
					WorkspaceID: "",
					Name:        "test goal",
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
			_, err := c.CreateGoal(tt.args.ctx, tt.args.goal)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateGoal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateGoal(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx  context.Context
		goal UpdateGoalRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success goal id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test goal update"}`
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
				goal: UpdateGoalRequest{
					GoalID: "Test GoalID",
					Name:   "test goal update",
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
				goal: UpdateGoalRequest{
					GoalID: "",
					Name:   "test goal",
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
			_, err := c.UpdateGoal(tt.args.ctx, tt.args.goal)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateGoal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GoalsForWorkspace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx              context.Context
		workspaceID      string
		includeCompleted bool
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
					body := `{"name": "test goal"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:              context.Background(),
				workspaceID:      "test workspace id",
				includeCompleted: false,
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
				ctx:              context.Background(),
				workspaceID:      "",
				includeCompleted: false,
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
			_, err := c.GoalsForWorkspace(tt.args.ctx, tt.args.workspaceID, tt.args.includeCompleted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GoalsForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_GoalForWorkSpace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx    context.Context
		goalID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GoalResponse
		wantErr bool
	}{
		{
			name: "Success goal id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test goal"}`
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
				goalID: "test goal id",
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
				ctx:    context.Background(),
				goalID: "",
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
			_, err := c.GoalForWorkSpace(tt.args.ctx, tt.args.goalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GoalForWorkSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteGoal(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx    context.Context
		goalID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success goal id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test goal"}`
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
				goalID: "test goal ID",
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
				ctx:    context.Background(),
				goalID: "",
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
			if err := c.DeleteGoal(tt.args.ctx, tt.args.goalID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteGoal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateKeyResultForGoal(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx       context.Context
		keyResult CreateKeyResultRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CreateKeyResultResponse
		wantErr bool
	}{
		{
			name: "Success goal id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test kr"}`
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
				keyResult: CreateKeyResultRequest{
					GoalID: "test goalID",
					Name:   "test kr",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing goal ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				keyResult: CreateKeyResultRequest{
					GoalID: "",
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
			_, err := c.CreateKeyResultForGoal(tt.args.ctx, tt.args.keyResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateKeyResultForGoal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateKeyResult(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx       context.Context
		keyResult UpdateKeyResultRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UpdateKeyResultResponse
		wantErr bool
	}{
		{
			name: "Success key result id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test kr"}`
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
				keyResult: UpdateKeyResultRequest{
					ID:   "test goalID",
					Name: "test kr",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing Key result ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				keyResult: UpdateKeyResultRequest{
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
			_, err := c.UpdateKeyResult(tt.args.ctx, tt.args.keyResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateKeyResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteKeyResult(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx         context.Context
		keyResultID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success key result id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test goal"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:         context.Background(),
				keyResultID: "test key result ID",
			},
			wantErr: false,
		},
		{
			name: "Fail missing key result ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:         context.Background(),
				keyResultID: "",
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
			if err := c.DeleteKeyResult(tt.args.ctx, tt.args.keyResultID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteKeyResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
