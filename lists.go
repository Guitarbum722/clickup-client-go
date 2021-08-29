package clickup

import (
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

type GetListsResponse struct {
	Lists []SingleList `json:"lists"`
}

func (c *Client) GetLists(folderID string, includeArchived bool) (*GetListsResponse, error) {
	var lists GetListsResponse

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	uri := fmt.Sprintf("/folder/%s/list/?%s", folderID, urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &lists); err != nil {
		return nil, err
	}

	return &lists, nil
}

func (c *Client) GetList(listID string) (*SingleList, error) {
	var list SingleList

	uri := fmt.Sprintf("/list/%s", listID)

	if err := c.call(http.MethodGet, uri, nil, &list); err != nil {
		return nil, err
	}

	return &list, nil
}
