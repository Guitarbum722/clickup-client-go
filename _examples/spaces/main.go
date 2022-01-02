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
	// workspaceID := os.Getenv("CLICKUP_WORKSPACE_ID")

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	// spaceReq := clickup.CreateSpaceRequest{
	// 	WorkspaceID:       workspaceID,
	// 	Name:              "First API Space",
	// 	MultipleAssignees: false,
	// }

	// newSpace, err := client.CreateSpaceForWorkspace(context.TODO(), spaceReq)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(newSpace)

	spaceReq := clickup.UpdateSpaceRequest{
		SpaceID: "43682028",
		Name:    "Updated API Space",
		Features: &clickup.Features{
			TimeTracking: &clickup.TimeTracking{
				Enabled: true,
			},
		},
	}

	updatedSpace, err := client.UpdateSpaceForWorkspace(context.TODO(), spaceReq)
	if err != nil {
		panic(err)
	}

	fmt.Println(updatedSpace)
}
