package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/Guitarbum722/clickup-client-go"
)

type Config struct {
	APIToken    string            `json:"api_token"`
	WorkspaceID string            `json:"workspace_id"`
	SpaceID     string            `json:"space_id"`
	FolderIDs   map[string]string `json:"folder_ids"`
	ListIDs     map[string]string `json:"list_ids"`
	TaskIDs     []string          `json:"task_ids"`
}

const maxBulkStatusRecords = 100

const usage = `clickup-time-status [-h] [-f]
Options:
 -h         help
 -f         configuration file path
`

func main() {
	var configFilePath string
	var help bool

	flag.StringVar(&configFilePath, "f", "", "")
	flag.BoolVar(&help, "h", false, "")
	flag.Parse()

	if help {
		fmt.Println(usage)
		os.Exit(0)
	}

	if configFilePath == "" {
		myDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configFilePath = fmt.Sprintf("%s/%s", myDir, "config.json")
	}

	configFile, err := os.Open(configFilePath)
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
			APIToken: config.APIToken,
		},
	})

	taskIDChunks := chunkSlice(config.TaskIDs, maxBulkStatusRecords)

	fmt.Printf("task_id,historic_status,status_start,status_end,status_weekdays_duration,status_order,current_status,current_status_start,current_status_end,current_status_weekdays_duration\n")
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
			currentStatusEnd := currentTimeSince.Add(time.Minute * time.Duration(statusHistory.CurrentStatus.TotalTime.ByMinute))

			for _, v := range statusHistory.StatusHistory {
				historyTimeSince, err := unixMillisToTime(v.TotalTime.Since)
				if err != nil {
					panic(err)
				}

				historyStatusEnd := historyTimeSince.Add(time.Minute * time.Duration(v.TotalTime.ByMinute))

				fmt.Printf("%s,%s,%s,%s,%v,%d,%s,%s,%s,%v\n",
					taskID,
					v.Status,
					historyTimeSince,
					historyStatusEnd,
					businessDaysBetweenTimes(historyTimeSince, historyStatusEnd),
					v.Orderindex,
					statusHistory.CurrentStatus.Status,
					currentTimeSince,
					currentStatusEnd,
					businessDaysBetweenTimes(currentTimeSince, currentStatusEnd),
				)
			}
		}
	}
}

func businessDaysBetweenTimes(start, end time.Time) int {
	offset := -int(start.Weekday())
	start = start.AddDate(0, 0, -int(start.Weekday()))

	offset += int(end.Weekday())
	if end.Weekday() == time.Sunday {
		offset++
	}
	end = end.AddDate(0, 0, -int(end.Weekday()))

	dif := end.Sub(start).Truncate(time.Hour * 24)
	weeks := float64((dif.Hours() / 24) / 7)
	return int(math.Round(weeks)*5) + offset
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
