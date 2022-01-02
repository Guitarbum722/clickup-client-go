// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClient_CreateTaskAttachment(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		taskID           string
		workspaceID      string
		useCustomTaskIDs bool
		params           *AttachmentParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful create attachment",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id": "test-attachment-id", "date": 1631484506345}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				taskID:           "test-task-id",
				workspaceID:      "test-workspace-id",
				useCustomTaskIDs: true,
				params: &AttachmentParams{
					FileName: "Test.txt",
					Reader:   bytes.NewBufferString("This is testing data."),
				},
			},
			wantErr: false,
		},
		{
			name: "TestFailure no workspace id",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					return nil, nil
				}),
			},
			args: args{
				taskID:           "test-task-id",
				useCustomTaskIDs: true,
				params: &AttachmentParams{
					FileName: "Test.txt",
					Reader:   bytes.NewBufferString("This is testing data."),
				},
			},
			wantErr: true,
		},
		{
			name: "TestFailure http bad request",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"ECODE": "ATTACH_044", "Err": "File must be named 'attachment'"}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				taskID: "test-task-id",
				params: &AttachmentParams{
					FileName: "Test.txt",
					Reader:   bytes.NewBufferString("This is testing data."),
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
			_, err := c.CreateTaskAttachment(context.Background(), tt.args.taskID, tt.args.workspaceID, tt.args.useCustomTaskIDs, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateTaskAttachment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
