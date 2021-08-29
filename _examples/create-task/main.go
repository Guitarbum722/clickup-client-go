// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Args[1],
		},
	})

	task := clickup.TaskRequest{
		Name: "Awesome",
	}

	newTask, err := client.CreateTask(os.Args[2], task)
	if err != nil {
		panic(err)
	}

	fmt.Println("New Task:", newTask.ID, newTask.Name)
}
