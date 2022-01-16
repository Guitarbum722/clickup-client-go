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
)

type Team struct {
	ID      string       `json:"id"`
	Name    string       `json:"name"`
	Color   string       `json:"color"`
	Avatar  string       `json:"avatar"`
	Members []TeamMember `json:"members"`
}

type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

type TeamUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
	Initials       string `json:"initials"`
	Role           int    `json:"role"`
	CustomRole     string `json:"custom_role"`
	LastActive     string `json:"last_active"`
	DateJoined     string `json:"date_joined"`
	DateInvited    string `json:"date_invited"`
}
type TeamMember struct {
	User      TeamUser `json:"user"`
	InvitedBy TeamUser `json:"invited_by"`
}

// Teams returns a listing of teams for the authenticated user (the Client).
func (c *Client) Teams(ctx context.Context) (*TeamsResponse, error) {
	endpoint := fmt.Sprintf("/team")

	var teamsResponse TeamsResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &teamsResponse); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &teamsResponse, nil
}

// TeamsForWorkspace is exactly the same as Teams() and simply calls it.  This
// method is to maintain consistent naming with other Teams related parts
// of the API.  I would do it differently if I went back in time.
func (c *Client) TeamsForWorkspace(ctx context.Context) (*TeamsResponse, error) {
	return c.Teams(ctx)
}

type Group struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"team_id"`
	Userid      int        `json:"userid"`
	Name        string     `json:"name"`
	Handle      string     `json:"handle"`
	DateCreated string     `json:"date_created"`
	Initials    string     `json:"initials"`
	Members     []TeamUser `json:"members"`
}

type GroupsQueryResponse struct {
	Groups []Group `json:"groups"`
}

// GroupsForWorkspace queries for any groups in a workspace and returns their corresponding data.
// optionalGroupIDs can be provided to narrow the data returned to the explicit groups quried.
func (c *Client) GroupsForWorkspace(ctx context.Context, workspaceID string, optionalGroupIDs ...string) (*GroupsQueryResponse, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("must provide a workspaceID: %w", ErrValidation)
	}

	urlValues := url.Values{}
	urlValues.Set("team_id", workspaceID)
	for _, v := range optionalGroupIDs {
		urlValues.Add("group_ids%5B%5D", v)
	}

	endpoint := fmt.Sprintf("/group/?%s", urlValues.Encode())

	var groups GroupsQueryResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &groups); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &groups, nil
}

type CreateGroupRequest struct {
	WorkspaceID string
	Name        string `json:"name"`
	MemberIDs   []int  `json:"member_ids"`
}

type CreateGroupResponse struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"team_id"`
	Userid      int        `json:"userid"`
	Name        string     `json:"name"`
	Handle      string     `json:"handle"`
	DateCreated string     `json:"date_created"`
	Initials    string     `json:"initials"`
	Members     []TeamUser `json:"members"`
}

// CreateGroup adds a new group to a workspace using group.WorkspaceID.
func (c *Client) CreateGroup(ctx context.Context, group CreateGroupRequest) (*CreateGroupResponse, error) {
	if group.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace ID: %w", ErrValidation)
	}

	b, err := json.Marshal(group)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new group: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/team/%s/group", group.WorkspaceID)

	var newGroup CreateGroupResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newGroup); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &newGroup, nil
}

type UpdateGroupRequest struct {
	ID      string
	Name    string `json:"name,omitempty"`
	Handle  string `json:"handle,omitempty"`
	Members []struct {
		Add    []int `json:"add,omitempty"`
		Remove []int `json:"rem,omitempty"`
	} `json:"members,omitempty`
}

type UpdateGroupResponse struct {
	ID          string     `json:"id"`
	TeamID      string     `json:"team_id"`
	Userid      int        `json:"userid"`
	Name        string     `json:"name"`
	Handle      string     `json:"handle"`
	DateCreated string     `json:"date_created"`
	Initials    string     `json:"initials"`
	Members     []TeamUser `json:"members"`
}

// UpdateGroup changes an existing group using group.ID.
func (c *Client) UpdateGroup(ctx context.Context, group UpdateGroupRequest) (*UpdateGroupResponse, error) {
	if group.ID == "" {
		return nil, fmt.Errorf("must provide an ID: %w", ErrValidation)
	}

	b, err := json.Marshal(group)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new group: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/group/%s", group.ID)

	var updatedGroup UpdateGroupResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedGroup); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &updatedGroup, nil
}

// DeleteGroup removes an existing group with an id of groupID.
func (c *Client) DeleteGroup(ctx context.Context, groupID string) error {
	if groupID == "" {
		return fmt.Errorf("must provide a group id to delete: %w", ErrValidation)
	}
	return c.call(ctx, http.MethodDelete, fmt.Sprintf("/group/%s", groupID), nil, &struct{}{})
}
