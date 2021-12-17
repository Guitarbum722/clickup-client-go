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

type SingleSpace struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Private  bool   `json:"private"`
	Statuses []struct {
		ID         string `json:"id"`
		Status     string `json:"status"`
		Type       string `json:"type"`
		Orderindex int    `json:"orderindex"`
		Color      string `json:"color"`
	} `json:"statuses"`
	MultipleAssignees bool `json:"multiple_assignees"`
	Features          struct {
		DueDates struct {
			Enabled            bool `json:"enabled"`
			StartDate          bool `json:"start_date"`
			RemapDueDates      bool `json:"remap_due_dates"`
			RemapClosedDueDate bool `json:"remap_closed_due_date"`
		} `json:"due_dates"`
		Sprints struct {
			Enabled bool `json:"enabled"`
		} `json:"sprints"`
		Points struct {
			Enabled bool `json:"enabled"`
		} `json:"points"`
		CustomItems struct {
			Enabled bool `json:"enabled"`
		} `json:"custom_items"`
		Priorities struct {
			Enabled    bool `json:"enabled"`
			Priorities []struct {
				ID         string `json:"id"`
				Priority   string `json:"priority"`
				Color      string `json:"color"`
				Orderindex string `json:"orderindex"`
			} `json:"priorities"`
		} `json:"priorities"`
		Tags struct {
			Enabled bool `json:"enabled"`
		} `json:"tags"`
		CheckUnresolved struct {
			Enabled  bool `json:"enabled"`
			Subtasks struct {
			} `json:"subtasks"`
			Checklists struct {
			} `json:"checklists"`
			Comments struct {
			} `json:"comments"`
		} `json:"check_unresolved"`
		Zoom struct {
			Enabled bool `json:"enabled"`
		} `json:"zoom"`
		Milestones struct {
			Enabled bool `json:"enabled"`
		} `json:"milestones"`
		CustomFields struct {
			Enabled bool `json:"enabled"`
		} `json:"custom_fields"`
		DependencyWarning struct {
			Enabled bool `json:"enabled"`
		} `json:"dependency_warning"`
	} `json:"features"`
	Archived bool `json:"archived"`
	Members  []struct {
		User struct {
			ID             int    `json:"id"`
			Username       string `json:"username"`
			Color          string `json:"color"`
			ProfilePicture string `json:"profilePicture"`
			Initials       string `json:"initials"`
		} `json:"user"`
	} `json:"members"`
}

type SpacesResponse struct {
	Spaces []SingleSpace `json:"spaces"`
}

func (c *Client) SpacesForWorkspace(ctx context.Context, teamID string, includeArchived bool) (*SpacesResponse, error) {
	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	endpoint := fmt.Sprintf("%s/team/%s/space/?%s", c.baseURL, teamID, urlValues.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("spaces request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make spaces request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var spacesResponse SpacesResponse

	if err := decoder.Decode(&spacesResponse); err != nil {
		return nil, fmt.Errorf("failed parse to spaces: %w", err)
	}

	return &spacesResponse, nil
}

func (c *Client) SpaceByID(ctx context.Context, spaceID string) (*SingleSpace, error) {

	endpoint := fmt.Sprintf("%s/space/%s", c.baseURL, spaceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("space request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make space request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var spaceResponse SingleSpace

	if err := decoder.Decode(&spaceResponse); err != nil {
		return nil, fmt.Errorf("failed parse to space: %w", err)
	}

	return &spaceResponse, nil
}
