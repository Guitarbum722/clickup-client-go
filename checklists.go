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

type ChecklistResponse struct {
	Checklist struct {
		ID          string `json:"id"`
		TaskID      string `json:"task_id"`
		Name        string `json:"name"`
		DateCreated string `json:"date_created"`
		Orderindex  int    `json:"orderindex"`
		Creator     int    `json:"creator"`
		Resolved    int    `json:"resolved"`
		Unresolved  int    `json:"unresolved"`
		Items       []struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			Orderindex  int      `json:"orderindex"`
			Assignee    TeamUser `json:"assignee"`
			Resolved    bool     `json:"resolved"`
			DateCreated string   `json:"date_created"`
		} `json:"items"`
	} `json:"checklist"`
}

type CreateChecklistRequest struct {
	TaskID           string `json:"-"`
	WorkspaceID      string `json:"-"`
	UseCustomTaskIDs bool   `json:"-"`
	Name             string `json:"name"`
}

type UpdateChecklistRequest struct {
	ChecklistID string `json:"-"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
}

type CreateChecklistItemRequest struct {
	ChecklistID string `json:"-"`
	Name        string `json:"name"`
}

type UpdateChecklistItemRequest struct {
	ChecklistID     string     `json:"checklist_id"`
	ChecklistItemID string     `json:"checklist_item_id"`
	Name            string     `json:"name"`
	Assignee        TeamMember `json:"assignee"`
	Resolved        bool       `json:"resolved"`
}

func (c *Client) CreateChecklist(ctx context.Context, request *CreateChecklistRequest) (*ChecklistResponse, error) {
	if request.TaskID == "" {
		return nil, fmt.Errorf("must provide task id: %w", ErrValidation)
	}
	if request.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide workspace id: %w", ErrValidation)
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize checklist request: %w", err)
	}
	buf := bytes.NewBuffer(b)

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(request.UseCustomTaskIDs))
	urlValues.Add("team_id", request.WorkspaceID)

	endpoint := fmt.Sprintf("/task/%s/checklist/?%s", request.TaskID, urlValues.Encode())

	var checklist ChecklistResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &checklist); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &checklist, nil
}

func (c *Client) UpdateChecklist(ctx context.Context, request *UpdateChecklistRequest) (*ChecklistResponse, error) {
	if request.ChecklistID == "" {
		return nil, fmt.Errorf("must provide checklist id: %w", ErrValidation)
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize checklist request: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("checklist/%s", request.ChecklistID)

	var checklist ChecklistResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &checklist); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &checklist, nil
}

func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error {
	if checklistID == "" {
		return fmt.Errorf("must provide checklistID: %w", ErrValidation)
	}

	endpoint := fmt.Sprintf("/checklist/%s", checklistID)

	return c.call(ctx, http.MethodDelete, endpoint, nil, &struct{}{})
}

func (c *Client) CreateChecklistItem(ctx context.Context, request *CreateChecklistItemRequest) (*ChecklistResponse, error) {
	if request.ChecklistID == "" {
		return nil, fmt.Errorf("must provide a checklist id: %w", ErrValidation)
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize checklist request: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/checklist/%s", request.ChecklistID)

	var checklist ChecklistResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &checklist); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &checklist, nil
}

func (c *Client) UpdateChecklistItem(ctx context.Context, request *UpdateChecklistItemRequest) (*ChecklistResponse, error) {
	if request.ChecklistID == "" {
		return nil, fmt.Errorf("must provide checklist id: %w", ErrValidation)
	}
	if request.ChecklistItemID == "" {
		return nil, fmt.Errorf("must provide checklist item id: %w", ErrValidation)
	}

	b, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize checklist item request: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/checklist/%s/checklist_item/%s", request.ChecklistID, request.ChecklistItemID)

	var checklist ChecklistResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &checklist); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &checklist, nil
}

func (c *Client) DeleteChecklistItem(ctx context.Context, checklistID, checklistItemID string) error {
	if checklistID == "" {
		return fmt.Errorf("must provide a checklist id: %w", ErrValidation)
	}
	if checklistItemID == "" {
		return fmt.Errorf("must provide a checklist item id: %w", ErrValidation)
	}

	endpoint := fmt.Sprintf("/checklist/%s/checklist_item/%s", checklistID, checklistItemID)

	return c.call(ctx, http.MethodDelete, endpoint, nil, &struct{}{})
}
