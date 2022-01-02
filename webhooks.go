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

func (c *Client) CreateWebhook(ctx context.Context, workspaceID string, webhook *CreateWebhookRequest) (*CreateWebhookResponse, error) {

	b, err := json.Marshal(webhook)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new webhook: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/team/%s/webhook", workspaceID)

	var newWebhook CreateWebhookResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newWebhook); err != nil {
		return nil, ErrCall
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

func (c *Client) UpdateWebhook(ctx context.Context, webhook *UpdateWebhookRequest) (*UpdateWebhookResponse, error) {
	if webhook.ID == "" {
		return nil, fmt.Errorf("must provide a webhook id: %w", ErrValidation)
	}

	b, err := json.Marshal(webhook)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize webhook: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/webhook/%s", webhook.ID)

	var updatedWebhook UpdateWebhookResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedWebhook); err != nil {
		return nil, ErrCall
	}

	return &updatedWebhook, nil
}

func (c *Client) DeleteWebhook(ctx context.Context, id string) error {
	return c.call(ctx, http.MethodGet, fmt.Sprintf("/webhook/%s", id), nil, &struct{}{})
}

func (c *Client) WebhooksFor(ctx context.Context, workspaceID string) (*WebhooksQueryResponse, error) {

	endpoint := fmt.Sprintf("/team/%s/webhook", workspaceID)

	var webhooks WebhooksQueryResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &webhooks); err != nil {
		return nil, ErrCall
	}

	return &webhooks, nil
}
