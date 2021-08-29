package main

import (
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {

	client := clickup.NewClient(&clickup.ClientOpts{
		APIToken:   os.Args[1],
		HTTPClient: nil,
	})

	task, err := client.GetTask(os.Args[2], os.Args[3], true, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Task", task.CustomID, task.Name)

	tasks, err := client.GetTasks(os.Args[4], false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, task := range tasks.Tasks {
		fmt.Println("Task: ", task.Name, task.ID)
	}
}
