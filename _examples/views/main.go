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
	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	view, err := client.ViewByID(context.Background(), os.Getenv("CLICKUP_VIEW_ID"))
	if err != nil {
		panic(err)
	}

	fmt.Println("View info:", view.View.ID, view.View.Name)

	views, err := client.ViewsFor(context.Background(), clickup.TypeSpace, os.Getenv("CLICKUP_SPACE_ID"))
	if err != nil {
		panic(err)
	}

	for _, v := range views.Views {
		fmt.Println(v.Name)
	}

	tasks, err := client.TasksForView(context.Background(), os.Getenv("CLICKUP_VIEW_ID"), 0)
	if err != nil {
		panic(err)
	}

	for _, v := range tasks.Tasks {
		fmt.Println(v.Name)
	}
}
