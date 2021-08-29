package clickup

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_ListsForFolder(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		folderID        string
		includeArchived bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful lists returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"lists":[{"id":"12345678","name":"Hippy Dippy List"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				folderID: "fakeFolderID",
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
				folderID: "fakeTeamID",
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
				folderID: "fakeTeamID",
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
				folderID: "fakeTeamID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:    tt.fields.doer,
				opts:    tt.fields.opts,
				baseURL: tt.fields.baseURL,
			}
			_, err := c.ListsForFolder(tt.args.folderID, tt.args.includeArchived)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListsForFolder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_ListByID(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		listID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SingleList
		wantErr bool
	}{
		{
			name: "TestSuccessful list by ID",
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
				listID: "fakeListID",
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
				listID: "fakeListID",
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
				listID: "fakeListID",
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
				listID: "fakeListID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:    tt.fields.doer,
				opts:    tt.fields.opts,
				baseURL: tt.fields.baseURL,
			}
			_, err := c.ListByID(tt.args.listID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
