// Copyright (c) 2022, John Moore
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

func TestClient_TagsForSpace(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		spaceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TagsQueryResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: &APITokenAuthenticator{},
				baseURL:       tt.fields.baseURL,
			}
			_, err := c.TagsForSpace(context.Background(), tt.args.spaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.TagsForSpace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateSpaceTag(t *testing.T) {
	type fields struct {
		doer    ClientDoer
		baseURL string
	}
	type args struct {
		spaceID string
		tag     Tag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestSuccessful Tag Create for Space",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				spaceID: "fakeSpaceID",
				tag: Tag{
					Name:  "TestTag",
					TagFg: "#000000",
					TagBg: "#000000",
				},
			},
			wantErr: false,
		},
		{
			name: "TestInvalidTag-Missing tag name",
			fields: fields{
				doer: newMockClientDoer(nil),
			},
			args: args{
				spaceID: "fakeSpaceID",
				tag: Tag{
					TagFg: "#000000",
					TagBg: "#000000",
				},
			},
			wantErr: true,
		},
		{
			name: "TestBadRequest - 400",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{}`
					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				})},
			args: args{
				spaceID: "fakeSpaceID",
				tag:     Tag{Name: "fakeName"},
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
			if err := c.CreateSpaceTag(context.Background(), tt.args.spaceID, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateSpaceTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateSpaceTag(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
		baseURL       string
	}
	type args struct {
		spacID Tag
		tag    Tag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				doer:          tt.fields.doer,
				authenticator: tt.fields.authenticator,
				baseURL:       tt.fields.baseURL,
			}
			if err := c.UpdateSpaceTag(context.Background(), tt.args.spacID, tt.args.tag); (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateSpaceTag() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
