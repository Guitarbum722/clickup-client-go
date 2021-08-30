package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

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
	fmt.Println(task.DateUpdated)
	fmt.Printf("Task:  %s %s %s\n", task.CustomID, task.Name, msToTime(task.DateUpdated))

	tasks, err := client.GetTasks(os.Args[4], false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, task := range tasks.Tasks {
		fmt.Println("Task: ", task.ID, task.Name)
	}
}

func msToTime(ms string) time.Time {
	msInt, _ := strconv.ParseInt(ms, 10, 64)

	return time.Unix(0, msInt*int64(time.Millisecond))
}
