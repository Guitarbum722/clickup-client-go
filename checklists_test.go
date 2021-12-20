package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_CreateChecklist(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		request *CreateChecklistRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful create checklist",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"checklist":{"id":"test-id","name": "test name", "task_id": "test-task-id"}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				request: &CreateChecklistRequest{
					TaskID:           "test-task-id",
					WorkspaceID:      "workspace-id",
					UseCustomTaskIDs: false,
					Name:             "test name",
				},
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing task id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &CreateChecklistRequest{
					Name:        "checklist item name",
					TaskID:      "",
					WorkspaceID: "test-workspace-id",
				},
			},
			wantErr: true,
		},
		{
			name: "TestFailure missing workspace id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &CreateChecklistRequest{
					Name:        "checklist item name",
					TaskID:      "test-task-id",
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
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.CreateChecklist(context.Background(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateChecklist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateChecklist(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		request *UpdateChecklistRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChecklistResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful update checklist",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"checklist":{"id":"test-id","name": "test name", "task_id": "test-task-id"}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				request: &UpdateChecklistRequest{
					ChecklistID: "test-checklist-id",
					Position:    0,
					Name:        "test name",
				},
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing checklist id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &UpdateChecklistRequest{
					ChecklistID: "",
				},
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
			_, err := c.UpdateChecklist(context.Background(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateChecklist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteChecklist(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		checklistID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful delete checklist",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				checklistID: "test-checklist-id",
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing checklist id",
			fields: fields{
				doer: nil,
			},
			args: args{
				checklistID: "",
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
			if err := c.DeleteChecklist(context.Background(), tt.args.checklistID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteChecklist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CreateChecklistItem(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		request *CreateChecklistItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChecklistResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful create checklist item",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"checklist":{"id":"test-id","name": "test name", "task_id": "test-task-id"}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				request: &CreateChecklistItemRequest{
					ChecklistID: "test-checklist-id",
					Name:        "checklist item name",
				},
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing checklist id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &CreateChecklistItemRequest{
					ChecklistID: "",
				},
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
			_, err := c.CreateChecklistItem(context.Background(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateChecklistItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateChecklistItem(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		ctx     context.Context
		request *UpdateChecklistItemRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ChecklistResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful update checklist item",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"checklist":{"id":"test-id","name": "test name", "task_id": "test-task-id"}}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				request: &UpdateChecklistItemRequest{
					ChecklistID:     "test-checklist-id",
					ChecklistItemID: "test-checklist-item-id",
					Resolved:        false,
					Name:            "test name",
				},
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing checklist id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &UpdateChecklistItemRequest{
					ChecklistID:     "",
					ChecklistItemID: "checklist-item-id",
				},
			},
			wantErr: true,
		},
		{
			name: "TestFailure missing checklist item id",
			fields: fields{
				doer: nil,
			},
			args: args{
				request: &UpdateChecklistItemRequest{
					ChecklistID:     "checklist-id",
					ChecklistItemID: "",
				},
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
			_, err := c.UpdateChecklistItem(context.Background(), tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateChecklistItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteChecklistItem(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx             context.Context
		checklistID     string
		checklistItemID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful delete checklist item",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				checklistID:     "test-checklist-id",
				checklistItemID: "test-checklist-item-id",
			},
			wantErr: false,
		},
		{
			name: "TestFailure missing checklist id",
			fields: fields{
				doer: nil,
			},
			args: args{
				checklistID:     "",
				checklistItemID: "checklist-item-id",
			},
			wantErr: true,
		},
		{
			name: "TestFailure missing checklist item id",
			fields: fields{
				doer: nil,
			},
			args: args{
				checklistID:     "checklist-id",
				checklistItemID: "",
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
			if err := c.DeleteChecklistItem(context.Background(), tt.args.checklistID, tt.args.checklistItemID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteChecklistItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
