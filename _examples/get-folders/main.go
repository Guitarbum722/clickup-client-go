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
	spaceID := os.Args[2]

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	folders, err := client.FoldersForSpace(context.Background(), spaceID, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, folder := range folders.Folders {
		fmt.Printf("Folder ID: %s\nName: %s\n\n", folder.ID, folder.Name)
	}

	folder, err := client.FolderByID(context.Background(), os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Single Folder\n\nSpace ID: %s\nName: %s\n", folder.ID, folder.Name)
}
