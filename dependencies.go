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

type AddDependencyRequest struct {
	TaskID           string
	DependsOn        string `json:"depends_on,omitempty"`
	DependencyOf     string `json:"dependency_of,omitempty"`
	WorkspaceID      string
	UseCustomTaskIDs bool
}

func (c *Client) AddDependencyForTask(ctx context.Context, dependency AddDependencyRequest) error {
	if dependency.TaskID == "" {
		return fmt.Errorf("must provide a task id to create a dependency: %w", ErrValidation)
	}
	if dependency.WorkspaceID == "" {
		return fmt.Errorf("must provide a workspace id to create a dependency: %w", ErrValidation)
	}
	if len(dependency.DependsOn) > 0 && len(dependency.DependencyOf) > 0 {
		return fmt.Errorf("must provide either a depends_on or dependency_of to create a dependency but not both: %w", ErrValidation)
	}
	if dependency.DependsOn == "" && dependency.DependencyOf == "" {
		return fmt.Errorf("must provide either a depends_on or dependency_of to create a dependency: %w", ErrValidation)
	}

	b, err := json.Marshal(dependency)
	if err != nil {
		return fmt.Errorf("unable to serialize new dependency: %w", err)
	}
	buf := bytes.NewBuffer(b)

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(dependency.UseCustomTaskIDs))
	urlValues.Add("team_id", dependency.WorkspaceID)

	endpoint := fmt.Sprintf("/task/%v/dependency/?%s", dependency.TaskID, urlValues.Encode())

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &struct{}{}); err != nil {
		return fmt.Errorf("failed to make clickup request: %w", err)
	}

	return nil
}

type AddTaskLinkRequest struct {
	TaskID           string
	LinksToTaskID    string
	WorkspaceID      string
	UseCustomTaskIDs bool
}

type TaskLinkResponse struct {
	Task *SingleTask `json:"task"`
}

func (c *Client) AddTaskLinkForTask(ctx context.Context, link AddTaskLinkRequest) (*TaskLinkResponse, error) {
	if link.TaskID == "" {
		return nil, fmt.Errorf("must provide a task id to create a task link: %w", ErrValidation)
	}
	if link.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id to create a task link: %w", ErrValidation)
	}
	if link.LinksToTaskID == "" {
		return nil, fmt.Errorf("must provide a task id to link to: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(link.UseCustomTaskIDs))
	urlValues.Add("team_id", link.WorkspaceID)

	endpoint := fmt.Sprintf("/task/%v/link/%v/?%v", link.TaskID, link.LinksToTaskID, urlValues.Encode())

	var linkedTask TaskLinkResponse

	if err := c.call(ctx, http.MethodPost, endpoint, &bytes.Buffer{}, &linkedTask); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &linkedTask, nil
}
