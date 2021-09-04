package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {
	if os.Args[1] == "" || os.Args[2] == "" || os.Args[3] == "" {
		panic("missing api key | task id | workspace id")
	}
	apiKey := os.Args[1]
	taskID := os.Args[2]
	workspaceID := os.Args[3]

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	timeInStatusResponse, err := client.TaskTimeInStatus(taskID, workspaceID, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Task:", taskID)
	fmt.Printf("Current Status: %s\n", timeInStatusResponse.CurrentStatus.Status)
	fmt.Printf("Current Since: %s\n", msToTime(timeInStatusResponse.CurrentStatus.TotalTime.Since))
	fmt.Printf("Current Duration: %v days\n", minsToDays(timeInStatusResponse.CurrentStatus.TotalTime.ByMinute))
	fmt.Println("-------------")

	for _, v := range timeInStatusResponse.StatusHistory {
		fmt.Printf("Status: %-15s\t%s\t%v days\tOrder: %v\n", v.Status, msToTime(v.TotalTime.Since), minsToDays(v.TotalTime.ByMinute), v.Orderindex)
	}
}

func msToTime(ms string) time.Time {
	msInt, _ := strconv.ParseInt(ms, 10, 64)

	return time.Unix(0, msInt*int64(time.Millisecond))
}

func minsToDays(mins int) int {
	return mins / 60 / 24
}
