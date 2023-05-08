// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestVerifyWebhookSignature(t *testing.T) {
	type args struct {
		webhookRequest *http.Request
		secret         string
		body           []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *webhookVerifyResult
		wantErr bool
	}{
		{
			name: "Successful validate signature",
			args: args{
				webhookRequest: &http.Request{
					Header: http.Header{
						"X-Signature": []string{"2831500d379c7e90a2c8b3ff55dec81a42889b8a91f6b97f8513d98ebb6b23bf"},
					},
				},
				secret: "imiO3dJZfIlyykAG",
				body:   []byte(`{"event":"taskUpdated"}`),
			},
			want: &webhookVerifyResult{
				validSignature:       true,
				signatureFromClickup: "2831500d379c7e90a2c8b3ff55dec81a42889b8a91f6b97f8513d98ebb6b23bf",
				signatureGenerated:   "2831500d379c7e90a2c8b3ff55dec81a42889b8a91f6b97f8513d98ebb6b23bf",
			},
			wantErr: false,
		},
		{
			name: "Invalid Signature",
			args: args{
				webhookRequest: &http.Request{
					Header: http.Header{
						"X-Signature": []string{"123456"},
					},
				},
				body:   []byte(`{"event":"taskUpdated"}`),
				secret: "imiO3dJZfIlyykAG",
			},
			want: &webhookVerifyResult{
				validSignature:       false,
				signatureFromClickup: "123456",
				signatureGenerated:   "2831500d379c7e90a2c8b3ff55dec81a42889b8a91f6b97f8513d98ebb6b23bf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.webhookRequest.Body = io.NopCloser(bytes.NewReader(tt.args.body))
			got, err := VerifyWebhookSignature(tt.args.webhookRequest, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyWebhookSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyWebhookSignature() = %v, want %v", got, tt.want)
			}

			// also check that the the request body may still be read elsewhere (such as another handler)
			body, err := io.ReadAll(tt.args.webhookRequest.Body)
			if err != nil {
				t.Errorf("unexpected error reading request body: %v", err)
				return
			}

			if !bytes.Equal(tt.args.body, body) {
				t.Errorf("request body after verification = %s, want %s", body, tt.args.body)
			}
		})
	}
}

func TestClient_WebhooksFor(t *testing.T) {
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
			name: "Success workspace id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"webhooks":[{"id":"4b67ac88-e506-4a29-9d42-26e504e3435e","userid":183,"team_id":108,"endpoint":"https://yourdomain.com/webhook","client_id":"QVOQP06ZXC6CMGVFKB0ZT7J9Y7APOYGO","events":["taskCreated"],"task_id":null,"list_id":null,"folder_id":null,"space_id":null,"health":{"status":"failing","fail_count":5},"secret":"O94IM25S7PXBPYTMNXLLET230SRP0S89COR7B1YOJ2ZIE8WQNK5UUKEF26W0Z5GA"}]}`
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
				workspaceID: "test id",
			},
			wantErr: false,
		},
		{
			name: "Fail missing workspace ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
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
			_, err := c.WebhooksFor(tt.args.ctx, tt.args.workspaceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.WebhooksFor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_DeleteWebhook(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx       context.Context
		webhookID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success webhook id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					if req.Method != http.MethodDelete {
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Body:       ioutil.NopCloser(strings.NewReader("")),
							Request:    req,
						}, nil
					}

					body := `{}`
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(strings.NewReader(body)),
						Request:    req,
					}, nil
				}),
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:       context.Background(),
				webhookID: "test id",
			},
			wantErr: false,
		},
		{
			name: "Fail missing webhook ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:       context.Background(),
				webhookID: "",
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
			if err := c.DeleteWebhook(tt.args.ctx, tt.args.webhookID); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UpdateWebhook(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx     context.Context
		webhook *UpdateWebhookRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success webhook id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id": "webhook id", "webhook": {"id": "webhook id", "endpoint": "https://endpoint.clickup"}}`
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
				webhook: &UpdateWebhookRequest{
					ID:       "webhook id",
					Endpoint: "https://endpoint.clickup",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing webhook ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx: context.Background(),
				webhook: &UpdateWebhookRequest{
					ID:       "",
					Endpoint: "https://endpoint.clickup",
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
			_, err := c.UpdateWebhook(tt.args.ctx, tt.args.webhook)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestClient_CreateWebhook(t *testing.T) {
	type fields struct {
		doer          ClientDoer
		authenticator Authenticator
	}
	type args struct {
		ctx         context.Context
		workspaceID string
		webhook     *CreateWebhookRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success workspace id provided",
			fields: fields{
				doer: newMockClientDoer(func(req *http.Request) (*http.Response, error) {
					body := `{"id": "webhook id", "webhook": {"id": "webhook id", "endpoint": "https://endpoint.clickup"}}`
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
				workspaceID: "workspace id",
				webhook: &CreateWebhookRequest{
					Endpoint: "https://endpoint.clickup",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail missing workspace ID",
			fields: fields{
				doer:          nil,
				authenticator: &APITokenAuthenticator{},
			},
			args: args{
				ctx:         context.Background(),
				workspaceID: "",
				webhook: &CreateWebhookRequest{
					Endpoint: "endpoint",
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
			_, err := c.CreateWebhook(tt.args.ctx, tt.args.workspaceID, tt.args.webhook)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
