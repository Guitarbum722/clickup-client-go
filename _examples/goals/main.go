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

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Getenv("CLICKUP_API_KEY"),
		},
	})

	goal := clickup.CreateGoalRequest{
		WorkspaceID: os.Getenv("CLICKUP_WORKSPACE_ID"),
		Name:        "Newer Goal",
	}

	res, err := client.CreateGoal(context.Background(), goal)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
