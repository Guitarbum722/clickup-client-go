// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {
	apiKey := os.Getenv("CLICKUP_API_KEY")
	workspaceID := os.Getenv("CLICKUP_WORKSPACE_ID")

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	groups, err := client.QueryTeams(context.TODO(), workspaceID)
	if err != nil {
		panic(err)
	}
	fmt.Println(groups)
}
