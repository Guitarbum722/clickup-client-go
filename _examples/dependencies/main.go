// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"context"
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

	dependency := clickup.AddDependencyRequest{
		TaskID:           os.Getenv("TASK_ID"),
		UseCustomTaskIDs: true,
		WorkspaceID:      os.Getenv("CLICKUP_WORKSPACE_ID"),
		DependencyOf:     os.Getenv("DEPENDENCY_OF_TASK_ID"),
	}

	err := client.AddDependencyForTask(context.Background(), dependency)
	if err != nil {
		panic(err)
	}
}
