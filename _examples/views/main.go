// Copyright (c) 2021, John Moore
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

	view, err := client.ViewByID(context.Background(), "view_id")
	if err != nil {
		panic(err)
	}

	fmt.Println("View info:", view.View.ID, view.View.Name)
}
