package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/Guitarbum722/clickup-client-go"
)

type Config struct {
	WorkspaceID string            `json:"workspace_id"`
	SpaceID     string            `json:"space_id"`
	FolderIDs   map[string]string `json:"folder_ids"`
	ListIDs     map[string]string `json:"list_ids"`
	TaskIDs     []string          `json:"task_ids"`
}

func main() {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	b, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	var config Config

	if err := json.Unmarshal(b, &config); err != nil {
		panic(err)
	}

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Args[1],
		},
	})

	// queriedTasks := map[string]clickup.SingleTask{}

	// for listName, listID := range config.ListIDs {
	// 	tasks, err := client.TasksForList(listID, clickup.TaskQueryOptions{
	// 		IncludeArchived: false,
	// 		IncludeSubtasks: true,
	// 		IncludeClosed:   true,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	for _, task := range tasks.Tasks {
	// 		queriedTasks[task.CustomID] = task
	// 		for _, subtask := range task.Subtasks {
	// 			queriedTasks[subtask.CustomID] = subtask
	// 		}
	// 	}
	// 	fmt.Printf("List ID: %s %s # Tasks: %v\n", listID, listName, len(queriedTasks)) // TODO: REMOVE!
	// }

	// taskIDs := make([]string, 0, len(queriedTasks))
	// for k := range queriedTasks {
	// 	taskIDs = append(taskIDs, k)
	// }

	// taskIDChunks := chunkSlice(taskIDs, 100)

	taskIDs := map[string]struct{}{}
	for _, v := range config.TaskIDs {
		taskIDs[v] = struct{}{}
	}

	taskIDChunks := chunkSlice(config.TaskIDs, 100)

	fmt.Printf("task_id,team_folder,historic_status,status_duration_mins,status_start,status_order,current_status,current_status_since,current_status_duration\n")
	for _, v := range taskIDChunks {
		bulkTimeInStatusResponse, err := client.BulkTaskTimeInStatus(v, config.WorkspaceID, true)
		if err != nil {
			panic(err)
		}

		for taskID, statusHistory := range bulkTimeInStatusResponse {
			currentTimeSince, err := unixMillisToTime(statusHistory.CurrentStatus.TotalTime.Since)
			if err != nil {
				panic(err)
			}
			for _, v := range statusHistory.StatusHistory {
				historyTimeSince, err := unixMillisToTime(v.TotalTime.Since)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s,%s,%s,%d,%s,%d,%s,%s,%d\n",
					taskID,
					// queriedTasks[taskID].Folder.Name,
					taskIDs[taskID],
					v.Status,
					v.TotalTime.ByMinute,
					historyTimeSince,
					v.Orderindex,
					statusHistory.CurrentStatus.Status,
					currentTimeSince,
					statusHistory.CurrentStatus.TotalTime.ByMinute,
				)
			}
		}
	}
}

func unixMillisToTime(m string) (time.Time, error) {
	i, err := strconv.ParseInt(m, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, i*int64(time.Millisecond)), nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
