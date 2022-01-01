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

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Getenv("CLICKUP_API_KEY"),
		},
	})

	comment := clickup.NewCreateTaskCommentRequest(
		os.Getenv("CLICKUP_TASK_ID"),
		true,
		os.Getenv("CLICKUP_WORKSPACE_ID"),
	)
	comment.BulletedListItem("Bullet Item 4", nil)
	comment.BulletedListItem("Bullet Item 5", nil)
	comment.BulletedListItem("Bullet Item 6", &clickup.Attributes{Italic: true})
	comment.NumberedListItem("Numbered Item 1", nil)
	comment.ChecklistItem("Checklist item 1", false, nil)
	comment.ChecklistItem("Checklist item 2", true, nil)

	res, err := client.CreateTaskComment(context.Background(), *comment)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
