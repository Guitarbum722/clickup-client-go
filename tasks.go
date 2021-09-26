package clickup

import (
	"bytes"
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

const (
	MaxPageSize = 100
)

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

// TODO: need to figure out paging!
func queryParamsFor(opts *TaskQueryOptions) *url.Values {
	urlValues := &url.Values{}

	urlValues.Add("page", strconv.Itoa(opts.Page))

	if opts.IncludeArchived {
		urlValues.Add("archived", "true")
	}
	if opts.IncludeSubtasks {
		urlValues.Add("subtasks", "true")
	}
	if opts.IncludeClosed {
		urlValues.Add("include_closed", "true")
	}
	if opts.Reverse {
		urlValues.Add("reverse", "true")
	}
	if len(opts.Statuses) > 0 {
		for _, v := range opts.Statuses {
			urlValues.Add("statuses%5B%5D", v)
		}
	}
	if opts.DueDateGreaterThan > 0 {
		urlValues.Add("due_date_gt", strconv.Itoa(opts.DueDateGreaterThan))
	}
	if opts.DueDateLessThan > 0 {
		urlValues.Add("due_date_lt", strconv.Itoa(opts.DueDateLessThan))
	}
	if opts.DateCreatedGreaterThan > 0 {
		urlValues.Add("date_created_gt", strconv.Itoa(opts.DateCreatedGreaterThan))
	}
	if opts.DateCreatedLessThan > 0 {
		urlValues.Add("date_created_lt", strconv.Itoa(opts.DateCreatedLessThan))
	}
	if opts.DateUpdatedGreaterThan > 0 {
		urlValues.Add("date_updated_gt", strconv.Itoa(opts.DateUpdatedGreaterThan))
	}
	if opts.DateUpdatedLessThan > 0 {
		urlValues.Add("date_updated_lt", strconv.Itoa(opts.DateUpdatedLessThan))
	}

	switch opts.OrderBy {
	case OrderByID:
		urlValues.Add("order_by", string(OrderByID))
	case OrderByDueDate:
		urlValues.Add("order_by", string(OrderByDueDate))
	case OrderByCreated:
		urlValues.Add("order_by", string(OrderByCreated))
	case OrderByUpdated:
		urlValues.Add("order_by", string(OrderByUpdated))
	default:

	}
	return urlValues
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
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

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
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

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

// TasksForList returns a listing of tasks that belong to the specified listID and fall withing the constraints of queryOpts.
// Clickup has some rather informal paging, so the caller is responsible for inspecting the count of tasks returned, and incrementing
// the Page in queryOpts if the number of tasks is 100.
// ie. if the current page returns 100 tasks (the maximum page size), then another query should be made to get the next page.
func (c *Client) TasksForList(listID string, queryOpts *TaskQueryOptions) (*GetTasksResponse, error) {

	endpoint := fmt.Sprintf("%s/list/%s/task/?%s", c.baseURL, listID, queryParamsFor(queryOpts).Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("tasks by list request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

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
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

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

type TaskRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Status        string   `json:"status,omitempty"`
	DueDate       int      `json:"due_date,omitempty"`
	DueDateTime   bool     `json:"due_date_time,omitempty"`
	StartDate     int      `json:"start_date,omitempty"`
	StartDateTime bool     `json:"start_date_time,omitempty"`
}

func (c *Client) CreateTask(listID string, task TaskRequest) (*SingleTask, error) {
	if listID == "" {
		return nil, fmt.Errorf("must provide a list id to create a task: %w", ErrValidation)
	}
	if task.Name == "" {
		return nil, fmt.Errorf("must provide a name for a new task: %w", ErrValidation)
	}

	b, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("%s/list/%s/task", c.baseURL, listID)

	req, err := http.NewRequest(http.MethodPost, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create task request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create task request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var newTask SingleTask

	if err := decoder.Decode(&newTask); err != nil {
		return nil, fmt.Errorf("failed to parse new task: %w", err)
	}

	return &newTask, nil
}

type TaskUpdateRequest struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Status        string   `json:"status,omitempty"`
	DueDate       int      `json:"due_date,omitempty"`
	DueDateTime   bool     `json:"due_date_time,omitempty"`
	StartDate     int      `json:"start_date,omitempty"`
	StartDateTime bool     `json:"start_date_time,omitempty"`
}

func (c *Client) UpdateTask(task *TaskUpdateRequest, workspaceID string, useCustomTaskIDs bool) (*SingleTask, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if updating by custom task id: %w", ErrValidation)
	}
	if task.ID == "" {
		return nil, fmt.Errorf("task to update must have an id provided: %w", ErrValidation)
	}

	b, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("%s/task/%s/?%s", c.baseURL, task.ID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodPost, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create task request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make update task request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var updatedTask SingleTask

	if err := decoder.Decode(&updatedTask); err != nil {
		return nil, fmt.Errorf("failed to parse new task: %w", err)
	}

	return &updatedTask, nil
}

func (c *Client) DeleteTask(taskID, workspaceID string, useCustomTaskIDs bool) error {
	if useCustomTaskIDs && workspaceID == "" {
		return fmt.Errorf("workspaceID must be provided if deleting by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("%s/task/%s/?%s", c.baseURL, taskID, urlValues.Encode())
	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete task request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete task request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
}
