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
	endpoint := fmt.Sprintf("%s/task/%s/checklist/?%s", c.baseURL, request.TaskID, urlValues.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create checklist request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create checklist request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var checklist ChecklistResponse

	if err := decoder.Decode(&checklist); err != nil {
		return nil, fmt.Errorf("failed to parse new checklist: %w", err)
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

	endpoint := fmt.Sprintf("%s/checklist/%s", c.baseURL, request.ChecklistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("update checklist request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make update checklist request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var checklist ChecklistResponse

	if err := decoder.Decode(&checklist); err != nil {
		return nil, fmt.Errorf("failed to parse updated checklist: %w", err)
	}

	return &checklist, nil
}

func (c *Client) DeleteChecklist(ctx context.Context, checklistID string) error {
	if checklistID == "" {
		return fmt.Errorf("must provide checklistID: %w", ErrValidation)
	}

	endpoint := fmt.Sprintf("%s/checklist/%s", c.baseURL, checklistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete checklist request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete checklist request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
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

	endpoint := fmt.Sprintf("%s/checklist/%s", c.baseURL, request.ChecklistID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("create checklist item request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create checklist item request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var checklist ChecklistResponse

	if err := decoder.Decode(&checklist); err != nil {
		return nil, fmt.Errorf("failed to parse new checklist: %w", err)
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

	endpoint := fmt.Sprintf("%s/checklist/%s/checklist_item/%s", c.baseURL, request.ChecklistID, request.ChecklistItemID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, buf)
	if err != nil {
		return nil, fmt.Errorf("update checklist item request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-type", "application/json")

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make update checklist request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var checklist ChecklistResponse

	if err := decoder.Decode(&checklist); err != nil {
		return nil, fmt.Errorf("failed to parse updated checklist: %w", err)
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

	endpoint := fmt.Sprintf("%s/checklist/%s/checklist_item/%s", c.baseURL, checklistID, checklistItemID)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete checklist item request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete checklist item request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
}
