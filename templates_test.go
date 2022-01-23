package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_TemplatesForWorkspace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx         context.Context
		workspaceID string
		page        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TemplatesResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful templates returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"templates":[{"id":"14865529","name":"John's Happy Template"}]}`
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
				workspaceID: "test workspace id",
				page:        0,
			},
			wantErr: false,
		},
		{
			name: "Fail missing workspace id",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:         context.Background(),
				workspaceID: "",
				page:        0,
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
			_, err := c.TemplatesForWorkspace(tt.args.ctx, tt.args.workspaceID, tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TemplatesForWorkspace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateTaskFromTemplate(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx  context.Context
		task TaskFromTemplateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful create task from template",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"name": "test task from template"}`
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
				task: TaskFromTemplateRequest{
					ListID:     "list id",
					TemplateID: "template id",
					Name:       "test task",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing template id",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				task: TaskFromTemplateRequest{
					ListID:     "list id",
					TemplateID: "",
					Name:       "test task",
				},
			},
			wantErr: true,
		},
		{
			name: "Fail missing list id",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				task: TaskFromTemplateRequest{
					ListID:     "",
					TemplateID: "template id",
					Name:       "test task",
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
			_, err := c.CreateTaskFromTemplate(tt.args.ctx, tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateTaskFromTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
