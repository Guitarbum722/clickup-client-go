// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WebhookEventMessage struct {
	Event        WebhookEvent `json:"event"`
	HistoryItems []struct {
		ID       string `json:"id"`
		Type     int    `json:"type"`
		Date     string `json:"date"`
		Field    string `json:"field"`
		ParentID string `json:"parent_id"`
		Data     struct {
			StatusType string `json:"status_type"`
		} `json:"data"`
		User struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Email          string `json:"email"`
			Color          string `json:"color"`
			Initials       string `json:"initials"`
			ProfilePicture string `json:"profilePicture"`
		} `json:"user"`
		Before struct {
			Status     string `json:"status"`
			Color      string `json:"color"`
			Orderindex int    `json:"orderindex"`
			Type       string `json:"type"`
		} `json:"before"`
		After struct {
			Status     string `json:"status"`
			Color      string `json:"color"`
			Orderindex int    `json:"orderindex"`
			Type       string `json:"type"`
		} `json:"after"`
	} `json:"history_items"`
	TaskID    string `json:"task_id"`
	WebhookID string `json:"webhook_id"`
}

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

type webhookVerifyResult struct {
	validSignature       bool
	signatureFromClickup string
	signatureGenerated   string
}

// VerifyWebhookSignature compares a generated signature using secret that
// is returned with the webhook CRUD operations with the x-signature http header
// that is sent with the http request to the webhook endpoint.
// It should be noted that err will be nil even if the signature is not valid,
// thus the WebhookVerifyResult.ValidSignature() should be called.
func VerifyWebhookSignature(webhookRequest *http.Request, secret string) (*webhookVerifyResult, error) {
	h := hmac.New(sha256.New, []byte(secret))
	b, err := ioutil.ReadAll(webhookRequest.Body)
	if err != nil {
		return nil, err
	}
	h.Write(b)
	sha := hex.EncodeToString(h.Sum(nil))

	sigHeader := webhookRequest.Header.Get("X-Signature")

	return &webhookVerifyResult{
		validSignature:       sigHeader == sha,
		signatureFromClickup: sigHeader,
		signatureGenerated:   sha,
	}, nil
}

func (w *webhookVerifyResult) Valid() bool {
	return w.validSignature
}

func (w *webhookVerifyResult) SignatureFromClickup() string {
	return w.signatureFromClickup
}

func (w *webhookVerifyResult) SignatureGenerated() string {
	return w.signatureGenerated
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
