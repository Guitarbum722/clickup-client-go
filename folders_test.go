// Copyright (c) 2021, John Moore
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

func TestClient_FoldersForSpace(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		spaceID         string
		includeArchived bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *FoldersResponse
		wantErr bool
	}{
		{
			name: "TestSuccessful folders returned",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"folders":[{"id":"12345678","name":"Hippy Dippy Folder"}]}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
			},
			args: args{
				spaceID: "fakeSpaceID",
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
				spaceID: "fakeSpaceID",
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
				spaceID: "fakeSpaceID",
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
			},
			args: args{
				spaceID: "fakeSpaceID",
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
			_, err := c.FoldersForSpace(context.Background(), tt.args.spaceID, tt.args.includeArchived)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FoldersForSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_FolderByID(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		folderID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *SingleFolder
		wantErr bool
	}{
		{
			name: "TestSuccessful folder by ID",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id":"123","name":"John's Happy Place"}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
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
			},
			args: args{
				folderID: "fakeFolderID",
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
			},
			args: args{
				folderID: "fakeFolderID",
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
			},
			args: args{
				folderID: "fakeFolderID",
			},
			wantErr: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.FolderByID(context.Background(), tt.args.folderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FolderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
