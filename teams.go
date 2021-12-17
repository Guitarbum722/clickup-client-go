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

func (c *Client) Teams(ctx context.Context) (*TeamsResponse, error) {
	endpoint := fmt.Sprintf("%s/team", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("get teams request failed: %w", err)
	}
	if err := c.AuthenticateFor(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make teams request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var teamsResponse TeamsResponse

	if err := decoder.Decode(&teamsResponse); err != nil {
		return nil, fmt.Errorf("failed parse to teams response: %w", err)
	}

	return &teamsResponse, nil
}
