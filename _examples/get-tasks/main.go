// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Args[1],
		},
	})

	queryOpts := &clickup.TaskQueryOptions{
		IncludeArchived: false,
	}

	for {
		tasks, err := client.TasksForList(context.Background(), os.Args[2], queryOpts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, task := range tasks.Tasks {
			fmt.Println("Task: ", task.CustomID, task.Name)
		}
		if len(tasks.Tasks) < clickup.MaxPageSize {
			return
		} else {
			queryOpts.Page++
		}
	}
}

func msToTime(ms string) time.Time {
	msInt, _ := strconv.ParseInt(ms, 10, 64)

	return time.Unix(0, msInt*int64(time.Millisecond))
}
