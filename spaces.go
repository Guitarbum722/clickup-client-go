package clickup

import (
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

type GetSpacesResponse struct {
	Spaces []SingleSpace `json:"spaces"`
}

func (c *Client) GetSpaces(teamID string, includeArchived bool) (*GetSpacesResponse, error) {
	var spacesResponse GetSpacesResponse

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	uri := fmt.Sprintf("/team/%s/space/?%s", teamID, urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &spacesResponse); err != nil {
		return nil, err
	}

	return &spacesResponse, nil
}

func (c *Client) GetSpace(spaceID string) (*SingleSpace, error) {
	var spaceResponse SingleSpace

	uri := fmt.Sprintf("/space/%s", spaceID)

	if err := c.call(http.MethodGet, uri, nil, &spaceResponse); err != nil {
		return nil, err
	}

	return &spaceResponse, nil
}
