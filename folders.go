// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SingleFolder struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Orderindex       int    `json:"orderindex"`
	OverrideStatuses bool   `json:"override_statuses"`
	Hidden           bool   `json:"hidden"`
	Space            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"space"`
	TaskCount string        `json:"task_count"`
	Archived  bool          `json:"archived"`
	Statuses  []interface{} `json:"statuses"`
	Lists     []struct {
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
			Orderindex int    `json:"orderindex"`
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

func (c *Client) FoldersForSpace(spaceID string, includeArchived bool) (*FoldersResponse, error) {

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	endpoint := fmt.Sprintf("%s/space/%s/folder/?%s", c.baseURL, spaceID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("folder by space request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make folders by space request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var folders FoldersResponse

	if err := decoder.Decode(&folders); err != nil {
		return nil, fmt.Errorf("failed to parse folders: %w", err)
	}

	return &folders, nil
}

func (c *Client) FolderByID(folderID string) (*SingleFolder, error) {
	endpoint := fmt.Sprintf("%s/folder/%s", c.baseURL, folderID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("folder request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make folder request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var folder SingleFolder

	if err := decoder.Decode(&folder); err != nil {
		return nil, fmt.Errorf("failed to parse folder: %w", err)
	}

	return &folder, nil
}
