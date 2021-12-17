// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"context"
	"encoding/json"
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

	endpoint := fmt.Sprintf("%s/folder/%s/list/?%s", c.baseURL, folderID, urlValues.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("lists for folder request failed: %w", err)
	}
	if err := c.AuthenticateFor(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make lists request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var lists ListsResponse

	if err := decoder.Decode(&lists); err != nil {
		return nil, fmt.Errorf("failed to parse lists: %w", err)
	}

	return &lists, nil
}

func (c *Client) ListByID(ctx context.Context, listID string) (*SingleList, error) {

	endpoint := fmt.Sprintf("%s/list/%s", c.baseURL, listID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("list request failed: %w", err)
	}
	if err := c.AuthenticateFor(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make list request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var list SingleList

	if err := decoder.Decode(&list); err != nil {
		return nil, fmt.Errorf("failed parse to list: %w", err)
	}

	return &list, nil
}
