package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_AddDependencyForTask(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx        context.Context
		dependency AddDependencyRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful create dependency link",
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
				ctx: context.Background(),
				dependency: AddDependencyRequest{
					TaskID:      "test task",
					DependsOn:   "dependency task",
					WorkspaceID: "workspace",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing Depends On or Dependency of",
			fields: fields{
				doer:          nil,
				authenticator: nil,
			},
			args: args{
				ctx: context.Background(),
				dependency: AddDependencyRequest{
					TaskID:      "test task",
					WorkspaceID: "workspace",
				},
			},
			wantErr: true,
		},
		{
			name: "Depends on and dependency of both provided",
			fields: fields{
				doer:          nil,
				authenticator: nil,
			},
			args: args{
				ctx: context.Background(),
				dependency: AddDependencyRequest{
					TaskID:       "test task",
					DependsOn:    "depends on task",
					DependencyOf: "dependency of task",
					WorkspaceID:  "workspace",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing workspace id",
			fields: fields{
				doer:          nil,
				authenticator: nil,
			},
			args: args{
				ctx: context.Background(),
				dependency: AddDependencyRequest{
					TaskID:       "test task",
					DependsOn:    "depends on task",
					DependencyOf: "dependency of task",
				},
			},
			wantErr: true,
		},
		{
			name: "Missing task id",
			fields: fields{
				doer:          nil,
				authenticator: nil,
			},
			args: args{
				ctx: context.Background(),
				dependency: AddDependencyRequest{
					DependsOn:    "depends on task",
					DependencyOf: "dependency of task",
					WorkspaceID:  "workspace",
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
			if err := c.AddDependencyForTask(tt.args.ctx, tt.args.dependency); (err != nil) != tt.wantErr {
				t.Errorf("Client.AddDependencyForTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_AddTaskLink(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx  context.Context
		link AddTaskLinkRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TaskLinkResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful create task link",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"task": {"id": "task id", "linked_tasks": [{"task_id": "linked task id"}]}}`
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
				link: AddTaskLinkRequest{
					TaskID:        "task id",
					LinksToTaskID: "linkedtask id",
					WorkspaceID:   "workspace id",
				},
			},
			wantErr: false,
		},
		{
			name: "Missing task id",
			fields: fields{
				doer:          nil,
				authenticator: nil,
			},
			args: args{
				ctx: context.Background(),
				link: AddTaskLinkRequest{
					WorkspaceID:   "workspace",
					LinksToTaskID: "linked task id",
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
			_, err := c.AddTaskLinkForTask(tt.args.ctx, tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.AddTaskLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
