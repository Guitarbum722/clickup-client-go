// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func Test_errorFromResponse(t *testing.T) {
	type args struct {
		res     *http.Response
		decoder *json.Decoder
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantType string
	}{
		{
			name: "Return rate limit error",
			args: args{
				res: &http.Response{
					StatusCode: http.StatusTooManyRequests,
					Header: http.Header{
						"x-ratelimit-limit":     []string{"10000"},
						"x-ratelimit-remaining": []string{"9999"},
						"x-ratelimit-reset":     []string{"1640818767"},
					},
				},
				decoder: json.NewDecoder(strings.NewReader("")),
			},
			wantErr:  true,
			wantType: "*clickup.RateLimitError",
		},
		{
			name: "Return rate clickup clickup error",
			args: args{
				res: &http.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Request: &http.Request{
						URL: &url.URL{Host: "https://mock.clickup.com"},
					},
				},
				decoder: json.NewDecoder(strings.NewReader(`{"ECODE": "fail", "err": "error"}`)),
			},
			wantErr:  true,
			wantType: "*clickup.ErrClickupResponse",
		},
		{
			name: "Return rate clickup clickup http error",
			args: args{
				res: &http.Response{
					StatusCode: http.StatusUnprocessableEntity,
					Request: &http.Request{
						URL: &url.URL{Host: "https://mock.clickup.com"},
					},
				},
				decoder: json.NewDecoder(strings.NewReader(`{{{badJSON}}}}`)),
			},
			wantErr:  true,
			wantType: "*clickup.HTTPError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorFromResponse(tt.args.res, tt.args.decoder)
			if (err != nil) != tt.wantErr {
				t.Errorf("errorFromResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if fmt.Sprintf("%T", err) != tt.wantType {
				t.Errorf("errorFromResponse() error = %T, wantType %v", err, tt.wantType)
			}
		})
	}
}
