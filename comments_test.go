// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestCreateCommentRequest_BulletedListItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build bulleted list starting with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build bulleted list starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build bulleted list with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.BulletedListItem(tt.args.text, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.BulletedListItem() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCreateCommentRequest_NumberedListItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build numbered list starting with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build ordered list starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build ordered list with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.NumberedListItem(tt.args.text, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.NumberedListItem() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCreateCommentRequest_ChecklistItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		checked    bool
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build checklist item with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    true,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "checked",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build checklist item starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    false,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build checklist item with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "checklist item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    false,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.ChecklistItem(tt.args.text, tt.args.checked, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.ChecklistItem() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestClient_CreateTaskComment(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx     context.Context
		comment CreateTaskCommentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success validated inputs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":458,"hist_id":"26508","date":1568036964079}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx: context.Background(),
				comment: CreateTaskCommentRequest{
					TaskID:           "123",
					UseCustomTaskIDs: false,
					WorkspaceID:      "444",
					CreateCommentRequest: CreateCommentRequest{
						CommentText: "test data",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing task ID",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":458,"hist_id":"26508","date":1568036964079}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx: context.Background(),
				comment: CreateTaskCommentRequest{
					TaskID:           "",
					UseCustomTaskIDs: false,
					WorkspaceID:      "444",
					CreateCommentRequest: CreateCommentRequest{
						CommentText: "test data",
					},
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
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.CreateTaskComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateTaskComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateListComment(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx     context.Context
		comment CreateListCommentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success validated inputs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":458,"hist_id":"26508","date":1568036964079}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx: context.Background(),
				comment: CreateListCommentRequest{
					ListID: "123",
					CreateCommentRequest: CreateCommentRequest{
						CommentText: "test data",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing list ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				comment: CreateListCommentRequest{
					ListID:               "",
					CreateCommentRequest: CreateCommentRequest{},
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
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.CreateListComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateListComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateChatViewComment(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx     context.Context
		comment CreateChatViewCommentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success validated inputs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":458,"hist_id":"26508","date":1568036964079}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx: context.Background(),
				comment: CreateChatViewCommentRequest{
					ViewID: "123",
					CreateCommentRequest: CreateCommentRequest{
						CommentText: "test data",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing view ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				comment: CreateChatViewCommentRequest{
					ViewID:               "",
					CreateCommentRequest: CreateCommentRequest{},
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
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.CreateChatViewComment(tt.args.ctx, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateChatViewComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_UpdateComment(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx     context.Context
		comment UpdateCommentRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success validated inputs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx: context.Background(),
				comment: UpdateCommentRequest{
					CommentID:   "123",
					CommentText: "test comment",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing comment ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				comment: UpdateCommentRequest{
					CommentID: "",
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
				baseURL:       tt.fields.baseURL,
			}
			if err := c.UpdateComment(tt.args.ctx, tt.args.comment); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_DeleteComment(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		ctx       context.Context
		commentID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success validated inputs provided",
			fields: fields{
				authenticator: &APITokenAuthenticator{},
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				ctx:       context.Background(),
				commentID: "123",
			},
			wantErr: false,
		},
		{
			name: "Fail missing comment ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "",
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
			if err := c.DeleteComment(tt.args.ctx, tt.args.commentID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
