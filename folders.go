package clickup

import (
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

type GetFoldersResponse struct {
	Folders []SingleFolder `json:"folders"`
}

func (c *Client) GetFolders(spaceID string, includeArchived bool) (*GetFoldersResponse, error) {
	var folders GetFoldersResponse

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	uri := fmt.Sprintf("/space/%s/folder/?%s", spaceID, urlValues.Encode())

	if err := c.call(http.MethodGet, uri, nil, &folders); err != nil {
		return nil, err
	}

	return &folders, nil
}

func (c *Client) GetFolder(folderID string) (*SingleFolder, error) {
	var folder SingleFolder

	uri := fmt.Sprintf("/folder/%s", folderID)

	if err := c.call(http.MethodGet, uri, nil, &folder); err != nil {
		return nil, err
	}

	return &folder, nil
}
