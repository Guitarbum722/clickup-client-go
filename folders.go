// Copyright (c) 2022, John Moore
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

type SingleFolder struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Orderindex       int    `json:"-"`
	OverrideStatuses bool   `json:"override_statuses"`
	Hidden           bool   `json:"hidden"`
	Space            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"space"`
	TaskCount string `json:"task_count"`
	Archived  bool   `json:"archived"`
	Lists     []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Orderindex int    `json:"-"`
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
		Space struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Access bool   `json:"access"`
		} `json:"space"`
		Archived         bool `json:"archived"`
		OverrideStatuses bool `json:"override_statuses"`
		Statuses         []struct {
			ID         string `json:"id"`
			Status     string `json:"status"`
			Orderindex int    `json:"-"`
			Color      string `json:"color"`
			Type       string `json:"type"`
		} `json:"statuses"`
		PermissionLevel string `json:"permission_level"`
	} `json:"lists"`
	PermissionLevel string `json:"permission_level"`
}

type FoldersResponse struct {
	Folders []SingleFolder `json:"folders"`
}

// FoldersForSpace queries all folders for a workspace.
func (c *Client) FoldersForSpace(ctx context.Context, workspaceID string, includeArchived bool) (*FoldersResponse, error) {

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	endpoint := fmt.Sprintf("/space/%s/folder/?%s", workspaceID, urlValues.Encode())

	var folders FoldersResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &folders); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &folders, nil
}

// FolderByID queries a single folder with folderID.
func (c *Client) FolderByID(ctx context.Context, folderID string) (*SingleFolder, error) {
	endpoint := fmt.Sprintf("/folder/%s", folderID)

	var folder SingleFolder

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &folder); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &folder, nil
}
