package clickup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookEvent string

const (
	EventAll                     WebhookEvent = "*"
	EventTaskCreated             WebhookEvent = "taskCreated"
	EventTaskUpdated             WebhookEvent = "taskUpdated"
	EventTaskDeleted             WebhookEvent = "taskDeleted"
	EventTaskPriorityUpdated     WebhookEvent = "taskPriorityUpdated"
	EventTaskStatusUpdated       WebhookEvent = "taskStatusUpdated"
	EventTaskAssigneeUpdated     WebhookEvent = "taskAssigneeUpdated"
	EventTaskDueDateUpdated      WebhookEvent = "taskDueDateUpdated"
	EventTaskTagUpdated          WebhookEvent = "taskTagUpdated"
	EventTaskMoved               WebhookEvent = "taskMoved"
	EventTaskCommentPosted       WebhookEvent = "taskCommentPosted"
	EventTaskCommentUpdated      WebhookEvent = "taskCommentUpdated"
	EventTaskTimeEstimateUpdated WebhookEvent = "taskTimeEstimateUpdated"
	EventTaskTimeTrackedUpdated  WebhookEvent = "taskTimeTrackedUpdated"
	EventListCreated             WebhookEvent = "listCreated"
	EventListUpdated             WebhookEvent = "listUpdated"
	EventListDeleted             WebhookEvent = "listDeleted"
	EventFolderCreated           WebhookEvent = "folderCreated"
	EventFolderUpdated           WebhookEvent = "folderUpdated"
	EventFolderDeleted           WebhookEvent = "folderDeleted"
	EventSpaceCreated            WebhookEvent = "spaceCreated"
	EventSpaceUpdated            WebhookEvent = "spaceUpdated"
	EventSpaceDeleted            WebhookEvent = "spaceDeleted"
	EventGoalCreated             WebhookEvent = "goalCreated"
	EventGoalUpdated             WebhookEvent = "goalUpdated"
	EventGoalDeleted             WebhookEvent = "goalDeleted"
	EventKeyResultCreated        WebhookEvent = "keyResultCreated"
	EventKeyResultUpdated        WebhookEvent = "keyResultUpdated"
	EventKeyResultDeleted        WebhookEvent = "keyResultDeleted"
)

type WebhookHealth struct {
	Status    string `json:"status"`
	FailCount int    `json:"fail_count"`
}

type Webhook struct {
	ID       string         `json:"id"`
	UserID   int            `json:"userid"`
	TeamID   int            `json:"team_id"`
	Endpoint string         `json:"endpoint"`
	ClientID string         `json:"client_id"`
	Events   []WebhookEvent `json:"events"`
	TaskID   int            `json:"task_id"`
	ListID   int            `json:"list_id"`
	FolderID int            `json:"folder_id"`
	SpaceID  int            `json:"space_id"`
	Health   *WebhookHealth `json:"health"`
	Secret   string         `json:"secret"`
}

type WebhooksQueryResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

type CreateWebhookResponse struct {
	ID      string `json:"id"`
	Webhook struct {
		ID       string         `json:"id"`
		UserID   int            `json:"userid"`
		TeamID   int            `json:"team_id"`
		Endpoint string         `json:"endpoint"`
		ClientID string         `json:"client_id"`
		Events   []WebhookEvent `json:"events"`
		TaskID   int            `json:"task_id"`
		ListID   int            `json:"list_id"`
		FolderID int            `json:"folder_id"`
		SpaceID  int            `json:"space_id"`
		Health   *WebhookHealth `json:"health"`
		Secret   string         `json:"secret"`
	} `json:"webhook"`
}

type UpdateWebhookResponse struct {
	CreateWebhookResponse
}

type CreateWebhookRequest struct {
	Endpoint string         `json:"endpoint,omitempty"`
	Events   []WebhookEvent `json:"events,omitempty"`
	TaskID   string         `json:"task_id,omitempty"`
	ListID   string         `json:"list_id,omitempty"`
	FolderID string         `json:"folder_id,omitempty"`
}

func (c *Client) CreateWebhook(workspaceID string, webhook *CreateWebhookRequest) (*CreateWebhookResponse, error) {

	b, err := json.Marshal(webhook)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new webhook: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("%s/team/%s/webhook", c.baseURL, workspaceID)

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
		return nil, fmt.Errorf("failed to make create webhook request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var newWebhook CreateWebhookResponse

	if err := decoder.Decode(&newWebhook); err != nil {
		return nil, fmt.Errorf("failed to parse new webhook - WARNING! WEBHOOK MIGHT HAVE BEEN CREATED: %w", err)
	}

	return &newWebhook, nil
}

type UpdateWebhookRequest struct {
	ID       string         `json:"id"`
	Endpoint string         `json:"endpoint,omitempty"`
	Events   []WebhookEvent `json:"events,omitempty"`
	TaskID   string         `json:"task_id,omitempty"`
	ListID   string         `json:"list_id,omitempty"`
	FolderID string         `json:"folder_id,omitempty"`
	Status   string         `json:"status,omitempty"`
}

func (c *Client) UpdateWebhook(webhook *UpdateWebhookRequest) (*UpdateWebhookResponse, error) {
	if webhook.ID == "" {
		return nil, fmt.Errorf("must provide a webhook id: %w", ErrValidation)
	}

	b, err := json.Marshal(webhook)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize webhook: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("%s/webhook/%s", c.baseURL, webhook.ID)

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
		return nil, fmt.Errorf("failed to make update webhook request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var updatedWebhook UpdateWebhookResponse

	if err := decoder.Decode(&updatedWebhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}

	return &updatedWebhook, nil
}

func (c *Client) DeleteWebhook(id string) error {
	endpoint := fmt.Sprintf("%s/webhook/%s", c.baseURL, id)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("webhooks delete request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete webhooks request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
}

func (c *Client) WebhooksFor(workspaceID string) (*WebhooksQueryResponse, error) {

	endpoint := fmt.Sprintf("%s/team/%s/webhook", c.baseURL, workspaceID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("webhooks request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make get webhooks request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var webhooks WebhooksQueryResponse

	if err := decoder.Decode(&webhooks); err != nil {
		return nil, fmt.Errorf("failed to parse webhooks: %w", err)
	}

	return &webhooks, nil
}
