// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SingleList struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Orderindex int    `json:"orderindex"`
	Status     struct {
	} `json:"status"`
	Priority struct {
	} `json:"priority"`
	Assignee struct {
	} `json:"assignee"`
	TaskCount int `json:"task_count"`
	DueDate   struct {
	} `json:"due_date"`
	StartDate struct {
	} `json:"start_date"`
	Folder struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"folder"`
	Space struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Access bool   `json:"access"`
	} `json:"space"`
	Archived         bool   `json:"archived"`
	OverrideStatuses bool   `json:"override_statuses"`
	PermissionLevel  string `json:"permission_level"`
}

type ListsResponse struct {
	Lists []SingleList `json:"lists"`
}

func (c *Client) ListsForFolder(ctx context.Context, folderID string, includeArchived bool) (*ListsResponse, error) {

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	endpoint := fmt.Sprintf("/folder/%s/list/?%s", folderID, urlValues.Encode())

	var lists ListsResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &lists); err != nil {
		return nil, ErrCall
	}

	return &lists, nil
}

func (c *Client) ListByID(ctx context.Context, listID string) (*SingleList, error) {

	endpoint := fmt.Sprintf("/list/%s", listID)

	var list SingleList

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &list); err != nil {
		return nil, ErrCall
	}

	return &list, nil
}
