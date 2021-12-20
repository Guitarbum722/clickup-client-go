package clickup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SingleView struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Parent struct {
		ID   string `json:"id"`
		Type int    `json:"type"`
	} `json:"parent"`
	Grouping struct {
		Field     string   `json:"field"`
		Dir       int      `json:"dir"`
		Collapsed []string `json:"collapsed"`
		Ignore    bool     `json:"ignore"`
	} `json:"grouping"`
	Filters struct {
		Op     string `json:"op"`
		Fields []struct {
			Field string `json:"field"`
			Op    string `json:"op"`
			Idx   int    `json:"idx"`
		} `json:"fields"`
		Search             string `json:"search"`
		SearchCustomFields bool   `json:"search_custom_fields"`
		SearchDescription  bool   `json:"search_description"`
		SearchName         bool   `json:"search_name"`
		ShowClosed         bool   `json:"show_closed"`
	} `json:"filters"`
	Columns struct {
		Fields []struct {
			Field  string `json:"field"`
			Idx    int    `json:"idx"`
			Width  int    `json:"width"`
			Hidden bool   `json:"hidden"`
		} `json:"fields"`
	} `json:"columns"`
	TeamSidebar struct {
		AssignedComments bool `json:"assigned_comments"`
		UnassignedTasks  bool `json:"unassigned_tasks"`
	} `json:"team_sidebar"`
	Settings struct {
		ShowTaskLocations      bool `json:"show_task_locations"`
		ShowSubtasks           int  `json:"show_subtasks"`
		ShowSubtaskParentNames bool `json:"show_subtask_parent_names"`
		ShowClosedSubtasks     bool `json:"show_closed_subtasks"`
		ShowAssignees          bool `json:"show_assignees"`
		ShowImages             bool `json:"show_images"`
		ShowTimer              bool `json:"show_timer"`
		MeComments             bool `json:"me_comments"`
		MeSubtasks             bool `json:"me_subtasks"`
		MeChecklists           bool `json:"me_checklists"`
		ShowEmptyStatuses      bool `json:"show_empty_statuses"`
		AutoWrap               bool `json:"auto_wrap"`
		TimeInStatusView       int  `json:"time_in_status_view"`
	} `json:"settings"`
	DateCreated string `json:"date_created"`
	Creator     int    `json:"creator"`
	Visibility  string `json:"visibility"`
	Protected   bool   `json:"protected"`
	Orderindex  int    `json:"orderindex"`
}

type GetViewResponse struct {
	View SingleView `json:"view"`
}

type GetViewsResponse struct {
	Views []SingleView `json:"views"`
}

type TasksForViewResponse struct {
	Tasks    []SingleTask `json:"tasks"`
	LastPage bool         `json:"last_page"`
}

func (c *Client) ViewByID(ctx context.Context, viewID string) (*GetViewResponse, error) {
	endpoint := fmt.Sprintf("%s/view/%s", c.baseURL, viewID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("view by id request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make view by id request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var view GetViewResponse

	if err := decoder.Decode(&view); err != nil {
		return nil, fmt.Errorf("failed to parse view: %w", err)

	}

	return &view, nil
}

type ViewListType int

const (
	TypeTeam ViewListType = iota
	TypeSpace
	TypeFolder
	TypeList
)

func (v ViewListType) String() string {
	switch v {
	case TypeTeam:
		return "team"
	case TypeSpace:
		return "space"
	case TypeFolder:
		return "folder"
	case TypeList:
		return "list"
	default:
		return "UNKNOWN_VIEW_LIST_TYPE"
	}
}
func (c *Client) ViewsFor(ctx context.Context, viewListType ViewListType, id string) (*GetViewsResponse, error) {
	viewsType := viewListType.String()
	if viewsType == "UNKNOWN_VIEW_LIST_TYPE" {
		return nil, errors.New("invalid ViewListType")
	}

	endpoint := fmt.Sprintf("%s/%s/%s/view", c.baseURL, viewsType, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("view request request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make views request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var views GetViewsResponse

	if err := decoder.Decode(&views); err != nil {
		return nil, fmt.Errorf("failed to parse views: %w", err)
	}

	return &views, nil
}

func (c *Client) DeleteView(ctx context.Context, viewID string) error {
	endpoint := fmt.Sprintf("%s/view/%s", c.baseURL, viewID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("delete view request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make delete view request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
}

// TasksForView requires possible pagination.  Clickup documents that a page will have a
// maximum of 30 tasks per page, defaulting to page 0.  This endpoint returns a boolean
// specifying whether or not the response consists of the last page (TasksForViewResponse.LastPage = true/false).
// func (c *Client) TasksForView(ctx context.Context, viewID string, page int) (*TasksForViewResponse, error) {
// 	panic("TODO: not implemented")
// }
