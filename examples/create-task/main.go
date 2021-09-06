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
