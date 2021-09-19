package clickup

import (
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
	TaskID   string         `json:"task_id"`
	ListID   string         `json:"list_id"`
	FolderID string         `json:"folder_id"`
	SpaceID  string         `json:"space_id"`
	Health   *WebhookHealth `json:"health"`
	Secret   string         `json:"secret"`
}

type WebhooksQueryResponse struct {
	Webhooks []Webhook `json:"webhooks"`
}

type CreateWebhookResponse struct {
	ID      string   `json:"id"`
	Webhook *Webhook `json:"webhook"`
}

type UpdateWebhookResponse struct {
	CreateWebhookResponse
}

func (c *Client) CreateWebhook(workspaceID string) (*CreateWebhookResponse, error) {
	panic("not implemented!")
	endpoint := fmt.Sprintf("%s/team/%s/webhook", c.baseURL, workspaceID)
}

func (c *Client) UpdateWebhook(id string) (*UpdateWebhookResponse, error) {
	panic("not implemented!")
	endpoint := fmt.Sprintf("%s/webhook/%s", c.baseURL, id)

}

func (c *Client) DeleteWebhook(id string) error {
	endpoint := fmt.Sprintf("%s/webhook/%s", c.baseURL, id)

	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("webhooks delete request failed: %w", err)
	}
	c.AuthenticateFor(req)

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
	c.AuthenticateFor(req)

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
