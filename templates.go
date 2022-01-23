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

type Template struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TemplatesResponse struct {
	Templates []Template `json:"templates"`
}

// TemplatesForWorkspace returns all templates for the workspace ID provided that the authenticated user
// is authorized to see.
// Specify the page number with page, starting with 0.
// TODO: ClickUp does not document what the max page is and there is no straightforward way to know.
// At this time, the page parameter doesn't seem to do anything.
func (c *Client) TemplatesForWorkspace(ctx context.Context, workspaceID string, page int) (*TemplatesResponse, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id to query templates: %w", ErrValidation)
	}
	urlValues := url.Values{}
	urlValues.Set("page", strconv.Itoa(page))

	endpoint := fmt.Sprintf("/team/%s/taskTemplate/?%s", workspaceID, urlValues.Encode())

	var templates TemplatesResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &templates); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &templates, nil
}

type TaskFromTemplateRequest struct {
	ListID     string
	TemplateID string
	Name       string `json:"name"`
}

type TaskFromTemplateResponse struct {
	ID string `json:"id"`
}

// CreateTaskFromTemplate creates a new task based on the template id specified in the TaskFromTemplateRequest.
// The new task will be created in the list specified by ListID in the TaskFromTemplateRequest.
func (c *Client) CreateTaskFromTemplate(ctx context.Context, task TaskFromTemplateRequest) (*TaskFromTemplateResponse, error) {
	if task.ListID == "" {
		return nil, fmt.Errorf("must provide a list id to create a task from template: %w", ErrValidation)
	}
	if task.TemplateID == "" {
		return nil, fmt.Errorf("must provide a template id to create a new task from template: %w", ErrValidation)
	}

	b, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/list/%s/taskTemplate/%s", task.ListID, task.TemplateID)

	var newTask TaskFromTemplateResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newTask); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &newTask, nil
}
