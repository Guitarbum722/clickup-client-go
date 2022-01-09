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

type CreateGoalRequest struct {
	WorkspaceID    string `json:"-"`
	Name           string `json:"name"`
	DueDate        int    `json:"due_date"`
	Description    string `json:"description"`
	MultipleOwners bool   `json:"multiple_owners"`
	Owners         []int  `json:"owners"`
	Color          string `json:"color"`
}

type CreateGoalResponse struct {
	Goal struct {
		ID               string      `json:"id"`
		Name             string      `json:"name"`
		TeamID           string      `json:"team_id"`
		DateCreated      string      `json:"date_created"`
		StartDate        string      `json:"start_date"`
		DueDate          string      `json:"due_date"`
		Description      string      `json:"description"`
		Private          bool        `json:"private"`
		Archived         bool        `json:"archived"`
		Creator          int         `json:"creator"`
		Color            string      `json:"color"`
		PrettyID         string      `json:"pretty_id"`
		MultipleOwners   bool        `json:"multiple_owners"`
		Members          []TeamUser  `json:"members"`
		Owners           []TeamUser  `json:"owners"`
		KeyResults       []KeyResult `json:"key_results"`
		PercentCompleted int         `json:"percent_completed"`
		PrettyURL        string      `json:"pretty_url"`
	} `json:"goal"`
}

func (c *Client) CreateGoal(ctx context.Context, goal CreateGoalRequest) (*CreateGoalResponse, error) {
	if goal.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id to create a goal: %w", ErrValidation)
	}

	b, err := json.Marshal(goal)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/team/%s/goal", goal.WorkspaceID)

	var newGoal CreateGoalResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newGoal); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &newGoal, nil
}

type UpdateGoalRequest struct {
	GoalID         string `json:"-"`
	Name           string `json:"name,omitempty"`
	DueDate        int    `json:"due_date,omitempty"`
	Description    string `json:"description,omitempty"`
	MultipleOwners bool   `json:"multiple_owners,omitempty"`
	Owners         []int  `json:"owners,omitempty"`
	Color          string `json:"color,omitempty"`
}

type UpdateGoalResponse struct {
	Goal struct {
		ID               string      `json:"id"`
		Name             string      `json:"name"`
		TeamID           string      `json:"team_id"`
		DateCreated      string      `json:"date_created"`
		StartDate        string      `json:"start_date"`
		DueDate          string      `json:"due_date"`
		Description      string      `json:"description"`
		Private          bool        `json:"private"`
		Archived         bool        `json:"archived"`
		Creator          int         `json:"creator"`
		Color            string      `json:"color"`
		PrettyID         string      `json:"pretty_id"`
		MultipleOwners   bool        `json:"multiple_owners"`
		Members          []TeamUser  `json:"members"`
		Owners           []TeamUser  `json:"owners"`
		KeyResults       []KeyResult `json:"key_results"`
		PercentCompleted int         `json:"percent_completed"`
		PrettyURL        string      `json:"pretty_url"`
	} `json:"goal"`
}

func (c *Client) UpdateGoal(ctx context.Context, goal UpdateGoalRequest) (*UpdateGoalResponse, error) {
	if goal.GoalID == "" {
		return nil, fmt.Errorf("must provide a goal id to update a goal: %w", ErrValidation)
	}

	b, err := json.Marshal(goal)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/goal/%s", goal.GoalID)

	var updatedGoal UpdateGoalResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedGoal); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &updatedGoal, nil
}

type GoalResponse struct {
	ID               string     `json:"id"`
	PrettyID         string     `json:"pretty_id"`
	Name             string     `json:"name"`
	TeamID           string     `json:"team_id"`
	Creator          int        `json:"creator"`
	Color            string     `json:"color"`
	DateCreated      string     `json:"date_created"`
	StartDate        string     `json:"start_date"`
	DueDate          string     `json:"due_date"`
	Description      string     `json:"description"`
	Private          bool       `json:"private"`
	Archived         bool       `json:"archived"`
	MultipleOwners   bool       `json:"multiple_owners"`
	EditorToken      string     `json:"editor_token"`
	DateUpdated      string     `json:"date_updated"`
	LastUpdate       string     `json:"last_update"`
	FolderID         string     `json:"folder_id"`
	Pinned           bool       `json:"pinned"`
	Owners           []TeamUser `json:"owners"`
	KeyResultCount   int        `json:"key_result_count"`
	PercentCompleted int        `json:"percent_completed"`
}

type GetGoalsResponse struct {
	Goals []GoalResponse `json:"goals"`
}

func (c *Client) GoalsForWorkspace(ctx context.Context, workspaceID string, includeCompleted bool) (*GetGoalsResponse, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id to retrieve goals for workspace: %w", ErrValidation)
	}
	urlValues := url.Values{}
	urlValues.Set("include_completed", strconv.FormatBool(includeCompleted))

	endpoint := fmt.Sprintf("/team/%s/goal/?%s", workspaceID, urlValues.Encode())

	var goals GetGoalsResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &goals); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &goals, nil
}

func (c *Client) GoalForWorkSpace(ctx context.Context, goalID string) (*GoalResponse, error) {
	if goalID == "" {
		return nil, fmt.Errorf("must provide a goal id to retrieve goal for workspace: %w", ErrValidation)
	}
	endpoint := fmt.Sprintf("/goal/%s", goalID)

	var goal GoalResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &goal); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &goal, nil
}

func (c *Client) DeleteGoal(ctx context.Context, goalID string) error {
	if goalID == "" {
		return fmt.Errorf("must provide a goal id to detelet: %w", ErrValidation)
	}
	return c.call(ctx, http.MethodDelete, fmt.Sprintf("/goal/%s", goalID), nil, &struct{}{})
}

type KeyResultType string

const (
	KeyResultNumber     KeyResultType = "number"
	KeyResultCurrency   KeyResultType = "currency"
	KeyResultBoolean    KeyResultType = "boolean"
	KeyResultPercentage KeyResultType = "percentage"
	KeyResultAutomatic  KeyResultType = "automatic"
)

type CreateKeyResultRequest struct {
	GoalID     string        `json:"-"`
	Name       string        `json:"name"`
	Owners     []int         `json:"owners"`
	Type       KeyResultType `json:"type"`
	StepsStart int           `json:"steps_start"`
	StepsEnd   int           `json:"steps_end"`
	Unit       string        `json:"unit"`
	TaskIds    []string      `json:"task_ids"`
	ListIds    []string      `json:"list_ids"`
}

type KeyResult struct {
	ID               string     `json:"id"`
	GoalID           string     `json:"goal_id"`
	Name             string     `json:"name"`
	Creator          int        `json:"creator"`
	Type             string     `json:"type"`
	DateCreated      string     `json:"date_created"`
	GoalPrettyID     string     `json:"goal_pretty_id"`
	PercentCompleted int        `json:"percent_completed"`
	Completed        bool       `json:"completed"`
	TaskIds          []string   `json:"task_ids"`
	Owners           []TeamUser `json:"owners"`
	LastAction       struct {
		ID           string `json:"id"`
		KeyResultID  string `json:"key_result_id"`
		Userid       int    `json:"userid"`
		Note         string `json:"note"`
		DateModified string `json:"date_modified"`
	} `json:"last_action"`
}

type CreateKeyResultResponse struct {
	KeyResult KeyResult `json:"key_result"`
}

func (c *Client) CreateKeyResultForGoal(ctx context.Context, keyResult CreateKeyResultRequest) (*CreateKeyResultResponse, error) {
	if keyResult.GoalID == "" {
		return nil, fmt.Errorf("must provide a goal id to create a key result: %w", ErrValidation)
	}

	b, err := json.Marshal(keyResult)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/goal/%s/key_result", keyResult.GoalID)

	var newKeyResult CreateKeyResultResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &newKeyResult); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &newKeyResult, nil
}

type UpdateKeyResultRequest struct {
	ID               string     `json:"-"`
	GoalID           string     `json:"goal_id,omitempty"`
	Name             string     `json:"name,omitempty"`
	Creator          int        `json:"creator,omitempty"`
	Type             string     `json:"type,omitempty"`
	DateCreated      string     `json:"date_created,omitempty"`
	GoalPrettyID     string     `json:"goal_pretty_id,omitempty"`
	PercentCompleted int        `json:"percent_completed,omitempty"`
	Completed        bool       `json:"completed,omitempty"`
	TaskIds          []string   `json:"task_ids,omitempty"`
	Owners           []TeamUser `json:"owners,omitempty"`
}

type UpdateKeyResultResponse struct {
	KeyResult KeyResult `json:"key_result"`
}

func (c *Client) UpdateKeyResult(ctx context.Context, keyResult UpdateKeyResultRequest) (*UpdateKeyResultResponse, error) {
	if keyResult.ID == "" {
		return nil, fmt.Errorf("must provide a key result id to update a goal: %w", ErrValidation)
	}

	b, err := json.Marshal(keyResult)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/key_result/%s", keyResult.ID)

	var updatedKeyResult UpdateKeyResultResponse

	if err := c.call(ctx, http.MethodPut, endpoint, buf, &updatedKeyResult); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &updatedKeyResult, nil
}

func (c *Client) DeleteKeyResult(ctx context.Context, keyResultID string) error {
	if keyResultID == "" {
		return fmt.Errorf("must provide key result id to delete: %w", ErrValidation)
	}
	return c.call(ctx, http.MethodDelete, fmt.Sprintf("/key_result/%s", keyResultID), nil, &struct{}{})
}
