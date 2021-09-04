package clickup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SingleTask struct {
	ID          string `json:"id"`
	CustomID    string `json:"custom_id"`
	Name        string `json:"name"`
	TextContent string `json:"text_content"`
	Description string `json:"description"`
	Status      struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Color      string `json:"color"`
		Orderindex int    `json:"orderindex"`
		Type       string `json:"type"`
	} `json:"status"`
	Orderindex  string      `json:"orderindex"`
	DateCreated string      `json:"date_created"`
	DateUpdated string      `json:"date_updated"`
	DateClosed  interface{} `json:"date_closed"`
	Archived    bool        `json:"archived"`
	Creator     struct {
		ID             int         `json:"id"`
		Username       string      `json:"username"`
		Color          string      `json:"color"`
		Email          string      `json:"email"`
		ProfilePicture interface{} `json:"profilePicture"`
	} `json:"creator"`
	Assignees []struct {
		ID             int    `json:"id"`
		Username       string `json:"username"`
		Color          string `json:"color"`
		Initials       string `json:"initials"`
		Email          string `json:"email"`
		ProfilePicture string `json:"profilePicture"`
	} `json:"assignees"`
	Watchers []struct {
		ID             int         `json:"id"`
		Username       string      `json:"username"`
		Color          string      `json:"color"`
		Initials       string      `json:"initials"`
		Email          string      `json:"email"`
		ProfilePicture interface{} `json:"profilePicture"`
	} `json:"watchers"`
	Checklists []interface{} `json:"checklists"`
	Tags       []interface{} `json:"tags"`
	Parent     string        `json:"parent"`
	Priority   struct {
		ID         string `json:"id"`
		Priority   string `json:"priority"`
		Color      string `json:"color"`
		Orderindex string `json:"orderindex"`
	} `json:"priority"`
	DueDate      interface{} `json:"due_date"`
	StartDate    interface{} `json:"start_date"`
	Points       interface{} `json:"points"`
	TimeEstimate interface{} `json:"time_estimate"`
	TimeSpent    int         `json:"time_spent"`
	CustomFields []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		TypeConfig struct {
			Default     int         `json:"default"`
			Placeholder interface{} `json:"placeholder"`
			Options     []struct {
				ID         string      `json:"id"`
				Name       string      `json:"name"`
				Color      interface{} `json:"color"`
				Orderindex int         `json:"orderindex"`
			} `json:"options"`
		} `json:"type_config"`
		DateCreated    string `json:"date_created"`
		HideFromGuests bool   `json:"hide_from_guests"`
		Required       bool   `json:"required"`
	} `json:"custom_fields"`
	Dependencies    []interface{} `json:"dependencies"`
	LinkedTasks     []interface{} `json:"linked_tasks"`
	TeamID          string        `json:"team_id"`
	URL             string        `json:"url"`
	PermissionLevel string        `json:"permission_level"`
	List            struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Access bool   `json:"access"`
	} `json:"list"`
	Project struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"project"`
	Folder struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"folder"`
	Space struct {
		ID string `json:"id"`
	} `json:"space"`
	Subtasks    []SingleTask  `json:"subtasks"`
	Attachments []interface{} `json:"attachments"`
}

type GetTasksResponse struct {
	Tasks []SingleTask `json:"tasks"`
}

type TaskTimeInStatusResponse struct {
	CurrentStatus struct {
		Status    string `json:"status"`
		Color     string `json:"color"`
		TotalTime struct {
			ByMinute int    `json:"by_minute"`
			Since    string `json:"since"`
		} `json:"total_time"`
	} `json:"current_status"`
	StatusHistory []struct {
		Status    string `json:"status"`
		Color     string `json:"color"`
		Type      string `json:"type"`
		TotalTime struct {
			ByMinute int    `json:"by_minute"`
			Since    string `json:"since"`
		} `json:"total_time"`
		Orderindex int `json:"orderindex"`
	} `json:"status_history"`
}

type OrderByVal string

const (
	OrderByID      OrderByVal = "id"
	OrderByCreated OrderByVal = "created"
	OrderByUpdated OrderByVal = "updated"
	OrderByDueDate OrderByVal = "due_date"
)

type TaskQueryOptions struct {
	IncludeArchived        bool
	Page                   int
	OrderBy                OrderByVal
	Reverse                bool
	IncludeSubtasks        bool
	Statuses               []string // statuses to query
	IncludeClosed          bool
	Assignees              []string
	DueDateGreaterThan     int
	DueDateLessThan        int
	DateCreatedGreaterThan int
	DateCreatedLessThan    int
	DateUpdatedGreaterThan int
	DateUpdatedLessThan    int
	// CustomFields map[string]interface{}
}

func (c *Client) TaskTimeInStatus(taskID, workspaceID string, useCustomTaskIDs bool) (*TaskTimeInStatusResponse, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("task_id", taskID)
	urlValues.Add("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("%s/task/%s/time_in_status/?%s", c.baseURL, taskID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("time in status request failed: %w", err)
	}
	c.AuthenticateFor(req)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make time in status request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var taskTimeInStatus TaskTimeInStatusResponse

	if err := decoder.Decode(&taskTimeInStatus); err != nil {
		return nil, fmt.Errorf("failed to parse time in status: %w", err)
	}

	return &taskTimeInStatus, nil
}

func (c *Client) BulkTaskTimeInStatus(taskIDs []string, workspaceID string, useCustomTaskIDs bool) (map[string]TaskTimeInStatusResponse, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	if len(taskIDs) < 2 || len(taskIDs) > 100 {
		return nil, fmt.Errorf("must provide between 2 and 100 tasks to retrieve bulk tasks time in status: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)
	for _, v := range taskIDs {
		urlValues.Add("task_ids", v)
	}

	endpoint := fmt.Sprintf("%s/task/bulk_time_in_status/task_ids/?%s", c.baseURL, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("bulk time in status request failed: %w", err)
	}
	c.AuthenticateFor(req)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make bulk time in status request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var bulkTaskTimeInStatus map[string]TaskTimeInStatusResponse
	if err := decoder.Decode(&bulkTaskTimeInStatus); err != nil {
		return nil, fmt.Errorf("failed to parse bulk time in status: %w", err)
	}

	return bulkTaskTimeInStatus, nil
}

func (c *Client) TasksForList(listID string, queryOpts TaskQueryOptions) (*GetTasksResponse, error) {

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(queryOpts.IncludeArchived))
	urlValues.Set("subtasks", strconv.FormatBool(queryOpts.IncludeSubtasks))
	urlValues.Set("include_closed", strconv.FormatBool(queryOpts.IncludeClosed))
	// urlValues.Set("date_updated_gt", strconv.Itoa(queryOpts.DateUpdatedGreaterThan))

	endpoint := fmt.Sprintf("%s/list/%s/task/?%s", c.baseURL, listID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("tasks by list request failed: %w", err)
	}
	c.AuthenticateFor(req)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make tasks by list request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var tasks GetTasksResponse

	if err := decoder.Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return &tasks, nil
}

// TODO: need a query options struct to inject because there are so many options for this call
func (c *Client) TaskByID(taskID, workspaceID string, useCustomTaskIDs, includeSubtasks bool) (*SingleTask, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("include_subtasks", strconv.FormatBool(includeSubtasks))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("%s/task/%s/?%s", c.baseURL, taskID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("task by id request failed: %w", err)
	}
	c.AuthenticateFor(req)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make task by id request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var task SingleTask

	if err := decoder.Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)

	}

	return &task, nil
}
