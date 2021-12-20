package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_ViewByID(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		viewID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *GetViewResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful view returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"view":{"id": "fake-view-id"}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				viewID: "fake-view-id",
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
			},
			args: args{
				viewID: "fake-view-id",
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
			},
			args: args{
				viewID: "fake-view-id",
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
			_, err := c.ViewByID(context.Background(), tt.args.viewID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ViewByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_ViewsFor(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		viewListType ViewListType
		id           string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful views returned with valid type",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"views":[{"id": "fake-view-id"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				viewListType: TypeFolder,
				id:           "fake-folder-id",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.ViewsFor(context.Background(), tt.args.viewListType, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ViewsFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
