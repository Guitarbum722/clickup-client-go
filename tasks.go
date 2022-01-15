// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Status struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Color      string `json:"color"`
	Orderindex int    `json:"orderindex"`
	Type       string `json:"type"`
}

type SingleTask struct {
	ID          string     `json:"id"`
	CustomID    string     `json:"custom_id"`
	Name        string     `json:"name"`
	TextContent string     `json:"text_content"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Orderindex  string     `json:"orderindex"`
	DateCreated string     `json:"date_created"`
	DateUpdated string     `json:"date_updated"`
	DateClosed  string     `json:"date_closed"`
	Archived    bool       `json:"archived"`
	Creator     TeamUser   `json:"creator"`
	Assignees   []TeamUser `json:"assignees"`
	Watchers    []TeamUser `json:"watchers"`
	Checklists  []struct {
		ID          string `json:"id"`
		TaskID      string `json:"task_id"`
		Name        string `json:"name"`
		DateCreated string `json:"date_created"`
		Orderindex  int    `json:"orderindex"`
		Creator     int    `json:"creator"`
		Resolved    int    `json:"resolved"`
		Unresolved  int    `json:"unresolved"`
		Items       []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Orderindex int    `json:"orderindex"`
			Assignee   struct {
				ID             int    `json:"id"`
				Username       string `json:"username"`
				Email          string `json:"email"`
				Color          string `json:"color"`
				Initials       string `json:"initials"`
				ProfilePicture string `json:"profilePicture"`
			} `json:"assignee"`
			Resolved    bool   `json:"resolved"`
			DateCreated string `json:"date_created"`
		} `json:"items"`
	} `json:"checklists"`
	Tags     []Tag  `json:"tags"`
	Parent   string `json:"parent"`
	Priority struct {
		ID         string `json:"id"`
		Priority   string `json:"priority"`
		Color      string `json:"color"`
		Orderindex string `json:"orderindex"`
	} `json:"priority"`
	DueDate      string `json:"due_date"`
	StartDate    string `json:"start_date"`
	Points       int    `json:"points"`
	TimeEstimate int    `json:"time_estimate"`
	TimeSpent    int    `json:"time_spent"`
	CustomFields []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Type       string `json:"type"`
		TypeConfig struct {
			Default     int    `json:"default"`
			Placeholder string `json:"placeholder"`
			Options     []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Color      string `json:"color"`
				Orderindex int    `json:"orderindex"`
			} `json:"options"`
		} `json:"type_config"`
		DateCreated    string      `json:"date_created"`
		HideFromGuests bool        `json:"hide_from_guests"`
		Required       bool        `json:"required"`
		Value          interface{} `json:"value"`
	} `json:"custom_fields"`
	Dependencies []struct {
		TaskID      string `json:"task_id"`
		DependsOn   string `json:"depends_on"`
		Type        int    `json:"type"`
		DateCreated string `json:"date_created"`
		Userid      string `json:"userid"`
	} `json:"dependencies"`
	LinkedTasks []struct {
		TaskID      string `json:"task_id"`
		LinkID      string `json:"link_id"`
		DateCreated string `json:"date_created"`
		Userid      string `json:"userid"`
	} `json:"linked_tasks"`
	TeamID          string `json:"team_id"`
	URL             string `json:"url"`
	PermissionLevel string `json:"permission_level"`
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
	Subtasks    []SingleTask `json:"subtasks"`
	Attachments []struct{}   `json:"attachments"`
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

// TaskTimeInStatus returns status history for taskID.  useCustomTaskIDs should be true if querying with a custom ID.
func (c *Client) TaskTimeInStatus(ctx context.Context, taskID, workspaceID string, useCustomTaskIDs bool) (*TaskTimeInStatusResponse, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("task_id", taskID)
	urlValues.Add("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("/task/%s/time_in_status/?%s", taskID, urlValues.Encode())

	var taskTimeInStatus TaskTimeInStatusResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &taskTimeInStatus); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &taskTimeInStatus, nil
}

// BulkTaskTimeInStatus returns task status history data for the provided taskIDs.  Must provide
// >= 2 and <= 100 task IDs at a time.
func (c *Client) BulkTaskTimeInStatus(ctx context.Context, taskIDs []string, workspaceID string, useCustomTaskIDs bool) (map[string]TaskTimeInStatusResponse, error) {
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

	endpoint := fmt.Sprintf("/task/bulk_time_in_status/task_ids/?%s", urlValues.Encode())

	var bulkTaskTimeInStatus map[string]TaskTimeInStatusResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &bulkTaskTimeInStatus); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return bulkTaskTimeInStatus, nil
}

// TasksForList returns a listing of tasks that belong to the specified listID and fall withing the constraints of queryOpts.
// Clickup has some rather informal paging, so the caller is responsible for inspecting the count of tasks returned, and incrementing
// the Page in queryOpts if the number of tasks is 100.
// ie. if the current page returns 100 tasks (the maximum page size), then another query should be made to get the next page.
func (c *Client) TasksForList(ctx context.Context, listID string, queryOpts *TaskQueryOptions) (*GetTasksResponse, error) {

	endpoint := fmt.Sprintf("/list/%s/task/?%s", listID, queryParamsFor(queryOpts).Encode())

	var tasks GetTasksResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &tasks); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &tasks, nil
}

// TaskByID queries a single task.
func (c *Client) TaskByID(ctx context.Context, taskID, workspaceID string, useCustomTaskIDs, includeSubtasks bool) (*SingleTask, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("include_subtasks", strconv.FormatBool(includeSubtasks))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("/task/%s/?%s", taskID, urlValues.Encode())

	var task SingleTask

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &task); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
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

// CreateTask inserts a new task into the specified list.
func (c *Client) CreateTask(ctx context.Context, listID string, task TaskRequest) (*SingleTask, error) {
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

	endpoint := fmt.Sprintf("/list/%s/task", listID)

	var newTask SingleTask

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newTask); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
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

// UpdateTask changes an existing task.
func (c *Client) UpdateTask(ctx context.Context, task *TaskUpdateRequest, workspaceID string, useCustomTaskIDs bool) (*SingleTask, error) {
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

	endpoint := fmt.Sprintf("/task/%s/?%s", task.ID, urlValues.Encode())

	var updatedTask SingleTask

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedTask); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &updatedTask, nil
}

// DeleteTask removes an existing task.
func (c *Client) DeleteTask(ctx context.Context, taskID, workspaceID string, useCustomTaskIDs bool) error {
	if useCustomTaskIDs && workspaceID == "" {
		return fmt.Errorf("workspaceID must be provided if deleting by custom task id: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("/task/%s/?%s", taskID, urlValues.Encode())

	return c.call(ctx, http.MethodDelete, endpoint, nil, &struct{}{})
}
