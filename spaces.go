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
		Orderindex int    `json:"-"`
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
				Orderindex string `json:"-"`
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

	endpoint := fmt.Sprintf("/team/%s/space/?%s", teamID, urlValues.Encode())

	var spacesResponse SpacesResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &spacesResponse); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &spacesResponse, nil
}

func (c *Client) SpaceByID(ctx context.Context, spaceID string) (*SingleSpace, error) {

	endpoint := fmt.Sprintf("/space/%s", spaceID)

	var spaceResponse SingleSpace

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &spaceResponse); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &spaceResponse, nil
}

type TimeTracking struct {
	Enabled bool `json:"enabled"`
}

type DueDates struct {
	Enabled            bool `json:"enabled"`
	StartDate          bool `json:"start_date"`
	RemapDueDates      bool `json:"remap_due_dates"`
	RemapClosedDueDate bool `json:"remap_closed_due_date"`
}

type Tags struct {
	Enabled bool `json:"enabled,omitempty"`
}

type TimeEstimates struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Checklists struct {
	Enabled bool `json:"enabled,omitempty"`
}

type CustomFields struct {
	Enabled bool `json:"enabled,omitempty"`
}

type RemapDependencies struct {
	Enabled bool `json:"enabled,omitempty"`
}

type DependencyWarning struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Portfolios struct {
	Enabled bool `json:"enabled,omitempty"`
}

type Features struct {
	DueDates          *DueDates          `json:"due_dates,omitempty"`
	TimeTracking      *TimeTracking      `json:"time_tracking,omitempty"`
	Tags              *Tags              `json:"tags,omitempty"`
	TimeEstimates     *TimeEstimates     `json:"time_estimates,omitempty"`
	Checklists        *Checklists        `json:"checklists,omitempty"`
	CustomFields      *CustomFields      `json:"custom_fields,omitempty"`
	RemapDependencies *RemapDependencies `json:"remap_dependencies,omitempty"`
	DependencyWarning *DependencyWarning `json:"dependency_warning,omitempty"`
	Portfolios        *Portfolios        `json:"portfolios,omitempty"`
}

type CreateSpaceRequest struct {
	WorkspaceID       string
	Name              string    `json:"name"`
	MultipleAssignees bool      `json:"multiple_assignees"`
	Features          *Features `json:"features,omitempty"`
}

// CreateSpaceForWorkspace uses the parameters from space to create a new space in the specified workspace.
func (c *Client) CreateSpaceForWorkspace(ctx context.Context, space CreateSpaceRequest) (*SingleSpace, error) {
	if space.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id: %w", ErrValidation)
	}

	b, err := json.Marshal(space)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new space: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/team/%s/space", space.WorkspaceID)

	var newSpace SingleSpace

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newSpace); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &newSpace, nil
}

type UpdateSpaceRequest struct {
	ID                string
	Name              string    `json:"name,omitempty"`
	MultipleAssignees bool      `json:"multiple_assignees,omitempty"`
	Features          *Features `json:"features,omitempty"`
}

// UpdateSpaceForWorkspace makes changes to an existing space using parameters specified in space.
func (c *Client) UpdateSpaceForWorkspace(ctx context.Context, space UpdateSpaceRequest) (*SingleSpace, error) {
	if space.ID == "" {
		return nil, fmt.Errorf("must provide a workspace id: %w", ErrValidation)
	}

	b, err := json.Marshal(space)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new space: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/space/%s", space.ID)

	var updatedSpace SingleSpace

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedSpace); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &updatedSpace, nil
}

// DeleteSpace removes an existing space using spaceID.
func (c *Client) DeleteSpace(ctx context.Context, spaceID string) error {
	if spaceID == "" {
		return fmt.Errorf("must provide a space id to delete: %w", ErrValidation)
	}
	return c.call(ctx, http.MethodGet, fmt.Sprintf("/space/%s", spaceID), nil, &struct{}{})
}
