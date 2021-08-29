package clickup

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_TaskTimeInStatus(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		taskID           string
		workspaceID      string
		useCustomTaskIDs bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful task time in status returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"current_status": {}, "status_history": [{}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
			},
			wantErr: true,
		},
		{
			name: "TestUseCustomTaskIDs No space ID provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) { return nil, nil }),
				opts: &ClientOpts{},
			},
			args: args{
				taskID:           "fakeTaskID",
				useCustomTaskIDs: true,
				workspaceID:      "",
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
			_, err := c.TaskTimeInStatus(tt.args.taskID, tt.args.workspaceID, tt.args.useCustomTaskIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TaskTimeInStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_BulkTaskTimeInStatus(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		taskIDs          []string
		workspaceID      string
		useCustomTaskIDs bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful bulk task time in status returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"fakeTaskID": {"status_history": [{}]}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				taskIDs: []string{"fakeTaskID", "fakeTaskID2"},
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
				taskIDs: []string{"fakeTaskID", "fakeTaskID2"},
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
				taskIDs: []string{"fakeTaskID", "fakeTaskID2"},
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
				taskIDs: []string{"fakeTaskID", "fakeTaskID2"},
			},
			wantErr: true,
		},
		{
			name: "TestUseCustomTaskIDs No space ID provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) { return nil, nil }),
				opts: &ClientOpts{},
			},
			args: args{
				taskIDs:          []string{"fakeTaskID", "fakeTaskID2"},
				useCustomTaskIDs: true,
				workspaceID:      "",
			},
			wantErr: true,
		},
		{
			name: "TestLessThan2IDsProvided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) { return nil, nil }),
				opts: &ClientOpts{},
			},
			args: args{
				taskIDs: []string{"fakeTaskID"},
			},
			wantErr: true,
		},
		{
			name: "TestMoreThan100IDsProvided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) { return nil, nil }),
				opts: &ClientOpts{},
			},
			args: args{
				taskIDs: []string{"fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID", "fakeTaskID"},
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
			_, err := c.BulkTaskTimeInStatus(tt.args.taskIDs, tt.args.workspaceID, tt.args.useCustomTaskIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.BulkTaskTimeInStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_TasksForList(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		listID    string
		queryOpts TaskQueryOptions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful tasks returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"tasks":[{"id":"14865529","name":"John's Happy Taks"}]}`
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
				listID: "fakeListID",
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
			_, err := c.TasksForList(tt.args.listID, tt.args.queryOpts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TasksForList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_TaskByID(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		opts    *ClientOpts
		baseURL string
	}
	type args struct {
		taskID           string
		workspaceID      string
		useCustomTaskIDs bool
		includeSubtasks  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful task returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":"14865529","name":"John's Happy Taks"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				opts: &ClientOpts{},
			},
			args: args{
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
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
				taskID: "fakeTaskID",
			},
			wantErr: true,
		},
		{
			name: "TestUseCustomTaskIDs No space ID provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) { return nil, nil }),
				opts: &ClientOpts{},
			},
			args: args{
				taskID:           "fakeTaskID",
				useCustomTaskIDs: true,
				workspaceID:      "",
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
			_, err := c.TaskByID(tt.args.taskID, tt.args.workspaceID, tt.args.useCustomTaskIDs, tt.args.includeSubtasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TaskByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
