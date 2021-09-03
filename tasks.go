package clickup

import (
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
	var taskTimeInStatus TaskTimeInStatusResponse

	urlValues := url.Values{}
	urlValues.Set("task_id", taskID)
	urlValues.Add("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	uri := fmt.Sprintf("/task/%s/time_in_status/?%s", taskID, urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &taskTimeInStatus); err != nil {
		return nil, err
	}

	return &taskTimeInStatus, nil
}

func (c *Client) BulkTaskTimeInStatus(taskIDs []string, workspaceID string, useCustomTaskIDs bool) (map[string]TaskTimeInStatusResponse, error) {
	if len(taskIDs) < 2 || len(taskIDs) > 100 {
		return nil, fmt.Errorf("must provide between 2 and 100 tasks to retrieve bulk tasks time in status: %w", ErrValidation)
	}

	var bulkTaskTimeInStatus map[string]TaskTimeInStatusResponse

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)
	for _, v := range taskIDs {
		urlValues.Add("task_ids", v)
	}

	uri := fmt.Sprintf("/task/bulk_time_in_status/task_ids/?%s", urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &bulkTaskTimeInStatus); err != nil {
		return nil, err
	}

	return bulkTaskTimeInStatus, nil
}

// TODO: need a query options struct to inject because there are so many options for this call
func (c *Client) TasksForList(listID string, queryOpts TaskQueryOptions) (*GetTasksResponse, error) {
	var tasks GetTasksResponse

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(queryOpts.IncludeArchived))
	urlValues.Set("subtasks", strconv.FormatBool(queryOpts.IncludeSubtasks))
	urlValues.Set("include_closed", strconv.FormatBool(queryOpts.IncludeClosed))
	// urlValues.Set("date_updated_gt", strconv.Itoa(queryOpts.DateUpdatedGreaterThan))

	uri := fmt.Sprintf("/list/%s/task/?%s", listID, urlValues.Encode())
	fmt.Println("uri: ", uri)
	if err := c.call(http.MethodGet, uri, nil, &tasks); err != nil {
		return nil, err
	}

	return &tasks, nil
}

// TODO: need a query options struct to inject because there are so many options for this call
func (c *Client) TaskByID(taskID, workspaceID string, useCustomTaskIDs, includeSubtasks bool) (*SingleTask, error) {
	var task SingleTask

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("include_subtasks", strconv.FormatBool(includeSubtasks))
	urlValues.Add("team_id", workspaceID)

	uri := fmt.Sprintf("/task/%s/?%s", taskID, urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &task); err != nil {
		return nil, err
	}

	return &task, nil
}
