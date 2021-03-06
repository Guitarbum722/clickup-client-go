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
	if os.Args[1] == "" || os.Args[2] == "" {
		panic("missing api key or team id")
	}
	apiKey := os.Args[1]
	teamID := os.Args[2]

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	spaces, err := client.SpacesForWorkspace(context.Background(), teamID, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, space := range spaces.Spaces {
		fmt.Printf("Space ID: %s\nName: %s\n\n", space.ID, space.Name)
	}

	singleSpace, err := client.SpaceByID(context.Background(), "14865529")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Single Space\n\nSpace ID: %s\nName: %s\n", singleSpace.ID, singleSpace.Name)
}
