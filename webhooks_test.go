// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestVerifyWebhookSignature(t *testing.T) {
	type args struct {
		webhookRequest *http.Request
		secret         string
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
					Body: io.NopCloser(strings.NewReader(`{"event":"taskUpdated"}`)),
				},
				secret: "imiO3dJZfIlyykAG",
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
					Body: io.NopCloser(strings.NewReader(`{"event":"taskUpdated"}`)),
				},
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
			got, err := VerifyWebhookSignature(tt.args.webhookRequest, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyWebhookSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyWebhookSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
