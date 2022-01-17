package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_SharedHierarchy(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx         context.Context
		workspaceID string
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
					body := `{"shared": {"lists": [{"id": "123"}]}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				ctx:         context.Background(),
				workspaceID: "test workspace id",
			},
			wantErr: false,
		},
		{
			name: "TestSuccess Retrieve Groups for Workspace with group IDs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer:          nil,
			},
			args: args{
				ctx:         context.Background(),
				workspaceID: "",
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
			_, err := c.SharedHierarchy(tt.args.ctx, tt.args.workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.SharedHierarchy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
